FROM 1.12.6-stretch

RUN apt update -y &&\
    apt install -y \
      git

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/zekroTJA/yuri2


