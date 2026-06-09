import type { CSSProperties } from 'react';
import { ScrollRegion } from './ScrollRegion';
import { defineWidget } from '../../../widgets/registry';
import type { ScrollRegionWidgetProps } from '../../../widgets/ir';

export const scrollRegionWidget = defineWidget<ScrollRegionWidgetProps>({
  type: 'ScrollRegion',
  module: 'ui.dsl',
  render: (props, children) => (
    <ScrollRegion className={props.className} axis={props.axis} style={props.style as CSSProperties | undefined}>
      {children}
    </ScrollRegion>
  ),
});
