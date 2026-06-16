__package__({ name: "database", short: "Database-oriented RAG experiment verbs" });

const db = require("db");
const markdown = require("markdown");
const fs = require("fs");
const yaml = require("yaml");

const DEFAULT_DB = "data/rag-eval.db";

function openDatabase(database) {
	const path = database || DEFAULT_DB;
	db.configure("sqlite3", path);
	return path;
}

function firstNonEmptyLine(text) {
	const lines = String(text || "").split(/\r?\n/);
	for (const line of lines) {
		const trimmed = line.trim();
		if (trimmed.length > 0) return trimmed.replace(/^#+\s*/, "");
	}
	return "";
}

function markdownTitle(text) {
	const ast = markdown.parse(String(text || ""));
	let title = "";
	markdown.walk(ast, (node) => {
		if (title) return;
		if (node.Type === "heading" && node.Level === 1) {
			title = markdown.textContent(node);
		}
	});
	return title || firstNonEmptyLine(text);
}

function topDocs(database, limit) {
	openDatabase(database);
	const n = limit || 10;
	return db.query(
		`
    SELECT
      d.id,
      d.title,
      d.source_id,
      s.name AS source_name,
      d.word_count,
      LENGTH(COALESCE(d.content_text, d.raw_content, '')) AS bytes,
      COUNT(c.id) AS chunk_count,
      d.status,
      d.updated_at
    FROM documents d
    LEFT JOIN sources s ON s.id = d.source_id
    LEFT JOIN chunks c ON c.document_id = d.id
    GROUP BY d.id
    ORDER BY d.word_count DESC, bytes DESC
    LIMIT ?
  `,
		n,
	);
}

__verb__("topDocs", {
	short: "List the top documents in the RAG SQLite database",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		limit: { type: "int", default: 10, help: "Maximum number of documents" },
	},
});

function docTitle(database, docId) {
	openDatabase(database);
	const rows = db.query(
		`
    SELECT id, title AS stored_title, source_id, content_text, raw_content, metadata_json
    FROM documents
    WHERE id = ?
  `,
		docId,
	);
	if (rows.length === 0) {
		throw new Error(`document not found: ${docId}`);
	}
	const doc = rows[0];
	const text = doc.content_text || doc.raw_content || "";
	const extractedTitle = markdownTitle(text);
	const ast = markdown.parse(text);
	const validation = markdown.validate(ast);
	return {
		id: doc.id,
		sourceId: doc.source_id,
		storedTitle: doc.stored_title,
		extractedMarkdownTitle: extractedTitle,
		titleMatches: String(doc.stored_title || "").trim() === String(extractedTitle || "").trim(),
		markdownValid: validation.Valid,
		contentBytes: text.length,
		preview: text.slice(0, 300),
	};
}

__verb__("docTitle", {
	short: "Extract a document title by parsing document Markdown from SQLite",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		docId: { argument: true, help: "documents.id value, for example ttc-article-9838" },
	},
});

function chunkPreview(database, docId, strategyId, limit) {
	openDatabase(database);
	return db.query(
		`
    SELECT
      c.id,
      c.document_id,
      c.strategy_id,
      c.chunk_index,
      c.token_count,
      c.start_offset,
      c.end_offset,
      substr(c.text, 1, 240) AS preview
    FROM chunks c
    WHERE c.document_id = ?
      AND (? = '' OR c.strategy_id = ?)
    ORDER BY c.strategy_id, c.chunk_index
    LIMIT ?
  `,
		docId,
		strategyId || "",
		strategyId || "",
		limit || 10,
	);
}

__verb__("chunkPreview", {
	short: "Show chunk previews for a document and optional strategy",
	fields: {
		database: { type: "string", default: "data/rag-eval.db", help: "SQLite database path" },
		docId: { argument: true, help: "documents.id value" },
		strategyId: { type: "string", default: "", help: "Optional chunks.strategy_id filter" },
		limit: { type: "int", default: 10, help: "Maximum chunks to return" },
	},
});

function configProbe(file) {
	const text = fs.readFileSync(file, "utf-8");
	const parsed = yaml.parse(text);
	return {
		file,
		valid: yaml.validate(text).valid,
		topLevelKeys: Object.keys(parsed || {}),
		parsed,
	};
}

__verb__("configProbe", {
	short: "Read and parse a YAML config file to prove fs + yaml are available",
	fields: {
		file: { argument: true, help: "YAML file to read" },
	},
});
