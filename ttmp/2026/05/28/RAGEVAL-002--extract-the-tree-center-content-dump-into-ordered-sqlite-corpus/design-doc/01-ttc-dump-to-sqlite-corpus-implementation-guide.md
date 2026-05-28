---
Title: ""
Ticket: ""
Status: ""
Topics: []
DocType: ""
Intent: ""
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/ttc/ttc/docker-compose.yml
      Note: Existing TTC compose setup reviewed for MySQL settings and port conflict avoidance
    - Path: ../../../../../../../../../../code/ttc/ttc/ttc_dev_dump.sql.bz2
      Note: Source MySQL dump for TTC corpus extraction
    - Path: ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/01-inspect-dump-schema.py
      Note: Bounded dump inspection script
    - Path: ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/02-load-dump-into-mysql.sh
      Note: Isolated MySQL dump import script
    - Path: ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/03-export-mysql-to-sqlite.py
      Note: Normalized SQLite corpus export script
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# TTC Dump to SQLite Corpus Implementation Guide

## Executive summary

RAGEVAL-002 extracts The Tree Center content directly from `/home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2` into a deterministic SQLite corpus for RAG evaluation. The dump is a compressed MySQL/WordPress dump. It contains the same public guide and blog-post counts observed through the sitemap, plus product records and WordPress/WooCommerce metadata that are not available cleanly through page scraping.

The proposed workflow is intentionally reproducible and bounded:

1. Inspect the compressed dump with bounded scripts that never print raw `INSERT` lines.
2. Load the dump into an isolated MySQL 8 Docker Compose service only when full SQL semantics are needed.
3. Export published `post`, `ttc_guide`, and `product` records into a normalized SQLite database.
4. Preserve taxonomy terms and selected product attributes in side tables.
5. Use the SQLite corpus as a source for RAG ingestion, chunking, embeddings, search, and evaluation.

All investigation and execution scripts live in this ticket's `scripts/` directory so the workflow can be replayed.

## Problem statement

The previous corpus workflow downloaded The Tree Center guides and posts from public URLs using Defuddle. That is useful for validating webpage extraction, but it misses several properties that matter for a richer RAG corpus:

- product metadata such as botanical name, hardiness zone, mature size, sunlight, soil, drought tolerance, SKU, price, and stock status;
- canonical WordPress IDs and slugs;
- draft/private/trash visibility information needed to decide what should enter the corpus;
- taxonomy relationships for categories and product categories;
- reproducible local extraction independent of website availability and page templates.

The compressed dump gives a better source of truth, but raw WordPress/WooCommerce tables are not shaped for RAG. The implementation must convert them into an ordered SQLite corpus with predictable IDs and inspectable fields.

## Evidence from dump inspection

The dump path is:

```text
/home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2
```

The compressed file size is about 43 MiB.

Bounded inspection script:

```text
ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/01-inspect-dump-schema.py
```

The script found these relevant WordPress/WooCommerce tables:

- `wp_posts`
- `wp_postmeta`
- `wp_terms`
- `wp_term_taxonomy`
- `wp_term_relationships`
- `wp_wc_product_meta_lookup`
- `search_products`

Relevant published content counts from `wp_posts`:

| post_type | status | count |
|---|---:|---:|
| `post` | `publish` | 483 |
| `ttc_guide` | `publish` | 19 |
| `product` | `publish` | 2594 |
| `page` | `publish` | 120 |
| `product_variation` | `publish` | 11913 |
| `faq` | `publish` | 35 |

The first corpus target is:

| kind | source post_type | include now | reason |
|---|---|---:|---|
| guide | `ttc_guide` | yes | The guide corpus is small and already validated through Defuddle. |
| article | `post` | yes | Blog posts are the main prose corpus. |
| product | `product` | yes | Product records add structured plant facts and commerce metadata. |
| page | `page` | no | Many pages are navigation/account/cart/system pages. Add later if needed. |
| variation | `product_variation` | no | Variations are SKU/size options, not primary RAG documents. |
| attachment | `attachment` | no | Media metadata should not become text corpus documents. |

## Target SQLite schema

The export creates a database at:

```text
data/corpus/ttc-dump/ttc-corpus.sqlite
```

This path is under `data/`, so it is intentionally ignored by Git.

### `content_items`

One row per article, guide, or product.

Important columns:

