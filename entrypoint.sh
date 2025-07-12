#!/bin/sh

# This script is executed as root when the container starts.
# It sets up permissions for the Docker socket and then drops privileges
# to run the main application as the non-root 'appuser'.
set -e # Exit immediately if a command exits with a non-zero status.

log() {
    # Prefix log messages for clarity
    echo "[entrypoint] $1"
}

log "Container entrypoint script started"

# Check if the Docker socket is mounted and is a socket file.
if [ -S /var/run/docker.sock ]; then
    log "Docker socket found at /var/run/docker.sock"
    # Get the Group ID (GID) of the Docker socket on the host.
    DOCKER_GID=$(stat -c '%g' /var/run/docker.sock)
    log "Host Docker socket GID: ${DOCKER_GID}"

    # Check if a group with the target GID already exists.
    if ! DOCKER_GROUP_NAME=$(getent group ${DOCKER_GID} | cut -d: -f1); then
        # If not, create a new 'docker' group with this GID.
        log "Group with GID ${DOCKER_GID} not found. Creating a new group 'docker'."
        DOCKER_GROUP_NAME=docker
        addgroup -S -g ${DOCKER_GID} ${DOCKER_GROUP_NAME}
    else
        log "Group '${DOCKER_GROUP_NAME}' with GID ${DOCKER_GID} already exists."
    fi

    # Add 'appuser' to the group that owns the docker socket.
    log "Adding user 'appuser' to group '${DOCKER_GROUP_NAME}'."
    adduser appuser ${DOCKER_GROUP_NAME}
else
    log "Docker socket not found at /var/run/docker.sock. Skipping GID setup."
fi

log "Switching to 'appuser' and executing command: $@"
# Drop root privileges and execute the original command (CMD) as 'appuser'.
# "$@" is the command passed to this script, which is ["./server"] from the Dockerfile CMD.
exec su-exec appuser "$@"
