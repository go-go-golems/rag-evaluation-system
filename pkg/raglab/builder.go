package raglab

import (
	"sort"
	"strconv"
)

type Fragment struct {
	Name      string
	Configure func(*ExperimentBuilder)
}

func NewFragment(name string, configure func(*ExperimentBuilder)) Fragment {
	return Fragment{Name: name, Configure: configure}
}

type ExperimentBuilder struct {
	name        string
	inputs      InputSpec
	provenance  Provenance
	retrieval   RetrievalPlan
	metrics     MetricsPlan
	buildIssues ValidationReport
}

func NewExperiment(name string) *ExperimentBuilder {
	return &ExperimentBuilder{
		name:       name,
		inputs:     InputSpec{Representations: []RepresentationSpec{{Name: "raw", Kind: RawChunksRepresentation}}},
		provenance: Provenance{Tags: map[string]string{}},
		retrieval:  RetrievalPlan{Collapse: CollapseNone, Results: 10},
	}
}

func (b *ExperimentBuilder) Use(fragment Fragment) *ExperimentBuilder {
	if fragment.Name == "" {
		b.buildIssues.add("RAG_INVALID_FRAGMENT", "$.provenance.fragments", "fragment name is required")
		return b
	}
	for _, name := range b.provenance.Fragments {
		if name == fragment.Name {
			return b
		}
	}
	b.provenance.Fragments = append(b.provenance.Fragments, fragment.Name)
	if fragment.Configure == nil {
		b.buildIssues.add("RAG_INVALID_FRAGMENT", "$.provenance.fragments", "fragment configurator is required")
		return b
	}
	fragment.Configure(b)
	return b
}

func (b *ExperimentBuilder) Corpus(ref ArtifactRef) *ExperimentBuilder {
	b.setRequiredArtifact("corpus snapshot", "$.inputs.corpusSnapshot", CorpusSnapshotArtifact, ref, &b.inputs.CorpusSnapshot)
	return b
}

func (b *ExperimentBuilder) Chunks(ref ArtifactRef) *ExperimentBuilder {
	b.setRequiredArtifact("chunk set", "$.inputs.chunkSet", ChunkSetArtifact, ref, &b.inputs.ChunkSet)
	return b
}

func (b *ExperimentBuilder) BM25(ref ArtifactRef) *ExperimentBuilder {
	var existing ArtifactRef
	if b.inputs.BM25Index != nil {
		existing = *b.inputs.BM25Index
	}
	b.setOptionalArtifact("BM25 index", "$.inputs.bm25Index", BM25IndexArtifact, ref, existing, func(value ArtifactRef) { b.inputs.BM25Index = &value })
	return b
}

func (b *ExperimentBuilder) Embeddings(ref ArtifactRef) *ExperimentBuilder {
	var existing ArtifactRef
	if b.inputs.EmbeddingSet != nil {
		existing = *b.inputs.EmbeddingSet
	}
	b.setOptionalArtifact("embedding set", "$.inputs.embeddingSet", EmbeddingSetArtifact, ref, existing, func(value ArtifactRef) { b.inputs.EmbeddingSet = &value })
	return b
}

func (b *ExperimentBuilder) Evaluation(ref ArtifactRef) *ExperimentBuilder {
	b.setRequiredArtifact("evaluation dataset", "$.inputs.evaluationDataset", EvaluationDatasetArtifact, ref, &b.inputs.EvaluationDataset)
	return b
}

func (b *ExperimentBuilder) Note(text string) *ExperimentBuilder {
	if text == "" {
		b.buildIssues.add("RAG_INVALID_NOTE", "$.provenance.notes", "note must not be empty")
		return b
	}
	b.provenance.Notes = append(b.provenance.Notes, text)
	return b
}

func (b *ExperimentBuilder) Tag(name, value string) *ExperimentBuilder {
	if name == "" {
		b.buildIssues.add("RAG_INVALID_TAG", "$.provenance.tags", "tag name is required")
		return b
	}
	if existing, ok := b.provenance.Tags[name]; ok && existing != value {
		b.buildIssues.add("RAG_CONFLICTING_TAG", "$.provenance.tags."+name, "tag already has a different value")
		return b
	}
	b.provenance.Tags[name] = value
	return b
}

func (b *ExperimentBuilder) Representations(configure func(*RepresentationBuilder)) *ExperimentBuilder {
	if configure == nil {
		b.buildIssues.add("RAG_INVALID_REPRESENTATIONS", "$.inputs.representations", "representation configurator is required")
		return b
	}
	configure(&RepresentationBuilder{experiment: b})
	return b
}

