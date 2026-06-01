import type { HTMLAttributes, ReactNode } from 'react';
import styles from './AppShell.module.css';

export interface AppShellProps extends HTMLAttributes<HTMLDivElement> {
  header?: ReactNode;
  sidebar?: ReactNode;
  children?: ReactNode;
}

export function AppShell({ header, sidebar, className, children, ...rest }: AppShellProps) {
  return (
    <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-layout="AppShell" {...rest}>
      {header && <header className={styles.header}>{header}</header>}
      <div className={styles.body}>
        {sidebar && <aside className={styles.sidebar}>{sidebar}</aside>}
        <main className={styles.main}>{children}</main>
      </div>
    </div>
  );
}
