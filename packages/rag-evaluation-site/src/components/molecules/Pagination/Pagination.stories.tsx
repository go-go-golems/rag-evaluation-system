import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Pagination } from "./Pagination";

const meta = {
	title: "Component Library/Molecules/Pagination",
	component: Pagination,
	args: { page: 3, pageCount: 9, onPageChange: () => {} },
} satisfies Meta<typeof Pagination>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const FirstPage: Story = {
	args: { page: 1 },
};

export const LastPage: Story = {
	args: { page: 9 },
};

export const WithTotals: Story = {
	args: { page: 1, pageCount: 9, pageSize: 24, totalItems: 210 },
};

export const SinglePage: Story = {
	args: { page: 1, pageCount: 1 },
};

export const Interactive: Story = {
	render: () => {
		const [page, setPage] = useState(1);
		return (
			<Pagination page={page} pageCount={9} pageSize={24} totalItems={210} onPageChange={setPage} />
		);
	},
};
