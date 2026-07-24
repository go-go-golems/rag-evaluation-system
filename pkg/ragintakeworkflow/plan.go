package ragintakeworkflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3product"
)

func CompilePlan(ctx context.Context, request Request, runtime RuntimeConfig) (workflowv3.WorkflowPlan, error) {
	normalized, err := NormalizeRequest(request, TargetDigest(runtime.DatabasePath))
	if err == nil && normalized.ProviderAuthorityDigest != runtime.ProviderAuthorityDigest {
		err = fmt.Errorf("RAG_INTAKE_PROVIDER_AUTHORITY")
	}
	if err != nil {
		return workflowv3.WorkflowPlan{}, err
	}
	environment, err := workflowv3product.NewAuthoringEnvironment([]string{PackageName}, NewPackage(runtime))
	if err != nil {
		return workflowv3.WorkflowPlan{}, err
	}
	authored, err := environment.Author(ctx, planSource(normalized, runtime))
	if err != nil {
		return workflowv3.WorkflowPlan{}, fmt.Errorf("RAG_INTAKE_PLAN: %w", err)
	}
	return authored.Plan, nil
}

func planSource(request Request, runtime RuntimeConfig) string {
	var body strings.Builder
	authority := runtime.ProviderAuthorityDigest
	if authority == "" {
		authority = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	}
	body.WriteString(`const workflow=require("workflow");const tasks=require("rag-intake-tasks");
const definition=workflow.define("rag-intake-v1",p=>{`)
	fmt.Fprintf(&body, `p.budget("provider",{limits:{requests:%d},policyDigest:%q});`, maxProviderOperations(runtime), authority)
	body.WriteString(`const request=p.input("request",{schema:"rag-intake-request/v1"});const all=[];const chunks=[];
`)
	if !request.SkipPreprocessing {
		for index := range request.DocumentIDs {
			fmt.Fprintf(&body, "const preprocess%03d=p.task(\"preprocess-%03d\",tasks.preprocessDocument({request}));all.push(preprocess%03d);\n", index, index, index)
		}
	}
	for index := range request.DocumentIDs {
		fmt.Fprintf(&body, "const chunk%03d=p.task(\"chunk-%03d\",tasks.chunkDocument({request}));chunks.push(chunk%03d);all.push(chunk%03d);\n", index, index, index, index)
	}
	dependencies := func(variable string, values []string) string {
		if len(values) == 0 {
			return ""
		}
		var chain strings.Builder
		fmt.Fprintf(&chain, ",job=>job")
		for _, value := range values {
			fmt.Fprintf(&chain, ".after(%s)", value)
		}
		return variable + chain.String()
	}
	chunkVars := make([]string, len(request.DocumentIDs))
	for i := range chunkVars {
		chunkVars[i] = fmt.Sprintf("chunk%03d", i)
	}
	if !request.SkipChunkEnrichment {
		for index := range request.ChunkIDs {
			fmt.Fprintf(&body, "const enrich%03d=p.task(\"enrich-%03d\",tasks.enrichChunk({request})%s);all.push(enrich%03d);\n", index, index, dependencies("", chunkVars), index)
		}
	}
	if !request.SkipEmbeddings {
		embedPolicy := dependencies("", chunkVars)
		if runtime.ProviderAuthorityDigest != "" {
			embedPolicy += `.budget({account:"provider",reserve:{requests:1},onExhausted:"fail-run"})`
		}
		fmt.Fprintf(&body, "const embed=p.task(\"embed\",tasks.computeEmbeddings({request})%s);all.push(embed);\n", embedPolicy)
	}
	if !request.SkipBM25 {
		fmt.Fprintf(&body, "const bm25=p.task(\"bm25\",tasks.buildBM25({request})%s);all.push(bm25);\n", dependencies("", chunkVars))
	}
	allVars := []string{}
	if !request.SkipPreprocessing {
		for i := range request.DocumentIDs {
			allVars = append(allVars, fmt.Sprintf("preprocess%03d", i))
		}
	}
	allVars = append(allVars, chunkVars...)
	if !request.SkipChunkEnrichment {
		for i := range request.ChunkIDs {
			allVars = append(allVars, fmt.Sprintf("enrich%03d", i))
		}
	}
	if !request.SkipEmbeddings {
		allVars = append(allVars, "embed")
	}
	if !request.SkipBM25 {
		allVars = append(allVars, "bm25")
	}
	fmt.Fprintf(&body, "const publish=p.task(\"publish\",tasks.publish({request})%s);p.output(\"result\",publish.output(\"result\"));});module.exports=workflow.compile(definition);\n", dependencies("", allVars))
	return body.String()
}
