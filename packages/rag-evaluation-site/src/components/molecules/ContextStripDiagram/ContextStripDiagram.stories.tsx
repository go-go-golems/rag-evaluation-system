import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { contextCobaltSandStyleSet, contextDefaultStyleSet, contextPaletteOptions, contextSignalOrangeStyleSet, contextSlateCoralStyleSet, contextStyleSetForPalette, contextWindowSnapshots, type ContextPaletteName, type ContextStyleSet } from '../../../context';
import { Button } from '../../atoms';
import { Inline, Panel, Stack } from '../../layout';
import { ContextLegend } from '../ContextLegend';
import { ContextStripDiagram, type ContextStripDiagramProps } from './ContextStripDiagram';

const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;

type PaletteControlsArgs = Omit<ContextStripDiagramProps, 'styleSet'> & { palette: ContextPaletteName; showLegend: boolean };

const meta = {
  title: 'Component Library/Molecules/ContextStripDiagram',
  component: ContextStripDiagram,
  args: { snapshot: deepBug!, styleSet: contextDefaultStyleSet },
} satisfies Meta<typeof ContextStripDiagram>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
  args: {
    snapshot: deepBug!,
    palette: 'Dusty Magenta / Blue',
    selectedPartId: 't14-file-reads',
    showLegend: true,
  },
  argTypes: {
    palette: { control: 'select', options: contextPaletteOptions },
    snapshot: { control: false },
    selectedPartId: { control: 'text' },
    showLabels: { control: 'boolean' },
    showLegend: { control: 'boolean' },
  },
  render: ({ palette, showLegend, ...args }) => {
    const styleSet = contextStyleSetForPalette(palette);
    return (
      <Panel title={`strip diagram · ${palette}`}>
        <Stack gap="sm">
          <ContextStripDiagram {...args} styleSet={styleSet} />
          {showLegend && <ContextLegend items={styleSet.legend} styles={styleSet.styles} size="sm" selectedId="result" />}
        </Stack>
      </Panel>
    );
  },
};

export const DenseSegments: Story = {
  render: () => (
    <Panel title="turn 14 context strip">
      <Stack gap="sm">
        <ContextStripDiagram snapshot={deepBug!} styleSet={contextDefaultStyleSet} />
        <ContextLegend items={contextDefaultStyleSet.legend} styles={contextDefaultStyleSet.styles} size="sm" />
      </Stack>
    </Panel>
  ),
};

export const SelectedSegment: Story = { render: () => <Panel title="selected file reads"><ContextStripDiagram snapshot={deepBug!} styleSet={contextDefaultStyleSet} selectedPartId="t14-file-reads" /></Panel> };

export const LimitComparison: Story = {
  render: () => (
    <Stack gap="md">
      <Panel title="at limit"><ContextStripDiagram snapshot={atLimit!} styleSet={contextDefaultStyleSet} /></Panel>
      <Panel title="before reclaim"><ContextStripDiagram snapshot={overBudget!} styleSet={contextDefaultStyleSet} /></Panel>
      <Panel title="at limit / signal orange"><ContextStripDiagram snapshot={atLimit!} styleSet={contextSignalOrangeStyleSet} /></Panel>
    </Stack>
  ),
};

const stripStyleSets: ContextStyleSet[] = [contextDefaultStyleSet, contextSignalOrangeStyleSet, contextSlateCoralStyleSet, contextCobaltSandStyleSet];

export const InteractiveStyleSwitcher: Story = {
  render: () => {
    const [styleSet, setStyleSet] = useState<ContextStyleSet>(contextDefaultStyleSet);
    return (
      <Panel title="switch the same strip between style sets">
        <Stack gap="sm">
          <Inline gap="xs" wrap>
            {stripStyleSets.map((candidate) => (
              <Button key={candidate.id ?? candidate.name} size="compact" selected={candidate.id === styleSet.id} onClick={() => setStyleSet(candidate)}>
                {candidate.name ?? candidate.id}
              </Button>
            ))}
          </Inline>
          <ContextStripDiagram snapshot={deepBug!} styleSet={styleSet} selectedPartId="t14-file-reads" />
          <ContextLegend items={styleSet.legend} styles={styleSet.styles} size="sm" selectedId="result" />
        </Stack>
      </Panel>
    );
  },
};