func (b *ExperimentBuilder) Retrieval(configure func(*RetrievalBuilder)) *ExperimentBuilder {
	if configure == nil {
		b.buildIssues.add("RAG_INVALID_RETRIEVAL", "$.retrieval", "retrieval configurator is required")
		return b
	}
	configure(&RetrievalBuilder{experiment: b, plan: &b.retrieval})
	return b
}

func (b *ExperimentBuilder) Metrics(configure func(*MetricsBuilder)) *ExperimentBuilder {
	if configure == nil {
		b.buildIssues.add("RAG_INVALID_METRICS", "$.metrics", "metrics configurator is required")
		return b
	}
	configure(&MetricsBuilder{experiment: b, plan: &b.metrics})
	return b
}

func (b *ExperimentBuilder) Validate() ValidationReport {
	report := b.buildIssues
	if b.name == "" {
		report.add("RAG_INVALID_NAME", "$.name", "experiment name is required")
	}
	validateRequiredRef(&report, "$.inputs.corpusSnapshot", CorpusSnapshotArtifact, b.inputs.CorpusSnapshot)
	validateRequiredRef(&report, "$.inputs.chunkSet", ChunkSetArtifact, b.inputs.ChunkSet)
	validateRequiredRef(&report, "$.inputs.evaluationDataset", EvaluationDatasetArtifact, b.inputs.EvaluationDataset)
	if len(b.retrieval.Channels) == 0 {
		report.add("RAG_MISSING_RETRIEVAL", "$.retrieval.channels", "at least one retrieval channel is required")
	}
	if b.retrieval.Results <= 0 {
		report.add("RAG_INVALID_RESULTS", "$.retrieval.results", "result count must be positive")
	}
	if b.retrieval.Collapse != CollapseNone && b.retrieval.Collapse != CollapseParentChunk && b.retrieval.Collapse != CollapseDocument {
		report.add("RAG_INVALID_COLLAPSE", "$.retrieval.collapse", "collapse must be none, parentChunk, or document")
	}
	representations := map[string]RepresentationSpec{}
	for i, representation := range b.inputs.Representations {
		path := "$.inputs.representations[" + itoa(i) + "]"
		if representation.Name == "" {
			report.add("RAG_INVALID_REPRESENTATION", path+".name", "representation name is required")
			continue
		}
		if _, exists := representations[representation.Name]; exists {
			report.add("RAG_DUPLICATE_REPRESENTATION", path+".name", "representation name must be unique")
			continue
		}
		if representation.Kind != RawChunksRepresentation && (representation.ArtifactID == "" || representation.Parent != "sourceChunk") {
			report.add("RAG_INVALID_REPRESENTATION", path, "materialized representations require an artifact and sourceChunk parent")
		}
		representations[representation.Name] = representation
	}
	seenChannels := map[string]bool{}
	for i, channel := range b.retrieval.Channels {
		path := "$.retrieval.channels[" + itoa(i) + "]"
		if channel.Name == "" {
			report.add("RAG_INVALID_CHANNEL", path+".name", "channel name is required")
		} else if seenChannels[channel.Name] {
			report.add("RAG_DUPLICATE_CHANNEL", path+".name", "channel name must be unique")
		} else {
			seenChannels[channel.Name] = true
		}
		if channel.Backend != BM25Backend && channel.Backend != VectorBackend {
			report.add("RAG_INVALID_CHANNEL", path+".backend", "channel must select exactly one backend")
		}
		if channel.TopK <= 0 {
			report.add("RAG_INVALID_CHANNEL", path+".topK", "channel topK must be positive")
		}
		if b.retrieval.Results > 0 && channel.TopK > 0 && channel.TopK < b.retrieval.Results {
			report.add("RAG_INVALID_CHANNEL", path+".topK", "channel topK must be at least final result count")
		}
		if _, exists := representations[channel.Representation]; !exists {
			report.add("RAG_UNKNOWN_REPRESENTATION", path+".representation", "channel references an undeclared representation")
		}
		if channel.Backend == BM25Backend && b.inputs.BM25Index == nil {
			report.add("RAG_MISSING_BM25", path+".backend", "BM25 channel requires a BM25 index")
		}
		if channel.Backend == VectorBackend && b.inputs.EmbeddingSet == nil {
			report.add("RAG_MISSING_EMBEDDINGS", path+".backend", "vector channel requires an embedding set")
		}
	}
	if len(b.retrieval.Channels) > 1 && b.retrieval.Fusion == nil {
		report.add("RAG_MISSING_FUSION", "$.retrieval.fusion", "multiple channels require an explicit fusion policy")
	}
	if b.retrieval.Fusion != nil {
		if b.retrieval.Fusion.Kind != "rrf" {
			report.add("RAG_INVALID_FUSION", "$.retrieval.fusion.kind", "only rrf fusion is supported")
		}
		if b.retrieval.Fusion.RankConstant <= 0 {
			report.add("RAG_INVALID_FUSION", "$.retrieval.fusion.rankConstant", "RRF rank constant must be positive")
		}
		for channel, weight := range b.retrieval.Fusion.Weights {
			if !seenChannels[channel] {
				report.add("RAG_INVALID_FUSION", "$.retrieval.fusion.weights."+channel, "weight references an unknown channel")
			}
			if weight <= 0 {
				report.add("RAG_INVALID_FUSION", "$.retrieval.fusion.weights."+channel, "channel weight must be positive")
			}
		}
	}
	if reranking := b.retrieval.Reranking; reranking != nil {
		if reranking.Kind != CrossEncoderReranking {
			report.add("RAG_INVALID_RERANKING", "$.retrieval.reranking.kind", "only crossEncoder reranking is supported")
		}
		if reranking.Model == "" {
			report.add("RAG_INVALID_RERANKING", "$.retrieval.reranking.model", "reranker model is required")
		}
		if reranking.CandidateCount <= 0 {
			report.add("RAG_INVALID_RERANKING", "$.retrieval.reranking.candidateCount", "candidate count must be positive")
		}
		if reranking.Results <= 0 || reranking.Results > reranking.CandidateCount {
			report.add("RAG_INVALID_RERANKING", "$.retrieval.reranking.results", "reranking results must be positive and no greater than candidate count")
		}
		if b.retrieval.Results > reranking.Results {
			report.add("RAG_INVALID_RERANKING", "$.retrieval.reranking.results", "reranking results must be at least final result count")
		}
	}
	validateFilter(&report, "$.retrieval.filter", b.retrieval.Filter)
	for i, channel := range b.retrieval.Channels {
		validateFilter(&report, "$.retrieval.channels["+itoa(i)+"].filter", channel.Filter)
	}
	validateMetrics(&report, b.metrics, b.retrieval.Results)
	report.Normalize()
	return report
}

