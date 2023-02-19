FROM golang:1.20.1-alpine3.17
RUN apk add --no-cache git
WORKDIR /home
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o go-mod-outdated .

FROM scratch
WORKDIR /home/
COPY --from=0 /home/go-mod-outdated .
ENTRYPOINT ["./go-mod-outdated"]
