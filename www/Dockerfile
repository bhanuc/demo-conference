FROM alpine:latest

ENV GOPATH /usr/local

RUN apk --no-cache add go git musl-dev && rm -rf /var/cache/apk/*
RUN go get -u github.com/m3ng9i/ran
COPY /certs/privkey.pem /root/leit.key
COPY /certs/cert.pem /root/leit.pem

WORKDIR /usr/local/src/github.com/pion/demo-conference/www
CMD ran -p 443\
        --cert=/root/leit.pem\
	--key=/root/leit.key
