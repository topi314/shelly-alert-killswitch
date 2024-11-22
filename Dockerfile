FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -ldflags="-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT'" -o shelly-alert-killswitch github.com/topi314/shelly-alert-killswitch

FROM alpine

RUN apk add --no-cache  \
    inkscape \
    ttf-freefont

COPY --from=build /build/shelly-alert-killswitch /bin/shelly-alert-killswitch

EXPOSE 80

ENTRYPOINT ["/bin/shelly-alert-killswitch"]

CMD ["-config", "/var/lib/shelly-alert-killswitch/config.toml"]