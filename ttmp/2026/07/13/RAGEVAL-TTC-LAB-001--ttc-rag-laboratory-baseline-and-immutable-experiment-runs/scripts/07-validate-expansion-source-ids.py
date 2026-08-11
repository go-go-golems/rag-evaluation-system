#!/usr/bin/env python3
"""Read-only validation of proposed TTC expansion source IDs.

This script proves only source identity: every referenced ``wp:*`` ID exists
in the rebuilt SQLite export. It does not assign relevance, inspect exact
evidence spans, resolve revisions, or freeze a dataset.
"""

from __future__ import annotations

import argparse
import re
import sqlite3
import sys
from pathlib import Path


SOURCE_ID = re.compile(r"\bwp:\d+\b")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--db", default="data/ttc-wordpress-rag.sqlite")
    parser.add_argument("cards", nargs="+", type=Path)
    args = parser.parse_args()

    source_ids: set[str] = set()
    for path in args.cards:
        text = path.read_text(encoding="utf-8")
        source_ids.update(SOURCE_ID.findall(text))

    with sqlite3.connect(args.db) as database:
        placeholders = ",".join("?" for _ in source_ids)
        rows = database.execute(
            f"SELECT doc_id, kind, title FROM documents WHERE doc_id IN ({placeholders})",
            sorted(source_ids),
        ).fetchall()

    found = {row[0]: row for row in rows}
    missing = sorted(source_ids - found.keys())
    print(f"files={len(args.cards)} unique_source_ids={len(source_ids)}")
    print(f"resolved={len(found)} missing={len(missing)}")
    for source_id in sorted(found):
        _, kind, title = found[source_id]
        print(f"FOUND {source_id}\t{kind}\t{title}")
    for source_id in missing:
        print(f"MISSING {source_id}", file=sys.stderr)

    return 1 if missing else 0


if __name__ == "__main__":
    raise SystemExit(main())
