import {
	cobaltSand,
	createContextStyleSetFromPalette,
	defaultContextStyleSet,
	dustyMagentaBlue,
	signalOrangeCyan,
	slateCoral,
	transcriptStyleSet,
} from "./styles";
import type {
	AnchoredComment,
	ContextCourse,
	ContextHandoutBundle,
	ContextSlide,
	ContextStyleSet,
	ContextWindowPart,
	ContextWindowSnapshot,
	TranscriptFixture,
} from "./types";

export const CONTEXT_WINDOW_LIMIT = 200_000;

export const contextDefaultStyleSet = defaultContextStyleSet(dustyMagentaBlue);
export const contextSignalOrangeStyleSet = defaultContextStyleSet(signalOrangeCyan);
export const contextSlateCoralStyleSet = defaultContextStyleSet(slateCoral);
export const contextCobaltSandStyleSet = defaultContextStyleSet(cobaltSand);

export const transcriptDefaultStyleSet = transcriptStyleSet(dustyMagentaBlue);
export const transcriptSignalOrangeStyleSet = transcriptStyleSet(signalOrangeCyan);
export const transcriptSlateCoralStyleSet = transcriptStyleSet(slateCoral);
export const transcriptCobaltSandStyleSet = transcriptStyleSet(cobaltSand);

export const contextThreeLabelStyleSets = [
	createContextStyleSetFromPalette({
		id: "three-label-dusty",
		name: "Three Label — Dusty Magenta / Blue",
		palette: dustyMagentaBlue,
		entries: [
			{ id: "prompt", label: "Prompt scaffolding", accent: "b", pattern: "checker" },
			{ id: "evidence", label: "Retrieved evidence", accent: "a", pattern: "stipple" },
			{ id: "answer", label: "Answer draft", accent: "a", pattern: "solid", solid: true },
			{ id: "free", label: "Free space", accent: "grid", pattern: "none", hidden: true },
		],
	}),
	createContextStyleSetFromPalette({
		id: "three-label-signal",
		name: "Three Label — Signal Orange / Cyan",
		palette: signalOrangeCyan,
		entries: [
			{ id: "prompt", label: "Prompt scaffolding", accent: "b", pattern: "checker" },
			{ id: "evidence", label: "Retrieved evidence", accent: "a", pattern: "stipple" },
			{ id: "answer", label: "Answer draft", accent: "a", pattern: "solid", solid: true },
			{ id: "free", label: "Free space", accent: "grid", pattern: "none", hidden: true },
		],
	}),
];

