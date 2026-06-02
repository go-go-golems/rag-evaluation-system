import type { SelectHTMLAttributes } from 'react';
import styles from './SelectInput.module.css';

export interface SelectInputProps extends SelectHTMLAttributes<HTMLSelectElement> {}

export function SelectInput({ className, children, ...rest }: SelectInputProps) {
  return <select className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-atom="SelectInput" {...rest}>{children}</select>;
}
