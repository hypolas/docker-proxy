# Example: Run docker-proxy with Docker Compose

This stack starts three services:

- `docker` (Docker-in-Docker) exposes the raw Docker API internally.
- `docker-proxy` uses the published image `hypolas/dockershield` to enforce ACL policies.
- `client` is a Docker CLI container that talks to the proxy.

## Usage

```bash
# from repository root
cd examples/docker-compose

# pull the published image
docker pull hypolas/dockershield:latest

# start the stack
docker-compose up -d

# run Docker commands through the proxy over TCP
docker-compose exec client docker ps

# run Docker commands through the Unix socket
docker-compose exec client env DOCKER_HOST=unix:///tmp/dockershield.sock docker ps

# stop and clean up
docker-compose down -v
```

The proxy is reachable on `tcp://127.0.0.1:12375` and exposes a socket at `/tmp/dockershield.sock`.
You can point any local CLI to the TCP endpoint:

```bash
export DOCKER_HOST=tcp://127.0.0.1:12375
docker ps

# or, if the socket is shared with your host
export DOCKER_HOST=unix:///tmp/dockershield.sock
docker ps
```
