/**
 * CSS extraction entry for standalone design-system references.
 *
 * Important: CSS Modules imported only for side effects can be removed from the
 * Vite/Rollup output. Keep bindings and export them so Rollup treats the module
 * objects as used and Vite emits the hashed CSS rules that match SSR output.
 */

import '../index.css';

import button from '../components/atoms/Button/Button.module.css';
import checkboxRow from '../components/atoms/CheckboxRow/CheckboxRow.module.css';
import errorCallout from '../components/atoms/ErrorCallout/ErrorCallout.module.css';
import iconButton from '../components/atoms/IconButton/IconButton.module.css';
import selectInput from '../components/atoms/SelectInput/SelectInput.module.css';
import textInput from '../components/atoms/TextInput/TextInput.module.css';

import caption from '../components/foundation/Caption/Caption.module.css';
import codeText from '../components/foundation/CodeText/CodeText.module.css';
import divider from '../components/foundation/Divider/Divider.module.css';
import statusText from '../components/foundation/StatusText/StatusText.module.css';
import text from '../components/foundation/Text/Text.module.css';
import visuallyHidden from '../components/foundation/VisuallyHidden/VisuallyHidden.module.css';

import appShell from '../components/layout/AppShell/AppShell.module.css';
import dashboardGrid from '../components/layout/DashboardGrid/DashboardGrid.module.css';
import formRow from '../components/layout/FormRow/FormRow.module.css';
import inline from '../components/layout/Inline/Inline.module.css';
import panel from '../components/layout/Panel/Panel.module.css';
import scrollRegion from '../components/layout/ScrollRegion/ScrollRegion.module.css';
import stack from '../components/layout/Stack/Stack.module.css';
import tabList from '../components/layout/TabList/TabList.module.css';

import appNav from '../components/molecules/AppNav/AppNav.module.css';
import dataTable from '../components/molecules/DataTable/DataTable.module.css';
import metadataGrid from '../components/molecules/MetadataGrid/MetadataGrid.module.css';

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
(globalThis as typeof globalThis & { __ragDesignReferenceCssModules?: unknown }).__ragDesignReferenceCssModules = cssModules;

export { cssModules };
