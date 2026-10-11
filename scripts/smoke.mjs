// Live-model evaluation against the synthetic store; no replay or fixed action script.
import { mkdir, writeFile } from 'node:fs/promises'
import assert from 'node:assert/strict'

const base = process.env.CHECKOUTGUARD_URL || 'http://127.0.0.1:8080'
const goal = 'Buy one tote using demo@example.test. Verify the checkout total and confirm exactly one order. Recover safely from any timeout.'
async function request(path, body) {
  const response = await fetch(base + path, body === undefined ? {} : { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) })
  const result = await response.json()
  if (!response.ok) throw new Error(result.error || `HTTP ${response.status}`)
  return result
}
async function wait(id) {
  const deadline = Date.now() + 210_000
  let previous = -1
  while (Date.now() < deadline) {
    const run = await request('/api/runs/' + id)
    if (run.steps.length !== previous) { console.log(`${run.id.slice(0, 8)} ${run.status}: ${run.steps.length} steps`); previous = run.steps.length }
    if (run.status !== 'running') return run
    await new Promise(resolve => setTimeout(resolve, 2000))
  }
  throw new Error('Evaluation deadline exceeded')
}
const records = []
async function evaluate(input, parent) {
  const created = parent ? await request('/api/runs/' + parent + '/retest', {}) : await request('/api/runs', { goal, ...input })
  const run = await wait(created.id)
  records.push(run)
  await mkdir('tmp', { recursive: true })
  await writeFile('tmp/live-evaluation.json', JSON.stringify({ type: 'live-model-synthetic-sandbox', evaluated_at: new Date().toISOString(), runs: records }, null, 2))
  assert.equal(run.order_count, 1, run.reason)
  assert.equal(run.total_verified, true, run.reason)
  assert.ok(run.usage.input_tokens > 0)
  const orders = await request('/api/runs/' + run.id + '/orders')
  assert.equal(orders.length, 1)
  assert.equal(orders[0].total_minor, run.displayed_total_minor)
  for (const step of run.steps) { assert.ok(step.evidence); const response = await fetch(base + step.evidence); assert.equal(response.status, 200); assert.equal(response.headers.get('content-type'), 'image/png'); const bytes = new Uint8Array(await response.arrayBuffer()); assert.deepEqual(Array.from(bytes.slice(0, 8)), [137, 80, 78, 71, 13, 10, 26, 10]) }
  console.log(`Completed ${run.scenario}/${run.variant}: ${run.status}, ${run.steps.length - 1} model calls, ${run.usage.input_tokens + run.usage.output_tokens} tokens, trace ${run.trace_id}`)
  return run
}
const existing = process.env.BASELINE_RUN_ID
let baseline
if (existing) { baseline = await wait(existing); records.push(baseline) }
else baseline = await evaluate({ scenario: 'total_bug', variant: 'standard' })
assert.equal(baseline.status, 'findings')
assert.ok(baseline.findings.some(f => f.kind === 'incorrect_total' && f.verified))
const retest = await evaluate({}, baseline.id)
assert.equal(retest.status, 'passed'); assert.equal(retest.parent_id, baseline.id)
assert.equal(retest.findings.filter(f => f.kind === 'incorrect_total').length, 0)
const timeout = await evaluate({ scenario: 'timeout', variant: 'standard' })
assert.equal(timeout.status, 'passed')
assert.ok(timeout.findings.some(f => f.kind === 'timeout_recovered' && f.verified))
assert.equal(timeout.steps.filter(s => s.action.selector === '#place').length, 1)
const variation = await evaluate({ scenario: 'fixed', variant: 'alternate' })
assert.equal(variation.status, 'passed'); assert.equal(variation.expected_total_minor, 9800)
console.log('Live evaluation passed: pricing defect, corrected retest, timeout recovery, alternate basket; zero duplicate orders.')
