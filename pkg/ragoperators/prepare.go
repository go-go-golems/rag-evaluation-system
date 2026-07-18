package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"unicode/utf8"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type unitOperator struct{ kind string }

func (o unitOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: o.kind, Version: "v1"}
}
func (o unitOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, _ *Environment) (map[string]any, error) {
	corpus, ok := inputs["corpus"].(Corpus)
	if !ok {
		return nil, fmt.Errorf("RAG_UNIT_INPUT: corpus")
	}
	records := append([]SourceRecord(nil), corpus.Records...)
	sort.Slice(records, func(i, j int) bool {
		if records[i].SessionID == records[j].SessionID {
			return records[i].Ordinal < records[j].Ordinal
		}
		return records[i].SessionID < records[j].SessionID
	})
	groups := [][]SourceRecord{}
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if o.kind == "transcript.units.agents-view-runs" && record.Role == "assistant" && len(groups) > 0 {
			last := groups[len(groups)-1]
			tail := last[len(last)-1]
			if tail.Role == "assistant" && tail.SessionID == record.SessionID && tail.Ordinal+1 == record.Ordinal {
				groups[len(groups)-1] = append(last, record)
				continue
			}
		}
		groups = append(groups, []SourceRecord{record})
	}
	units := make([]Unit, 0, len(groups))
	for _, group := range groups {
		text := ""
		ranges := []ragcontract.SourceRange{}
		ids := []string{}
		for _, r := range group {
			text += r.Text
			ranges = append(ranges, ragcontract.SourceRange{SourceID: r.ID, ByteStart: 0, ByteEnd: int64(len(r.Text)), OrdinalStart: r.Ordinal, OrdinalEnd: r.Ordinal + 1})
			ids = append(ids, r.ID)
		}
		digest, _ := ragcontract.Digest(text)
		idDigest, _ := ragcontract.Digest(struct {
			Kind   string
			IDs    []string
			Digest string
		}{o.kind, ids, digest})
		units = append(units, Unit{Record: ragcontract.UnitRecord{ID: "unit:" + idDigest[7:23], MemberSourceIDs: ids, Ranges: ranges, ContentDigest: digest}, Text: text, Records: group})
	}
	data, digest := materializationData(units)
	_, corpusDigest := materializationData(corpus.Records)
	manifest := ragcontract.UnitSetManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.UnitSetManifestSchema, Digest: digest, Parents: parent("corpus", corpusDigest, ragcontract.CorpusManifestSchema), Production: &ragcontract.Production{Operator: o.Ref(), Config: node.Config}}, UnitCount: int64(len(units)), IdentitySchema: "rag-unit-identity/v2"}
	for index := range units {
		units[index].ManifestDigest = digest
	}
	artifact := materializedArtifact("unit-set", "rag-unit-set", node.ID+".json", ragcontract.UnitSetManifestSchema, data, manifest)
	return map[string]any{"units": units, "artifact": artifact}, nil
}

type chunkOperator struct{ kind string }

