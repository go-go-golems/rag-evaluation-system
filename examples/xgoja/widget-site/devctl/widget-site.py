#!/usr/bin/env python3
"""devctl plugin for examples/xgoja/widget-site.

This plugin is intentionally scoped to the standalone generated xgoja example.
It lets a developer run the usual loop from the example directory:

  devctl plugins list
  devctl plan
  devctl up --force --timeout 10m
  devctl status
  devctl logs --service widget-site --stderr --follow
  devctl down

The plugin computes a local port, validates the sibling go-go-goja checkout and
RAG app assets, rebuilds the generated binary, and lets devctl supervise the
resulting HTTP server.
"""
from __future__ import annotations

import json
import os
from pathlib import Path
import shutil
import socket
import subprocess
import sys
import time
from typing import Any, Dict, List, Optional

PLUGIN_NAME = "rag-widget-xgoja-site"
DEFAULT_PORT = 18791


def emit(obj: Dict[str, Any]) -> None:
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()


def log(msg: str) -> None:
    sys.stderr.write(f"[{PLUGIN_NAME}] {msg}\n")
    sys.stderr.flush()


def unsupported(request_id: str, op: str) -> None:
    emit({
        "type": "response",
        "request_id": request_id,
        "ok": False,
        "error": {"code": "E_UNSUPPORTED", "message": f"unsupported op: {op}"},
    })


def find_free_port(preferred: int) -> int:
    def available(port: int) -> bool:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        try:
            sock.bind(("127.0.0.1", port))
            return True
        except OSError:
            return False
        finally:
            sock.close()

    if available(preferred):
        return preferred

    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.bind(("127.0.0.1", 0))
    port = sock.getsockname()[1]
    sock.close()
    return int(port)


def cfg_get(config: Dict[str, Any], dotted: str, default: Any = None) -> Any:
    node: Any = config
    for part in dotted.split("."):
        if isinstance(node, dict) and part in node:
            node = node[part]
        else:
            return default
    return node


def example_root(ctx: Dict[str, Any]) -> Path:
    return Path(ctx.get("repo_root") or os.getcwd()).resolve()


def rag_root(root: Path) -> Path:
    return (root / "../../..").resolve()


def xgoja_root(root: Path) -> Path:
    return (root / "../../../../go-go-goja").resolve()


def run_step(name: str, command: List[str], cwd: Path, timeout_s: int, dry_run: bool) -> Dict[str, Any]:
    log(f"build step {name}: {' '.join(command)} (cwd={cwd})")
    if dry_run:
        return {"name": name, "ok": True, "status": "skipped", "output": "dry-run"}

    started = time.time()
    try:
        proc = subprocess.run(
            command,
            cwd=str(cwd),
            text=True,
            capture_output=True,
            timeout=timeout_s,
        )
    except subprocess.TimeoutExpired:
        return {"name": name, "ok": False, "status": "failed", "output": f"timeout after {timeout_s}s"}

    elapsed = time.time() - started
    output = (proc.stderr or proc.stdout or "")[-4000:]
    ok = proc.returncode == 0
    status = "succeeded" if ok else "failed"
    if proc.returncode != 0:
        log(f"step {name} failed with exit code {proc.returncode}: {output[-500:]}")
    return {"name": name, "ok": ok, "status": status, "output": output, "duration_ms": int(elapsed * 1000)}


emit({
    "type": "handshake",
    "protocol_version": "v2",
    "plugin_name": PLUGIN_NAME,
    "capabilities": {
        "ops": ["config.mutate", "validate.run", "build.run", "launch.plan", "command.run"],
        "commands": [
            {"name": "sync-app", "help": "Copy packages/rag-evaluation-site/app-dist into embedded xgoja assets"},
            {"name": "smoke", "help": "Run the example Makefile smoke, including action endpoint assertions"},
            {"name": "clean", "help": "Remove the generated xgoja binary dist directory"},
        ],
    },
})

