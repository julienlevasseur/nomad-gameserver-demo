FROM golang:1.13.8-alpine
MAINTAINER Julien Levasseur <contact@julienlevasseur.ca>

COPY $PWD /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
WORKDIR /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
RUN apk update && \
apk add git sshd && \
go get && \
go build && \
apk del git

CMD ["./api"]
