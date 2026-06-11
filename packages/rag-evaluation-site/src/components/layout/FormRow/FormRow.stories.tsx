import type { Meta, StoryObj } from '@storybook/react-vite';
import { Button, CheckboxRow, SelectInput, TextInput } from '../../atoms';
import { StatusText } from '../../foundation';
import { Inline, Panel, Stack } from '../../layout';
import { FormRow } from './FormRow';

const meta = { title: 'Design System/Layout/FormRow', component: FormRow } satisfies Meta<typeof FormRow>;
export default meta;
type Story = StoryObj<typeof meta>;

const baseArgs = { label: 'Example', control: <TextInput defaultValue="example" /> };

export const Basic: Story = {
  args: {
    label: 'Limit',
    control: <TextInput defaultValue="10" />,
    hint: 'Maximum number of results.',
  },
};

export const States: Story = {
  args: baseArgs,
  render: () => (
    <Panel title="Form row states" density="condensed">
      <Stack gap="md">
        <FormRow
          label="Display name"
          control={<TextInput defaultValue="Manuel" />}
          description="Shown next to uploads from this browser."
          hint="Use a short recognizable name."
          counter="6 / 40"
        />
        <FormRow
          label="Required"
          required
          control={<TextInput placeholder="Required field" required />}
          hint="Required marker is visual; native validation stays on the control."
        />
        <FormRow
          label="Saved"
          control={<TextInput defaultValue="Manuel" />}
          success={<StatusText status="succeeded" icon>Display name saved.</StatusText>}
        />
        <FormRow
          label="Invalid"
          control={<TextInput defaultValue="Invalid name" aria-invalid="true" />}
          error={<StatusText status="error" icon>Display name contains unsupported characters.</StatusText>}
        />
        <FormRow
          label="Read-only"
          control={<TextInput defaultValue="sess_8f2c…" readOnly />}
          hint="System-managed value."
        />
        <FormRow
          label="Disabled"
          control={<TextInput defaultValue="Locked" disabled />}
          hint="Disabled controls use the same row anatomy."
        />
      </Stack>
    </Panel>
  ),
};

export const Stacked: Story = {
  args: baseArgs,
  render: () => (
    <Panel title="Stacked form rows" density="condensed">
      <Stack gap="md">
        <FormRow
          orientation="stacked"
          label="Long helper text"
          control={<TextInput defaultValue="Context workshop uploads" />}
          description="Stacked rows are useful for narrow panels or longer explanatory copy."
          hint="The helper text spans the full row instead of starting after a fixed label column."
          counter="24 / 80"
        />
        <FormRow
          orientation="stacked"
          label="Visibility"
          control={(
            <SelectInput defaultValue="public">
              <option value="public">Public</option>
              <option value="private">Private</option>
            </SelectInput>
          )}
          error="Choose carefully; this affects where the upload appears."
        />
      </Stack>
    </Panel>
  ),
};

export const MixedControls: Story = {
  args: baseArgs,
  render: () => (
    <Panel title="Mixed controls" density="condensed">
      <Stack gap="md">
        <FormRow
          label="Name"
          control={<TextInput defaultValue="Anonymous visitor" />}
          hint="Text input control."
        />
        <FormRow
          label="Scope"
          control={(
            <SelectInput defaultValue="mine">
              <option value="all">All uploads</option>
              <option value="mine">My uploads</option>
            </SelectInput>
          )}
          hint="Select control."
        />
        <FormRow
          label="Options"
          control={(
            <Stack gap="xs">
              <CheckboxRow checked={true} onChange={() => {}}>Show my display name</CheckboxRow>
              <CheckboxRow checked={false} onChange={() => {}}>Hide uploads by default</CheckboxRow>
            </Stack>
          )}
          hint="Multiple controls can live in the control slot."
        />
      </Stack>
    </Panel>
  ),
};

export const DenseSettingsSketch: Story = {
  args: baseArgs,
  render: () => (
    <Panel
      title="Settings anatomy sketch"
      density="condensed"
      actions={(
        <Inline gap="xs" wrap>
          <Button type="button" size="compact">Cancel</Button>
          <Button type="button" size="compact" variant="primary">Save</Button>
        </Inline>
      )}
    >
      <Stack gap="sm">
        <FormRow
          label="Browser identity"
          control={<TextInput defaultValue="Anonymous visitor" readOnly />}
          description="This is stored in a long-lived browser cookie, not an authenticated account."
        />
        <FormRow
          label="Display name"
          required
          control={<TextInput defaultValue="Manuel" maxLength={40} />}
          hint="Shown on future uploads from this browser."
          counter="6 / 40"
        />
        <FormRow
          label="Status"
          control={<StatusText status="succeeded" icon>Ready to save</StatusText>}
          success="All fields are valid."
        />
      </Stack>
    </Panel>
  ),
};
