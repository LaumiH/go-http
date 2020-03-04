FROM golang:alpine

ADD ./src/http-servers /go/src/http-servers
WORKDIR /go/src/http-servers

ENV PORT=3001

CMD ["go", "run", "http-servers.go"]
