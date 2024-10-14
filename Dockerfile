############################
# Backend build
############################
FROM golang:1.23-alpine3.20 AS build-base
RUN apk add --update --no-cache git ca-certificates build-base
WORKDIR /build
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download -x

FROM build-base AS backend-builder
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/bin/parrot main.go

############################
# Frontend build
############################
FROM node:20-alpine3.20 AS frotnend-builder
WORKDIR /app

RUN npm install -g pnpm

COPY ./dashboard/package.json .
COPY ./dashboard/pnpm-lock.yaml .
RUN pnpm install --frozen-lockfile

COPY ./dashboard /app
RUN pnpm build

############################
# Image base
############################
FROM nginx:1.27.2-alpine3.20 AS image-base
COPY ./docker/nginx.conf /etc/nginx
COPY ./docker/entrypoint.sh /
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

############################
# Finalize image
############################
FROM image-base
WORKDIR /
COPY LICENSE .
COPY --from=backend-builder /build/bin/parrot /usr/local/bin/parrot
COPY --from=frotnend-builder /app/dist /usr/share/nginx/html
