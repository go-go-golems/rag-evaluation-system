import type { InputHTMLAttributes } from 'react';
import styles from './TextInput.module.css';

export interface TextInputProps extends InputHTMLAttributes<HTMLInputElement> {}

export function TextInput({ className, ...rest }: TextInputProps) {
  return <input className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-atom="TextInput" {...rest} />;
}
