import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { CmsAsset } from "../../../cms";
import { cmsAssetFixtures } from "../../../cms";
import { Button } from "../../atoms";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { AssetPickerDialog } from "./AssetPickerDialog";

const meta = {
	title: "Component Library/Organisms/AssetPickerDialog",
	component: AssetPickerDialog,
	args: {
		open: true,
		mode: "inline",
		assets: cmsAssetFixtures.filter((asset) => asset.kind === "image"),
		onConfirm: () => {},
		onCancel: () => {},
	},
} satisfies Meta<typeof AssetPickerDialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Static: Story = {};

export const WithSearch: Story = {
	args: {
		query: "",
		onQueryChange: () => {},
		kindFilter: "image",
		onKindFilterChange: () => {},
	},
};

export const Interactive: Story = {
	render: () => {
		const [open, setOpen] = useState(false);
		const [picked, setPicked] = useState<CmsAsset>();
		return (
			<Stack gap="sm">
				<Button onClick={() => setOpen(true)}>Insert image…</Button>
				<Caption>picked: {picked ? picked.filename : "(none)"}</Caption>
				<AssetPickerDialog
					open={open}
					assets={cmsAssetFixtures.filter((asset) => asset.kind === "image")}
					onConfirm={(asset) => {
						setPicked(asset);
						setOpen(false);
					}}
					onCancel={() => setOpen(false)}
				/>
			</Stack>
		);
	},
};
