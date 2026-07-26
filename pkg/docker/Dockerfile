FROM node:25-alpine AS web_ui

RUN npm install -g pnpm

WORKDIR /web/

COPY ui/package.json ui/pnpm-*.yaml ./

RUN pnpm ci

COPY ui/ .

RUN pnpm run build

FROM golang:1.26-alpine AS go_builder

RUN apk update && apk add --no-cache \
    gcc \
    musl-dev \
    libc-dev \
    build-base \
    curl \
    git

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

RUN task --help

WORKDIR /build

COPY core/go.mod core/go.sum core/

RUN cd core && go mod download

COPY .git .git
COPY core core
COPY Taskfile.yml .

RUN task go:b:docker OUT_EXE=../glacier PROD=1

FROM alpine:latest

RUN apk add --no-cache tzdata su-exec

COPY pkg/docker/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /app/

COPY --from=go_builder build/glacier glacier

COPY --from=web_ui /web/build/ web

ENV GLACIER_CONFIG_DIR="/config"
ENV GLACIER_GAME_DIR="/gamestop"

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget -q --spider http://localhost:${GLACIER_SERVER_PORT:-6699}/api/server/public/ping || exit 1

ENTRYPOINT ["entrypoint.sh"]

CMD ["./glacier"]
