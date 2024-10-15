############################
# Backend build
############################
FROM golang:1.23-alpine3.20 AS backend-builder
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /build/bin/parrot main.go

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
# Finalize image
############################
FROM nginx:1.27.2-alpine3.20 AS image-base
WORKDIR /
ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

COPY LICENSE .

COPY ./docker/nginx.conf /etc/nginx
COPY ./docker/entrypoint.sh /

COPY --from=backend-builder /build/bin/parrot /usr/local/bin/parrot
COPY --from=frotnend-builder /app/dist /usr/share/nginx/html
