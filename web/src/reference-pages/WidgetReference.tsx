/**
 * WidgetReference.tsx
 *
 * Renders every widget in every key state, plus composition patterns.
 * This is the "visual language" page for external designers.
 *
 * Rendered server-side → bundled into a self-contained HTML file.
 */

import { Panel } from '../components/layout';
import { Stack } from '../components/layout';
import { Inline } from '../components/layout';
import { DashboardGrid } from '../components/layout';
import { FormRow } from '../components/layout';
import { TabList } from '../components/layout';
import { AppShell } from '../components/layout';
import { AppNav } from '../components/molecules';
import { DataTable } from '../components/molecules';
import { MetadataGrid } from '../components/molecules';
import { StatusText } from '../components/foundation';
import { Caption } from '../components/foundation';
import { Divider } from '../components/foundation';
import { Text } from '../components/foundation';
import { Button } from '../components/atoms';
import { TextInput } from '../components/atoms';
import { SelectInput } from '../components/atoms';
import { CheckboxRow } from '../components/atoms';
import { ErrorCallout } from '../components/atoms';

/* ── sample data ──────────────────────────────────────────── */

const sampleRows = [
  { id: '1', name: 'retrieval-a', status: 'succeeded', latency: '243ms', docs: 47 },
  { id: '2', name: 'retrieval-b', status: 'running',   latency: '—',     docs: 0 },
  { id: '3', name: 'retrieval-c', status: 'failed',    latency: '1201ms', docs: 0 },
  { id: '4', name: 'retrieval-d', status: 'pending',   latency: '—',     docs: 0 },
];

const sampleColumns = [
  { id: 'id',   header: 'ID',   cell: (r: typeof sampleRows[0]) => r.id },
  { id: 'name', header: 'Name', cell: (r: typeof sampleRows[0]) => r.name },
  { id: 'status', header: 'Status', cell: (r: typeof sampleRows[0]) => <StatusText status={r.status} icon>{r.status}</StatusText> },
  { id: 'latency', header: 'Latency', cell: (r: typeof sampleRows[0]) => r.latency },
];

const metaItems = [
  { key: 'Corpus', value: 'wikipedia-en' },
  { key: 'Retriever', value: 'bm25' },
  { key: 'Top-K', value: '10' },
  { key: 'Embedding', value: 'text-embedding-3-small' },
  { key: 'Created', value: '2026-05-27 14:32' },
  { key: 'Status', value: <StatusText status="succeeded" icon>Completed</StatusText> },
];

const navItems = [
  { id: 'dashboard', label: 'Dashboard' },
  { id: 'corpus', label: 'Corpus' },
  { id: 'queries', label: 'Queries' },
];

/* ── component ────────────────────────────────────────────── */

