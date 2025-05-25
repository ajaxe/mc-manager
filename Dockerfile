FROM golang:alpine AS builder

ARG InstallFolder=/go/src/github.com/ajaxe/mc-manager

RUN mkdir -p $InstallFolder

WORKDIR $InstallFolder

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=bind,target=. \
    GOARCH=wasm GOOS=js go build -tags wasm -o /root/app/web/app.wasm ./cmd/webapp \
    && go build -tags unix -o /root/app/server ./cmd/webapp/ \
    && cp -a ./web/* /root/app/web/

FROM alpine:latest AS runner

ARG InstallFolder=/root/app \
    Port=8000

RUN apk add --no-cache curl

# local folder to shared minecraft-worlds volume
RUN mkdir -p /minecraft/worlds

RUN mkdir -p /home/app/
WORKDIR /home/app/

COPY --from=builder "${InstallFolder}/." .

ENV APP_ENV=production \
    APP_SERVER_PORT=$Port

EXPOSE $Port

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl --fail "http://localhost:8000/healthcheck" || exit

CMD ["./server"]
