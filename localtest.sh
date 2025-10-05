#!/bin/bash

export DOCKER_HOST=unix:///var/run/docker.sock
export LISTEN_SOCKET=unix:///tmp/docker_proxy.sock
export CONTAINERS=1
export IMAGES=1

go run ./cmd/docker-proxy
