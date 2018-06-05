#### Build Step ####
FROM golang:1.10-alpine as builder

WORKDIR src/git.containerum.net/ch/api-gateway
COPY . .

RUN go build -v -ldflags="-w -s" -o /bin/api-gateway ./cmd/api-gateway

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
RUN mkdir -p /app
COPY --from=builder /go/src/git.containerum.net/ch/api-gateway/charts/api-gateway/env/config.toml /app
COPY --from=builder /go/src/git.containerum.net/ch/api-gateway/charts/api-gateway/env/routes /app/routes
COPY --from=builder /bin/api-gateway /app

# Copy certs
COPY --from=generator /cert /cert

# Set envs
ENV GATEWAY_DEBUG=false \
    GRPC_AUTH_ADDRESS="127.0.0.1:1112" \
    CONFIG_FILE="config.toml" \
    ROUTES_FILE="routes/routes.toml" \
    TLS_CERT="/cert/cert.pem" \
    TLS_KEY="/cert/key.pem" \
    SERVICE_HOST_PREFIX=""

EXPOSE 8082 8282

# run app
WORKDIR "/app"
CMD "./api-gateway"
