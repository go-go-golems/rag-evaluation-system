#!/usr/bin/env python3
import argparse
import json
from pathlib import Path


def render(root: Path) -> str:
    manifest = json.loads((root / "manifest.json").read_text())
    cases = []
    for case in manifest["cases"]:
        domain = json.loads((root / case["domainConfig"]["path"]).read_text())
        refs = []
        for key, kind, ident in [
            ("execution", "rag-execution", "execution-" + case["id"]),
            ("corpus", "rag-corpus", "corpus-" + case["id"]),
            ("queries", "rag-query-set", "queries-" + case["id"]),
        ]:
            value = case[key]
            refs.append({
                "role": "workflow-input", "kind": kind, "id": ident,
                "digest": value["digest"], "sizeBytes": value["sizeBytes"],
                "schemaVersion": value["schemaVersion"],
                "uri": "inputs/rag-workflow/" + value["path"],
                "mediaType": "application/json",
            })
        cases.append((case["id"], domain, refs))
    lines = ['const research = require("researchctl");', '']
    for index, (_, domain, _) in enumerate(cases):
        lines.append(f"const execution{index} = " + json.dumps(domain, sort_keys=True, separators=(",", ":")) + ";")
    lines += [
        '', 'function specification(label, execution, inputs) {', '  return {',
        '    canonicalIdentity: {', '      schemaVersion: "researchctl-execution-spec/v1",',
        '      identityScheme: "researchctl-execution-identity/v1",', '      domain: "scraper-workflow",',
        '      domainSchemaVersion: "scraper-workflow-execution/v2",', '      inputs,',
        '      domainConfig: execution,', '      requestedMeasures: [',
        '        {name: "rag.mrr", valueKind: "number", unit: "ratio", required: true},',
        '        {name: "workflow.elapsed", valueKind: "number", unit: "microseconds", required: true},',
        '        {name: "workflow.retries", valueKind: "number", unit: "count", required: true}',
        '      ],', '      factors: {ragCase: label}', '    },',
        '    displayName: `RAG Workflow ${label}`,',
        '    provenance: {authoring: "examples/rag-workflow/researchctl-plan.js"},',
        '    labels: {fixture: "rag-v2-workflow-lowering"}', '  };', '}', '',
        'module.exports = research.experimentPlan("rag-workflow-matrix", plan => plan',
        '  .experiment("EXP-RAG-WORKFLOW")',
    ]
    for index, (case_id, _, refs) in enumerate(cases):
        lines.append(f'  .case("{case_id}", value => value.specification(specification("{case_id}", execution{index}, {json.dumps(refs, sort_keys=True, separators=(",", ":"))})).factors({{ragCase: "{case_id}"}}).replicates(2))')
    lines += ['  .ordering({strategy: "blocked"})', '  .execution({maxConcurrent: 2, failFast: false})', ');', '']
    return "\n".join(lines)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--fixture-root", required=True)
    parser.add_argument("--out", default="-")
    args = parser.parse_args()
    body = render(Path(args.fixture_root))
    if args.out == "-":
        print(body, end="")
    else:
        Path(args.out).write_text(body)

if __name__ == "__main__":
    main()
