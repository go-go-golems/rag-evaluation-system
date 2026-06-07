import type { HTMLAttributes } from 'react';
import { contextKindOrder, getContextKindLabel, type ContextDiagramStyle, type ContextPartKind } from '../../../context';
import { ContextKindSwatch } from '../../atoms';
import styles from './ContextLegend.module.css';

export interface ContextLegendProps extends HTMLAttributes<HTMLDivElement> {
  kinds?: ContextPartKind[];
  mode?: ContextDiagramStyle;
  compact?: boolean;
  selectedKind?: ContextPartKind;
}

export function ContextLegend({
  kinds = contextKindOrder.filter((kind) => kind !== 'other'),
  mode = 'pattern',
  compact = false,
  selectedKind,
  className,
  ...rest
}: ContextLegendProps) {
  return (
    <div
      className={[styles.root, compact ? styles.compact : styles.default, className ?? ''].filter(Boolean).join(' ')}
      data-rag-molecule="ContextLegend"
      data-mode={mode}
      {...rest}
    >
      {kinds.map((kind) => (
        <div key={kind} className={styles.item} data-selected={selectedKind === kind ? 'true' : undefined}>
          <ContextKindSwatch kind={kind} mode={mode} size={compact ? 'sm' : 'md'} selected={selectedKind === kind} />
          <span className={styles.label}>{getContextKindLabel(kind)}</span>
        </div>
      ))}
    </div>
  );
}
