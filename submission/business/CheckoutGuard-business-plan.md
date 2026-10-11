# CheckoutGuard — Business plan

Prepared 11 October 2026. USD. Pre-revenue, assumption-led plan; not validated traction.

## Decision summary

Recommend a capped Team paid pilot, not unlimited usage or automated purchasing. Known catalog-priced token estimate: $0.00315467/successful synthetic run; mean supplied runtime: 22.2933 seconds. Modeled delivery cost including provisional retry, browser, database, evidence, telemetry and run support allowances: $0.0190856/completed run. Neither is a confirmed charge.

Base break-even: approximately 96 subscriptions at the assumed plan mix, or $8,650.67 monthly revenue. Fixed costs of $6,800/month include a $6,000 founder-pay assumption. Base first nonnegative operating month: March 2028; peak accumulated operating loss before that is approximately $59,175. This is a model-implied operating deficit, not verified funding need, cash balance or runway. Collections, tax, capex and liabilities remain unknown.

## Proposed pricing — assumptions

| Plan | Monthly price | Included workspaces | Pooled completed runs | Contribution at 50% use | Margin at full quota |
|---|---:|---:|---:|---:|---:|
| Developer | $29 | 1 | 100 | $21.88 | 72.1% |
| Team | $99 | 3 | 500 | $78.46 | 74.4% |
| Agency | $299 | 10 | 2,000 | $240.64 | 74.1% |

Support and fees are assumptions; the margins do not include fixed overhead.

## Revenue scenarios — all assumptions

| Case | Month 12 revenue (Oct 2027) | Month 24 revenue (Oct 2028) | Month 24 subscriptions (expected) | Month 24 operating result | First nonnegative operating month |
|---|---:|---:|---:|---:|---|
| Conservative | $1,864 | $2,871 | 31.9 | $-4,527 | Not reached through October 2028 |
| Base | $5,586 | $13,907 | 154.3 | $4,209 | 2028-03 |
| Growth | $12,342 | $40,765 | 458.0 | $25,464 | 2027-08 |

No implied actual users or sales. Fractional subscriptions represent expected counts. Month 24 revenue is monthly, not annual revenue.

## Known supplied runtime observations

| Synthetic run | Decisions | Screenshots | Seconds | Input tokens | Output tokens | Catalog estimate |
|---|---:|---:|---:|---:|---:|---:|
| 1 | 7 | 8 | 21.30 | 12,227 | 202 | $0.00315775 |
| 2 | 7 | 8 | 28.92 | 12,128 | 179 | $0.00312150 |
| 3 | 7 | 8 | 16.66 | 12,335 | 202 | $0.00318475 |

Evidence.md was referenced but not attached. No independent runtime verification is claimed. Configured request: Apertus; response label: Qwen SEA-LION v4 32B IT. Identity not independently verified. Earlier timed-out request usage remains unknown.

## Cost formula

Cost/completed run = known token estimate × provider price multiplier × (1 + extra-attempt ratio × extra-attempt token allowance) + browser seconds/attempt ÷ 3,600 × browser/hour × (1 + extra-attempt ratio) + incremental database/run + evidence/run + telemetry/run + support minutes/run ÷ 60 × support/hour.

Contribution = subscription revenue − completed runs × cost/run − workspace support − assumed payment fees. Operating result = contribution − fixed monthly costs.

## One-way sensitivity — base October 2028

| Test | Monthly operating result | Change from base |
|---|---:|---:|
| Operating result | $4,209 | $0 |
| Provider 10x operating result | $2,923 | $-1,285 |
| Extra attempts 100 percent operating result | $4,068 | $-140 |
| Browser 3x operating result | $4,073 | $-136 |
| Support 2x operating result | $2,117 | $-2,092 |
| Full quota usage operating result | $3,489 | $-720 |

Provider pricing ×10 reduces that month’s operating result by about $1,285; doubling all support reduces it by about $2,092. Do not equate cheap token usage with a cheap service.

---

# CheckoutGuard — Business plan

# CheckoutGuard

**Pre-revenue developer checkout QA — assumption-led business plan, 11 October 2026.**

Goal-driven checkout QA for developers and ecommerce agencies. Prototype: Go, React/TypeScript, Redux Toolkit, PostgreSQL and Chrome. Checks totals against integer price breakdowns, confirms stored synthetic orders, keeps screenshots and compares corrected retests.

**Evidence boundary:** no validated paying customers or revenue. All pricing, customer counts, conversion, churn and revenue projections are assumptions—not traction. Runtime observations below are supplied by the user; EVIDENCE.md has not been attached. No funding, cash balance, runway, provider invoice or production purchase authority is established.

**Safety boundary:** owned or explicitly authorized staging environments and synthetic orders only. Experimental Shopee/Tokopedia adapters prepare checkout and stop before final order confirmation. No automatic goods purchases or operation of merchant payment systems is assumed.

