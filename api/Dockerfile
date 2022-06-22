FROM golang:1.18-alpine
LABEL maintainer="julien.levasseur@free.fr"

COPY $PWD /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api
WORKDIR /go/src/github.com/julienlevasseur/nomad-gameserver-demo/api/api
RUN apk update && \
apk add git && \
pwd && ls -la && \
go install && \
go build && \
apk del git

CMD ["./api"]
