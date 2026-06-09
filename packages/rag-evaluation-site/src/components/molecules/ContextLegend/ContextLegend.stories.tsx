import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextDefaultStyleSet, contextPaletteOptions, contextStyleSetForPalette, contextThreeLabelStyleSets, type ContextPaletteName } from '../../../context';
import { Caption } from '../../foundation';
import { Panel, Stack } from '../../layout';
import { ContextLegend, type ContextLegendProps } from './ContextLegend';

type PaletteControlsArgs = Omit<ContextLegendProps, 'items' | 'styles'> & { palette: ContextPaletteName };

const meta = {
  title: 'Component Library/Molecules/ContextLegend',
  component: ContextLegend,
  args: { items: contextDefaultStyleSet.legend, styles: contextDefaultStyleSet.styles, size: 'md' },
} satisfies Meta<typeof ContextLegend>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
  args: {
    palette: 'Dusty Magenta / Blue',
    selectedId: 'active',
    size: 'md',
  },
  argTypes: {
    palette: { control: 'select', options: contextPaletteOptions },
    selectedId: { control: 'text' },
    size: { control: 'select', options: ['xs', 'sm', 'md', 'lg'] },
  },
  render: ({ palette, ...args }) => {
    const styleSet = contextStyleSetForPalette(palette);
    return <Panel title={`legend · ${palette}`}><ContextLegend {...args} items={styleSet.legend} styles={styleSet.styles} /></Panel>;
  },
};

export const DefaultStyleSet: Story = {
  render: () => <Panel title="default context-window legend"><ContextLegend items={contextDefaultStyleSet.legend} styles={contextDefaultStyleSet.styles} selectedId="active" /></Panel>,
};

export const Sizes: Story = {
  render: () => (
    <Stack gap="md">
      {(['xs', 'sm', 'md', 'lg'] as const).map((size) => (
        <Panel key={size} title={`${size} legend`}>
          <ContextLegend items={contextDefaultStyleSet.legend.slice(0, 8)} styles={contextDefaultStyleSet.styles} size={size} selectedId="result" />
        </Panel>
      ))}
    </Stack>
  ),
};

export const CustomThreeLabelLegend: Story = {
  render: () => {
    const styleSet = contextThreeLabelStyleSets[0]!;
    return (
      <Stack gap="sm">
        <Caption>Caller-defined legend labels and palette size (3 visible labels)</Caption>
        <ContextLegend items={styleSet.legend} styles={styleSet.styles} selectedId="evidence" />
      </Stack>
    );
  },
};
