import { useState, useCallback } from 'react';
import { Panel, ScrollRegion } from '@go-go-golems/rag-evaluation-site';
import { DashboardGrid } from '@go-go-golems/rag-evaluation-site';
import { TranscriptReaderPanel, AnnotationRailPanel, type TranscriptFixture, transcriptFixture } from '@go-go-golems/rag-evaluation-site';
import styles from './TranscriptAnnotationPage.module.css';

const FIXTURE: TranscriptFixture = transcriptFixture;

export function TranscriptAnnotationPage() {
  const [selectedAnnotationId, setSelectedAnnotationId] = useState<string | undefined>(undefined);

  const handleSelectAnnotation = useCallback((id: string) => {
    setSelectedAnnotationId((prev) => (prev === id ? undefined : id));
  }, []);

  return (
    <DashboardGrid recipe="twoColumn" className={styles.root} data-rag-page="TranscriptAnnotationPage">
      <ScrollRegion>
        <TranscriptReaderPanel
          title={FIXTURE.title}
          subtitle={FIXTURE.subtitle}
          messages={FIXTURE.messages}
          annotations={FIXTURE.annotations}
          selectedAnnotationId={selectedAnnotationId}
          onAnnotationSelect={handleSelectAnnotation}
        />
      </ScrollRegion>
      <Panel className={styles.rail} title="Annotations" density="condensed">
        <ScrollRegion>
          <AnnotationRailPanel
            annotations={FIXTURE.annotations}
            selectedAnnotationId={selectedAnnotationId}
            onAnnotationSelect={handleSelectAnnotation}
          />
        </ScrollRegion>
      </Panel>
    </DashboardGrid>
  );
}

export default TranscriptAnnotationPage;
