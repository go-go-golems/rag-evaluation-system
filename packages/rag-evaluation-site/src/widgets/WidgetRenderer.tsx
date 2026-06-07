import { createElement, type CSSProperties, type ReactNode } from 'react';
import { Button, ErrorCallout, SelectInput, TextInput } from '../components/atoms';
import { Caption, StatusText } from '../components/foundation';
import { AppShell, DashboardGrid, FormRow, Inline, Panel, ScrollRegion, Stack, TabList } from '../components/layout';
import { AppNav, DataTable, MetadataGrid } from '../components/molecules';
import type {
  ActionSpec,
  AppNavWidgetProps,
  AppShellWidgetProps,
  ButtonWidgetProps,
  CaptionWidgetProps,
  ComponentNode,
  DashboardGridWidgetProps,
  DataTableWidgetProps,
  ElementNode,
  FormRowWidgetProps,
  InlineWidgetProps,
  JsonObject,
  MetadataGridWidgetProps,
  PanelWidgetProps,
  ScrollRegionWidgetProps,
  SelectInputWidgetProps,
  StackWidgetProps,
  StatusTextWidgetProps,
  TabListWidgetProps,
  TextInputWidgetProps,
  WidgetNode,
} from './ir';
import { bindAction, type WidgetActionContext, type WidgetActionHandler } from './actions';
import { renderCell, renderRenderable, rowKey } from './cellRenderers';

export interface WidgetRendererProps {
  node: WidgetNode;
  onAction?: WidgetActionHandler;
}

export function WidgetRenderer({ node, onAction }: WidgetRendererProps) {
  return <>{renderWidgetNode(node, onAction)}</>;
}

function renderWidgetNode(node: WidgetNode, onAction?: WidgetActionHandler): ReactNode {
  if (node.kind === 'text') return node.text;
  if (node.kind === 'element') return renderElementNode(node, onAction);
  return renderComponentNode(node, onAction);
}

function renderChildren(children: WidgetNode[] | undefined, onAction?: WidgetActionHandler): ReactNode[] {
  return (children ?? []).map((child, index) => <WidgetRenderer key={index} node={child} onAction={onAction} />);
}

function renderRenderableValue(value: unknown, onAction?: WidgetActionHandler): ReactNode {
  return renderRenderable(value as never, (node) => renderWidgetNode(node, onAction));
}

function renderElementNode(node: ElementNode, onAction?: WidgetActionHandler): ReactNode {
  const attrs = { ...(node.attrs ?? {}) } as Record<string, unknown>;
  if (attrs.style && typeof attrs.style === 'object') attrs.style = attrs.style as CSSProperties;
  return createElement(node.tag, attrs, ...renderChildren(node.children, onAction));
}

function renderComponentNode(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  switch (node.type) {
    case 'AppShell':
      return renderAppShell(node, onAction);
    case 'AppNav':
      return renderAppNav(node, onAction);
    case 'Button':
      return renderButton(node, onAction);
    case 'Caption':
      return renderCaption(node, onAction);
    case 'DashboardGrid':
      return renderDashboardGrid(node, onAction);
    case 'DataTable':
      return renderDataTable(node, onAction);
    case 'FormRow':
      return renderFormRow(node, onAction);
    case 'Inline':
      return renderInline(node, onAction);
    case 'MetadataGrid':
      return renderMetadataGrid(node, onAction);
    case 'Panel':
      return renderPanel(node, onAction);
    case 'ScrollRegion':
      return renderScrollRegion(node, onAction);
    case 'SelectInput':
      return renderSelectInput(node, onAction);
    case 'Stack':
      return renderStack(node, onAction);
    case 'StatusText':
      return renderStatusText(node, onAction);
    case 'TabList':
      return renderTabList(node, onAction);
    case 'TextInput':
      return renderTextInput(node);
    default:
      return <ErrorCallout>Unknown widget: {node.type}</ErrorCallout>;
  }
}

function renderAppShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AppShellWidgetProps;
  return (
    <AppShell
      className={props.className}
      header={props.header ? renderWidgetNode(props.header, onAction) : undefined}
      sidebar={props.sidebar ? renderWidgetNode(props.sidebar, onAction) : undefined}
    >
      {renderChildren(node.children, onAction)}
    </AppShell>
  );
}

function renderAppNav(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AppNavWidgetProps;
  return (
    <AppNav
      brand={renderRenderableValue(props.brand, onAction)}
      items={props.items.map((item) => ({ id: item.id, label: renderRenderableValue(item.label, onAction) }))}
      activeItemId={props.activeItemId}
      onItemSelect={(itemId) => {
        const item = props.items.find((candidate) => candidate.id === itemId);
        if (item?.action) bindAndRun(item.action, { value: itemId, componentType: 'AppNav' }, onAction);
      }}
    />
  );
}

