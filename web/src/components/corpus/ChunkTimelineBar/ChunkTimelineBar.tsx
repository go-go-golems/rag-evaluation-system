import React from 'react';
import { CorpusChunk } from '../../../services/api';
import styles from './ChunkTimelineBar.module.css';

interface ChunkTimelineBarProps {
  chunks: CorpusChunk[];
  selectedIdx: number | null;
  onSelect: (idx: number | null) => void;
}

export const ChunkTimelineBar: React.FC<ChunkTimelineBarProps> = ({ chunks, selectedIdx, onSelect }) => {
  if (chunks.length === 0) return null;

  const maxEnd = Math.max(...chunks.map((chunk) => chunk.end_offset));
  if (maxEnd <= 0) return null;

  return (
    <div className={styles.root} data-rag-component="ChunkTimelineBar">
      {chunks.map((chunk, idx) => {
        const left = (chunk.start_offset / maxEnd) * 100;
        const width = Math.max(((chunk.end_offset - chunk.start_offset) / maxEnd) * 100, 0.5);
        const embedded = Boolean(chunk.embedding?.present);
        return (
          <button
            key={chunk.id}
            type="button"
            className={[
              styles.chunk,
              embedded ? styles.embedded : styles.notEmbedded,
              selectedIdx === idx ? styles.selected : '',
            ].filter(Boolean).join(' ')}
            style={{
              left: `${left}%`,
              width: `${width}%`,
            }}
            title={`#${chunk.chunk_index} ${chunk.start_offset}–${chunk.end_offset} ${embedded ? '●' : '○'}`}
            aria-label={`Select chunk ${chunk.chunk_index}, offsets ${chunk.start_offset} to ${chunk.end_offset}, ${embedded ? 'embedded' : 'not embedded'}`}
            aria-pressed={selectedIdx === idx}
            onClick={() => onSelect(idx === selectedIdx ? null : idx)}
          />
        );
      })}
    </div>
  );
};
