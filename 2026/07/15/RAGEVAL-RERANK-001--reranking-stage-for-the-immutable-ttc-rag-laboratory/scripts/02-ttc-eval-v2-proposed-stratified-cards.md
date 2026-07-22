# TTC evaluation expansion: proposed stratified cards

Status: **authoring draft only**. This is deliberately not a dataset, not a
source-of-truth judgment file, and not input to an immutable run. It proposes
72 new cards to add to the TTC development/holdout pipeline after independent
source validation, candidate pooling, and blind human adjudication.

Each `expected_document_ids` entry is a source-document target discovered in
the current SQLite export. `?` means the document title/source target is known
but the claimed answer facet has not been manually checked in its full text.
Every card needs revision resolution, exact evidence ranges, and relevance
grading before it can be scored. Keep policy-conflict cards withheld.

## Proposed cards

| ID | Category / intent | Query | Expected document IDs | Answerability | Difficulty / rationale |
|---|---|---|---|---|---|
| v2-001 | support / bulk-order | Do you offer discounts for large wholesale tree orders? | `wp:4123` | answerable | easy; direct FAQ wording |
| v2-002 | support / gift-order | Can I send a tree as a gift to someone else? | `wp:4122` | answerable | easy; direct FAQ |
| v2-003 | support / order-status | How can I check the status of my order? | `wp:4124`, `wp:16771` | answerable | medium; FAQ versus account/status page |
| v2-004 | support / shipment-notification | How will I know when my order has shipped? | `wp:4126` | answerable | easy; lexical directness |
| v2-005 | support / shipping-cost | How much are shipping and handling charges? | `wp:4121` | answerable | easy; direct FAQ |
| v2-006 | support / shipping-method | How are trees packed and shipped to customers? | `wp:76492`, `wp:4235` | answerable | medium; overlapping shipping pages |
| v2-007 | support / order-edit | Can I edit an order after placing it? | `wp:69806` | answerable | easy; direct FAQ |
| v2-008 | support / delayed-shipping | Can I delay shipment of an order I already placed? | `wp:4127` | answerable | easy; near cancellation/scheduling confusion |
| v2-009 | support / order-timing | When will my order ship? | `wp:4125`, `wp:217834` | answerable | medium; ordinary versus preorder timing |
| v2-010 | support / presale | What does it mean when a plant is on pre-sale? | `wp:76498`, `wp:846300` | answerable | medium; FAQ and editorial explanation |
| v2-011 | support / wrong-item | I received the wrong plant—what should I do? | `wp:4138` | answerable | easy; direct incident FAQ |
| v2-012 | support / wilted-arrival | My shrubs arrived wilted; what should I do first? | `wp:4137`, `wp:270766` | answerable | medium; symptom versus guarantee evidence |
| v2-013 | support / size-expectation | The plant I received is smaller than expected. Is that normal? | `wp:4139`, `wp:76490` | answerable | medium; product-size semantics |
| v2-014 | support / gallon-size | When a plant is sold in gallons, does that tell me its height? | `wp:76490` | answerable | easy; controlled terminology |
| v2-015 | support / payment | Which forms of payment does The Tree Center accept? | `wp:4120`, `wp:398551` | answerable | medium; FAQ versus phone-order page |
| v2-016 | support / catalog | Does The Tree Center have a printed catalog? | `wp:4143` | answerable | easy; direct FAQ |
| v2-017 | support / location | Where is The Tree Center located? | `wp:76336`, `wp:4160` | answerable | medium; FAQ versus About page |
| v2-018 | support / escalation | My question is not listed in the FAQ. How do I get help? | `wp:4136`, `wp:4144` | answerable | easy; support-routing |
| v2-019 | support / extension | Where can I find my local USDA Extension Office? | `wp:4146` | answerable | easy; named resource |
| v2-020 | support / policy-conflict | What cancellation fee applies to my order? | `wp:398597`, `wp:4128`, `wp:4129` | **withhold: policy conflict** | hard; conflicting policy sources require owner precedence |
| v2-021 | care / staking | Does a newly planted tree need to be staked? | `wp:4132`, `wp:10052`, `wp:39184` | answerable | medium; FAQ, how-to, and explanatory article |
| v2-022 | care / staking | Why can staking a tree be necessary, and when is it harmful? | `wp:39184`, `wp:10052` | answerable | hard; causal explanation, not mere procedure |
| v2-023 | care / pests | I think my tree has a pest or fungus—what should I do? | `wp:4135`, `wp:753964` | answerable | medium; triage versus general diagnosis |
| v2-024 | care / sick-plant | My tree looks sick. What information should I collect before asking for help? | `wp:4142`, `wp:4135` | answerable | medium; nearby support pages |
| v2-025 | care / leaf-curl | Why are my tree leaves curling and turning brown? | `wp:4134`, `wp:627148` | answerable | medium; adversarial opposite of yellow-leaf overwatering |
| v2-026 | care / tree-death | My tree died after planting. What support or guarantee path applies? | `wp:4140`, `wp:270766` | answerable | medium; navigation plus policy |
| v2-027 | care / deciduous-planting | How do I plant a deciduous tree? | `wp:405420`, `wp:9167` | answerable | medium; guide/page duplicate family |
| v2-028 | care / evergreen-planting | How do I plant an evergreen tree? | `wp:405433`, `wp:9175` | answerable | medium; guide/page duplicate family |
| v2-029 | care / bare-root | What is the correct procedure for planting a bare-root tree? | `wp:405431`, `wp:9171` | answerable | medium; discriminates root treatments |
| v2-030 | care / bamboo-planting | How should I plant bamboo trees? | `wp:398536`, `wp:9204` | answerable | medium; guide/page duplicate family |
| v2-031 | care / citrus-planting | How should I plant a citrus tree? | `wp:418611`, `wp:9212` | answerable | medium; guide/page duplicate family |
| v2-032 | care / rose-planting | How should I plant roses? | `wp:418613`, `wp:9214` | answerable | medium; guide/page duplicate family |
| v2-033 | care / hydrangea-planting | How should I plant hydrangeas? | `wp:418605`, `wp:9202` | answerable | medium; guide/page duplicate family |
| v2-034 | care / japanese-maple-planting | How should I plant a Japanese maple? | `wp:418603`, `wp:9194` | answerable | medium; named-plant procedure |
| v2-035 | care / apple-pear-planting | How should I plant apple and pear trees? | `wp:418609`, `wp:9210` | answerable | medium; multi-entity task |
| v2-036 | care / peach-nectarine-planting | How should I plant peach and nectarine trees? | `wp:398553`, `wp:9208` | answerable | medium; multi-entity task |
| v2-037 | care / mulch | How should mulch be used around trees and shrubs? | `wp:37719` | answerable | medium; general practice and likely distractors |
| v2-038 | care / watering-outdoors | What are the basic principles for watering plants outdoors? | `wp:627148` | answerable | medium; broad explanatory retrieval |
| v2-039 | care / watering-containers | How is watering pots and planters different from watering outdoor plants? | `wp:627151`, `wp:627148` | answerable | hard; complementary two-document comparison |
| v2-040 | care / drought | What should I do when I cannot water my garden? | `wp:693642` | answerable | medium; paraphrase-heavy |
| v2-041 | care / winter-containers | How can I keep shrubs and trees in planters alive through winter? | `wp:704216`, `wp:642332` | answerable | medium; adjacent overwintering sources |
| v2-042 | care / winter-roses | How do I protect roses and hydrangeas over winter? | `wp:540727`, `wp:643793` | answerable | medium; overlapping rose-specific source |
| v2-043 | care / panicle-hydrangea | When and how should I prune panicle hydrangeas? | `wp:690769` | answerable | easy; highly specific title |
| v2-044 | care / wisteria | How should I prune wisteria? | `wp:507449` | answerable | easy; highly specific title |
| v2-045 | care / grape-pruning | How should grape vines be pruned correctly? | `wp:691485` | answerable | easy; precise terminology |
| v2-046 | care / spring-shrubs | Which shrubs should be pruned in spring? | `wp:16177`, `wp:23006`, `wp:23259` | answerable | hard; series plus categorical question |
| v2-047 | care / young-tree-pruning | How should young trees be pruned for long-term strength? | `wp:751617`, `wp:751625`, `wp:751629` | answerable | hard; three-part series / multi-hop |
| v2-048 | editorial / soil-basics | What are the basic things a gardener should understand about soil? | `wp:694920`, `wp:27842` | answerable | medium; broad topical overlap |
| v2-049 | editorial / fall-soil | How can I improve soil before fall planting? | `wp:51194` | answerable | medium; seasonal constraint |
| v2-050 | editorial / alkaline-soil | What trees and shrubs are suitable for alkaline soil? | `wp:545207`, `wp:548585`, `wp:26974` | answerable | hard; series and an adversarial acid-soil item |
| v2-051 | editorial / acidic-soil | Can rhododendrons be grown in alkaline soil? | `wp:26974`, `wp:418694` | answerable | hard; contradiction/prescription distinction |
| v2-052 | editorial / coastal-plants | What plants or trees are good choices near the coast or beach? | `wp:15584`, `wp:63934` | answerable | medium; synonymous corpus vocabulary |
| v2-053 | editorial / deer | What plants and shrubs are resistant to deer? | `wp:6625` | answerable | easy; direct topic |
| v2-054 | editorial / shade | How do I choose a fast-growing shade tree? | `wp:24788`, `wp:4355`, `wp:8231` | answerable | hard; selection trade-offs and multiple candidate articles |
| v2-055 | editorial / hosta | Which hostas can grow in hot zones? | `wp:691134`, `wp:704419` | answerable | medium; niche botanical constraint |
| v2-056 | editorial / xeric | What is xeric gardening and how does it reduce water use? | `wp:38062`, `wp:539163` | answerable | medium; definition plus recommendations |
| v2-057 | editorial / privacy | Which trees make good backyard privacy screens? | `wp:4332`, `wp:8017`, `wp:640339` | answerable | hard; high duplicate/near-neighbor concentration |
| v2-058 | editorial / leyland | Is Leyland Cypress a good privacy tree, and how should it be planted? | `wp:6516`, `wp:10424`, `wp:10577`, `wp:3701` | answerable | hard; product plus editorial/procedure evidence |
| v2-059 | editorial / bamboo-care | What basic care does bamboo need? | `wp:9624`, `wp:9549` | answerable | medium; care versus general overview |
| v2-060 | editorial / citrus-containers | Can citrus trees be grown in containers or pots? | `wp:9557`, `wp:15027` | answerable | medium; near-duplicate instructional posts |
| v2-061 | editorial / citrus-pollination | How can I pollinate citrus trees for a better crop? | `wp:28222` | answerable | easy; specific process |
| v2-062 | editorial / apple-growing | What do I need to know to grow apple trees? | `wp:9559`, `wp:10056` | answerable | medium; care article versus diseases distractor |
| v2-063 | editorial / apple-disease | How can I identify or manage apple tree diseases? | `wp:10056` | answerable | easy; disease-specific target |
| v2-064 | editorial / maple-disease | How can I identify diseases of maple trees? | `wp:10049` | answerable | easy; exact disease target |
| v2-065 | editorial / pine-pests | What insect pests affect pine trees? | `wp:179026`, `wp:37142` | answerable | medium; pest versus disease distinction |
| v2-066 | editorial / neem | Can neem oil control garden pests? | `wp:555358` | answerable | easy; exact remedy term |
| v2-067 | editorial / japanese-maple | What are the main types or groups of Japanese maples? | `wp:755891`, `wp:9947`, `wp:4498` | answerable | hard; taxonomy across broad guides |
| v2-068 | editorial / maple-selection | How do I choose the right red maple for my garden? | `wp:427455` | answerable | medium; selection query |
| v2-069 | product / exact-entity | What are the product details for Thuja Green Giant? | `wp:3699` | answerable | easy; product-fact control |
| v2-070 | product / exact-entity | What are the product details for Bald Cypress Tree? | `wp:3717` | answerable | easy; product-fact control |
| v2-071 | product / exact-entity | What are the product details for Leyland Cypress? | `wp:3701` | answerable | easy; product-fact control with strong editorial distractors |
| v2-072 | negative / out-of-corpus | Do you offer landscape-design consultations at my home? | none expected | unanswerable-in-declared-corpus | hard; must validate negative against FTS and policy pages before freeze |

## Suggested split before adjudication

- Development candidates: `v2-001` through `v2-058` (58 cards; use only after validation).
- Untouched holdout candidates: `v2-059` through `v2-071` (13 answerable cards).
- Negative control candidate: `v2-072`; validate its absence before inclusion.
- Withheld: `v2-020`; never score until policy precedence is explicitly adjudicated.

This split is a proposal, not a randomization. Before freezing, rebalance for
intent, difficulty, product/editorial/support mix, and avoid placing
near-paraphrases of a development query into the holdout.

## Validation and adjudication checklist

1. Confirm every source ID against the frozen TTC snapshot and resolve it to a document revision ID.
2. Inspect each cited document manually; capture exact evidence ranges and set 0/1/2/3 relevance grades.
3. Construct pools using source/SQL, BM25, vector, and hybrid candidates, then judge blind to system.
4. Validate `v2-072` with both FTS and a human search of policy/support content; absence of a lexical hit alone is insufficient.
5. Keep `v2-020` withheld unless a policy owner establishes a source-precedence rule.
