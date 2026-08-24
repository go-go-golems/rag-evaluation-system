const task = require("workflow/task");
const intake = require("rag:intake");

async function run(ctx) {
  ctx.checkpoint();
  const result = intake.run();
  if (result.providerUsage) {
    for (const dimension of Object.keys(result.providerUsage).sort()) {
      ctx.usage.report(dimension, result.providerUsage[dimension]);
    }
  }
  if (!result.ok) {
    throw task.failure(result.failure);
  }
  ctx.checkpoint();
  const output = await ctx.outputs.putJSON("result", {
    schema: result.schema,
    value: result.value,
  });
  return task.success({result: output});
}

exports.run = task.implementation(run);
