# BASE DEP BUILD STAGE
FROM golang:alpine3.12 as build
LABEL maintainer="Richard James<richjames11@gmail.com>"
ENV GO111MODULE=on

RUN apk add make git

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY . ./

RUN GO111MODULE=on go build -o /app/bin


# APP IMAGE BUILD STAGE
FROM golang:alpine3.12

COPY --from=build /app /app

EXPOSE 8080
ENTRYPOINT ["/app/bin", "gitlab.com/r1chjames/photobox/api"]
