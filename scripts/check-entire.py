#!/usr/bin/env python3
"""Read Entire prerequisites without invoking Git or changing repository state.

This check deliberately does not claim that a session or checkpoint was captured.
Use Entire's session and checkpoint diagnostics to verify the caller session.
"""

import json
import re
import shutil
import sys
from pathlib import Path


def check() -> int:
    root = Path(__file__).resolve().parent.parent
    results = []

    def record(name: str, passed: bool, detail: str) -> None:
        results.append({"name": name, "passed": passed, "detail": detail})

    record("entire_cli", shutil.which("entire") is not None, "CLI must be on PATH")
    git_dir = root / ".git"
    if git_dir.is_file():
        pointer = git_dir.read_text().strip()
        if pointer.startswith("gitdir: "):
            git_dir = (root / pointer[8:]).resolve()
    common_dir = git_dir
    if (git_dir / "commondir").is_file():
        common_dir = (git_dir / (git_dir / "commondir").read_text().strip()).resolve()

    head = (git_dir / "HEAD").read_text().strip() if (git_dir / "HEAD").is_file() else ""
    commit = head
    if head.startswith("ref: "):
        ref = head[5:]
        loose = common_dir / ref
        commit = loose.read_text().strip() if loose.is_file() else ""
        packed = common_dir / "packed-refs"
        if not commit and packed.is_file():
            for line in packed.read_text().splitlines():
                if line.endswith(" " + ref):
                    commit = line.split(" ", 1)[0]
                    break
    has_head = bool(re.fullmatch(r"[0-9a-f]{40}|[0-9a-f]{64}", commit))
    record("head_reference", has_head, commit if has_head else "No anchor commit; Entire cannot initialize capture")

    settings = {}
    for name in ("settings.json", "settings.local.json"):
        path = root / ".entire" / name
        if path.is_file():
            settings.update(json.loads(path.read_text()))
    record("entire_enabled", settings.get("enabled") is True, "Effective repository setting")

    path = root / ".codex/hooks.json"
    hooks = json.loads(path.read_text()).get("hooks", {}) if path.is_file() else {}
    for event, command in (
        ("SessionStart", "session-start"),
        ("UserPromptSubmit", "user-prompt-submit"),
        ("PostToolUse", "post-tool-use"),
        ("Stop", "stop"),
        ("SessionEnd", "session-end"),
    ):
        present = any(
            "entire hooks codex " + command in hook.get("command", "")
            for group in hooks.get(event, [])
            for hook in group.get("hooks", [])
        )
        record("codex_" + event, present, "Entire hook configured")

    for name in ("prepare-commit-msg", "commit-msg", "post-commit"):
        path = common_dir / "hooks" / name
        present = path.is_file() and "entire hooks git " + name in path.read_text()
        record("git_hook_" + name, present, "Default Git hook configured; custom core.hooksPath is not checked")

    passed = all(item["passed"] for item in results)
    print(json.dumps({
        "prerequisites": "pass" if passed else "blocked",
        "capture_verified": None,
        "checkpoint_verified": None,
        "checks": results,
        "next": "Inspect current caller session and checkpoint diagnostics; this preflight does not verify capture" if passed else "Resolve failing prerequisites before application features",
    }, indent=2))
    return 0 if passed else 1


if __name__ == "__main__":
    try:
        sys.exit(check())
    except (OSError, ValueError) as error:
        print(json.dumps({"prerequisites": "error", "reason": type(error).__name__}))
        sys.exit(2)
