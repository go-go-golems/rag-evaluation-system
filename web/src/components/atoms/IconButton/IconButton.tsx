import type { ButtonHTMLAttributes, ReactNode } from 'react';
import styles from './IconButton.module.css';

export type IconButtonSize = 'compact' | 'normal';

export interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  size?: IconButtonSize;
  label?: string;
  children?: ReactNode;
}

export function IconButton({ size = 'compact', label, className, children, ...rest }: IconButtonProps) {
  return (
    <button
      type="button"
      aria-label={label}
      title={rest.title ?? label}
      className={[
        styles.root,
        size === 'compact' ? styles.compact : styles.normal,
        className ?? '',
      ].filter(Boolean).join(' ')}
      data-rag-atom="IconButton"
      {...rest}
    >
      {children}
    </button>
  );
}
