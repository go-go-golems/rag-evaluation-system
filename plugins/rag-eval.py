#!/usr/bin/env python3
"""devctl plugin for the RAG Evaluation System.

Runs two services:
  - backend: Go API server (free port, health-checked)
  - web:     Vite dev server (free port, HMR + API proxy to backend)

Port allocation:
  config.mutate finds two free ports and publishes them as
  services.backend.port / services.web.port.  All downstream ops
  (build.run, launch.plan) read from the merged config so the
  same ports are used everywhere.

Usage:
  devctl plan
  devctl up
  devctl status
  devctl logs --service backend --follow
  devctl logs --service web --follow
  devctl down
"""
import json
import os
import shutil
import socket
import subprocess
import sys

# ─── Port allocation ───────────────────────────────────────────────────────

def find_free_port(preferred: int) -> int:
    """Return a free TCP port on 127.0.0.1.

    Tries the preferred port first; if taken, asks the OS for any
    free port.  The socket is opened and closed so the port is
    available by the time the caller uses it (there is a small
    TOCTOU window, but devctl runs sequentially).
    """
    def _port_free(p: int) -> bool:
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            s.bind(("127.0.0.1", p))
            s.close()
            return True
        except OSError:
            return False

    if _port_free(preferred):
        return preferred

    # Ask the OS for any free port
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind(("127.0.0.1", 0))
    port = s.getsockname()[1]
    s.close()
    return port


# ─── Protocol helpers ──────────────────────────────────────────────────────

