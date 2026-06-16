__package__({ name: "explorer", short: "RAG database explorer and navigation verbs" });

const db = require("db");
const fs = require("fs");
const yaml = require("yaml");
const markdown = require("markdown");
const sanitize = require("sanitize");
const extract = require("extract");
const express = require("express");

const DEFAULT_DB = "data/rag-eval.db";

function openDatabase(database) {
	const path = database || DEFAULT_DB;
	db.configure("sqlite3", path);
	return path;
}

function sources(database) {
	openDatabase(database);
	return db.query(`
    SELECT id, name, type, config_json, created_at, updated_at
    FROM sources
    ORDER BY created_at DESC
  `);
}

__verb__("sources", {
	short: "List all corpus sources in the RAG database",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
	},
});

function documents(database, sourceId, status, limit, offset) {
	openDatabase(database);
	const n = limit || 20;
	const off = offset || 0;
	let sql = `
    SELECT d.id, d.source_id, s.name AS source_name,
           d.external_id, d.title, d.author, d.url,
           d.content_type, d.word_count, d.language, d.status,
           d.created_at, d.updated_at,
           COUNT(c.id) AS chunk_count
    FROM documents d
    LEFT JOIN sources s ON s.id = d.source_id
    LEFT JOIN chunks c ON c.document_id = d.id
   `;
	const conditions = [];
	const args = [];
	if (sourceId) {
		conditions.push("d.source_id = ?");
		args.push(sourceId);
	}
	if (status) {
		conditions.push("d.status = ?");
		args.push(status);
	}
	if (conditions.length > 0) {
		sql += " WHERE " + conditions.join(" AND ");
	}
	sql += ` GROUP BY d.id ORDER BY d.created_at DESC LIMIT ? OFFSET ?`;
	args.push(n, off);
	return db.query(sql, args);
}

__verb__("documents", {
	short: "List documents with optional source/status filtering",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		sourceId: { type: "string", default: "", help: "Filter by source ID" },
		status: { type: "string", default: "", help: "Filter by document status" },
		limit: { type: "int", default: 20, help: "Max documents to return" },
		offset: { type: "int", default: 0, help: "Skip first N documents" },
	},
});

function docDetail(database, docId) {
	openDatabase(database);
	const docRows = db.query(
		`
    SELECT d.*, s.name AS source_name
    FROM documents d
    LEFT JOIN sources s ON s.id = d.source_id
    WHERE d.id = ?
  `,
		docId,
	);
	if (docRows.length === 0) {
		throw new Error(`document not found: ${docId}`);
	}
	const doc = docRows[0];
	const chunkRows = db.query(
		`
    SELECT c.id, c.strategy_id, cs.name AS strategy_name,
           c.chunk_index, c.token_count, c.start_offset, c.end_offset,
           substr(c.text, 1, 200) AS preview
    FROM chunks c
    LEFT JOIN chunking_strategies cs ON cs.id = c.strategy_id
    WHERE c.document_id = ?
    ORDER BY c.strategy_id, c.chunk_index
  `,
		docId,
	);
	const strategies = db.query(
		`
    SELECT strategy_id, COUNT(*) AS chunk_count, SUM(token_count) AS total_tokens
    FROM chunks WHERE document_id = ? GROUP BY strategy_id
  `,
		docId,
	);
	const embeddings = db.query(
		`
    SELECT ce.strategy_id, ce.provider, ce.model, ce.dimensions, COUNT(*) AS count
    FROM chunk_embeddings ce
    JOIN chunks c ON c.id = ce.chunk_id AND c.strategy_id = ce.strategy_id
    WHERE c.document_id = ?
    GROUP BY ce.strategy_id, ce.provider, ce.model, ce.dimensions
  `,
		docId,
	);
	const artifacts = db.query(
		`
    SELECT artifact_type, prompt_version, provider, model, status, updated_at
    FROM document_processing_artifacts
    WHERE document_id = ?
    ORDER BY artifact_type, updated_at DESC
  `,
		docId,
	);
	return {
		document: doc,
		chunks: chunkRows,
		strategies: strategies,
		embeddings: embeddings,
		artifacts: artifacts,
	};
}

__verb__("docDetail", {
	short: "Show full document detail with chunks, strategies, embeddings, and artifacts",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		docId: { argument: true, help: "Document ID" },
	},
});

function chunkStrategies(database) {
	openDatabase(database);
	return db.query(`
    SELECT id, name, type, config_json, description, created_at
    FROM chunking_strategies
    ORDER BY created_at DESC
  `);
}

__verb__("chunkStrategies", {
	short: "List all chunking strategies",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
	},
});

