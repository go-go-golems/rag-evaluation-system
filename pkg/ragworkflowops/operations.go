// Package ragworkflowops decorates RAG provider interfaces with durable
// Scraper Workflow V3 external-operation custody.
package ragworkflowops

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const (
	GenerateOperation = "provider.generate"
	EmbedOperation    = "provider.embed"
	RerankOperation   = "provider.rerank"
)

type Clock interface{ Now() time.Time }
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now().UTC() }

// ProviderResultError reports that provider contact succeeded but domain result
// validation failed. The operation remains succeeded while the task fails.
type ProviderResultError struct{ Err error }

func (e *ProviderResultError) Error() string { return "RAG_PROVIDER_RESULT_INVALID" }
func (e *ProviderResultError) Unwrap() error { return e.Err }

func ProviderSucceededWithInvalidResult(err error) error {
	if err == nil {
		return nil
	}
	return &ProviderResultError{Err: err}
}

type Policy struct {
	AuthorityDigest string
	MaxPerAttempt   int
	FinishTimeout   time.Duration
	Reservations    map[string][]workflowv3.ExternalOperationCounter
}

type Decorator struct {
	recorder    workflowv3.ExternalOperationRecorder
	descriptors map[string]workflowv3.ExternalOperationDescriptor
	policy      Policy
	clock       Clock
}

func NewDecorator(recorder workflowv3.ExternalOperationRecorder, policy Policy, clock Clock) (*Decorator, error) {
	if recorder == nil || policy.MaxPerAttempt < 1 || policy.FinishTimeout <= 0 {
		return nil, fmt.Errorf("RAG_PROVIDER_OPERATION_CONFIG")
	}
	if clock == nil {
		clock = wallClock{}
	}
	descriptors, err := NewDescriptors(policy.AuthorityDigest, policy.MaxPerAttempt)
	if err != nil {
		return nil, err
	}
	byName := make(map[string]workflowv3.ExternalOperationDescriptor, len(descriptors))
	for _, descriptor := range descriptors {
		byName[descriptor.Kind.Name] = descriptor
	}
	for name, counters := range policy.Reservations {
		descriptor, ok := byName[name]
		if !ok {
			return nil, fmt.Errorf("RAG_PROVIDER_OPERATION_RESERVATION_KIND")
		}
		if err := workflowv3.ValidateExternalOperationSpec(descriptor, workflowv3.ExternalOperationSpec{DescriptorDigest: descriptor.Digest, Reservation: counters}); err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_OPERATION_RESERVATION: %w", err)
		}
	}
	return &Decorator{recorder: recorder, descriptors: byName, policy: policy, clock: clock}, nil
}

func NewDescriptors(authorityDigest string, maxPerAttempt int) ([]workflowv3.ExternalOperationDescriptor, error) {
	rolesMeasureUsage := []workflowv3.ExternalOperationCounterRole{workflowv3.ExternalOperationCounterMeasure, workflowv3.ExternalOperationCounterUsage}
	rolesReservationUsage := []workflowv3.ExternalOperationCounterRole{workflowv3.ExternalOperationCounterReservation, workflowv3.ExternalOperationCounterUsage}
	counter := func(name, unit string, roles []workflowv3.ExternalOperationCounterRole) workflowv3.ExternalOperationCounterDescriptor {
		return workflowv3.ExternalOperationCounterDescriptor{Name: name, Unit: unit, Roles: roles}
	}
	definitions := []struct {
		name     string
		counters []workflowv3.ExternalOperationCounterDescriptor
	}{
		{GenerateOperation, []workflowv3.ExternalOperationCounterDescriptor{counter("cost_microunits", "microunits", rolesReservationUsage), counter("input_items", "items", rolesMeasureUsage), counter("input_tokens", "tokens", rolesReservationUsage), counter("output_items", "items", rolesMeasureUsage), counter("output_tokens", "tokens", rolesReservationUsage), counter("requests", "requests", rolesReservationUsage)}},
		{EmbedOperation, []workflowv3.ExternalOperationCounterDescriptor{counter("cost_microunits", "microunits", rolesReservationUsage), counter("embedding_tokens", "tokens", rolesReservationUsage), counter("input_items", "items", rolesMeasureUsage), counter("output_items", "items", rolesMeasureUsage), counter("requests", "requests", rolesReservationUsage)}},
		{RerankOperation, []workflowv3.ExternalOperationCounterDescriptor{counter("cost_microunits", "microunits", rolesReservationUsage), counter("input_items", "items", rolesMeasureUsage), counter("input_tokens", "tokens", rolesReservationUsage), counter("output_items", "items", rolesMeasureUsage), counter("requests", "requests", rolesReservationUsage)}},
	}
	ret := make([]workflowv3.ExternalOperationDescriptor, 0, len(definitions))
	for _, definition := range definitions {
		descriptor, err := workflowv3.NewExternalOperationDescriptor(workflowv3.ExternalOperationDescriptor{Kind: workflowv3.ExternalOperationKind{Name: definition.name, Version: "v1"}, AuthorityDigest: authorityDigest, Counters: definition.counters, MaxPerAttempt: maxPerAttempt})
		if err != nil {
			return nil, err
		}
		ret = append(ret, descriptor)
	}
	return ret, nil
}

