const { rag } = require("./common");
let calls = 0;
const preparation = rag.fragment("immediate", (p) => {
	calls += 1;
	p.units(rag.transcript.units.agentsViewRuns()).chunks(rag.chunks.recursive({ maxRunes: 800 }));
});
module.exports = { calls, enumerableKeys: Object.keys(preparation) };
