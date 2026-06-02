import type { HTMLAttributes, ReactNode } from 'react';
import styles from './CodeText.module.css';

export interface CodeTextProps extends HTMLAttributes<HTMLElement> {
  as?: 'code' | 'span' | 'div';
  tone?: 'primary' | 'muted' | 'accent';
  display?: 'inline' | 'block';
  copyable?: boolean;
  children?: ReactNode;
}

const toneClass = {
  primary: '',
  muted: styles.muted,
  accent: styles.accent,
};

export function CodeText({
  as: Component = 'code',
  tone = 'primary',
  display = 'inline',
  copyable = false,
  className,
  children,
  onClick,
  ...rest
}: CodeTextProps) {
  const rootClassName = [
    styles.root,
    display === 'inline' ? styles.inline : styles.block,
    toneClass[tone],
    copyable ? styles.copyable : '',
    className ?? '',
  ].filter(Boolean).join(' ');

  const handleClick: React.MouseEventHandler<HTMLElement> = (event) => {
    if (copyable && typeof children === 'string') {
      navigator.clipboard?.writeText(children).catch(() => {});
    }
    onClick?.(event);
  };

  return <Component className={rootClassName} onClick={handleClick} {...rest}>{children}</Component>;
}
