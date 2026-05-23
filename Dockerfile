# Stage 1: Build Go API binary
FROM golang:1.25-alpine AS api-builder

ARG VERSION=dev

WORKDIR /app

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/. ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" -o api .

# Stage 2: Build Webapp
FROM node:22-alpine AS webapp-builder

RUN apk add --no-cache python3 make g++

WORKDIR /app

COPY webapp/package*.json ./
RUN npm ci --legacy-peer-deps

COPY webapp/. ./

ENV VITE_API_URL=/api
RUN npm run build

# Stage 3: Runtime
FROM nginx:alpine

RUN apk add --no-cache ca-certificates tzdata ffmpeg perl-image-exiftool

# Copy API binary
COPY --from=api-builder /app/api /app/api

# Copy webapp build output
COPY --from=webapp-builder /app/build /usr/share/nginx/html

# Copy nginx config (proxies /api/ to localhost:8080)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
