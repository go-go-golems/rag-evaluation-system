package ragoperators

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type collapseOperator struct{ kind string }

func (o collapseOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: o.kind, Version: "v1"}
}
func (o collapseOperator) Execute(_ context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	var config struct{ Scope, Representative string }
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	groups := map[string][]RankedRecord{}
	if hits, ok := inputs["hits"].([]RankedRecord); ok {
		for _, hit := range hits {
			key := collapseKey(hit.Representation, config.Scope)
			groups[key] = append(groups[key], hit)
		}
	} else if parents, ok := inputs["parents"].([]RankedParent); ok {
		for _, parent := range parents {
			record := RankedRecord{Rank: parent.Rank, Representation: parent.Representative, Score: parent.Score}
			key := collapseKey(parent.Representative, config.Scope)
			groups[key] = append(groups[key], record)
		}
	} else {
		return nil, fmt.Errorf("RAG_COLLAPSE_INPUT")
	}
	result := []RankedParent{}
	traceGroups := []ragcontract.CollapseGroup{}
	for key, members := range groups {
		sort.Slice(members, func(i, j int) bool {
			if members[i].Score == members[j].Score {
				return members[i].Representation.Record.ID < members[j].Representation.Record.ID
			}
			return members[i].Score > members[j].Score
		})
		selected := members[0]
		result = append(result, RankedParent{Identity: ragcontract.CollapseIdentity{Scope: config.Scope, ID: key}, Score: selected.Score, Representative: selected.Representation, Members: members})
		trace := ragcontract.CollapseGroup{Key: ragcontract.CollapseIdentity{Scope: config.Scope, ID: key}, SelectedRepresentationID: selected.Representation.Record.ID}
		for _, member := range members {
			trace.Members = append(trace.Members, ragcontract.CollapseMember{RepresentationID: member.Representation.Record.ID, Rank: member.Rank, Score: member.Score})
		}
		traceGroups = append(traceGroups, trace)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Identity.ID < result[j].Identity.ID
		}
		return result[i].Score > result[j].Score
	})
	for i := range result {
		result[i].Rank = i + 1
	}
	for i := range traceGroups {
		for _, parent := range result {
			if traceGroups[i].Key == parent.Identity {
				traceGroups[i].Rank = parent.Rank
			}
		}
	}
	sort.Slice(traceGroups, func(i, j int) bool { return traceGroups[i].Rank < traceGroups[j].Rank })
	if env.Trace != nil {
		stage := "channel"
		if o.kind == "collapse.final" {
			stage = "final"
		}
		env.Trace.Collapses = append(env.Trace.Collapses, ragcontract.CollapseTrace{Stage: stage, Operator: o.Ref(), Groups: traceGroups})
	}
	return map[string]any{"parents": result}, nil
}
func collapseKey(rep Representation, scope string) string {
	if scope == "chunk" {
		return rep.Record.ParentChunkID
	}
	return rep.Record.ParentUnitID
}

type fusionOperator struct{}

