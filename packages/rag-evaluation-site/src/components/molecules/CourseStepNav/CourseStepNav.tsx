import type { HTMLAttributes } from 'react';
import type { ContextCourseAgendaItem } from '../../../context';
import { Caption } from '../../foundation';
import styles from './CourseStepNav.module.css';

export interface CourseStepNavProps extends HTMLAttributes<HTMLOListElement> {
  items: ContextCourseAgendaItem[];
  activeItemId?: string;
  onItemSelect?: (itemId: string) => void;
}

export function CourseStepNav({ items, activeItemId, onItemSelect, className, ...rest }: CourseStepNavProps) {
  return (
    <ol className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-molecule="CourseStepNav" {...rest}>
      {items.map((item) => {
        const active = item.id === activeItemId;
        return (
          <li key={item.id} className={styles.item} data-active={active ? 'true' : undefined}>
            <button type="button" className={styles.button} onClick={() => onItemSelect?.(item.id)}>
              <span className={styles.number}>{item.number}</span>
              <span className={styles.body}>
                <span className={styles.title}>{item.title}</span>
                <span className={styles.description}>{item.description}</span>
              </span>
              {item.duration && <Caption>{item.duration}</Caption>}
            </button>
          </li>
        );
      })}
    </ol>
  );
}
