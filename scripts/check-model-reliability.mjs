// Starts and polls one real sandbox agent run. It never supplies browser actions.
import assert from "node:assert/strict";
import { mkdir, writeFile } from "node:fs/promises";
import { parseArgs } from "node:util";

const { values, positionals } = parseArgs({
  options: { retest: { type: "string" } },
  allowPositionals: true,
});
assert.ok(positionals.length <= 1, "Choose one scenario");
const scenario = values.retest ? "fixed" : positionals[0] || "timeout";
assert.ok(
  !values.retest || !positionals.length,
  "Use --retest ID without a scenario",
);
assert.ok(
  ["fixed", "timeout", "total_bug"].includes(scenario),
  "Choose fixed, timeout, or total_bug",
);
const base = "http://127.0.0.1:8080";
async function api(path, body) {
  const response = await fetch(
    base + path,
    body === undefined
      ? undefined
      : {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(body),
        },
  );
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `HTTP ${response.status}`);
  return result;
}
const runs = await api("/api/runs");
assert.ok(
  !runs.some((r) => r.status === "running"),
  "Another agent run is active",
);
const baseline = values.retest
  ? await api("/api/runs/" + encodeURIComponent(values.retest))
  : null;
if (baseline) {
  assert.equal(baseline.environment, "sandbox");
  assert.equal(
    baseline.total_verified,
    true,
    "Baseline must have a verified expected amount",
  );
  await mkdir("tmp/model-reliability", { recursive: true });
  await writeFile(
    "tmp/model-reliability/" + baseline.id + ".json",
    JSON.stringify(baseline, null, 2) + "\n",
    { mode: 0o600 },
  );
}
let run = baseline
  ? await api("/api/runs/" + encodeURIComponent(baseline.id) + "/retest", {})
  : await api("/api/runs", {
      environment: "sandbox",
      scenario,
      variant: "standard",
      url: base + "/sandbox",
      goal: "Buy exactly one tote using demo@example.test. Verify the checkout total and confirm exactly one persisted order. If payment times out, inspect order status before any retry.",
    });
console.log(
  JSON.stringify({
    id: run.id,
    scenario,
    parent_id: run.parent_id,
    status: run.status,
  }),
);
const deadline = Date.now() + 200_000;
while (run.status === "running" && Date.now() < deadline) {
  await new Promise((resolve) => setTimeout(resolve, 2000));
  run = await api("/api/runs/" + run.id);
}
const orders = await api("/api/runs/" + run.id + "/orders");
const directory = "tmp/model-reliability";
await mkdir(directory, { recursive: true });
const file = directory + "/" + run.id + ".json";
await writeFile(file, JSON.stringify(run, null, 2) + "\n", { mode: 0o600 });
let screenshots = 0;
for (const step of run.steps) {
  if (!step.evidence) continue;
  const response = await fetch(base + step.evidence);
  const bytes = new Uint8Array(await response.arrayBuffer());
  assert.equal(response.status, 200);
  assert.equal(response.headers.get("content-type"), "image/png");
  assert.deepEqual(
    Array.from(bytes.slice(0, 8)),
    [137, 80, 78, 71, 13, 10, 26, 10],
  );
  screenshots++;
}
const placements = run.steps.filter(
  (s) => s.action.tool === "click" && s.action.selector === "#place",
);
console.log(
  JSON.stringify({
    id: run.id,
    status: run.status,
    reason: run.reason,
    orders: orders.length,
    total: orders[0]?.total_minor,
    total_verified: run.total_verified,
    placements: placements.length,
    model_errors: run.steps.filter((s) => s.action.tool === "model_error")
      .length,
    screenshots,
    trace_id: run.trace_id,
    usage: run.usage,
    file,
  }),
);
assert.equal(
  run.status,
  scenario === "total_bug" ? "findings" : "passed",
  "Actual run failed; saved result remains available",
);
assert.equal(orders.length, 1);
const expectedTotal = baseline ? baseline.expected_total_minor : 12400;
assert.equal(
  orders[0].total_minor,
  scenario === "total_bug" ? 14400 : expectedTotal,
);
assert.equal(run.expected_total_minor, expectedTotal);
assert.equal(run.total_verified, true);
assert.equal(placements.length, 1);
if (baseline) {
  assert.equal(run.parent_id, baseline.id);
  assert.equal(run.goal, baseline.goal);
  assert.equal(run.variant, baseline.variant);
  assert.equal(run.displayed_total_minor, expectedTotal);
  assert.equal(
    run.findings.filter((f) => f.kind !== "timeout_recovered").length,
    0,
  );
}
if (scenario === "timeout")
  assert.ok(
    run.findings.some((f) => f.kind === "timeout_recovered" && f.verified),
  );
if (scenario === "total_bug")
  assert.ok(
    run.findings.some((f) => f.kind === "incorrect_total" && f.verified),
  );
