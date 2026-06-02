import type { ButtonHTMLAttributes, ReactNode } from 'react';
import styles from './Button.module.css';

export type ButtonVariant = 'default' | 'primary';
export type ButtonSize = 'normal' | 'compact';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  selected?: boolean;
  children?: ReactNode;
}

export function Button({ variant = 'default', size = 'normal', selected = false, className, children, ...rest }: ButtonProps) {
  return (
    <button
      type="button"
      className={[
        styles.root,
        variant === 'primary' ? styles.primary : '',
        selected ? styles.selected : '',
        size === 'compact' ? styles.compact : styles.normal,
        className ?? '',
      ].filter(Boolean).join(' ')}
      data-rag-atom="Button"
      {...rest}
    >
      {children}
    </button>
  );
}
