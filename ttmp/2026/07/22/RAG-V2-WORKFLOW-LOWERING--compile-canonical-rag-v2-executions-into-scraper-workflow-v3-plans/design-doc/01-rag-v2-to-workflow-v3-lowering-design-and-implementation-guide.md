---
Title: RAG v2 to Workflow V3 lowering design and implementation guide
Ticket: RAG-V2-WORKFLOW-LOWERING
Status: complete
Topics:
    - rag
    - rag-eval
    - workflow
    - scraper
    - javascript
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/gojamodules/workflow/authoring.go
      Note: Target Workflow V3 authoring and plan model
    - Path: repo://pkg/ragcompiler/targets.go
      Note: RAG study and execution compilation
    - Path: repo://pkg/ragcontract/types.go
      Note: Canonical RAG execution contracts
    - Path: repo://pkg/ragengine/engine.go
      Note: Semantic parity source
    - Path: repo://pkg/ragworkflow/lower.go
      Note: Production deterministic lowering backend
    - Path: repo://pkg/ragworkflow/registry.go
      Note: Closed versioned operator lowering registry
    - Path: repo://pkg/ragworkflow/runtime.go
      Note: Provider-free task runtime and preparation/query boundary
    - Path: repo://pkg/ragworkflow/projector.go
      Note: Bounded privacy-safe RAG observation projection
ExternalSources: []
Summary: Design for preserving canonical RAG v2 semantics while compiling executions into durable Workflow V3 plans.
LastUpdated: 2026-07-22T23:15:00-04:00
WhatFor: Make Scraper the execution backend for RAG without leaking RAG concepts into Scraper.
WhenToUse: Read before implementing RAG Workflow V3 tasks, lowering, or replacing ragengine execution.
---


# RAG v2 to Workflow V3 lowering design and implementation guide

## Program context

This ticket belongs to **EXPERIMENT-PLATFORM-CONVERGENCE**. Siblings are `RESEARCHCTL-EXPERIMENT-PLANS`, `SCRAPER-WORKFLOW-V3-PRODUCT-CUTOVER`, `EXPERIMENT-PLATFORM-SCRAPER-RUNNER`, `SCRAPER-WORKFLOW-OBSERVATIONS`, `RAG-GEPPETTO-WORKFLOW-OPERATIONS`, `RESEARCHCTL-EXPERIMENT-ANALYSIS`, `RAG-V2-EXECUTION-CUTOVER`, and `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE`. Their guides are under the Researchctl, Scraper, and RAG repositories' `ttmp/2026/07/22/` directories.

## Executive summary

RAG v2 already has a clean authoring/compiler boundary. JavaScript creates Go-backed pipeline, product, and study values; `pkg/ragcompiler` emits canonical data-only contracts in `pkg/ragcontract`; `pkg/ragengine` executes operators. This ticket keeps the authoring and canonical IR, but adds a lowering backend that maps a `rag-pipeline-execution/v2` into a Scraper Workflow V3 plan.

The first slice is deterministic and provider-free: corpus input, chunking, raw representation, fixture embeddings, indexing, retrieval, and evaluation. Geppetto calls are added by the separate operations ticket. Researchctl integration is already generic through the Scraper runner.

## Implementation outcome

This design is implemented and accepted. `pkg/ragworkflow` now provides the exact operator registry, deterministic lowerer, versioned `rag-v2-provider-free@1.0.0` package, trusted `rag:workflow` runtime module, preparation/index evidence, bounded query map and reduction, result contract, and privacy-safe domain projector. The RAG-owned process runner links those capabilities into generic `scraper-workflow-execution/v2` without adding RAG semantics to Scraper.

Deterministic examples are under `examples/rag-workflow`. The built-binary acceptance, generator, and boundary guards are under this ticket's `scripts/`; the durable summary is `sources/smoke/01-summary.json`; the requirement-to-evidence audit is `analysis/01-rag-workflow-lowering-acceptance-and-boundary-audit.md`. Real provider operations remain explicitly rejected and belong to `RAG-GEPPETTO-WORKFLOW-OPERATIONS`.

## RAG v2 orientation

Read these packages in order:

1. `pkg/ragcontract/README.md` and `types.go`: sole wire-contract owner.
2. `pkg/ragmodel/`: fluent Go model behind JavaScript.
3. `pkg/gojamodules/rag/`: pure `require("rag")` authoring.
4. `pkg/ragcompiler/`: normalization, validation, factor expansion, executions.
5. `pkg/ragoperators/`: operator implementations, artifacts, metrics.
6. `pkg/ragengine/`: current in-process execution semantics.
7. `cmd/rag-worker/`: current Researchctl runner adapter.

