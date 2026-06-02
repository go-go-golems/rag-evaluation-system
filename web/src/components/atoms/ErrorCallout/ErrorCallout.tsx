import type { HTMLAttributes, ReactNode } from 'react';
import styles from './ErrorCallout.module.css';

export interface ErrorCalloutProps extends HTMLAttributes<HTMLDivElement> {
  children?: ReactNode;
}

export function ErrorCallout({ className, children, ...rest }: ErrorCalloutProps) {
  return <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} role="alert" data-rag-atom="ErrorCallout" {...rest}>{children}</div>;
}
