import type { InputHTMLAttributes, LabelHTMLAttributes, ReactNode } from 'react';
import styles from './CheckboxRow.module.css';

export interface CheckboxRowProps extends Omit<LabelHTMLAttributes<HTMLLabelElement>, 'onChange'> {
  checked: boolean;
  onChange: InputHTMLAttributes<HTMLInputElement>['onChange'];
  children?: ReactNode;
}

export function CheckboxRow({ checked, onChange, className, children, ...rest }: CheckboxRowProps) {
  return (
    <label className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-atom="CheckboxRow" {...rest}>
      <input type="checkbox" checked={checked} onChange={onChange} />
      {children}
    </label>
  );
}
