#!/usr/bin/env bash
# Validates source facts used by the TTC baseline evaluation-card draft.
# This is deliberately read-only: it validates the rebuilt source export before
# document revisions, a corpus snapshot, or a frozen dataset exist.
set -euo pipefail

database_path="${1:-data/ttc-wordpress-rag.sqlite}"

if [[ ! -f "$database_path" ]]; then
    echo "error: SQLite source export not found: $database_path" >&2
    exit 2
fi

assert_zero() {
    local description="$1"
    local sql="$2"
    local result
    result="$(sqlite3 "$database_path" "$sql")"
    if [[ "$result" != "0" ]]; then
        echo "FAIL: $description (unexpected count: $result)" >&2
        exit 1
    fi
    echo "PASS: $description"
}

assert_one() {
    local description="$1"
    local sql="$2"
    local result
    result="$(sqlite3 "$database_path" "$sql")"
    if [[ "$result" != "1" ]]; then
        echo "FAIL: $description (expected one result, got: $result)" >&2
        exit 1
    fi
    echo "PASS: $description"
}

echo "Validating TTC evaluation-card source facts against $database_path"

# Every document used as a provisional positive or an intentional near miss must
# exist and remain in the expected document kind. Titles are reported for human
# review below, while ID/kind are stable machine checks.
assert_zero "all draft-card document IDs resolve to expected kinds" "
WITH expected(doc_id, kind) AS (
  VALUES
    ('wp:3699','product'), ('wp:3701','product'),
    ('wp:549614','product'), ('wp:3709','product'), ('wp:552438','product'),
    ('wp:15947','product'), ('wp:3703','product'),
    ('wp:7347','product'), ('wp:26028','product'),
    ('wp:3717','product'), ('wp:10069','product'),
    ('wp:812290','ttc_guide'), ('wp:4131','faq'),
    ('wp:627148','post'), ('wp:4133','faq'), ('wp:4134','faq'),
    ('wp:9892','post'), ('wp:28084','post'),
    ('wp:398454','ttc_guide'), ('wp:405509','ttc_guide'),
    ('wp:405437','ttc_guide'), ('wp:15288','post'),
    ('wp:418694','ttc_guide'), ('wp:19387','post'), ('wp:9688','post'),
    ('wp:224522','post'), ('wp:4237','page'), ('wp:4116','faq'),
    ('wp:76495','faq'), ('wp:76497','faq'), ('wp:456943','page'),
    ('wp:558351','page'), ('wp:398600','page'),
    ('wp:270766','page'), ('wp:4140','faq'), ('wp:398551','page')
)
SELECT COUNT(*)
FROM expected e
LEFT JOIN documents d ON d.doc_id = e.doc_id
WHERE d.doc_id IS NULL OR d.kind <> e.kind;
"

# Source-text anchors. These are not retrieval tests; they prove that the
# proposed cards are grounded in source material rather than model memory.
assert_zero "all required source evidence phrases are present" "
WITH anchors(doc_id, phrase) AS (
  VALUES
    ('wp:812290','3 times as wide'),
    ('wp:812290','twice per week'),
    ('wp:627148','Light watering is worse than no watering at all'),
    ('wp:9892','after the last frost'),
    ('wp:398454','Ball and Burlap'),
    ('wp:405509','hedge'),
    ('wp:15288','below 7'),
    ('wp:19387','dense and permanent'),
    ('wp:224522','oxygen'),
    ('wp:4237','minimum'),
    ('wp:76495','unable to ship Citrus Trees to Florida'),
    ('wp:76497','not the date it will arrive'),
    ('wp:558351','cash refunds'),
    ('wp:270766','picture and description')
)
SELECT COUNT(*)
FROM anchors a
JOIN documents d ON d.doc_id = a.doc_id
WHERE instr(lower(d.search_text), lower(a.phrase)) = 0;
"

# Exact product-attribute cards must remain constrained, not merely plausible.
assert_one "Blue Ice Arizona Cypress constrained discovery is unique" "
SELECT COUNT(*)
FROM view_products
WHERE categories LIKE '%Privacy Trees%'
  AND drought_tolerance = 'Very Drought Resistant'
  AND sunlight = 'Full Sun'
  AND mature_height = '15-25'
  AND mature_width = '6-8';
"

assert_one "Danica Globe Thuja dimensions and taxonomy identify one product" "
SELECT COUNT(*)
FROM view_products
WHERE categories LIKE '%Thuja Trees%'
  AND mature_height = '1-2'
  AND mature_width = '1-2';
"

assert_one "Bald Cypress wet-soil height constraint identifies one product" "
SELECT COUNT(*)
FROM view_products
WHERE categories LIKE '%Cypress Trees%'
  AND attributes LIKE '%pa_special-features: Tolerates Wet Soil%'
  AND mature_height = '50-70';
"

assert_zero "Bitcoin has no corpus FTS hit for the explicitly unanswerable card" "
SELECT COUNT(*) FROM documents_fts WHERE documents_fts MATCH 'bitcoin';
"

echo
echo "Source identity review (titles are intentionally displayed for an adjudicator):"
sqlite3 -header -column "$database_path" "
SELECT doc_id, kind, title
FROM documents
WHERE doc_id IN (
  'wp:3699','wp:549614','wp:15947','wp:7347','wp:3717',
  'wp:812290','wp:9892','wp:398454','wp:405509','wp:15288',
  'wp:19387','wp:224522','wp:4237','wp:76495','wp:76497',
  'wp:558351','wp:270766'
)
ORDER BY doc_id;
"

echo "PASS: TTC baseline evaluation-card source validation completed"
