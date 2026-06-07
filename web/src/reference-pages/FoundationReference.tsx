/**
 * FoundationReference.tsx
 *
 * Renders every design token as a visual specimen:
 *   - Colour palette swatches
 *   - Typography role specimens
 *   - Spacing scale
 *   - Border / shadow / radius rules
 *   - Status indicator colours
 *
 * This component is rendered server-side and the output is
 * bundled into a self-contained HTML file.
 */

import { StatusText } from '../components/foundation';
import { Caption } from '../components/foundation';
import { Button } from '../components/atoms';

/* ── colour palette ────────────────────────────────────────── */

const colours = [
  { token: '--mac-bg',        hex: '#FFFFFF', label: 'Page background' },
  { token: '--mac-bg-dark',   hex: '#000000', label: 'Dark bg / stripe' },
  { token: '--mac-text',      hex: '#000000', label: 'Primary text' },
  { token: '--mac-text-dim',  hex: '#666666', label: 'Muted text' },
  { token: '--mac-surface',   hex: '#FFFFFF', label: 'Panel surface' },
  { token: '--mac-surface-2', hex: '#EEEEEE', label: 'Alt surface' },
  { token: '--mac-accent',    hex: '#0000CC', label: 'Primary action' },
  { token: '--mac-accent-2',  hex: '#CC0000', label: 'Destructive / error' },
  { token: '--mac-green',     hex: '#007722', label: 'Success' },
  { token: '--mac-amber',    hex: '#AA7700', label: 'Warning / pending' },
];

/* ── typography ─────────────────────────────────────────────── */

const fontRoles = [
  { token: '--rag-font-role-body',     sample: 'The quick brown fox jumps over the lazy dog.' },
  { token: '--rag-font-role-compact',  sample: 'Compact secondary or toolbar text' },
  { token: '--rag-font-role-metadata', sample: 'key: value  —  400 11px Monaco' },
  { token: '--rag-font-role-label',    sample: 'SECTION HEADER — 700 11px Monaco' },
  { token: '--rag-font-role-metric',   sample: '2,847' },
  { token: '--rag-font-role-code',     sample: 'const x = await fetch("/api/data")' },
];

/* ── spacing ─────────────────────────────────────────────────── */

const spacing = [
  { name: 'xxs', px: 2 },
  { name: 'xs',  px: 4 },
  { name: 'sm',  px: 8 },
  { name: 'md',  px: 16 },
  { name: 'lg',  px: 24 },
];

/* ── status states ───────────────────────────────────────────── */

const statusStates: Array<{ status: string; label: string }> = [
  { status: 'succeeded', label: 'Succeeded' },
  { status: 'ok',        label: 'OK / Done' },
  { status: 'running',   label: 'Running' },
  { status: 'pending',   label: 'Pending' },
  { status: 'warning',   label: 'Warning' },
  { status: 'failed',    label: 'Failed' },
  { status: 'error',     label: 'Error' },
  { status: 'canceled',  label: 'Canceled' },
];

/* ── caption tones ───────────────────────────────────────────── */

const captionTones = ['muted', 'accent', 'success', 'warning', 'danger', 'inherit'] as const;

/* ── button variants ────────────────────────────────────────── */

/* ── component ──────────────────────────────────────────────── */

