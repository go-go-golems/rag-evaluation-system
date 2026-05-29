package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/engine/runner"
)

// EchoRunner is a Phase 0 compatibility runner. It proves rag-eval can register
// a Go-native scraper runner without wiring production intake services yet.
type EchoRunner struct{}

type EchoInput struct {
	Operation string          `json:"operation"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type EchoOutput struct {
	Operation  string          `json:"operation"`
	WorkflowID string          `json:"workflow_id"`
	OpID       string          `json:"op_id"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

func (r *EchoRunner) Kind() string { return IntakeRunnerKind }

func (r *EchoRunner) Run(ctx context.Context, runCtx runner.RunContext) (*model.OpResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var input EchoInput
	if len(runCtx.Op.Input) > 0 {
		if err := json.Unmarshal(runCtx.Op.Input, &input); err != nil {
			return nil, fmt.Errorf("decode intake op input: %w", err)
		}
	}
	if input.Operation == "" {
		input.Operation = "echo"
	}
	if input.Operation != "echo" {
		data, _ := json.Marshal(map[string]any{"operation": input.Operation})
		return &model.OpResult{
			OpID: runCtx.Op.ID,
			Data: data,
			Error: &model.OpError{
				Code:      "unknown_operation",
				Message:   fmt.Sprintf("unknown intake operation %q", input.Operation),
				Retryable: false,
			},
		}, nil
	}

	output := EchoOutput{
		Operation:  input.Operation,
		WorkflowID: string(runCtx.Workflow.ID),
		OpID:       string(runCtx.Op.ID),
		Payload:    input.Payload,
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
}
