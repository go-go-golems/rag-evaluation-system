package workflowv3ttc

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
)

func TestProductionWorkflowPinsPublicationGate(t *testing.T) {
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := registry.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	authored, err := workflowmodule.Author(context.Background(), ProductionWorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		t.Fatal(err)
	}
	if len(authored.Plan.Gates) != 4 {
		t.Fatalf("gates=%#v", authored.Plan.Gates)
	}
	publicationGateFound := false
	for _, gate := range authored.Plan.Gates {
		if gate.Key == "approve-publication" && gate.Policy.RequiredRole == "rag.ttc.publisher" {
			publicationGateFound = true
		}
	}
	if !publicationGateFound {
		t.Fatalf("publication gate missing: %#v", authored.Plan.Gates)
	}
	if len(authored.Plan.Nodes) != 2 || authored.Plan.Nodes[1].Key != "publish-prepared" || authored.Plan.Nodes[0].Bindings["shard"].Source != "reduction-output" {
		t.Fatalf("nodes=%#v", authored.Plan.Nodes)
	}
}

type fixtureEvaluation struct{}

func (*fixtureEvaluation) Evaluate(_ context.Context, _ PublicationReceipt, query QueryEnvelope) (QueryEvidence, error) {
	return QueryEvidence{SchemaVersion: QueryEvidenceSchema, QueryID: query.Query.ID, DatasetDigest: query.DatasetDigest, CitationChunkIDs: append([]string(nil), query.Query.RelevantIDs...), Usage: []Usage{{Dimension: "cost_microunits", Units: 0}, {Dimension: "embedding_tokens", Units: 0}, {Dimension: "input_tokens", Units: 0}, {Dimension: "output_tokens", Units: 0}, {Dimension: "requests", Units: 2}}}, nil
}

type fixturePublication struct {
	mu        sync.Mutex
	publishes int
}

func (p *fixturePublication) Validate(shard PreparedShard) (ValidationReceipt, error) {
	return ValidationReceipt{SchemaVersion: ValidationReceiptSchema, ShardDigest: shard.Digest, ItemCount: len(shard.Items), RepresentationCount: len(shard.Items), EmbeddingCount: len(shard.Items)}, nil
}
func (p *fixturePublication) Publish(_ context.Context, shard PreparedShard, decision PublicationDecision) (PublicationReceipt, error) {
	if !decision.Approved {
		return PublicationReceipt{}, fmt.Errorf("not approved")
	}
	p.mu.Lock()
	p.publishes++
	p.mu.Unlock()
	return PublicationReceipt{SchemaVersion: PublicationReceiptSchema, ShardDigest: shard.Digest, PreparedDigest: digestOf("f"), Identity: ragengine.PreparedCorpusIdentity{SchemaVersion: "rag-prepared-corpus-manifest/v1"}, ItemCount: len(shard.Items)}, nil
}

