package researchctladapter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-go-golems/researchctl/pkg/lab"
	"github.com/go-go-golems/researchctl/pkg/labprogress"
)

func TestPreparationProgressEmitterUsesMonotonicProducerSequences(t *testing.T) {
	var events []lab.EventInput
	emitter, err := NewPreparationProgressEmitter(func(event lab.EventInput) error {
		events = append(events, event)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	envelope := labprogress.Envelope{SchemaVersion: labprogress.SchemaVersionV1, Type: "rag.preparation.workflow.progress/v1", OccurredAt: time.Now().UTC(), WorkflowID: "workflow-1", WorkflowIdentityDigest: "sha256:identity", Counts: labprogress.Counts{Succeeded: 1, Pending: 1, Total: 2}}
	if err := emitter.Emit(envelope); err != nil {
		t.Fatal(err)
	}
	envelope.Counts.Succeeded = 2
	envelope.Counts.Pending = 0
	if err := emitter.Emit(envelope); err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].ProducerSequence == nil || events[1].ProducerSequence == nil || *events[0].ProducerSequence != 1 || *events[1].ProducerSequence != 2 {
		t.Fatalf("events=%#v", events)
	}
}

func TestPreparationProgressEventIsRedactedTypedLabEvent(t *testing.T) {
	envelope := labprogress.Envelope{SchemaVersion: labprogress.SchemaVersionV1, Type: "rag.preparation.workflow.progress/v1", OccurredAt: time.Now().UTC(), WorkflowID: "workflow-1", WorkflowIdentityDigest: "sha256:identity", Counts: labprogress.Counts{Succeeded: 2, Total: 2}, ProviderCalls: 1}
	event, err := PreparationProgressEvent(envelope, 7)
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != PreparationProgressEventType || event.ProducerSequence == nil || *event.ProducerSequence != 7 {
		t.Fatalf("event=%#v", event)
	}
	var got labprogress.Envelope
	if err := json.Unmarshal(event.Payload, &got); err != nil {
		t.Fatal(err)
	}
	if got.WorkflowID != envelope.WorkflowID || got.Counts != envelope.Counts {
		t.Fatalf("payload=%#v", got)
	}
}
