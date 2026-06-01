import React from 'react';
import { CorpusChunk } from '../../../services/api';

interface ChunkTimelineBarProps {
  chunks: CorpusChunk[];
  selectedIdx: number | null;
  onSelect: (idx: number | null) => void;
}

export const ChunkTimelineBar: React.FC<ChunkTimelineBarProps> = ({ chunks, selectedIdx, onSelect }) => {
  if (chunks.length === 0) return null;

  const maxEnd = Math.max(...chunks.map((c) => c.end_offset));
  if (maxEnd <= 0) return null;

  return (
    <div className="chunk-bar-container">
      {chunks.map((chunk, idx) => {
        const left = (chunk.start_offset / maxEnd) * 100;
        const width = Math.max(((chunk.end_offset - chunk.start_offset) / maxEnd) * 100, 0.5);
        return (
          <div
            key={chunk.id}
            className={`chunk-bar ${chunk.embedding?.present ? 'embedded' : 'not-embedded'} ${selectedIdx === idx ? 'selected' : ''}`}
            style={{
              left: `${left}%`,
              width: `${width}%`,
            }}
            title={`#${chunk.chunk_index} ${chunk.start_offset}–${chunk.end_offset} ${chunk.embedding?.present ? '●' : '○'}`}
            onClick={() => onSelect(idx === selectedIdx ? null : idx)}
          />
        );
      })}
    </div>
  );
};
