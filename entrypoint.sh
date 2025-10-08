#!/bin/sh
set -e

# Default to /var/run/docker.sock if DOCKER_SOCKET is not set
DOCKER_SOCKET="${DOCKER_SOCKET:-'unix:///var/run/docker.sock'}"
echo "Using DOCKER_SOCKET: $DOCKER_SOCKET"

# Add dkrproxy user to Docker group if /var/run/docker.sock exists
SOCK="${DOCKER_SOCKET#unix://}"
DOCKER_SOCK_GID=$(stat -c '%g' "$SOCK")
echo "Adding dkrproxy user to GID $DOCKER_SOCK_GID for Docker access..."

# Create dockerhost group with Docker socket GID if it doesn't exist
addgroup -g "$DOCKER_SOCK_GID" dockerhost 2>/dev/null || true

# Add dkrproxy user to dockerhost group
addgroup dkrproxy dockerhost 2>/dev/null || true

# Default to dockershield when the first argument looks like an option
if [ "${1#-}" != "$1" ]; then
    set -- dockershield "$@"
fi

if [ "$1" = "dockershield" ]; then
    shift
    DOCKERHOST_GID=$(getent group dockerhost | cut -d: -f3)
    exec setpriv --reuid=dkrproxy --regid=dkrproxy --groups=$DOCKERHOST_GID /app/dockershield "$@"
fi

exec "$@"
