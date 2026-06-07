import type { HTMLAttributes } from 'react';
import { contextWindowTokenTotal, getContextKindLabel, type ContextWindowSnapshot } from '../../../context';
import styles from './ContextTreemap.module.css';
export interface ContextTreemapProps extends HTMLAttributes<HTMLDivElement> { snapshot: ContextWindowSnapshot; selectedPartId?: string; }
function formatTokens(tokens: number) { return `${tokens.toLocaleString()} tok`; }
export function ContextTreemap({ snapshot, selectedPartId, className, ...rest }: ContextTreemapProps) {
  const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
  const parts = snapshot.parts.filter((part) => part.kind !== 'empty' && part.tokens > 0);
  const total = contextWindowTokenTotal(snapshot) || 1;
  return <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="ContextTreemap" {...rest}>
    <div className={styles.map} role="img" aria-label={`${snapshot.title} token treemap`}>
      {parts.map((part) => {
        const area = Math.max(7, (part.tokens / total) * 100);
        const selected = effectiveSelectedPartId === part.id;
        return <div key={part.id} className={[styles.tile, styles[`kind_${part.kind}`] ?? styles.kind_other, selected ? styles.selected : ''].filter(Boolean).join(' ')} style={{ flexBasis: `${area}%`, flexGrow: part.tokens }} title={`${part.label}: ${formatTokens(part.tokens)} (${getContextKindLabel(part.kind)})`}>
          <span className={styles.label}>{part.label}</span><span className={styles.tokens}>{formatTokens(part.tokens)}</span>
        </div>;
      })}
    </div>
    <div className={styles.caption}>{formatTokens(total)} in use · {snapshot.title}</div>
  </div>;
}
