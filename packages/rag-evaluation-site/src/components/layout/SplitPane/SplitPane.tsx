import type { HTMLAttributes, ReactNode } from 'react';
import styles from './SplitPane.module.css';

export interface SplitPaneProps extends HTMLAttributes<HTMLDivElement> {
  left: ReactNode;
  right: ReactNode;
  ratio?: 'balanced' | 'leftNarrow' | 'rightNarrow' | 'course';
  divider?: boolean;
}

export function SplitPane({ left, right, ratio = 'balanced', divider = false, className, ...rest }: SplitPaneProps) {
  return (
    <div className={[styles.root, styles[ratio], divider ? styles.divider : '', className ?? ''].filter(Boolean).join(' ')} data-rag-layout="SplitPane" {...rest}>
      <div className={styles.pane}>{left}</div>
      <div className={styles.pane}>{right}</div>
    </div>
  );
}
