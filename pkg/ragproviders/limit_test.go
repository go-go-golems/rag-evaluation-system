package ragproviders

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type blockingGenerator struct {
	mu      sync.Mutex
	active  int
	maximum int
	reached chan struct{}
	gate    <-chan struct{}
}

func (g *blockingGenerator) Generate(ctx context.Context, _ ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	g.mu.Lock()
	g.active++
	if g.active > g.maximum {
		g.maximum = g.active
	}
	if g.active == 3 {
		select {
		case <-g.reached:
		default:
			close(g.reached)
		}
	}
	g.mu.Unlock()
	select {
	case <-g.gate:
	case <-ctx.Done():
		return ragoperators.GenerationResult{}, ctx.Err()
	}
	g.mu.Lock()
	g.active--
	g.mu.Unlock()
	return ragoperators.GenerationResult{}, nil
}

func (g *blockingGenerator) max() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.maximum
}

func TestLimitedGeneratorHonorsConfiguredMaximumInFlight(t *testing.T) {
	gate := make(chan struct{})
	inner := &blockingGenerator{reached: make(chan struct{}), gate: gate}
	limited, err := newLimitedGenerator(inner, 3)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := limited.Generate(context.Background(), ragoperators.GenerationRequest{}); err != nil {
				t.Errorf("Generate() error = %v", err)
			}
		}()
	}
	select {
	case <-inner.reached:
	case <-time.After(time.Second):
		t.Fatal("did not reach configured concurrency")
	}
	if got := inner.max(); got != 3 {
		t.Fatalf("maximum in flight = %d, want 3", got)
	}
	close(gate)
	wg.Wait()
}

type namedGenerator struct{ name string }

func (g namedGenerator) Generate(_ context.Context, _ ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	return ragoperators.GenerationResult{Text: g.name}, nil
}

func TestGeneratorRouterDispatchesOnlyKnownModelIDs(t *testing.T) {
	router := generatorRouter{byModel: map[string]ragoperators.TextGenerator{"flash": namedGenerator{name: "flash"}, "primary": namedGenerator{name: "primary"}}}
	value, err := router.Generate(context.Background(), ragoperators.GenerationRequest{Model: "flash"})
	if err != nil || value.Text != "flash" {
		t.Fatalf("flash route: %#v %v", value, err)
	}
	if _, err := router.Generate(context.Background(), ragoperators.GenerationRequest{Model: "unknown"}); err == nil {
		t.Fatal("unknown model route unexpectedly accepted")
	}
}

func TestProviderConcurrencyLimitDefaultsToOne(t *testing.T) {
	if got := providerConcurrencyLimit(ProviderSpec{}); got != 1 {
		t.Fatalf("default limit = %d, want 1", got)
	}
	if got := providerConcurrencyLimit(ProviderSpec{Concurrency: ConcurrencyConfig{MaxInFlight: 3}}); got != 3 {
		t.Fatalf("configured limit = %d, want 3", got)
	}
}
