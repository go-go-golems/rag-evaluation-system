import type { ReactNode } from 'react';
import styles from './AppNav.module.css';

export interface AppNavItem {
  id: string;
  label: ReactNode;
}

export interface AppNavProps {
  brand: ReactNode;
  items: AppNavItem[];
  activeItemId: string;
  onItemSelect: (itemId: string) => void;
}

export function AppNav({ brand, items, activeItemId, onItemSelect }: AppNavProps) {
  return (
    <nav className={styles.root} aria-label="Primary" data-rag-component="AppNav">
      <span className={styles.brand}>{brand}</span>
      {items.map((item) => {
        const isActive = item.id === activeItemId;
        return (
          <button
            key={item.id}
            type="button"
            className={[styles.item, isActive ? styles.active : ''].filter(Boolean).join(' ')}
            aria-current={isActive ? 'page' : undefined}
            onClick={() => onItemSelect(item.id)}
          >
            {item.label}
          </button>
        );
      })}
    </nav>
  );
}
