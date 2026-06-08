import type { CSSProperties, HTMLAttributes, ReactNode } from 'react';
import styles from './SidebarShell.module.css';

export interface SidebarShellProps extends HTMLAttributes<HTMLDivElement> {
  sidebar: ReactNode;
  children?: ReactNode;
  sidebarWidth?: number | string;
  header?: ReactNode;
  footer?: ReactNode;
}

function toCssSize(value: number | string | undefined) {
  if (typeof value === 'number') return `${value}px`;
  return value;
}

export function SidebarShell({ sidebar, sidebarWidth = 220, header, footer, className, children, style, ...rest }: SidebarShellProps) {
  const shellStyle = {
    ...style,
    '--rag-sidebar-width': toCssSize(sidebarWidth),
  } as CSSProperties;

  return (
    <div className={[styles.root, className ?? ''].filter(Boolean).join(' ')} style={shellStyle} data-rag-layout="SidebarShell" {...rest}>
      <aside className={styles.sidebar} aria-label="Course navigation">
        {header && <div className={styles.sidebarHeader}>{header}</div>}
        <div className={styles.sidebarBody}>{sidebar}</div>
        {footer && <div className={styles.sidebarFooter}>{footer}</div>}
      </aside>
      <main className={styles.content}>{children}</main>
    </div>
  );
}
