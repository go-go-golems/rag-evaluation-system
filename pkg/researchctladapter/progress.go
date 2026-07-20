package researchctladapter

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/researchctl/pkg/lab"
	"github.com/go-go-golems/researchctl/pkg/labprogress"
)

const PreparationProgressEventType = "rag.preparation.progress/v1"

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
