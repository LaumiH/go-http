FROM golang:latest

RUN go get github.com/gorilla/sessions

ADD ./src/http-servers /go/src/http-servers
WORKDIR /go/src/http-servers

CMD ["go", "run", "http-servers.go"]
