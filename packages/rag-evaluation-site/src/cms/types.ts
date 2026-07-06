import type { ArticleBlock } from "../context";

export type CmsContentStatus = "draft" | "published" | "scheduled" | "archived";

export type CmsAssetKind = "image" | "file";

export interface CmsAsset {
	id: string;
	kind: CmsAssetKind;
	title: string;
	filename: string;
	mime: string;
	size: number;
	src: string;
	thumbSrc?: string;
	width?: number;
	height?: number;
	alt?: string;
	tags: string[];
	status: CmsContentStatus;
	createdAt: string;
	updatedAt: string;
}

export interface CmsArticleSummary {
	id: string;
	slug: string;
	title: string;
	status: CmsContentStatus;
	author?: string;
	tags: string[];
	excerpt?: string;
	updatedAt: string;
}

export interface CmsArticleDetail extends CmsArticleSummary {
	blocks: ArticleBlock[];
	coverAssetId?: string;
}

export type UploadItemStatus = "queued" | "uploading" | "done" | "error" | "canceled";

export interface UploadQueueItem {
	id: string;
	filename: string;
	size: number;
	progress: number;
	status: UploadItemStatus;
	error?: string;
}
