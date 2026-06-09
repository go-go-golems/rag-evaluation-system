import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextDefaultStyleSet, contextPaletteOptions, contextSignalOrangeStyleSet, contextStyleSetForPalette, contextWindowSnapshots, type ContextPaletteName } from '../../../context';
import { Panel, Stack } from '../../layout';
import { ContextStackDiagram, type ContextStackDiagramProps } from './ContextStackDiagram';

const [, deepBug, atLimit] = contextWindowSnapshots;

type PaletteControlsArgs = Omit<ContextStackDiagramProps, 'styleSet'> & { palette: ContextPaletteName };

const meta = { title: 'Component Library/Molecules/ContextStackDiagram', component: ContextStackDiagram, args: { snapshot: deepBug!, styleSet: contextDefaultStyleSet } } satisfies Meta<typeof ContextStackDiagram>;
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
    <Panel title={`stack diagram · ${palette}`}>
      <ContextStackDiagram {...args} styleSet={contextStyleSetForPalette(palette)} />
    </Panel>
  ),
};

export const GroupedContextWindow: Story = { render: () => <Panel title="layered call"><ContextStackDiagram snapshot={deepBug!} styleSet={contextDefaultStyleSet} /></Panel> };
export const SelectedLayer: Story = { render: () => <Panel title="selected scratchpad"><ContextStackDiagram snapshot={atLimit!} styleSet={contextDefaultStyleSet} selectedPartId="t31-scratchpad" /></Panel> };
export const Comparison: Story = { render: () => <Stack gap="md"><Panel title="turn 14"><ContextStackDiagram snapshot={deepBug!} styleSet={contextDefaultStyleSet} /></Panel><Panel title="turn 31"><ContextStackDiagram snapshot={atLimit!} styleSet={contextDefaultStyleSet} /></Panel><Panel title="turn 31 / signal orange"><ContextStackDiagram snapshot={atLimit!} styleSet={contextSignalOrangeStyleSet} /></Panel></Stack> };
