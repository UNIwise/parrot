############################
# STEP 1 Build base
############################
FROM --platform=$BUILDPLATFORM golang:1.26-alpine3.23 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /build

# Cache dependencies to speed up subsequent builds
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x

# Copy source code and build the binary
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /build/bin/parrot main.go

############################
# STEP 2 Finalize image
############################
FROM gcr.io/distroless/static-debian13:nonroot AS runtime
WORKDIR /app
COPY LICENSE .

COPY --from=builder /build/bin/parrot /usr/bin/parrot

USER 10000:10000

# serve stays in the entrypoint: the deployments pass --config as args, which replace CMD.
ENTRYPOINT [ "parrot", "serve" ]
