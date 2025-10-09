# Example: Use docker-proxy with a remote Docker CLI

This example assumes you already launched `docker-proxy` (for example with the Compose stack or in another host) and it listens on `tcp://proxy-host:2375`.

## One-off commands

```bash
export DOCKER_HOST=tcp://proxy-host:2375
export DOCKER_TLS_VERIFY=0        # remove if you enable TLS

# read-only operations use GET/HEAD and should succeed
docker ps
docker images

# write operations (POST/DELETE) will be filtered depending on proxy config
docker run --rm alpine:3.19 echo "hello"
```

## Switching between proxy and local dockerd

```bash
# save current socket
export ORIGINAL_DOCKER_HOST=${DOCKER_HOST}

# use the proxy for a task
export DOCKER_HOST=tcp://proxy-host:2375
run-some-command.sh

# restore previous configuration
export DOCKER_HOST=${ORIGINAL_DOCKER_HOST}
```

The CLI only needs the environment variable `DOCKER_HOST` to point at the proxy. All ACLs are enforced server-side by docker-proxy.