export const contextWindowSnapshots: ContextWindowSnapshot[] = [
	{
		id: "t03",
		title: "Turn 3 — warming up",
		subtitle: "Agent has read the repo map and the failing test.",
		limit: CONTEXT_WINDOW_LIMIT,
		parts: [
			{
				id: "t03-system",
				label: "system + tools",
				styleKey: "system",
				tokens: 7200,
				note: "instructions + tool schemas",
			},
			{
				id: "t03-project",
				label: "project context",
				styleKey: "context",
				tokens: 4100,
				note: "CLAUDE.md, repo map",
			},
			{
				id: "t03-conversation",
				label: "conversation",
				styleKey: "context",
				tokens: 5200,
				note: "3 turns so far",
			},
			{
				id: "t03-file-reads",
				label: "file reads",
				styleKey: "result",
				tokens: 9800,
				note: "test + 2 source files",
			},
			{
				id: "t03-task",
				label: "current task",
				styleKey: "active",
				tokens: 1400,
				note: "fix failing test",
			},
			{ id: "t03-empty", label: "free space", styleKey: "empty", tokens: 162300 },
		],
	},
	{
		id: "t14",
		title: "Turn 14 — deep in the bug",
		subtitle: "Many tool results accumulated; scratchpad reasoning growing.",
		limit: CONTEXT_WINDOW_LIMIT,
		selectedPartId: "t14-file-reads",
		parts: [
			{
				id: "t14-system",
				label: "system + tools",
				styleKey: "system",
				tokens: 7200,
				note: "instructions + tool schemas",
			},
			{
				id: "t14-project",
				label: "project context",
				styleKey: "context",
				tokens: 4100,
				note: "CLAUDE.md, repo map",
			},
			{
				id: "t14-summary",
				label: "summary",
				styleKey: "summary",
				tokens: 3400,
				note: "turns 1–6 compressed",
			},
			{
				id: "t14-conversation",
				label: "conversation",
				styleKey: "context",
				tokens: 21800,
				note: "turns 7–14",
			},
			{
				id: "t14-file-reads",
				label: "file reads",
				styleKey: "result",
				tokens: 38600,
				note: "6 files + grep output",
			},
			{
				id: "t14-test-output",
				label: "test output",
				styleKey: "result",
				tokens: 12400,
				note: "stack traces",
			},
			{
				id: "t14-scratchpad",
				label: "scratchpad",
				styleKey: "generated",
				tokens: 9200,
				note: "hypotheses, plan",
			},
			{
				id: "t14-task",
				label: "current task",
				styleKey: "active",
				tokens: 1900,
				note: "patch the parser",
			},
			{ id: "t14-empty", label: "free space", styleKey: "empty", tokens: 91700 },
		],
	},
	{
		id: "t31",
		title: "Turn 31 — at the limit",
		subtitle: "Window full. Old turns evicted; results re-summarized to reclaim space.",
		limit: CONTEXT_WINDOW_LIMIT,
		parts: [
			{
				id: "t31-system",
				label: "system + tools",
				styleKey: "system",
				tokens: 7200,
				note: "instructions + tool schemas",
			},
			{
				id: "t31-project",
				label: "project context",
				styleKey: "context",
				tokens: 4100,
				note: "CLAUDE.md, repo map",
			},
			{
				id: "t31-evicted",
				label: "evicted",
				styleKey: "evicted",
				tokens: 28400,
				note: "turns 1–18 dropped",
			},
			{
				id: "t31-summary",
				label: "summary",
				styleKey: "summary",
				tokens: 8600,
				note: "rolling summary",
			},
			{
				id: "t31-conversation",
				label: "conversation",
				styleKey: "context",
				tokens: 41200,
				note: "turns 19–31",
			},
			{
				id: "t31-file-reads",
				label: "file reads",
				styleKey: "result",
				tokens: 64800,
				note: "12 files in working set",
			},
			{
				id: "t31-scratchpad",
				label: "scratchpad",
				styleKey: "generated",
				tokens: 22600,
				note: "long reasoning trace",
			},
			{
				id: "t31-task",
				label: "current task",
				styleKey: "active",
				tokens: 2100,
				note: "verify the fix",
			},
			{ id: "t31-empty", label: "free space", styleKey: "empty", tokens: 21000 },
		],
	},
	{
		id: "over",
		title: "Turn 31 — before reclaim (over budget)",
		subtitle: "What the naïve append-everything window would have looked like.",
		limit: CONTEXT_WINDOW_LIMIT,
		parts: [
			{ id: "over-system", label: "system + tools", styleKey: "system", tokens: 7200 },
			{ id: "over-project", label: "project context", styleKey: "context", tokens: 4100 },
			{ id: "over-conversation", label: "conversation (all)", styleKey: "context", tokens: 88200 },
			{ id: "over-file-reads", label: "file reads (all)", styleKey: "result", tokens: 96400 },
			{ id: "over-scratchpad", label: "scratchpad", styleKey: "generated", tokens: 31600 },
			{ id: "over-task", label: "current task", styleKey: "active", tokens: 2100 },
		],
	},
];

