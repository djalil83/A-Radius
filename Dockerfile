# syntax=docker/dockerfile:1.7

FROM golang:1.25-bookworm AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /out/profile-api ./cmd/profile-api

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install --no-install-recommends -y ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --no-create-home --uid 10001 app

COPY --from=builder /out/profile-api /usr/local/bin/profile-api
USER app
EXPOSE 8080

ENV PROFILE_API_ADDR=:8080
ENTRYPOINT ["/usr/local/bin/profile-api"]