func (d *Decorator) Descriptors() []workflowv3.ExternalOperationDescriptor {
	ret := make([]workflowv3.ExternalOperationDescriptor, 0, len(d.descriptors))
	for _, value := range d.descriptors {
		ret = append(ret, value)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Kind.Name < ret[j].Kind.Name })
	return workflowv3.CloneExternalOperationDescriptors(ret)
}

func (d *Decorator) Generator(next ragoperators.TextGenerator) (ragoperators.TextGenerator, error) {
	if next == nil {
		return nil, fmt.Errorf("RAG_PROVIDER_GENERATOR_REQUIRED")
	}
	return generationDecorator{decorator: d, next: next}, nil
}
func (d *Decorator) Embedder(next ragoperators.Embedder) (ragoperators.Embedder, error) {
	if next == nil {
		return nil, fmt.Errorf("RAG_PROVIDER_EMBEDDER_REQUIRED")
	}
	return embeddingDecorator{decorator: d, next: next}, nil
}
func (d *Decorator) Reranker(next ragoperators.Reranker) (ragoperators.Reranker, error) {
	if next == nil {
		return nil, fmt.Errorf("RAG_PROVIDER_RERANKER_REQUIRED")
	}
	return rerankDecorator{decorator: d, next: next}, nil
}

type generationDecorator struct {
	decorator *Decorator
	next      ragoperators.TextGenerator
}

func (w generationDecorator) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	correlation, err := safeDigest(struct {
		SchemaVersion, Kind, Model, Prompt, OutputSchema, ParentID string
		Count                                                      int
	}{"rag-provider-generation-correlation/v1", request.Kind, request.Model, request.Prompt, request.OutputSchema, request.ParentID, request.Count})
	if err != nil {
		return ragoperators.GenerationResult{}, err
	}
	measures := counters(counter("input_items", 1))
	if request.Count > 0 {
		measures = append(measures, counter("output_items", int64(request.Count)))
		sortCounters(measures)
	}
	return executeOperation(ctx, w.decorator, GenerateOperation, correlation, measures, func(callCtx context.Context) (ragoperators.GenerationResult, []workflowv3.ExternalOperationCounter, error) {
		result, callErr := w.next.Generate(callCtx, request)
		if callErr != nil {
			return result, nil, callErr
		}
		observed := []workflowv3.ExternalOperationCounter{counter("requests", 1)}
		if result.InputTokens > 0 {
			observed = append(observed, counter("input_tokens", result.InputTokens))
		}
		if result.OutputTokens > 0 {
			observed = append(observed, counter("output_tokens", result.OutputTokens))
		}
		if result.Cost != nil {
			units, costErr := costMicrounits(*result.Cost)
			if costErr != nil {
				return result, observed, ProviderSucceededWithInvalidResult(costErr)
			}
			observed = append(observed, counter("cost_microunits", units))
		}
		sortCounters(observed)
		return result, observed, nil
	})
}

type embeddingDecorator struct {
	decorator *Decorator
	next      ragoperators.Embedder
}

func (w embeddingDecorator) Embed(ctx context.Context, model string, texts []string) ([][]float64, ragoperators.Usage, error) {
	correlation, err := safeDigest(struct {
		SchemaVersion, Model string
		InputItems           int
	}{"rag-provider-embedding-correlation/v1", model, len(texts)})
	if err != nil {
		return nil, ragoperators.Usage{}, err
	}
	type value struct {
		vectors [][]float64
		usage   ragoperators.Usage
	}
	result, err := executeOperation(ctx, w.decorator, EmbedOperation, correlation, counters(counter("input_items", int64(len(texts)))), func(callCtx context.Context) (value, []workflowv3.ExternalOperationCounter, error) {
		vectors, usage, callErr := w.next.Embed(callCtx, model, texts)
		if callErr != nil {
			return value{}, nil, callErr
		}
		observed := []workflowv3.ExternalOperationCounter{counter("requests", 1), counter("output_items", int64(len(vectors)))}
		if usage.EmbeddingTokens > 0 {
			observed = append(observed, counter("embedding_tokens", usage.EmbeddingTokens))
		}
		sortCounters(observed)
		return value{vectors, usage}, observed, nil
	})
	return result.vectors, result.usage, err
}

type rerankDecorator struct {
	decorator *Decorator
	next      ragoperators.Reranker
}

