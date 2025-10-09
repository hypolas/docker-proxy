#!/usr/bin/env python3
"""Integration test runner for docker-proxy in Unix socket mode."""

from __future__ import annotations

import argparse
import os
import re
import subprocess
import sys
import time
import uuid
from pathlib import Path
from typing import Dict, Iterable, List, Tuple

from scenarios import SCENARIOS

BASE_DIR = Path(__file__).resolve().parent
COMPOSE_FILE = BASE_DIR / "docker-compose.socket.yml"
DOCKER_COMPOSE_CMD = ["docker-compose"]


class Colors:
    BLUE = "\033[0;34m"
    GREEN = "\033[0;32m"
    RED = "\033[0;31m"
    YELLOW = "\033[1;33m"
    RESET = "\033[0m"


def compose_cmd(*args: str) -> List[str]:
    return DOCKER_COMPOSE_CMD + ["-f", str(COMPOSE_FILE)] + list(args)


def run_subprocess(
    args: Iterable[str],
    *,
    env: Dict[str, str],
    cwd: Path = BASE_DIR,
    capture_output: bool = True,
    check: bool = False,
) -> Tuple[int, str]:
    proc = subprocess.run(  # noqa: S603, S607 - docker compose command
        list(args),
        cwd=str(cwd),
        env=env,
        capture_output=capture_output,
        text=True,
    )
    if check:
        proc.check_returncode()
    output = (proc.stdout or "") + (proc.stderr or "")
    return proc.returncode, output.strip()


def exec_in_client(command: str, *, env: Dict[str, str], project: str) -> Tuple[int, str]:
    args = compose_cmd("-p", project, "exec", "-T", "test-client", "sh", "-lc", command)
    return run_subprocess(args, env=env)


def bring_up(env: Dict[str, str], project: str) -> None:
    rc, output = run_subprocess(
        compose_cmd("-p", project, "up", "-d"),
        env=env,
    )
    if rc != 0:
        raise RuntimeError(f"docker-compose up failed:\n{output}")
    time.sleep(8)


def bring_down(env: Dict[str, str], project: str) -> None:
    run_subprocess(
        compose_cmd("-p", project, "down", "-v", "--remove-orphans"),
        env=env,
        capture_output=False,
    )


def contains_any(haystack: str, needles: Iterable[str]) -> bool:
    text = haystack.lower()
    return any(needle.lower() in text for needle in needles)


def contains_all(haystack: str, needles: Iterable[str]) -> bool:
    text = haystack.lower()
    return all(needle.lower() in text for needle in needles)


def regex_any(haystack: str, patterns: Iterable[str]) -> bool:
    return any(re.search(pattern, haystack, flags=re.MULTILINE) for pattern in patterns)


def regex_all(haystack: str, patterns: Iterable[str]) -> bool:
    return all(re.search(pattern, haystack, flags=re.MULTILINE) for pattern in patterns)


def evaluate(expect: Dict[str, List[str] | List[int]], output: str, returncode: int) -> bool:
    success = True
    if "returncodes" in expect:
        success &= returncode in expect["returncodes"]  # type: ignore[index]
    if "contains_any" in expect:
        success &= contains_any(output, expect["contains_any"])  # type: ignore[index]
    if "contains_all" in expect:
        success &= contains_all(output, expect["contains_all"])  # type: ignore[index]
    if "regex_any" in expect:
        success &= regex_any(output, expect["regex_any"])  # type: ignore[index]
    if "regex_all" in expect:
        success &= regex_all(output, expect["regex_all"])  # type: ignore[index]
    if "not_contains" in expect:
        success &= not contains_any(output, expect["not_contains"])  # type: ignore[index]
    if "not_regex" in expect:
        success &= not regex_any(output, expect["not_regex"])  # type: ignore[index]
    return success


def run_test_case(case: Dict[str, object], *, env: Dict[str, str], project: str) -> bool:
    case_name = case["name"]  # type: ignore[index]
    print(f"  {Colors.YELLOW}➤{Colors.RESET} {case_name}")

    command = case["command"]  # type: ignore[index]
    rc, output = exec_in_client(command, env=env, project=project)
    expect = case.get("expect", {})  # type: ignore[assignment]
    success = evaluate(expect, output, rc)

    if success:
        message = case.get("messages", {}).get("success", "Test passed.")  # type: ignore[call-arg, index]
        print(f"    {Colors.GREEN}✔{Colors.RESET} {message}")
    else:
        message = case.get("messages", {}).get("failure", "Test failed.")  # type: ignore[call-arg, index]
        print(f"    {Colors.RED}✖{Colors.RESET} {message}")
        print("    --- command output ---")
        for line in output.splitlines():
            print(f"    {line}")
        print("    ----------------------")

    cleanup_commands = case.get("cleanup", [])  # type: ignore[assignment]
    for clean_cmd in cleanup_commands:
        exec_in_client(clean_cmd, env=env, project=project)

    return success


def run_scenario(name: str, scenario: Dict[str, object], *, project: str) -> bool:
    print(f"{Colors.BLUE}{'━'*70}{Colors.RESET}")
    print(f"{Colors.YELLOW}Scenario:{Colors.RESET} {name}")
    print(f"  {scenario.get('description', '').strip()}")

    env = os.environ.copy()
    env.update({k: str(v) for k, v in scenario.get("environment", {}).items()})

    try:
        bring_up(env, project)
        results = []
        for case in scenario.get("tests", []):  # type: ignore[assignment]
            results.append(run_test_case(case, env=env, project=project))
        passed = sum(1 for result in results if result)
        failed = len(results) - passed
        print(f"  {Colors.YELLOW}Summary:{Colors.RESET} {passed} passed, {failed} failed")
        return failed == 0
    finally:
        bring_down(env, project)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run docker-proxy socket integration tests")
    parser.add_argument(
        "--scenario",
        "-s",
        action="append",
        help="Scenario name to run (default: all). Can be provided multiple times.",
    )
    parser.add_argument(
        "--list",
        action="store_true",
        help="List available scenarios and exit.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()

    if args.list:
        for name, scenario in SCENARIOS.items():
            print(f"{name}: {scenario.get('description', '').strip()}")
        return 0

    scenario_names = args.scenario or list(SCENARIOS.keys())
    exit_code = 0
    for name in scenario_names:
        if name not in SCENARIOS:
            print(f"{Colors.RED}Unknown scenario:{Colors.RESET} {name}")
            exit_code = 1
            continue
        project = f"dkproxy_socket_{name}_{uuid.uuid4().hex[:6]}"
        success = run_scenario(name, SCENARIOS[name], project=project)
        if not success:
            exit_code = 1
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
