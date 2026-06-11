import type { ReactNode } from 'react';
import { Button } from '../components/atoms';
import { Caption, StatusText } from '../components/foundation';
import type { ActionSpec, CellSpec, JsonObject, RenderableValue, WidgetNode } from './ir';
import { isWidgetNode } from './ir';

export type RenderWidgetNode = (node: WidgetNode) => ReactNode;

export function renderCell(spec: CellSpec, row: JsonObject, renderWidgetNode: RenderWidgetNode, dispatchAction?: (action: ActionSpec, context: JsonObject) => void): ReactNode {
  switch (spec.kind) {
    case 'field':
      return stringify(getPath(row, spec.field), spec.fallback);
    case 'number':
      return formatNumber(getPath(row, spec.field), spec);
    case 'status': {
      const status = stringify(getPath(row, spec.field), spec.fallback ?? 'pending');
      return <StatusText status={status} icon={spec.icon}>{status}</StatusText>;
    }
    case 'caption':
      return <Caption tone={spec.tone}>{stringify(getPath(row, spec.field), spec.fallback)}</Caption>;
    case 'template':
      return renderTemplate(spec.template, row);
    case 'link': {
      const href = stringify(getPath(row, spec.hrefField), '#');
      const label = stringify(getPath(row, spec.labelField), spec.fallbackLabel ?? href);
      return <a href={href} target={spec.target} rel={spec.target === '_blank' ? 'noreferrer' : undefined}>{label}</a>;
    }
    case 'linkButton': {
      const href = stringify(getPath(row, spec.hrefField), '#');
      const label = stringify(getPath(row, spec.labelField), spec.fallbackLabel ?? href);
      return (
        <Button
          variant={spec.variant}
          size={spec.size ?? 'compact'}
          onClick={(event) => {
            event.stopPropagation();
            if (!href || href === '#') return;
            window.history.pushState({}, '', href);
            window.dispatchEvent(new PopStateEvent('popstate'));
          }}
        >
          {label}
        </Button>
      );
    }
    case 'actionButton':
      return (
        <Button
          variant={spec.variant}
          size={spec.size ?? 'compact'}
          disabled={spec.disabled}
          onClick={(event) => {
            event.stopPropagation();
            dispatchAction?.(spec.action, { row, rowKey: rowKey(row, 'file'), componentType: 'DataTableCell' });
          }}
        >
          {renderRenderable(spec.label, renderWidgetNode)}
        </Button>
      );
    case 'constant':
      return renderRenderable(spec.value, renderWidgetNode);
    default:
      return <Caption tone="danger">Unsupported cell</Caption>;
  }
}

export function renderRenderable(value: RenderableValue | undefined, renderWidgetNode: RenderWidgetNode): ReactNode {
  if (value === undefined || value === null || value === false) return null;
  if (isWidgetNode(value)) return renderWidgetNode(value);
  return String(value);
}

export function rowKey(row: JsonObject, spec: string | { field: string } | { template: string }): string {
  if (typeof spec === 'string') return stringify(getPath(row, spec), '');
  if ('field' in spec) return stringify(getPath(row, spec.field), '');
  return renderTemplate(spec.template, row);
}

function getPath(row: JsonObject, path: string): unknown {
  const parts = path.split('.').filter(Boolean);
  let current: unknown = row;
  for (const part of parts) {
    if (!current || typeof current !== 'object') return undefined;
    current = (current as Record<string, unknown>)[part];
  }
  return current;
}

function stringify(value: unknown, fallback = ''): string {
  if (value === undefined || value === null) return fallback;
  return String(value);
}

function formatNumber(value: unknown, spec: Extract<CellSpec, { kind: 'number' }>): string {
  const n = typeof value === 'number' ? value : Number(value);
  if (!Number.isFinite(n)) return spec.fallback ?? '';
  if (spec.format === 'fixed') return n.toFixed(spec.digits ?? 2);
  return Math.round(n).toLocaleString();
}

function renderTemplate(template: string, row: JsonObject): string {
  return template.replace(/\$\{([^}]+)\}|\$([A-Za-z0-9_.-]+)/g, (_match, braced: string | undefined, bare: string | undefined) => {
    const path = braced ?? bare ?? '';
    return stringify(getPath(row, path), '');
  });
}
