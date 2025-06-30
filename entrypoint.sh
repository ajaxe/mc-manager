#!/bin/sh
# This script is executed as root when the container starts.
# It sets up permissions for the Docker socket and then drops privileges
# to run the main application as the non-root 'appuser'.
set -e

# Check if the Docker socket is mounted and is a socket file.
if [ -S /var/run/docker.sock ]; then
    # Get the Group ID (GID) of the Docker socket on the host.
    DOCKER_GID=$(stat -c '%g' /var/run/docker.sock)

    # Check if a group with the target GID already exists.
    if ! DOCKER_GROUP_NAME=$(getent group ${DOCKER_GID} | cut -d: -f1); then
        # If not, create a new 'docker' group with this GID.
        DOCKER_GROUP_NAME=docker
        addgroup -S -g ${DOCKER_GID} ${DOCKER_GROUP_NAME}
    fi

    # Add 'appuser' to the group that owns the docker socket.
    adduser appuser ${DOCKER_GROUP_NAME}
fi

# Drop root privileges and execute the original command (CMD) as 'appuser'.
# "$@" is the command passed to this script, which is ["./server"] from the Dockerfile CMD.
exec su-exec appuser "$@"