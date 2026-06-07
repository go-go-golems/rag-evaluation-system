import { ErrorCallout } from '@go-go-golems/rag-evaluation-site';
import { Caption } from '@go-go-golems/rag-evaluation-site';
import { DataTable, MetadataGrid, type DataTableColumn } from '@go-go-golems/rag-evaluation-site';
import type { OpResult } from '../../../services/api';
import styles from './WorkflowOpResultPanel.module.css';

export interface WorkflowOpResultPanelProps {
  result?: OpResult | null;
  isLoading?: boolean;
  error?: unknown;
}

type RecordRow = OpResult['Records'][number];
type ArtifactRow = OpResult['Artifacts'][number];

function isNotFound(error: unknown) {
  return Boolean(error && typeof error === 'object' && 'status' in error && (error as { status: number }).status === 404);
}

export function WorkflowOpResultPanel({ result, isLoading = false, error }: WorkflowOpResultPanelProps) {
  if (isLoading) return <Caption>Loading result…</Caption>;

  if (error) {
    if (isNotFound(error)) {
      return (
        <section className={styles.root} data-rag-component="WorkflowOpResultPanel">
          <div className={styles.sectionTitle}>Result</div>
          <Caption>No result data recorded for this op.</Caption>
        </section>
      );
    }
    return <ErrorCallout className={styles.error}>Failed to load result: {String(error)}</ErrorCallout>;
  }

  if (!result) return null;

  const recordColumns: DataTableColumn<RecordRow>[] = [
    { id: 'table', header: 'Table', cell: (record) => record.Table },
    { id: 'pk', header: 'PK', cell: (record) => record.PK },
  ];
  const artifactColumns: DataTableColumn<ArtifactRow>[] = [
    { id: 'name', header: 'Name', cell: (artifact) => artifact.Name },
    { id: 'kind', header: 'Kind', cell: (artifact) => artifact.Kind },
    { id: 'contentType', header: 'Content Type', cell: (artifact) => <Caption>{artifact.ContentType}</Caption> },
  ];

  return (
    <section className={styles.root} data-rag-component="WorkflowOpResultPanel">
      <div className={styles.sectionTitle}>Result</div>
      <MetadataGrid items={[{ key: 'Completed', value: result.CompletedAt || '—' }]} />

      {result.Data && Object.keys(result.Data).length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionTitle}>Data</div>
          <pre className={styles.pre}>{JSON.stringify(result.Data, null, 2)}</pre>
        </section>
      )}

      {result.Records && result.Records.length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionTitle}>Records Written ({result.Records.length})</div>
          <DataTable rows={result.Records} columns={recordColumns} getRowKey={(record) => `${record.Table}:${record.PK}`} />
        </section>
      )}

      {result.Artifacts && result.Artifacts.length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionTitle}>Artifacts ({result.Artifacts.length})</div>
          <DataTable rows={result.Artifacts} columns={artifactColumns} getRowKey={(artifact) => `${artifact.Name}:${artifact.Kind}`} />
        </section>
      )}

      {result.EmittedIDs && result.EmittedIDs.length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionTitle}>Emitted Ops ({result.EmittedIDs.length})</div>
          <div className={styles.emittedList}>
            {result.EmittedIDs.slice(0, 20).map(id => <div key={id}>{id}</div>)}
            {result.EmittedIDs.length > 20 && <div>… and {result.EmittedIDs.length - 20} more</div>}
          </div>
        </section>
      )}

      {result.Error && (
        <ErrorCallout>
          [{result.Error.Code}] {result.Error.Message}
          {result.Error.Retryable && <Caption>(retryable)</Caption>}
        </ErrorCallout>
      )}
    </section>
  );
}
