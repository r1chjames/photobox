# face-engine — YuNet detection + SFace embeddings (onnxruntime CPU)
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Download models at build time so the image is self-contained (offline at runtime).
# The app re-downloads them at startup if they are missing, but the build-time copy
# keeps first boot fast and works in air-gapped environments.
RUN mkdir -p /app/models \
    && python -c "import urllib.request;urllib.request.urlretrieve('https://github.com/opencv/opencv_zoo/raw/main/models/face_detection_yunet/face_detection_yunet_2023mar.onnx','/app/models/face_detection_yunet_2023mar.onnx')" \
    && python -c "import urllib.request;urllib.request.urlretrieve('https://github.com/opencv/opencv_zoo/raw/main/models/face_recognition_sface/face_recognition_sface_2021dec.onnx','/app/models/face_recognition_sface_2021dec.onnx')"

COPY app.py .

ENV HOST=0.0.0.0 \
    PORT=5010 \
    OMP_NUM_THREADS=2

EXPOSE 5010

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s \
    CMD python -c "import urllib.request;urllib.request.urlopen('http://127.0.0.1:5010/health', timeout=4)" || exit 1

CMD ["python", "app.py"]
