#!/usr/bin/env python3
"""Read-only audit for TTC evaluation-card drafts and future split metadata.

Usage:
  python3 scripts/08-audit-ttc-evaluation-cards.py \
    --cards reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md \
    --cards reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md

  # Add a future machine-readable manifest to enforce split/family checks.
  python3 scripts/08-audit-ttc-evaluation-cards.py \
    --cards reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md \
    --metadata path/to/ttc-baseline-eval-v2-draft.json

The script never writes files and never promotes draft evidence to frozen
labels. It accepts source-document IDs in either `expected_source_ids` or
`source_document_ids`; it does not adjudicate their relevance.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from collections import Counter, defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Any


LIST_ID_RE = re.compile(r"^\s*-\s+\{?\s*id:\s*['\"]?([A-Za-z0-9][A-Za-z0-9._-]*)", re.MULTILINE)
QUERY_RE = re.compile(r"\bquery:\s*['\"]([^'\"]+)['\"]")
SOURCE_LIST_RE = re.compile(r"\b(?:expected_source_ids|source_document_ids):\s*\[([^]]*)\]")
ANSWERABILITY_RE = re.compile(r"\banswerability:\s*([A-Za-z0-9_-]+)")
WP_RE = re.compile(r"\bwp:\d+\b")
HEADING_CARD_RE = re.compile(r"^####\s+`([^`]+)`", re.MULTILINE)


@dataclass(frozen=True)
class Card:
    card_id: str
    query: str | None
    source_ids: tuple[str, ...]
    answerability: str | None
    path: Path


def source_ids(text: str) -> tuple[str, ...]:
    return tuple(sorted(set(WP_RE.findall(text))))


def parse_cards(path: Path) -> list[Card]:
    """Parse supported TTC markdown draft conventions without YAML dependency."""
    text = path.read_text(encoding="utf-8")
    cards: list[Card] = []

    # YAML list records are used by the expansion authoring queue. Fields may
    # wrap over several lines, so inspect the complete record to the next ID.
    matches = list(LIST_ID_RE.finditer(text))
    for index, match in enumerate(matches):
        card_id = match.group(1)
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        record = text[match.start() : end]
        query_match = QUERY_RE.search(record)
        sources_match = SOURCE_LIST_RE.search(record)
        answerability_match = ANSWERABILITY_RE.search(record)
        cards.append(
            Card(
                card_id=card_id,
                query=query_match.group(1).strip() if query_match else None,
                source_ids=source_ids(sources_match.group(1)) if sources_match else (),
                answerability=answerability_match.group(1) if answerability_match else None,
                path=path,
            )
        )

    # The v1 candidate-card document expresses IDs as headings and evidence as
    # prose below a fenced query block. Parse the section up to the next heading.
    for match in HEADING_CARD_RE.finditer(text):
        card_id = match.group(1)
        if any(card.card_id == card_id for card in cards):
            continue
        end = text.find("#### ", match.end())
        section = text[match.end() : end if end != -1 else len(text)]
        query_match = QUERY_RE.search(section)
        answerability_match = ANSWERABILITY_RE.search(section)
        cards.append(
            Card(
                card_id=card_id,
                query=query_match.group(1).strip() if query_match else None,
                source_ids=source_ids(section),
                answerability=answerability_match.group(1) if answerability_match else None,
                path=path,
            )
        )
    return cards


def error(errors: list[str], message: str) -> None:
    errors.append(message)


def audit_drafts(cards: list[Card], errors: list[str]) -> None:
    by_id: dict[str, list[Card]] = defaultdict(list)
    for card in cards:
        by_id[card.card_id].append(card)
        if not card.query:
            error(errors, f"{card.path}: {card.card_id} has no parseable query")
        if not card.source_ids and not (card.answerability or "").startswith("unanswerable"):
            error(errors, f"{card.path}: {card.card_id} has no declared source IDs")
    for card_id, duplicates in sorted(by_id.items()):
        if len(duplicates) > 1:
            locations = ", ".join(str(card.path) for card in duplicates)
            error(errors, f"duplicate card ID {card_id}: {locations}")


def expect_string_list(value: Any, field: str, card_id: str, errors: list[str]) -> list[str]:
    if not isinstance(value, list) or not value or not all(isinstance(item, str) and item for item in value):
        error(errors, f"metadata {card_id}: {field} must be a non-empty list of strings")
        return []
    return value


def audit_metadata(cards: list[Card], metadata_path: Path, errors: list[str]) -> None:
    try:
        document = json.loads(metadata_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        error(errors, f"cannot parse metadata {metadata_path}: {exc}")
        return
    records = document.get("cards") if isinstance(document, dict) else None
    if not isinstance(records, list):
        error(errors, f"metadata {metadata_path}: top-level cards must be an array")
        return

    draft_ids = {card.card_id for card in cards}
    by_id: dict[str, dict[str, Any]] = {}
    family_partitions: dict[str, set[str]] = defaultdict(set)
    source_partitions: dict[str, set[str]] = defaultdict(set)
    partitions = {"development", "holdout", "regression"}

    for record in records:
        if not isinstance(record, dict):
            error(errors, f"metadata {metadata_path}: each cards item must be an object")
            continue
        card_id = record.get("id")
        if not isinstance(card_id, str) or not card_id:
            error(errors, f"metadata {metadata_path}: card has missing id")
            continue
        if card_id in by_id:
            error(errors, f"metadata {metadata_path}: duplicate card ID {card_id}")
            continue
        by_id[card_id] = record
        partition = record.get("partition")
        family = record.get("evidence_family_id")
        if partition not in partitions:
            error(errors, f"metadata {card_id}: partition must be development, holdout, or regression")
        if not isinstance(family, str) or not family:
            error(errors, f"metadata {card_id}: evidence_family_id is required")
        elif partition in partitions:
            family_partitions[family].add(partition)
        sources = record.get("source_document_ids", record.get("expected_source_ids"))
        for source in expect_string_list(sources, "source_document_ids", card_id, errors):
            if partition in partitions:
                source_partitions[source].add(partition)

    for card_id in sorted(draft_ids - set(by_id)):
        error(errors, f"metadata missing draft card {card_id}")
    for card_id in sorted(set(by_id) - draft_ids):
        error(errors, f"metadata references unknown card {card_id}")
    for family, family_sets in sorted(family_partitions.items()):
        if len(family_sets) > 1:
            error(errors, f"leakage: evidence family {family} crosses partitions {sorted(family_sets)}")
    for source, source_sets in sorted(source_partitions.items()):
        if len(source_sets) > 1:
            error(errors, f"leakage: source document {source} crosses partitions {sorted(source_sets)}; merge families or record an exception")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--cards", action="append", required=True, type=Path, help="candidate-card Markdown file; repeatable")
    parser.add_argument("--metadata", type=Path, help="optional future JSON split/family metadata manifest")
    args = parser.parse_args()

    errors: list[str] = []
    cards: list[Card] = []
    for path in args.cards:
        if not path.is_file():
            error(errors, f"candidate-card file does not exist: {path}")
            continue
        cards.extend(parse_cards(path))
    if not cards:
        error(errors, "no candidate cards parsed; supported forms are Markdown headings or one-line YAML list entries")
    audit_drafts(cards, errors)
    if args.metadata:
        audit_metadata(cards, args.metadata, errors)

    counts = Counter(card.path for card in cards)
    print(f"parsed {len(cards)} cards from {len(counts)} file(s)")
    for path, count in sorted(counts.items(), key=lambda item: str(item[0])):
        print(f"  {path}: {count}")
    if errors:
        print("AUDIT FAILED:", file=sys.stderr)
        for item in errors:
            print(f"- {item}", file=sys.stderr)
        return 1
    print("PASS: candidate IDs, queries, and declared source IDs are complete")
    if args.metadata:
        print("PASS: metadata partitions and evidence/source leakage families are complete")
    else:
        print("NOTE: no --metadata supplied; partition and leakage-family checks deferred")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