Workflow V3 orientation is in the Scraper cutover guide.

## Compiler layering

```text
RAG JS source
   -> ragmodel values
   -> ragcompiler canonical PipelineIR / PipelineExecution
   -> NEW ragworkflow compiler
   -> Workflow V3 IR
   -> Workflow V3 compiler
   -> WorkflowPlan
```

The RAG IR remains authoritative for domain identity. Workflow plans are execution artifacts derived from that identity. A change in RAG operator semantics changes the operator version and therefore the derived plan digest.

## Task catalog

Initial task keys:

```text
rag.corpus.load/v1
rag.units.prepare/v1
rag.chunks.create/v1
rag.represent.raw/v1
rag.embed.batch/v1              // fixture first
rag.index.build/v1
rag.retrieve/v1
rag.rank/v1
rag.evaluate/v1
rag.results.publish/v1
```

Each task has strict input/output artifact schemas and a host implementation. The companion descriptor module allows JS workflow authors to reference tasks, but ordinary RAG users should use the automatic lowering compiler.

```go
type Lowerer interface {
    Lower(ctx context.Context, execution ragcontract.PipelineExecution) (
        workflowv3.WorkflowIR, error,
    )
}

type OperatorLowerer interface {
    OperatorID() ragcontract.OperatorRef
    LowerNode(LoweringContext, ragcontract.PipelineNode) (NodeFragment, error)
}
```

The registry is closed over explicit operator versions. Unknown operators fail compilation; they do not become generic callback tasks.

## Example lowering

RAG intent:

```text
corpus -> units -> chunks -> raw representations -> embeddings
                                             \-> lexical index
embeddings -> vector index
query -> retrieve both -> weighted RRF -> evaluate
```

Workflow result:

```text
[load corpus]
      |
[prepare units]
      |
[chunk] -> set<chunk>
      |
+-----+------------------+
| map raw representation |
+-----+------------------+
      | set<representation>
      +---->[map embedding batches]---->[vector index]
      +-------------------------------->[lexical index]
                                            |
query set -> map retrieve -> map rank -> reduce evaluation -> publish
```

Static preparation and per-query execution should be separate subgraphs when identity permits reuse. The compiler, not the workflow engine, decides which RAG nodes are static.

## Lowering pseudocode

```pseudo
validateCanonicalExecution(execution)
builder = new WorkflowIR(name=execution.pipeline.digest)
inputs = bindCorpusAndEvaluationArtifacts(builder, execution.bindings)
values = map[inputNodeIDs]inputs

for node in topologicalOrder(execution.pipeline.nodes):
    lowerer = registry.require(node.operator.id)
    fragment = lowerer.lower(node, values)
    builder.append(fragment.nodes, fragment.maps, fragment.reductions)
    values[node.id] = fragment.outputs

for requestedMeasure in execution.measures:
    attachEvaluationOrProjection(requestedMeasure, values, builder)

builder.output("rag-result", publish(values))
return workflow.validateAndCompile(builder)
```

## Artifact rules

Workflow control state carries only `ArtifactRef`. RAG domain outputs remain canonical manifests:

- corpus, unit, chunk, representation, embedding, index, evaluation, and query trace manifests;
- schema version and semantic digest;
- content-addressed locator;
- explicit source lineage.

Do not serialize live indexes or provider clients into plans. Live values are rebuilt or opened by task implementations from immutable artifacts.

## Decisions

### Decision: RAG owns lowering

- **Context:** Scraper cannot know whether a RAG node is batchable, static, cacheable, or query-dependent.
- **Decision:** A RAG package compiles canonical RAG executions into generic Workflow V3 IR.
- **Consequences:** Scraper remains domain neutral; RAG must test semantic parity across backends during cutover.
- **Status:** accepted.

### Decision: registered Go tasks, not serialized JS callbacks

- **Decision:** Plans reference versioned task keys and immutable configuration.
- **Rationale:** Durable replay and isolated execution require stable implementations.
- **Status:** accepted.

### Decision: no legacy DTO bridge

- **Decision:** Lower directly from active RAG v2 contracts. Delete older workflow representations during the execution cutover.
- **Status:** accepted.

## Implementation phases

1. Add `pkg/ragworkflow` with task specs, schemas, and lowering registry.
2. Lower corpus, units, chunks, and raw representation.
3. Add fixture embedding and indexes.
4. Add retrieval, ranking, and evaluation.
5. Emit RAG artifacts and metrics through task context.
6. Add static/query subgraph separation and reuse rules.
7. Add Researchctl/Scraper end-to-end fixture.
8. Prove parity with `ragengine` fixture goldens.
9. Hand off real-provider tasks to `RAG-GEPPETTO-WORKFLOW-OPERATIONS`.

