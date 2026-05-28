#!/usr/bin/env python3
"""Inspect selected MySQL dump tables without printing huge INSERT lines.

This script is intentionally bounded: it streams a .sql.bz2 dump, extracts CREATE
TABLE blocks, counts INSERT statements, and samples post-type/status counts from
wp_posts INSERT tuples without echoing raw rows.
"""

from __future__ import annotations

import argparse
import bz2
import csv
import io
import re
from collections import Counter, defaultdict
from pathlib import Path

DEFAULT_TABLES = {
    "wp_posts",
    "wp_postmeta",
    "wp_terms",
    "wp_term_taxonomy",
    "wp_term_relationships",
    "search_products",
    "wp_wc_product_meta_lookup",
}

CREATE_RE = re.compile(r"^CREATE TABLE `([^`]+)`")
INSERT_RE = re.compile(r"^INSERT INTO `([^`]+)`")


def split_insert_tuples(values_sql: str):
    """Yield tuple payload strings from an INSERT VALUES suffix.

    Handles quoted strings sufficiently for mysqldump output so we can pass each
    tuple through Python's csv parser.
    """
    depth = 0
    in_quote = False
    escape = False
    start = None
    for i, ch in enumerate(values_sql):
        if in_quote:
            if escape:
                escape = False
            elif ch == "\\":
                escape = True
            elif ch == "'":
                in_quote = False
            continue
        if ch == "'":
            in_quote = True
        elif ch == "(":
            if depth == 0:
                start = i + 1
            depth += 1
        elif ch == ")":
            depth -= 1
            if depth == 0 and start is not None:
                yield values_sql[start:i]
                start = None


def parse_tuple(payload: str) -> list[str]:
    # MySQL string quoting is close enough to CSV when using single quote and backslash escape.
    reader = csv.reader(io.StringIO(payload), delimiter=",", quotechar="'", escapechar="\\")
    return next(reader)


def inspect_dump(path: Path, tables: set[str], sample_limit: int) -> None:
    create_blocks: dict[str, list[str]] = {}
    insert_counts: Counter[str] = Counter()
    wp_post_type_status: Counter[tuple[str, str]] = Counter()
    wp_post_samples: defaultdict[str, list[tuple[str, str, str, str]]] = defaultdict(list)

    current_table = None
    current_block: list[str] = []

    with bz2.open(path, "rt", encoding="utf-8", errors="replace") as f:
        for line in f:
            m = CREATE_RE.match(line)
            if m:
                current_table = m.group(1)
                current_block = [line.rstrip("\n")]
                continue
            if current_table:
                current_block.append(line.rstrip("\n"))
                if line.startswith(") ENGINE"):
                    if current_table in tables:
                        create_blocks[current_table] = current_block
                    current_table = None
                    current_block = []
                continue

            m = INSERT_RE.match(line)
            if not m:
                continue
            table = m.group(1)
            insert_counts[table] += 1
            if table != "wp_posts":
                continue

            values_idx = line.find(" VALUES ")
            if values_idx < 0:
                continue
            for payload in split_insert_tuples(line[values_idx + len(" VALUES ") :].rstrip(";\n")):
                try:
                    fields = parse_tuple(payload)
                except Exception:
                    continue
                if len(fields) < 21:
                    continue
                post_id = fields[0]
                title = fields[5]
                status = fields[7]
                slug = fields[11]
                post_type = fields[20]
                wp_post_type_status[(post_type, status)] += 1
                samples = wp_post_samples[post_type]
                if len(samples) < sample_limit:
                    samples.append((post_id, status, slug, title[:120]))

    print("# Dump inspection")
    print(f"path: {path}")
    print(f"size_bytes: {path.stat().st_size}")
    print("\n## Selected CREATE TABLE blocks")
    for table in sorted(tables):
        block = create_blocks.get(table)
        if not block:
            print(f"\n### {table}\nNOT FOUND")
            continue
        print(f"\n### {table}")
        for line in block[:80]:
            print(line)
        if len(block) > 80:
            print(f"... truncated {len(block)-80} lines ...")

    print("\n## INSERT statement counts")
    for table, count in sorted(insert_counts.items()):
        if table in tables or table.startswith("wp_"):
            print(f"{table}\t{count}")

    print("\n## wp_posts post_type/status tuple counts")
    for (post_type, status), count in sorted(wp_post_type_status.items()):
        print(f"{post_type}\t{status}\t{count}")

    print("\n## wp_posts samples by post_type")
    for post_type in sorted(wp_post_samples):
        print(f"\n### {post_type}")
        for post_id, status, slug, title in wp_post_samples[post_type]:
            print(f"{post_id}\t{status}\t{slug}\t{title}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("dump", nargs="?", default="/home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2")
    parser.add_argument("--tables", default=",".join(sorted(DEFAULT_TABLES)))
    parser.add_argument("--sample-limit", type=int, default=5)
    args = parser.parse_args()
    tables = {t.strip() for t in args.tables.split(",") if t.strip()}
    inspect_dump(Path(args.dump), tables, args.sample_limit)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
