FROM alpine:latest

ENV GOPATH /usr/local
ENV REALM leit.cc
ENV SERVERIP 206.189.135.92

RUN apk --no-cache add go git musl-dev && rm -rf /var/cache/apk/*
RUN go get -u -v github.com/cespare/reflex
RUN go get -u -v github.com/pion/turn

WORKDIR /usr/local/src/github.com/pion/demo-conference/turn
CMD reflex -r . -s go run main.go
