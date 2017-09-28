FROM alpine:3.6

ENV CYQLDOG_VERSION=0.1.1
WORKDIR /app
RUN mkdir -p /app/bin

RUN apk --no-cache add curl ca-certificates && update-ca-certificates
RUN curl -fsSL https://github.com/crowdworks/cyqldog/releases/download/v${CYQLDOG_VERSION}/cyqldog_${CYQLDOG_VERSION}_linux_amd64.tar.gz \
    | tar -xzC /app/bin && chmod +x /app/bin/cyqldog
RUN apk del --purge curl

ENTRYPOINT ["/app/bin/cyqldog"]
