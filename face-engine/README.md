# face-engine

Standalone Python microservice for face detection and recognition, used by the
Photobox Go API (issue #80). Local/offline inference — no cloud, no external
APIs. Images never leave this process.

- **Detection**: [YuNet](https://github.com/opencv/opencv_zoo/tree/main/models/face_detection_yunet)
  (`face_detection_yunet_2023mar.onnx`) — FPN face detector, raw NCNN-style
  decode (12 outputs).
- **Embeddings**: [SFace](https://github.com/opencv/opencv_zoo/tree/main/models/face_recognition_sface)
  (`face_recognition_sface_2021dec.onnx`) — 128-d, L2-normalized.
- **Runtime**: Python 3.11, `onnxruntime` (CPU), stdlib `http.server`. No
  framework dependencies.

## HTTP API

| Method | Path      | Body                                                        | Response |
|--------|-----------|-------------------------------------------------------------|----------|
| GET    | `/health` | —                                                           | `{"status": "ok"}` |
| POST   | `/detect` | `{"image": "<base64 image bytes>", "min_score": 0.5}`       | see below |

`POST /detect` response:

```json
{
  "faces": [
    {
      "box": [x, y, w, h],
      "score": 0.87,
      "landmarks": [[x, y], [x, y], [x, y], [x, y], [x, y]],
      "embedding": [0.12, -0.03, "... 128 floats, L2-normalized"]
    }
  ]
}
```

- `box` and `landmarks` are in **original image pixel coordinates** (the
  service downscales internally to a 640 px longest side for inference and
  maps results back).
- `min_score` (default 0.5) is the minimum detection confidence
  (`sqrt(cls * obj)` per YuNet convention).
- Faces smaller than 20 px on a side are dropped; NMS at IoU 0.3.
- Errors: `400` for bad JSON / undecodable image, `500` for inference failure.

## Run locally

```bash
python3 -m venv .venv && .venv/bin/pip install -r requirements.txt
.venv/bin/python app.py            # listens on 0.0.0.0:5010
```

Models are downloaded automatically to `./models/` on first run (from
opencv_zoo GitHub raw) if missing. Env vars: `PORT` (default 5010), `HOST`.

Smoke test:

```bash
curl -s localhost:5010/health
python3 generate_testdata.py      # writes privacy-safe synthetic faces to testdata/
IMG=$(base64 < testdata/alice_0.png)
curl -s localhost:5010/detect -d "{\"image\": \"$IMG\"}" | python3 -m json.tool
```

## Docker

```bash
docker build -t photobox-face-engine .
docker run --rm -p 5010:5010 photobox-face-engine
```

Models are baked into the image at build time (the app still re-downloads them
at startup if absent, so a runtime without network works as long as `models/`
is present). Multi-arch (linux/amd64 + linux/arm64) via onnxruntime wheels.

## Self test

```bash
python3 self_test.py
```

Starts the service in-process and verifies: `/health`, detection of ≥1 face in
every synthetic `testdata/` image, 128-d embeddings, and that same-person
embeddings are more similar (cosine) than different-person ones. Exit code 0 =
PASS.

## Test data

`generate_testdata.py` draws 3 distinct cartoon "people" × 3 poses each into
`testdata/*.png` (seeded, deterministic). No real people — safe to commit.
YuNet is trained on real photos, so the style is deliberately minimal (plain
background, simple head/hair/eyes); do not add complex scenery or small faces
or detection may fail.

## Integration with the Go API

The Go service (`api/internal/core/service/face.go`) is the only client: it
base64-encodes a photo and POSTs to `FACE_ENGINE_URL` (default
`http://localhost:5010`), then stores boxes/scores and encrypts the 128-float
embedding for clustering in Postgres. This service is stateless — all
persistence and person clustering live in Go.