## Test strategy

- Canonical IR to Workflow IR goldens.
- Unknown operator/version rejection.
- Deterministic plan digest.
- Artifact schema mismatch rejection.
- Map batching and bounded materialization.
- Static preparation reuse only with matching pipeline/corpus/profile fingerprints.
- Retrieval metric parity with current operator fixtures.
- Crash/restart between preparation and query phases.
- No provider contact in deterministic fixture tests.

## Intern guidance

Do not rewrite chunking or retrieval algorithms. Treat `ragoperators` as the semantic source and extract task-sized entry points. Start with one raw BM25 pipeline; avoid the full operator matrix. A good first milestone is a plan JSON whose nodes can be understood by comparing each one to the RAG pipeline node that generated it.

## Completion criteria

A canonical fixture RAG execution compiles to a deterministic Workflow V3 plan, executes through Scraper, produces the same RAG metrics and artifacts as the fixture engine, and is recorded by Researchctl without a RAG-specific experiment runner loop.

## Technology primer: RAG data and execution stages

A RAG pipeline transforms source records into searchable evidence and then uses that evidence to answer or evaluate queries. Preparation and querying have different reuse boundaries. Preparation creates units, chunks, representations, embeddings, and indexes from a corpus. Query execution retrieves and ranks evidence, optionally generates an answer, and computes evaluation metrics. A workflow compiler must preserve these semantics rather than merely turn every operator into a sequential task.

```text
Corpus preparation (reusable when identity matches)
source records -> units -> chunks -> representations -> embeddings -> indexes

Query execution (repeated per dataset/query)
query -> candidate retrieval -> fusion/reranking -> evidence hydration
      -> answer generation -> evaluation
```

A corpus digest, pipeline digest, prompt/profile fingerprint, and operator versions determine whether prepared artifacts are reusable. Changing concurrency alone may not change semantic output identity, but changing batching can affect provider prompts or ordering and therefore must be modeled explicitly.

## Canonical RAG v2 model

`ragcontract.PipelineIR` is data. Nodes refer to versioned operators and bind typed ports. `ragcompiler.Normalize` validates and canonicalizes the graph. `PipelineExecution` binds a pipeline to immutable artifacts, selected factors, requested measures, and a dataset. This execution is the correct input to lowering because all authoring callbacks have already disappeared.

The distinction between semantic and execution identity matters. The RAG pipeline says what result should be computed. The Workflow V3 plan says how durable tasks compute it. Two workflow backends could implement the same RAG execution, but during this program Workflow V3 becomes the sole production backend.

## Guided source tour

In `pkg/gojamodules/rag/module.go`, follow `pipelineFactory`, the builder object, and `compileStudyValue`. Notice that runtime-private symbols protect Go-backed values and `compileStudy` returns plain canonical data. In `pkg/ragcompiler/targets.go`, follow factor selections into expanded `PipelineExecution` values. In `pkg/ragengine/engine.go`, trace static node execution, query nodes, observer metrics, and artifacts. The lowering compiler must preserve those externally observable semantics.

`pkg/ragoperators/registry.go` provides operator lookup. Individual operator files show the domain types needed by tasks. Do not have workflow tasks exchange arbitrary `map[string]any`; use the canonical manifest schemas already defined in `pkg/ragcontract`.

## Worked lowering: raw hybrid retrieval

For a pipeline with recursive chunks, raw representations, embeddings, lexical/vector indexes, weighted fusion, and MRR, lowering might create:

```text
load-corpus
  -> create-units
  -> create-chunks (output set manifest)
      -> map raw-representation per bounded page
          -> build-lexical-index
          -> map embedding batches -> build-vector-index
query-set
  -> map retrieve-lexical
  -> map retrieve-vector
  -> map weighted-fusion
  -> reduce evaluation metrics
  -> publish evaluation artifact
```

The workflow task `build-vector-index` consumes an embedding-set manifest, not a Go slice held in memory by the compiler. Large item sets remain behind manifests so restart and isolation do not require placing all values in workflow control rows.

## Node-lowering contract in detail

```go
type LoweringContext interface {
    Value(nodeID, port string) (workflowv3.ValueRef, bool)
    Set(nodeID, port string) (workflowv3.SetRef, bool)
    AddTask(key workflowv3.NodeKey, task workflowv3.TaskInvocation) NodeHandle
    AddMap(...)
    AddReduce(...)
    BindOutput(nodeID, port string, value workflowv3.ValueRef)
}
```

