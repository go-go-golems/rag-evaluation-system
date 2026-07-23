package ragworkflowops

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/stretchr/testify/require"
)

type fakeRecorder struct {
	descriptors         map[string]workflowv3.ExternalOperationDescriptor
	specs               []workflowv3.ExternalOperationSpec
	completions         []workflowv3.ExternalOperationCompletion
	finishContextErrors []error
	beginErr, finishErr error
}

func (r *fakeRecorder) BeginExternalOperation(_ context.Context, spec workflowv3.ExternalOperationSpec) (workflowv3.ExternalOperationTicket, error) {
	if r.beginErr != nil {
		return workflowv3.ExternalOperationTicket{}, r.beginErr
	}
	r.specs = append(r.specs, spec)
	return workflowv3.ExternalOperationTicket{OperationID: "operation", CompletionKey: "secret-ticket"}, nil
}
func (r *fakeRecorder) FinishExternalOperation(ctx context.Context, _ workflowv3.ExternalOperationTicket, completion workflowv3.ExternalOperationCompletion) error {
	r.finishContextErrors = append(r.finishContextErrors, ctx.Err())
	r.completions = append(r.completions, completion)
	return r.finishErr
}

type fixedClock struct {
	values []time.Time
	index  int
}

func (c *fixedClock) Now() time.Time { value := c.values[c.index]; c.index++; return value }

type fakeGenerator struct {
	result ragoperators.GenerationResult
	err    error
	calls  int
}

func (f *fakeGenerator) Generate(ctx context.Context, _ ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	f.calls++
	if f.err != nil {
		return ragoperators.GenerationResult{}, f.err
	}
	if err := ctx.Err(); err != nil {
		return ragoperators.GenerationResult{}, err
	}
	return f.result, nil
}

type rejectingPreflightGenerator struct{ calls int }

func (p *rejectingPreflightGenerator) Generate(context.Context, ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	p.calls++
	return ragoperators.GenerationResult{}, nil
}
func (p *rejectingPreflightGenerator) PrepareGeneration(ragoperators.GenerationRequest) (PreparedGeneration, error) {
	return nil, errors.New("RAG_PREFLIGHT_REJECTED")
}

type fakeEmbedder struct {
	vectors [][]float64
	usage   ragoperators.Usage
	err     error
}

func (f fakeEmbedder) Embed(context.Context, string, []string) ([][]float64, ragoperators.Usage, error) {
	return f.vectors, f.usage, f.err
}

type fakeReranker struct {
	result ragoperators.RerankResult
	err    error
}

func (f fakeReranker) Rerank(context.Context, ragoperators.RerankRequest) (ragoperators.RerankResult, error) {
	return f.result, f.err
}

func newTestDecorator(t *testing.T, recorder *fakeRecorder, reservations map[string][]workflowv3.ExternalOperationCounter) *Decorator {
	t.Helper()
	clock := &fixedClock{values: []time.Time{time.Date(2026, 7, 24, 1, 2, 3, 0, time.UTC), time.Date(2026, 7, 24, 1, 2, 3, 2500000, time.UTC)}}
	decorator, err := NewDecorator(recorder, Policy{AuthorityDigest: "sha256:" + strings.Repeat("a", 64), MaxPerAttempt: 8, FinishTimeout: time.Second, Reservations: reservations}, clock)
	require.NoError(t, err)
	for _, descriptor := range decorator.Descriptors() {
		if recorder.descriptors == nil {
			recorder.descriptors = map[string]workflowv3.ExternalOperationDescriptor{}
		}
		recorder.descriptors[descriptor.Digest] = descriptor
	}
	return decorator
}

