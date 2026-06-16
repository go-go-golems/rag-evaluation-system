import { cp, rm, mkdir } from "node:fs/promises";
import { resolve } from "node:path";

const appDist = resolve("app-dist");
const target = resolve("../../pkg/defaultspa/dist");

await rm(target, { recursive: true, force: true });
await mkdir(target, { recursive: true });
await cp(appDist, target, { recursive: true });
console.log(`synced ${appDist} -> ${target}`);
