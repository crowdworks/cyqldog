FROM golang:1.8 AS builder

WORKDIR /go/src/github.com/crowdworks/cyqldog
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cyqldog

FROM alpine:3.6
WORKDIR /app
COPY --from=builder /go/src/github.com/crowdworks/cyqldog/bin/cyqldog ./
COPY --from=builder /go/src/github.com/crowdworks/cyqldog/config/ ./config/
ENTRYPOINT ["./cyqldog"]


