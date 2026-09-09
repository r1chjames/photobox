#!/usr/bin/env python3
"""face-engine — face detection (YuNet) + recognition embeddings (SFace) over HTTP.

Stdlib-only HTTP server (ThreadingHTTPServer); inference via onnxruntime.

Endpoints:
  GET  /health  -> {"status": "ok"}
  POST /detect  body {"image": "<base64 image bytes>", "min_score": 0.5}
               -> {"faces": [{"box": [x, y, w, h], "score": float,
                              "landmarks": [[x, y] x 5], "embedding": [128 floats]}]}

Boxes and landmarks are in ORIGINAL image pixel coordinates (unscaled).
Models are downloaded to ./models on first run if missing.
"""
from __future__ import annotations

import base64
import io
import json
import math
import os
import sys
import threading
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

import numpy as np
import onnxruntime as ort
from PIL import Image

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
# The Helm chart sets FACE_MODELS_DIR; fall back to ./models for local runs.
MODELS_DIR = os.environ.get("FACE_MODELS_DIR") or os.path.join(BASE_DIR, "models")

YUNET_PATH = os.path.join(MODELS_DIR, "face_detection_yunet_2023mar.onnx")
SFACE_PATH = os.path.join(MODELS_DIR, "face_recognition_sface_2021dec.onnx")
YUNET_URL = ("https://github.com/opencv/opencv_zoo/raw/main/models/"
             "face_detection_yunet/face_detection_yunet_2023mar.onnx")
SFACE_URL = ("https://github.com/opencv/opencv_zoo/raw/main/models/"
             "face_recognition_sface/face_recognition_sface_2021dec.onnx")

MAX_SIDE = 640          # downscale longest side to this before inference
STRIDES = (8, 16, 32)   # YuNet FPN strides
MIN_FACE_PX = 20        # drop boxes smaller than this (quality filter)
NMS_IOU = 0.3           # NMS IoU threshold
EMBED_SIZE = 112        # SFace input size
EMBED_MARGIN = 1.3      # square crop margin around detection box

_log_lock = threading.Lock()


def log(msg: str) -> None:
    with _log_lock:
        print(f"[face-engine] {msg}", file=sys.stderr, flush=True)


def download(url: str, dest: str) -> None:
    """Download a model file atomically (tmp + rename), with retries."""
    tmp = dest + ".part"
    last_exc: Exception | None = None
    for attempt in range(1, 4):
        try:
            log(f"downloading model (attempt {attempt}): {url}")
            with urllib.request.urlopen(url, timeout=60) as resp, open(tmp, "wb") as f:
                while True:
                    chunk = resp.read(1 << 20)
                    if not chunk:
                        break
                    f.write(chunk)
            os.replace(tmp, dest)
            return
        except Exception as exc:  # noqa: BLE001 - retry any network failure
            last_exc = exc
            log(f"download failed: {exc}")
    raise RuntimeError(f"could not download model from {url}: {last_exc}")


def ensure_models() -> None:
    os.makedirs(MODELS_DIR, exist_ok=True)
    for url, path in ((YUNET_URL, YUNET_PATH), (SFACE_URL, SFACE_PATH)):
        if not os.path.exists(path):
            download(url, path)


def _iou(a, b) -> float:
    ax1, ay1, aw, ah = a
    bx1, by1, bw, bh = b
    ix1, iy1 = max(ax1, bx1), max(ay1, by1)
    ix2, iy2 = min(ax1 + aw, bx1 + bw), min(ay1 + ah, by1 + bh)
    inter = max(0.0, ix2 - ix1) * max(0.0, iy2 - iy1)
    union = aw * ah + bw * bh - inter
    return inter / union if union > 0 else 0.0


