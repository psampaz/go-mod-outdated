FROM golang:1.11.12-alpine3.10

RUN apk add --no-cache git

RUN go get -u github.com/psampaz/go-mod-outdated

WORKDIR /tmp

CMD go list -u -m -json all | go-mod-outdated
