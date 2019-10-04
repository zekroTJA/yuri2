FROM golang:1.12.6-stretch

ENV LAVALINK_ADDR="localhost:2333"
ENV LAVALINK_PW="lavalink_pw"
ENV LAVALINK_LOC="/etc/sounds"

RUN apt update -y &&\
    apt install -y \
      git

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - &&\
        apt-get install -y nodejs &&\
        npm install -g @angular/cli

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/zekroTJA/yuri2

ADD . .

RUN dep ensure -v

RUN go build -v -o ./bin/yuri -ldflags "\
            -X github.com/zekroTJA/yuri2/internal/static.AppVersion=$(git describe --tags) \
            -X github.com/zekroTJA/yuri2/internal/static.AppCommit=$(git rev-parse HEAD) \
            -X github.com/zekroTJA/yuri2/internal/static.Release=TRUE" \
        ./cmd/yuri/*.go

RUN cd ./web &&\
        npm i &&\
        ng build --prod=true

RUN mkdir -p /etc/yuri/config &&\
    mkdir -p /etc/yuri/cert &&\
    mkdir -p /etc/yuri/db &&\
    mkdir -p ${LAVALINK_LOC}

EXPOSE 8080

CMD ./bin/yuri \
        -c "/etc/yuri/config/config.yml" \
        -addr ":8080" \
        -db-dsn "file:/etc/yuri/db/db.sqlite3" \
        -lavalink-address "${LAVALINK_ADDR}" \
        -lavalink-password "${LAVALINK_PW}" \
        -lavalink-location "${LAVALINK_LOC}"
