import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextDefaultStyleSet, contextPaletteOptions, contextStyleSetForPalette, contextWindowSnapshots, type ContextPaletteName } from '../../../context';
import { ContextLegend, ContextStripDiagram } from '../index';
import { FigureBlock, type FigureBlockProps } from './FigureBlock';

const snapshot = contextWindowSnapshots[1]!;
type PaletteControlsArgs = Omit<FigureBlockProps, 'legend' | 'children'> & { palette: ContextPaletteName };

const meta = {
  title: 'Component Library/Molecules/FigureBlock',
  component: FigureBlock,
  args: {
    caption: 'context window — 200,000 tokens',
    legend: <ContextLegend items={contextDefaultStyleSet.legend} styles={contextDefaultStyleSet.styles} size="sm" />,
    children: <ContextStripDiagram snapshot={snapshot} styleSet={contextDefaultStyleSet} />,
  },
} satisfies Meta<typeof FigureBlock>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
  args: {
    palette: 'Dusty Magenta / Blue',
    caption: 'context window — 200,000 tokens',
    frame: 'none',
  },
  argTypes: {
    palette: { control: 'select', options: contextPaletteOptions },
    frame: { control: 'select', options: ['none', 'bordered'] },
  },
  render: ({ palette, ...args }) => {
    const styleSet = contextStyleSetForPalette(palette);
    return (
      <FigureBlock {...args} legend={<ContextLegend items={styleSet.legend} styles={styleSet.styles} size="sm" selectedId="result" />}>
        <ContextStripDiagram snapshot={snapshot} styleSet={styleSet} selectedPartId="t14-file-reads" />
      </FigureBlock>
    );
  },
};

export const DiagramFigure: Story = {};
export const BorderedFigure: Story = { args: { frame: 'bordered' } };
