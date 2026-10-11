# CheckoutGuard evidence for reviewers

Team **FLdevBoo** · [public repository](https://github.com/fadillaoegi/checkoutGuard).
These are synthetic sandbox orders. Model decisions, Chrome actions and
PostgreSQL verification were real; no marketplace purchase was made.

| Case | Result | Expected / displayed | Orders / submits | Duration |
|---|---|---|---|---|
| Pricing defect | Verified finding | 9,800 / 11,050 IDR | 1 / 1 | 21.30 s |
| Corrected retest | Passed | 9,800 / 9,800 IDR | 1 / 1 | 28.92 s |
| Payment timeout | Passed after checking saved orders | 12,400 / 12,400 IDR | 1 / 1 | 16.66 s |

The pricing calculation is **8,750 − 1,250 + 2,300 = 9,800 IDR**. The defective
fixture displays 11,050 IDR. Retest keeps the same goal and basket and selects
the corrected fixture; the agent does not edit a merchant's source code.

## Browser evidence

The exported PNGs are byte-for-byte copies of final browser evidence. All
32 screenshots across these four runs were checked for PNG signatures and
HTTP MIME. Only these four reviewed sandbox images are published; private
browser profiles, raw runs, model inputs and account screenshots stay ignored.

### Pricing defect

![Synthetic checkout with an incorrect total of 11,050 IDR](evidence/pricing-defect.png)

### Corrected retest

![Corrected synthetic checkout with one confirmed order of 9,800 IDR](evidence/corrected-retest.png)

### Payment timeout recovery

![Original payment timeout message retained in browser evidence](evidence/payment-timeout.png)

The last screenshot deliberately retains the original timeout message. The
`order_status` tool checks PostgreSQL directly, confirms one saved order under
the original correlation ID, and the agent finishes without another submit.
This demonstrates payment timeout recovery, not recovery from an unavailable
model provider.

## Hosted neatlogs readback

| Case | Trace ID | Readback |
|---|---|---|
| Pricing defect | `000d9deb6a723905d68a787ae9926c4f` | [verified summary](evidence/pricing-defect-trace.json) |
| Corrected retest | `7a97b78bbc758fb1cbadf231ccb0f210` | [verified summary](evidence/corrected-retest-trace.json) |
| Payment timeout | `94f12c0013ff1e04832776d908767f78` | [verified summary](evidence/payment-timeout-trace.json) |

Each trace had 15 spans: one root, seven LLM calls, and seven tools. Hosted
hierarchy, input/output availability and provider-reported token totals matched
local records. A successful defect-detection workflow still has application
status `findings`; it does not imply the displayed amount is correct.

[Run summaries](evidence/runs.json) preserve the original unknown recorded cost
and separately label posthoc estimates from the public model catalog. Those
estimates are not verified invoices. [SHA-256 manifest](evidence/manifest.json)
identifies the published artifacts; it is an integrity record, not a provider
signature. Reviewer-accessible hosted trace URLs/screenshots still need to be
added from the project dashboard before final submission.

## Cost-observability improvement

The earlier corrected retest recorded cost as unknown. After configuring rates
for the exact provider/returned label, a new retest of the same baseline/goal/
basket passed with one order, 12,361 input + 199 output tokens, and an estimate
of **US$0.00318975**: `(12,361 × 0.25 + 199 × 0.50) / 1,000,000`.
The published catalog was checked on 11 October 2026. This is a known-usage
estimate, not an invoice; missing request usage remains unknown and partial
totals are labelled. Historical records were not rewritten.

Trace `db1fd6a1e4cb8d7159c6cefc1c0b3daa` was read back with 15 valid spans and
matching token totals. [Verified summary](evidence/priced-retest-trace.json) ·
[Original screenshot](evidence/priced-retest.png).

This before/after evidence measures observability, not inference speed or a
reduction in cost: the provider-reported token counts differ. Actual feature
commit/checkpoint provenance remains pending user publication.

## Model identity and reliability

The configured provider was Public AI; requests selected
`swiss-ai/apertus-v1.5-70b`, while responses were labelled
`aisingapore/Qwen-SEA-LION-v4-32B-IT`. We report both and do not claim independent
verification of model weights or successful Apertus inference. These three runs
had no model timeouts; previous failed runs remain documented in
[EVALUATION.md](../docs/EVALUATION.md). Provider reliability is still unresolved.

## Corrected model observability

Retest `9689d682-97ac-4982-bfd6-538aff957d4e` passed in 35.15 seconds with one
9800 IDR order, one submit, eight valid screenshots and no model timeout.
[Hosted readback](evidence/model-diagnostics-trace.json) verifies all seven
LLM calls in trace `869bdb1c0d98cd794c3a895292fe23c7`: requested Apertus labels
remain stable; returned SEA-LION labels, 45000 ms client deadlines and error
categories are separately correlated by span ID. Usage 12306 input + 197 output;
estimated cost USD 0.003175. This demonstrates observability, not independent
verification of serving model weights or provider stability. Older trace
labels are historical and are not rewritten.

## Reproduction

Follow the [README](../README.md) to configure private local credentials and
start PostgreSQL, backend and dashboard. Run an incorrect-total check with
the alternate basket, then select **Retest corrected store**. For payment
timeout choose the standard basket. The agent chooses the browser actions;
the scripts only start/poll runs and verify actual results.

```sh
node scripts/check-model-reliability.mjs timeout
node scripts/check-model-reliability.mjs --retest YOUR_BASELINE_RUN_ID
```

Published run IDs belong to the author's local database. A new clone starts
with its own run IDs; do not expect these historical runs to appear in an empty
database. Evidence and safe summaries above remain inspectable without the DB.

## Competition status

Runtime evidence is ready. The user has published the implementation and the
existing hosted CI passes. Entire feature-checkpoint contents/access still
require verification. The model diagnostics update requires user publication
before its new CI checks can be verified on GitHub.
Original cfo.ai PDF/XLSX/Markdown exports are now included; see
[CFO-RESULTS.md](CFO-RESULTS.md) and [provenance](business/artifact-provenance.json).
The three supplied runtime observations match this file's curated run records.
The local [independent review](evidence/cfo-validation.json) matched 409 numeric
comparisons, including 75 forecast periods and 5 sensitivities, and inspected
all seven PDF pages. Original cfo.ai outputs are preserved byte-for-byte;
cfo.ai did not receive EVIDENCE.md, so this local comparison is a separate check.
Runway method is included; numerical cash runway requires actual dated cash
inputs. All pricing, customer and financial projections remain assumptions.
Demo video, building-in-public links and the official X quote submission remain
pending. See [SUBMISSION-CHECKLIST.md](SUBMISSION-CHECKLIST.md).
