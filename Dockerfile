#### Build Step ####
FROM golang:1.9-alpine as builder

WORKDIR src/git.containerum.net/ch/api-gateway
COPY . .

RUN CGO_ENABLED=0 go build -v -o /bin/ch-gateway cmd/*

#### Generate Cert Step ####
FROM alpine as generator

RUN apk update && \
  apk add --no-cache openssl && \
  rm -rf /var/cache/apk/*

WORKDIR /cert

RUN openssl req -subj '/CN=containerum.io/O=Containerum/C=LV' -new -newkey rsa:2048 -sha256 -days 365 -nodes -x509 -keyout key.pem -out cert.pem

#### Run Step ####
FROM alpine

# Copy bin and migrations
COPY --from=builder /go/src/git.containerum.net/ch/api-gateway/charts/api-gateway/env/config.toml /
COPY --from=builder /go/src/git.containerum.net/ch/api-gateway/charts/api-gateway/env/routes/ routes/
COPY --from=builder /bin/ch-gateway /

# Copy certs
COPY --from=generator /cert /cert

# Set envs
ENV GATEWAY_DEBUG=false \
    GRPC_AUTH_ADDRESS="127.0.0.1" \
    GRPC_AUTH_PORT="1112" \
    CONFIG_FILE="config.toml" \
    ROUTES_FILE="/routes/routes.toml" \
    TLS_CERT="/cert/cert.pem" \
    TLS_KEY="/cert/key.pem"

# run app
ENTRYPOINT ["/ch-gateway"]

EXPOSE 8082 8282
