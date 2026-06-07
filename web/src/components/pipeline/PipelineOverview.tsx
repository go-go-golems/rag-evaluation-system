import { Caption, CodeText, StatusText } from '@go-go-golems/rag-evaluation-site';
import { DashboardGrid, Panel } from '@go-go-golems/rag-evaluation-site';
import { DataTable, type DataTableColumn } from '@go-go-golems/rag-evaluation-site';
import type { Document, Source } from '../../services/api';

export interface PipelineOverviewProps {
  sources: Source[];
  documents: Document[];
  sourcesLoading?: boolean;
  documentsLoading?: boolean;
}

const sourceColumns: DataTableColumn<Source>[] = [
  { id: 'name', header: 'Name', cell: (source) => source.name },
  { id: 'type', header: 'Type', cell: (source) => <CodeText>{source.type}</CodeText> },
  { id: 'created', header: 'Created', cell: (source) => <CodeText>{source.created_at}</CodeText> },
];

const documentColumns: DataTableColumn<Document>[] = [
  { id: 'title', header: 'Title', cell: (document) => document.title },
  { id: 'status', header: 'Status', cell: (document) => <StatusText status={document.status === 'indexed' ? 'done' : document.status}>{document.status}</StatusText> },
  { id: 'words', header: 'Words', align: 'end', cell: (document) => document.word_count },
];

export function PipelineOverview({ sources, documents, sourcesLoading = false, documentsLoading = false }: PipelineOverviewProps) {
  return (
    <DashboardGrid recipe="twoColumn" data-rag-component="PipelineOverview">
      <Panel title="Sources" density="condensed" fill>
        {sourcesLoading ? (
          <Caption>Loading...</Caption>
        ) : sources.length > 0 ? (
          <DataTable rows={sources} columns={sourceColumns} getRowKey={(source) => source.id} />
        ) : (
          <Caption>No sources yet.</Caption>
        )}
      </Panel>

      <Panel title="Documents" density="condensed" fill>
        {documentsLoading ? (
          <Caption>Loading...</Caption>
        ) : documents.length > 0 ? (
          <DataTable rows={documents} columns={documentColumns} getRowKey={(document) => document.id} />
        ) : (
          <Caption>No documents indexed yet.</Caption>
        )}
      </Panel>
    </DashboardGrid>
  );
}
