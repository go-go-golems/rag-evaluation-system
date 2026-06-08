import type { HTMLAttributes, KeyboardEvent } from 'react';
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
  onPartSelect?: (partId: string) => void;
}

function formatTokens(tokens: number) {
  return `${Math.round(tokens).toLocaleString()} tok`;
}

function usedParts(parts: ContextWindowPart[]) {
  return parts.filter((part) => part.kind !== 'empty' && part.tokens > 0);
}

function handlePartKeyDown(event: KeyboardEvent<HTMLDivElement>, partId: string, onPartSelect?: (partId: string) => void) {
  if (!onPartSelect) return;
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    onPartSelect(partId);
  }
}

export function ContextBudgetBar({ snapshot, mode = 'pattern', showLegend = true, selectedPartId, onPartSelect, className, ...rest }: ContextBudgetBarProps) {
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
          const selected = effectiveSelectedPartId === part.id;
          const interactive = Boolean(onPartSelect);
          return (
            <div
              key={part.id}
              className={[styles.segment, styles[`kind_${part.kind}`] ?? styles.kind_other, selected ? styles.selected : ''].filter(Boolean).join(' ')}
              style={{ width: `${width}%` }}
              title={`${part.label}: ${formatTokens(part.tokens)} (${getContextKindLabel(part.kind)})`}
              role={interactive ? 'button' : undefined}
              tabIndex={interactive ? 0 : undefined}
              aria-pressed={interactive ? selected : undefined}
              onClick={interactive ? () => onPartSelect?.(part.id) : undefined}
              onKeyDown={interactive ? (event) => handlePartKeyDown(event, part.id, onPartSelect) : undefined}
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
