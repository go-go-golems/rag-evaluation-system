import type { HTMLAttributes, ReactNode } from "react";
import styles from "./MarkdownArticle.module.css";

export interface MarkdownArticleProps extends HTMLAttributes<HTMLElement> {
	source: string;
}

type InlineToken = string | ReactNode;

const HAS_SCHEME = /^[a-z][a-z0-9+.-]*:/i;
const SAFE_SCHEME = /^(https?:|mailto:)/i;
const EXTERNAL_SCHEME = /^https?:/i;

// Markdown text is escaped by React, so URLs are the only injection surface.
// Allow http(s)/mailto and relative paths; reject javascript:, data:, and
// protocol-relative URLs so CMS-authored content stays inert.
function sanitizeUrl(value: string): string | null {
	const url = value.trim();
	if (!url || url.startsWith("//")) return null;
	if (HAS_SCHEME.test(url)) {
		return SAFE_SCHEME.test(url) ? url : null;
	}
	return url;
}

function renderInline(value: string, keyPrefix: string): InlineToken[] {
	const parts: InlineToken[] = [];
	let rest = value;
	let index = 0;
	const re = /(\*\*([^*]+)\*\*|`([^`]+)`|\[([^\]]+)\]\(([^)]+)\))/;
	let match: RegExpMatchArray | null;

	while ((match = rest.match(re))) {
		if (match.index && match.index > 0) parts.push(rest.slice(0, match.index));
		if (match[2] != null) {
			parts.push(<strong key={`${keyPrefix}-strong-${index}`}>{match[2]}</strong>);
		} else if (match[3] != null) {
			parts.push(
				<code key={`${keyPrefix}-code-${index}`} className={styles.inlineCode}>
					{match[3]}
				</code>,
			);
		} else if (match[4] != null && match[5] != null) {
			const href = sanitizeUrl(match[5]);
			if (href) {
				parts.push(
					<a
						key={`${keyPrefix}-link-${index}`}
						href={href}
						rel={EXTERNAL_SCHEME.test(href) ? "noopener noreferrer" : undefined}
						className={styles.link}
					>
						{match[4]}
					</a>,
				);
			} else {
				parts.push(match[4]);
			}
		}
		rest = rest.slice((match.index ?? 0) + match[0].length);
		index += 1;
	}
	if (rest) parts.push(rest);
	return parts;
}

function isTableSeparator(cells: string[]) {
	return cells.every((cell) => /^:?-{3,}:?$/.test(cell.trim()));
}

function tableCells(line: string) {
	return line
		.split("|")
		.slice(1, -1)
		.map((cell) => cell.trim());
}

function parseImageMarkdown(value: string) {
	const match = value.trim().match(/^!\[([^\]]*)\]\((\S+?)(?:\s+"([^"]*)")?\)$/);
	if (!match) return null;
	return { alt: match[1] ?? "", src: match[2] ?? "", title: match[3] };
}

export function MarkdownArticle({ source, className, ...rest }: MarkdownArticleProps) {
	const lines = source.replace(/\r\n/g, "\n").split("\n");
	const nodes: ReactNode[] = [];
	let i = 0;
	let key = 0;

	while (i < lines.length) {
		const line = lines[i] ?? "";
		const trimmed = line.trim();

		if (!trimmed) {
			i += 1;
			continue;
		}

		if (trimmed.startsWith("```")) {
			const language = trimmed.slice(3).trim();
			const body: string[] = [];
			i += 1;
			while (i < lines.length && !(lines[i] ?? "").trim().startsWith("```")) {
				body.push(lines[i] ?? "");
				i += 1;
			}
			i += 1;
			nodes.push(
				<pre key={key++} className={styles.codeBlock} data-language={language || undefined}>
					<code>{body.join("\n")}</code>
				</pre>,
			);
			continue;
		}

		if (trimmed === "---" || trimmed === "***") {
			nodes.push(<hr key={key++} className={styles.hr} />);
			i += 1;
			continue;
		}

		const image = parseImageMarkdown(trimmed);
		if (image) {
			const src = sanitizeUrl(image.src);
			nodes.push(
				<figure key={key++} className={styles.figure} data-rag-article-block="image">
					{src && (
						<img
							className={styles.image}
							src={src}
							alt={image.alt}
							title={image.title}
							loading="lazy"
						/>
					)}
					{(image.title || image.alt) && (
						<figcaption className={styles.figcaption}>{image.title || image.alt}</figcaption>
					)}
				</figure>,
			);
			i += 1;
			continue;
		}

		if (trimmed.startsWith("### ")) {
			nodes.push(
				<h3 key={key++} className={styles.h3}>
					{renderInline(trimmed.slice(4), `h3-${key}`)}
				</h3>,
			);
			i += 1;
			continue;
		}
		if (trimmed.startsWith("## ")) {
			nodes.push(
				<h2 key={key++} className={styles.h2}>
					{renderInline(trimmed.slice(3), `h2-${key}`)}
				</h2>,
			);
			i += 1;
			continue;
		}
		if (trimmed.startsWith("# ")) {
			nodes.push(
				<h1 key={key++} className={styles.h1}>
					{renderInline(trimmed.slice(2), `h1-${key}`)}
				</h1>,
			);
			i += 1;
			continue;
		}

		if (trimmed.startsWith("> ")) {
			const quote: string[] = [];
			while (i < lines.length && (lines[i] ?? "").trim().startsWith("> ")) {
				quote.push((lines[i] ?? "").trim().slice(2));
				i += 1;
			}
			nodes.push(
				<blockquote key={key++} className={styles.blockquote}>
					{renderInline(quote.join(" "), `quote-${key}`)}
				</blockquote>,
			);
			continue;
		}

		if (trimmed.startsWith("|")) {
			const rows: string[][] = [];
			while (i < lines.length && (lines[i] ?? "").trim().startsWith("|")) {
				rows.push(tableCells((lines[i] ?? "").trim()));
				i += 1;
			}
			const filtered = rows.filter((row) => !isTableSeparator(row));
			const [head, ...body] = filtered;
			if (head) {
				nodes.push(
					<table key={key++} className={styles.table}>
						<thead>
							<tr>
								{head.map((cell, cellIndex) => (
									<th key={cellIndex}>{renderInline(cell, `th-${key}-${cellIndex}`)}</th>
								))}
							</tr>
						</thead>
						<tbody>
							{body.map((row, rowIndex) => (
								<tr key={rowIndex}>
									{row.map((cell, cellIndex) => (
										<td key={cellIndex}>
											{renderInline(cell, `td-${key}-${rowIndex}-${cellIndex}`)}
										</td>
									))}
								</tr>
							))}
						</tbody>
					</table>,
				);
			}
			continue;
		}

		if (/^\s*-\s+\[[ xX]\]\s+/.test(line)) {
			const items: Array<{ checked: boolean; text: string }> = [];
			while (i < lines.length && /^\s*-\s+\[[ xX]\]\s+/.test(lines[i] ?? "")) {
				const current = lines[i] ?? "";
				items.push({
					checked: /\[[xX]\]/.test(current),
					text: current.replace(/^\s*-\s+\[[ xX]\]\s+/, ""),
				});
				i += 1;
			}
			nodes.push(
				<ul key={key++} className={styles.checkList}>
					{items.map((item, itemIndex) => (
						<li key={itemIndex}>
							<span className={styles.check}>{item.checked ? "☑" : "☐"}</span>
							<span>{renderInline(item.text, `check-${key}-${itemIndex}`)}</span>
						</li>
					))}
				</ul>,
			);
			continue;
		}

		if (/^\s*-\s+/.test(line)) {
			const items: string[] = [];
			while (i < lines.length && /^\s*-\s+/.test(lines[i] ?? "")) {
				items.push((lines[i] ?? "").replace(/^\s*-\s+/, ""));
				i += 1;
			}
			nodes.push(
				<ul key={key++} className={styles.ul}>
					{items.map((item, itemIndex) => (
						<li key={itemIndex}>{renderInline(item, `ul-${key}-${itemIndex}`)}</li>
					))}
				</ul>,
			);
			continue;
		}

		if (/^\s*\d+\.\s+/.test(line)) {
			const items: string[] = [];
			while (i < lines.length && /^\s*\d+\.\s+/.test(lines[i] ?? "")) {
				items.push((lines[i] ?? "").replace(/^\s*\d+\.\s+/, ""));
				i += 1;
			}
			nodes.push(
				<ol key={key++} className={styles.ol}>
					{items.map((item, itemIndex) => (
						<li key={itemIndex}>{renderInline(item, `ol-${key}-${itemIndex}`)}</li>
					))}
				</ol>,
			);
			continue;
		}

		const paragraph: string[] = [];
		while (i < lines.length) {
			const current = lines[i] ?? "";
			const currentTrimmed = current.trim();
			if (
				!currentTrimmed ||
				/^(#{1,3}\s|>|```|---|\*\*\*|\||!\[[^\]]*\]\(|\s*-\s+|\s*\d+\.\s+)/.test(currentTrimmed)
			)
				break;
			paragraph.push(currentTrimmed);
			i += 1;
		}
		nodes.push(
			<p key={key++} className={styles.p}>
				{renderInline(paragraph.join(" "), `p-${key}`)}
			</p>,
		);
	}

	return (
		<article
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="MarkdownArticle"
			{...rest}
		>
			{nodes}
		</article>
	);
}
