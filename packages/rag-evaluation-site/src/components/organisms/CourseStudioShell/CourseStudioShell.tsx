import type { HTMLAttributes, ReactNode } from 'react';
import { Caption } from '../../foundation';
import { SidebarShell } from '../../layout';
import { SidebarNav, type SidebarNavSection } from '../../molecules';
import styles from './CourseStudioShell.module.css';

export interface CourseStudioShellProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  sections: SidebarNavSection[];
  activeItemId?: string;
  onNavigate?: (itemId: string) => void;
  title?: ReactNode;
  subtitle?: ReactNode;
  sidebarFooter?: ReactNode;
  children?: ReactNode;
}

export function CourseStudioShell({
  sections,
  activeItemId,
  onNavigate,
  title = 'Context Window Engineering',
  subtitle = 'Course studio',
  sidebarFooter,
  className,
  children,
  ...rest
}: CourseStudioShellProps) {
  const header = (
    <div className={styles.header}>
      <div className={styles.title}>{title}</div>
      {subtitle && <Caption>{subtitle}</Caption>}
    </div>
  );

  return (
    <SidebarShell
      className={className}
      sidebarWidth={214}
      header={header}
      footer={sidebarFooter}
      sidebar={<SidebarNav sections={sections} activeItemId={activeItemId} onItemSelect={onNavigate} />}
      data-rag-organism="CourseStudioShell"
      {...rest}
    >
      <div className={styles.content}>{children}</div>
    </SidebarShell>
  );
}
