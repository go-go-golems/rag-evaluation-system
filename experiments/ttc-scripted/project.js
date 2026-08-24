const research = require("researchctl");

module.exports = research
	.project("PRJ-TTC-SCRIPTED", "Thin scripted TTC acceptance on the converged experiment platform")
	.goal("Prove TTC has no bespoke lifecycle", (goal) =>
		goal.id("GOAL-TTC-SCRIPTED").status("active").priority("P0"),
	)
	.experiment("TTC scripted acceptance", (experiment) =>
		experiment.id("EXP-TTC-SCRIPTED").status("active"),
	);