func (b *ExperimentBuilder) Build() (ExperimentSpecification, error) {
	report := b.Validate()
	if !report.OK() {
		return ExperimentSpecification{}, &ValidationError{Report: report}
	}
	spec := ExperimentSpecification{
		SchemaVersion: AuthoringSchemaVersion,
		Name:          b.name,
		Provenance:    normalizeProvenance(b.provenance),
		Inputs:        normalizeInputs(b.inputs),
		Retrieval:     normalizeRetrieval(b.retrieval),
		Metrics:       normalizeMetrics(b.metrics),
	}
	return spec, nil
}

func (b *ExperimentBuilder) setRequiredArtifact(label, path string, expected ArtifactKind, incoming ArtifactRef, target *ArtifactRef) {
	if incoming.Kind != expected || incoming.ID == "" {
		b.buildIssues.add("RAG_INVALID_ARTIFACT", path, "expected non-empty "+string(expected)+" for "+label)
		return
	}
	if target.ID != "" && target.ID != incoming.ID {
		b.buildIssues.add("RAG_CONFLICTING_FRAGMENT", path, label+" is already set to a different artifact")
		return
	}
	*target = incoming
}

func (b *ExperimentBuilder) setOptionalArtifact(label, path string, expected ArtifactKind, incoming, existing ArtifactRef, set func(ArtifactRef)) {
	if incoming.Kind != expected || incoming.ID == "" {
		b.buildIssues.add("RAG_INVALID_ARTIFACT", path, "expected non-empty "+string(expected)+" for "+label)
		return
	}
	if existing.ID != "" && existing.ID != incoming.ID {
		b.buildIssues.add("RAG_CONFLICTING_FRAGMENT", path, label+" is already set to a different artifact")
		return
	}
	set(incoming)
}

type RepresentationBuilder struct{ experiment *ExperimentBuilder }

func (b *RepresentationBuilder) RawChunks(name string) *RepresentationBuilder {
	if name == "" {
		name = "raw"
	}
	b.experiment.inputs.Representations = appendRepresentation(b.experiment.inputs.Representations, RepresentationSpec{Name: name, Kind: RawChunksRepresentation})
	return b
}