export default function FoundationReference() {
  return (
    <>
      <h1 style={{ borderBottom: '3px solid #000', paddingBottom: 8, marginBottom: '0.5em' }}>
        Foundation &amp; Tokens
      </h1>
      <p style={{ color: '#666', marginBottom: '2rem' }}>
        Self-contained visual reference. Every colour, font, spacing value, and
        component state is rendered from the real React components and CSS tokens.
      </p>

      {/* ── 1. Colour Palette ────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">1. Colour Palette</h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))', gap: 8 }}>
          {colours.map(c => (
            <div key={c.token} style={{ border: '1px solid #ccc', padding: 8 }}>
              <div style={{ height: 48, background: c.hex, border: c.hex === '#FFFFFF' ? '1px solid #ccc' : 'none' }} />
              <div style={{ fontFamily: 'Monaco, monospace', fontSize: 11, marginTop: 4 }}>
                <code>{c.token}</code>
              </div>
              <div style={{ fontSize: 12, color: '#666' }}>{c.hex} — {c.label}</div>
            </div>
          ))}
        </div>
      </div>

      {/* ── 2. Typography ─────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">2. Typography</h2>
        {fontRoles.map(role => (
          <div key={role.token} style={{ margin: '0.75em 0', borderBottom: '1px solid #eee', paddingBottom: '0.5em' }}>
            <div className="ds-ref-label"><code>{role.token}</code></div>
            <div style={{ font: `var(${role.token})` }}>{role.sample}</div>
          </div>
        ))}
      </div>

      {/* ── 3. Spacing ────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">3. Spacing Scale</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {spacing.map(s => (
            <div key={s.name} style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <div style={{ width: 80, fontFamily: 'Monaco, monospace', fontSize: 11 }}>{s.name}</div>
              <div style={{ width: s.px, height: 16, background: '#000' }} />
              <span style={{ fontFamily: 'Monaco, monospace', fontSize: 11, color: '#666' }}>{s.px}px</span>
            </div>
          ))}
        </div>
      </div>

      {/* ── 4. Status Indicators ──────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">4. Status Indicators</h2>
        <p className="ds-ref-label">StatusText component in every state:</p>
        <div className="ds-ref-row">
          {statusStates.map(s => (
            <div key={s.status} className="ds-ref-item">
              <StatusText status={s.status} icon>{s.label}</StatusText>
              <div className="ds-ref-code" style={{ marginTop: 4 }}>{s.status}</div>
            </div>
          ))}
        </div>
        <p className="ds-ref-label" style={{ marginTop: '1em' }}>Without icon:</p>
        <div className="ds-ref-row">
          {statusStates.map(s => (
            <div key={s.status + '-noicon'} className="ds-ref-item">
              <StatusText status={s.status}>{s.label}</StatusText>
            </div>
          ))}
        </div>
      </div>

      {/* ── 5. Caption Tones ──────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">5. Caption Tones</h2>
        <div className="ds-ref-row">
          {captionTones.map(tone => (
            <div key={tone} className="ds-ref-item">
              <Caption tone={tone}>{tone}</Caption>
            </div>
          ))}
        </div>
        <p className="ds-ref-label" style={{ marginTop: '1em' }}>Uppercase:</p>
        <div className="ds-ref-row">
          {captionTones.map(tone => (
            <div key={tone + '-upper'} className="ds-ref-item">
              <Caption tone={tone} transform="uppercase">{tone} label</Caption>
            </div>
          ))}
        </div>
      </div>

      {/* ── 6. Buttons ────────────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">6. Buttons</h2>
        <p className="ds-ref-label">Variant × Size × State:</p>
        <div className="ds-ref-row">
          <div className="ds-ref-item"><Button variant="default" size="normal">Default Normal</Button></div>
          <div className="ds-ref-item"><Button variant="primary" size="normal">Primary Normal</Button></div>
          <div className="ds-ref-item"><Button variant="default" size="compact">Default Compact</Button></div>
          <div className="ds-ref-item"><Button variant="primary" size="compact">Primary Compact</Button></div>
        </div>
        <p className="ds-ref-label" style={{ marginTop: '1em' }}>Selected state:</p>
        <div className="ds-ref-row">
          <div className="ds-ref-item"><Button variant="default" size="normal" selected>Selected</Button></div>
          <div className="ds-ref-item"><Button variant="primary" size="normal" selected>Selected Primary</Button></div>
        </div>
        <p className="ds-ref-label" style={{ marginTop: '1em' }}>Disabled:</p>
        <div className="ds-ref-row">
          <div className="ds-ref-item"><Button variant="default" size="normal" disabled>Disabled</Button></div>
          <div className="ds-ref-item"><Button variant="primary" size="normal" disabled>Disabled Primary</Button></div>
        </div>
      </div>

      {/* ── 7. Border Rules ───────────────────────────────── */}
      <div className="ds-ref-section">
        <h2 className="ds-ref-section-title">7. Visual Rules</h2>
        <ul style={{ paddingLeft: '1.5em', lineHeight: 1.8 }}>
          <li>1px solid <code>#000</code> borders on all components</li>
          <li>Border radius: <strong>0</strong> — no rounding anywhere</li>
          <li>No box-shadow or gradients</li>
          <li>Focus ring: 2px inset outline in <code>--mac-accent</code></li>
          <li>Hover: light grey background on interactive elements</li>
        </ul>
      </div>
    </>
  );
}

export { FoundationReference };
