import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  css: {
    modules: {
      // Readable class names in dev (Button_root, Button_normal)
      // Short hashes in production
      generateScopedName:
        process.env.NODE_ENV === 'production'
          ? '[hash:base64:5]'
          : '[name]_[local]',
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: false,
    sourcemap: true,
    lib: {
      entry: {
        index: 'src/index.ts',
        ir: 'src/widgets/ir.ts',
        'app/index': 'src/app/index.ts',
      },
      formats: ['es'],
      cssFileName: 'styles',
    },
    rollupOptions: {
      external: ['react', 'react-dom', 'react-dom/client'],
    },
  },
});
