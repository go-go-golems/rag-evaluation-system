import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const backendPort = parseInt(process.env.RAG_EVAL_BACKEND_PORT || "8772", 10);

export default defineConfig({
	plugins: [react(), tailwindcss()],
	resolve: {
		alias: [
			{
				find: "@go-go-golems/rag-evaluation-site/app",
				replacement: fileURLToPath(
					new URL("../packages/rag-evaluation-site/src/app/index.ts", import.meta.url),
				),
			},
			{
				find: "@go-go-golems/rag-evaluation-site/ir",
				replacement: fileURLToPath(
					new URL("../packages/rag-evaluation-site/src/widgets/ir/index.ts", import.meta.url),
				),
			},
			{
				find: "@go-go-golems/rag-evaluation-site",
				replacement: fileURLToPath(
					new URL("../packages/rag-evaluation-site/src/index.ts", import.meta.url),
				),
			},
		],
	},
	server: {
		port: 5173,
		proxy: {
			"/api": {
				target: `http://127.0.0.1:${backendPort}`,
				changeOrigin: true,
			},
		},
	},
	build: {
		outDir: "../internal/web/dist",
		emptyOutDir: true,
		sourcemap: true,
	},
});
