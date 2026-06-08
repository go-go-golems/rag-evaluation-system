import type { HTMLAttributes, KeyboardEvent } from 'react';
import { getContextKindLabel, type ContextDiagramStyle, type ContextWindowSnapshot } from '../../../context';
import styles from './ContextStripDiagram.module.css';

export interface ContextStripDiagramProps extends HTMLAttributes<HTMLDivElement> {
  snapshot: ContextWindowSnapshot;
  mode?: ContextDiagramStyle;
  selectedPartId?: string;
  showLabels?: boolean;
  onPartSelect?: (partId: string) => void;
}

function formatTokens(tokens: number) {
  return `${tokens.toLocaleString()} tok`;
}

function handlePartKeyDown(event: KeyboardEvent<HTMLDivElement>, partId: string, onPartSelect?: (partId: string) => void) {
  if (!onPartSelect) return;
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault();
    onPartSelect(partId);
  }
}

export function ContextStripDiagram({ snapshot, mode = 'pattern', selectedPartId, showLabels = true, onPartSelect, className, ...rest }: ContextStripDiagramProps) {
  const effectiveSelectedPartId = selectedPartId ?? snapshot.selectedPartId;
  const totalWidth = snapshot.limit || snapshot.parts.reduce((sum, part) => sum + part.tokens, 0);

  return (
    <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="ContextStripDiagram" data-mode={mode} {...rest}>
      <div className={styles.strip} role="img" aria-label={`${snapshot.title} context strip`}>
        {snapshot.parts.map((part) => {
          const width = Math.max(2, (part.tokens / totalWidth) * 100);
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
            >
              {showLabels && width >= 7 && <span className={styles.label}>{part.label}</span>}
            </div>
          );
        })}
      </div>
      <div className={styles.caption}>{snapshot.title} · {snapshot.limit.toLocaleString()} token window</div>
    </div>
  );
}
