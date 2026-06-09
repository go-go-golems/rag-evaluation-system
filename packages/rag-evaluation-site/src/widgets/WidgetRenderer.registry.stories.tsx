import type { Meta, StoryObj } from '@storybook/react-vite';
import { WidgetRenderer } from './WidgetRenderer';
import { defaultWidgetRegistry } from './defaultRegistry';
import { component, text } from './ir';

const meta = {
  title: 'Widget IR/Renderer/Registry Dispatch',
  component: WidgetRenderer,
  args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

export const RegisteredPanelLookup: Story = {
  args: {
    node: component('Panel', { title: 'Registered Panel', density: 'condensed' }, [
      component('Stack', { gap: 'sm' }, [
        text('This panel and stack render through the explicit widget registry.'),
        component('Button', { variant: 'primary' }, [text('Registry-backed button')]),
      ]),
    ]),
  },
};

export const UnknownWidgetFallback: Story = {
  args: {
    node: component('NotRegisteredYet', { reason: 'Phase 4 has not converted every widget yet.' }, [
      text('This child is intentionally ignored because the widget type is unknown.'),
    ]),
  },
};
