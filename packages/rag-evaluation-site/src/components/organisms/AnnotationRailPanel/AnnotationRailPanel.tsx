import type { HTMLAttributes } from 'react';
import type { TranscriptAnnotation } from '../../../context';
import { Caption } from '../../foundation';
import { Panel, Stack } from '../../layout';
import { AnnotationNoteCard } from '../../molecules';

export interface AnnotationRailPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  title?: string;
  annotations: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  onAnnotationSelect?: (annotationId: string) => void;
}

export function AnnotationRailPanel({ title = 'Annotations', annotations, selectedAnnotationId, onAnnotationSelect, ...rest }: AnnotationRailPanelProps) {
  return (
    <Panel title={title} fill data-rag-organism="AnnotationRailPanel" {...rest}>
      <Stack gap="sm">
        <Caption>{annotations.length} annotation{annotations.length === 1 ? '' : 's'}</Caption>
        {annotations.map((annotation) => (
          <button key={annotation.id} type="button" style={{ appearance: 'none', border: 0, background: 'transparent', padding: 0, textAlign: 'left', cursor: 'pointer' }} onClick={() => onAnnotationSelect?.(annotation.id)}>
            <AnnotationNoteCard annotation={annotation} selected={annotation.id === selectedAnnotationId} />
          </button>
        ))}
      </Stack>
    </Panel>
  );
}
