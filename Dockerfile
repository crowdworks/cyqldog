FROM golang:1.22.5-bookworm AS builder

WORKDIR /go/src/github.com/crowdworks/cyqldog

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /go/src/github.com/crowdworks/cyqldog/bin/cyqldog ./bin/cyqldog
ENTRYPOINT ["/app/bin/cyqldog"]
