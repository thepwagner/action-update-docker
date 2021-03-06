FROM golang:1.16.6 AS builder

WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app
RUN go build -o /update-docker /app

FROM debian:buster-slim
COPY --from=builder /update-docker /update-docker
ENTRYPOINT ["/update-docker"]
