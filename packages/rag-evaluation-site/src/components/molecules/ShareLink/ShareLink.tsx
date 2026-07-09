import type { HTMLAttributes, ReactNode } from "react";
import { Button } from "../../atoms";
import { Caption, CodeText, Text } from "../../foundation";
import styles from "./ShareLink.module.css";

export interface ShareLinkProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	label?: ReactNode;
	description?: ReactNode;
	href: string;
	displayHref?: ReactNode;
	copyLabel?: ReactNode;
	copiedLabel?: ReactNode;
	copied?: boolean;
	onCopy?: () => void;
}

export function ShareLink({
	label = "Share link",
	description,
	href,
	displayHref,
	copyLabel = "Copy link",
	copiedLabel = "Copied",
	copied = false,
	onCopy,
	className,
	...rest
}: ShareLinkProps) {
	const visibleHref = displayHref ?? href;
	const currentCopyLabel = copied ? copiedLabel : copyLabel;
	const copyButtonLabel =
		typeof currentCopyLabel === "string" ? currentCopyLabel : copied ? "Copied" : "Copy link";
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="ShareLink"
			{...rest}
		>
			<div className={styles.header}>
				<Text size="label" tone="muted">
					{label}
				</Text>
				{description ? <Caption>{description}</Caption> : null}
			</div>
			<div className={styles.body}>
				<a
					className={styles.link}
					href={href}
					aria-label={typeof label === "string" ? label : "Share link"}
				>
					<CodeText as="span" display="block" tone="accent">
						{visibleHref}
					</CodeText>
				</a>
				{onCopy ? (
					<Button
						size="compact"
						className={styles.copyButton}
						onClick={onCopy}
						aria-label={copyButtonLabel}
						title={copyButtonLabel}
					>
						{copied ? (
							<span className={styles.copyIcon} aria-hidden="true">
								✓
							</span>
						) : (
							<CopyIcon />
						)}
					</Button>
				) : null}
			</div>
		</div>
	);
}

function CopyIcon() {
	return (
		<svg
			className={styles.copyIcon}
			viewBox="0 0 16 16"
			width="14"
			height="14"
			aria-hidden="true"
			focusable="false"
		>
			<path d="M6 2.5h6.5v8H11v-6H6z" fill="none" stroke="currentColor" strokeWidth="1.25" />
			<path d="M3.5 5.5H10v8H3.5z" fill="none" stroke="currentColor" strokeWidth="1.25" />
		</svg>
	);
}
