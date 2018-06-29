FROM golang:1.10.3-alpine as builder
WORKDIR $GOPATH/src/github.com/aracki/countgo
COPY . .
ENV GOBIN $GOPATH/bin
RUN GOOS=linux GOARCH=386 go install cmd/aracki/main.go

FROM scratch as appgo
COPY --from=builder /go/bin/main /go/bin/main
COPY mongo_config.yml mongo_config.yml
ENTRYPOINT ["/go/bin/main"]
