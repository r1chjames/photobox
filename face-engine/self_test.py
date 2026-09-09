#!/usr/bin/env python3
"""End-to-end self test for face-engine.

Starts the service in-process on a random port, then verifies:
  1. GET /health returns {"status": "ok"}
  2. POST /detect with every synthetic testdata image detects >= 1 face
  3. Same-person embeddings are more similar (cosine) than different-person ones

Exit code 0 = PASS, 1 = FAIL.
"""
from __future__ import annotations

import base64
import json
import os
import sys
import threading
import urllib.request

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)

import app as fe  # noqa: E402


def http_json(method: str, url: str, payload=None):
    data = json.dumps(payload).encode() if payload is not None else None
    req = urllib.request.Request(url, data=data, method=method)
    if data:
        req.add_header("Content-Type", "application/json")
    with urllib.request.urlopen(req, timeout=120) as resp:
        return json.loads(resp.read().decode())


def main() -> int:
    failures = []

    # 1. start server in-process on a random free port
    engine = fe.FaceEngine()
    import socket
    s = socket.socket()
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    server = fe.create_server(engine, "127.0.0.1", port)
    t = threading.Thread(target=server.serve_forever, daemon=True)
    t.start()
    base = f"http://127.0.0.1:{port}"

    # 2. health
    try:
        r = http_json("GET", f"{base}/health")
        if r.get("status") != "ok":
            failures.append(f"health returned {r!r}")
        else:
            print("PASS /health ->", r)
    except Exception as exc:  # noqa: BLE001
        failures.append(f"health request failed: {exc}")

    # 3. detect on all testdata images
    testdata = os.path.join(HERE, "testdata")
    if not os.path.isdir(testdata):
        print("generating testdata ...")
        import generate_testdata
        generate_testdata.main()

    by_person: dict[str, list] = {}
    all_faces = []
    for fn in sorted(os.listdir(testdata)):
        if not fn.endswith(".png"):
            continue
        with open(os.path.join(testdata, fn), "rb") as f:
            b64 = base64.b64encode(f.read()).decode()
        try:
            r = http_json("POST", f"{base}/detect", {"image": b64, "min_score": 0.5})
        except Exception as exc:  # noqa: BLE001
            failures.append(f"detect {fn} failed: {exc}")
            continue
        faces = r.get("faces", [])
        if len(faces) < 1:
            failures.append(f"{fn}: expected >=1 face, got {len(faces)}")
            print(f"FAIL {fn}: no faces detected")
            continue
        person = fn.rsplit("_", 1)[0]
        for face in faces:
            assert len(face["embedding"]) == 128, "embedding must be 128-d"
            assert len(face["box"]) == 4 and len(face["landmarks"]) == 5
            by_person.setdefault(person, []).append(face)
            all_faces.append((fn, face))
        print(f"PASS {fn}: {len(faces)} face(s), "
              f"scores={[f['score'] for f in faces]}")

    # 4. embedding similarity: same-person > different-person (median comparison)
    def cosine(a, b):
        return sum(x * y for x, y in zip(a, b))

    sims = {"same": [], "diff": []}
    persons = list(by_person.keys())
    for i in range(len(persons)):
        for j in range(i + 1, len(persons)):
            for fa in by_person[persons[i]]:
                for fb in by_person[persons[j]]:
                    sims["diff"].append(cosine(fa["embedding"], fb["embedding"]))
    for p in persons:
        faces = by_person[p]
        for i in range(len(faces)):
            for j in range(i + 1, len(faces)):
                sims["same"].append(
                    cosine(faces[i]["embedding"], faces[j]["embedding"]))

    if sims["same"] and sims["diff"]:
        med = lambda xs: sorted(xs)[len(xs) // 2]
        same_med, diff_med = med(sims["same"]), med(sims["diff"])
        print(f"median cosine same-person: {same_med:.3f} | "
              f"different-person: {diff_med:.3f}")
        if same_med > diff_med:
            print("PASS embeddings separate identities")
        else:
            failures.append(
                f"same-person median {same_med:.3f} <= different-person "
                f"{diff_med:.3f}")
    else:
        failures.append("not enough faces to compare embeddings")

    server.shutdown()

    if failures:
        print("\nFAILURES:")
        for f in failures:
            print(" -", f)
        return 1
    print(f"\nPASS ({len(all_faces)} faces across {len(persons)} identities)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
