#!/bin/sh
set -e

log() {
    echo "[entrypoint] $1"
}

log "Container entrypoint script started"

# 1. FIX VOLUME PERMISSIONS
# Ensures appuser (999) owns the mounted volume for unlinkat operations
if [ -d /minecraft/worlds ]; then
    log "Ensuring appuser (999) owns /minecraft/worlds..."
    chown -R appuser:appgroup /minecraft/worlds
fi

# 2. DOCKER SOCKET GID SETUP
if [ -S /var/run/docker.sock ]; then
    log "Docker socket found at /var/run/docker.sock"
    DOCKER_GID=$(stat -c '%g' /var/run/docker.sock)

    if getent group ${DOCKER_GID} > /dev/null 2>&1; then
        DOCKER_GROUP_NAME=$(getent group ${DOCKER_GID} | cut -d: -f1)
    else
        log "Creating group 'docker' with GID ${DOCKER_GID}"
        DOCKER_GROUP_NAME=docker
        addgroup -S -g ${DOCKER_GID} ${DOCKER_GROUP_NAME}
    fi

    log "Adding appuser to group ${DOCKER_GROUP_NAME}"
    adduser appuser ${DOCKER_GROUP_NAME}
fi

log "Dropping root privileges and starting server..."
exec su-exec appuser "$@"