func (b *RepresentationBuilder) Summaries(name string, artifact ArtifactRef) *RepresentationBuilder {
	b.experiment.inputs.Representations = appendRepresentation(b.experiment.inputs.Representations, RepresentationSpec{Name: name, Kind: SummaryRepresentation, ArtifactID: artifact.ID, Parent: "sourceChunk"})
	if artifact.Kind != RepresentationSetArtifact {
		b.experiment.buildIssues.add("RAG_INVALID_ARTIFACT", "$.inputs.representations", "summary representation requires a representationSet artifact")
	}
	return b
}

func (b *RepresentationBuilder) Questions(name string, artifact ArtifactRef) *RepresentationBuilder {
	b.experiment.inputs.Representations = appendRepresentation(b.experiment.inputs.Representations, RepresentationSpec{Name: name, Kind: QuestionRepresentation, ArtifactID: artifact.ID, Parent: "sourceChunk"})
	if artifact.Kind != RepresentationSetArtifact {
		b.experiment.buildIssues.add("RAG_INVALID_ARTIFACT", "$.inputs.representations", "question representation requires a representationSet artifact")
	}
	return b
}

type RetrievalBuilder struct {
	experiment *ExperimentBuilder
	plan       *RetrievalPlan
}

func (b *RetrievalBuilder) Channel(name string, configure func(*ChannelBuilder)) *RetrievalBuilder {
	channel := ChannelSpec{Name: name, Representation: "raw"}
	b.plan.Channels = append(b.plan.Channels, channel)
	index := len(b.plan.Channels) - 1
	if configure == nil {
		b.experiment.buildIssues.add("RAG_INVALID_CHANNEL", "$.retrieval.channels["+itoa(index)+"]", "channel configurator is required")
		return b
	}
	configure(&ChannelBuilder{experiment: b.experiment, channel: &b.plan.Channels[index]})
	return b
}

func (b *RetrievalBuilder) Filter(configure func(*FilterBuilder)) *RetrievalBuilder {
	if configure == nil {
		b.experiment.buildIssues.add("RAG_INVALID_FILTER", "$.retrieval.filter", "filter configurator is required")
		return b
	}
	configure(&FilterBuilder{experiment: b.experiment, filter: &b.plan.Filter, path: "$.retrieval.filter"})
	return b
}

func (b *RetrievalBuilder) FuseRRF(rankConstant int) *RetrievalBuilder {
	if b.plan.Fusion != nil && b.plan.Fusion.Kind != "rrf" {
		b.experiment.buildIssues.add("RAG_CONFLICTING_FUSION", "$.retrieval.fusion", "fusion is already configured")
		return b
	}
	if rankConstant == 0 {
		rankConstant = 60
	}
	b.plan.Fusion = &FusionSpec{Kind: "rrf", RankConstant: rankConstant, Weights: map[string]float64{}}
	return b
}

func (b *RetrievalBuilder) Weight(channel string, value float64) *RetrievalBuilder {
	if b.plan.Fusion == nil {
		b.experiment.buildIssues.add("RAG_INVALID_FUSION", "$.retrieval.fusion", "configure RRF before channel weights")
		return b
	}
	b.plan.Fusion.Weights[channel] = value
	return b
}

// RerankCrossEncoder configures a bounded cross-encoder stage. The model name
// participates in canonical experiment identity; provider URL and credentials
// are supplied separately when execution begins.
func (b *RetrievalBuilder) RerankCrossEncoder(model string, candidateCount, results int) *RetrievalBuilder {
	if b.plan.Reranking != nil {
		b.experiment.buildIssues.add("RAG_CONFLICTING_RERANKING", "$.retrieval.reranking", "reranking is already configured")
		return b
	}
	b.plan.Reranking = &RerankingSpec{Kind: CrossEncoderReranking, Model: model, CandidateCount: candidateCount, Results: results}
	return b
}

func (b *RetrievalBuilder) Collapse(scope CollapseScope) *RetrievalBuilder {
	b.plan.Collapse = scope
	return b
}
func (b *RetrievalBuilder) Results(count int) *RetrievalBuilder { b.plan.Results = count; return b }

type ChannelBuilder struct {
	experiment *ExperimentBuilder
	channel    *ChannelSpec
}

