import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { FormRow, Inline, Panel, Stack } from '../../layout';
import { StatusText } from '../../foundation';
import { SelectInput } from './SelectInput';

const meta = { title: 'Design System/Atoms/SelectInput', component: SelectInput } satisfies Meta<typeof SelectInput>;
export default meta;
type Story = StoryObj<typeof meta>;

function ControlledSelect() {
  const [value, setValue] = useState('mine');

  return (
    <SelectInput value={value} onChange={(event) => setValue(event.currentTarget.value)}>
      <option value="all">All uploads</option>
      <option value="mine">My uploads</option>
      <option value="anonymous">Anonymous uploads</option>
    </SelectInput>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <Panel title={title} density="condensed">
      <Stack gap="sm">{children}</Stack>
    </Panel>
  );
}

export const Basic: Story = {
  render: () => (
    <SelectInput defaultValue="hybrid">
      <option value="bm25">bm25</option>
      <option value="vector">vector</option>
      <option value="hybrid">hybrid</option>
    </SelectInput>
  ),
};

export const States: Story = {
  render: () => (
    <Stack gap="md">
      <Section title="Common states">
        <FormRow
          label="Default"
          control={(
            <SelectInput defaultValue="hybrid">
              <option value="bm25">bm25</option>
              <option value="vector">vector</option>
              <option value="hybrid">hybrid</option>
            </SelectInput>
          )}
          hint="Normal native select input."
        />
        <FormRow
          label="Disabled"
          control={(
            <SelectInput defaultValue="all" disabled>
              <option value="all">All uploads</option>
              <option value="mine">My uploads</option>
            </SelectInput>
          )}
          hint="Native disabled state."
        />
        <FormRow
          label="Invalid"
          control={(
            <SelectInput defaultValue="" aria-invalid="true" required>
              <option value="" disabled>Choose one…</option>
              <option value="public">Public</option>
              <option value="private">Private</option>
            </SelectInput>
          )}
          hint={<StatusText status="error" icon>Select a visibility before saving.</StatusText>}
        />
        <FormRow
          label="Controlled"
          control={<ControlledSelect />}
          hint="Controlled selection for future settings/filter forms."
        />
      </Section>
      <Section title="Inline density">
        <Inline gap="sm" wrap>
          <SelectInput defaultValue="all">
            <option value="all">All uploads</option>
            <option value="mine">My uploads</option>
          </SelectInput>
          <SelectInput defaultValue="jsonl">
            <option value="jsonl">jsonl</option>
            <option value="markdown">markdown</option>
            <option value="text">text</option>
          </SelectInput>
          <SelectInput defaultValue="newest">
            <option value="newest">Newest first</option>
            <option value="oldest">Oldest first</option>
          </SelectInput>
        </Inline>
      </Section>
    </Stack>
  ),
};

export const InFormRows: Story = {
  render: () => (
    <Panel title="Filter form sketch" density="condensed">
      <Stack gap="sm">
        <FormRow
          label="Scope"
          control={(
            <SelectInput defaultValue="mine">
              <option value="all">All uploads</option>
              <option value="mine">My uploads</option>
            </SelectInput>
          )}
          hint="Example reusable filter row."
        />
        <FormRow
          label="Visibility"
          control={(
            <SelectInput defaultValue="public">
              <option value="public">Public</option>
              <option value="private">Private</option>
            </SelectInput>
          )}
          hint="A settings-style select field."
        />
      </Stack>
    </Panel>
  ),
};
