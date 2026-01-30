FROM node:25-alpine AS web_ui

WORKDIR /web/

COPY ui/package* .

RUN npm ci

COPY ui/ .

RUN npm run build

FROM golang:1-alpine AS go_builder

# for sqlite
ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache gcc musl-dev curl git

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

ENTRYPOINT ["entrypoint.sh"]

CMD ["./glacier"]
