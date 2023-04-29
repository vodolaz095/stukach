FROM golang:1.19 AS builder

VOLUME /build/

ADD go.mod /app/go.mod
ADD go.sum /app/go.sum
WORKDIR /app/
RUN go mod download
RUN go mod verify
ADD . /app/

CMD go build -o /build/stukach main.go

