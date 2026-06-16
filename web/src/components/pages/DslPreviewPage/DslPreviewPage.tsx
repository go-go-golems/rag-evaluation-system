import { ErrorCallout } from "@go-go-golems/rag-evaluation-site";
import { Caption } from "@go-go-golems/rag-evaluation-site";
import { Panel } from "@go-go-golems/rag-evaluation-site";
import { useGetDslPageQuery } from "../../../services/api";
import { WidgetRenderer, defaultWidgetRegistry } from "@go-go-golems/rag-evaluation-site";
import styles from "./DslPreviewPage.module.css";

export interface DslPreviewPageProps {
	pageId?: string;
}

export function DslPreviewPage({ pageId = "demo" }: DslPreviewPageProps) {
	const { data, error, isLoading } = useGetDslPageQuery(pageId);

	if (isLoading) {
		return (
			<Panel className={styles.state} title="DSL Preview" density="condensed">
				<Caption>Loading Widget IR…</Caption>
			</Panel>
		);
	}

	if (error) {
		return (
			<ErrorCallout className={styles.state}>Failed to load Widget IR demo page.</ErrorCallout>
		);
	}

	if (!data) {
		return (
			<Panel className={styles.state} title="DSL Preview" density="condensed">
				<Caption>No Widget IR returned.</Caption>
			</Panel>
		);
	}

	return (
		<div className={styles.root} data-rag-page="DslPreviewPage">
			<WidgetRenderer node={data.root} registry={defaultWidgetRegistry} />
		</div>
	);
}
