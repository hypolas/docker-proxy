# Usage Examples for docker-proxy

This directory shows three typical deployment patterns:

| Example | Description |
|---------|-------------|
| `docker-compose/` | Full stack with Docker-in-Docker, the proxy, and a test client managed via Compose. |
| `docker-cli/` | How to target an existing proxy using only the Docker CLI and environment variables. |
| `binary/` | Run the compiled `dockershield` binary directly on a host or via systemd. |

Use these as starting points and adapt the configuration to your environment (registries, ACL rules, network ports, etc.).
