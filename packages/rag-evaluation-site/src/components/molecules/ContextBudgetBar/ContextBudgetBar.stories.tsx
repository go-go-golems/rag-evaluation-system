import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextCobaltSandStyleSet, contextDefaultStyleSet, contextPaletteOptions, contextSignalOrangeStyleSet, contextStyleSetForPalette, contextWindowSnapshots, type ContextPaletteName } from '../../../context';
import { Panel, Stack } from '../../layout';
import { ContextBudgetBar, type ContextBudgetBarProps } from './ContextBudgetBar';

const [underBudget, selectedBudget, nearLimit, overBudget] = contextWindowSnapshots;

type PaletteControlsArgs = Omit<ContextBudgetBarProps, 'styleSet'> & { palette: ContextPaletteName };

const meta = {
  title: 'Component Library/Molecules/ContextBudgetBar',
  component: ContextBudgetBar,
  args: { snapshot: selectedBudget!, styleSet: contextDefaultStyleSet },
} satisfies Meta<typeof ContextBudgetBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
  args: {
    snapshot: selectedBudget!,
    palette: 'Dusty Magenta / Blue',
    selectedPartId: 't14-file-reads',
    showLegend: true,
  },
  argTypes: {
    palette: { control: 'select', options: contextPaletteOptions },
    snapshot: { control: false },
    selectedPartId: { control: 'text' },
    showLegend: { control: 'boolean' },
  },
  render: ({ palette, ...args }) => (
    <Panel title={`budget bar · ${palette}`}>
      <ContextBudgetBar {...args} styleSet={contextStyleSetForPalette(palette)} />
    </Panel>
  ),
};

export const BudgetStates: Story = {
  render: () => (
    <Stack gap="md">
      <Panel title="under budget"><ContextBudgetBar snapshot={underBudget!} styleSet={contextDefaultStyleSet} /></Panel>
      <Panel title="near limit"><ContextBudgetBar snapshot={nearLimit!} styleSet={contextDefaultStyleSet} /></Panel>
      <Panel title="over budget"><ContextBudgetBar snapshot={overBudget!} styleSet={contextDefaultStyleSet} /></Panel>
    </Stack>
  ),
};

export const SelectedSegment: Story = {
  render: () => <Panel title="selected context tenant"><ContextBudgetBar snapshot={selectedBudget!} styleSet={contextDefaultStyleSet} selectedPartId="t14-file-reads" /></Panel>,
};

export const SameBudgetDifferentPalettes: Story = {
  render: () => (
    <Stack gap="md">
      <Panel title="Dusty Magenta / Blue"><ContextBudgetBar snapshot={selectedBudget!} styleSet={contextDefaultStyleSet} /></Panel>
      <Panel title="Signal Orange / Cyan"><ContextBudgetBar snapshot={selectedBudget!} styleSet={contextSignalOrangeStyleSet} /></Panel>
      <Panel title="Cobalt / Sand"><ContextBudgetBar snapshot={selectedBudget!} styleSet={contextCobaltSandStyleSet} /></Panel>
    </Stack>
  ),
};
