# syntax=docker/dockerfile:1.7

FROM golang:1.25.13-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/profile-api ./cmd/profile-api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/profile-api /app/profile-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/profile-api"]