function embeddingsCoverage(database, strategyId, provider, model, dimensions) {
	openDatabase(database);
	const strat = strategyId || "";
	const prov = provider || "";
	const mod = model || "";
	const dim = dimensions || 0;

	let sql = `
    SELECT d.source_id, s.name AS source_name,
           COUNT(c.id) AS chunk_count,
           SUM(CASE WHEN ce.chunk_id IS NOT NULL THEN 1 ELSE 0 END) AS embedded_count
    FROM chunks c
    JOIN documents d ON d.id = c.document_id
    LEFT JOIN sources s ON s.id = d.source_id
    LEFT JOIN chunk_embeddings ce ON ce.chunk_id = c.id
      AND ce.strategy_id = c.strategy_id
  `;
	const conditions = [];
	const args = [];
	if (strat) {
		conditions.push("c.strategy_id = ?");
		args.push(strat);
	}
	if (prov) {
		conditions.push("ce.provider = ?");
		args.push(prov);
	}
	if (mod) {
		conditions.push("ce.model = ?");
		args.push(mod);
	}
	if (dim) {
		conditions.push("ce.dimensions = ?");
		args.push(dim);
	}
	if (conditions.length > 0) {
		sql += " WHERE " + conditions.join(" AND ");
	}
	sql += ` GROUP BY d.source_id ORDER BY d.source_id`;
	return db.query(sql, args);
}

__verb__("embeddingsCoverage", {
	short: "Show embedding coverage grouped by source",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		strategyId: { type: "string", default: "", help: "Chunking strategy ID" },
		provider: { type: "string", default: "", help: "Embedding provider" },
		model: { type: "string", default: "", help: "Embedding model" },
		dimensions: { type: "int", default: 0, help: "Embedding dimensions" },
	},
});

function searchIndexes(database) {
	openDatabase(database);
	return db.query(`
    SELECT id, name, strategy_id, provider, model, dimensions,
           index_type, index_path, document_count, chunk_count,
           last_rebuild_at, status, created_at, updated_at
    FROM search_indexes
    ORDER BY created_at DESC
  `);
}

__verb__("searchIndexes", {
	short: "List all search indexes",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
	},
});

function docArtifacts(database, docId) {
	openDatabase(database);
	return db.query(
		`
    SELECT document_id, artifact_type, prompt_version,
           provider, model, status, error_code, error_message, updated_at
    FROM document_processing_artifacts
    WHERE document_id = ?
    ORDER BY artifact_type, updated_at DESC
  `,
		docId,
	);
}

__verb__("docArtifacts", {
	short: "List document processing artifacts for a given document",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		docId: { argument: true, help: "Document ID" },
	},
});

function exportDocs(database, sourceId, status, outDir) {
	openDatabase(database);
	const dir = outDir || "exports";
	fs.mkdirSync(dir, { recursive: true });

	let sql = `SELECT id, source_id, external_id, title, author, url,
                    content_type, word_count, language, status,
                    content_text, raw_content, metadata_json,
                    created_at, updated_at
             FROM documents`;
	const conditions = [];
	const args = [];
	if (sourceId) {
		conditions.push("source_id = ?");
		args.push(sourceId);
	}
	if (status) {
		conditions.push("status = ?");
		args.push(status);
	}
	if (conditions.length > 0) {
		sql += " WHERE " + conditions.join(" AND ");
	}
	const rows = db.query(sql, args);
	const manifest = [];
	for (const doc of rows) {
		const fileName = `${dir}/${doc.id}.json`;
		fs.writeFileSync(fileName, JSON.stringify(doc, null, 2), "utf-8");
		manifest.push({ id: doc.id, title: doc.title, file: fileName });
	}
	const manifestPath = `${dir}/manifest.json`;
	fs.writeFileSync(manifestPath, JSON.stringify(manifest, null, 2), "utf-8");
	return {
		exported: rows.length,
		directory: dir,
		manifest: manifestPath,
	};
}

__verb__("exportDocs", {
	short: "Export documents matching filters to JSON files on disk",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		sourceId: { type: "string", default: "", help: "Filter by source ID" },
		status: { type: "string", default: "", help: "Filter by document status" },
		outDir: { type: "string", default: "exports", help: "Output directory" },
	},
});

