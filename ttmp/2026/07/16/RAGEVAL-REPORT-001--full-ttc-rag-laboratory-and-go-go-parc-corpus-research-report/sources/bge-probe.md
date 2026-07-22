# llama.cpp BGE reranker probe results

Date: 2026-07-16

## Environment

| Property | Observed value |
| --- | --- |
| Mac host | `mimimi-2.local` |
| Server binary | `/Applications/Ollama.app/Contents/Resources/llama-server` |
| Server version | `1 (cb295bf59)`, AppleClang 21.0.0.21000099, Darwin arm64 |
| Model blob | `qllama/bge-reranker-v2-m3:q4_k_m`, Ollama model blob SHA-256 `10a8e2b5…a9c44cd` |
| Server flags | `--embedding --pooling rank --rerank --host 127.0.0.1 --port 8012` |
| Workstation tunnel | tmux `rag-reranker-mimimi`, `127.0.0.1:18012 -> mimimi-2.local:127.0.0.1:8012` |
| Request route | `POST /v1/rerank` |

The server health endpoint returned `{"status":"ok"}`. Startup log confirms
the model loaded and that it listens only on the Mac loopback interface.

## Full ranking probe

The script `01-probe-llamacpp-bge-reranker.sh` submitted three documents and
requested `top_n: 3`. The relevant payroll-adjustment candidate was input index
zero; the cypress and weather documents were indexes one and two.

```json
{
  "model": "qllama/bge-reranker-v2-m3:q4_k_m",
  "object": "list",
  "usage": {"prompt_tokens": 96, "total_tokens": 96},
  "results": [
    {"index": 0, "relevance_score": -3.32784366607666},
    {"index": 1, "relevance_score": -9.837879180908203},
    {"index": 2, "relevance_score": -11.012685775756836}
  ]
}
```

Wall-clock request time through the SSH tunnel was `0.202822` seconds. Scores
are descending and negative. The Go adapter must accept any finite score and
sort descending; it must not assume scores are probabilities or constrained to
the interval `[0,1]`.

## Truncated probe

With the same request and `top_n: 2`, the response returned only indices zero
and one. This confirms that `top_n` changes response cardinality. The initial
adapter should request a score for every submitted candidate by setting
`top_n == len(documents)`; it can then enforce a complete one-to-one response
before it applies the laboratory's final result limit.

## Decoder contract derived from the probe

```text
request.documents[i] -> response.results[*].index == i -> immutable chunk ID

required validation:
  index in [0, len(documents))
  no duplicate indexes
  all indexes returned when top_n == len(documents)
  relevance_score is finite
  response is sorted by score only after mapping to candidate identity
```

The server log showed three independent tasks/slots for the three documents.
Response order was already score order in this probe, but the adapter must use
the returned `index` rather than infer identity from completion or result order.
