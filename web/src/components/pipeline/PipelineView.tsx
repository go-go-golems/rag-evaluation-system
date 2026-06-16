import React from "react";
import { useListDocumentsQuery, useListSourcesQuery } from "../../services/api";
import { PipelinePage } from "../pages/PipelinePage";

export const PipelineView: React.FC = () => {
	const { data: sources = [], isLoading: sourcesLoading } = useListSourcesQuery();
	const { data: documents = [], isLoading: docsLoading } = useListDocumentsQuery();

	return (
		<PipelinePage
			sources={sources}
			documents={documents}
			sourcesLoading={sourcesLoading}
			documentsLoading={docsLoading}
		/>
	);
};
