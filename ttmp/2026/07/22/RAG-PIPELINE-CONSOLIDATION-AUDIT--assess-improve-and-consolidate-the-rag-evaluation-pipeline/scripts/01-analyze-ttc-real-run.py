#!/usr/bin/env python3
"""Derive reproducible performance evidence from compact TTC cell/operation artifacts.

The script reads only the retained evidence, per-cell JSON, and bounded operation
JSONL files. It does not read source text, provider payloads, credentials, or the
removed runtime database.
"""

from __future__ import annotations

import argparse
import csv
import json
import math
import statistics
from collections import Counter
from datetime import datetime
from pathlib import Path
from typing import Any, Iterable

import matplotlib.pyplot as plt
from matplotlib.lines import Line2D


def parse_time(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def percentile(values: list[float], quantile: float) -> float:
    if not values:
        return 0.0
    ordered = sorted(values)
    position = (len(ordered) - 1) * quantile
    lower = math.floor(position)
    upper = math.ceil(position)
    if lower == upper:
        return ordered[lower]
    return ordered[lower] * (upper - position) + ordered[upper] * (position - lower)


def interval_union(intervals: Iterable[tuple[datetime, datetime]]) -> tuple[float, int]:
    events: list[tuple[datetime, int]] = []
    for start, end in intervals:
        events.extend(((start, 1), (end, -1)))
    events.sort(key=lambda item: (item[0], item[1]))
    active = 0
    peak = 0
    previous: datetime | None = None
    seconds = 0.0
    for stamp, delta in events:
        if active > 0 and previous is not None:
            seconds += (stamp - previous).total_seconds()
        active += delta
        peak = max(peak, active)
        previous = stamp
    return seconds, peak


def active_series(
    intervals: Iterable[tuple[datetime, datetime]], origin: datetime
) -> tuple[list[float], list[int]]:
    events: list[tuple[datetime, int]] = []
    for start, end in intervals:
        events.extend(((start, 1), (end, -1)))
    events.sort(key=lambda item: (item[0], item[1]))
    x = [0.0]
    y = [0]
    active = 0
    for stamp, delta in events:
        offset = (stamp - origin).total_seconds()
        x.extend((offset, offset))
        y.extend((active, active + delta))
        active += delta
    return x, y


def operation_measure(operation: dict[str, Any], name: str) -> int:
    for measure in operation.get("measures", []):
        if measure.get("name") == name:
            return int(measure.get("units", 0))
    return 0


def load_cells(run_root: Path) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    cells: list[dict[str, Any]] = []
    all_operations: list[dict[str, Any]] = []
    for cell_path in sorted((run_root / "cells").glob("*.json")):
        root = json.loads(cell_path.read_text())
        evidence = root["cell"]
        params = evidence["cell"]
        operations: list[dict[str, Any]] = []
        operation_path = run_root / evidence["operations"]["jsonlPath"]
        for line in operation_path.read_text().splitlines():
            record = json.loads(line)["operation"]
            completion = record["completion"]
            row = {
                "cell": cell_path.stem,
                "batch_size": params["chunksPerRequest"],
                "concurrency": params["concurrency"],
                "replicate": params["replicate"],
                "operation_id": record["operationId"],
                "node_key": record["nodeKey"],
                "attempt": record["attempt"],
                "kind": record["kind"]["name"],
                "outcome": completion["outcome"],
                "failure_class": completion.get("failure", {}).get("class", ""),
                "failure_code": completion.get("failure", {}).get("code", ""),
                "admitted_at": record["admittedAt"],
                "provider_started_at": completion["providerStartedAt"],
                "completed_at": completion["completedAt"],
                "elapsed_seconds": completion["elapsedMicros"] / 1_000_000,
                "queue_seconds": (
                    parse_time(completion["providerStartedAt"])
                    - parse_time(record["admittedAt"])
                ).total_seconds(),
                "chunk_count": operation_measure(record, "chunk_count"),
                "representation_count": operation_measure(record, "representation_count"),
            }
            operations.append(row)
            all_operations.append(row)

        attempts = evidence["attempts"] + evidence["embeddingAttempts"]
        canonical_start = min(parse_time(item["startedAt"]) for item in attempts)
        canonical_end = max(parse_time(item["finishedAt"]) for item in attempts)
        operation_start = min(parse_time(item["admitted_at"]) for item in operations)
        operation_end = max(parse_time(item["completed_at"]) for item in operations)
        inclusive_start = min(canonical_start, operation_start)
        inclusive_end = max(canonical_end, operation_end)
        inclusive_seconds = (inclusive_end - inclusive_start).total_seconds()

        generation = [item for item in operations if item["kind"] == "provider.generate"]
        embedding = [item for item in operations if item["kind"] == "provider.embed"]
        successful_generation = [item for item in generation if item["outcome"] == "succeeded"]
        failed_generation = [item for item in generation if item["outcome"] != "succeeded"]
        generation_intervals = [
            (parse_time(item["provider_started_at"]), parse_time(item["completed_at"]))
            for item in generation
        ]
        embedding_intervals = [
            (parse_time(item["provider_started_at"]), parse_time(item["completed_at"]))
            for item in embedding
        ]
        provider_union_seconds, provider_peak = interval_union(
            generation_intervals + embedding_intervals
        )
        generation_union_seconds, generation_peak = interval_union(generation_intervals)
        embedding_union_seconds, embedding_peak = interval_union(embedding_intervals)
        generation_latencies = [item["elapsed_seconds"] for item in successful_generation]
        embedding_latencies = [item["elapsed_seconds"] for item in embedding]
        generation_work_seconds = sum(item["elapsed_seconds"] for item in generation)
        successful_generation_work_seconds = sum(generation_latencies)
        embedding_work_seconds = sum(embedding_latencies)
        mean_generation = statistics.mean(generation_latencies)
        generation_cv = (
            statistics.pstdev(generation_latencies) / mean_generation
            if len(generation_latencies) > 1 and mean_generation
            else 0.0
        )
        canonical_seconds = evidence["makespanMicros"] / 1_000_000
        cells.append(
            {
                "cell": cell_path.stem,
                "batch_size": params["chunksPerRequest"],
                "concurrency": params["concurrency"],
                "replicate": params["replicate"],
                "started_at": inclusive_start.isoformat().replace("+00:00", "Z"),
                "canonical_seconds": canonical_seconds,
                "inclusive_seconds": inclusive_seconds,
                "boundary_undercount_seconds": inclusive_seconds - canonical_seconds,
                "chunks": evidence["chunks"],
                "planned_generation_requests": evidence["requests"],
                "generation_admissions": len(generation),
                "generation_succeeded": len(successful_generation),
                "generation_failed": len(failed_generation),
                "retry_rate": len(failed_generation) / len(generation),
                "embedding_requests": len(embedding),
                "generation_p50_seconds": statistics.median(generation_latencies),
                "generation_p95_seconds": percentile(generation_latencies, 0.95),
                "generation_min_seconds": min(generation_latencies),
                "generation_max_seconds": max(generation_latencies),
                "generation_cv": generation_cv,
                "embedding_p50_seconds": statistics.median(embedding_latencies),
                "embedding_p95_seconds": percentile(embedding_latencies, 0.95),
                "generation_work_seconds": generation_work_seconds,
                "successful_generation_work_seconds": successful_generation_work_seconds,
                "embedding_work_seconds": embedding_work_seconds,
                "generation_union_seconds": generation_union_seconds,
                "embedding_union_seconds": embedding_union_seconds,
                "provider_union_seconds": provider_union_seconds,
                "provider_coverage": provider_union_seconds / inclusive_seconds,
                "provider_idle_seconds": inclusive_seconds - provider_union_seconds,
                "generation_slot_occupancy": generation_work_seconds
                / (inclusive_seconds * params["concurrency"]),
                "generation_peak": generation_peak,
                "embedding_peak": embedding_peak,
                "provider_peak": provider_peak,
                "chunks_per_second": evidence["chunks"] / inclusive_seconds,
                "successful_requests_per_second": len(successful_generation) / inclusive_seconds,
                "queue_p95_seconds": percentile(
                    [item["queue_seconds"] for item in operations], 0.95
                ),
                "operations": operations,
                "inclusive_start": inclusive_start,
                "inclusive_end": inclusive_end,
            }
        )
    cells.sort(key=lambda item: item["inclusive_start"])
    for index, cell in enumerate(cells, start=1):
        cell["execution_order"] = index
    return cells, all_operations


def write_csv(path: Path, rows: list[dict[str, Any]], fields: list[str]) -> None:
    with path.open("w", newline="") as handle:
        writer = csv.DictWriter(
            handle, fieldnames=fields, extrasaction="ignore", lineterminator="\n"
        )
        writer.writeheader()
        writer.writerows(rows)


def public_cell(cell: dict[str, Any]) -> dict[str, Any]:
    return {
        key: value
        for key, value in cell.items()
        if key not in {"operations", "inclusive_start", "inclusive_end"}
    }


def style_axis(axis: Any) -> None:
    axis.grid(True, alpha=0.25, axis="y")
    axis.spines[["top", "right"]].set_visible(False)


def render_graphs(cells: list[dict[str, Any]], graph_dir: Path) -> list[str]:
    graph_dir.mkdir(parents=True, exist_ok=True)
    colors = {1: "#2563eb", 2: "#ea580c"}
    labels = [f"b{item['batch_size']}/c{item['concurrency']}" for item in cells]
    produced: list[str] = []

    # 1. Canonical vs operation-inclusive elapsed time.
    ordered = sorted(cells, key=lambda item: (item["batch_size"], item["concurrency"]))
    x = list(range(len(ordered)))
    fig, axis = plt.subplots(figsize=(11, 5.5))
    canonical = [item["canonical_seconds"] for item in ordered]
    inclusive = [item["inclusive_seconds"] for item in ordered]
    width = 0.38
    axis.bar([item - width / 2 for item in x], canonical, width=width, color="#94a3b8", label="published canonical makespan")
    axis.bar([item + width / 2 for item in x], inclusive, width=width, color="#dc2626", alpha=0.75, label="operation-inclusive elapsed")
    for index, item in enumerate(ordered):
        if item["boundary_undercount_seconds"] > 0.01:
            axis.text(index + width / 2, item["inclusive_seconds"] + 4, f"+{item['boundary_undercount_seconds']:.1f}s", ha="center", color="#991b1b", fontweight="bold")
    axis.set_xticks(x, [f"b{item['batch_size']} / c{item['concurrency']}" for item in ordered])
    axis.set_ylabel("Elapsed seconds")
    axis.set_title("TTC real run: published vs operation-inclusive elapsed boundary")
    axis.legend()
    style_axis(axis)
    fig.tight_layout()
    for suffix in ("png", "svg"):
        name = f"01-elapsed-boundary.{suffix}"
        fig.savefig(graph_dir / name, dpi=180)
        produced.append(name)
    plt.close(fig)

    # 2. Throughput and concurrency speedup.
    fig, (left, right) = plt.subplots(1, 2, figsize=(12, 5))
    for concurrency in (1, 2):
        series = sorted(
            [item for item in cells if item["concurrency"] == concurrency],
            key=lambda item: item["batch_size"],
        )
        left.plot(
            [item["batch_size"] for item in series],
            [item["chunks_per_second"] for item in series],
            marker="o",
            linewidth=2,
            color=colors[concurrency],
            label=f"generation concurrency {concurrency}",
        )
    by_key = {(item["batch_size"], item["concurrency"]): item for item in cells}
    batches = [1, 2, 4, 8]
    speedups = [
        by_key[(batch, 1)]["inclusive_seconds"] / by_key[(batch, 2)]["inclusive_seconds"]
        for batch in batches
    ]
    right.bar([str(value) for value in batches], speedups, color="#7c3aed")
    right.axhline(2.0, color="#334155", linestyle="--", label="ideal 2×")
    for index, speedup in enumerate(speedups):
        right.text(index, speedup + 0.04, f"{speedup:.2f}×", ha="center")
    left.set(xlabel="Chunks per generation request", ylabel="Completed chunks / second", xticks=batches)
    left.set_title("End-to-end throughput")
    left.legend()
    right.set(xlabel="Chunks per generation request", ylabel="Speedup: c1 elapsed / c2 elapsed")
    right.set_title("Observed concurrency speedup")
    right.legend()
    style_axis(left)
    style_axis(right)
    fig.suptitle("Batching and concurrency effects (operation-inclusive boundary)")
    fig.tight_layout()
    for suffix in ("png", "svg"):
        name = f"02-throughput-speedup.{suffix}"
        fig.savefig(graph_dir / name, dpi=180)
        produced.append(name)
    plt.close(fig)

    # 3. Generation latency distributions with failed operations overlaid.
    ordered = sorted(cells, key=lambda item: (item["batch_size"], item["concurrency"]))
    success_values: list[list[float]] = []
    failed_points: list[tuple[int, float]] = []
    for index, cell in enumerate(ordered, start=1):
        generation = [item for item in cell["operations"] if item["kind"] == "provider.generate"]
        success_values.append(
            [item["elapsed_seconds"] for item in generation if item["outcome"] == "succeeded"]
        )
        failed_points.extend(
            (index, item["elapsed_seconds"])
            for item in generation
            if item["outcome"] != "succeeded"
        )
    fig, axis = plt.subplots(figsize=(12, 5.5))
    boxes = axis.boxplot(
        success_values,
        patch_artist=True,
        showmeans=True,
        meanprops={"marker": "^", "markerfacecolor": "#15803d", "markeredgecolor": "#15803d"},
    )
    for patch, item in zip(boxes["boxes"], ordered):
        patch.set_facecolor(colors[item["concurrency"]])
        patch.set_alpha(0.55)
    if failed_points:
        axis.scatter(
            [item[0] + (0.12 if index % 2 else -0.12) for index, item in enumerate(failed_points)],
            [item[1] for item in failed_points],
            marker="x",
            color="#dc2626",
            s=85,
            linewidth=2.5,
            label="failed transport operation",
            zorder=3,
        )
    axis.set_xticks(range(1, len(ordered) + 1), [f"b{i['batch_size']}/c{i['concurrency']}" for i in ordered])
    axis.set_ylabel("Provider generation latency (seconds)")
    axis.set_title("Successful generation latency distribution and failed calls")
    axis.legend(
        handles=[
            Line2D([0], [0], marker="^", color="#15803d", linestyle="None", markersize=8, label="successful mean"),
            Line2D([0], [0], marker="x", color="#dc2626", linestyle="None", markersize=9, label="failed transport operation"),
        ]
    )
    style_axis(axis)
    fig.tight_layout()
    for suffix in ("png", "svg"):
        name = f"03-generation-latency.{suffix}"
        fig.savefig(graph_dir / name, dpi=180)
        produced.append(name)
    plt.close(fig)

    # 4. Provider occupancy, coverage, and retries use separate denominators.
    fig, (left, middle, right) = plt.subplots(1, 3, figsize=(16, 5))
    labels = [f"b{i['batch_size']}/c{i['concurrency']}" for i in ordered]
    occupancy = [100 * item["generation_slot_occupancy"] for item in ordered]
    coverage = [100 * item["provider_coverage"] for item in ordered]
    retry_rates = [100 * item["retry_rate"] for item in ordered]
    left.bar(labels, occupancy, color="#0f766e")
    left.set_ylim(80, 101)
    left.set_ylabel("Occupied configured generation capacity (%)")
    left.set_title("Generation slot occupancy")
    middle.bar(labels, coverage, color="#4f46e5")
    middle.set_ylim(98, 100.1)
    middle.set_ylabel("Inclusive elapsed with any provider active (%)")
    middle.set_title("Provider-time coverage")
    right.bar(labels, retry_rates, color="#dc2626")
    right.set_ylim(0, max(retry_rates) * 1.25)
    right.set_ylabel("Failed generation admissions (%)")
    right.set_title("Transport retry incidence")
    for axis, values in ((left, occupancy), (middle, coverage), (right, retry_rates)):
        for index, value in enumerate(values):
            if axis is right:
                axis.text(index, value + 0.4, f"{value:.1f}%", ha="center", va="bottom", fontsize=8)
            else:
                axis.text(index, value - 0.12, f"{value:.1f}%", ha="center", va="top", fontsize=8)
        axis.tick_params(axis="x", rotation=35)
        style_axis(axis)
    fig.suptitle("Provider occupancy and retries (single replicate; truncated occupancy axes are labeled)")
    fig.tight_layout()
    for suffix in ("png", "svg"):
        name = f"04-occupancy-retries.{suffix}"
        fig.savefig(graph_dir / name, dpi=180)
        produced.append(name)
    plt.close(fig)

    # 5. One active-provider timeline per cell, preserving execution order.
    fig, axes = plt.subplots(len(cells), 1, figsize=(12, 15), sharex=True)
    common_x_max = max(item["inclusive_seconds"] for item in cells)
    common_y_max = max(item["provider_peak"] for item in cells)
    for axis, cell in zip(axes, cells):
        origin = cell["inclusive_start"]
        generation_intervals = [
            (parse_time(item["provider_started_at"]), parse_time(item["completed_at"]))
            for item in cell["operations"]
            if item["kind"] == "provider.generate"
        ]
        embedding_intervals = [
            (parse_time(item["provider_started_at"]), parse_time(item["completed_at"]))
            for item in cell["operations"]
            if item["kind"] == "provider.embed"
        ]
        gx, gy = active_series(generation_intervals, origin)
        ex, ey = active_series(embedding_intervals, origin)
        axis.plot(gx, gy, color="#2563eb", linewidth=1.8, label="generation provider spans")
        axis.plot(ex, ey, color="#ea580c", linewidth=2.2, linestyle=":", label="embedding provider spans")
        axis.axhline(cell["concurrency"], color="#334155", linestyle="--", linewidth=0.8)
        axis.set_ylabel(f"b{cell['batch_size']}/c{cell['concurrency']}\nactive calls")
        axis.set_xlim(0, common_x_max)
        axis.set_ylim(0, common_y_max + 0.2)
        axis.set_yticks(range(0, common_y_max + 1))
        axis.grid(True, alpha=0.2)
    axes[0].legend(ncol=2, loc="upper right")
    axes[-1].set_xlabel("Seconds from operation-inclusive cell start")
    fig.suptitle("Provider activity timeline for every cell, in execution order")
    fig.tight_layout(rect=(0, 0, 1, 0.98))
    for suffix in ("png", "svg"):
        name = f"05-provider-timelines.{suffix}"
        fig.savefig(graph_dir / name, dpi=180)
        produced.append(name)
    plt.close(fig)

    # Matplotlib emits path-data lines with trailing spaces. Normalize generated
    # text artifacts so repository diff checks remain a meaningful gate.
    for svg in graph_dir.glob("*.svg"):
        svg.write_text("\n".join(line.rstrip() for line in svg.read_text().splitlines()) + "\n")

    return produced


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--run-root", type=Path, required=True)
    parser.add_argument("--output-dir", type=Path, required=True)
    args = parser.parse_args()
    args.output_dir.mkdir(parents=True, exist_ok=True)

    cells, operations = load_cells(args.run_root)
    public_cells = [public_cell(item) for item in cells]
    cell_fields = list(public_cells[0].keys())
    operation_fields = list(operations[0].keys())
    write_csv(args.output_dir / "cells-derived.csv", public_cells, cell_fields)
    write_csv(args.output_dir / "operations-derived.csv", operations, operation_fields)

    failure_counts = Counter(
        (item["failure_class"], item["failure_code"])
        for item in operations
        if item["outcome"] != "succeeded"
    )
    summary = {
        "schemaVersion": "rag-ttc-pipeline-audit/v1",
        "sourceRunRoot": str(args.run_root),
        "cellCount": len(cells),
        "executionOrder": [item["cell"] for item in cells],
        "totals": {
            "chunksAcrossCells": sum(item["chunks"] for item in cells),
            "plannedGenerationRequests": sum(item["planned_generation_requests"] for item in cells),
            "generationAdmissions": sum(item["generation_admissions"] for item in cells),
            "generationSucceeded": sum(item["generation_succeeded"] for item in cells),
            "generationFailed": sum(item["generation_failed"] for item in cells),
            "embeddingRequests": sum(item["embedding_requests"] for item in cells),
            "failedProviderSeconds": sum(
                item["elapsed_seconds"] for item in operations if item["outcome"] != "succeeded"
            ),
        },
        "failureCounts": [
            {"class": key[0], "code": key[1], "count": count}
            for key, count in sorted(failure_counts.items())
        ],
        "cells": public_cells,
    }
    (args.output_dir / "summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    graphs = render_graphs(cells, args.output_dir / "graphs")
    manifest = {
        "schemaVersion": "rag-ttc-pipeline-audit-graphs/v1",
        "sourceRunRoot": str(args.run_root),
        "graphs": graphs,
    }
    (args.output_dir / "graphs" / "manifest.json").write_text(
        json.dumps(manifest, indent=2) + "\n"
    )
    print(
        f"cells={len(cells)} operations={len(operations)} "
        f"graphs={len(graphs)} output={args.output_dir}"
    )


if __name__ == "__main__":
    main()