func TestProductionWorkflowWaitsWithoutLeaseThenPublishesAfterAuthorizedDecision(t *testing.T) {
	ctx := context.Background()
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	catalog, _ := registry.Catalog()
	authored, err := workflowmodule.Author(ctx, ProductionWorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	artifacts, err := workflowv3.NewFileArtifactStore(filepath.Join(root, "artifacts"), 1<<26)
	if err != nil {
		t.Fatal(err)
	}
	items := make([]workflowv3.ManifestItem, 2)
	for index := range items {
		key := fmt.Sprintf("chunk-%04d", index)
		text := "gate fixture source"
		digest, _ := workflowv3.Digest(text)
		chunk := Chunk{Key: key, Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: key, ParentUnitID: key, TextDigest: digest, Citation: ragcontract.CitationRef{SourceID: key}}, Text: text}, CitationIDs: []string{key}, SourceDigest: digestOf("d")}
		body, _ := workflowv3.CanonicalJSON(chunk)
		ref, putErr := artifacts.Put(ctx, ChunkSchema, "application/json", body)
		if putErr != nil {
			t.Fatal(putErr)
		}
		items[index] = workflowv3.ManifestItem{Key: key, Value: ref}
	}
	manifest, _ := workflowv3.NewItemManifest(ChunkSchema, items)
	body, _ := workflowv3.EncodeItemManifest(manifest)
	manifestRef, _ := artifacts.Put(ctx, workflowv3.ItemManifestSchemaV1, "application/json", body)
	query := QueryEnvelope{SchemaVersion: QuerySchema, DatasetDigest: digestOf("9"), Query: ragoperators.Query{ID: "query-1", Text: "What is bounded?", RelevantIDs: []string{"chunk-0000"}}}
	queryBody, _ := workflowv3.CanonicalJSON(query)
	queryRef, _ := artifacts.Put(ctx, QuerySchema, "application/json", queryBody)
	queryManifest, _ := workflowv3.NewItemManifest(QuerySchema, []workflowv3.ManifestItem{{Key: "query-1", Value: queryRef}})
	queryManifestBody, _ := workflowv3.EncodeItemManifest(queryManifest)
	queryManifestRef, _ := artifacts.Put(ctx, workflowv3.ItemManifestSchemaV1, "application/json", queryManifestBody)
	provider := &fixtureProvider{calls: map[string]int{}}
	publication := &fixturePublication{}
	modules, err := workflowv3runtime.NewTaskModuleRegistry(ModuleWithAuthorities(provider, publication, &fixtureEvaluation{}))
	if err != nil {
		t.Fatal(err)
	}
	databasePath := filepath.Join(root, "workflow.sqlite")
	store, err := workflowv3sqlite.Open(ctx, databasePath)
	if err != nil {
		t.Fatal(err)
	}
	engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules}
	if err := engine.Submit(ctx, "production-gate", authored.Plan, map[string]workflowv3.ArtifactRef{"chunks": manifestRef, "queries": queryManifestRef}); err != nil {
		t.Fatal(err)
	}
	runDispatcher := func() (context.CancelFunc, <-chan error) {
		dispatchCtx, cancel := context.WithCancel(ctx)
		done := make(chan error, 1)
		dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{ResourceGeneration: 2, ResourceEmbedding: 2, ResourceLocal: 1, ResourceEvaluation: 1}, PollInterval: time.Millisecond}
		go func() { done <- dispatcher.Run(dispatchCtx) }()
		return cancel, done
	}
	waitForGate := func(gateKey workflowv3.NodeKey) {
		t.Helper()
		cancel, done := runDispatcher()
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			select {
			case dispatchErr := <-done:
				t.Fatalf("dispatcher stopped waiting for %s: %v", gateKey, dispatchErr)
			default:
			}
			operational, snapshotErr := store.OperationalSnapshot(ctx, nil, registry, nil, time.Now().UTC())
			if snapshotErr == nil && len(operational.Queue.ActiveByResource) == 0 {
				for _, gate := range operational.Gates {
					if gate.GateKey == gateKey && gate.Status == "waiting" {
						cancel()
						<-done
						return
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
		cancel()
		<-done
		t.Fatalf("gate %s did not become ready", gateKey)
	}
	waitForGate("approve-generation-spend")
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = workflowv3sqlite.Open(ctx, databasePath)
	if err != nil {
		t.Fatal(err)
	}
	engine.Store = store
	publication.mu.Lock()
	before := publication.publishes
	publication.mu.Unlock()
	if before != 0 {
		t.Fatal("publication ran before approval")
	}
	budgetRef, err := artifacts.Put(ctx, "rag-ttc-budget-decision/v1", "application/json", []byte(`{"approved":true}`))
	if err != nil {
		t.Fatal(err)
	}
	increase := func(account string, dimensions []struct {
		name  string
		delta int64
	}) {
		for index, dimension := range dimensions {
			if err := store.IncreaseBudget(ctx, "production-gate", account, dimension.name, dimension.delta, int64(index+1), "budget-operator", time.Now().UTC()); err != nil {
				t.Fatal(err)
			}
		}
	}
	increase("generation", []struct {
		name  string
		delta int64
	}{{"cost_microunits", 20_000}, {"input_tokens", 2_048}, {"output_tokens", 2_048}, {"requests", 1}})
	if err := store.DecideGate(ctx, workflowv3.GateDecisionCommand{RunID: "production-gate", GateKey: "approve-generation-spend", ExpectedVersion: 1, Decision: "approve", DecisionCode: "TTC_BUDGET_APPROVED", ActorID: "operator", AuthorizedRole: "rag.ttc.budget-approver", DecisionRef: &budgetRef}, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	waitForGate("approve-embedding-spend")
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = workflowv3sqlite.Open(ctx, databasePath)
	if err != nil {
		t.Fatal(err)
	}
	engine.Store = store
	increase("embedding", []struct {
		name  string
		delta int64
	}{{"embedding_tokens", 4_096}})
	if err := store.DecideGate(ctx, workflowv3.GateDecisionCommand{RunID: "production-gate", GateKey: "approve-embedding-spend", ExpectedVersion: 1, Decision: "approve", DecisionCode: "TTC_BUDGET_APPROVED", ActorID: "operator", AuthorizedRole: "rag.ttc.budget-approver", DecisionRef: &budgetRef}, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	waitForGate("approve-publication")
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = workflowv3sqlite.Open(ctx, databasePath)
	if err != nil {
		t.Fatal(err)
	}
	engine.Store = store
	decisionBody, _ := workflowv3.CanonicalJSON(PublicationDecision{SchemaVersion: PublicationDecisionSchema, Approved: true, ShardDigest: digestOf("a"), PolicyDigest: digestOf("b")})
	decisionRef, err := artifacts.Put(ctx, PublicationDecisionSchema, "application/json", decisionBody)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.DecideGate(ctx, workflowv3.GateDecisionCommand{RunID: "production-gate", GateKey: "approve-publication", ExpectedVersion: 1, Decision: "approve", DecisionCode: "TTC_PUBLICATION_APPROVED", ActorID: "operator", AuthorizedRole: "rag.ttc.publisher", DecisionRef: &decisionRef}, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	cancel, done := runDispatcher()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		snapshot, snapshotErr := engine.Snapshot(ctx, "production-gate")
		if snapshotErr == nil && snapshot.Status == "failed" {
			cancel()
			<-done
			for _, attempt := range snapshot.Attempts {
				if attempt.Failure != nil {
					t.Fatalf("production attempt %s failed: %+v", attempt.NodeKey, *attempt.Failure)
				}
			}
			t.Fatal("production workflow failed without attempt evidence")
		}
		if snapshotErr == nil && snapshot.Status == "succeeded" {
			cancel()
			<-done
			publication.mu.Lock()
			count := publication.publishes
			publication.mu.Unlock()
			if count != 1 {
				t.Fatalf("publishes=%d", count)
			}
			_ = store.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	<-done
	snapshot, _ := engine.Snapshot(ctx, "production-gate")
	queue, _ := store.QueueSnapshot(ctx, registry, map[string]int{ResourceGeneration: 2, ResourceEmbedding: 2, ResourceLocal: 1, ResourceEvaluation: 1}, time.Now().UTC())
	t.Fatalf("production workflow did not publish: status=%s maps=%+v blocked=%+v", snapshot.Status, queue.Maps, queue.BlockedByReason)
}
