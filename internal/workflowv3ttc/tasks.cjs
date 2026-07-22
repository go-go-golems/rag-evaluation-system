const rag = require("rag:ttc");
const task = require("workflow/task");

function finish(ctx, result, port, schema) {
  for (const usage of (result.usage || [])) ctx.usage.report(usage.dimension, usage.units);
  if (!result.ok) throw task.failure(result.failure);
  return ctx.outputs.putJSON(port, {schema, value: result.value});
}

exports.generate = task.implementation(ctx => {
  const output = finish(ctx, rag.generate(), "generated", "rag-ttc-generated/v1");
  return task.success({generated: output});
});

exports.embed = task.implementation(ctx => {
  const output = finish(ctx, rag.embed(), "embedded", "rag-ttc-prepared-shard/v1");
  return task.success({embedded: output});
});

exports.merge = task.implementation(ctx => {
  const output = finish(ctx, rag.merge(), "shard", "rag-ttc-prepared-shard/v1");
  return task.success({shard: output});
});

exports.validatePublication = task.implementation(ctx => {
  const output = finish(ctx, rag.validatePublication(), "receipt", "rag-ttc-validation-receipt/v1");
  return task.success({receipt: output});
});

exports.publish = task.implementation(ctx => {
  const output = finish(ctx, rag.publish(), "publication", "rag-ttc-publication-receipt/v1");
  return task.success({publication: output});
});

exports.evaluate = task.implementation(ctx => {
  const output = finish(ctx, rag.evaluate(), "evidence", "rag-ttc-query-evidence/v1");
  return task.success({evidence: output});
});
