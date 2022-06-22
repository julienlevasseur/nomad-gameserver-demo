FROM golang:1.18-alpine
LABEL maintainer="julien.levasseur@free.fr"

COPY $PWD /go/src/github.com/julienlevasseur/nomad-gameserver-demo
WORKDIR /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
RUN apk update && \
apk add git && \
go install && \
go build && \
apk del git

CMD ["./api"]
