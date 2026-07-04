import type { CSSProperties, HTMLAttributes } from "react";
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
				if (block.kind === "gallery") {
					return (
						<figure
							key={block.id}
							className={[styles.block, styles.galleryFigure].join(" ")}
							data-rag-article-block="gallery"
						>
							<div
								className={styles.galleryGrid}
								style={{ "--rag-gallery-columns": block.columns ?? 3 } as CSSProperties}
							>
								{block.images.map((image, index) => (
									<figure key={`${block.id}-${index}`} className={styles.galleryItem}>
										<img className={styles.image} src={image.src} alt={image.alt} loading="lazy" />
										{(image.caption || image.alt) && (
											<figcaption className={styles.caption}>
												{image.caption || image.alt}
											</figcaption>
										)}
									</figure>
								))}
							</div>
							{block.caption && <figcaption className={styles.caption}>{block.caption}</figcaption>}
						</figure>
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
