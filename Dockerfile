FROM golang:1.20.7-bullseye AS builder

WORKDIR /go/src/github.com/crowdworks/cyqldog

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /go/src/github.com/crowdworks/cyqldog/bin/cyqldog ./bin/cyqldog
ENTRYPOINT ["/app/bin/cyqldog"]