An `OperatorLowerer` validates operator configuration, requires input references, appends one or more workflow constructs, and binds outputs. It must not open artifacts or call providers. Lowering is pure and deterministic.

Some RAG operators lower to more than one task. Combined preparation may partition chunks, invoke batched generation, validate outputs, flatten representations, and publish a manifest. This is acceptable: one domain operator is not required to equal one scheduler node.

## Artifact lifecycle

Take an embedding set as an example. Batch tasks produce bounded batch-result artifacts. A deterministic publish/reduce task verifies cardinality and source identities, orders records, and writes the canonical embedding manifest. Downstream index tasks consume only that publication. If one batch retries, previously successful batches remain reusable within the run.

```text
chunk manifest
  -> batch partition manifest
      -> embedding batch artifacts
          -> verified embedding-set publication
              -> vector index artifact
```

The publication step is where domain completeness becomes authoritative. A directory containing some batch files is not a completed embedding set.

## Semantic parity strategy

Before replacing `ragengine`, run the same fixture execution through both backends. Compare canonical domain outputs and metrics, not scheduler traces. Expected differences in workflow-specific artifacts should be excluded. Any semantic difference requires either a lowering fix or an explicit versioned contract change.

```pseudo
direct = ragengine.execute(fixtureExecution)
workflow = executeThroughWorkflow(lower(fixtureExecution))
assert canonical(direct.ragArtifacts) == canonical(workflow.ragArtifacts)
assert canonical(direct.ragMetrics) == canonical(workflow.ragMetrics)
```

## First-week implementation route

Start with `examples/rag-v2/06-raw-study.js`, which avoids provider complexity. Compile it and inspect the execution JSON. Manually sketch the expected workflow graph. Implement only corpus load, chunks, raw representation, lexical index, retrieval, and MRR. Execute with fixture inputs. Once parity holds, add embedding as a fixture operation. Do not begin with combined provider preparation.

## Common mistakes

- Lowering from Goja builder values couples execution to one VM.
- Making one workflow node per corpus item at compile time prevents bounded runtime expansion.
- Storing live indexes in workflow state breaks restart and isolation.
- Reusing prepared artifacts without prompt/profile fingerprints can return semantically stale data.
- Letting Scraper understand RAG operator IDs reverses dependency ownership.
- Comparing scheduler artifacts in semantic parity tests creates false failures.

## API and schema reference

Every task spec should name exact port schemas. For example:

```go
workflowv3.TaskSpec{
    Key: "rag.chunks.create/v1",
    Inputs: map[string]string{"units": ragcontract.UnitSetManifestSchema},
    Outputs: map[string]string{"chunks": ragcontract.ChunkSetManifestSchema},
}
```

The task implementation decodes the unit-set manifest strictly, applies the canonical chunk operator, validates the resulting chunk manifest, and stores it through `ArtifactStore.Put`. It returns only the `ArtifactRef`. The descriptor module constructs an invocation with a workflow value matching the input schema. The lowering compiler uses the same task key directly. These three surfaces—spec, descriptor, and implementation—must remain in parity tests.

Lowering errors should identify the RAG node and operator:

```text
RAG_WORKFLOW_UNSUPPORTED_OPERATOR: pipeline.nodes[7] operator rerank.foo/v2 has no registered lowerer
RAG_WORKFLOW_PORT_SCHEMA: node embed input representations expected rag-representation-set/v2
```

## Review exercise

Choose one pipeline from `examples/rag-v2`, write its canonical node list, then write the expected Workflow V3 task/map/reduce list beside it. Mark where artifact publications occur and where restart can happen. If an intermediate value has no artifact schema or cannot survive process loss, the lowering design is incomplete.

## Intern onboarding checklist

The intern should compile one RAG JS example, identify every canonical pipeline node and port, explain static versus query-dependent work, locate the matching operator implementation, draw its proposed workflow fragment, and prove deterministic fixture parity for a provider-free pipeline.

## References

- Program: Researchctl ticket `EXPERIMENT-PLATFORM-CONVERGENCE`.
- `pkg/ragcontract/README.md`
- `pkg/gojamodules/rag/README.md`
- `pkg/ragcompiler/`
- `pkg/ragengine/`
- `pkg/ragoperators/`
- Scraper `pkg/workflowv3/` and `pkg/gojamodules/workflow/`.
- Depends on `EXPERIMENT-PLATFORM-SCRAPER-RUNNER`.
- Enables `RAG-GEPPETTO-WORKFLOW-OPERATIONS` and `RAG-V2-EXECUTION-CUTOVER`.
