import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextSlides, contextWindowSnapshots } from '../../../context';
import { ContextStackDiagram, FigureBlock, KeyPointList } from '../../molecules';
import { SlideShell } from './SlideShell';

const slide = contextSlides[1]!;
const snapshot = contextWindowSnapshots.find((item) => item.id === slide.snapshotId)!;

const figure = (
  <FigureBlock caption="context window — 200,000 tokens" frame="none">
    <ContextStackDiagram snapshot={snapshot} />
  </FigureBlock>
);

const points = <KeyPointList items={slide.notes} />;

const meta = {
  title: 'Design System/Layout/SlideShell',
  component: SlideShell,
  args: {
    eyebrow: slide.kicker,
    counter: '02 / 06',
    title: slide.title,
    primary: figure,
    secondary: points,
  },
} satisfies Meta<typeof SlideShell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const VisualFirst: Story = {};

export const ContentFirst: Story = {
  args: {
    primary: points,
    secondary: figure,
    primarySide: 'left',
    ratio: 'secondaryWide',
  },
};

export const VisualOnRight: Story = {
  args: {
    primary: figure,
    secondary: points,
    primarySide: 'right',
    ratio: 'primaryWide',
  },
};
