# # BASE API BUILD STAGE
# FROM golang:alpine3.15 as api-build
# LABEL maintainer="Richard James<richjames11@gmail.com>"
# ENV GO111MODULE=on

# RUN apk add make git

# WORKDIR /app

# COPY api/go.mod api/go.sum api/Makefile ./
# RUN go mod download

# COPY api/. ./

# RUN GO111MODULE=on go build -o /app/bin/api


# # BASE WEBAPP BUILD STAGE
# FROM node:current-alpine3.15 as webapp-build
# LABEL maintainer="Richard James<richjames11@gmail.com>"

# RUN apk add make

# WORKDIR /app

# COPY webapp/. ./

# RUN make build_app


# APP IMAGE BUILD STAGE
FROM alpine:3.13.5

RUN apk add file git tzdata

# COPY --from=api-build /app/bin/api /app
# COPY --from=webapp-build /app/build/. /web

COPY api/bin/api /app/api
# COPY app/build/ /app/web

EXPOSE 8080
ENTRYPOINT ["/app/api"]
