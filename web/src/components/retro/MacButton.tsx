import React from "react";

interface MacButtonProps {
	label: string;
	onClick: () => void;
	primary?: boolean;
	disabled?: boolean;
}

export const MacButton: React.FC<MacButtonProps> = ({
	label,
	onClick,
	primary = false,
	disabled = false,
}) => {
	return (
		<button
			className={`mac-btn ${primary ? "mac-btn-primary" : ""}`}
			onClick={onClick}
			disabled={disabled}
		>
			{label}
		</button>
	);
};
