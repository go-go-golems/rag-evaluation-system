package geppetto

import (
	"context"
	"errors"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

func classifyProviderError(err error) error {
	switch {
	case errors.Is(err, context.Canceled):
		return ragworkflowops.NewProviderCallError(err, "canceled", "PROVIDER_CANCELED", workflowv3.ExternalOperationOutcomeCanceled)
	case errors.Is(err, context.DeadlineExceeded):
		return ragworkflowops.NewProviderCallError(err, "timeout", "PROVIDER_TIMEOUT", workflowv3.ExternalOperationOutcomeTimedOut)
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "status=429"), strings.Contains(message, "status code: 429"), strings.Contains(message, "too many requests"):
		return ragworkflowops.NewProviderCallError(err, "rate-limit", "PROVIDER_RATE_LIMITED", workflowv3.ExternalOperationOutcomeFailed)
	case strings.Contains(message, "status=5"), strings.Contains(message, "status code: 5"), strings.Contains(message, "bad gateway"), strings.Contains(message, "service unavailable"):
		return ragworkflowops.NewProviderCallError(err, "provider-5xx", "PROVIDER_SERVER_ERROR", workflowv3.ExternalOperationOutcomeFailed)
	default:
		return ragworkflowops.NewProviderCallError(err, "transport", "PROVIDER_TRANSPORT", workflowv3.ExternalOperationOutcomeFailed)
	}
}