func (b *ChannelBuilder) BM25() *ChannelBuilder {
	if b.channel.Backend != "" && b.channel.Backend != BM25Backend {
		b.experiment.buildIssues.add("RAG_CHANNEL_BACKEND_CONFLICT", "$.retrieval.channels."+b.channel.Name+".backend", "channel already selected a different backend")
		return b
	}
	b.channel.Backend = BM25Backend
	return b
}
func (b *ChannelBuilder) Vector() *ChannelBuilder {
	if b.channel.Backend != "" && b.channel.Backend != VectorBackend {
		b.experiment.buildIssues.add("RAG_CHANNEL_BACKEND_CONFLICT", "$.retrieval.channels."+b.channel.Name+".backend", "channel already selected a different backend")
		return b
	}
	b.channel.Backend = VectorBackend
	return b
}
func (b *ChannelBuilder) Representation(name string) *ChannelBuilder {
	b.channel.Representation = name
	return b
}
func (b *ChannelBuilder) TopK(count int) *ChannelBuilder { b.channel.TopK = count; return b }
func (b *ChannelBuilder) Filter(configure func(*FilterBuilder)) *ChannelBuilder {
	if configure == nil {
		b.experiment.buildIssues.add("RAG_INVALID_FILTER", "$.retrieval.channels."+b.channel.Name+".filter", "filter configurator is required")
		return b
	}
	configure(&FilterBuilder{experiment: b.experiment, filter: &b.channel.Filter, path: "$.retrieval.channels." + b.channel.Name + ".filter"})
	return b
}

type FilterBuilder struct {
	experiment *ExperimentBuilder
	filter     *FilterSpec
	path       string
}

func (b *FilterBuilder) SourceIDs(ids ...string) *FilterBuilder {
	b.filter.SourceIDs = append(b.filter.SourceIDs, ids...)
	return b
}
func (b *FilterBuilder) DocumentIDs(ids ...string) *FilterBuilder {
	b.filter.DocumentIDs = append(b.filter.DocumentIDs, ids...)
	return b
}
func (b *FilterBuilder) ContentTypes(types ...string) *FilterBuilder {
	b.filter.ContentTypes = append(b.filter.ContentTypes, types...)
	return b
}
func (b *FilterBuilder) MetadataEquals(key, value string) *FilterBuilder {
	if b.filter.MetadataEquals == nil {
		b.filter.MetadataEquals = map[string]string{}
	}
	if existing, ok := b.filter.MetadataEquals[key]; ok && existing != value {
		b.experiment.buildIssues.add("RAG_CONFLICTING_FILTER", b.path+".metadataEquals."+key, "metadata key already has a different value")
		return b
	}
	b.filter.MetadataEquals[key] = value
	return b
}

type MetricsBuilder struct {
	experiment *ExperimentBuilder
	plan       *MetricsPlan
}

func (b *MetricsBuilder) RelevanceAt(grade RelevanceGrade) *MetricsBuilder {
	b.plan.RelevanceAt = &grade
	return b
}
func (b *MetricsBuilder) PrecisionAt(cutoffs ...int) *MetricsBuilder {
	b.plan.PrecisionAt = append(b.plan.PrecisionAt, cutoffs...)
	return b
}
func (b *MetricsBuilder) RecallAt(cutoffs ...int) *MetricsBuilder {
	b.plan.RecallAt = append(b.plan.RecallAt, cutoffs...)
	return b
}
func (b *MetricsBuilder) HitRateAt(cutoffs ...int) *MetricsBuilder {
	b.plan.HitRateAt = append(b.plan.HitRateAt, cutoffs...)
	return b
}
func (b *MetricsBuilder) NDCGAt(cutoff int) *MetricsBuilder { b.plan.NDCGAt = cutoff; return b }
func (b *MetricsBuilder) MRR() *MetricsBuilder              { b.plan.MRR = true; return b }
func (b *MetricsBuilder) MeanRelevantRecallAt(cutoff int) *MetricsBuilder {
	b.plan.MeanRelevantRecall = cutoff
	return b
}
func (b *MetricsBuilder) Abstention() *MetricsBuilder { b.plan.Abstention = true; return b }

