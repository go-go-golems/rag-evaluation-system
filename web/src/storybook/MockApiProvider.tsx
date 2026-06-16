import { configureStore } from "@reduxjs/toolkit";
import type { ReactNode } from "react";
import { useEffect, useMemo } from "react";
import { Provider } from "react-redux";
import { ragApi } from "../services/api";
import type { WidgetNode } from "@go-go-golems/rag-evaluation-site/ir";
import { component, text } from "@go-go-golems/rag-evaluation-site/ir";

interface MockApiProviderProps {
	children: ReactNode;
}

const now = "2026-06-07T12:00:00Z";

const sources = [
	{
		id: "src-tree-guides",
		name: "Tree Center Guides",
		type: "web_crawl",
		created_at: now,
		updated_at: now,
	},
	{
		id: "src-product-pages",
		name: "Product Pages",
		type: "sitemap",
		created_at: now,
		updated_at: now,
	},
];

const documents = [
	{
		id: "doc-privacy-trees",
		source_id: "src-tree-guides",
		external_id: "privacy-trees",
		title: "Fast Growing Trees for Privacy Screens",
		author: "Tree Center Editorial",
		url: "https://example.test/privacy-trees",
		content_type: "text/html",
		word_count: 1824,
		language: "en",
		status: "indexed",
		created_at: now,
		updated_at: now,
	},
	{
		id: "doc-hydrangea-pruning",
		source_id: "src-tree-guides",
		external_id: "hydrangea-pruning",
		title: "Hydrangea Pruning by Bloom Type",
		author: "Tree Center Editorial",
		url: "https://example.test/hydrangea-pruning",
		content_type: "text/html",
		word_count: 1420,
		language: "en",
		status: "chunked",
		created_at: now,
		updated_at: now,
	},
];

const chunks = [
	{
		id: "chunk_privacy_001",
		document_id: "doc-privacy-trees",
		strategy_id: "fixed-1200-150",
		chunk_index: 1,
		text: "Choosing fast growing trees for privacy depends on mature width, climate, soil drainage, and maintenance expectations.",
		token_count: 220,
		start_offset: 0,
		end_offset: 1100,
	},
	{
		id: "chunk_privacy_002",
		document_id: "doc-privacy-trees",
		strategy_id: "fixed-1200-150",
		chunk_index: 2,
		text: "Thuja Green Giant and selected arborvitae cultivars are common options for evergreen screening.",
		token_count: 245,
		start_offset: 950,
		end_offset: 2100,
	},
];

const strategies = [
	{
		id: "fixed-1200-150",
		name: "Fixed 1200 / 150",
		type: "fixed",
		description: "Fixed-size chunks with overlap.",
		created_at: now,
	},
	{
		id: "semantic-512",
		name: "Semantic 512",
		type: "semantic",
		description: "Semantic paragraph grouping.",
		created_at: now,
	},
];

const corpusSources = [
	{
		source_id: "src-tree-guides",
		source_name: "Tree Center Guides",
		source_type: "web_crawl",
		document_count: 128,
		chunk_count: 4120,
		embedded_count: 3988,
	},
	{
		source_id: "src-product-pages",
		source_name: "Product Pages",
		source_type: "sitemap",
		document_count: 56,
		chunk_count: 1472,
		embedded_count: 920,
	},
];

const corpusDocuments = documents.map((document, index) => ({
	...document,
	chunk_count: index === 0 ? 42 : 31,
	embedded_count: index === 0 ? 37 : 19,
}));

const documentDetail = {
	document: {
		...documents[0],
		metadata: { category: "Guides", audience: "homeowners", region: "US" },
		content_text:
			"Choosing fast growing trees for privacy depends on mature width, climate, soil drainage, and maintenance expectations. Thuja Green Giant, Leyland Cypress, and selected arborvitae cultivars are common options for evergreen screening.",
	},
	chunks: chunks.map((chunk, index) => ({
		...chunk,
		embedding: {
			present: index === 0,
			provider_type: "openai",
			model: "text-embedding-3-small",
			dimensions: 1536,
		},
		enrichment:
			index === 0
				? {
						present: true,
						prompt_version: "summary-v1",
						short_summary: "Privacy tree selection factors.",
						quality_score: 0.9,
					}
				: { present: false },
	})),
};

const workflow = {
	ID: "wf_intake_001",
	Site: "default",
	Name: "Corpus Intake",
	Status: "running",
	Input: {
		strategy_id: "fixed-1200-150",
		source_ids: ["src-tree-guides"],
		document_ids: ["doc-privacy-trees"],
	},
	Metadata: { strategy: "fixed-1200-150" },
	CreatedAt: now,
	UpdatedAt: now,
};

const workflowGroups = [
	{
		operation: "fetch_document",
		queue: "io",
		status: "succeeded",
		count: 32,
		sample: {
			op: {
				ID: "op_fetch_001",
				WorkflowID: workflow.ID,
				Kind: "fetch_document",
				Queue: "io",
				DedupKey: "fetch/doc-privacy-trees",
				Input: {},
				DependsOn: [],
				Retry: { MaxAttempts: 3, BackoffKind: "exponential", InitialBackoff: 1000 },
				RetryState: { Attempt: 1, NextAttemptAt: null, LastError: "" },
				Metadata: null,
			},
			status: "succeeded",
			createdAt: now,
			updatedAt: now,
		},
	},
	{
		operation: "embed_chunks",
		queue: "embeddings",
		status: "running",
		count: 8,
		sample: {
			op: {
				ID: "op_embed_001",
				WorkflowID: workflow.ID,
				Kind: "embed_chunks",
				Queue: "embeddings",
				DedupKey: "embed/doc-privacy-trees",
				Input: {},
				DependsOn: [],
				Retry: { MaxAttempts: 3, BackoffKind: "exponential", InitialBackoff: 1000 },
				RetryState: { Attempt: 1, NextAttemptAt: null, LastError: "" },
				Metadata: null,
			},
			status: "running",
			createdAt: now,
			updatedAt: now,
		},
	},
];

