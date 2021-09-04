# syntax=docker/dockerfile:1.3
ARG GO_VERSION

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache git
COPY --from=goreleaser-xx / /
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
ARG GIT_REF
RUN --mount=type=bind,target=/src,rw \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=target=/go/pkg/mod,type=cache \
  goreleaser-xx --debug \
    --name "go-mod-outdated" \
    --dist "/out" \
    --hooks="go mod tidy" \
    --hooks="go mod download" \
    --ldflags="-s -w" \
    --files="CHANGELOG.md" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifacts
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM scratch
WORKDIR /home/
COPY --from=build /usr/local/bin/go-mod-outdated .
ENTRYPOINT ["./go-mod-outdated"]