func (w rerankDecorator) Rerank(ctx context.Context, request ragoperators.RerankRequest) ([]ragoperators.RerankScore, error) {
	correlation, err := safeDigest(struct {
		SchemaVersion, Model string
		InputItems           int
	}{"rag-provider-rerank-correlation/v1", request.Model, len(request.Candidates)})
	if err != nil {
		return nil, err
	}
	return executeOperation(ctx, w.decorator, RerankOperation, correlation, counters(counter("input_items", int64(len(request.Candidates)))), func(callCtx context.Context) ([]ragoperators.RerankScore, []workflowv3.ExternalOperationCounter, error) {
		scores, callErr := w.next.Rerank(callCtx, request)
		if callErr != nil {
			return nil, nil, callErr
		}
		observed := []workflowv3.ExternalOperationCounter{counter("output_items", int64(len(scores))), counter("requests", 1)}
		sortCounters(observed)
		return scores, observed, nil
	})
}

func executeOperation[T any](ctx context.Context, d *Decorator, name, correlation string, measures []workflowv3.ExternalOperationCounter, call func(context.Context) (T, []workflowv3.ExternalOperationCounter, error)) (T, error) {
	var zero T
	descriptor, ok := d.descriptors[name]
	if !ok {
		return zero, fmt.Errorf("RAG_PROVIDER_OPERATION_DESCRIPTOR")
	}
	ticket, err := d.recorder.BeginExternalOperation(ctx, workflowv3.ExternalOperationSpec{DescriptorDigest: descriptor.Digest, CorrelationDigest: correlation, Reservation: append([]workflowv3.ExternalOperationCounter(nil), d.policy.Reservations[name]...), Measures: measures})
	if err != nil {
		return zero, fmt.Errorf("RAG_PROVIDER_OPERATION_ADMISSION")
	}
	started := d.clock.Now()
	value, observed, callErr := call(ctx)
	finished := d.clock.Now()
	completion := completionFor(callErr, started, finished, observed, len(d.policy.Reservations[name]) > 0)
	finishCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), d.policy.FinishTimeout)
	defer cancel()
	if finishErr := d.recorder.FinishExternalOperation(finishCtx, ticket, completion); finishErr != nil {
		return zero, fmt.Errorf("RAG_PROVIDER_OPERATION_FINISH")
	}
	if callErr != nil {
		return zero, callErr
	}
	return value, nil
}

func completionFor(err error, started, finished time.Time, counters []workflowv3.ExternalOperationCounter, reserved bool) workflowv3.ExternalOperationCompletion {
	elapsed := finished.Sub(started).Microseconds()
	if elapsed < 0 {
		elapsed = 0
	}
	completion := workflowv3.ExternalOperationCompletion{ProviderStartedAt: started.UTC(), ElapsedMicros: elapsed, Outcome: workflowv3.ExternalOperationOutcomeSucceeded, AccountingMode: workflowv3.ExternalOperationAccountingActual, Counters: counters}
	if err == nil {
		return completion
	}
	var resultErr *ProviderResultError
	if errors.As(err, &resultErr) {
		return completion
	}
	completion.Counters = nil
	if reserved {
		completion.AccountingMode = workflowv3.ExternalOperationAccountingConservative
	} else {
		completion.AccountingMode = workflowv3.ExternalOperationAccountingNone
	}
	class, code, outcome := "transport", "PROVIDER_TRANSPORT", workflowv3.ExternalOperationOutcomeFailed
	if errors.Is(err, context.Canceled) {
		class, code, outcome = "canceled", "PROVIDER_CANCELED", workflowv3.ExternalOperationOutcomeCanceled
	} else if errors.Is(err, context.DeadlineExceeded) {
		class, code, outcome = "timeout", "PROVIDER_TIMEOUT", workflowv3.ExternalOperationOutcomeTimedOut
	}
	completion.Outcome = outcome
	completion.Failure = &workflowv3.ExternalOperationFailure{Class: class, Code: code}
	return completion
}
func safeDigest(value any) (string, error) { return ragcontract.Digest(value) }
func counter(name string, units int64) workflowv3.ExternalOperationCounter {
	return workflowv3.ExternalOperationCounter{Name: name, Units: units}
}
func counters(values ...workflowv3.ExternalOperationCounter) []workflowv3.ExternalOperationCounter {
	ret := values[:0]
	for _, v := range values {
		if v.Units > 0 {
			ret = append(ret, v)
		}
	}
	sortCounters(ret)
	return ret
}
func sortCounters(values []workflowv3.ExternalOperationCounter) {
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
}
func costMicrounits(value float64) (int64, error) {
	if value < 0 || math.IsNaN(value) || math.IsInf(value, 0) || value > float64(math.MaxInt64)/1_000_000 {
		return 0, fmt.Errorf("RAG_PROVIDER_COST_INVALID")
	}
	return int64(math.Round(value * 1_000_000)), nil
}
