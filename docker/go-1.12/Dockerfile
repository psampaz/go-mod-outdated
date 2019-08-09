FROM golang:1.12.7-alpine3.10

RUN apk add --no-cache git

RUN go get -u github.com/psampaz/go-mod-outdated

WORKDIR /tmp

CMD go list -u -m -json all | go-mod-outdated
