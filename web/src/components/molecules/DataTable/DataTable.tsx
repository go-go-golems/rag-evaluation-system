import type { ReactNode } from 'react';
import styles from './DataTable.module.css';

export interface DataTableColumn<T> {
  id: string;
  header: ReactNode;
  cell: (row: T) => ReactNode;
  align?: 'start' | 'end' | 'center';
  maxWidth?: number | string;
}

export interface DataTableProps<T> {
  columns: DataTableColumn<T>[];
  rows: T[];
  getRowKey: (row: T) => string;
  selectedKey?: string | null;
  onRowSelect?: (row: T) => void;
  emptyMessage?: ReactNode;
  className?: string;
}

export function DataTable<T>({ columns, rows, getRowKey, selectedKey, onRowSelect, emptyMessage, className }: DataTableProps<T>) {
  return (
    <table className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-component="DataTable">
      <thead>
        <tr>
          {columns.map((column) => <th key={column.id} className={styles[column.align ?? 'start']}>{column.header}</th>)}
        </tr>
      </thead>
      <tbody>
        {rows.length === 0 && emptyMessage && (
          <tr><td colSpan={columns.length} className={styles.empty}>{emptyMessage}</td></tr>
        )}
        {rows.map((row) => {
          const key = getRowKey(row);
          return (
            <tr key={key} className={[onRowSelect ? styles.selectable : '', selectedKey === key ? styles.selected : ''].filter(Boolean).join(' ')} onClick={onRowSelect ? () => onRowSelect(row) : undefined}>
              {columns.map((column) => (
                <td key={column.id} className={styles[column.align ?? 'start']} style={column.maxWidth ? { maxWidth: column.maxWidth } : undefined}>
                  {column.cell(row)}
                </td>
              ))}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