func (o chunkOperator) Ref() ragcontract.OperatorRef {
	kind := o.kind
	if kind == "" {
		kind = "chunks.recursive"
	}
	return ragcontract.OperatorRef{Kind: kind, Version: "v1"}
}
func (o chunkOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, _ *Environment) (map[string]any, error) {
	units, ok := inputs["units"].([]Unit)
	if !ok {
		return nil, fmt.Errorf("RAG_CHUNK_INPUT: units")
	}
	var config struct {
		Size    int `json:"size"`
		Overlap int `json:"overlap"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	if o.kind == "chunks.identity" {
		config.Size = int(^uint(0) >> 1)
		config.Overlap = 0
	}
	if config.Size <= 0 {
		return nil, fmt.Errorf("RAG_CHUNK_SIZE")
	}
	if config.Overlap < 0 || config.Overlap >= config.Size {
		return nil, fmt.Errorf("RAG_CHUNK_OVERLAP")
	}
	chunks := []Chunk{}
	for _, unit := range units {
		runes := []rune(unit.Text)
		for start := 0; start < len(runes); {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			end := start + config.Size
			if end > len(runes) {
				end = len(runes)
			}
			text := string(runes[start:end])
			byteStart := int64(len(string(runes[:start])))
			byteEnd := byteStart + int64(len(text))
			digest, _ := ragcontract.Digest(text)
			identity, _ := ragcontract.Digest(struct {
				Unit       string
				Start, End int
				Digest     string
			}{unit.Record.ID, start, end, digest})
			ranges, citation := exactSourceRanges(unit, byteStart, byteEnd)
			record := ragcontract.ChunkRecord{ID: "chunk:" + identity[7:23], ParentUnitID: unit.Record.ID, ByteStart: byteStart, ByteEnd: byteEnd, LogicalStart: int64(start), LogicalEnd: int64(end), TextDigest: digest, Chunker: o.Ref(), Citation: citation}
			chunks = append(chunks, Chunk{Record: record, Text: text, Ranges: ranges})
			if end == len(runes) {
				break
			}
			start = end - config.Overlap
		}
	}
	data, digest := materializationData(chunks)
	unitDigest := ""
	if len(units) > 0 {
		unitDigest = units[0].ManifestDigest
	}
	manifest := ragcontract.ChunkSetManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ChunkSetManifestSchema, Digest: digest, Parents: parent("unit-set", unitDigest, ragcontract.UnitSetManifestSchema), Production: &ragcontract.Production{Operator: o.Ref(), Config: node.Config}}, ChunkCount: int64(len(chunks)), RangeUnit: "bytes-and-runes", UnicodePolicy: "utf8-rune-boundaries", EmptyInputPolicy: "omit"}
	for index := range chunks {
		chunks[index].ManifestDigest = digest
	}
	artifact := materializedArtifact("chunk-set", "rag-chunk-set", node.ID+".json", ragcontract.ChunkSetManifestSchema, data, manifest)
	return map[string]any{"chunks": chunks, "artifact": artifact}, nil
}

func exactSourceRanges(unit Unit, chunkStart, chunkEnd int64) ([]ragcontract.SourceRange, ragcontract.CitationRef) {
	ranges := []ragcontract.SourceRange{}
	offset := int64(0)
	for _, record := range unit.Records {
		recordStart := offset
		recordEnd := recordStart + int64(len(record.Text))
		offset = recordEnd
		start := max(chunkStart, recordStart)
		end := min(chunkEnd, recordEnd)
		if start >= end {
			continue
		}
		ranges = append(ranges, ragcontract.SourceRange{SourceID: record.ID, ByteStart: start - recordStart, ByteEnd: end - recordStart, OrdinalStart: record.Ordinal, OrdinalEnd: record.Ordinal + 1})
	}
	citation := ragcontract.CitationRef{}
	if len(ranges) > 0 {
		citation = ragcontract.CitationRef{SourceID: ranges[0].SourceID, ByteStart: ranges[0].ByteStart, ByteEnd: ranges[0].ByteEnd, OrdinalStart: ranges[0].OrdinalStart, OrdinalEnd: ranges[len(ranges)-1].OrdinalEnd}
	}
	return ranges, citation
}

func decodeConfig(raw json.RawMessage, target any) error {
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("RAG_OPERATOR_CONFIG: %w", err)
	}
	return nil
}
func validateUTF8Ranges(chunks []Chunk) error {
	for _, chunk := range chunks {
		if !utf8.ValidString(chunk.Text) {
			return fmt.Errorf("RAG_CHUNK_UTF8")
		}
		if chunk.Record.ByteEnd-chunk.Record.ByteStart != int64(len(chunk.Text)) {
			return fmt.Errorf("RAG_CHUNK_RANGE: %s", chunk.Record.ID)
		}
	}
	return nil
}
