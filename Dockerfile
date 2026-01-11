# --- Builder Stage ---
FROM golang:1.24.4-alpine3.22 AS builder
ARG APP_DIR=/app
WORKDIR ${APP_DIR}
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /dist/web && \
    GOARCH=wasm GOOS=js go build -ldflags="-s -w" -tags wasm -o /dist/web/app.wasm ./cmd/webapp && \
    go build -ldflags="-s -w" -tags unix -o /dist/server ./cmd/webapp/ && \
    cp -a ./web/* /dist/web/

# --- Runner Stage ---
FROM alpine:3.19 AS runner
ARG Port=8000

# 1. Delete the existing group with GID 999 (usually 'ping') 
# 2. Create appgroup with GID 999
# 3. Create appuser with UID 999
RUN delgroup $(getent group 999 | cut -d: -f1) || true && \
    addgroup -S -g 999 appgroup && \
    adduser -S -u 999 -G appgroup appuser

# Install tools and prepare directories
RUN apk add --no-cache curl su-exec && \
    mkdir -p /minecraft/worlds /home/app && \
    chown -R appuser:appgroup /minecraft/worlds /home/app

WORKDIR /home/app

# Copy artifacts from builder
COPY --chown=appuser:appgroup --from=builder /dist .
COPY --chown=root:root entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV APP_ENV=production \
    APP_SERVER_PORT=${Port}

EXPOSE ${Port}

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl --fail "http://localhost:${APP_SERVER_PORT}/healthcheck" || exit 1

# Start as root so the entrypoint can fix volume permissions
USER root 

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["./server"]