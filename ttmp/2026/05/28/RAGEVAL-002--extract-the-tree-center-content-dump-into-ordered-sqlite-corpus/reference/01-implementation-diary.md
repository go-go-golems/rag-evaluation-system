---
Title: Implementation Diary
Ticket: RAGEVAL-002
Status: active
Topics: [rag, ingestion, database, corpus]
DocType: reference
Intent: long-term
---

# Diary

## Goal

This diary records the implementation workflow for extracting The Tree Center content from a compressed MySQL/WordPress dump into an ordered SQLite corpus for RAG evaluation.

## Step 1: Ticket Creation, Bounded Dump Inspection, and Export Design

I created a new docmgr ticket for the dump-to-SQLite corpus work and moved the investigation into ticket-local scripts so every command can be replayed. The first pass established the dump shape without loading it into MySQL and without printing raw multi-megabyte `INSERT` lines.

The key result is an implementation guide plus three replayable scripts: a bounded dump inspector, an isolated MySQL Docker Compose import path, and a MySQL-to-SQLite export script. The dump contains exactly the content types we need: published posts, guides, and products.

### Prompt Context

**User prompt (verbatim):** "so look at /home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2 and potentially setup a local docker-compose setup to load it and dump out all the articles sand guides and products straight from the db into a more ordered sqlite setup. 

Create a new ticket for that. analyze and create an implementation guide and then execute, keeping a diary, committing at appropriate intervals."

**Assistant interpretation:** Create a separate ticket for extracting TTC content directly from the compressed database dump, analyze its schema/content counts, design the import/export workflow, and then implement it with scripts and diary/changelog bookkeeping.

**Inferred user intent:** Replace or complement webpage scraping with a reproducible database-backed corpus extraction path that captures articles, guides, and products in a structured SQLite format.

### What I did

- Created ticket `RAGEVAL-002` titled "Extract The Tree Center content dump into ordered SQLite corpus".
- Added design and diary documents.
- Added tasks for dump analysis, workflow design, script implementation, and validation.
- Checked the dump path and existing TTC repository Docker Compose setup.
- Initially tried `bzgrep` for table/insert inspection, but raw dump lines are enormous and produced too much terminal output.
- Created `scripts/01-inspect-dump-schema.py` to stream the compressed dump and print bounded schema/count/sample summaries.
- Captured inspection output in `sources/01-dump-schema-inspection.md`.
- Created `scripts/02-docker-compose.mysql.yml` for an isolated MySQL 8 service on port `3347`.
- Created `scripts/02-load-dump-into-mysql.sh` to reset/import the dump into the isolated MySQL service.
- Created `scripts/03-export-mysql-to-sqlite.py` to export published `post`, `ttc_guide`, and `product` rows into `data/corpus/ttc-dump/ttc-corpus.sqlite`.
- Wrote the implementation guide at `design-doc/01-ttc-dump-to-sqlite-corpus-implementation-guide.md`.
- Related relevant source files and scripts with docmgr.
- Marked tasks 1 and 2 complete.

### Why

The database dump contains canonical WordPress IDs, product metadata, taxonomy relationships, and content status fields that are not reliably available from Defuddle scraping. A normalized SQLite corpus gives the RAG system a better source for reproducible ingestion and evaluation.

### What worked

- Bounded dump inspection worked and produced manageable output.
- The dump contains the expected public content counts:
  - `post/publish`: 483
  - `ttc_guide/publish`: 19
  - `product/publish`: 2594
- The relevant schema tables are present: `wp_posts`, `wp_postmeta`, `wp_terms`, `wp_term_taxonomy`, `wp_term_relationships`, `wp_wc_product_meta_lookup`, and `search_products`.
- The design now separates source dump import from corpus export and later RAG app ingestion.

### What didn't work

- `bzgrep` against raw multi-row inserts printed a very long line from `wp_postmeta`, creating noisy terminal output. I stopped using raw grep for INSERT inspection and replaced it with a bounded Python parser.
- I have not yet run the full MySQL import/export scripts in this step; they are implemented and ready for execution next.

### What I learned

- The compressed dump is about 43 MiB but expands into very large multi-row INSERT statements, especially for `wp_postmeta`.
- Published primary content counts in the dump align with sitemap observations for posts and guides, while products are richer in the dump than the first web-scraped corpus.
- Product variations should not be first-class RAG documents in the initial export; parent `product` rows are the correct first target.

### What was tricky to build

The main issue was output control. MySQL dumps encode many rows in single physical lines, so common tools like `grep` and `head` can still print huge chunks. The bounded inspector avoids this by parsing table creation blocks and tuple counts without echoing raw insert payloads.

The export script also avoids adding a Python MySQL dependency by asking the MySQL CLI to emit one JSON object per row. Python then reads JSON lines and writes SQLite through the standard library.

### What warrants a second pair of eyes

- Review the selected product metadata fields in `scripts/03-export-mysql-to-sqlite.py`; more `_treeinfo_*` or WooCommerce keys may be useful.
- Review whether raw `post_content` plus HTML-stripped text is enough, or whether WordPress shortcode expansion is required before chunking.
- Review the decision to exclude pages, FAQs, and product variations from the first corpus export.

### What should be done in the future

- Run the MySQL import script and capture duration/counts.
- Run the SQLite export script and validate the resulting row counts.
- Add a bridge from `ttc-corpus.sqlite` into `rag-eval` app `sources`/`documents`.
- Chunk and embed a bounded sample before processing all products.

### Code review instructions

- Start with the design guide:
  - `ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/design-doc/01-ttc-dump-to-sqlite-corpus-implementation-guide.md`
- Then review scripts in order:
  - `scripts/01-inspect-dump-schema.py`
  - `scripts/02-docker-compose.mysql.yml`
  - `scripts/02-load-dump-into-mysql.sh`
  - `scripts/03-export-mysql-to-sqlite.py`
- Re-run bounded inspection:
  - `scripts/01-inspect-dump-schema.py > sources/01-dump-schema-inspection.md`

### Technical details

- Dump path: `/home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2`.
- Isolated MySQL compose port: `3347`.
- Planned SQLite output: `data/corpus/ttc-dump/ttc-corpus.sqlite`.
- Initial corpus types: published `post`, `ttc_guide`, and `product`.
