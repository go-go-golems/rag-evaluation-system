import { TranscriptRoleBadge } from './TranscriptRoleBadge';
import { defineWidget } from '../../../widgets/registry';
import type { TranscriptRoleBadgeWidgetProps } from '../../../widgets/ir';

export const transcriptRoleBadgeWidget = defineWidget<TranscriptRoleBadgeWidgetProps>({
  type: 'TranscriptRoleBadge',
  module: 'context_window.dsl',
  render: (props) => <TranscriptRoleBadge className={props.className} role={props.role} name={props.name} />,
});
