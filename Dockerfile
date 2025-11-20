############################
# STEP 1 build base
############################
FROM 097788601444.dkr.ecr.eu-west-1.amazonaws.com/proxy/library/golang:1.19-alpine3.17 as build-base
RUN apk add --update --no-cache git ca-certificates build-base
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download -x

############################
# STEP 2 image base
############################
FROM 097788601444.dkr.ecr.eu-west-1.amazonaws.com/proxy/library/alpine:3.17 as image-base
WORKDIR /app
COPY --from=build-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/parrot", "serve"]

############################
# STEP 3 build executable
############################
FROM build-base AS builder
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/parrot main.go

############################
# STEP 4 Finalize image
############################
FROM image-base
COPY LICENSE .
COPY --from=builder /build/bin/parrot parrot