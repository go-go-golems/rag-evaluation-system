import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextCobaltSandStyleSet, contextDefaultStyleSet, contextPaletteOptions, contextStyleSetForPalette, contextWindowSnapshots, type ContextPaletteName } from '../../../context';
import { Panel, Stack } from '../../layout';
import { ContextTreemap, type ContextTreemapProps } from './ContextTreemap';

const [, deepBug, atLimit] = contextWindowSnapshots;

type PaletteControlsArgs = Omit<ContextTreemapProps, 'styleSet'> & { palette: ContextPaletteName };

const meta = { title: 'Component Library/Molecules/ContextTreemap', component: ContextTreemap, args: { snapshot: deepBug!, styleSet: contextDefaultStyleSet } } satisfies Meta<typeof ContextTreemap>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
  args: {
    snapshot: deepBug!,
    palette: 'Dusty Magenta / Blue',
    selectedPartId: 't14-file-reads',
  },
  argTypes: {
    palette: { control: 'select', options: contextPaletteOptions },
    snapshot: { control: false },
    selectedPartId: { control: 'text' },
  },
  render: ({ palette, ...args }) => (
    <Panel title={`treemap · ${palette}`}>
      <ContextTreemap {...args} styleSet={contextStyleSetForPalette(palette)} />
    </Panel>
  ),
};

export const ProportionalTokens: Story = { render: () => <Panel title="where the tokens went"><ContextTreemap snapshot={atLimit!} styleSet={contextDefaultStyleSet} /></Panel> };
export const SelectedTile: Story = { render: () => <Panel title="selected file reads"><ContextTreemap snapshot={deepBug!} styleSet={contextDefaultStyleSet} selectedPartId="t14-file-reads" /></Panel> };
export const Comparison: Story = { render: () => <Stack gap="md"><Panel title="turn 14"><ContextTreemap snapshot={deepBug!} styleSet={contextDefaultStyleSet} /></Panel><Panel title="turn 31"><ContextTreemap snapshot={atLimit!} styleSet={contextDefaultStyleSet} /></Panel><Panel title="turn 31 / cobalt sand"><ContextTreemap snapshot={atLimit!} styleSet={contextCobaltSandStyleSet} /></Panel></Stack> };
