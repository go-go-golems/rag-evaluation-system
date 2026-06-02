import type { HTMLAttributes } from 'react';
import styles from './Divider.module.css';

export interface DividerProps extends HTMLAttributes<HTMLHRElement> {
  orientation?: 'horizontal' | 'vertical';
}

export function Divider({ orientation = 'horizontal', className, ...rest }: DividerProps) {
  return (
    <hr
      aria-orientation={orientation}
      className={[styles.root, orientation === 'vertical' ? styles.vertical : '', className ?? ''].filter(Boolean).join(' ')}
      {...rest}
    />
  );
}