func TestGenerationOperationRecordsActualUsageAndSafeIdentity(t *testing.T) {
	recorder := &fakeRecorder{}
	decorator := newTestDecorator(t, recorder, nil)
	cost := 0.000123
	provider := &fakeGenerator{result: ragoperators.GenerationResult{Text: "SECRET-OUTPUT", InputTokens: 17, OutputTokens: 9, Cost: &cost}}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	result, err := wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "model-v1", Prompt: "prompt-v1", OutputSchema: "answer/v1", ParentID: "q1", Text: "SECRET-PROMPT", Count: 1})
	require.NoError(t, err)
	require.Equal(t, "SECRET-OUTPUT", result.Text)
	require.Equal(t, 1, provider.calls)
	require.Len(t, recorder.specs, 1)
	require.Len(t, recorder.completions, 1)
	spec, completion := recorder.specs[0], recorder.completions[0]
	descriptor := recorder.descriptors[spec.DescriptorDigest]
	require.Equal(t, GenerateOperation, descriptor.Kind.Name)
	require.NotContains(t, spec.CorrelationDigest, "SECRET")
	require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, completion.Outcome)
	require.Equal(t, workflowv3.ExternalOperationAccountingActual, completion.AccountingMode)
	require.Equal(t, int64(2500), completion.ElapsedMicros)
	require.Equal(t, []workflowv3.ExternalOperationCounter{{Name: "cost_microunits", Units: 123}, {Name: "input_tokens", Units: 17}, {Name: "output_tokens", Units: 9}, {Name: "requests", Units: 1}}, completion.Counters)
}

func TestCanceledCallFinishesWithDetachedContextAndConservativeReservation(t *testing.T) {
	recorder := &fakeRecorder{}
	reservations := map[string][]workflowv3.ExternalOperationCounter{GenerateOperation: {{Name: "requests", Units: 1}}}
	decorator := newTestDecorator(t, recorder, reservations)
	provider := &fakeGenerator{err: context.Canceled}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	_, err = wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "m", Prompt: "p", OutputSchema: "s"})
	require.ErrorIs(t, err, context.Canceled)
	require.Len(t, recorder.completions, 1)
	completion := recorder.completions[0]
	require.Equal(t, workflowv3.ExternalOperationOutcomeCanceled, completion.Outcome)
	require.Equal(t, workflowv3.ExternalOperationAccountingConservative, completion.AccountingMode)
	require.Equal(t, &workflowv3.ExternalOperationFailure{Class: "canceled", Code: "PROVIDER_CANCELED"}, completion.Failure)
	require.NoError(t, recorder.finishContextErrors[0])
}

func TestClassifiedProviderFailureRecordsSafeTaxonomyOnly(t *testing.T) {
	recorder := &fakeRecorder{}
	decorator := newTestDecorator(t, recorder, nil)
	raw := errors.New("SECRET-PROVIDER-BODY")
	provider := &fakeGenerator{err: NewProviderCallError(raw, "rate-limit", "PROVIDER_RATE_LIMITED", workflowv3.ExternalOperationOutcomeFailed)}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	_, err = wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "m", Prompt: "p", OutputSchema: "s"})
	require.EqualError(t, err, "PROVIDER_RATE_LIMITED")
	completion := recorder.completions[0]
	require.Equal(t, &workflowv3.ExternalOperationFailure{Class: "rate-limit", Code: "PROVIDER_RATE_LIMITED"}, completion.Failure)
	require.NotContains(t, completion.Failure.Code, "SECRET")
}

func TestAdmissionFailurePreventsProviderContact(t *testing.T) {
	recorder := &fakeRecorder{beginErr: errors.New("budget exhausted")}
	decorator := newTestDecorator(t, recorder, nil)
	provider := &fakeGenerator{}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	_, err = wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "m", Prompt: "p", OutputSchema: "s"})
	require.ErrorContains(t, err, "RAG_PROVIDER_OPERATION_ADMISSION")
	require.Zero(t, provider.calls)
	require.Empty(t, recorder.completions)
}

func TestPreflightFailureCreatesNoOperationAndNoProviderContact(t *testing.T) {
	recorder := &fakeRecorder{}
	decorator := newTestDecorator(t, recorder, nil)
	provider := &rejectingPreflightGenerator{}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	_, err = wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "m", Prompt: "p", OutputSchema: "s"})
	require.ErrorContains(t, err, "RAG_PREFLIGHT_REJECTED")
	require.Zero(t, provider.calls)
	require.Empty(t, recorder.specs)
}

