package ragworkflow

import (
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

const (
	ProviderPackageName     = "rag-v2-geppetto"
	ProviderPackageVersion  = "1.0.0"
	ProviderAuthoritySchema = "rag-workflow-provider-authority/v1"
)

type WorkflowProviderIdentity struct {
	Role                               string `json:"role"`
	ProfileSlug                        string `json:"profileSlug,omitempty"`
	ModelManifestDigest                string `json:"modelManifestDigest"`
	ModelID                            string `json:"modelId"`
	SettingsFingerprint                string `json:"settingsFingerprint"`
	ConcurrencyLimit                   int    `json:"concurrencyLimit"`
	MaxResponseTokens                  int    `json:"maxResponseTokens,omitempty"`
	PricingConfigured                  bool   `json:"pricingConfigured"`
	InputCostMicrounitsPerMillion      int64  `json:"inputCostMicrounitsPerMillion,omitempty"`
	OutputCostMicrounitsPerMillion     int64  `json:"outputCostMicrounitsPerMillion,omitempty"`
	CacheReadCostMicrounitsPerMillion  int64  `json:"cacheReadCostMicrounitsPerMillion,omitempty"`
	CacheWriteCostMicrounitsPerMillion int64  `json:"cacheWriteCostMicrounitsPerMillion,omitempty"`
}

type ProviderAuthority struct {
	SchemaVersion         string                     `json:"schemaVersion"`
	ProfileID             string                     `json:"profileId"`
	Capabilities          []string                   `json:"capabilities"`
	ModelManifestDigests  []string                   `json:"modelManifestDigests"`
	PromptManifestDigests []string                   `json:"promptManifestDigests"`
	Providers             []WorkflowProviderIdentity `json:"providers"`
	Digest                string                     `json:"digest"`
}

func ProviderAuthorityFromSet(set *ragproviders.ProviderSet) (ProviderAuthority, error) {
	if set == nil {
		return ProviderAuthority{}, fmt.Errorf("RAG_WORKFLOW_PROVIDER_SET")
	}
	capabilities := set.CapabilityDescriptor()
	authority := ProviderAuthority{SchemaVersion: ProviderAuthoritySchema, ProfileID: capabilities.ProfileID, Capabilities: append([]string(nil), capabilities.Capabilities...), ModelManifestDigests: append([]string(nil), capabilities.ModelManifestDigests...), PromptManifestDigests: append([]string(nil), capabilities.PromptManifestDigests...)}
	for _, identity := range capabilities.EffectiveProviderIdentities {
		authority.Providers = append(authority.Providers, WorkflowProviderIdentity{Role: identity.Role, ProfileSlug: identity.ProfileSlug, ModelManifestDigest: identity.ModelManifestDigest, ModelID: identity.ModelID, SettingsFingerprint: identity.SettingsFingerprint, ConcurrencyLimit: identity.ConcurrencyLimit, MaxResponseTokens: identity.MaxResponseTokens, PricingConfigured: identity.PricingConfigured, InputCostMicrounitsPerMillion: identity.InputCostMicrounitsPerMillion, OutputCostMicrounitsPerMillion: identity.OutputCostMicrounitsPerMillion, CacheReadCostMicrounitsPerMillion: identity.CacheReadCostMicrounitsPerMillion, CacheWriteCostMicrounitsPerMillion: identity.CacheWriteCostMicrounitsPerMillion})
	}
	canonicalizeAuthority(&authority)
	digest, err := providerAuthorityDigest(authority)
	if err != nil {
		return ProviderAuthority{}, err
	}
	authority.Digest = digest
	if err := ValidateProviderAuthority(authority); err != nil {
		return ProviderAuthority{}, err
	}
	return authority, nil
}

func ValidateProviderAuthority(authority ProviderAuthority) error {
	if authority.SchemaVersion != ProviderAuthoritySchema || authority.ProfileID == "" || len(authority.ProfileID) > 128 || len(authority.Capabilities) == 0 || len(authority.Capabilities) > 32 || len(authority.ModelManifestDigests) == 0 || len(authority.ModelManifestDigests) > 256 || len(authority.PromptManifestDigests) > 256 || len(authority.Providers) == 0 || len(authority.Providers) > 64 {
		return fmt.Errorf("RAG_WORKFLOW_PROVIDER_AUTHORITY")
	}
	copyValue := authority
	canonicalizeAuthority(&copyValue)
	copyValue.Digest = ""
	original := authority
	original.Digest = ""
	left, _ := ragcontract.CanonicalJSON(copyValue)
	right, _ := ragcontract.CanonicalJSON(original)
	if string(left) != string(right) {
		return fmt.Errorf("RAG_WORKFLOW_PROVIDER_AUTHORITY_ORDER")
	}
	for _, values := range [][]string{authority.Capabilities, authority.ModelManifestDigests, authority.PromptManifestDigests} {
		previous := ""
		for _, value := range values {
			if value == "" || value <= previous {
				return fmt.Errorf("RAG_WORKFLOW_PROVIDER_AUTHORITY_ORDER")
			}
			previous = value
		}
	}
	for _, digest := range append(append([]string(nil), authority.ModelManifestDigests...), authority.PromptManifestDigests...) {
		if !providerDigestValid(digest) {
			return fmt.Errorf("RAG_WORKFLOW_PROVIDER_MANIFEST_DIGEST")
		}
	}
	want, err := providerAuthorityDigest(authority)
	if err != nil || authority.Digest != want {
		return fmt.Errorf("RAG_WORKFLOW_PROVIDER_AUTHORITY_DIGEST")
	}
	previousRole := ""
	for _, provider := range authority.Providers {
		if provider.Role == "" || provider.Role <= previousRole || provider.ModelID == "" || provider.ConcurrencyLimit < 1 || provider.ConcurrencyLimit > 100_000 || !providerDigestValid(provider.ModelManifestDigest) || !providerDigestValid(provider.SettingsFingerprint) {
			return fmt.Errorf("RAG_WORKFLOW_PROVIDER_IDENTITY")
		}
		if provider.PricingConfigured && (provider.InputCostMicrounitsPerMillion < 0 || provider.OutputCostMicrounitsPerMillion < 0 || provider.CacheReadCostMicrounitsPerMillion < 0 || provider.CacheWriteCostMicrounitsPerMillion < 0) {
			return fmt.Errorf("RAG_WORKFLOW_PROVIDER_IDENTITY")
		}
		previousRole = provider.Role
		if provider.MaxResponseTokens < 0 {
			return fmt.Errorf("RAG_WORKFLOW_PROVIDER_IDENTITY")
		}
	}
	return nil
}
func providerDigestValid(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil
}

func providerAuthorityDigest(authority ProviderAuthority) (string, error) {
	authority.Digest = ""
	return ragcontract.Digest(authority)
}
func canonicalizeAuthority(authority *ProviderAuthority) {
	sort.Strings(authority.Capabilities)
	sort.Strings(authority.ModelManifestDigests)
	sort.Strings(authority.PromptManifestDigests)
	sort.Slice(authority.Providers, func(i, j int) bool { return authority.Providers[i].Role < authority.Providers[j].Role })
}

type ProviderServices struct {
	Authority             ProviderAuthority
	Manifests             ragoperators.ManifestResolver
	Schemas               ragoperators.OutputSchemaValidator
	Generator             ragoperators.TextGenerator
	Embedder              ragoperators.Embedder
	Reranker              ragoperators.Reranker
	Cache                 ragoperators.Cache
	GenerationConcurrency int
}

func ProviderServicesFromSet(set *ragproviders.ProviderSet) (ProviderServices, error) {
	authority, err := ProviderAuthorityFromSet(set)
	if err != nil {
		return ProviderServices{}, err
	}
	return ProviderServices{Authority: authority, Manifests: set.Manifests, Schemas: set.Schemas, Generator: set.Generator, Embedder: set.Embedder, Reranker: set.Reranker, Cache: set.Cache, GenerationConcurrency: set.GenerationConcurrency}, nil
}

type ProviderPackage struct {
	services    ProviderServices
	authority   ProviderAuthority
	policy      ragworkflowops.Policy
	descriptors []workflowv3.ExternalOperationDescriptor
}

func NewProviderPackage(services ProviderServices, policy ragworkflowops.Policy) (*ProviderPackage, error) {
	authority := services.Authority
	if err := ValidateProviderAuthority(authority); err != nil {
		return nil, err
	}
	policy.AuthorityDigest = authority.Digest
	if policy.Reservations == nil {
		policy.Reservations = map[string][]workflowv3.ExternalOperationCounter{
			ragworkflowops.GenerateOperation: {{Name: "requests", Units: 1}},
			ragworkflowops.EmbedOperation:    {{Name: "requests", Units: 1}},
			ragworkflowops.RerankOperation:   {{Name: "requests", Units: 1}},
		}
	}
	if policy.MaxPerAttempt < 1 {
		policy.MaxPerAttempt = 10_000
	}
	if policy.FinishTimeout <= 0 {
		policy.FinishTimeout = 5 * time.Second
	}
	descriptors, err := ragworkflowops.NewDescriptors(authority.Digest, policy.MaxPerAttempt)
	if err != nil {
		return nil, err
	}
	if services.Generator == nil || services.Embedder == nil || services.Reranker == nil || services.Manifests == nil || services.Schemas == nil {
		return nil, fmt.Errorf("RAG_WORKFLOW_PROVIDER_CAPABILITIES")
	}
	return &ProviderPackage{services: services, authority: authority, policy: policy, descriptors: descriptors}, nil
}
func (p *ProviderPackage) Name() string              { return ProviderPackageName }
func (p *ProviderPackage) Version() string           { return ProviderPackageVersion }
func (p *ProviderPackage) RequiredModules() []string { return []string{ModuleAlias} }
func (p *ProviderPackage) DescriptorModules() []workflowmodule.DescriptorModule {
	return []workflowmodule.DescriptorModule{DescriptorModule()}
}
func (p *ProviderPackage) Bundle() (*workflowv3.Bundle, error) {
	authority, err := ragcontract.CanonicalJSON(p.authority)
	if err != nil {
		return nil, err
	}
	maximums, err := p.taskBudgetMaximums()
	if err != nil {
		return nil, err
	}
	return bundleFor(ProviderPackageName, ProviderPackageVersion, map[string][]byte{"task.cjs": taskSource, "provider-authority.json": authority}, maximums)
}
func (p *ProviderPackage) TaskModuleFactories() []workflowv3runtime.TaskModuleFactory {
	return []workflowv3runtime.TaskModuleFactory{newTaskModuleFactory(ModuleAlias, p.descriptors, p.authority.Digest, p.environment)}
}
func (p *ProviderPackage) Authority() ProviderAuthority {
	ret := p.authority
	ret.Capabilities = append([]string(nil), p.authority.Capabilities...)
	ret.ModelManifestDigests = append([]string(nil), p.authority.ModelManifestDigests...)
	ret.PromptManifestDigests = append([]string(nil), p.authority.PromptManifestDigests...)
	ret.Providers = append([]WorkflowProviderIdentity(nil), p.authority.Providers...)
	return ret
}
func (p *ProviderPackage) taskBudgetMaximums() (map[workflowv3.TaskKey]*workflowv3.BudgetClaim, error) {
	represent, err := p.operationBudgetClaim(map[string]int{ragworkflowops.GenerateOperation: p.policy.MaxPerAttempt})
	if err != nil {
		return nil, err
	}
	embed, err := p.operationBudgetClaim(map[string]int{ragworkflowops.EmbedOperation: p.policy.MaxPerAttempt})
	if err != nil {
		return nil, err
	}
	query, err := p.operationBudgetClaim(map[string]int{ragworkflowops.GenerateOperation: 1, ragworkflowops.EmbedOperation: 1, ragworkflowops.RerankOperation: 1})
	if err != nil {
		return nil, err
	}
	return map[workflowv3.TaskKey]*workflowv3.BudgetClaim{TaskRepresent: represent, TaskEmbed: embed, TaskQuery: query}, nil
}

func (p *ProviderPackage) operationBudgetClaim(operationCounts map[string]int) (*workflowv3.BudgetClaim, error) {
	amounts := map[string]int64{}
	for operation, count := range operationCounts {
		if count <= 0 {
			return nil, fmt.Errorf("RAG_WORKFLOW_PROVIDER_BUDGET_COUNT")
		}
		for _, reservation := range p.policy.Reservations[operation] {
			if reservation.Units > math.MaxInt64/int64(count) || amounts[reservation.Name] > math.MaxInt64-reservation.Units*int64(count) {
				return nil, fmt.Errorf("RAG_WORKFLOW_PROVIDER_BUDGET_OVERFLOW")
			}
			amounts[reservation.Name] += reservation.Units * int64(count)
		}
	}
	if len(amounts) == 0 {
		return nil, nil
	}
	dimensions := make([]string, 0, len(amounts))
	for dimension := range amounts {
		dimensions = append(dimensions, dimension)
	}
	sort.Strings(dimensions)
	reserve := make([]workflowv3.BudgetAmount, 0, len(dimensions))
	for _, dimension := range dimensions {
		reserve = append(reserve, workflowv3.BudgetAmount{Dimension: dimension, Units: amounts[dimension]})
	}
	return &workflowv3.BudgetClaim{Account: "provider", Reserve: reserve, OnExhausted: workflowv3.BudgetExhaustFailRun}, nil
}

func (p *ProviderPackage) applyBudgets(ir *workflowv3.WorkflowIR) error {
	if ir == nil {
		return fmt.Errorf("RAG_WORKFLOW_PROVIDER_BUDGET_IR")
	}
	maximums, err := p.taskBudgetMaximums()
	if err != nil {
		return err
	}
	totals := map[string]int64{}
	add := func(claim *workflowv3.BudgetClaim, multiplier int) error {
		if claim == nil {
			return nil
		}
		if multiplier <= 0 {
			return fmt.Errorf("RAG_WORKFLOW_PROVIDER_BUDGET_COUNT")
		}
		for _, amount := range claim.Reserve {
			if amount.Units > math.MaxInt64/int64(multiplier) || totals[amount.Dimension] > math.MaxInt64-amount.Units*int64(multiplier) {
				return fmt.Errorf("RAG_WORKFLOW_PROVIDER_BUDGET_OVERFLOW")
			}
			totals[amount.Dimension] += amount.Units * int64(multiplier)
		}
		return nil
	}
	for index := range ir.Nodes {
		claim := maximums[ir.Nodes[index].Task]
		ir.Nodes[index].Budget = cloneBudgetClaim(claim)
		if err := add(claim, 1); err != nil {
			return err
		}
	}
	for index := range ir.Maps {
		claim := maximums[ir.Maps[index].ItemTask]
		ir.Maps[index].Budget = cloneBudgetClaim(claim)
		if err := add(claim, ir.Maps[index].Policy.MaxItems); err != nil {
			return err
		}
	}
	if len(totals) == 0 {
		return nil
	}
	dimensions := make([]string, 0, len(totals))
	for dimension := range totals {
		dimensions = append(dimensions, dimension)
	}
	sort.Strings(dimensions)
	limits := make([]workflowv3.BudgetAmount, 0, len(dimensions))
	for _, dimension := range dimensions {
		limits = append(limits, workflowv3.BudgetAmount{Dimension: dimension, Units: totals[dimension]})
	}
	ir.Budgets = []workflowv3.BudgetAccount{{Account: "provider", Limits: limits, PolicyDigest: p.authority.Digest}}
	return nil
}

func (p *ProviderPackage) environment(moduleContext workflowv3runtime.TaskModuleContext, _ ragcontract.PipelineExecution) (*ragoperators.Environment, error) {
	decorator, err := ragworkflowops.NewDecorator(moduleContext.ExternalOperations, p.policy, nil)
	if err != nil {
		return nil, err
	}
	generator, err := decorator.Generator(p.services.Generator)
	if err != nil {
		return nil, err
	}
	embedder, err := decorator.Embedder(p.services.Embedder)
	if err != nil {
		return nil, err
	}
	reranker, err := decorator.Reranker(p.services.Reranker)
	if err != nil {
		return nil, err
	}
	return &ragoperators.Environment{Manifests: p.services.Manifests, Schemas: p.services.Schemas, Generator: generator, Embedder: embedder, Reranker: reranker, Cache: p.services.Cache, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: p.services.GenerationConcurrency, GenerationSettingsFingerprint: p.authority.Digest}, nil
}
