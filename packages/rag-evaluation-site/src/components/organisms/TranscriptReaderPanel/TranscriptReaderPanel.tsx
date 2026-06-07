import type { HTMLAttributes } from 'react';
import type { TranscriptAnnotation, TranscriptMessage } from '../../../context';
import { Caption } from '../../foundation';
import { Panel, Stack } from '../../layout';
import { TranscriptMessageCard } from '../../molecules';

export interface TranscriptReaderPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  title?: string;
  subtitle?: string;
  messages: TranscriptMessage[];
  annotations?: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  onAnnotationSelect?: (annotationId: string) => void;
}

export function TranscriptReaderPanel({ title = 'Transcript', subtitle, messages, annotations = [], selectedAnnotationId, onAnnotationSelect, ...rest }: TranscriptReaderPanelProps) {
  return (
    <Panel title={title} fill data-rag-organism="TranscriptReaderPanel" {...rest}>
      <Stack gap="sm">
        {subtitle && <Caption>{subtitle}</Caption>}
        {messages.map((message) => (
          <TranscriptMessageCard
            key={message.id}
            message={message}
            annotations={annotations}
            selectedAnnotationId={selectedAnnotationId}
            onAnnotationSelect={onAnnotationSelect}
          />
        ))}
      </Stack>
    </Panel>
  );
}
