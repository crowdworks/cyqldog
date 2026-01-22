FROM golang:1.25.6-trixie AS builder

WORKDIR /go/src/github.com/crowdworks/cyqldog

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM debian:trixie-slim
WORKDIR /app
COPY --from=builder /go/src/github.com/crowdworks/cyqldog/bin/cyqldog ./bin/cyqldog
ENTRYPOINT ["/app/bin/cyqldog"]