func TestProviderResultErrorRejectsArbitraryText(t *testing.T) {
	err := ProviderSucceededWithInvalidResult("SECRET-PROVIDER-BODY")
	require.EqualError(t, err, "RAG_PROVIDER_RESULT_INVALID")
}

func TestProviderSuccessWithInvalidDomainObservationRemainsSucceeded(t *testing.T) {
	recorder := &fakeRecorder{}
	decorator := newTestDecorator(t, recorder, nil)
	invalid := math.NaN()
	provider := &fakeGenerator{result: ragoperators.GenerationResult{Text: "ok", Cost: &invalid}}
	wrapped, err := decorator.Generator(provider)
	require.NoError(t, err)
	_, err = wrapped.Generate(context.Background(), ragoperators.GenerationRequest{Kind: "generate.answer", Model: "m", Prompt: "p", OutputSchema: "s"})
	require.ErrorContains(t, err, "RAG_PROVIDER_COST_INVALID")
	require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, recorder.completions[0].Outcome)
	require.Equal(t, []workflowv3.ExternalOperationCounter{{Name: "requests", Units: 1}}, recorder.completions[0].Counters)
}

func TestRerankOperationRecordsProviderUsageAndCost(t *testing.T) {
	recorder := &fakeRecorder{}
	decorator := newTestDecorator(t, recorder, nil)
	cost := 0.000003
	wrapped, err := decorator.Reranker(fakeReranker{result: ragoperators.RerankResult{Scores: []ragoperators.RerankScore{{ChunkID: "a", Score: 1}}, InputTokens: 5, Cost: &cost}})
	require.NoError(t, err)
	_, err = wrapped.Rerank(context.Background(), ragoperators.RerankRequest{Model: "rerank-v1", Query: "SECRET-Q", Candidates: []ragoperators.Evidence{{}}})
	require.NoError(t, err)
	require.Equal(t, []workflowv3.ExternalOperationCounter{{Name: "cost_microunits", Units: 3}, {Name: "input_tokens", Units: 5}, {Name: "output_items", Units: 1}, {Name: "requests", Units: 1}}, recorder.completions[0].Counters)
}

func TestEmbeddingAndRerankDescriptorsPreserveCardinalityWithoutPayloads(t *testing.T) {
	for _, test := range []struct {
		name          string
		run           func(*Decorator) error
		kind          string
		input, output int64
	}{
		{"embed", func(d *Decorator) error {
			wrapped, err := d.Embedder(fakeEmbedder{vectors: [][]float64{{1, 0}, {0, 1}}, usage: ragoperators.Usage{EmbeddingTokens: 4}})
			if err != nil {
				return err
			}
			_, _, err = wrapped.Embed(context.Background(), "embed-v1", []string{"SECRET-A", "SECRET-B"})
			return err
		}, EmbedOperation, 2, 2},
		{"rerank", func(d *Decorator) error {
			wrapped, err := d.Reranker(fakeReranker{result: ragoperators.RerankResult{Scores: []ragoperators.RerankScore{{ChunkID: "a", Score: 1}}}})
			if err != nil {
				return err
			}
			_, err = wrapped.Rerank(context.Background(), ragoperators.RerankRequest{Model: "rerank-v1", Query: "SECRET-Q", Candidates: []ragoperators.Evidence{{}}})
			return err
		}, RerankOperation, 1, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := &fakeRecorder{}
			decorator := newTestDecorator(t, recorder, nil)
			require.NoError(t, test.run(decorator))
			descriptor := recorder.descriptors[recorder.specs[0].DescriptorDigest]
			require.Equal(t, test.kind, descriptor.Kind.Name)
			require.Equal(t, []workflowv3.ExternalOperationCounter{{Name: "input_items", Units: test.input}}, recorder.specs[0].Measures)
			require.Contains(t, recorder.completions[0].Counters, workflowv3.ExternalOperationCounter{Name: "output_items", Units: test.output})
		})
	}
}
