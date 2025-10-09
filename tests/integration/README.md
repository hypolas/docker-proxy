# Unix Socket Integration Tests

These integration tests validate the docker-proxy application when exposed through a Unix socket.
The workflow is driven by a single Python runner backed by dictionaries that specify the environment
variables to override, the expected command output (via substring/regex checks), and the message to
print for success or failure.

## File Layout

```
tests/integration/
├── docker-compose.socket.yml   # Test environment (DinD, proxy, client)
├── run_socket_tests.py         # Main runner
└── scenarios.py                # Dictionary of scenarios and test cases
```

## Requirements

- Docker Engine with `docker-compose`
- Python 3.10+
- Ability to run privileged containers (required by `docker:dind`)

## How It Works

`scenarios.py` defines a `SCENARIOS` dictionary. Each entry provides:

- `environment`: key/value pairs injected into the docker-compose stack (e.g. `TEST_CONTAINERS=1`).
- `tests`: ordered test cases, each with:
  - `command`: shell command executed inside the test client (`docker ps`, `docker run`, etc.).
  - `expect`: expectations for the command (allowed return codes, substrings, regex patterns).
  - `messages`: success/failure messages displayed depending on the result.
  - optional `cleanup`: commands executed after the test to remove temporary containers or volumes.

The runner spins up the compose stack once per scenario, executes the commands inside the
`test-client` service, evaluates expectations, reports a summary, and tears the stack down.

## Usage

From the repository root:

```bash
cd tests/integration
python3 run_socket_tests.py --list           # Show available scenarios
python3 run_socket_tests.py                  # Run all scenarios
python3 run_socket_tests.py -s readonly      # Run a specific scenario
python3 run_socket_tests.py -s volume_filters -s socket_permissions
```

Typical output:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Scenario: readonly
  Read-only access: list commands allowed, mutations blocked.
  ➤ Container listing is permitted
    ✔ Container listing succeeds in read-only mode.
  ➤ Container creation is blocked
    ✔ Mutating container operations are blocked as expected.
  Summary: 2 passed, 0 failed
```

## Adding New Tests

1. Add a new entry to `SCENARIOS` or extend an existing one with a new dictionary item.
2. Provide any additional environment overrides under `environment`.
3. Append new `tests` dictionaries describing the command and expectations.
4. Re-run `python3 run_socket_tests.py` to validate your additions.

## Cleanup

The runner automatically executes `docker-compose ... down -v --remove-orphans` after each
scenario. If a run is interrupted, you can clean up manually:

```bash
cd tests/integration
docker-compose -f docker-compose.socket.yml down -v --remove-orphans
```

## Troubleshooting

- **docker-compose not found**: install the `docker-compose` plugin or binary.
- **Permission denied / privileged mode**: ensure Docker can run privileged containers (DinD).
- **Containers keep running after tests**: run the cleanup command above.

