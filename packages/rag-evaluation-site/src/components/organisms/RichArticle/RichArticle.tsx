import type { HTMLAttributes } from "react";
import type { ArticleBlock, ContextStyleSet } from "../../../context";
import { contextDefaultStyleSet } from "../../../context";
import { MarkdownArticle } from "../../molecules";
import { ContextDiagramPanel } from "../ContextDiagramPanel";
import styles from "./RichArticle.module.css";

export interface RichArticleProps extends HTMLAttributes<HTMLElement> {
	blocks: ArticleBlock[];
	styleSet?: ContextStyleSet;
}

export function RichArticle({
	blocks,
	styleSet = contextDefaultStyleSet,
	className,
	...rest
}: RichArticleProps) {
	return (
		<article
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="RichArticle"
			{...rest}
		>
			{blocks.map((block) => {
				if (block.kind === "markdown") {
					return (
						<MarkdownArticle
							key={block.id}
							className={styles.block}
							source={block.source}
							data-rag-article-block="markdown"
						/>
					);
				}
				if (block.kind === "context-window") {
					return (
						<div
							key={block.id}
							className={[styles.block, styles.diagramBlock].join(" ")}
							data-rag-article-block="context-window"
						>
							<ContextDiagramPanel
								snapshot={block.snapshot}
								styleSet={styleSet}
								initialView={block.view ?? "budget"}
								views={["budget", "strip", "stack"]}
								chrome="inline"
								showLegend
								showPartDetails={false}
							/>
						</div>
					);
				}
				return (
					<figure
						key={block.id}
						className={[styles.block, styles.imageFigure].join(" ")}
						data-rag-article-block="image"
					>
						<img className={styles.image} src={block.src} alt={block.alt} loading="lazy" />
						{(block.caption || block.alt) && (
							<figcaption className={styles.caption}>{block.caption || block.alt}</figcaption>
						)}
					</figure>
				);
			})}
		</article>
	);
}
