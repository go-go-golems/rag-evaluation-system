import { mkdtemp, writeFile, mkdir } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { execFileSync } from 'node:child_process';

const cwd = process.cwd();
const smokeDir = await mkdtemp(join(tmpdir(), 'rag-evaluation-site-consumer-'));
const tarballOutput = execFileSync('npm', ['pack', './dist', '--pack-destination', smokeDir], {
  cwd,
  encoding: 'utf8',
  stdio: ['ignore', 'pipe', 'inherit'],
}).trim().split('\n').at(-1);

if (!tarballOutput) {
  throw new Error('npm pack did not report a tarball filename');
}

const tarball = join(smokeDir, tarballOutput);

await writeFile(join(smokeDir, 'package.json'), JSON.stringify({
  private: true,
  type: 'module',
  scripts: {
    typecheck: 'tsc --noEmit',
    build: 'vite build',
  },
  dependencies: {
    '@go-go-golems/rag-evaluation-site': `file:${tarball}`,
    react: '^19.2.7',
    'react-dom': '^19.2.7',
  },
  devDependencies: {
    '@types/react': '^19.2.16',
    '@types/react-dom': '^19.2.3',
    '@vitejs/plugin-react': '^6.0.2',
    typescript: '~6.0.3',
    vite: '^8.0.16',
  },
}, null, 2));

await writeFile(join(smokeDir, 'tsconfig.json'), JSON.stringify({
  compilerOptions: {
    target: 'ES2020',
    lib: ['ES2020', 'DOM', 'DOM.Iterable'],
    module: 'ESNext',
    moduleResolution: 'bundler',
    jsx: 'react-jsx',
    strict: true,
    skipLibCheck: true,
    noEmit: true,
  },
  include: ['src'],
}, null, 2));

await writeFile(join(smokeDir, 'index.html'), '<div id="root"></div><script type="module" src="/src/main.tsx"></script>\n');
await mkdir(join(smokeDir, 'src'));
await writeFile(join(smokeDir, 'src/vite-env.d.ts'), '/// <reference types="vite/client" />\n');
await writeFile(join(smokeDir, 'src/main.tsx'), `
import { createRoot } from 'react-dom/client';
import { WidgetRenderer, defaultWidgetRegistry, type WidgetNode } from '@go-go-golems/rag-evaluation-site';
import '@go-go-golems/rag-evaluation-site/styles.css';
import type { RagEvaluationSiteAppProps } from '@go-go-golems/rag-evaluation-site/app';

const node: WidgetNode = { kind: 'text', text: 'consumer smoke' };
const props: RagEvaluationSiteAppProps = { apiBase: '/api/widget' };
void props;

createRoot(document.getElementById('root')!).render(<WidgetRenderer node={node} registry={defaultWidgetRegistry} />);
`);

execFileSync('npm', ['install', '--silent'], { cwd: smokeDir, stdio: 'inherit' });
execFileSync('npm', ['run', 'typecheck', '--silent'], { cwd: smokeDir, stdio: 'inherit' });
execFileSync('npm', ['run', 'build', '--silent'], { cwd: smokeDir, stdio: 'inherit' });

console.log(`clean consumer smoke passed in ${resolve(smokeDir)}`);
