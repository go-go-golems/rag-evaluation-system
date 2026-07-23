const research = require("researchctl");

module.exports = research
	.project("PRJ-RAG-WORKFLOW", "RAG v2 provider-free Workflow V3 parity fixture")
	.goal("Prove RAG semantic parity through durable workflows", (goal) =>
		goal.id("GOAL-RAG-WORKFLOW").status("active").priority("P0"),
	)
	.experiment("RAG Workflow V3 matrix", (experiment) =>
		experiment.id("EXP-RAG-WORKFLOW").status("active"),
	);
