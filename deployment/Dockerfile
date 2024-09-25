FROM golang:1.22-alpine

RUN apk update && \
    apk add bash git && \
    apk add gcc && \
    apk add musl-dev && \
    apk add curl && \
    apk add openssh && \
    apk add --update make


COPY . /home/golang/src
WORKDIR /home/golang/src
