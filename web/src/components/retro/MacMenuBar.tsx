import React from "react";

interface MacMenuBarProps {
	activeView: string;
	onViewChange: (view: string) => void;
}

const views = [
	{ id: "pipeline", label: "Pipeline" },
	{ id: "embeddings", label: "Embeddings" },
	{ id: "search", label: "Search" },
	{ id: "evaluation", label: "Evaluation" },
];

export const MacMenuBar: React.FC<MacMenuBarProps> = ({ activeView, onViewChange }) => {
	return (
		<div className="mac-menu-bar">
			<span style={{ fontWeight: "bold", marginRight: 8 }}>◉ RAG Eval</span>
			{views.map((view) => (
				<span
					key={view.id}
					className={`mac-menu-item ${activeView === view.id ? "underline" : ""}`}
					onClick={() => onViewChange(view.id)}
				>
					{view.label}
				</span>
			))}
		</div>
	);
};
