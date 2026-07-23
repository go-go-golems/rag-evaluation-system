const task = require("workflow/task");
const rag = require("rag:workflow");

async function run(ctx, operation, port, schema) {
  ctx.checkpoint();
  const result = rag[operation]();
  if (!result.ok) {
    ctx.checkpoint();
    throw task.failure(result.failure);
  }
  const value = result.value;
  ctx.checkpoint();
  const output = await ctx.outputs.putJSON(port, {schema, value});
  return task.success({[port]: output});
}

exports.loadCorpus = task.implementation((ctx) => run(ctx, "prepare", "prepared", "rag-workflow-prepared/v1"));
exports.prepare = task.implementation((ctx) => run(ctx, "prepare", "prepared", "rag-workflow-prepared/v1"));
exports.query = task.implementation((ctx) => run(ctx, "query", "result", "rag-workflow-result-partition/v1"));
exports.merge = task.implementation((ctx) => run(ctx, "merge", "result", "rag-workflow-result-partition/v1"));
exports.publish = task.implementation((ctx) => run(ctx, "publish", "result", "rag-workflow-result/v1"));
