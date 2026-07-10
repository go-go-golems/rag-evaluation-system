import type { HTMLAttributes, ReactNode } from "react";
import { Caption, Text } from "../../foundation";
import { Panel, SplitPane, Stack } from "../../layout";
import styles from "./RecordShell.module.css";

export interface RecordShellIdentity {
	name: ReactNode;
	subtitle?: ReactNode;
	/** Avatar image; falls back to `avatarText` (usually initials). */
	avatarUrl?: string;
	avatarText?: string;
}

export interface RecordShellProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	identity: RecordShellIdentity;
	/** Header action bar (buttons/menus). */
	actions?: ReactNode;
	/** Left column — the record's fields (usually a RecordFieldList). */
	details: ReactNode;
	detailsTitle?: ReactNode;
	/** Right column top — the timeline (usually an ActivityFeed). */
	activity?: ReactNode;
	activityTitle?: ReactNode;
	activityActions?: ReactNode;
	/** Right column below the timeline — related-list panels. */
	related?: ReactNode;
}

/**
 * The record page shell shared by every object type (contact, company, deal): a
 * header (identity + actions), a left column of fields, and a right column of a
 * timeline and related lists. Presentational — presets (contactRecord,
 * dealRecord) supply which fields and related lists appear. Composes the existing
 * layout primitives (Panel, SplitPane, Stack) rather than inventing new layout.
 */
export function RecordShell({
	identity,
	actions,
	details,
	detailsTitle = "Details",
	activity,
	activityTitle = "Activity",
	activityActions,
	related,
	className,
	...rest
}: RecordShellProps) {
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="RecordShell"
			{...rest}
		>
			<header className={styles.header}>
				<span className={styles.avatar} aria-hidden="true">
					{identity.avatarUrl ? (
						<img src={identity.avatarUrl} alt="" />
					) : (
						(identity.avatarText ?? "?")
					)}
				</span>
				<div className={styles.identity}>
					<Text as="div" size="body" weight="bold">
						{identity.name}
					</Text>
					{identity.subtitle != null ? <Caption>{identity.subtitle}</Caption> : null}
				</div>
				{actions != null ? <div className={styles.actions}>{actions}</div> : null}
			</header>

			<SplitPane
				ratio="leftNarrow"
				gutter="lg"
				left={
					<Panel title={detailsTitle} density="condensed">
						{details}
					</Panel>
				}
				right={
					<Stack gap="md">
						{activity != null ? (
							<Panel title={activityTitle} density="condensed" actions={activityActions}>
								{activity}
							</Panel>
						) : null}
						{related}
					</Stack>
				}
			/>
		</div>
	);
}
