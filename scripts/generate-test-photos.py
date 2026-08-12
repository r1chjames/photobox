#!/usr/bin/env python3
"""Generate synthetic test photos with EXIF for Photobox perf testing.

Usage:
    python3 generate-test-photos.py [count] [output_dir]

Generates `count` JPEG photos (default 2000) across several album directories,
each with EXIF DateTimeOriginal (spread over ~3 years), GPS coordinates, and a
camera model, so the timeline scrubber, map view, and photo grid all have data.
"""
import os
import random
import sys
import math
from datetime import datetime, timedelta

from PIL import Image, ImageDraw
import piexif

# Album layout: (dirname, weight) — weight controls how many photos land there.
ALBUMS = [
    ("Album1", 3),
    ("Album2", 3),
    ("Vacation/2024", 2),
    ("Vacation/2023", 2),
    ("Family", 2),
    ("Random", 2),
]

CAMERA_MODELS = [
    "Canon EOS R5",
    "Sony A7 IV",
    "iPhone 15 Pro",
    "Nikon Z8",
    "Google Pixel 9",
]

# Date range: 2023-01-01 .. 2026-01-01 (3 years)
START = datetime(2023, 1, 1)
END = datetime(2026, 1, 1)
TOTAL_SECONDS = (END - START).total_seconds()

# GPS range: somewhere in the UK/London area with some spread.
BASE_LAT, BASE_LNG = 51.5074, -0.1278


def random_date() -> datetime:
    delta = random.randint(0, int(TOTAL_SECONDS))
    return START + timedelta(seconds=delta)


def make_jpeg(path: str, width: int, height: int, date: datetime, model: str) -> None:
    # Build a simple but non-trivial image (gradient + shapes) so thumbnails vary.
    r0, g0, b0 = random.randint(0, 255), random.randint(0, 255), random.randint(0, 255)
    r1, g1, b1 = random.randint(0, 255), random.randint(0, 255), random.randint(0, 255)
    img = Image.new("RGB", (width, height))
    draw = ImageDraw.Draw(img)
    for y in range(height):
        t = y / max(height - 1, 1)
        row = (int(r0 + (r1 - r0) * t), int(g0 + (g1 - g0) * t), int(b0 + (b1 - b0) * t))
        draw.line([(0, y), (width, y)], fill=row)

    # A couple of geometric accents.
    for _ in range(3):
        x0 = random.randint(0, width // 2)
        y0 = random.randint(0, height // 2)
        x1 = random.randint(x0, width)
        y1 = random.randint(y0, height)
        fill = (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))
        if random.random() < 0.5:
            draw.rectangle([x0, y0, x1, y1], fill=fill)
        else:
            draw.ellipse([x0, y0, x1, y1], fill=fill)

    # EXIF: camera + capture time + GPS.
    lat = BASE_LAT + random.uniform(-2.0, 2.0)
    lng = BASE_LNG + random.uniform(-2.0, 2.0)

    def to_deg(value: float):
        d = int(value)
        minutes = (abs(value) - abs(d)) * 60
        m = int(minutes)
        s = int(round((minutes - m) * 60))
        return ((abs(d), 1), (m, 1), (s, 1))

    exif_dict = {"0th": {}, "Exif": {}, "GPS": {}, "1st": {}, "thumbnail": None}
    exif_dict["0th"][piexif.ImageIFD.Make] = model.split()[0].encode()
    exif_dict["0th"][piexif.ImageIFD.Model] = model.encode()
    exif_dict["Exif"][piexif.ExifIFD.DateTimeOriginal] = date.strftime("%Y:%m:%d %H:%M:%S").encode()
    exif_dict["GPS"][piexif.GPSIFD.GPSLatitude] = to_deg(lat)
    exif_dict["GPS"][piexif.GPSIFD.GPSLatitudeRef] = "N" if lat >= 0 else "S"
    exif_dict["GPS"][piexif.GPSIFD.GPSLongitude] = to_deg(lng)
    exif_dict["GPS"][piexif.GPSIFD.GPSLongitudeRef] = "E" if lng >= 0 else "W"

    try:
        exif_bytes = piexif.dump(exif_dict)
        img.save(path, "JPEG", quality=82, exif=exif_bytes)
    except Exception:
        # Fall back to no EXIF if encoding fails for any reason.
        img.save(path, "JPEG", quality=82)


def main() -> None:
    count = int(sys.argv[1]) if len(sys.argv) > 1 else 2000
    out_dir = sys.argv[2] if len(sys.argv) > 2 else "photos"

    total_weight = sum(w for _, w in ALBUMS)
    rng = random.Random(12345)  # deterministic

    made = 0
    for album, weight in ALBUMS:
        album_path = os.path.join(out_dir, album)
        os.makedirs(album_path, exist_ok=True)
        album_count = round(count * weight / total_weight)
        for i in range(album_count):
            width = rng.choice([1200, 1600, 2000, 2400])
            height = rng.choice([900, 1200, 1500, 1600])
            date = START + timedelta(seconds=rng.randint(0, int(TOTAL_SECONDS)))
            model = rng.choice(CAMERA_MODELS)
            fname = os.path.join(album_path, f"{date.strftime('%Y%m%d_%H%M%S')}_{i:04d}.jpg")
            make_jpeg(fname, width, height, date, model)
            made += 1

    print(f"Generated {made} photos in {out_dir}")


if __name__ == "__main__":
    main()
