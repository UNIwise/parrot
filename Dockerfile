############################
# STEP 1 build base
############################
FROM golang:1.15-alpine3.12 as build-base
RUN apk add --update --no-cache git ca-certificates build-base
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download -x

############################
# STEP 2 image base
############################
FROM alpine:3.11 as image-base
# Set workdir
WORKDIR /app
# Copy certificates
COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Run the main binary
ENTRYPOINT ["/app/service"]

############################
# STEP 3 build executable
############################
FROM build-base AS builder
# Copy src
COPY . .
# Build the binary.
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/service cmd/service/main.go

############################
# STEP 4 Finalize image
############################
FROM image-base
# Copy our static executable
COPY --from=builder /build/bin/service service