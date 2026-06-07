import { Caption } from '@go-go-golems/rag-evaluation-site';
import { Panel } from '@go-go-golems/rag-evaluation-site';
import { DataTable, MetadataGrid, type DataTableColumn } from '@go-go-golems/rag-evaluation-site';
import type { QueueStatus } from '../../../services/api';

export interface QueueHealthPanelProps {
  queues: QueueStatus[];
  isLoading?: boolean;
}

export function QueueHealthPanel({ queues, isLoading = false }: QueueHealthPanelProps) {
  const totalReady = queues.reduce((sum, queue) => sum + queue.ready, 0);
  const totalRunning = queues.reduce((sum, queue) => sum + queue.running, 0);
  const totalFailed = queues.reduce((sum, queue) => sum + queue.failed, 0);

  const columns: DataTableColumn<QueueStatus>[] = [
    { id: 'queue', header: 'Queue', cell: (queue) => queue.queue.replace('rag-eval:', '') },
    { id: 'ready', header: 'R', align: 'end', cell: (queue) => queue.ready },
    { id: 'running', header: 'Run', align: 'end', cell: (queue) => queue.running },
    { id: 'failed', header: 'F', align: 'end', cell: (queue) => <Caption tone={queue.failed > 0 ? 'danger' : 'muted'}>{queue.failed}</Caption> },
  ];

  return (
    <Panel title="Queue Health" density="condensed" data-rag-component="QueueHealthPanel">
      {isLoading ? (
        <Caption>Loading...</Caption>
      ) : queues.length === 0 ? (
        <Caption>No queues active</Caption>
      ) : (
        <>
          <MetadataGrid
            items={[
              { key: 'Ready', value: totalReady },
              { key: 'Running', value: totalRunning },
              { key: 'Failed', value: <Caption tone={totalFailed > 0 ? 'danger' : 'muted'}>{totalFailed}</Caption> },
            ]}
          />
          <DataTable rows={queues} columns={columns} getRowKey={(queue) => queue.queue} />
        </>
      )}
    </Panel>
  );
}