for raw in sys.stdin:
    raw = raw.strip()
    if not raw:
        continue

    req = json.loads(raw)
    request_id = req.get("request_id", "")
    op = req.get("op", "")
    ctx = req.get("ctx", {}) or {}
    inp = req.get("input", {}) or {}
    config = inp.get("config", {}) or {}
    dry_run = bool(ctx.get("dry_run", False))

    root = example_root(ctx)
    repo = rag_root(root)
    goja = xgoja_root(root)
    bin_path = root / "dist" / "rag-widget-xgoja-site"

    try:
        if op == "config.mutate":
            port = find_free_port(DEFAULT_PORT)
            url = f"http://127.0.0.1:{port}"
            emit({
                "type": "response",
                "request_id": request_id,
                "ok": True,
                "output": {
                    "config_patch": {
                        "set": {
                            "services.widgetSite.port": port,
                            "services.widgetSite.url": url,
                            "services.widgetSite.healthUrl": f"{url}/healthz",
                            "services.widgetSite.actionsUrl": f"{url}/pages/actions",
                            "artifacts.binary": str(bin_path.relative_to(root)),
                            "paths.ragRoot": str(repo),
                            "paths.xgojaRoot": str(goja),
                            "env.HTTP_ADDR": f"127.0.0.1:{port}",
                        },
                        "unset": [],
                    }
                },
            })

        elif op == "validate.run":
            errors: List[Dict[str, str]] = []
            warnings: List[Dict[str, str]] = []

            for tool in ["go", "make", "curl", "python3"]:
                if not shutil.which(tool):
                    errors.append({"code": "E_MISSING_TOOL", "message": f"{tool} not found in PATH"})

            if not (root / "xgoja.yaml").is_file():
                errors.append({"code": "E_MISSING_FILE", "message": "xgoja.yaml not found; run devctl from examples/xgoja/widget-site or pass --repo-root"})
            if not (root / "Makefile").is_file():
                errors.append({"code": "E_MISSING_FILE", "message": "Makefile not found in widget-site example"})
            if not (root / "verbs" / "sites.js").is_file():
                errors.append({"code": "E_MISSING_FILE", "message": "verbs/sites.js not found"})
            if not (repo / "packages" / "rag-evaluation-site" / "app-dist" / "index.html").is_file():
                warnings.append({"code": "W_MISSING_APP_DIST", "message": "packages/rag-evaluation-site/app-dist/index.html not found; run the package app build before sync-app"})
            if not (root / "assets" / "public" / "index.html").is_file():
                warnings.append({"code": "W_MISSING_SYNCED_APP", "message": "assets/public/index.html not found; build.run will execute make sync-app"})
            if not (goja / "cmd" / "xgoja").is_dir():
                errors.append({"code": "E_MISSING_XGOJA", "message": f"xgoja checkout not found at {goja}"})
            if not bin_path.exists():
                warnings.append({"code": "W_BINARY_NOT_BUILT", "message": "dist/rag-widget-xgoja-site is not built yet; devctl build or devctl up will build it"})

            emit({
                "type": "response",
                "request_id": request_id,
                "ok": True,
                "output": {"valid": len(errors) == 0, "errors": errors, "warnings": warnings},
            })

        elif op == "build.run":
            deadline_ms = int(ctx.get("deadline_ms") or 0)
            timeout_s = 600 if deadline_ms <= 0 else max(30, min(600, int(deadline_ms / 1000) - 2))
            steps = [
                ("sync-app", ["make", "sync-app"], root),
                ("build", ["make", "build"], root),
            ]
            results = [run_step(name, cmd, cwd, timeout_s, dry_run) for name, cmd, cwd in steps]
            all_ok = all(result["status"] in ("succeeded", "skipped") for result in results)
            emit({
                "type": "response",
                "request_id": request_id,
                "ok": all_ok,
                "output": {
                    "steps": results,
                    "artifacts": {
                        "binary": str(bin_path.relative_to(root)),
                        "embeddedAssets": "assets/public",
                    },
                },
                **({} if all_ok else {"error": {"code": "E_BUILD_FAILED", "message": "one or more widget-site build steps failed"}}),
            })

        elif op == "launch.plan":
            port = int(cfg_get(config, "services.widgetSite.port", DEFAULT_PORT))
            addr = f"127.0.0.1:{port}"
            url = f"http://{addr}"
            emit({
                "type": "response",
                "request_id": request_id,
                "ok": True,
                "output": {
                    "services": [
                        {
                            "name": "widget-site",
                            "cwd": ".",
                            "command": [str(bin_path.relative_to(root)), "serve", "sites", "demo", "--http-listen", addr],
                            "env": {"HTTP_ADDR": addr},
                            "health": {"type": "http", "url": f"{url}/healthz", "timeout_ms": 30000},
                        }
                    ]
                },
            })

        elif op == "command.run":
            name = inp.get("name", "")
            command_map = {
                "sync-app": ["make", "sync-app"],
                "smoke": ["make", "smoke"],
                "clean": ["make", "clean"],
            }
            if name not in command_map:
                unsupported(request_id, f"command.run:{name}")
                continue
            if dry_run:
                log(f"dry-run: would execute {' '.join(command_map[name])}")
                emit({"type": "response", "request_id": request_id, "ok": True, "output": {"exit_code": 0}})
                continue
            proc = subprocess.run(command_map[name], cwd=str(root), stdout=sys.stderr, stderr=sys.stderr)
            emit({"type": "response", "request_id": request_id, "ok": proc.returncode == 0, "output": {"exit_code": proc.returncode}})

        else:
            unsupported(request_id, op)

    except Exception as exc:  # keep protocol errors explicit
        log(f"ERROR: {exc}")
        emit({
            "type": "response",
            "request_id": request_id,
            "ok": False,
            "error": {"code": "E_PLUGIN", "message": str(exc)},
        })
