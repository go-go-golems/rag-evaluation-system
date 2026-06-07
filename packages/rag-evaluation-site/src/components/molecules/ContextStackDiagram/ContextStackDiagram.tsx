import type { HTMLAttributes } from 'react';
import { getContextKindLabel, type ContextWindowSnapshot } from '../../../context';
import styles from './ContextStackDiagram.module.css';

export interface ContextStackDiagramProps extends HTMLAttributes<HTMLDivElement> {
  snapshot: ContextWindowSnapshot;
  selectedPartId?: string;
}

function formatTokens(tokens: number) { return `${tokens.toLocaleString()} tok`; }

export function ContextStackDiagram({ snapshot, selectedPartId, className, ...rest }: ContextStackDiagramProps) {
  const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
  const visible = snapshot.parts.filter((part) => part.tokens > 0);
  const maxTokens = Math.max(...visible.map((part) => part.tokens), 1);
  return (
    <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="ContextStackDiagram" {...rest}>
      <div className={styles.layers} role="img" aria-label={`${snapshot.title} context stack`}>
        {visible.map((part) => {
          const selected = effectiveSelectedPartId === part.id;
          const height = Math.max(26, Math.min(74, 22 + (part.tokens / maxTokens) * 52));
          return (
            <div key={part.id} className={[styles.layer, styles[`kind_${part.kind}`] ?? styles.kind_other, selected ? styles.selected : ''].filter(Boolean).join(' ')} style={{ minHeight: height }} title={`${part.label}: ${formatTokens(part.tokens)} (${getContextKindLabel(part.kind)})`}>
              <span className={styles.label}>{part.label}</span>
              <span className={styles.tokens}>{formatTokens(part.tokens)}</span>
            </div>
          );
        })}
      </div>
      <div className={styles.caption}>{snapshot.title}</div>
    </div>
  );
}
