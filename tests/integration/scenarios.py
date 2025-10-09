"""Scenario definitions for socket-mode integration tests."""

SCENARIOS = {
    "readonly": {
        "description": "Read-only access: list commands allowed, mutations blocked.",
        "environment": {
            "TEST_CONTAINERS": "1",
            "TEST_IMAGES": "1",
            "TEST_POST": "0",
            "TEST_DELETE": "0",
            "TEST_PUT": "0",
        },
        "tests": [
            {
                "name": "Container listing is permitted",
                "command": "docker ps",
                "expect": {
                    "returncodes": [0],
                    "regex_any": [r"CONTAINER\s+ID"],
                },
                "messages": {
                    "success": "Container listing succeeds in read-only mode.",
                    "failure": "docker ps should succeed when CONTAINERS=1.",
                },
            },
            {
                "name": "Container creation is blocked",
                "command": "docker run --rm --name socket-readonly-denied nginx:alpine",
                "expect": {
                    "returncodes": [126, 127, 125, 1],
                    "contains_any": [
                        "denied",
                        "forbidden",
                        "blocked",
                        "HTTP 403",
                        "not allowed",
                    ],
                },
                "messages": {
                    "success": "Mutating container operations are blocked as expected.",
                    "failure": "docker run should be denied when POST=0.",
                },
            },
        ],
    },
    "volume_filters": {
        "description": "Volume filters block sensitive mounts while allowing others.",
        "environment": {
            "TEST_CONTAINERS": "1",
            "TEST_VOLUMES": "1",
            "TEST_POST": "1",
            "TEST_IMAGES": "1",
            "TEST_DENIED_PATHS": "/var/run/docker.sock,^/var",
        },
        "tests": [
            {
                "name": "Mounting Docker socket is denied",
                "command": (
                    "docker run --rm --name socket-volume-denied "
                    "-v /var/run/docker.sock:/docker.sock debian:12-slim"
                ),
                "expect": {
                    "returncodes": [126, 127, 125, 1],
                    "contains_any": [
                        "denied",
                        "blocked",
                        "not allowed",
                    ],
                },
                "messages": {
                    "success": "Docker socket mount is correctly blocked by filters.",
                    "failure": "Docker socket mount should be rejected by volume filters.",
                },
            },
            {
                "name": "Transient container without sensitive mounts works",
                "command": "docker run -d --name socket-volume-okay alpine:3.19 sleep 30",
                "expect": {
                    "returncodes": [0],
                    "regex_any": [r"[0-9a-f]{12}"],
                },
                "messages": {
                    "success": "Container without restricted volumes starts successfully.",
                    "failure": "Container without restricted volumes should be allowed.",
                },
                "cleanup": [
                    "docker rm -f socket-volume-okay",
                ],
            },
        ],
    },
    "socket_permissions": {
        "description": "Socket file is created with the expected permissions and Docker CLI works.",
        "environment": {
            "TEST_CONTAINERS": "1",
            "TEST_IMAGES": "1",
        },
        "tests": [
            {
                "name": "Socket file exists",
                "command": "ls -l /tmp/docker-proxy.sock",
                "expect": {
                    "returncodes": [0],
                    "contains_any": ["docker-proxy.sock"],
                },
                "messages": {
                    "success": "Socket file is present on shared volume.",
                    "failure": "Socket file should exist at /tmp/docker-proxy.sock.",
                },
            },
            {
                "name": "Socket permissions are 0666",
                "command": "stat -c '%a' /tmp/docker-proxy.sock",
                "expect": {
                    "returncodes": [0],
                    "contains_any": ["666"],
                },
                "messages": {
                    "success": "Socket permissions expose proxy to clients as configured.",
                    "failure": "Socket permissions should match SOCKET_PERMS=0666.",
                },
            },
            {
                "name": "Docker CLI connects through the socket",
                "command": "docker version",
                "expect": {
                    "returncodes": [0],
                    "contains_any": ["Server:"],
                },
                "messages": {
                    "success": "Client can reach Docker API via the proxy socket.",
                    "failure": "docker version should succeed via proxy socket.",
                },
            },
        ],
    },
}
