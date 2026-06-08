import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots } from '../../../context';
import { ContextLegend, ContextStripDiagram } from '../index';
import { FigureBlock } from './FigureBlock';
const snapshot = contextWindowSnapshots[1]!;
const meta = { title: 'Component Library/Molecules/FigureBlock', component: FigureBlock, args: { caption: 'context window — 200,000 tokens', legend: <ContextLegend compact />, children: <ContextStripDiagram snapshot={snapshot} /> } } satisfies Meta<typeof FigureBlock>;
export default meta; type Story = StoryObj<typeof meta>;
export const DiagramFigure: Story = {};
export const BorderedFigure: Story = { args: { frame: 'bordered' } };