function markdownStats(database, docId) {
	openDatabase(database);
	const rows = db.query(
		`
    SELECT id, title, content_text, raw_content
    FROM documents WHERE id = ?
  `,
		docId,
	);
	if (rows.length === 0) {
		throw new Error(`document not found: ${docId}`);
	}
	const doc = rows[0];
	const text = doc.content_text || doc.raw_content || "";
	const ast = markdown.parse(text);

	let headings = 0;
	let paragraphs = 0;
	let codeBlocks = 0;
	let links = 0;
	let images = 0;
	let maxHeadingLevel = 0;

	markdown.walk(ast, (node) => {
		switch (node.Type) {
			case "heading":
				headings++;
				if (node.Level > maxHeadingLevel) maxHeadingLevel = node.Level;
				break;
			case "paragraph":
				paragraphs++;
				break;
			case "fencedCodeBlock":
			case "codeBlock":
				codeBlocks++;
				break;
			case "link":
				links++;
				break;
			case "image":
				images++;
				break;
		}
	});

	const h1Title = markdown.textContent(
		ast.Children.find((n) => n.Type === "heading" && n.Level === 1) || {},
	);

	return {
		id: doc.id,
		storedTitle: doc.title,
		extractedH1: h1Title || "",
		headings,
		maxHeadingLevel,
		paragraphs,
		codeBlocks,
		links,
		images,
		totalChars: text.length,
		totalLines: text.split(/\r?\n/).length,
	};
}

__verb__("markdownStats", {
	short: "Analyze Markdown structure of a document using goja-text markdown parser",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		docId: { argument: true, help: "Document ID" },
	},
});

function sanitizeProbe(input) {
	const text = input || "name: Alice\n  age: 30\n";
	const yResult = sanitize.yaml.sanitize(text);
	const jResult = sanitize.json.sanitize('{"ok": true, "count": 5}');
	const yFixes = yResult.Fixes || [];
	const yIssues = (yResult.Issues || yResult.LintIssues || []).filter((i) => i && i.Message);
	const jFixes = jResult.Fixes || [];
	return {
		yaml: {
			sanitized: yResult.Sanitized,
			fixes: yFixes.map((f) => f.Rule),
			issues: yIssues.map((i) => i.Message),
		},
		json: {
			sanitized: jResult.Sanitized,
			strictParseClean: jResult.StrictParseClean,
			fixes: jFixes.map((f) => f.Rule),
		},
	};
}

__verb__("sanitizeProbe", {
	short: "Probe the goja-text sanitize module (YAML + JSON repair)",
	fields: {
		input: { type: "string", default: "", help: "YAML-like input to sanitize" },
	},
});

function extractProbe(input) {
	const text = input || `Some text\n~~~json\n{"found": true}\n~~~\n<yaml>name: test</yaml>`;
	const candidates = extract.all(text);
	return candidates.map((c) => ({
		kind: c.Kind,
		format: c.Format,
		text: c.Text,
		confidence: c.Confidence,
		diagnostics: c.Diagnostics,
	}));
}

__verb__("extractProbe", {
	short: "Probe the goja-text extract module for structured data candidates",
	fields: {
		input: {
			type: "string",
			default: "",
			help: "Input text containing structured data candidates",
		},
	},
});

function serveDbBrowser(database, port) {
	openDatabase(database);
	const app = express.app();
	const p = port || 8787;

	app.get("/api/sources", (_req, res) => {
		const rows = db.query(
			`SELECT id, name, type, created_at FROM sources ORDER BY created_at DESC`,
		);
		res.json({ sources: rows });
	});

	app.get("/api/documents", (_req, res) => {
		const rows = db.query(`
      SELECT d.id, d.title, d.source_id, s.name AS source_name, d.status, d.word_count
      FROM documents d
      LEFT JOIN sources s ON s.id = d.source_id
      ORDER BY d.updated_at DESC LIMIT 50
    `);
		res.json({ documents: rows });
	});

	app.get("/api/documents/:id", (req, res) => {
		const rows = db.query(`SELECT * FROM documents WHERE id = ?`, req.params.id);
		if (rows.length === 0) {
			res.json({ error: "not found" });
			return;
		}
		const chunks = db.query(
			`
      SELECT c.id, c.strategy_id, c.chunk_index, c.token_count,
             substr(c.text, 1, 240) AS preview
      FROM chunks c WHERE c.document_id = ?
      ORDER BY c.strategy_id, c.chunk_index
    `,
			req.params.id,
		);
		res.json({ document: rows[0], chunks: chunks });
	});

	app.get("/api/search-indexes", (_req, res) => {
		const rows = db.query(`SELECT * FROM search_indexes ORDER BY created_at DESC`);
		res.json({ indexes: rows });
	});

	app.get("/", (_req, res) => {
		res.json({
			ok: true,
			endpoints: ["/api/sources", "/api/documents", "/api/documents/:id", "/api/search-indexes"],
		});
	});

	return {
		ok: true,
		port: p,
		note: `HTTP server is running. Use --keep-alive or run this verb inside a long-lived script. Endpoints: /api/sources, /api/documents, /api/documents/:id, /api/search-indexes`,
	};
}

__verb__("serveDbBrowser", {
	short: "Start an Express HTTP server that exposes key DB tables as JSON REST endpoints",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		port: { type: "int", default: 8787, help: "HTTP listen port" },
	},
});