- `id`: stable local ID, e.g. `ttc-guide-398454`.
- `wp_id`: original WordPress post ID.
- `kind`: `article`, `guide`, or `product`.
- `post_type`: original WordPress post type.
- `status`: original WordPress status, currently filtered to `publish`.
- `slug`: WordPress slug.
- `title`: title.
- `url_path`: canonical path approximation.
- `published_at`, `modified_at`.
- `excerpt`.
- `content_html`: original post content.
- `content_text`: HTML-stripped text plus title/excerpt.
- `word_count`.
- `metadata_json`: selected extra metadata as JSON.

### `item_terms`

Taxonomy terms attached to each exported item.

Important columns:

- `item_id`
- `taxonomy`
- `term_id`
- `term_slug`
- `term_name`

### `product_meta`

Selected product-specific fields.

Important columns:

- `item_id`
- `sku`
- `min_price`
- `max_price`
- `stock_status`
- `botanical_name`
- `hardiness_zone`
- `mature_height`
- `mature_width`
- `sunlight`
- `soil_conditions`
- `drought_tolerance`

## Workflow scripts

All scripts are stored in this ticket workspace.

### 1. Inspect dump safely

```bash
ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/01-inspect-dump-schema.py \
  > ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/sources/01-dump-schema-inspection.md
```

This script is the replacement for unsafe `bzgrep` over raw insert lines. It prints bounded summaries only.

### 2. Start/load isolated MySQL

Compose file:

```text
ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/02-docker-compose.mysql.yml
```

Import command:

```bash
ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/02-load-dump-into-mysql.sh \
  /home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2
```

The isolated MySQL service uses port `3347` to avoid the existing TTC project port `3336`.

### 3. Export normalized SQLite

```bash
ttmp/2026/05/28/RAGEVAL-002--extract-the-tree-center-content-dump-into-ordered-sqlite-corpus/scripts/03-export-mysql-to-sqlite.py \
  --out data/corpus/ttc-dump/ttc-corpus.sqlite
```

## Design decisions

### Use MySQL for dump execution, not a direct SQL translator

The dump is MySQL-specific and includes InnoDB table definitions, MySQL comments, character set declarations, indexes, and large multi-row inserts. A direct MySQL-to-SQLite translator would need to emulate too much. Loading into MySQL first is simpler and more faithful.

### Export only published primary documents first

The first export includes `post`, `ttc_guide`, and `product` with `post_status='publish'`. Draft/private/trash records remain available in MySQL if later analysis needs them, but they should not enter the initial RAG corpus.

### Keep product variations out of `content_items`

Product variations are useful commerce rows, but they are not primary text documents. They can be modeled later as structured product options. The first corpus keeps one document per published parent product.

### Store HTML and derived text

`content_html` preserves the source. `content_text` provides a downstream chunking input. This makes text extraction inspectable and reversible enough for early experiments.

### Keep SQLite export independent from the RAG app schema

The export database is a corpus database, not the app's operational `data/rag-eval.db`. A later ingestion step can read `content_items` and insert into the app's `sources`/`documents` tables. Keeping these separate avoids coupling dump normalization to app runtime migrations.

## Implementation plan

1. Create ticket workspace and tasks.
2. Create bounded dump inspection script.
3. Capture schema/count evidence into `sources/01-dump-schema-inspection.md`.
4. Create isolated MySQL Compose file.
5. Create dump import script.
6. Create MySQL-to-SQLite export script.
7. Import dump into MySQL.
8. Export SQLite corpus.
9. Validate counts:
   - expected `guide`: 19;
   - expected `article`: 483;
   - expected `product`: 2594.
10. Add an ingestion bridge from `ttc-corpus.sqlite` into `rag-eval` app documents.
11. Chunk and embed a bounded sample before scaling.

## Validation commands

```bash
sqlite3 data/corpus/ttc-dump/ttc-corpus.sqlite \
  "SELECT kind, COUNT(*), SUM(word_count) FROM content_items GROUP BY kind ORDER BY kind;"

sqlite3 data/corpus/ttc-dump/ttc-corpus.sqlite \
  "SELECT taxonomy, COUNT(*) FROM item_terms GROUP BY taxonomy ORDER BY taxonomy;"

sqlite3 data/corpus/ttc-dump/ttc-corpus.sqlite \
  "SELECT title, botanical_name, hardiness_zone, mature_height, sunlight FROM product_meta JOIN content_items ON content_items.id=product_meta.item_id LIMIT 10;"
```

## Risks and open questions

- WordPress shortcodes may remain in `content_text`; later extraction may need WordPress-aware shortcode cleanup.
- Product descriptions may depend on theme/plugin rendering; raw `post_content` may not match the public page exactly.
- Product metadata keys may have additional useful fields beyond the initial selected `_treeinfo_*` subset.
- The full product corpus is much larger than guides/posts and should be embedded only after coverage/cost controls exist.
