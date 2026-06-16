import React from "react";

interface MacWindowProps {
	title: string;
	children: React.ReactNode;
	className?: string;
	onClose?: () => void;
}

export const MacWindow: React.FC<MacWindowProps> = ({
	title,
	children,
	className = "",
	onClose,
}) => {
	return (
		<div className={`mac-window ${className}`}>
			<div className="mac-title-bar">
				{onClose && (
					<div className="mac-window-button" onClick={onClose} title="Close">
						■
					</div>
				)}
				<span className="mac-title-bar-text">{title}</span>
			</div>
			<div className="p-2" style={{ background: "var(--mac-window-white)" }}>
				{children}
			</div>
		</div>
	);
};
