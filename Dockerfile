FROM golang:alpine3.12
LABEL maintainer="Richard James<richjames11@gmail.com>"

RUN apk add make

WORKDIR /app

COPY * ./

RUN make compile

EXPOSE 8080

ENTRYPOINT ["/app/bin", "api"]

