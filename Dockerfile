#### Build Step ####
FROM golang:1.9-alpine as builder

WORKDIR src/bitbucket.org/exonch/ch-gateway
COPY . .

RUN CGO_ENABLED=0 go build -v -o /bin/ch-gateway -ldflags "-X bitbucket.org/exonch/ch-gateway/version.version=${APP_VERSION}" cmd/*

#### Run Step ####
FROM scratch

# Copy bin
COPY --from=builder /bin/ch-gateway /

# Set envs
ENV GATEWAY_DEBUG=false \
    PG_USER="pg" \
    PG_PASSWORD="123456789" \
    PG_DATABASE="postgres" \
    PG_ADDRESS="x1.containerum.io:36519" \
    STATSD_ADDRESS="213.239.208.25:8125" \
    STATSD-PREFIX="ch-gateway" \
    STATSD-BUFFER-TIME=300

# run app
ENTRYPOINT ["/ch-gateway"]

EXPOSE 8080