class FaceEngine:
    """YuNet detection + SFace embeddings on CPU (onnxruntime)."""

    def __init__(self, models_dir: str = MODELS_DIR) -> None:
        self.models_dir = models_dir
        yunet = os.path.join(models_dir, "face_detection_yunet_2023mar.onnx")
        sface = os.path.join(models_dir, "face_recognition_sface_2021dec.onnx")
        opts = ort.SessionOptions()
        opts.intra_op_num_threads = max(1, (os.cpu_count() or 2) // 2)
        self._det = ort.InferenceSession(yunet, opts, providers=["CPUExecutionProvider"])
        self._rec = ort.InferenceSession(sface, opts, providers=["CPUExecutionProvider"])
        self._det_in = self._det.get_inputs()[0].name
        self._rec_in = self._rec.get_inputs()[0].name

    # -- image prep ---------------------------------------------------------

    @staticmethod
    def _to_bgr(pil_img: Image.Image) -> np.ndarray:
        return np.array(pil_img.convert("RGB"))[:, :, ::-1].copy()

    def _prep(self, img_bgr: np.ndarray):
        """Downscale longest side to MAX_SIDE (keep aspect), pad bottom/right to 640x640."""
        h, w = img_bgr.shape[:2]
        scale = 1.0
        if max(h, w) > MAX_SIDE:
            scale = MAX_SIDE / max(h, w)
            nw, nh = int(w * scale), int(h * scale)
            pil = Image.fromarray(img_bgr[:, :, ::-1]).resize((nw, nh), Image.BILINEAR)
            img_bgr = np.array(pil)[:, :, ::-1].copy()
            h, w = nh, nw
        pad = np.zeros((MAX_SIDE, MAX_SIDE, 3), dtype=np.uint8)
        pad[:h, :w] = img_bgr
        return pad, scale

    # -- YuNet decode ---------------------------------------------------------

    def detect(self, img_bgr: np.ndarray, min_score: float = 0.5):
        """Return [(box [x,y,w,h], score, landmarks [[x,y]x5])] in original coords."""
        orig_h, orig_w = img_bgr.shape[:2]
        padded, scale = self._prep(img_bgr)
        ph, pw = padded.shape[:2]
        blob = padded.transpose(2, 0, 1)[None].astype(np.float32)  # raw 0..255 BGR
        outs = self._det.run(None, {self._det_in: blob})
        cands = []
        for si, st in enumerate(STRIDES):
            rows, cols = ph // st, pw // st
            cls = outs[si][0, :, 0]
            obj = outs[si + 3][0, :, 0]
            score = np.sqrt(np.clip(cls, 0, 1) * np.clip(obj, 0, 1))
            bb = outs[si + 6][0]
            kp = outs[si + 9][0]
            for j in np.where(score >= min_score)[0]:
                r, c = j // cols, j % cols
                v = bb[j]
                cx = (c + v[0]) * st
                cy = (r + v[1]) * st
                bw = math.exp(v[2]) * st
                bh = math.exp(v[3]) * st
                x1, y1 = cx - bw / 2, cy - bh / 2
                # unscale back to original image coordinates
                x1 /= scale; y1 /= scale; bw /= scale; bh /= scale
                x1 = max(0.0, min(x1, float(orig_w)))
                y1 = max(0.0, min(y1, float(orig_h)))
                bw = min(bw, orig_w - x1)
                bh = min(bh, orig_h - y1)
                if bw < MIN_FACE_PX or bh < MIN_FACE_PX:  # quality filter tiny faces
                    continue
                lm = [[(kp[j][2 * n] + c) * st / scale,
                       (kp[j][2 * n + 1] + r) * st / scale] for n in range(5)]
                cands.append(([float(x1), float(y1), float(bw), float(bh)],
                              float(score[j]), lm))
        if not cands:
            return []
        cands.sort(key=lambda x: -x[1])
        keep = []
        for cand in cands:
            if all(_iou(cand[0], k[0]) <= NMS_IOU for k in keep):
                keep.append(cand)
        return keep

    # -- SFace embedding ------------------------------------------------------

    def embed(self, img_bgr: np.ndarray, box):
        """128-d L2-normalized embedding for the face at `box` (original coords)."""
        H, W = img_bgr.shape[:2]
        x, y, w, h = box
        cx, cy = x + w / 2, y + h / 2
        s = max(w, h) * EMBED_MARGIN
        x0, y0 = int(max(0, cx - s / 2)), int(max(0, cy - s / 2))
        x1, y1 = int(min(W, cx + s / 2)), int(min(H, cy + s / 2))
        crop = img_bgr[y0:y1, x0:x1]
        if crop.size == 0:
            return None
        pil = Image.fromarray(crop[:, :, ::-1])  # back to RGB for PIL
        pw, ph = pil.size
        side = max(pw, ph)
        canvas = Image.new("RGB", (side, side), (0, 0, 0))  # pad black
        canvas.paste(pil, ((side - pw) // 2, (side - ph) // 2))
        face = np.array(canvas.resize((EMBED_SIZE, EMBED_SIZE), Image.BILINEAR))
        inp = (face[None].astype(np.float32) / 127.5 - 1.0).transpose(0, 3, 1, 2)
        e = self._rec.run(None, {self._rec_in: inp})[0][0]
        n = np.linalg.norm(e)
        return (e / n).tolist() if n > 0 else None

    # -- combined -------------------------------------------------------------

    def detect_image(self, pil_img: Image.Image, min_score: float = 0.5):
        """Detect + embed; returns list of face dicts ready for JSON."""
        img_bgr = self._to_bgr(pil_img)
        out = []
        for box, score, lm in self.detect(img_bgr, min_score):
            emb = self.embed(img_bgr, box)
            if emb is None:
                continue
            out.append({
                "box": [round(v, 2) for v in box],
                "score": round(score, 4),
                "landmarks": [[round(a, 2), round(b, 2)] for a, b in lm],
                "embedding": [round(float(v), 6) for v in emb],
            })
        return out


class Handler(BaseHTTPRequestHandler):
    server_version = "face-engine/1.0"

    def log_message(self, fmt, *args):  # route access logs through log()
        log(f"{self.address_string()} {fmt % args}")

    def _send_json(self, code: int, obj) -> None:
        body = json.dumps(obj).encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):  # noqa: N802 (http.server API)
        if self.path == "/health":
            self._send_json(200, {"status": "ok"})
        else:
            self._send_json(404, {"error": "not found"})

    def do_POST(self):  # noqa: N802 (http.server API)
        if self.path != "/detect":
            return self._send_json(404, {"error": "not found"})
        try:
            length = int(self.headers.get("Content-Length") or 0)
            raw = self.rfile.read(length) if length > 0 else b""
            body = json.loads(raw.decode("utf-8"))
            image_b64 = body["image"]
            min_score = float(body.get("min_score", 0.5))
        except (ValueError, KeyError, UnicodeDecodeError) as exc:
            return self._send_json(400, {"error": f"bad request: {exc}"})
        try:
            img_bytes = base64.b64decode(image_b64, validate=True)
            pil_img = Image.open(io.BytesIO(img_bytes))
            pil_img.load()
        except Exception as exc:  # noqa: BLE001 - any decode failure is a 400
            return self._send_json(400, {"error": f"could not decode image: {exc}"})
        try:
            faces = self.server.engine.detect_image(pil_img, min_score)
        except Exception as exc:  # noqa: BLE001
            log(f"detect failed: {exc}")
            return self._send_json(500, {"error": str(exc)})
        self._send_json(200, {"faces": faces})


def create_server(engine: FaceEngine, host: str = "0.0.0.0", port: int = 5010):
    server = ThreadingHTTPServer((host, port), Handler)
    server.daemon_threads = True
    server.engine = engine
    return server


def main() -> None:
    port = int(os.environ.get("FACE_PORT") or os.environ.get("PORT", "5010"))
    host = os.environ.get("HOST", "0.0.0.0")
    log(f"loading models from {MODELS_DIR}")
    ensure_models()
    engine = FaceEngine()
    server = create_server(engine, host, port)
    log(f"listening on http://{host}:{port} (endpoints: /health, /detect)")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass


if __name__ == "__main__":
    main()
