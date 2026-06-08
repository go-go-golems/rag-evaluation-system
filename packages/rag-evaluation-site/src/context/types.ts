export type ContextPartKind =
  | 'system'
  | 'instruction'
  | 'context'
  | 'conversation'
  | 'summary'
  | 'retrieval'
  | 'tool'
  | 'result'
  | 'generated'
  | 'annotation'
  | 'course'
  | 'active'
  | 'evicted'
  | 'empty'
  | 'other';

export type TranscriptRole =
  | 'system'
  | 'developer'
  | 'user'
  | 'assistant'
  | 'tool'
  | 'other';

export type ContextJsonPrimitive = string | number | boolean | null;
export type ContextJsonValue = ContextJsonPrimitive | ContextJsonValue[] | { [key: string]: ContextJsonValue };
export type ContextJsonRecord = Record<string, ContextJsonValue>;

export interface ContextWindowPart {
  id: string;
  label: string;
  kind: ContextPartKind;
  tokens: number;
  note?: string;
  contentPreview?: string;
  content?: string;
  sourceId?: string;
  startToken?: number;
  endToken?: number;
  metadata?: ContextJsonRecord;
}

export interface ContextWindowSnapshot {
  id: string;
  title: string;
  subtitle?: string;
  limit: number;
  parts: ContextWindowPart[];
  selectedPartId?: string;
  metadata?: ContextJsonRecord;
}

export type ContextDiagramView = 'strip' | 'stack' | 'budget' | 'treemap';
export type ContextDiagramStyle = 'pattern' | 'tone' | 'outline';

export interface ContextDiagramSegment {
  id: string;
  label: string;
  kind: ContextPartKind;
  tokens: number;
  note?: string;
  metadata?: ContextJsonRecord;
}

export interface TranscriptAnnotation {
  id: string;
  targetMessageId: string;
  kind: ContextPartKind;
  label: string;
  text: string;
  confidence?: number;
  metadata?: ContextJsonRecord;
}

export interface TranscriptMessage {
  id: string;
  role: TranscriptRole;
  text: string;
  tokens?: number;
  name?: string;
  timestamp?: string;
  annotationIds?: string[];
  metadata?: ContextJsonRecord;
}

export interface TranscriptFixture {
  id: string;
  title: string;
  subtitle?: string;
  messages: TranscriptMessage[];
  annotations: TranscriptAnnotation[];
  selectedAnnotationId?: string;
}

export interface ContextCourseInstructor {
  name: string;
  role?: string;
  bio?: string;
}

export interface ContextCourseAgendaItem {
  id: string;
  number: string;
  title: string;
  description: string;
  duration?: string;
}

export interface ContextCourse {
  id: string;
  kicker?: string;
  title: string;
  tagline: string;
  when?: string;
  where?: string;
  format?: string;
  price?: string;
  instructor?: ContextCourseInstructor;
  blurb?: string;
  outcomes: string[];
  agenda: ContextCourseAgendaItem[];
}

export interface ContextSlide {
  id: string;
  number: string;
  kicker?: string;
  title: string;
  view: ContextDiagramView;
  snapshotId: string;
  notes: string[];
}

export interface ContextHandoutDocument {
  id: string;
  title: string;
  file: string;
  format: string;
  size?: string;
  description: string;
  body: string;
}

export type AnchoredCommentStatus = 'open' | 'resolved';

export interface AnchoredComment {
  id: string;
  anchorX: number;
  anchorY: number;
  author: string;
  text: string;
  time?: string;
  status?: AnchoredCommentStatus;
  metadata?: ContextJsonRecord;
}

export interface ContextHandoutBundle {
  intro: string;
  docs: ContextHandoutDocument[];
}
