import { createElement } from 'react';
import { ContextStudioNavIcon, type ContextStudioNavIconId } from '../../atoms';
import type { SidebarNavSection } from '../../molecules';

function icon(id: ContextStudioNavIconId) {
  return createElement(ContextStudioNavIcon, { id });
}

export const courseStudioNavSections: SidebarNavSection[] = [
  { id: 'session', label: 'Session', items: [{ id: 'upload', label: 'Upload', icon: icon('upload') }, { id: 'visualize', label: 'Visualize', icon: icon('visualize') }, { id: 'transcript', label: 'Transcript', icon: icon('transcript') }] },
];
