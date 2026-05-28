#!/usr/bin/env python3
"""Export TTC WordPress content from local MySQL into an ordered SQLite corpus.

Reads JSON rows produced by the mysql CLI to avoid adding Python MySQL driver
dependencies. Output is deterministic and optimized for downstream RAG ingestion:
articles, guides, and products become rows in content_items; product attributes
and taxonomy terms are broken out into side tables.
"""

from __future__ import annotations

import argparse
import html
import json
import re
import sqlite3
import subprocess
from html.parser import HTMLParser
from pathlib import Path
from typing import Any

SCRIPT_DIR = Path(__file__).resolve().parent
COMPOSE_FILE = SCRIPT_DIR / "02-docker-compose.mysql.yml"


class TextExtractor(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.parts: list[str] = []
        self.skip_depth = 0

    def handle_starttag(self, tag: str, attrs):
        if tag in {"script", "style", "noscript"}:
            self.skip_depth += 1
        elif tag in {"p", "br", "div", "section", "article", "li", "h1", "h2", "h3", "h4", "h5", "h6"}:
            self.parts.append("\n")

    def handle_endtag(self, tag: str):
        if tag in {"script", "style", "noscript"} and self.skip_depth:
            self.skip_depth -= 1
        elif tag in {"p", "div", "section", "article", "li", "h1", "h2", "h3", "h4", "h5", "h6"}:
            self.parts.append("\n")

    def handle_data(self, data: str):
        if not self.skip_depth:
            self.parts.append(data)

    def text(self) -> str:
        text = html.unescape("".join(self.parts))
        text = re.sub(r"\r\n?", "\n", text)
        text = re.sub(r"[ \t]+", " ", text)
        text = re.sub(r"\n{3,}", "\n\n", text)
        return text.strip()


def html_to_text(value: str) -> str:
    parser = TextExtractor()
    parser.feed(value or "")
    return parser.text()


def run_mysql_json(query: str) -> list[dict[str, Any]]:
    cmd = [
        "docker", "compose", "-f", str(COMPOSE_FILE), "exec", "-T", "mysql",
        "mysql", "-uroot", "-psomewordpress", "-N", "-B", "--default-character-set=utf8mb4", "ttc", "-e", query,
    ]
    result = subprocess.run(cmd, check=True, text=True, capture_output=True)
    rows: list[dict[str, Any]] = []
    for line in result.stdout.splitlines():
        if not line.strip():
            continue
        rows.append(json.loads(line))
    return rows


def init_sqlite(path: Path) -> sqlite3.Connection:
    path.parent.mkdir(parents=True, exist_ok=True)
    if path.exists():
        path.unlink()
    conn = sqlite3.connect(path)
    conn.executescript(
        """
        PRAGMA journal_mode=WAL;
        PRAGMA foreign_keys=ON;

        CREATE TABLE content_items (
          id TEXT PRIMARY KEY,
          wp_id INTEGER NOT NULL UNIQUE,
          kind TEXT NOT NULL,
          post_type TEXT NOT NULL,
          status TEXT NOT NULL,
          slug TEXT NOT NULL,
          title TEXT NOT NULL,
          url_path TEXT NOT NULL,
          published_at TEXT,
          modified_at TEXT,
          excerpt TEXT NOT NULL DEFAULT '',
          content_html TEXT NOT NULL DEFAULT '',
          content_text TEXT NOT NULL DEFAULT '',
          word_count INTEGER NOT NULL DEFAULT 0,
          parent_wp_id INTEGER NOT NULL DEFAULT 0,
          metadata_json TEXT NOT NULL DEFAULT '{}'
        );

        CREATE TABLE item_terms (
          item_id TEXT NOT NULL REFERENCES content_items(id) ON DELETE CASCADE,
          taxonomy TEXT NOT NULL,
          term_id INTEGER NOT NULL,
          term_slug TEXT NOT NULL,
          term_name TEXT NOT NULL,
          PRIMARY KEY (item_id, taxonomy, term_id)
        );

        CREATE TABLE product_meta (
          item_id TEXT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
          sku TEXT,
          min_price REAL,
          max_price REAL,
          stock_status TEXT,
          botanical_name TEXT,
          hardiness_zone TEXT,
          mature_height TEXT,
          mature_width TEXT,
          sunlight TEXT,
          soil_conditions TEXT,
          drought_tolerance TEXT
        );

        CREATE INDEX idx_content_items_kind ON content_items(kind);
        CREATE INDEX idx_content_items_post_type ON content_items(post_type);
        CREATE INDEX idx_content_items_slug ON content_items(slug);
        CREATE INDEX idx_item_terms_taxonomy ON item_terms(taxonomy, term_slug);
        """
    )
    return conn


def content_query() -> str:
    return r"""
SELECT JSON_OBJECT(
  'wp_id', p.ID,
  'post_type', p.post_type,
  'status', p.post_status,
  'slug', p.post_name,
  'title', p.post_title,
  'published_at', p.post_date,
  'modified_at', p.post_modified,
  'excerpt', p.post_excerpt,
  'content_html', p.post_content,
  'parent_wp_id', p.post_parent,
  'guid', p.guid,
  'metadata', JSON_OBJECT(
    'botanical_name', botanical.meta_value,
    'hardiness_zone', hardiness.meta_value,
    'mature_height', height.meta_value,
    'mature_width', width.meta_value,
    'sunlight', sunlight.meta_value,
    'soil_conditions', soil.meta_value,
    'drought_tolerance', drought.meta_value,
    'thumbnail_id', thumbnail.meta_value
  )
)
FROM wp_posts p
LEFT JOIN wp_postmeta botanical ON botanical.post_id=p.ID AND botanical.meta_key='_treeinfo_botanical_name'
LEFT JOIN wp_postmeta hardiness ON hardiness.post_id=p.ID AND hardiness.meta_key='_treeinfo_hardiness_zone'
LEFT JOIN wp_postmeta height ON height.post_id=p.ID AND height.meta_key='_treeinfo_mature_height'
LEFT JOIN wp_postmeta width ON width.post_id=p.ID AND width.meta_key='_treeinfo_mature_width'
LEFT JOIN wp_postmeta sunlight ON sunlight.post_id=p.ID AND sunlight.meta_key='_treeinfo_sunlight'
LEFT JOIN wp_postmeta soil ON soil.post_id=p.ID AND soil.meta_key='_treeinfo_soil_conditions'
LEFT JOIN wp_postmeta drought ON drought.post_id=p.ID AND drought.meta_key='_treeinfo_drought_tolerance'
LEFT JOIN wp_postmeta thumbnail ON thumbnail.post_id=p.ID AND thumbnail.meta_key='_thumbnail_id'
WHERE p.post_status='publish' AND p.post_type IN ('post','ttc_guide','product')
ORDER BY FIELD(p.post_type,'ttc_guide','post','product'), p.post_date, p.ID;
"""


def terms_query() -> str:
    return r"""
SELECT JSON_OBJECT(
  'wp_id', p.ID,
  'taxonomy', tt.taxonomy,
  'term_id', t.term_id,
  'term_slug', t.slug,
  'term_name', t.name
)
FROM wp_posts p
JOIN wp_term_relationships tr ON tr.object_id=p.ID
JOIN wp_term_taxonomy tt ON tt.term_taxonomy_id=tr.term_taxonomy_id
JOIN wp_terms t ON t.term_id=tt.term_id
WHERE p.post_status='publish' AND p.post_type IN ('post','ttc_guide','product')
ORDER BY p.ID, tt.taxonomy, t.name;
"""


def product_lookup_query() -> str:
    return r"""
SELECT JSON_OBJECT(
  'wp_id', product_id,
  'sku', sku,
  'min_price', min_price,
  'max_price', max_price,
  'stock_status', stock_status
)
FROM wp_wc_product_meta_lookup
ORDER BY product_id;
"""


def kind_for_post_type(post_type: str) -> str:
    return {"post": "article", "ttc_guide": "guide", "product": "product"}.get(post_type, post_type)


def url_path_for(row: dict[str, Any]) -> str:
    post_type = row["post_type"]
    slug = row["slug"]
    if post_type == "ttc_guide":
        return f"/guides/{slug}/"
    return f"/{slug}/"


def export(out_path: Path) -> None:
    conn = init_sqlite(out_path)
    product_lookup = {int(r["wp_id"]): r for r in run_mysql_json(product_lookup_query())}
    rows = run_mysql_json(content_query())

    for row in rows:
        wp_id = int(row["wp_id"])
        post_type = row["post_type"]
        kind = kind_for_post_type(post_type)
        content_html = row.get("content_html") or ""
        text_parts = [row.get("title") or "", row.get("excerpt") or "", html_to_text(content_html)]
        content_text = "\n\n".join(part.strip() for part in text_parts if part and part.strip())
        item_id = f"ttc-{kind}-{wp_id}"
        metadata = row.get("metadata") or {}
        conn.execute(
            """
            INSERT INTO content_items (
              id, wp_id, kind, post_type, status, slug, title, url_path,
              published_at, modified_at, excerpt, content_html, content_text,
              word_count, parent_wp_id, metadata_json
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                item_id, wp_id, kind, post_type, row["status"], row["slug"], row["title"], url_path_for(row),
                row.get("published_at"), row.get("modified_at"), row.get("excerpt") or "", content_html,
                content_text, len(content_text.split()), int(row.get("parent_wp_id") or 0), json.dumps(metadata, ensure_ascii=False),
            ),
        )
        if kind == "product":
            lookup = product_lookup.get(wp_id, {})
            conn.execute(
                """
                INSERT INTO product_meta (
                  item_id, sku, min_price, max_price, stock_status, botanical_name,
                  hardiness_zone, mature_height, mature_width, sunlight, soil_conditions, drought_tolerance
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    item_id, lookup.get("sku"), lookup.get("min_price"), lookup.get("max_price"), lookup.get("stock_status"),
                    metadata.get("botanical_name"), metadata.get("hardiness_zone"), metadata.get("mature_height"),
                    metadata.get("mature_width"), metadata.get("sunlight"), metadata.get("soil_conditions"), metadata.get("drought_tolerance"),
                ),
            )

    for term in run_mysql_json(terms_query()):
        kind = None
        wp_id = int(term["wp_id"])
        existing = conn.execute("SELECT id FROM content_items WHERE wp_id=?", (wp_id,)).fetchone()
        if not existing:
            continue
        conn.execute(
            "INSERT OR IGNORE INTO item_terms (item_id, taxonomy, term_id, term_slug, term_name) VALUES (?, ?, ?, ?, ?)",
            (existing[0], term["taxonomy"], int(term["term_id"]), term["term_slug"], term["term_name"]),
        )

    conn.commit()
    print(f"wrote {out_path}")
    for row in conn.execute("SELECT kind, COUNT(*), SUM(word_count) FROM content_items GROUP BY kind ORDER BY kind"):
        print(f"{row[0]}\titems={row[1]}\twords={row[2]}")
    conn.close()


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--out", default="data/corpus/ttc-dump/ttc-corpus.sqlite")
    args = parser.parse_args()
    export(Path(args.out))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
