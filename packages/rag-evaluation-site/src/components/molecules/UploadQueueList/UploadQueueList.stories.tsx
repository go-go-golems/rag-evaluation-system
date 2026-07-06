import type { Meta, StoryObj } from "@storybook/react-vite";
import { useEffect, useState } from "react";
import type { UploadQueueItem } from "../../../cms";
import { cmsUploadQueueFixtures } from "../../../cms";
import { UploadQueueList } from "./UploadQueueList";

const meta = {
	title: "Component Library/Molecules/UploadQueueList",
	component: UploadQueueList,
	args: {
		items: cmsUploadQueueFixtures,
		onCancel: () => {},
		onRetry: () => {},
		onDismiss: () => {},
		style: { maxWidth: 560 },
	},
} satisfies Meta<typeof UploadQueueList>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Mixed: Story = {};

export const AllUploading: Story = {
	args: {
		items: cmsUploadQueueFixtures.map((item, index) => ({
			...item,
			status: "uploading",
			progress: 0.15 + index * 0.18,
			error: undefined,
		})),
	},
};

export const WithErrors: Story = {
	args: {
		items: cmsUploadQueueFixtures.filter((item) => item.status === "error"),
	},
};

export const Empty: Story = {
	args: { items: [] },
};

export const Interactive: Story = {
	render: () => {
		const [items, setItems] = useState<UploadQueueItem[]>(
			cmsUploadQueueFixtures.map((item) => ({ ...item })),
		);
		useEffect(() => {
			const timer = setInterval(() => {
				setItems((current) =>
					current.map((item) =>
						item.status === "uploading"
							? {
									...item,
									progress: Math.min(1, item.progress + 0.08),
									status: item.progress >= 0.92 ? "done" : "uploading",
								}
							: item,
					),
				);
			}, 400);
			return () => clearInterval(timer);
		}, []);
		return (
			<UploadQueueList
				style={{ maxWidth: 560 }}
				items={items}
				onCancel={(id) =>
					setItems((current) =>
						current.map((item) => (item.id === id ? { ...item, status: "canceled" } : item)),
					)
				}
				onRetry={(id) =>
					setItems((current) =>
						current.map((item) =>
							item.id === id
								? { ...item, status: "uploading", progress: 0, error: undefined }
								: item,
						),
					)
				}
				onDismiss={(id) => setItems((current) => current.filter((item) => item.id !== id))}
			/>
		);
	},
};
