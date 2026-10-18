# syntax=docker/dockerfile:1
# Multi-stage: build with the Go toolchain, ship a static binary on distroless as non-root.
FROM --platform=$BUILDPLATFORM golang:1.27-bookworm AS build
WORKDIR /src
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
ARG TARGETOS TARGETARCH
ARG VERSION=dev COMMIT=none DATE=unknown
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" -o /out/ladder-app .

FROM gcr.io/distroless/static-debian12:nonroot
ARG VERSION=dev COMMIT=none DATE=unknown
LABEL org.opencontainers.image.title="ladder-app" \
      org.opencontainers.image.version="$VERSION" \
      org.opencontainers.image.revision="$COMMIT" \
      org.opencontainers.image.created="$DATE"
COPY --from=build /out/ladder-app /ladder-app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/ladder-app"]
