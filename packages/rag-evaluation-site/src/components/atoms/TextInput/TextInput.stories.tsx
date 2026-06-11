import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { FormRow, Inline, Panel, Stack } from '../../layout';
import { Caption, StatusText } from '../../foundation';
import { TextInput } from './TextInput';

const meta = { title: 'Design System/Atoms/TextInput', component: TextInput } satisfies Meta<typeof TextInput>;
export default meta;
type Story = StoryObj<typeof meta>;

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <Panel title={title} density="condensed">
      <Stack gap="sm">{children}</Stack>
    </Panel>
  );
}

function ControlledExample() {
  const [value, setValue] = useState('Manuel');
  const remaining = 40 - value.length;

  return (
    <Stack gap="xs">
      <TextInput
        aria-describedby="display-name-counter"
        maxLength={40}
        value={value}
        onChange={(event) => setValue(event.currentTarget.value)}
      />
      <Caption id="display-name-counter" tone={remaining < 5 ? 'warning' : 'muted'}>
        {remaining} characters remaining
      </Caption>
    </Stack>
  );
}

export const Examples: Story = {
  render: () => (
    <Stack gap="xs">
      <TextInput defaultValue="fixed-1200-150" />
      <TextInput type="number" defaultValue={10} />
      <TextInput placeholder="Placeholder" />
    </Stack>
  ),
};

export const States: Story = {
  render: () => (
    <Stack gap="md">
      <Section title="Common states">
        <FormRow label="Default" control={<TextInput defaultValue="sample-pi-session.jsonl" />} hint="Normal editable text input." />
        <FormRow label="Placeholder" control={<TextInput placeholder="Enter display name" />} hint="Empty input with placeholder copy." />
        <FormRow label="Disabled" control={<TextInput defaultValue="Locked by system" disabled />} hint="Native disabled state; not submitted by forms." />
        <FormRow label="Read-only" control={<TextInput defaultValue="session-9a7c" readOnly />} hint="Readable but intentionally not editable." />
        <FormRow
          label="Invalid"
          control={<TextInput defaultValue="A name with control character" aria-invalid="true" />}
          hint={<StatusText status="error" icon>Display name contains unsupported characters.</StatusText>}
        />
        <FormRow label="Required" control={<TextInput required placeholder="Required field" />} hint="Uses native required validation." />
      </Section>
      <Section title="Input types">
        <Inline gap="sm" wrap>
          <TextInput type="text" defaultValue="text" />
          <TextInput type="number" min={1} max={100} defaultValue={25} />
          <TextInput type="email" placeholder="name@example.com" />
          <TextInput type="url" placeholder="https://example.com" />
          <TextInput type="password" defaultValue="secret" />
        </Inline>
      </Section>
    </Stack>
  ),
};

export const ControlledWithCounter: Story = {
  render: () => (
    <Section title="Controlled display-name field">
      <FormRow
        label="Display name"
        control={<ControlledExample />}
        hint="This mirrors the future MINIVIZ settings field without adding app-specific behavior."
      />
    </Section>
  ),
};

export const InFormRows: Story = {
  render: () => (
    <Panel title="Settings form sketch" density="condensed">
      <Stack gap="sm">
        <FormRow label="Display name" control={<TextInput defaultValue="Anonymous visitor" />} hint="Shown next to uploads from this browser." />
        <FormRow label="Max rows" control={<TextInput type="number" min={1} max={500} defaultValue={100} />} hint="Numeric input in the same row anatomy." />
        <FormRow label="Session id" control={<TextInput defaultValue="sess_8f2c…" readOnly />} hint="Read-only system value." />
      </Stack>
    </Panel>
  ),
};
