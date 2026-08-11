#!/usr/bin/env python3
"""Score provisional TTC expansion traces against the candidate manifest."""

from __future__ import annotations

import argparse
import json
import statistics
from pathlib import Path


def percentile(values: list[int], fraction: float) -> int:
    if not values:
        return 0
    ordered = sorted(values)
    index = min(len(ordered) - 1, round((len(ordered) - 1) * fraction))
    return ordered[index]


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--traces", type=Path, required=True)
    parser.add_argument("--out", type=Path, required=True)
    args = parser.parse_args()

    manifest = json.loads(args.manifest.read_text(encoding="utf-8"))
    trace_file = json.loads(args.traces.read_text(encoding="utf-8"))
    cards = {card["id"]: card for card in manifest["cards"]}
    traces = [trace for trace in trace_file["traces"] if trace["id"] in cards]
    methods = {"bm25": "bm25", "vector": "vector", "hybrid": "hybrid"}
    metrics: dict[str, dict[str, float | int]] = {}

    for name, field in methods.items():
        answerable = 0
        mrr = 0.0
        recall = {1: 0, 3: 0, 10: 0}
        relevant_recall = 0.0
        for trace in traces:
            relevant = set(cards[trace["id"]]["relevantDocumentRevisionIds"])
            if not relevant:
                continue
            answerable += 1
            hits = trace[field]
            ranks = [index + 1 for index, hit in enumerate(hits) if hit["document_revision_id"] in relevant]
            if ranks:
                first = min(ranks)
                mrr += 1.0 / first
                for cutoff in recall:
                    if first <= cutoff:
                        recall[cutoff] += 1
            found = {hit["document_revision_id"] for hit in hits[:10]} & relevant
            relevant_recall += len(found) / len(relevant)
        denominator = max(1, answerable)
        metrics[name] = {
            "queries": len(traces),
            "answerable_queries": answerable,
            "recall_at_1": recall[1] / denominator,
            "recall_at_3": recall[3] / denominator,
            "recall_at_10": recall[10] / denominator,
            "mean_reciprocal_rank": mrr / denominator,
            "mean_relevant_recall_at_10": relevant_recall / denominator,
        }

    latencies = [trace["total_duration_ms"] for trace in traces]
    result = {
        "schema_version": "rag-eval-expansion-candidate-metrics/v1",
        "dataset_id": "candidate:ttc-expansion-v0",
        "dataset_status": "candidate",
        "trace_count": len(traces),
        "provisional": True,
        "methods": metrics,
        "latency_ms": {
            "mean": statistics.mean(latencies) if latencies else 0,
            "p50": percentile(latencies, 0.50),
            "p95": percentile(latencies, 0.95),
            "min": min(latencies) if latencies else 0,
            "max": max(latencies) if latencies else 0,
        },
        "embedding_cost": {"billed": 0, "note": "User-owned Ollama service; hardware and energy excluded."},
        "human_adjudication_required": True,
    }
    args.out.write_text(json.dumps(result, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