export const transcriptFixture: TranscriptFixture = {
	id: "parser-session",
	title: "Coding-agent transcript — parser fix",
	subtitle: "Annotated transcript fixture from the context-viewer prototype.",
	selectedAnnotationId: "task-framing",
	messages: [
		{
			id: "m1",
			role: "user",
			tokens: 320,
			text: "The test `parser.test.ts › handles nested arrays` is failing on main. Can you fix it?",
		},
		{
			id: "m2",
			role: "assistant",
			tokens: 180,
			text: "I'll start by reading the failing test and the parser it exercises.",
			annotationIds: ["task-framing"],
		},
		{
			id: "m3",
			role: "tool",
			name: "read_file",
			tokens: 2400,
			text: "parser.test.ts  (148 lines)\nsrc/parser.ts  (412 lines)",
			annotationIds: ["tool-result-cost"],
		},
		{
			id: "m4",
			role: "assistant",
			tokens: 540,
			text: "The test expects `parseArray` to recurse into nested brackets, but the current implementation only handles one level. Let me run the test to confirm the failure mode.",
		},
		{
			id: "m5",
			role: "tool",
			name: "run_test",
			tokens: 1900,
			text: 'FAIL  parser.test.ts\n  ● handles nested arrays\n    Expected: [[1,2],[3]]\n    Received: ["[1","2]","[3]"]\n    at parseArray (src/parser.ts:218)',
			annotationIds: ["stack-trace-noise"],
		},
		{
			id: "m6",
			role: "assistant",
			tokens: 880,
			text: "Confirmed — `parseArray` splits on commas without tracking bracket depth. I'll rewrite it to maintain a depth counter and recurse.",
			annotationIds: ["scratchpad-reasoning"],
		},
		{
			id: "m7",
			role: "tool",
			name: "edit_file",
			tokens: 1200,
			text: "src/parser.ts  — rewrote parseArray (lines 198–240)",
		},
		{
			id: "m8",
			role: "tool",
			name: "run_test",
			tokens: 760,
			text: "PASS  parser.test.ts  (12 tests)",
		},
		{
			id: "m9",
			role: "assistant",
			tokens: 240,
			text: "All 12 tests pass. The fix tracks bracket depth and recurses on nested arrays. Want me to add a regression test for deeply-nested input?",
			annotationIds: ["window-rebuilt"],
		},
		{ id: "m10", role: "user", tokens: 60, text: "Yes please." },
	],
	annotations: [
		{
			id: "task-framing",
			targetMessageId: "m2",
			styleKey: "context",
			label: "Task framing",
			text: "The current task is small but high-value — it stays pinned at the bottom of every window for the rest of the session.",
		},
		{
			id: "tool-result-cost",
			targetMessageId: "m3",
			styleKey: "result",
			label: "Tool result cost",
			text: "Two file reads are cheap now, but reads compound — by turn 14 the working set is 38k.",
		},
		{
			id: "stack-trace-noise",
			targetMessageId: "m5",
			styleKey: "result",
			label: "Stack traces are noisy",
			text: "Test output carries stack traces. Keep the assertion and first frame; the rest is a summarization candidate.",
		},
		{
			id: "scratchpad-reasoning",
			targetMessageId: "m6",
			styleKey: "generated",
			label: "Scratchpad reasoning",
			text: "Reasoning tokens are generated, not retrieved. Useful in-turn, but rarely worth carrying forward verbatim.",
		},
		{
			id: "window-rebuilt",
			targetMessageId: "m9",
			styleKey: "active",
			label: "Window rebuilt here",
			text: "Before this reply, the window was rebuilt: the failing run_test output was summarized to one line now that it's resolved.",
		},
	],
};

export const contextCourseFixture: ContextCourse = {
	id: "context-window-engineering",
	kicker: "LIVE WORKSHOP · SESSION 04",
	title: "Context Window Engineering",
	tagline:
		"What fits, what falls out, and why. A hands-on session on designing the prompt that the model actually sees.",
	when: "Thursday, June 18, 2026 · 6:30 – 8:30 PM",
	where: "South Park Commons · San Francisco",
	format: "In person · 24 seats · laptops required",
	price: "Free for members",
	instructor: {
		name: "M. Calder",
		role: "Applied ML, Developer Tools",
		bio: "Builds agent runtimes; has spent more nights staring at token budgets than is healthy.",
	},
	blurb:
		"Every model call is a window of fixed size. The difference between an agent that finishes the job and one that loops forever is usually not the model — it's what you packed into that window.",
	outcomes: [
		"Read any context window as a budget of competing claims on space",
		"Diagram a window four ways — strip, stack, budget, and treemap",
		"Decide what to summarize, what to evict, and what to pin",
		"Instrument a coding agent and watch its window churn turn by turn",
	],
	agenda: [
		{
			id: "agenda-01",
			number: "01",
			title: "The window as a budget",
			description: "Fixed size, competing tenants. Why everything is a trade.",
			duration: "20 min",
		},
		{
			id: "agenda-02",
			number: "02",
			title: "Anatomy of a call",
			description: "System, context, results, scratchpad, the current task.",
			duration: "25 min",
		},
		{
			id: "agenda-03",
			number: "03",
			title: "Growth & truncation",
			description: "What falls out when the window fills, and how to choose.",
			duration: "20 min",
		},
		{
			id: "agenda-04",
			number: "04",
			title: "The agent loop",
			description: "Rebuild, don't append. Perceive → think → act → observe.",
			duration: "25 min",
		},
		{
			id: "agenda-05",
			number: "05",
			title: "Reclaiming space",
			description: "Summarize, evict, retrieve, pin. Live on a real session.",
			duration: "25 min",
		},
		{
			id: "agenda-06",
			number: "06",
			title: "Workshop",
			description: "Instrument your own agent and annotate its window.",
			duration: "30 min",
		},
	],
};

