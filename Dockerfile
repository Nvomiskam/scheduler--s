FROM golang:1.24.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY pkg ./pkg
COPY web ./web

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-scheduler main.go

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

RUN mkdir -p /data

COPY --from=builder /todo-scheduler /usr/local/bin/todo-scheduler
COPY --from=builder /app/web /web

VOLUME /data

CMD ["todo-scheduler"]