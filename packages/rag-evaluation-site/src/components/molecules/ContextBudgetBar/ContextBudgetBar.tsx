import type { HTMLAttributes } from 'react';
import {
  contextWindowFillRatio,
  contextWindowTokenTotal,
  getContextKindLabel,
  type ContextDiagramStyle,
  type ContextWindowPart,
  type ContextWindowSnapshot,
} from '../../../context';
import { Caption } from '../../foundation';
import styles from './ContextBudgetBar.module.css';

export interface ContextBudgetBarProps extends HTMLAttributes<HTMLDivElement> {
  snapshot: ContextWindowSnapshot;
  mode?: ContextDiagramStyle;
  showLegend?: boolean;
  selectedPartId?: string;
}

function formatTokens(tokens: number) {
  return `${Math.round(tokens).toLocaleString()} tok`;
}

function usedParts(parts: ContextWindowPart[]) {
  return parts.filter((part) => part.kind !== 'empty' && part.tokens > 0);
}

export function ContextBudgetBar({ snapshot, mode = 'pattern', showLegend = true, selectedPartId, className, ...rest }: ContextBudgetBarProps) {
  const total = contextWindowTokenTotal(snapshot);
  const ratio = contextWindowFillRatio(snapshot);
  const overBudget = total > snapshot.limit;
  const nearBudget = !overBudget && ratio >= 0.8;
  const parts = usedParts(snapshot.parts);
  const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;

  return (
    <div
      className={[styles.root, overBudget ? styles.overBudget : nearBudget ? styles.nearBudget : '', className ?? ''].filter(Boolean).join(' ')}
      data-rag-molecule="ContextBudgetBar"
      data-mode={mode}
      {...rest}
    >
      <div className={styles.header}>
        <Caption transform="uppercase">{snapshot.title}</Caption>
        <Caption tone={overBudget ? 'danger' : nearBudget ? 'warning' : 'muted'}>
          {formatTokens(total)} / {formatTokens(snapshot.limit)} · {Math.round(ratio * 100)}%
        </Caption>
      </div>
      <div className={styles.track} role="img" aria-label={`${snapshot.title}: ${formatTokens(total)} of ${formatTokens(snapshot.limit)} used`}>
        {parts.map((part) => {
          const width = snapshot.limit > 0 ? Math.max(0.5, (part.tokens / snapshot.limit) * 100) : 0;
          return (
            <div
              key={part.id}
              className={[styles.segment, styles[`kind_${part.kind}`] ?? styles.kind_other, effectiveSelectedPartId === part.id ? styles.selected : ''].filter(Boolean).join(' ')}
              style={{ width: `${width}%` }}
              title={`${part.label}: ${formatTokens(part.tokens)} (${getContextKindLabel(part.kind)})`}
            />
          );
        })}
        {overBudget && <div className={styles.limitMarker} title="model context limit" />}
      </div>
      {showLegend && (
        <div className={styles.legend}>
          {parts.map((part) => (
            <span key={part.id} className={styles.legendItem} data-selected={effectiveSelectedPartId === part.id ? 'true' : undefined}>
              <span className={[styles.dot, styles[`kind_${part.kind}`] ?? styles.kind_other].join(' ')} />
              <span>{part.label}</span>
              <span className={styles.tokens}>{formatTokens(part.tokens)}</span>
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
