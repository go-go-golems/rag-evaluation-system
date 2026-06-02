import type { HTMLAttributes, ReactNode } from 'react';
import styles from './FormRow.module.css';

export interface FormRowProps extends HTMLAttributes<HTMLDivElement> {
  label: ReactNode;
  control: ReactNode;
  hint?: ReactNode;
}

export function FormRow({ label, control, hint, className, ...rest }: FormRowProps) {
  return (
    <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-layout="FormRow" {...rest}>
      <div className={styles.label}>{label}</div>
      <div className={styles.control}>{control}</div>
      {hint && <div className={styles.hint}>{hint}</div>}
    </div>
  );
}
