import { Caption } from '../../foundation';
import { Panel, Stack } from '../../layout';
import styles from './QueryPresetList.module.css';

export interface QueryPresetListProps {
  queries: string[];
  onSelect: (query: string) => void;
}

export function QueryPresetList({ queries, onSelect }: QueryPresetListProps) {
  return (
    <Panel title="Test Queries" density="condensed" data-rag-component="QueryPresetList">
      <Stack gap="xs">
        {queries.map(query => (
          <button key={query} type="button" className={styles.item} onClick={() => onSelect(query)}>
            <Caption tone="inherit">{query}</Caption>
          </button>
        ))}
      </Stack>
    </Panel>
  );
}
