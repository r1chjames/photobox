#!/usr/bin/env python3
"""Generate privacy-safe synthetic cartoon faces for face-engine testing.

Draws 3 distinct "people" x 3 poses each as PNG images in ./testdata.
Deterministic (seeded). No real people, no external assets.

Each person has fixed identity features (skin tone, hair color, eye spacing,
mouth width) that vary across their poses (position, scale), so a face
detector should find >=1 face per image and same-person embeddings should be
more similar than different-person ones.

Style notes: YuNet is trained on real photos; the style below is deliberately
minimal (plain background, simple head oval, dark hair cap, large eyes with
white/iris/pupil, thin nose line, mouth arc) — this is what the detector
responds to reliably on synthetic input.
"""
from __future__ import annotations

import os
import random

from PIL import Image, ImageDraw

HERE = os.path.dirname(os.path.abspath(__file__))
OUT_DIR = os.path.join(HERE, "testdata")

CANVAS = 512          # square canvas
BG = (170, 168, 175)  # neutral gray background


def _face_layer(person: dict, scale: float) -> Image.Image:
    """Draw one cartoon face on a transparent layer, centered."""
    s = scale
    w = h = CANVAS
    img = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)

    skin = person["skin"]
    shadow = person["skin_shadow"]
    hair = person["hair"]
    iris = person["iris"]

    cx = w / 2
    cy = h / 2 - h * 0.04
    fw = w * 0.56 * s   # head width
    fh = w * 0.68 * s   # head height (taller than wide)

    # head
    d.ellipse([cx - fw / 2, cy - fh / 2, cx + fw / 2, cy + fh / 2], fill=skin)

    # hair: dark cap over the top of the head
    cap_h = person["hair_drop"] * s
    d.pieslice([cx - fw / 2 - fw * 0.03, cy - fh / 2 - fw * 0.05,
                cx + fw / 2 + fw * 0.03, cy - fh / 2 + cap_h],
               start=180, end=360, fill=hair)

    # eyes (white sclera + iris + pupil)
    eye_dx = person["eye_dx"] * s
    eye_dy = -person["eye_dy"] * s
    er = person["eye_r"] * s
    for sx in (-1, 1):
        ex = cx + sx * eye_dx
        ey = cy + eye_dy
        d.ellipse([ex - er * 1.6, ey - er * 1.2, ex + er * 1.6, ey + er * 1.2],
                  fill=(250, 248, 244))
        d.ellipse([ex - er * 0.9, ey - er * 0.9, ex + er * 0.9, ey + er * 0.9],
                  fill=iris)
        d.ellipse([ex - er * 0.35, ey - er * 0.35, ex + er * 0.35, ey + er * 0.35],
                  fill=(20, 18, 16))
        # eyebrow
        d.line([ex - er * 1.7, ey - er * 2.3, ex + er * 1.7, ey - er * 2.5],
               fill=hair, width=max(3, int(er * 0.8)))

    # nose: short vertical line with a small base arc
    nx, ny = cx, cy + fh * 0.08
    d.line([nx, cy - fh * 0.04, nx, ny], fill=shadow, width=max(3, int(5 * s)))
    d.arc([nx - fw * 0.07, ny - fh * 0.03, nx + fw * 0.07, ny + fh * 0.05],
          start=20, end=160, fill=shadow, width=max(2, int(3 * s)))

    # mouth: person-specific width arc
    mw = person["mouth_w"] * s
    my = cy + fh * 0.28
    d.arc([cx - mw / 2, my - fh * 0.09, cx + mw / 2, my + fh * 0.15],
          start=15, end=165, fill=(140, 60, 60), width=max(3, int(6 * s)))

    return img


def make_image(person: dict, pose: dict) -> Image.Image:
    base = Image.new("RGB", (CANVAS, CANVAS), BG)
    layer = _face_layer(person, pose["scale"])
    dx = int((pose["cx"] - 0.5) * CANVAS)
    dy = int((pose["cy"] - 0.5) * CANVAS)
    base.paste(layer, (dx, dy), layer)
    return base


def build_people(rng: random.Random) -> dict:
    """Three distinct identities with fixed features."""
    return {
        "alice": dict(
            skin=(238, 198, 168), skin_shadow=(195, 145, 115),
            hair=(70, 45, 32), iris=(60, 105, 160),
            eye_dx=64, eye_dy=30, eye_r=18, mouth_w=92, hair_drop=100,
        ),
        "bob": dict(
            skin=(180, 132, 96), skin_shadow=(135, 92, 64),
            hair=(28, 24, 22), iris=(62, 46, 32),
            eye_dx=74, eye_dy=26, eye_r=16, mouth_w=108, hair_drop=70,
        ),
        "carol": dict(
            skin=(248, 224, 208), skin_shadow=(208, 178, 160),
            hair=(200, 150, 65), iris=(85, 135, 95),
            eye_dx=58, eye_dy=34, eye_r=20, mouth_w=78, hair_drop=130,
        ),
    }


def main() -> None:
    rng = random.Random(80)
    os.makedirs(OUT_DIR, exist_ok=True)
    # remove stale generated fixtures so testdata matches this script exactly
    for fn in os.listdir(OUT_DIR):
        if fn.endswith(".png"):
            os.remove(os.path.join(OUT_DIR, fn))

    written = []
    for name, person in build_people(rng).items():
        for i in range(3):
            pose = dict(
                scale=rng.uniform(0.95, 1.08),
                cx=rng.uniform(0.47, 0.53),
                cy=rng.uniform(0.50, 0.56),
            )
            img = make_image(person, pose)
            path = os.path.join(OUT_DIR, f"{name}_{i}.png")
            img.save(path)
            written.append(path)
    print(f"wrote {len(written)} synthetic faces to {OUT_DIR}")
    for p in written:
        print(" ", os.path.basename(p))


if __name__ == "__main__":
    main()
