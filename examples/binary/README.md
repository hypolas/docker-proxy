# Example: Run docker-proxy binary directly

Use this approach when you want to run docker-proxy on a host that already has access to the Docker socket.

## 1. Build or download the binary

```bash
# from repository root
go build -o dockershield ./cmd/dockershield
# or download from your release artifacts into /usr/local/bin/dockershield
docker build -t docker-proxy:dev .   # optional image build
```

## 2. Export the desired configuration

```bash
export DOCKER_SOCKET=unix:///var/run/docker.sock
export LISTEN_ADDR=:2375
export CONTAINERS=1
export IMAGES=1
export POST=0
export DELETE=0
```

All configuration keys can also be stored in a `.env` file and loaded with `env $(cat .env) ./dockershield`.

## 3. Run the proxy

```bash
./dockershield
```

The binary now listens on `tcp://0.0.0.0:2375`. Point your Docker CLI to the proxy:

```bash
export DOCKER_HOST=tcp://127.0.0.1:2375
docker version
```

Press `Ctrl+C` to stop the proxy.

## Systemd service snippet

```
[Unit]
Description=Docker Proxy
After=docker.service

[Service]
EnvironmentFile=/etc/docker-proxy.env
ExecStart=/usr/local/bin/dockershield
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Create `/etc/docker-proxy.env` with the environment variables shown earlier, then run:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now docker-proxy.service
```
