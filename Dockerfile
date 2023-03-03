FROM golang

ADD . /go/src/app
WORKDIR /go/src/app

RUN make build

ENTRYPOINT /go/src/app/run.sh

EXPOSE 8080
