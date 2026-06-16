/**
 * CSS extraction entry for standalone design-system references.
 *
 * Important: CSS Modules imported only for side effects can be removed from the
 * Vite/Rollup output. Keep bindings and export them so Rollup treats the module
 * objects as used and Vite emits the hashed CSS rules that match SSR output.
 */

import "../index.css";

import button from "../../../packages/rag-evaluation-site/src/components/atoms/Button/Button.module.css";
import checkboxRow from "../../../packages/rag-evaluation-site/src/components/atoms/CheckboxRow/CheckboxRow.module.css";
import errorCallout from "../../../packages/rag-evaluation-site/src/components/atoms/ErrorCallout/ErrorCallout.module.css";
import iconButton from "../../../packages/rag-evaluation-site/src/components/atoms/IconButton/IconButton.module.css";
import selectInput from "../../../packages/rag-evaluation-site/src/components/atoms/SelectInput/SelectInput.module.css";
import textInput from "../../../packages/rag-evaluation-site/src/components/atoms/TextInput/TextInput.module.css";

import caption from "../../../packages/rag-evaluation-site/src/components/foundation/Caption/Caption.module.css";
import codeText from "../../../packages/rag-evaluation-site/src/components/foundation/CodeText/CodeText.module.css";
import divider from "../../../packages/rag-evaluation-site/src/components/foundation/Divider/Divider.module.css";
import statusText from "../../../packages/rag-evaluation-site/src/components/foundation/StatusText/StatusText.module.css";
import text from "../../../packages/rag-evaluation-site/src/components/foundation/Text/Text.module.css";
import visuallyHidden from "../../../packages/rag-evaluation-site/src/components/foundation/VisuallyHidden/VisuallyHidden.module.css";

import appShell from "../../../packages/rag-evaluation-site/src/components/layout/AppShell/AppShell.module.css";
import dashboardGrid from "../../../packages/rag-evaluation-site/src/components/layout/DashboardGrid/DashboardGrid.module.css";
import formRow from "../../../packages/rag-evaluation-site/src/components/layout/FormRow/FormRow.module.css";
import inline from "../../../packages/rag-evaluation-site/src/components/layout/Inline/Inline.module.css";
import panel from "../../../packages/rag-evaluation-site/src/components/layout/Panel/Panel.module.css";
import scrollRegion from "../../../packages/rag-evaluation-site/src/components/layout/ScrollRegion/ScrollRegion.module.css";
import stack from "../../../packages/rag-evaluation-site/src/components/layout/Stack/Stack.module.css";
import tabList from "../../../packages/rag-evaluation-site/src/components/layout/TabList/TabList.module.css";

import appNav from "../../../packages/rag-evaluation-site/src/components/molecules/AppNav/AppNav.module.css";
import dataTable from "../../../packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.module.css";
import metadataGrid from "../../../packages/rag-evaluation-site/src/components/molecules/MetadataGrid/MetadataGrid.module.css";

const cssModules = {
	button,
	checkboxRow,
	errorCallout,
	iconButton,
	selectInput,
	textInput,
	caption,
	codeText,
	divider,
	statusText,
	text,
	visuallyHidden,
	appShell,
	dashboardGrid,
	formRow,
	inline,
	panel,
	scrollRegion,
	stack,
	tabList,
	appNav,
	dataTable,
	metadataGrid,
};

// Runtime side effect: prevents Rollup from deleting all CSS Module imports.
(
	globalThis as typeof globalThis & { __ragDesignReferenceCssModules?: unknown }
).__ragDesignReferenceCssModules = cssModules;

export { cssModules };
