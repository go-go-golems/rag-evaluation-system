const task = require("workflow/task");
const rag = require("rag:workflow");

async function run(ctx, operation, port, schema) {
  ctx.checkpoint();
  let value;
  try {
    value = rag[operation]();
  } catch (_) {
    ctx.checkpoint();
    throw task.failure({
      class: "execution",
      code: "RAG_WORKFLOW_TASK_FAILED",
      retryable: false,
      message: "RAG workflow task failed",
    });
  }
  ctx.checkpoint();
  const output = await ctx.outputs.putJSON(port, {schema, value});
  return task.success({[port]: output});
}

exports.loadCorpus = task.implementation((ctx) => run(ctx, "prepare", "prepared", "rag-workflow-prepared/v1"));
exports.prepare = task.implementation((ctx) => run(ctx, "prepare", "prepared", "rag-workflow-prepared/v1"));
exports.query = task.implementation((ctx) => run(ctx, "query", "result", "rag-workflow-result-partition/v1"));
exports.merge = task.implementation((ctx) => run(ctx, "merge", "result", "rag-workflow-result-partition/v1"));
exports.publish = task.implementation((ctx) => run(ctx, "publish", "result", "rag-workflow-result/v1"));