function renderButton(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as ButtonWidgetProps;
  return (
    <Button
      className={props.className}
      variant={props.variant}
      size={props.size}
      selected={props.selected}
      disabled={props.disabled}
      type={props.type ?? 'button'}
      onClick={bindAction(props.action, { componentType: 'Button' }, onAction)}
    >
      {renderChildren(node.children, onAction)}
    </Button>
  );
}

function renderCaption(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CaptionWidgetProps;
  return <Caption className={props.className} tone={props.tone} transform={props.transform} truncate={props.truncate}>{renderChildren(node.children, onAction)}</Caption>;
}

function renderDashboardGrid(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DashboardGridWidgetProps;
  return <DashboardGrid className={props.className} recipe={props.recipe}>{renderChildren(node.children, onAction)}</DashboardGrid>;
}

function renderDataTable(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DataTableWidgetProps;
  const rowSelectAction = props.onRowSelect;
  return (
    <DataTable<JsonObject>
      className={props.className}
      rows={props.rows}
      columns={props.columns.map((column) => ({
        id: column.id,
        header: renderRenderableValue(column.header, onAction),
        align: column.align,
        maxWidth: column.maxWidth,
        cell: (row) => renderCell(column.cell, row, (child) => renderWidgetNode(child, onAction)),
      }))}
      getRowKey={(row) => rowKey(row, props.getRowKey)}
      selectedKey={props.selectedKey == null ? props.selectedKey : String(props.selectedKey)}
      emptyMessage={renderRenderableValue(props.emptyMessage, onAction)}
      onRowSelect={rowSelectAction ? (row) => bindAndRun(rowSelectAction, { row, rowKey: rowKey(row, props.getRowKey), componentType: 'DataTable' }, onAction) : undefined}
    />
  );
}

function renderFormRow(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as FormRowWidgetProps;
  return (
    <FormRow
      className={props.className}
      label={renderRenderableValue(props.label, onAction)}
      control={renderWidgetNode(props.control, onAction)}
      hint={renderRenderableValue(props.hint, onAction)}
    />
  );
}

function renderInline(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as InlineWidgetProps;
  return <Inline className={props.className} gap={props.gap} justify={props.justify} wrap={props.wrap}>{renderChildren(node.children, onAction)}</Inline>;
}

function renderMetadataGrid(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as MetadataGridWidgetProps;
  return (
    <MetadataGrid
      className={props.className}
      density={props.density}
      items={props.items.map((item) => ({
        key: renderRenderableValue(item.key, onAction),
        value: renderRenderableValue(item.value, onAction),
        copyValue: item.copyValue,
      }))}
    />
  );
}

function renderPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as PanelWidgetProps;
  return (
    <Panel
      className={props.className}
      title={renderRenderableValue(props.title, onAction)}
      actions={renderRenderableValue(props.actions, onAction)}
      density={props.density}
      fill={props.fill}
    >
      {renderChildren(node.children, onAction)}
    </Panel>
  );
}

function renderScrollRegion(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as ScrollRegionWidgetProps;
  return <ScrollRegion className={props.className} axis={props.axis} style={props.style as CSSProperties | undefined}>{renderChildren(node.children, onAction)}</ScrollRegion>;
}

function renderSelectInput(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SelectInputWidgetProps;
  return (
    <SelectInput className={props.className} name={props.name} value={props.value} disabled={props.disabled}>
      {(props.options ?? []).map((option) => (
        <option key={String(option.value)} value={option.value} disabled={option.disabled}>
          {renderRenderableValue(option.label, onAction)}
        </option>
      ))}
      {renderChildren(node.children, onAction)}
    </SelectInput>
  );
}

function renderStack(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as StackWidgetProps;
  return <Stack className={props.className} gap={props.gap} align={props.align}>{renderChildren(node.children, onAction)}</Stack>;
}

function renderStatusText(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as StatusTextWidgetProps;
  return <StatusText className={props.className} status={props.status} icon={props.icon}>{renderChildren(node.children, onAction)}</StatusText>;
}

function renderTabList(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TabListWidgetProps;
  return (
    <TabList
      items={props.items.map((item) => ({ id: item.id, label: renderRenderableValue(item.label, onAction) }))}
      activeId={props.activeId}
      ariaLabel={props.ariaLabel}
      onChange={(id) => {
        if (props.onChange) bindAndRun(props.onChange, { value: id, componentType: 'TabList' }, onAction);
      }}
    />
  );
}

function renderTextInput(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as TextInputWidgetProps;
  return (
    <TextInput
      className={props.className}
      name={props.name}
      value={props.value}
      placeholder={props.placeholder}
      type={props.type}
      disabled={props.disabled}
      min={props.min}
      max={props.max}
      readOnly
    />
  );
}

function bindAndRun(action: ActionSpec, context: WidgetActionContext, onAction?: WidgetActionHandler): void {
  bindAction(action, context, onAction)?.();
}