func (fusionOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "fusion.weighted-rrf", Version: "v1"}
}
func (fusionOperator) Execute(_ context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	var config struct {
		RankConstant         int                `json:"rankConstant"`
		Weights              map[string]float64 `json:"weights"`
		MissingChannelPolicy string             `json:"missingChannelPolicy"`
		TieBreak             string             `json:"tieBreak"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	if config.RankConstant <= 0 {
		return nil, fmt.Errorf("RAG_RRF_CONSTANT")
	}
	if config.MissingChannelPolicy == "reject" {
		for channel := range config.Weights {
			if _, ok := inputs["channel."+channel]; !ok {
				return nil, fmt.Errorf("RAG_RRF_MISSING_CHANNEL: %s", channel)
			}
		}
	}
	byKey := map[string]*RankedParent{}
	for port, value := range inputs {
		parents, ok := value.([]RankedParent)
		if !ok {
			return nil, fmt.Errorf("RAG_RRF_INPUT: %s", port)
		}
		channel := port
		if len(channel) > 8 && channel[:8] == "channel." {
			channel = channel[8:]
		}
		weight := 1.0
		if w, ok := config.Weights[channel]; ok {
			weight = w
		}
		if weight <= 0 {
			return nil, fmt.Errorf("RAG_RRF_WEIGHT: %s", channel)
		}
		for _, parent := range parents {
			entry := byKey[parent.Identity.ID]
			if entry == nil {
				cloned := parent
				cloned.Score = 0
				cloned.Contributions = nil
				entry = &cloned
				byKey[parent.Identity.ID] = entry
			} else {
				seen := map[string]bool{}
				for _, member := range entry.Members {
					seen[member.Channel+"\x00"+member.Representation.Record.ID] = true
				}
				for _, member := range parent.Members {
					key := member.Channel + "\x00" + member.Representation.Record.ID
					if !seen[key] {
						seen[key] = true
						entry.Members = append(entry.Members, member)
					}
				}
			}
			contribution := weight / float64(config.RankConstant+parent.Rank)
			entry.Score += contribution
			entry.Contributions = append(entry.Contributions, ragcontract.FusionContribution{Channel: channel, Rank: parent.Rank, Weight: weight, Value: contribution})
			if contribution > bestContribution(entry.Contributions[:len(entry.Contributions)-1]) || (contribution == bestContribution(entry.Contributions[:len(entry.Contributions)-1]) && parent.Representative.Record.ID < entry.Representative.Record.ID) {
				entry.Representative = parent.Representative
			}
		}
	}
	result := make([]RankedParent, 0, len(byKey))
	for _, entry := range byKey {
		sort.Slice(entry.Contributions, func(i, j int) bool { return entry.Contributions[i].Channel < entry.Contributions[j].Channel })
		result = append(result, *entry)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Identity.ID < result[j].Identity.ID
		}
		return result[i].Score > result[j].Score
	})
	trace := ragcontract.FusionTrace{Operator: fusionOperator{}.Ref(), MissingChannelPolicy: config.MissingChannelPolicy, TieBreak: "scoreThenCollapseId"}
	for i := range result {
		result[i].Rank = i + 1
		trace.Results = append(trace.Results, ragcontract.FusionResult{Rank: i + 1, Identity: result[i].Identity, Contributions: result[i].Contributions, Score: result[i].Score})
	}
	if env.Trace != nil {
		env.Trace.Fusion = &trace
	}
	return map[string]any{"parents": result}, nil
}
func bestContribution(values []ragcontract.FusionContribution) float64 {
	best := 0.0
	for _, v := range values {
		if v.Value > best {
			best = v.Value
		}
	}
	return best
}

type hydrateOperator struct{}

func (hydrateOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "hydrate.source-evidence", Version: "v1"}
}
func (hydrateOperator) Execute(_ context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	parents, ok := inputs["parents"].([]RankedParent)
	if !ok {
		return nil, fmt.Errorf("RAG_HYDRATE_PARENTS")
	}
	chunks, ok := inputs["chunks"].([]Chunk)
	if !ok {
		return nil, fmt.Errorf("RAG_HYDRATE_CHUNKS")
	}
	byID := map[string]Chunk{}
	for _, chunk := range chunks {
		byID[chunk.Record.ID] = chunk
	}
	var config struct {
		Results int `json:"results"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	if config.Results > 0 && len(parents) > config.Results {
		parents = parents[:config.Results]
	}
	result := []Evidence{}
	trace := ragcontract.HydrationTrace{Operator: hydrateOperator{}.Ref()}
	for _, parent := range parents {
		chunk, exists := byID[parent.Representative.Record.ParentChunkID]
		if !exists {
			return nil, fmt.Errorf("RAG_HYDRATE_LINEAGE: %s", parent.Representative.Record.ParentChunkID)
		}
		evidence := Evidence{Rank: parent.Rank, Collapse: parent.Identity, Chunk: chunk, Score: parent.Score, Contributions: parent.Contributions, Matched: append([]RankedRecord(nil), parent.Members...)}
		result = append(result, evidence)
		identity := ragcontract.EvidenceIdentity{ChunkID: chunk.Record.ID, Digest: chunk.Record.TextDigest, Citation: chunk.Record.Citation}
		trace.Candidates = append(trace.Candidates, ragcontract.HydrationCandidate{Collapse: parent.Identity, Evidence: identity, Contribution: bestContribution(parent.Contributions)})
		trace.Selected = append(trace.Selected, identity)
	}
	if env.Trace != nil {
		env.Trace.Hydration = &trace
	}
	return map[string]any{"evidence": result}, nil
}

type rerankOperator struct{}

func (rerankOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "rerank.cross-encoder", Version: "v1"}
}
func (rerankOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	evidence, ok := inputs["evidence"].([]Evidence)
	if !ok {
		return nil, fmt.Errorf("RAG_RERANK_INPUT")
	}
	if env.Reranker == nil {
		return nil, fmt.Errorf("RAG_RERANKER_UNAVAILABLE")
	}
	var config struct {
		Model                                   string `json:"model"`
		CandidateCount                          int    `json:"candidateCount"`
		Results                                 int    `json:"results"`
		Truncation, Tokenization, InputTemplate string
		TimeoutMilliseconds                     int64 `json:"timeoutMilliseconds"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	modelManifest, err := resolveModel(env, config.Model)
	if err != nil {
		return nil, err
	}
	if config.Truncation != modelManifest.Truncation || config.Tokenization != modelManifest.Tokenization {
		return nil, fmt.Errorf("RAG_RERANK_MODEL_POLICY: configured truncation/tokenization do not match manifest")
	}
	if config.CandidateCount < len(evidence) {
		evidence = evidence[:config.CandidateCount]
	}
	query := env.QueryText
	if config.TimeoutMilliseconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(config.TimeoutMilliseconds)*time.Millisecond)
		defer cancel()
	}
	scores, err := env.Reranker.Rerank(ctx, RerankRequest{Model: modelManifest.ModelID, InputTemplate: config.InputTemplate, Truncation: config.Truncation, Tokenization: config.Tokenization, Query: query, Candidates: evidence, Results: config.Results})
	if err != nil {
		return nil, fmt.Errorf("RAG_RERANK_FAILED: %w", err)
	}
	byID := map[string]float64{}
	for _, score := range scores {
		if _, duplicate := byID[score.ChunkID]; duplicate {
			return nil, fmt.Errorf("RAG_RERANK_DUPLICATE: %s", score.ChunkID)
		}
		byID[score.ChunkID] = score.Score
	}
	if len(byID) != len(evidence) {
		return nil, fmt.Errorf("RAG_RERANK_INCOMPLETE: got %d want %d", len(byID), len(evidence))
	}
	for i := range evidence {
		score := byID[evidence[i].Chunk.Record.ID]
		evidence[i].RerankerScore = &score
	}
	sort.Slice(evidence, func(i, j int) bool {
		if *evidence[i].RerankerScore == *evidence[j].RerankerScore {
			return evidence[i].Chunk.Record.ID < evidence[j].Chunk.Record.ID
		}
		return *evidence[i].RerankerScore > *evidence[j].RerankerScore
	})
	if config.Results > 0 && len(evidence) > config.Results {
		evidence = evidence[:config.Results]
	}
	trace := ragcontract.RerankingTrace{Operator: rerankOperator{}.Ref(), ModelManifestDigest: modelManifest.Digest, InputPolicy: "source-evidence", InputTemplate: config.InputTemplate, Truncation: config.Truncation, Tokenization: config.Tokenization, CandidateCount: config.CandidateCount, ResultsLimit: config.Results, TimeoutMilliseconds: config.TimeoutMilliseconds}
	for i := range evidence {
		before := evidence[i].Rank
		evidence[i].Rank = i + 1
		trace.Entries = append(trace.Entries, ragcontract.RerankingEntry{Evidence: ragcontract.EvidenceIdentity{ChunkID: evidence[i].Chunk.Record.ID, Digest: evidence[i].Chunk.Record.TextDigest, Citation: evidence[i].Chunk.Record.Citation}, BeforeRank: before, AfterRank: i + 1, RetrievalScore: evidence[i].Score, RerankerScore: *evidence[i].RerankerScore})
	}
	if env.Trace != nil {
		env.Trace.Reranking = &trace
	}
	return map[string]any{"evidence": evidence}, nil
}

type answerOperator struct{}

func (answerOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "generate.answer", Version: "v1"}
}
func (answerOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	evidence, ok := inputs["evidence"].([]Evidence)
	if !ok {
		return nil, fmt.Errorf("RAG_ANSWER_INPUT")
	}
	if env.Generator == nil {
		return nil, fmt.Errorf("RAG_GENERATOR_UNAVAILABLE: answer")
	}
	var config struct {
		Model, Prompt, Citations string
		ContextBudgetTokens      int
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	modelManifest, err := resolveModel(env, config.Model)
	if err != nil {
		return nil, err
	}
	promptManifest, err := resolvePrompt(env, config.Prompt)
	if err != nil {
		return nil, err
	}
	result, err := env.Generator.Generate(ctx, GenerationRequest{Kind: "answer", Model: modelManifest.ModelID, Prompt: promptManifest.PromptID, Evidence: evidence})
	if err != nil {
		return nil, fmt.Errorf("RAG_ANSWER_FAILED: %w", err)
	}
	if config.Citations == "source" && len(evidence) > 0 && len(result.CitationChunkIDs) == 0 && !result.Abstained {
		return nil, fmt.Errorf("RAG_ANSWER_CITATION_REQUIRED")
	}
	valid := map[string]bool{}
	for _, item := range evidence {
		valid[item.Chunk.Record.ID] = true
	}
	for _, id := range result.CitationChunkIDs {
		if !valid[id] {
			return nil, fmt.Errorf("RAG_ANSWER_CITATION: %s", id)
		}
	}
	if config.Citations == "required" && len(result.CitationChunkIDs) == 0 {
		return nil, fmt.Errorf("RAG_ANSWER_CITATIONS_REQUIRED")
	}
	answer := Answer{Text: result.Text, CitationChunkIDs: result.CitationChunkIDs, FinishReason: result.FinishReason, Abstained: result.Abstained, InputTokens: result.InputTokens, OutputTokens: result.OutputTokens}
	env.Usage.InputTokens += result.InputTokens
	env.Usage.OutputTokens += result.OutputTokens
	if env.Trace != nil {
		identities := make([]ragcontract.EvidenceIdentity, len(evidence))
		for index, item := range evidence {
			identities[index] = ragcontract.EvidenceIdentity{ChunkID: item.Chunk.Record.ID, Digest: item.Chunk.Record.TextDigest, Citation: item.Chunk.Record.Citation}
		}
		inputDigest, _ := ragcontract.Digest(identities)
		outputDigest, _ := ragcontract.Digest(answer)
		env.Trace.Generation = &ragcontract.GenerationTrace{Operator: answerOperator{}.Ref(), ModelManifestDigest: modelManifest.Digest, PromptManifestDigest: promptManifest.Digest, Evidence: identities, InputArtifactDigest: inputDigest, OutputArtifactDigest: outputDigest, FinishReason: result.FinishReason, CitationsValid: true}
	}
	return map[string]any{"answer": answer}, nil
}