export const contextSlides: ContextSlide[] = [
	{
		id: "slide-01",
		number: "01",
		kicker: "MODULE 01",
		title: "The window is a fixed budget",
		view: "budget",
		snapshotId: "t14",
		notes: [
			"Every model call gets a window of fixed size.",
			"Instructions, history, results, reasoning and the task all compete for the same space.",
			"If it doesn't fit, the model can't see it. Headroom is the whole game.",
		],
	},
	{
		id: "slide-02",
		number: "02",
		kicker: "MODULE 02",
		title: "Anatomy of a single call",
		view: "stack",
		snapshotId: "t14",
		notes: [
			"A call is a layered document: durable layers on top, volatile ones below.",
			"System and project context are stable. Results and scratchpad churn every turn.",
			"The current task sits at the bottom and is never evicted.",
		],
	},
	{
		id: "slide-03",
		number: "03",
		kicker: "MODULE 03",
		title: "Composition at a glance",
		view: "strip",
		snapshotId: "t14",
		notes: [
			"Laid flat, the window reads left to right as a strip of segments.",
			"Separators mark where one tenant ends and the next begins.",
			"This is the fastest way to sanity-check what you packed.",
		],
	},
	{
		id: "slide-04",
		number: "04",
		kicker: "MODULE 04",
		title: "Where the tokens actually go",
		view: "treemap",
		snapshotId: "t31",
		notes: [
			"Area equals tokens. The biggest box is usually not the one you think.",
			"Tool results and file reads dominate long sessions — not the conversation.",
			"Treemaps make the expensive tenants impossible to miss.",
		],
	},
	{
		id: "slide-05",
		number: "05",
		kicker: "MODULE 05",
		title: "Growth & truncation",
		view: "budget",
		snapshotId: "t31",
		notes: [
			"As the session runs, new information is added every turn.",
			"When the limit is reached, the oldest, lowest-value content is dropped.",
			"Evicted turns are gone unless they were summarized first.",
		],
	},
	{
		id: "slide-06",
		number: "06",
		kicker: "MODULE 06",
		title: "Reclaim, don't append",
		view: "budget",
		snapshotId: "over",
		notes: [
			"Append-everything blows the budget — here the same session runs 26% over.",
			"Summarize resolved results, evict stale turns, retrieve on demand.",
			"Rebuild the window every turn around what matters now.",
		],
	},
];

export const anchoredCommentFixtures: AnchoredComment[] = [
	{
		id: "comment-file-reads",
		anchorX: 0.3,
		anchorY: 0.52,
		author: "you",
		text: "This file-reads block is 38k tokens — biggest single tenant. Candidate to summarize once the bug is found.",
		time: "2m",
		status: "open",
	},
	{
		id: "comment-scratchpad",
		anchorX: 0.74,
		anchorY: 0.5,
		author: "Priya",
		text: "Why is the scratchpad carried forward? Drop it after the step completes.",
		time: "just now",
		status: "open",
	},
	{
		id: "comment-summary",
		anchorX: 0.18,
		anchorY: 0.3,
		author: "Manuel",
		text: "Resolved tool output can collapse into the rolling summary after the assertion is understood.",
		time: "resolved",
		status: "resolved",
	},
];