const dslRoot: WidgetNode = component("Panel", { title: "DSL Preview" }, [
	component("Stack", { gap: "sm" }, [
		component("Caption", { tone: "accent" }, [text("Widget IR loaded from mocked Storybook API")]),
		component("Button", { variant: "primary" }, [text("Primary action")]),
	]),
]);

function jsonResponse(body: unknown): Response {
	return new Response(JSON.stringify(body), {
		status: 200,
		headers: { "content-type": "application/json" },
	});
}

function mockPayload(url: URL, method: string): unknown {
	const path = url.pathname.replace(/^\/api\/v1\/?/, "");

	if (method === "POST") {
		if (path === "embeddings/coverage")
			return { total_chunks: 5592, embedded_chunks: 4908, missing_chunks: 684 };
		if (path === "embeddings/compute")
			return {
				strategy_id: "fixed-1200-150",
				provider_type: "ollama",
				model: "nomic-embed-text",
				dimensions: 768,
				considered: 16,
				computed: 12,
				skipped_fresh: 4,
			};
		if (path === "embeddings/similarity")
			return {
				strategy_id: "fixed-1200-150",
				provider_type: "ollama",
				model: "nomic-embed-text",
				dimensions: 768,
				source: {
					chunk_id: "chunk_privacy_001",
					document_id: "doc-privacy-trees",
					strategy_id: "fixed-1200-150",
					chunk_index: 1,
					text_preview: "Choosing fast growing trees for privacy depends on mature width.",
				},
				matches: [
					{
						chunk_id: "chunk_privacy_002",
						document_id: "doc-privacy-trees",
						strategy_id: "fixed-1200-150",
						chunk_index: 2,
						score: 0.8124,
						text_preview: "Thuja Green Giant and selected arborvitae cultivars are common options.",
					},
				],
				considered: 2,
				candidate_limit: 200,
			};
		if (path.startsWith("search/"))
			return { query: "privacy trees", retriever: path.split("/").pop(), items: [], took_ms: 12 };
		if (path === "workflows/intake") return { workflow_id: workflow.ID };
		return {};
	}

	if (path === "sources") return { items: sources };
	if (path === "documents") return { items: documents };
	if (path.endsWith("/chunks")) return { items: chunks };
	if (path === "chunking-strategies") return { items: strategies };
	if (path === "corpus/sources") return { items: corpusSources };
	if (path === "corpus/documents") return { items: corpusDocuments };
	if (path.startsWith("corpus/documents/")) return documentDetail;
	if (path === "artifacts/document-processing/identities")
		return {
			items: [
				{
					artifact_type: "clean_text",
					prompt_version: "v1",
					provider: "storybook",
					model: "mock-cleaner",
				},
			],
		};
	if (path.startsWith("artifacts/document-processing/coverage"))
		return {
			items: [
				{ source_id: "src-tree-guides", total_documents: 128, processed_documents: 121 },
				{ source_id: "src-product-pages", total_documents: 56, processed_documents: 42 },
			],
		};
	if (path.startsWith("dsl/pages/")) return { id: "demo", title: "Demo DSL Page", root: dslRoot };
	if (path === "queues")
		return {
			queues: [
				{
					site: "default",
					queue: "io",
					pending: 3,
					ready: 2,
					running: 1,
					succeeded: 32,
					failed: 0,
				},
			],
		};
	if (path === "workflows") return { workflows: [{ workflow, opTotal: 40, opDone: 32 }], total: 1 };
	if (path === `workflows/${workflow.ID}`)
		return {
			workflow,
			stats: { Total: 40, Pending: 3, Ready: 2, Running: 3, Succeeded: 32, Failed: 0, Canceled: 0 },
		};
	if (path === `workflows/${workflow.ID}/ops`)
		return { workflow_id: workflow.ID, total: 40, groups: workflowGroups };
	if (path.includes("/results/"))
		return {
			OpID: "op_fetch_001",
			Data: { bytes: 4096 },
			Records: [],
			Artifacts: [],
			Emitted: [],
			EmittedIDs: [],
			Error: null,
			CompletedAt: now,
		};

	return { items: [] };
}

export function MockApiProvider({ children }: MockApiProviderProps) {
	const store = useMemo(
		() =>
			configureStore({
				reducer: { [ragApi.reducerPath]: ragApi.reducer },
				middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(ragApi.middleware),
			}),
		[],
	);

	useEffect(() => {
		const originalFetch = globalThis.fetch;
		globalThis.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
			const url = new URL(
				typeof input === "string" || input instanceof URL ? input.toString() : input.url,
				window.location.origin,
			);
			const method =
				init?.method ??
				(typeof input === "string" || input instanceof URL ? "GET" : input.method) ??
				"GET";
			return jsonResponse(mockPayload(url, method.toUpperCase()));
		};
		return () => {
			globalThis.fetch = originalFetch;
		};
	}, []);

	return <Provider store={store}>{children}</Provider>;
}
