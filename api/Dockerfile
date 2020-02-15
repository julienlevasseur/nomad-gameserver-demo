FROM golang:1.13.8-alpine
MAINTAINER Julien Levasseur <contact@julienlevasseur.ca>

COPY $PWD/api /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
WORKDIR /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
RUN apk update && \
apk add git && \
go get && \
go build && \
apk del git

CMD ["./api"]
