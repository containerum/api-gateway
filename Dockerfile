#### Build Step ####
FROM golang:1.9-alpine as builder

WORKDIR src/git.containerum.net/ch/api-gateway
COPY . .

RUN CGO_ENABLED=0 go build -v -o /bin/ch-gateway \
    -ldflags "-X git.containerum.net/ch/api-gateway/main.version=${APP_VERSION}" cmd/*

COPY pkg/store/migrations /pkg/store/migrations

#### Run Step ####
# FROM scratch
FROM ubuntu

# Copy bin
COPY --from=builder /bin/ch-gateway /
COPY --from=builder /pkg/store/migrations /pkg/store/migrations

# Set envs
ENV GATEWAY_DEBUG=false \
    PG_USER="pg" \
    PG_PASSWORD="123456789" \
    PG_DATABASE="postgres" \
    PG_ADDRESS="x1.containerum.io" \
    PG_PORT="36519" \
    PG_MIGRATIONS="true" \
    STATSD_ADDRESS="213.239.208.25:8125" \
    STATSD-PREFIX="ch-gateway" \
    STATSD-BUFFER-TIME=300 \
    GRPC_AUTH_ADDRESS="192.168.88.200" \
    GRPC_AUTH_PORT="1112" \
    REDIS_ADDRESS="192.168.88.200:6379" \
    REDIS_PASSWORD="" \
    RATE_LIMIT="3" \
    CLICKHOUSE_LOGGER="88.99.160.131:7777" \
    TLS_CERT="cert.pem" \
    TLS_KEY="key.pem"

# run app
ENTRYPOINT ["/ch-gateway"]

EXPOSE 8082
