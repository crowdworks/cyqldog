FROM debian:bullseye-slim

ENV CYQLDOG_VERSION=0.1.3
WORKDIR /app
RUN mkdir -p /app/bin

RUN apt update && \
    apt install -y curl ca-certificates && \
    curl -fsSL https://github.com/crowdworks/cyqldog/releases/download/v${CYQLDOG_VERSION}/cyqldog_${CYQLDOG_VERSION}_linux_amd64.tar.gz | tar -xzC /app/bin && \
    chmod +x /app/bin/cyqldog

ENTRYPOINT ["/app/bin/cyqldog"]