export const contextHandoutFixture: ContextHandoutBundle = {
	intro:
		"Everything from tonight, to take home. Slides, the diagram source, and a one-page field guide.",
	docs: [
		{
			id: "fieldguide",
			title: "The Context Window — Field Guide",
			file: "context-window-field-guide.md",
			format: "Markdown",
			size: "11 KB",
			description:
				"One-page reference: the six ideas, the four diagrams, and the reclaim checklist.",
			body: `# The Context Window — Field Guide

Every model call is a **fixed-size document**. The model can only reason over the instructions, history, retrieved context, tool results, summaries, and task framing that fit in that document.

## The one idea

Treat the context window as a workbench, not a filing cabinet. Everything on the workbench must earn its place right now.

- **Pin the current task** — requirements and acceptance criteria stay visible.
- **Keep exact evidence only while it is active** — resolved tool output becomes summary plus source pointer.
- **Rebuild every turn** — choose the next window deliberately instead of appending forever.

## Budget table

| Tenant | Keep verbatim when | Reclaim strategy |
| --- | --- | --- |
| System instructions | Always relevant | Keep pinned |
| Tool output | It contains an unresolved error | Summarize after resolution |
| File reads | The next patch depends on exact lines | Keep excerpts with paths |
| Conversation | It changes requirements | Summarize decisions |

## The four views

1. Strip — read the whole window left to right.
2. Stack — compare durable and volatile layers.
3. Budget — check headroom and overflow.
4. Treemap — find the biggest token tenants.

> If it is not in the prompt, the model cannot use it. Saved files matter only after retrieval.

### Minimal instrumentation

Track \`tokens_by_kind\` for system, context, conversation, results, generated text, active task, evicted text, and empty space.`,
		},
		{
			id: "slides-pdf",
			title: "Slide Deck",
			file: "context-window-engineering.pdf",
			format: "PDF",
			size: "1.4 MB",
			description: "Printable deck for the six workshop modules.",
			body: "PDF preview is not available in this static Storybook fixture.",
		},
		{
			id: "diagram-json",
			title: "Diagram source (JSON)",
			file: "context-window-diagram.json",
			format: "JSON",
			size: "6 KB",
			description: "Token segment data used by the visualizer examples.",
			body: '{ "limit": 200000, "parts": [] }',
		},
		{
			id: "checklist",
			title: "Reclaim Checklist",
			file: "reclaim-checklist.md",
			format: "Markdown",
			size: "3 KB",
			description: "The order to cut things when the window fills.",
			body: `# Reclaim Checklist

- [ ] Is the current task pinned?
- [ ] Have resolved tool results been summarized?
- [ ] Is the rolling summary still accurate?
- [ ] Are stale file reads replaced by path + line references?
- [ ] Did the next action survive compaction?`,
		},
	],
};

export const DEFAULT_CONTEXT_HEADROOM_STYLE_KEYS = ["empty", "free", "headroom"] as const;

export function contextWindowHeadroomStyleKeys(styleSet?: ContextStyleSet) {
	const keys = new Set<string>(DEFAULT_CONTEXT_HEADROOM_STYLE_KEYS);
	for (const item of styleSet?.legend ?? []) {
		if (item.hidden) {
			keys.add(item.styleKey ?? item.id);
		}
	}
	return keys;
}

export function contextWindowIsHeadroomPart(
	part: ContextWindowPart,
	options: { includeEmpty?: boolean; headroomStyleKeys?: Iterable<string> } = {},
) {
	if (options.includeEmpty) {
		return false;
	}
	return new Set(options.headroomStyleKeys ?? DEFAULT_CONTEXT_HEADROOM_STYLE_KEYS).has(
		part.styleKey,
	);
}

export function contextWindowTokenTotal(
	snapshot: ContextWindowSnapshot,
	options: { includeEmpty?: boolean; headroomStyleKeys?: Iterable<string> } = {},
) {
	return snapshot.parts.reduce((total, part) => {
		if (contextWindowIsHeadroomPart(part, options)) {
			return total;
		}
		return total + part.tokens;
	}, 0);
}

export function contextWindowFillRatio(
	snapshot: ContextWindowSnapshot,
	options: { includeEmpty?: boolean; headroomStyleKeys?: Iterable<string> } = {},
) {
	if (snapshot.limit <= 0) {
		return 0;
	}
	return contextWindowTokenTotal(snapshot, options) / snapshot.limit;
}
