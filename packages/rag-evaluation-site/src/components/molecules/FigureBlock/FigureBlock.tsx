import type { HTMLAttributes, ReactNode } from 'react';
import { Caption } from '../../foundation';
import styles from './FigureBlock.module.css';

export interface FigureBlockProps extends HTMLAttributes<HTMLElement> {
  as?: 'figure' | 'div';
  children: ReactNode;
  caption?: ReactNode;
  legend?: ReactNode;
  frame?: 'none' | 'bordered' | 'inset';
}

export function FigureBlock({ as: Element = 'figure', children, caption, legend, frame = 'none', className, ...rest }: FigureBlockProps) {
  return (
    <Element className={[styles.root, styles[frame], className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="FigureBlock" {...rest}>
      <div className={styles.media}>{children}</div>
      {(caption || legend) && (
        <div className={styles.supporting}>
          {caption && <Caption className={styles.caption}>{caption}</Caption>}
          {legend && <div className={styles.legend}>{legend}</div>}
        </div>
      )}
    </Element>
  );
}
