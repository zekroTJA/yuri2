FROM golang:1.13-alpine AS build

RUN apk add --no-cache \
        git curl nodejs nodejs-npm build-base

WORKDIR /build

ADD . .

RUN go build -v -o ./bin/yuri -ldflags "\
            -X github.com/zekroTJA/yuri2/internal/static.AppVersion=$(git describe --tags) \
            -X github.com/zekroTJA/yuri2/internal/static.AppCommit=$(git rev-parse HEAD) \
            -X github.com/zekroTJA/yuri2/internal/static.Release=TRUE" \
        ./cmd/yuri/*.go

RUN cd ./web \
    && npm i \
    && npx ng build --prod=true \
    && mkdir -p ../bin/web/dist \
    && mv dist/web ../bin/web/dist/web


FROM alpine:latest AS final

COPY --from=build /build/bin /app

WORKDIR /app

ENV LAVALINK_ADDR="localhost:2333"
ENV LAVALINK_PW="lavalink_pw"
ENV LAVALINK_LOC="/etc/sounds"

RUN mkdir -p /etc/yuri/config \
    && mkdir -p /etc/yuri/cert \
    && mkdir -p /etc/yuri/db \
    && mkdir -p ${LAVALINK_LOC}

EXPOSE 8080

CMD ./yuri \
        -c "/etc/yuri/config/config.yml" \
        -addr ":8080" \
        -db-dsn "file:/etc/yuri/db/db.sqlite3" \
        -lavalink-address "${LAVALINK_ADDR}" \
        -lavalink-password "${LAVALINK_PW}" \
        -lavalink-location "${LAVALINK_LOC}"
