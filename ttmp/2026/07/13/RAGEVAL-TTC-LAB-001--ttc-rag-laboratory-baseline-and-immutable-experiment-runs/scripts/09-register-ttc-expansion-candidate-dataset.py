#!/usr/bin/env python3
"""Register source-grounded TTC expansion cards as an immutable candidate set.

The script supports the two ticket-local YAML-list forms used by the 70-card
and 80-card drafts. It resolves stable ``wp:*`` IDs against one immutable
corpus snapshot and inserts a candidate evaluation dataset into ``rag-eval.db``.
It never updates an existing dataset ID with different content.
"""

from __future__ import annotations

import argparse
import json
import re
import sqlite3
from pathlib import Path


CARD_START = re.compile(r"^\s*-\s+\{?\s*id:\s*['\"]?([A-Za-z0-9][A-Za-z0-9._-]*)", re.MULTILINE)
QUERY = re.compile(r"\bquery:\s*['\"]([^'\"]+)['\"]")
SOURCES = re.compile(r"\bexpected_source_ids:\s*\[([^]]*)\]")
ANSWERABILITY = re.compile(r"\banswerability:\s*([A-Za-z0-9_-]+)")
WP_ID = re.compile(r"\bwp:\d+\b")


def parse_cards(paths: list[Path]) -> list[dict[str, object]]:
    cards: list[dict[str, object]] = []
    seen: set[str] = set()
    for path in paths:
        text = path.read_text(encoding="utf-8")
        starts = list(CARD_START.finditer(text))
        for index, match in enumerate(starts):
            end = starts[index + 1].start() if index + 1 < len(starts) else len(text)
            record = text[match.start() : end]
            card_id = match.group(1)
            if card_id in seen:
                raise ValueError(f"duplicate candidate card ID: {card_id}")
            query_match = QUERY.search(record)
            if not query_match:
                raise ValueError(f"card has no query: {card_id}")
            sources_match = SOURCES.search(record)
            source_ids = sorted(set(WP_ID.findall(sources_match.group(1)))) if sources_match else []
            answerability_match = ANSWERABILITY.search(record)
            cards.append(
                {
                    "id": card_id,
                    "query": query_match.group(1).strip(),
                    "relevantDocumentRevisionIds": [],
                    "provenance": {
                        "sourceDocumentIds": source_ids,
                        "answerability": answerability_match.group(1) if answerability_match else "unknown",
                        "origin": str(path),
                    },
                }
            )
            seen.add(card_id)
    if not cards:
        raise ValueError("no candidate cards parsed")
    cards.sort(key=lambda card: str(card["id"]))
    return cards


def resolve_revision_ids(database: sqlite3.Connection, snapshot_id: str, cards: list[dict[str, object]]) -> None:
    for card in cards:
        provenance = card["provenance"]
        assert isinstance(provenance, dict)
        revision_ids: list[str] = []
        for source_id in provenance["sourceDocumentIds"]:
            row = database.execute(
                """
                SELECT dr.id
                FROM corpus_snapshot_documents csd
                JOIN document_revisions dr ON dr.id = csd.document_revision_id
                WHERE csd.snapshot_id = ? AND dr.stable_document_id = ?
                """,
                (snapshot_id, source_id),
            ).fetchone()
            if row is None:
                raise ValueError(f"source {source_id} is absent from snapshot {snapshot_id}")
            revision_ids.append(row[0])
        card["relevantDocumentRevisionIds"] = sorted(set(revision_ids))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--db", default="data/rag-eval.db")
    parser.add_argument("--dataset-id", default="candidate:ttc-expansion-v0")
    parser.add_argument(
        "--snapshot-id",
        default="sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409",
    )
    parser.add_argument("--exclude-card", action="append", default=[], help="candidate ID withheld from this snapshot-bound set")
    parser.add_argument("--manifest-out", type=Path, help="optional path for the canonical manifest JSON")
    parser.add_argument("cards", nargs="+", type=Path)
    args = parser.parse_args()

    cards = parse_cards(args.cards)
    excluded = set(args.exclude_card)
    cards = [card for card in cards if card["id"] not in excluded]
    if excluded:
        print(f"excluded_cards={sorted(excluded)}")
    with sqlite3.connect(args.db) as database:
        resolve_revision_ids(database, args.snapshot_id, cards)
        manifest = {
            "schemaVersion": "rag-eval-evaluation-dataset/v1",
            "datasetStatus": "candidate",
            "binaryRelevantAtOrAbove": "2_SUBSTANTIAL",
            "cards": cards,
        }
        manifest_json = json.dumps(manifest, sort_keys=True, separators=(",", ":"))
        if args.manifest_out:
            args.manifest_out.write_text(manifest_json + "\n", encoding="utf-8")
        existing = database.execute(
            "SELECT manifest_json FROM evaluation_datasets WHERE id = ?", (args.dataset_id,)
        ).fetchone()
        if existing is not None:
            if existing[0] != manifest_json:
                raise ValueError(f"immutable dataset ID already has different content: {args.dataset_id}")
            print(f"already registered dataset={args.dataset_id} cards={len(cards)}")
            return 0
        database.execute(
            """
            INSERT INTO evaluation_datasets
              (id, schema_version, corpus_snapshot_id, status, manifest_json, query_count)
            VALUES (?, ?, ?, ?, ?, ?)
            """,
            (args.dataset_id, manifest["schemaVersion"], args.snapshot_id, "candidate", manifest_json, len(cards)),
        )
        database.commit()
    print(f"registered dataset={args.dataset_id} cards={len(cards)} snapshot={args.snapshot_id}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
