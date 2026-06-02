import type { HTMLAttributes, ReactNode } from 'react';
import styles from './DashboardGrid.module.css';

export type DashboardGridRecipe = 'searchWorkbench' | 'corpusExplorer' | 'twoColumn';

export interface DashboardGridProps extends HTMLAttributes<HTMLDivElement> {
  recipe?: DashboardGridRecipe;
  children?: ReactNode;
}

const recipeClass: Record<DashboardGridRecipe, string> = {
  searchWorkbench: styles.searchWorkbench,
  corpusExplorer: styles.corpusExplorer,
  twoColumn: styles.twoColumn,
};

export function DashboardGrid({ recipe = 'twoColumn', className, children, ...rest }: DashboardGridProps) {
  return (
    <div className={[styles.root, recipeClass[recipe], className ?? ''].filter(Boolean).join(' ')} data-rag-layout="DashboardGrid" data-rag-layout-recipe={recipe} {...rest}>
      {children}
    </div>
  );
}
