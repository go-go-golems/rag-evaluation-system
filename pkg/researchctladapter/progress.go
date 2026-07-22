package researchctladapter

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-go-golems/researchctl/pkg/lab"
	"github.com/go-go-golems/researchctl/pkg/labprogress"
)

const PreparationProgressEventType = "rag.preparation.progress/v1"

// PreparationProgressEmitter serializes redacted progress envelopes into the
// existing researchctl runner event shape. Its sink is normally the worker's
// stdio frame writer; researchctl commits the resulting EventInput to lab_events.
type PreparationProgressEmitter struct {
	mu       sync.Mutex
	sequence int64
	emit     func(lab.EventInput) error
}

func NewPreparationProgressEmitter(emit func(lab.EventInput) error) (*PreparationProgressEmitter, error) {
	if emit == nil {
		return nil, fmt.Errorf("preparation progress event emitter is required")
	}
	return &PreparationProgressEmitter{emit: emit}, nil
}

func (e *PreparationProgressEmitter) Emit(envelope labprogress.Envelope) error {
	if e == nil {
		return fmt.Errorf("preparation progress event emitter is nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sequence++
	event, err := PreparationProgressEvent(envelope, e.sequence)
	if err != nil {
		e.sequence--
		return err
	}
	return e.emit(event)
}

// PreparationProgressEvent converts an already-redacted typed preparation
// progress fact into the existing researchctl runner event contract. The
// caller still commits it through ObservationSink; this function performs no
// persistence and never carries provider input/output text or credentials.
func PreparationProgressEvent(envelope labprogress.Envelope, producerSequence int64) (lab.EventInput, error) {
	if producerSequence < 1 {
		return lab.EventInput{}, fmt.Errorf("preparation progress producer sequence must be positive")
	}
	if err := envelope.Validate(); err != nil {
		return lab.EventInput{}, err
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return lab.EventInput{}, fmt.Errorf("marshal preparation progress: %w", err)
	}
	return lab.EventInput{Type: PreparationProgressEventType, ProducerSequence: &producerSequence, ProducerOccurredAt: envelope.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"), Payload: payload}, nil
}