func validateRequiredRef(report *ValidationReport, path string, kind ArtifactKind, ref ArtifactRef) {
	if ref.Kind != kind || ref.ID == "" {
		report.add("RAG_MISSING_INPUT", path, "required "+string(kind)+" artifact is missing")
	}
}
func validateFilter(report *ValidationReport, path string, filter FilterSpec) {
	for _, item := range append(append(append([]string{}, filter.SourceIDs...), filter.DocumentIDs...), filter.ContentTypes...) {
		if item == "" {
			report.add("RAG_INVALID_FILTER", path, "filter values must not be empty")
		}
	}
	for key, value := range filter.MetadataEquals {
		if key == "" || value == "" {
			report.add("RAG_INVALID_FILTER", path+".metadataEquals", "metadata key and value must not be empty")
		}
	}
}
func validateMetrics(report *ValidationReport, metrics MetricsPlan, results int) {
	if !metrics.RequiresRelevance() && !metrics.Abstention {
		report.add("RAG_MISSING_METRICS", "$.metrics", "select at least one metric")
	}
	if metrics.RequiresRelevance() && metrics.RelevanceAt == nil {
		report.add("RAG_MISSING_RELEVANCE_THRESHOLD", "$.metrics.relevanceAt", "graded metrics require a relevance threshold")
	}
	for _, cutoff := range append(append(append([]int{}, metrics.PrecisionAt...), metrics.RecallAt...), metrics.HitRateAt...) {
		if cutoff <= 0 || (results > 0 && cutoff > results) {
			report.add("RAG_INVALID_CUTOFF", "$.metrics", "metric cutoff must be positive and no greater than final results")
		}
	}
	if metrics.NDCGAt < 0 || (metrics.NDCGAt > 0 && results > 0 && metrics.NDCGAt > results) {
		report.add("RAG_INVALID_CUTOFF", "$.metrics.ndcgAt", "nDCG cutoff must be positive and no greater than final results")
	}
	if metrics.MeanRelevantRecall < 0 || (metrics.MeanRelevantRecall > 0 && results > 0 && metrics.MeanRelevantRecall > results) {
		report.add("RAG_INVALID_CUTOFF", "$.metrics.meanRelevantRecallAt", "mean relevant recall cutoff must be positive and no greater than final results")
	}
}
func normalizeInputs(inputs InputSpec) InputSpec {
	inputs.Representations = append([]RepresentationSpec(nil), inputs.Representations...)
	return inputs
}
func normalizeProvenance(p Provenance) Provenance {
	p.Fragments = uniqueStringsPreservingOrder(p.Fragments)
	p.Notes = append([]string(nil), p.Notes...)
	if len(p.Tags) == 0 {
		p.Tags = nil
	}
	return p
}
func normalizeRetrieval(plan RetrievalPlan) RetrievalPlan {
	plan.Filter = normalizeFilter(plan.Filter)
	plan.Channels = append([]ChannelSpec(nil), plan.Channels...)
	for i := range plan.Channels {
		plan.Channels[i].Filter = normalizeFilter(plan.Channels[i].Filter)
	}
	if plan.Fusion != nil {
		cloned := *plan.Fusion
		if len(cloned.Weights) == 0 {
			cloned.Weights = nil
		}
		plan.Fusion = &cloned
	}
	if plan.Reranking != nil {
		cloned := *plan.Reranking
		plan.Reranking = &cloned
	}
	return plan
}
func normalizeMetrics(m MetricsPlan) MetricsPlan {
	m.PrecisionAt = uniqueInts(m.PrecisionAt)
	m.RecallAt = uniqueInts(m.RecallAt)
	m.HitRateAt = uniqueInts(m.HitRateAt)
	return m
}
func normalizeFilter(filter FilterSpec) FilterSpec {
	filter.SourceIDs = uniqueStrings(filter.SourceIDs)
	filter.DocumentIDs = uniqueStrings(filter.DocumentIDs)
	filter.ContentTypes = uniqueStrings(filter.ContentTypes)
	if len(filter.MetadataEquals) == 0 {
		filter.MetadataEquals = nil
	}
	return filter
}
func appendRepresentation(existing []RepresentationSpec, incoming RepresentationSpec) []RepresentationSpec {
	for _, item := range existing {
		if item.Name == incoming.Name {
			return existing
		}
	}
	return append(existing, incoming)
}
func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	set := map[string]struct{}{}
	for _, v := range values {
		set[v] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for v := range set {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}
func uniqueStringsPreservingOrder(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
func uniqueInts(values []int) []int {
	if len(values) == 0 {
		return nil
	}
	set := map[int]struct{}{}
	for _, v := range values {
		set[v] = struct{}{}
	}
	result := make([]int, 0, len(set))
	for v := range set {
		result = append(result, v)
	}
	sort.Ints(result)
	return result
}
func itoa(value int) string {
	return strconv.Itoa(value)
}
