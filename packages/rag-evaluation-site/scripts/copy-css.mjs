import { copyFile, mkdir } from "node:fs/promises";
import { dirname, resolve } from "node:path";

const copies = [["src/theme.css", "dist/theme.css"]];

for (const [from, to] of copies) {
	const target = resolve(to);
	await mkdir(dirname(target), { recursive: true });
	await copyFile(resolve(from), target);
}
