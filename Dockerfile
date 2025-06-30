# --- Builder Stage ---
# Use a specific version of golang:alpine for reproducibility
FROM golang:1.24.4-alpine3.22 AS builder

# Define the application directory
ARG APP_DIR=/app
WORKDIR ${APP_DIR}

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
# Download dependencies. This creates a separate layer for dependencies.
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go applications (server and wasm)
# Using BuildKit cache mounts for faster builds.
# -ldflags="-s -w" strips debug information, reducing binary size.
# Output artifacts to a dedicated /dist directory for easy copying.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /dist/web && \
    GOARCH=wasm GOOS=js go build -ldflags="-s -w" -tags wasm -o /dist/web/app.wasm ./cmd/webapp && \
    go build -ldflags="-s -w" -tags unix -o /dist/server ./cmd/webapp/ && \
    cp -a ./web/* /dist/web/

# --- Runner Stage ---
# Use a specific, minimal base image
FROM alpine:3.19 AS runner

# Port for the server to listen on
ARG Port=8000

# Create a non-root user and group for security
# It's good practice to run applications as a non-root user.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install curl for healthcheck and create necessary directories.
# Combining these into one RUN command reduces image layers.
# su-exec is a lightweight tool to drop root privileges.
RUN apk add --no-cache curl su-exec && \
    mkdir -p /minecraft/worlds /home/app && \
    chown -R appuser:appgroup /minecraft/worlds /home/app

# Set the working directory
WORKDIR /home/app

# Copy artifacts and the entrypoint script.
COPY --chown=appuser:appgroup --from=builder /dist .
COPY --chown=root:root entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

# The container will start as root, and the entrypoint script will
# perform setup before dropping privileges to 'appuser'.

# Set environment variables
ENV APP_ENV=production \
    APP_SERVER_PORT=${Port}

# Expose the application port
EXPOSE ${Port}

# Healthcheck to ensure the server is running.
# Using ${Port} makes it dynamic. `|| exit 1` ensures a non-zero exit code on failure.
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl --fail "http://localhost:${APP_SERVER_PORT}/healthcheck" || exit 1

# Set the entrypoint to our script.
# The CMD will be passed as arguments to this script.
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Command to run the application
CMD ["./server"]
