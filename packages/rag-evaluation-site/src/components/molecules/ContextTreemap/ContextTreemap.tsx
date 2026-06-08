import type { HTMLAttributes, KeyboardEvent } from 'react';
import { contextWindowTokenTotal, getContextKindLabel, type ContextWindowSnapshot } from '../../../context';
import styles from './ContextTreemap.module.css';

export interface ContextTreemapProps extends HTMLAttributes<HTMLDivElement> {
  snapshot: ContextWindowSnapshot;
  selectedPartId?: string;
  onPartSelect?: (partId: string) => void;
}

function formatTokens(tokens: number) { return `${tokens.toLocaleString()} tok`; }

function handlePartKeyDown(event: KeyboardEvent<HTMLDivElement>, partId: string, onPartSelect?: (partId: string) => void) {
  if (!onPartSelect) return;
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    onPartSelect(partId);
  }
}

export function ContextTreemap({ snapshot, selectedPartId, onPartSelect, className, ...rest }: ContextTreemapProps) {
  const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
  const parts = snapshot.parts.filter((part) => part.kind !== 'empty' && part.tokens > 0);
  const total = contextWindowTokenTotal(snapshot) || 1;
  return <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="ContextTreemap" {...rest}>
    <div className={styles.map} role="img" aria-label={`${snapshot.title} token treemap`}>
      {parts.map((part) => {
        const area = Math.max(7, (part.tokens / total) * 100);
        const selected = effectiveSelectedPartId === part.id;
        const interactive = Boolean(onPartSelect);
        return <div
          key={part.id}
          className={[styles.tile, styles[`kind_${part.kind}`] ?? styles.kind_other, selected ? styles.selected : ''].filter(Boolean).join(' ')}
          style={{ flexBasis: `${area}%`, flexGrow: part.tokens }}
          title={`${part.label}: ${formatTokens(part.tokens)} (${getContextKindLabel(part.kind)})`}
          role={interactive ? 'button' : undefined}
          tabIndex={interactive ? 0 : undefined}
          aria-pressed={interactive ? selected : undefined}
          onClick={interactive ? () => onPartSelect?.(part.id) : undefined}
          onKeyDown={interactive ? (event) => handlePartKeyDown(event, part.id, onPartSelect) : undefined}
        >
          <span className={styles.label}>{part.label}</span><span className={styles.tokens}>{formatTokens(part.tokens)}</span>
        </div>;
      })}
    </div>
    <div className={styles.caption}>{formatTokens(total)} in use · {snapshot.title}</div>
  </div>;
}