## Strategy: checkout-specific QA, not a general purchasing agent

Buyer: agency technical lead, ecommerce developer or QA owner responsible for checkout regressions. Start with owned or authorized staging sites and seeded synthetic orders. Core value: integer-safe totals assertions, persisted-order checks, screenshot evidence and corrected-retest comparisons. Alternative workflows are manual checklists, Playwright/Cypress scripts and general browser agents; CheckoutGuard must prove less setup and clearer fault evidence than these alternatives, not just that it can click through a checkout.

MVP priorities: reproducible fixtures, environment allowlists, synthetic test identities, idempotency and cleanup, correlation IDs, CI artifact links, timeout/retry budgets, model-routing/usage visibility and deletion controls. React/TypeScript and Redux manage workflow; Go orchestrates browser checks and PostgreSQL stores run/order assertions. Do not promise production merchant access. Experimental marketplace adapters are outside the initial paid offering.

Distribution: recruit design partners through ecommerce developer communities and agency leads; run founder-led demos on their staging checkout, then integrate into CI. Budgeted sales spend is an unvalidated capacity assumption. No validated demand, sales, funding or users are claimed.

## Validation plan — proposed gates, not results

**Weeks 1–2:** interview 15 qualified developers/agency QA owners about their last checkout incident, current scripts, time cost, approval/security rules and evidence needs. Gate: at least 8 describe a recent recurring problem and 5 grant access to an authorized staging test. Record refusals as well as enthusiasm.

**Weeks 3–4:** 5 design partners run at least 20 representative cases each, including discounts, taxes, shipping, currency rounding, broken totals, persisted-order failures and corrected retests. Instrument completion, false positives/negatives, tokens per attempt, timeout billing, screenshots/bytes, browser time and support minutes. Proposed gates: ≥95% task completion on supported fixtures, ≥90% seeded-defect detection, ≤5% false alarms, and p95 completion ≤60 seconds. These targets are not established by three synthetic runs. Human review confirms all seeded outcomes.

**Weeks 5–6:** present real paid offers—not only willingness-to-pay surveys. Test Team $99 with 500 pooled runs against $149 with the same quota and an Agency $299 offer. Aim for 3 paid pilots from 10 qualified trial organizations, usage in at least 3 of 4 weeks and two renewal commitments after a full month. No sale is recorded until payment is real. Stop expanding acquisition if activation/renewal fail; diagnose trust, integration effort, false alarms and security objections first.

**Weeks 7–8:** reconcile invoices and run logs; test p95 high-token/retry cases, support and full-quota margins. Gate: ≥70% contribution margin at the intended mix, sustainable support burden, and documented authorization/data-retention controls. Replace assumptions in this model with observed cohort figures before claiming a repeatable business.

## Runway — deliberately uncalculated

Simple steady-burn runway = unrestricted available cash ÷ monthly net cash burn, only when burn is positive. A better answer is the first month when opening unrestricted cash + cumulative actual cash receipts − cash payments + committed financing falls below minimum cash. This plan’s cumulative operating result is not a cash balance or a runway claim.

Needed: cash amount and as-of date/currency; restricted cash; minimum cash buffer; founder pay timing; actual hosting/provider/payment invoices; taxes; collections/refunds; payables; capex; debt service; and only documented committed financing. Variable burn requires a dated cash-flow schedule, not dividing by one optimistic month. No funding, balance or runway has been invented.

---

# CheckoutGuard — Assumptions

## Proposed prices — test, don’t promise

All values on this page are editable planning assumptions in USD, excluding tax. Developer: $29/month, one workspace, 100 completed runs. Team: $99/month, three workspaces, 500 runs pooled. Agency: $299/month, ten workspaces, 2,000 runs pooled. Charge per organization subscription, not per developer seat. Start with Team-focused paid pilots; do not launch unlimited usage. No overage revenue is modeled: default to a hard cap plus explicit upgrade approval; if validated, test $0.10/completed-run overages with a monthly spending cap. Corrected retests consume a run. Provider retries should not consume a second customer quota unit.

Prices buy auditable checkout checks, regression evidence and less QA work—not tokens. Retention assumption: 30 days for screenshots, longer only by explicit agreement. The $0.001 evidence/run allowance is provisional: 8 screenshots × assumed 0.25 MB × retention and storage/request rates plus export/egress; measure real files before contracting.

## Cost basis and unknowns

Public AI catalog lists Qwen SEA-LION v4 32B IT at $0.25/input million and $0.50/output million: https://platform.publicai.co/models . Checked 11 October 2026. Requested Apertus and response labels disagree; model identity is not independently verified. These rates estimate known reported usage; they do not establish actual model routing or charges.

