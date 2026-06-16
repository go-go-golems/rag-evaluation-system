import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RagEvaluationSiteApp } from "./App";
import "../styles.css";

const rootElement = document.getElementById("root");

if (!rootElement) {
	throw new Error('RagEvaluationSiteApp requires an element with id="root"');
}

createRoot(rootElement).render(
	<StrictMode>
		<RagEvaluationSiteApp />
	</StrictMode>,
);
