import type { Meta, StoryObj } from "@storybook/react-vite";
import type { ReactNode } from "react";
import { TextInput } from "../atoms";
import { Caption } from "./Caption";
import { CodeText } from "./CodeText";
import { Divider } from "./Divider";
import { StatusText } from "./StatusText";
import { Text } from "./Text";
import { VisuallyHidden } from "./VisuallyHidden";
import styles from "./Foundation.stories.module.css";

const meta = {
	title: "Design System/Foundation/Overview",
	parameters: {
		layout: "padded",
	},
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

const colorTokens = [
	["Background", "--mac-bg"],
	["Text", "--mac-text"],
	["Text dim", "--mac-text-dim"],
	["Inverse text", "--mac-text-inv"],
	["Border", "--mac-border"],
	["Surface", "--mac-surface"],
	["Surface 2", "--mac-surface-2"],
	["Accent", "--mac-accent"],
	["Accent 2", "--mac-accent-2"],
	["Success green", "--mac-green"],
	["Warning amber", "--mac-amber"],
];

const spacingTokens = [
	["2px", 2],
	["4px", 4],
	["6px", 6],
	["8px", 8],
	["12px", 12],
	["16px", 16],
	["24px", 24],
	["32px", 32],
];

function Section({ title, children }: { title: string; children: ReactNode }) {
	return (
		<section className={styles.section}>
			<h3 className={styles.sectionTitle}>{title}</h3>
			{children}
		</section>
	);
}

export const Colors: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Monochrome dashboard palette">
				<div className={styles.swatchGrid}>
					{colorTokens.map(([label, token]) => (
						<div key={token} className={styles.swatch}>
							<div className={styles.swatchColor} style={{ background: `var(${token})` }} />
							<div className={styles.swatchBody}>
								<Caption tone="inherit">{label}</Caption>
								<CodeText>{token}</CodeText>
							</div>
						</div>
					))}
				</div>
			</Section>
		</div>
	),
};

export const Typography: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Type roles">
				<div className={styles.specimenStack}>
					<Text>Body text — used for readable dashboard copy and placeholder descriptions.</Text>
					<Text size="compact">
						Compact text — used in dense panels where vertical space is constrained.
					</Text>
					<Caption>Caption — metadata, hints, helper text, and small dashboard labels.</Caption>
					<Caption transform="uppercase">Uppercase caption — compact categorical labels.</Caption>
					<CodeText>CodeText — chunk_01 / openai/text-embedding-3-small / fixed-1200-150</CodeText>
				</div>
			</Section>
		</div>
	),
};

export const StatusTones: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Status tones">
				<div className={styles.statusGrid}>
					{[
						"pending",
						"ready",
						"running",
						"succeeded",
						"done",
						"partial",
						"warning",
						"failed",
						"error",
						"canceled",
					].map((status) => (
						<StatusText key={status} status={status} icon>
							{status}
						</StatusText>
					))}
				</div>
			</Section>
		</div>
	),
};

export const Spacing: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Dashboard spacing examples">
				<div className={styles.specimenStack}>
					{spacingTokens.map(([label, width]) => (
						<div key={label} className={styles.spacingRow}>
							<CodeText>{label}</CodeText>
							<div className={styles.spacingBar} style={{ width }} />
						</div>
					))}
				</div>
			</Section>
		</div>
	),
};

export const BordersAndRadii: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Retro borders and surfaces">
				<div className={styles.radiusGrid}>
					<div>
						<div className={styles.radiusSample} />
						<Caption>Default one-pixel panel border</Caption>
					</div>
					<div>
						<div className={`${styles.radiusSample} ${styles.borderStrong}`} />
						<Caption>Strong border for emphasis/active frames</Caption>
					</div>
					<div>
						<div className={styles.radiusSample} style={{ background: "var(--mac-surface-2)" }} />
						<Caption>Secondary surface fill</Caption>
					</div>
				</div>
			</Section>
			<Divider />
			<Caption>
				RAG intentionally uses square Mac-style surfaces; add radius tokens only if the visual
				language changes.
			</Caption>
		</div>
	),
};

export const Accessibility: Story = {
	render: () => (
		<div className={styles.page}>
			<Section title="Accessibility primitives">
				<div className={styles.accessibilityDemo}>
					<Text>
						The next label is visually hidden but remains available to assistive technology.
					</Text>
					<label>
						<VisuallyHidden>Search query</VisuallyHidden>
						<TextInput placeholder="Visible input placeholder" />
					</label>
					<Caption>Use VisuallyHidden instead of page-local sr-only CSS.</Caption>
				</div>
			</Section>
		</div>
	),
};
