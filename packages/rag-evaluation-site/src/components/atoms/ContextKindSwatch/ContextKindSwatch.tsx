import type { HTMLAttributes } from 'react';
import { getContextKindLabel, getContextKindSpec, type ContextDiagramStyle, type ContextPartKind } from '../../../context';
import styles from './ContextKindSwatch.module.css';

export type ContextKindSwatchSize = 'sm' | 'md' | 'lg';

export interface ContextKindSwatchProps extends HTMLAttributes<HTMLSpanElement> {
  kind: ContextPartKind;
  mode?: ContextDiagramStyle;
  size?: ContextKindSwatchSize;
  selected?: boolean;
}

export function ContextKindSwatch({ kind, mode = 'pattern', size = 'md', selected = false, className, style, ...rest }: ContextKindSwatchProps) {
  const spec = getContextKindSpec(kind, mode);
  const label = getContextKindLabel(kind);

  return (
    <span
      className={[
        styles.root,
        styles[size],
        spec.pattern ? styles[spec.pattern] : '',
        spec.dashed ? styles.dashed : '',
        spec.dotted ? styles.dotted : '',
        spec.inner ? styles.inner : '',
        selected ? styles.selected : '',
        className ?? '',
      ].filter(Boolean).join(' ')}
      data-rag-atom="ContextKindSwatch"
      data-kind={kind}
      data-mode={mode}
      aria-label={label}
      title={label}
      style={{
        '--rag-context-swatch-fill': spec.fill,
        '--rag-context-swatch-stroke': spec.stroke,
        '--rag-context-swatch-stroke-width': `${spec.strokeWidth ?? 1}px`,
        ...style,
      } as React.CSSProperties}
      {...rest}
    />
  );
}