def emit(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()

def log(msg):
    sys.stderr.write(f"[rag-eval] {msg}\n")
    sys.stderr.flush()

def unsupported(rid, op):
    emit({
        "type": "response",
        "request_id": rid,
        "ok": False,
        "error": {"code": "E_UNSUPPORTED", "message": f"unsupported op: {op}"},
    })

# ─── Handshake ─────────────────────────────────────────────────────────────

emit({
    "type": "handshake",
    "protocol_version": "v2",
    "plugin_name": "rag-eval",
    "capabilities": {
        "ops": [
            "config.mutate",
            "validate.run",
            "build.run",
            "launch.plan",
            "command.run",
        ],
        "commands": [
            {"name": "build", "help": "Build Go binary with embedded SPA"},
            {"name": "web-build", "help": "Build Vite frontend only"},
        ],
    },
})

# ─── Request loop ──────────────────────────────────────────────────────────

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue

    req = json.loads(line)
    rid = req.get("request_id", "")
    op = req.get("op", "")
    ctx = req.get("ctx", {}) or {}
    inp = req.get("input", {}) or {}

    repo_root = ctx.get("repo_root", "")
    dry_run = bool(ctx.get("dry_run", False))
    # config is the merged config from config.mutate — available in all ops
    config = inp.get("config", {}) or {}

    # Helper: read a nested config value using dotted key path
    def cfg_get(key: str, default=None):
        """Traverse nested dict by dotted key like 'services.backend.port'."""
        parts = key.split('.')
        node = config
        for p in parts:
            if isinstance(node, dict) and p in node:
                node = node[p]
            else:
                return default
        return node

    def cfg_port(key: str, default: int) -> int:
        v = cfg_get(key)
        if v is not None:
            return int(v)
        return default

    try:
        # ─── config.mutate ────────────────────────────────────────────
        if op == "config.mutate":
            backend_port = find_free_port(8772)
            vite_port = find_free_port(5173)
            backend_url = f"http://127.0.0.1:{backend_port}"

            log(f"allocated ports: backend={backend_port}, web={vite_port}")

            emit({
                "type": "response",
                "request_id": rid,
                "ok": True,
                "output": {
                    "config_patch": {
                        "set": {
                            "services.backend.port": backend_port,
                            "services.backend.url": backend_url,
                            "services.web.port": vite_port,
                            "env.RAG_EVAL_ADDRESS": f"127.0.0.1:{backend_port}",
                            "env.RAG_EVAL_LOG_LEVEL": "debug",
                        },
                        "unset": [],
                    },
                },
            })

        # ─── validate.run ─────────────────────────────────────────────
        elif op == "validate.run":
            errors = []
            warnings = []

            if not shutil.which("go"):
                errors.append({
                    "code": "E_MISSING_TOOL",
                    "message": "go not found in PATH",
                    "hint": "Install Go 1.22+ from https://go.dev/dl/",
                })

            if not shutil.which("pnpm"):
                errors.append({
                    "code": "E_MISSING_TOOL",
                    "message": "pnpm not found in PATH",
                    "hint": "Install pnpm: npm install -g pnpm",
                })

            web_dir = os.path.join(repo_root, "web")
            if not os.path.isdir(os.path.join(web_dir, "node_modules")):
                warnings.append({
                    "code": "W_MISSING_DEPS",
                    "message": "web/node_modules not found — run pnpm install",
                    "hint": "cd web && pnpm install",
                })

            state_dir = os.path.join(repo_root, "state")
            if not os.path.isdir(state_dir):
                warnings.append({
                    "code": "W_NO_STATE_DIR",
                    "message": "state/ directory does not exist — will be created on first run",
                })

            emit({
                "type": "response",
                "request_id": rid,
                "ok": True,
                "output": {
                    "valid": len(errors) == 0,
                    "errors": errors,
                    "warnings": warnings,
                },
            })

        # ─── build.run ────────────────────────────────────────────────
        elif op == "build.run":
            steps = []
            web_dir = os.path.join(repo_root, "web")

            if os.path.isdir(web_dir) and not os.path.isdir(os.path.join(web_dir, "node_modules")):
                steps.append({
                    "name": "pnpm-install",
                    "description": "Install frontend dependencies",
                    "command": ["pnpm", "install"],
                    "cwd": "web",
                })

            steps.append({
                "name": "web-build",
                "description": "Build Vite SPA for embedding",
                "command": ["pnpm", "build"],
                "cwd": "web",
            })

            steps.append({
                "name": "go-build",
                "description": "Build rag-eval binary with embedded frontend",
                "command": ["go", "build", "./cmd/rag-eval"],
                "cwd": ".",
            })

            if dry_run:
                log(f"dry-run: would run {len(steps)} build steps")
                for s in steps:
                    log(f"  - {s['name']}: {' '.join(s['command'])} (cwd={s['cwd']})")

            results = []
            for step in steps:
                if dry_run:
                    results.append({"name": step["name"], "status": "skipped", "output": ""})
                    continue

                step_cwd = os.path.join(repo_root, step["cwd"])
                log(f"build step: {step['name']}")
                try:
                    proc = subprocess.run(
                        step["command"],
                        cwd=step_cwd,
                        capture_output=True,
                        text=True,
                        timeout=120,
                    )
                    status = "succeeded" if proc.returncode == 0 else "failed"
                    output = proc.stderr[-500:] if proc.stderr else proc.stdout[-500:]
                    if proc.returncode != 0:
                        log(f"  FAILED: {output[:200]}")
                    results.append({"name": step["name"], "status": status, "output": output})
                except subprocess.TimeoutExpired:
                    results.append({"name": step["name"], "status": "failed", "output": "timeout after 120s"})
                    log(f"  TIMEOUT: {step['name']}")

            all_ok = all(r["status"] in ("succeeded", "skipped") for r in results)

            emit({
                "type": "response",
                "request_id": rid,
                "ok": all_ok,
                "output": {
                    "steps": results,
                    "artifacts": {
                        "rag-eval-bin": "rag-eval",
                        "web-dist": "internal/web/dist",
                    },
                } if all_ok else {},
                **({} if all_ok else {"error": {"code": "E_BUILD_FAILED", "message": "one or more build steps failed"}}),
            })

        # ─── launch.plan ──────────────────────────────────────────────
        elif op == "launch.plan":
            # Read ports from merged config (set by config.mutate)
            backend_port = cfg_port("services.backend.port", 8772)
            vite_port = cfg_port("services.web.port", 5173)
            backend_url = f"http://127.0.0.1:{backend_port}"

            if dry_run:
                log(f"dry-run: would launch backend on :{backend_port}, web on :{vite_port}")

            # Ensure state directory exists for engine DB
            state_dir = os.path.join(repo_root, "state")
            if not dry_run:
                os.makedirs(state_dir, exist_ok=True)

            log(f"launching backend=:{backend_port} web=:{vite_port}")

            emit({
                "type": "response",
                "request_id": rid,
                "ok": True,
                "output": {
                    "services": [
                        {
                            "name": "backend",
                            "command": [
                                "bash", "--noprofile", "--norc", "-lc",
                                f"mkdir -p state && exec ./rag-eval serve --address 127.0.0.1:{backend_port} --db state/rag-eval.db --engine-db state/rag-eval-workflows.db --log-level debug",
                            ],
                            "env": {
                                "RAG_EVAL_ADDRESS": f"127.0.0.1:{backend_port}",
                                "RAG_EVAL_DB": "state/rag-eval.db",
                                "RAG_EVAL_ENGINE_DB": "state/rag-eval-workflows.db",
                            },
                            "health": {
                                "type": "http",
                                "url": f"{backend_url}/api/v1/health",
                                "timeout_ms": 15000,
                            },
                        },
                        {
                            "name": "web",
                            "command": [
                                "bash", "--noprofile", "--norc", "-lc",
                                f"exec pnpm dev --port {vite_port}",
                            ],
                            "cwd": "web",
                            "env": {
                                "RAG_EVAL_BACKEND_PORT": str(backend_port),
                            },
                        },
                    ],
                },
            })

        # ─── command.run ──────────────────────────────────────────────
        elif op == "command.run":
            cmd_name = inp.get("name", "")
            cmd_argv = inp.get("argv", [])

            if cmd_name == "build":
                log("command: build (Go binary + embedded SPA)")
                if not dry_run:
                    web_dir = os.path.join(repo_root, "web")
                    proc = subprocess.run(["pnpm", "build"], cwd=web_dir, timeout=120)
                    if proc.returncode != 0:
                        emit({"type": "response", "request_id": rid, "ok": False,
                              "error": {"code": "E_BUILD_FAILED", "message": "pnpm build failed"}})
                        continue
                    proc = subprocess.run(["go", "build", "./cmd/rag-eval"], cwd=repo_root, timeout=120)
                    emit({"type": "response", "request_id": rid, "ok": proc.returncode == 0,
                          "output": {"exit_code": proc.returncode}})
                else:
                    emit({"type": "response", "request_id": rid, "ok": True, "output": {"exit_code": 0}})

            elif cmd_name == "web-build":
                log("command: web-build (Vite only)")
                if not dry_run:
                    web_dir = os.path.join(repo_root, "web")
                    proc = subprocess.run(["pnpm", "build"], cwd=web_dir, timeout=120)
                    emit({"type": "response", "request_id": rid, "ok": proc.returncode == 0,
                          "output": {"exit_code": proc.returncode}})
                else:
                    emit({"type": "response", "request_id": rid, "ok": True, "output": {"exit_code": 0}})

            else:
                unsupported(rid, f"command.run:{cmd_name}")

        else:
            unsupported(rid, op)

    except Exception as e:
        log(f"ERROR: {e}")
        emit({
            "type": "response",
            "request_id": rid,
            "ok": False,
            "error": {"code": "E_PLUGIN", "message": str(e)},
        })
