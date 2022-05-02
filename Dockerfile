FROM golang:1.18.1-alpine3.15
RUN apk add git
WORKDIR /go/src/Zenith

COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/zenith ./main.go

EXPOSE 1323
ENTRYPOINT ./bin/zenith --web-mode --use-proxy