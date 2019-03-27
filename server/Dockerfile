FROM golang:1.12.1

WORKDIR /ws-server

ADD main.go go.sum go.mod /ws-server/

RUN go build

ENTRYPOINT ["/ws-server/product-service"]