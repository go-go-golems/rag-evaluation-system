#!/usr/bin/env python3
"""Report legacy Widget DSL module usage in JavaScript/TypeScript sources.

This is a migration helper for RAGEVAL-SCHEDULE-WIDGETS Phase 11. It does not
rewrite files; it gives reviewers a focused list of scripts that still need the
legacy split modules selected in xgoja runtime config.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import Iterable

LEGACY_MODULES = [
    "ui.dsl",
    "data.dsl",
    "data.v2.dsl",
    "context_window.dsl",
    "course.dsl",
    "cms.dsl",
]

# Raw component usage can exist in v3 too, but every occurrence should be an
# explicit migration exception instead of a silent compatibility crutch.
RAW_COMPONENT_RE = re.compile(r"\b(?:raw|widget\.raw)\.component\s*\(")

SOURCE_SUFFIXES = {".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"}
DEFAULT_IGNORED_DIRS = {
    ".git",
    "node_modules",
    "dist",
    "app-dist",
    "storybook-static",
    "coverage",
    ".next",
    ".turbo",
}


@dataclass(frozen=True)
class Finding:
    path: str
    line: int
    kind: str
    value: str
    text: str


def iter_source_files(paths: Iterable[Path]) -> Iterable[Path]:
    for path in paths:
        if not path.exists():
            continue
        if path.is_file():
            if path.suffix in SOURCE_SUFFIXES:
                yield path
            continue
        for child in path.rglob("*"):
            if child.is_dir():
                continue
            if child.suffix not in SOURCE_SUFFIXES:
                continue
            if any(part in DEFAULT_IGNORED_DIRS for part in child.parts):
                continue
            yield child


def module_patterns(module_name: str) -> list[re.Pattern[str]]:
    escaped = re.escape(module_name)
    return [
        re.compile(rf"\brequire\(\s*['\"]{escaped}['\"]\s*\)"),
        re.compile(rf"\bfrom\s+['\"]{escaped}['\"]"),
        re.compile(rf"\bimport\(\s*['\"]{escaped}['\"]\s*\)"),
    ]


MODULE_PATTERNS = {
    module_name: module_patterns(module_name) for module_name in LEGACY_MODULES
}


def scan_file(path: Path, root: Path) -> list[Finding]:
    findings: list[Finding] = []
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except UnicodeDecodeError:
        return findings

    display_path = str(path.relative_to(root)) if path.is_relative_to(root) else str(path)
    for line_no, line in enumerate(lines, start=1):
        stripped = line.strip()
        for module_name, patterns in MODULE_PATTERNS.items():
            if any(pattern.search(line) for pattern in patterns):
                findings.append(
                    Finding(
                        path=display_path,
                        line=line_no,
                        kind="legacy-module-import",
                        value=module_name,
                        text=stripped,
                    )
                )
        if RAW_COMPONENT_RE.search(line):
            findings.append(
                Finding(
                    path=display_path,
                    line=line_no,
                    kind="raw-component-escape-hatch",
                    value="raw.component",
                    text=stripped,
                )
            )
    return findings


def default_paths(root: Path) -> list[Path]:
    candidates = [
        root / "go-go-course" / "cmd" / "go-go-course" / "lib" / "pages",
        root / "pkg" / "widgetdsl" / "testdata" / "v3" / "examples",
        root / "examples",
    ]
    return [path for path in candidates if path.exists()]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "paths",
        nargs="*",
        type=Path,
        help="Files/directories to scan. Defaults to known first-party widget script locations when present.",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Emit findings as JSON instead of a human-readable report.",
    )
    parser.add_argument(
        "--fail-on-findings",
        action="store_true",
        help="Exit 1 when any finding is reported; useful for migration gates after a host claims v3-only status.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    root = Path.cwd().resolve()
    paths = [path.resolve() for path in args.paths] if args.paths else default_paths(root)

    findings: list[Finding] = []
    for source in sorted(set(iter_source_files(paths))):
        findings.extend(scan_file(source.resolve(), root))

    findings.sort(key=lambda item: (item.path, item.line, item.kind, item.value))

    if args.json:
        print(json.dumps([asdict(finding) for finding in findings], indent=2, sort_keys=True))
    else:
        if not findings:
            print("No legacy Widget DSL imports or raw component escape hatches found.")
        else:
            print(f"Found {len(findings)} migration finding(s):")
            for finding in findings:
                print(
                    f"{finding.path}:{finding.line}: {finding.kind} "
                    f"{finding.value}: {finding.text}"
                )

    if findings and args.fail_on_findings:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