Three successful synthetic runs are not a workload distribution or reliability estimate. Earlier timeout token usage is unknown. **Extra attempt ratio = additional attempts per completed run**, assumed 20%, not measured failure probability. Timeout token allowance ratio assumes each extra attempt uses 100% of a successful run’s tokens; test 0.5×–2× and obtain provider logs/invoices. Never replace unknown timeout usage with zero. Browser assumption is 45 seconds/attempt at $0.12/hour, above the supplied 22.29-second mean to allow cold-start/setup time. Browser/Go worker cost is included in that rate; fixed React/API hosting is separate.

DB, storage, telemetry and support costs are placeholders, not invoices. Database $0.002/completed run and telemetry $0.0005/run are incremental allowances above the separate fixed service baselines. Support has both 0.02 minutes/run and 10/25/60 minutes/subscription/month at $30/hour, as incremental capacity cost. Founder compensation $6,000/month covers development/product time, excluding that support capacity; avoid counting the same hours twice. Other fixed costs: hosting $150, DB $60, telemetry $30, overhead $60, sales $500/month: total $6,800 including founder compensation.

## Acquisition and forecast conventions

Base qualified prospects/month at launch: Developer 40, Team 20, Agency 5; assumed conversion 5%/10%/10%; pipeline grows 5% monthly and churn is 3% monthly. Conservative: conversion 2.5%/5%/5%, no pipeline growth, 5% churn. Growth: conversion 10%/15%/20%, 8% monthly pipeline growth, 2% churn. These are explicit targets, not observed lead flow, conversion or retention.

October 2026 is a zero-subscription **planning seed**, not a historical customer census; launch is assumed November 2026, with 24 forecast months through October 2028. Fractions are expected subscriptions, not actual partial customers. Revenue uses full-month ending subscriptions, collected in that month by assumption; no annual prepay, expansion, refunds, taxes or receivable delay are included. Payment processing uses an unverified 3% + $0.30/subscription/month assumption. Customer mix changes only through acquisition and churn. No financing or cash balance is modeled.

---

# CheckoutGuard — Forecast

## Reading this forecast

Main Scenario is the base case. Conservative and Growth change only assumed acquisition conversion, monthly prospect growth and churn; pricing, usage, support and fixed costs stay the same. New subscriptions = qualified prospects × conversion; lost subscriptions = prior ending subscriptions × churn; ending = prior + new − lost. Ending subscriptions receive a full month of revenue, an optimistic billing convention that should be replaced with real subscription start dates.

Contribution subtracts browser/model delivery, incremental database/evidence/telemetry, run-level and workspace support, and **assumed payment fees of 3% + $0.30/subscription/month**. They are not verified merchant/provider charges. Fixed $6,800/month is then deducted once to give operating result, before tax, depreciation and financing.

Break-even subscriptions = ceiling(fixed monthly cost ÷ weighted contribution/subscription). Break-even revenue uses that same customer mix. It is a steady-state threshold, not proof that customers can be acquired. Cumulative operating result starts from zero as a profit accumulation seed, not a cash balance. Cash burn may differ with collections, prepayments, tax, capex and liabilities.

One-way sensitivities hold customer demand and prices fixed: model token prices ×10; extra attempts increase from 20% to 100% per completed run; browser compute ×3; all support ×2; or run quota usage 50%→100%. Effects show in operating result, not confirmed spending. Provider cost and support changes do not silently change conversion or churn. A usage increase requires delivery capacity and may affect service quality; those second-order effects are not modeled.

## Deliverable verification and limits

- base: all 25 periods matched independent calculations for subscriptions, revenue, contribution, operating result and cumulative result.
- conservative: all 25 periods matched independent calculations for subscriptions, revenue, contribution, operating result and cumulative result.
- growth: all 25 periods matched independent calculations for subscriptions, revenue, contribution, operating result and cumulative result.
- Five one-way sensitivities matched independent October 2028 calculations; signs reduce operating results as expected.
- Annual native reads checked: flows sum monthly values; 2028 year-end balance is intentionally absent because the horizon stops in October. No full-year 2028 balance is asserted.
- Actual export: scanned XML/relationships for personal/account/credential identifiers, checked external links and Excel error cells. No matched private identifiers or external workbook links found; export preserved byte-for-byte.

- Native model simulator rendered without reported errors or warnings in light and dark at 1,200px width. Visible screenshots were reviewed in chat. Controls/click actions and narrower widths were not browser-tested.
- The reviewed Canvas screenshots were not returned as downloadable image files by cfo.ai, so they are **not included in this artifact**. Do not claim this package preserves those images. Attach the screenshots from the app/chat to complete that evidence requirement.
- This package is prepared for public review, but has not been published or given a public share link. No account, payment-method, subscription or purchasing operation was performed.
