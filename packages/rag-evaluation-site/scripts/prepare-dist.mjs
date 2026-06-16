import { copyFile, readFile, writeFile } from "node:fs/promises";
import { resolve } from "node:path";

const sourcePackage = JSON.parse(await readFile(resolve("package.json"), "utf8"));

const distPackage = {
	name: sourcePackage.name,
	version: sourcePackage.version,
	private: false,
	type: sourcePackage.type,
	license: sourcePackage.license,
	description: sourcePackage.description,
	keywords: sourcePackage.keywords,
	homepage: sourcePackage.homepage,
	bugs: sourcePackage.bugs,
	repository: sourcePackage.repository,
	sideEffects: sourcePackage.sideEffects,
	publishConfig: sourcePackage.publishConfig,
	peerDependencies: sourcePackage.peerDependencies,
	dependencies: sourcePackage.dependencies,
	exports: sourcePackage.exports,
	files: ["**/*.js", "**/*.d.ts", "**/*.css", "**/*.json", "README.md"],
};

await writeFile(resolve("dist/package.json"), `${JSON.stringify(distPackage, null, 2)}\n`);
await copyFile(resolve("README.md"), resolve("dist/README.md"));