export default function WidgetReference() {
  return (
    <>
      <h1 style={{ borderBottom: '3px solid #000', paddingBottom: 8, marginBottom: '0.5em' }}>
        Widgets &amp; Composition
      </h1>
      <p style={{ color: '#666', marginBottom: '2rem' }}>
        Visual reference for every widget in the design system, rendered from
        the real React components. Shows all key states and how widgets compose.
      </p>

      {/* ── Panel ────────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Panel</h2>
        <p className="ds-ref-label">The primary container. Bordered card with optional header and actions.</p>

        <h3 style={{ margin: '1em 0 0.5em' }}>Panel with title</h3>
        <Panel title="Server Health">
          <StatusText status="succeeded" icon>All systems operational</StatusText>
        </Panel>

        <h3 style={{ margin: '1em 0 0.5em' }}>Panel with title + actions</h3>
        <Panel title="Query Queue" actions={<Button variant="primary" size="compact">Add query</Button>}>
          <StatusText status="running" icon>4 queries in queue</StatusText>
        </Panel>

        <h3 style={{ margin: '1em 0 0.5em' }}>Panel condensed</h3>
        <Panel title="Metrics" density="condensed">
          <Caption tone="muted">Condensed padding for dense dashboard tiles</Caption>
        </Panel>

        <h3 style={{ margin: '1em 0 0.5em' }}>Panel without title</h3>
        <Panel>
          <Text>A bare container — just a border and padding.</Text>
        </Panel>
      </div>

      {/* ── Stack ────────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Stack</h2>
        <p className="ds-ref-label">Vertical flex layout. Gap values: xs(2px), sm(4px), md(8px), lg(16px), xl(24px).</p>
        <div className="ds-ref-row">
          {(['xs', 'sm', 'md', 'lg', 'xl'] as const).map(gap => (
            <div key={gap} className="ds-ref-item">
              <div className="ds-ref-label">gap={gap}</div>
              <Stack gap={gap}>
                <div style={{ background: '#eee', padding: 4 }}>A</div>
                <div style={{ background: '#eee', padding: 4 }}>B</div>
                <div style={{ background: '#eee', padding: 4 }}>C</div>
              </Stack>
            </div>
          ))}
        </div>
      </div>

      {/* ── Inline ───────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Inline</h2>
        <p className="ds-ref-label">Horizontal flex. Gap values: xs, sm, md, lg. Justify: start, between, end.</p>
        <div className="ds-ref-label">justify=start (default):</div>
        <Inline gap="sm" justify="start">
          <Button size="compact">First</Button>
          <Button size="compact">Second</Button>
          <Button size="compact">Third</Button>
        </Inline>
        <div className="ds-ref-label" style={{ marginTop: '0.5em' }}>justify=between:</div>
        <Inline gap="sm" justify="between">
          <Button size="compact">Left</Button>
          <Button variant="primary" size="compact">Right</Button>
        </Inline>
        <div className="ds-ref-label" style={{ marginTop: '0.5em' }}>justify=end:</div>
        <Inline gap="sm" justify="end">
          <Caption tone="muted">Cancel</Caption>
          <Button variant="primary" size="compact">Confirm</Button>
        </Inline>
      </div>

      {/* ── DashboardGrid ────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">DashboardGrid</h2>
        <p className="ds-ref-label">Responsive CSS grid. Recipe: twoColumn (default), searchWorkbench, corpusExplorer.</p>
        <DashboardGrid recipe="twoColumn">
          <Panel title="Metrics" density="condensed">
            <StatusText status="succeeded" icon>142 queries</StatusText>
          </Panel>
          <Panel title="Health" density="condensed">
            <StatusText status="ok" icon>99.8% uptime</StatusText>
          </Panel>
          <Panel title="Queue" density="condensed">
            <StatusText status="running" icon>3 running</StatusText>
          </Panel>
          <Panel title="Errors" density="condensed">
            <StatusText status="error" icon>2 failures</StatusText>
          </Panel>
        </DashboardGrid>
      </div>

      {/* ── DataTable ────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">DataTable</h2>
        <p className="ds-ref-label">Sortable, selectable data table with StatusText columns.</p>
        <Panel title="Query Results">
          <DataTable
            columns={sampleColumns}
            rows={sampleRows}
            getRowKey={(r) => r.id}
          />
        </Panel>

        <h3 style={{ margin: '1em 0 0.5em' }}>Empty state</h3>
        <Panel title="Empty Table">
          <DataTable
            columns={sampleColumns}
            rows={[]}
            getRowKey={(r) => r.id}
            emptyMessage="No queries yet. Click 'Add query' to start."
          />
        </Panel>
      </div>

      {/* ── MetadataGrid ──────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">MetadataGrid</h2>
        <div className="ds-ref-row">
          <div style={{ flex: 1 }}>
            <div className="ds-ref-label">density=normal (default):</div>
            <MetadataGrid items={metaItems} />
          </div>
          <div style={{ flex: 1 }}>
            <div className="ds-ref-label">density=compact:</div>
            <MetadataGrid items={metaItems} density="compact" />
          </div>
        </div>
      </div>

      {/* ── FormRow + inputs ──────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">FormRow + Inputs</h2>
        <Panel title="Settings">
          <Stack gap="md">
            <FormRow label="Name" control={<TextInput defaultValue="" placeholder="Enter name" />} />
            <FormRow label="Retriever" control={
              <SelectInput defaultValue="bm25">
                <option value="bm25">BM25</option>
                <option value="dense">Dense</option>
                <option value="hybrid">Hybrid</option>
              </SelectInput>
            } />
            <FormRow label="Verbose" control={<CheckboxRow checked={false} onChange={() => {}}>Enable verbose logging</CheckboxRow>} />
            <FormRow label="Max Results" control={<TextInput defaultValue="10" type="number" min={1} max={100} />} />
            <Divider />
            <Inline gap="sm" justify="end">
              <Button>Cancel</Button>
              <Button variant="primary">Save</Button>
            </Inline>
          </Stack>
        </Panel>
      </div>

      {/* ── TabList ───────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">TabList</h2>
        <TabList
          items={[
            { id: 'overview', label: 'Overview' },
            { id: 'artifacts', label: 'Artifacts' },
            { id: 'chunks', label: 'Chunks' },
          ]}
          activeId="overview"
          onChange={() => {}}
        />
      </div>

      {/* ── ErrorCallout ──────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">ErrorCallout</h2>
        <ErrorCallout>Something went wrong: the retrieval service returned HTTP 503.</ErrorCallout>
      </div>

      {/* ── Composition: Dashboard ────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Composition: Dashboard</h2>
        <p className="ds-ref-label">AppShell → AppNav → DashboardGrid → Panel → Stack → StatusText</p>
        <div style={{ border: '1px solid #ccc', overflow: 'hidden' }}>
          <AppShell header={
            <AppNav brand="RAG Eval" items={navItems} activeItemId="dashboard" onItemSelect={() => {}} />
          }>
            <DashboardGrid recipe="twoColumn">
              <Panel title="Queue Health" density="condensed">
                <StatusText status="succeeded" icon>All clear</StatusText>
              </Panel>
              <Panel title="Query Queue" density="condensed">
                <StatusText status="running" icon>3 running</StatusText>
              </Panel>
              <Panel title="Recent Queries">
                <DataTable
                  columns={sampleColumns.slice(0, 3)}
                  rows={sampleRows}
                  getRowKey={(r) => r.id}
                />
              </Panel>
              <Panel title="Configuration">
                <MetadataGrid items={metaItems.slice(0, 4)} density="compact" />
              </Panel>
            </DashboardGrid>
          </AppShell>
        </div>
      </div>

      {/* ── Composition: Master-Detail ────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Composition: Master-Detail</h2>
        <p className="ds-ref-label">DashboardGrid → DataTable (master) + Panel/MetadataGrid (detail)</p>
        <DashboardGrid recipe="twoColumn">
          <Panel title="Queries">
            <DataTable
              columns={sampleColumns.slice(0, 2)}
              rows={sampleRows}
              getRowKey={(r) => r.id}
            />
          </Panel>
          <Panel title="Query Detail">
            <Stack gap="sm">
              <MetadataGrid items={metaItems} density="compact" />
              <Divider />
              <Inline gap="sm">
                <Button size="compact">Re-run</Button>
                <Button size="compact" variant="primary">Edit</Button>
              </Inline>
            </Stack>
          </Panel>
        </DashboardGrid>
      </div>

      {/* ── Composition: Settings Form ───────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">Composition: Settings Form</h2>
        <p className="ds-ref-label">Panel → Stack → FormRow × N + Inline → Button</p>
        <Panel title="Retrieval Configuration">
          <Stack gap="sm">
            <FormRow label="Corpus" control={<TextInput defaultValue="wikipedia-en" />} />
            <FormRow label="Retriever" control={
              <SelectInput defaultValue="bm25">
                <option value="bm25">BM25</option>
                <option value="dense">Dense</option>
              </SelectInput>
            } />
            <FormRow label="Top-K" control={<TextInput defaultValue="10" type="number" />} />
            <FormRow label="Verbose" control={<CheckboxRow checked={false} onChange={() => {}}>Log retrieval scores</CheckboxRow>} />
            <Divider />
            <Inline gap="sm" justify="end">
              <Button>Reset</Button>
              <Button variant="primary">Apply</Button>
            </Inline>
          </Stack>
        </Panel>
      </div>

    </>
  );
}

export { WidgetReference };
