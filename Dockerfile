FROM golang
MAINTAINER Ivan Aracki <aracki.ivan@gmail.com>

ADD . /go/src/github.com/aracki/countgo
COPY ./config.yml /etc/countgo/config.yml

RUN go get ./...
ENV GOBIN=/go/bin
RUN go install /go/src/github.com/aracki/countgo/cmd/aracki/main.go
ENTRYPOINT /go/bin/main

EXPOSE 8080