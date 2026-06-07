import type { HTMLAttributes } from 'react';
import { useState } from 'react';
import type { ContextDiagramView, ContextWindowSnapshot } from '../../../context';
import { Button } from '../../atoms';
import { Caption } from '../../foundation';
import { Inline, Panel, Stack } from '../../layout';
import { ContextBudgetBar, ContextLegend, ContextStackDiagram, ContextStripDiagram, ContextTreemap } from '../../molecules';
import styles from './ContextDiagramPanel.module.css';
export interface ContextDiagramPanelProps extends HTMLAttributes<HTMLDivElement> { snapshot: ContextWindowSnapshot; initialView?: ContextDiagramView; selectedPartId?: string; }
const views: ContextDiagramView[] = ['strip','stack','budget','treemap'];
export function ContextDiagramPanel({ snapshot, initialView='strip', selectedPartId, className, ...rest }: ContextDiagramPanelProps) {
  const [view,setView]=useState<ContextDiagramView>(initialView);
  const selected = selectedPartId ?? snapshot.selectedPartId;
  return <Panel title={snapshot.title} actions={<Inline gap="xs">{views.map(v=><Button key={v} size="compact" selected={view===v} aria-pressed={view===v} onClick={()=>setView(v)}>{v}</Button>)}</Inline>} className={className} data-rag-organism="ContextDiagramPanel" {...rest}>
    <Stack gap="sm">
      {snapshot.subtitle && <Caption>{snapshot.subtitle}</Caption>}
      <div className={styles.viewport}>
        {view==='strip' && <ContextStripDiagram snapshot={snapshot} selectedPartId={selected} />}
        {view==='stack' && <ContextStackDiagram snapshot={snapshot} selectedPartId={selected} />}
        {view==='budget' && <ContextBudgetBar snapshot={snapshot} selectedPartId={selected} />}
        {view==='treemap' && <ContextTreemap snapshot={snapshot} selectedPartId={selected} />}
      </div>
      <ContextLegend compact kinds={['system','context','summary','result','generated','evicted','active','empty']} selectedKind={selected ? snapshot.parts.find(p=>p.id===selected)?.kind : undefined} />
    </Stack>
  </Panel>;
}
