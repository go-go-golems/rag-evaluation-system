# `require("rag")` v2 authoring module

The module exposes pure Go-backed descriptors and fluent builders. JavaScript composes intent; `pkg/ragmodel`, `pkg/ragcompiler`, and `pkg/ragcontract` own semantics, normalization, validation, and wire values.

```javascript
const rag = require("rag");
const pipeline = rag.pipeline("raw", p => p
  .corpus(rag.inputs.corpus("corpus"))
  .units(rag.units.identity())
  .chunks(rag.chunks.recursive({ maxRunes: 800 }))
  .representations(rag.representations.raw("raw"))
  .index("representations", rag.indexes.bleveMulti({ lexical: true })));
console.log(rag.explain(pipeline));
```

Go embeddings register `NewRegistrar()` with an `engine.RuntimeFactoryBuilder`, or import the module and use `NewLoader()` with a goja-nodejs registry. Generated xgoja binaries register `pkg/xgoja/providers/rag`.

### Decision: Go-backed hidden values and immediate configurators

- **Context:** authoring needs fluent typed values, reusable fragments, and nested lambdas, but functions and Goja values cannot cross into normalized IR or workers.
- **Options considered:** plain mutable JavaScript maps, serialized callback source, or Go-backed values attached by private symbols.
- **Decision:** descriptors, fragments, pipelines, query plans, products, studies, variants, and factor references are Go-owned pointers attached through runtime-private `goja.Symbol` keys. Configurator callbacks run exactly once while their parent factory/method is executing.
- **Rationale:** symbols are absent from `Object.keys` and JSON; Go builders enforce type boundaries; compiled targets contain only data.
- **Consequences:** values are runtime-local and cannot be forged by plain maps. Reuse happens through explicit descriptors/fragments, while portable output is canonical JSON.
- **Status:** accepted.

### Decision: pure authoring capability boundary

- **Context:** safe project loading must not trigger providers, files, databases, indexes, or laboratory lifecycle.
- **Decision:** this package imports only Goja adapter APIs and pure RAG model/contract/compiler packages. It exposes `validate`, `explain`, `compileProduct`, `compileStudy`, and `preview`; it exposes no execution method.
- **Consequences:** compilation requires explicit immutable binding data but performs no I/O. Execution belongs to later RAG-owned CLI/worker layers.
- **Status:** accepted.
