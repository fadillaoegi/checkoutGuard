# CheckoutGuard

<img src="checkoutguard_FE/public/brand/checkoutguard-mark.svg" width="72" alt="CheckoutGuard logo">

**Evidence before release.** A goal-driven browser agent that tests a checkout,
checks price arithmetic, confirms saved orders, and preserves evidence.
Built for neatHack by **FLdevBoo**.

[Public repository](https://github.com/fadillaoegi/checkoutGuard) ·
[Reviewer evidence](submission/EVIDENCE.md) ·
[Architecture](docs/ARCHITECTURE.md) ·
[Bahasa Indonesia](README-ID.md)

![CheckoutGuard dashboard with goals, environments and saved verification results](submission/evidence/dashboard.png)

## Verified synthetic demo

| Case | Verified result |
|---|---|
| Incorrect total | Expected 9,800 IDR, displayed 11,050 IDR; one saved order and a pricing finding |
| Corrected retest | Same goal and basket; total 9,800 IDR, one saved order, passed |
| Payment timeout | Saved order checked before retry; total 12,400 IDR, one submit and one order |

These runs used real model calls, Chrome and PostgreSQL. The model chose the
browser actions; evaluation scripts only started/polled runs and checked results.
All 24 screenshots were valid PNGs and all three hosted neatlogs traces were
read back with correct hierarchy and matching token totals. Inspect the actual
[images, trace summaries and manifest](submission/EVIDENCE.md).

Retest selects a corrected **synthetic store fixture**. The MVP does not edit a
merchant's source code or repair a production store. All sandbox orders and
payments are simulated.

## What the agent does

1. Receives a goal and environment.
2. Observes the browser and selects a schema-validated action.
3. Executes the approved action and saves a screenshot/result.
4. Verifies amounts with integer arithmetic and orders with persistent data.
5. Records failures and performs bounded recovery before choosing the next step.

The dashboard shows steps, findings, expected/displayed amounts, saved orders,
usage, cost estimates and linked baseline/retest results. English is the default;
Bahasa Indonesia is available in the header. The language preference persists
in the browser; custom goals and original evidence remain unchanged.

## Stack and architecture

Go 1.26.3 · React/TypeScript · Redux Toolkit · PostgreSQL 18 · Chrome/chromedp.
Domain and application layers depend inward. HTTP, PostgreSQL, model, browser
and tracing implementations are adapters; startup injects them. Frontend state
is Redux-based, with use cases behind repository ports. Dependency-direction
checks run in backend tests and the frontend build. See
[ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Run locally

Requirements: Go 1.26.3, Node 24.15.0, pnpm 11.1.3, Docker, Chrome/Chromium.
Use your own provider/telemetry credentials. From the repository root:

```sh
cp docs/examples/neatlogs.env.example .env.neatlogs.local
cp docs/examples/publicai-model.env.example .env.model.local
chmod 600 .env.neatlogs.local .env.model.local
```

Fill the key values locally; never put them in the frontend or a public commit.
The backend reads these files without executing shell assignments. Environment
values take precedence. Existing parent-directory `.neatlogs-api.env` and
`.apertus-api.env` configuration remains supported. `NEATLOGS_ENV_FILE` and
`MODEL_ENV_FILE` can override the paths.

Optional known-usage estimates from the checked public model catalog:

```sh
cp docs/examples/publicai-pricing.env.example .env.pricing.local
```

Prices apply only when both provider and returned model label match. They are
estimates, not verified invoices. Unknown usage stays unknown; partial estimates
are labelled as a lower bound. Historical records are not rewritten.

```sh
docker compose up -d db
sh scripts/dev-backend.sh
```

In another terminal:

```sh
cd checkoutguard_FE
pnpm install --frozen-lockfile
pnpm dev
```

Open http://127.0.0.1:5173. Backend: http://127.0.0.1:8080. Dedicated database:
127.0.0.1:55433. Set `CHROME_PATH` when Chrome is not found automatically.
Choose Sandbox, use the local URL `http://127.0.0.1:8080/sandbox`, then run
**Incorrect total** or **Payment timeout**. Use **Retest corrected store** to
compare the same goal/basket. A fresh clone has an empty run history; reviewer
artifacts are available separately in `submission/evidence/`.

## Model reliability and scope

The local demo profile and example select `swiss-ai/apertus-v1.5-8b` through
Public AI, using JSON-object responses and short action summaries. The previous
70B profile remains supported; historical requests for
`swiss-ai/apertus-v1.5-70b` received successful responses labelled
`aisingapore/Qwen-SEA-LION-v4-32B-IT`.
New LLM spans keep the requested model as their label and record the response
label separately in persisted `step.usage.model_call` diagnostics. Workflow
output carries the same records, correlated by span ID for hosted readback.
Historical spans may use the response label on success; this is not independent
proof of the serving model weights.
Provider timeouts remain intermittent. The three successful demo runs do not
establish production reliability. [Diagnosis](docs/MODEL-CONNECTION-DIAGNOSIS.md)
and [candidate provider setup](docs/PROVIDER-CANDIDATE.md) document the evidence
and remaining work. A backup provider has not been activated.

One run/browser at a time; at most 18 actions and 180 seconds. Model requests
have a default 45-second deadline (`MODEL_REQUEST_TIMEOUT_SECONDS`, 10–60).
Transient failures use bounded retries and Retry-After; permanent errors stop.
Planning retries never replay browser actions. Sandbox submit retry requires a
prior order-status check and is forbidden once an order is confirmed.
Each decision includes the latest browser observation, with the original goal
and complete action/result history. Superseded page snapshots are removed from
later requests; persisted steps and screenshots are retained.
Terminal status is published after browser cleanup, with run admission released
alongside final persistence. Trace export does not block the next checkout or
retest after the completed result becomes visible.

Shopee/Tokopedia are **experimental checkout preparation** adapters. Open the
dedicated browser and log in manually, then supply a full HTTPS product URL.
The agent stops before final order/payment confirmation. Marketplace amounts
and orders are not independently verified. Logged-in marketplace success is
not demonstrated; Shopee access verification has blocked manual tests.

CAPTCHA/login/unavailable pages require user action. The app reports blocked
access, preserves failed page checks, and offers Sandbox selection without
starting a run automatically. It does not complete access challenges.

## Tests and diagnostics

Safe hosted diagnostics parser (no credentials or network required):

```sh
python3 scripts/test-read-trace-diagnostics.py
```


Backend with optional real PostgreSQL and synthetic Chrome fixtures:

```sh
cd checkoutguard_BE
CHECKOUTGUARD_BROWSER_TEST=1 CHECKOUTGUARD_TEST_DATABASE_URL='postgres://checkoutguard:checkoutguard_local@127.0.0.1:55433/checkoutguard?sslmode=disable' go test -race -timeout 60s ./...
go vet ./...
```

Frontend:

```sh
cd checkoutguard_FE
pnpm format:check
pnpm test
pnpm build
```

Live sandbox evaluation from the root, with the app idle (uses model/telemetry
access and creates synthetic orders):

```sh
node scripts/check-model-reliability.mjs timeout
node scripts/check-model-reliability.mjs --retest YOUR_BASELINE_RUN_ID
```

From the backend directory, inference-only diagnostics:

```sh
go run ./cmd/modelprobe -catalog
go run ./cmd/modelprobe -attempts 2
go run ./cmd/modelprobe -env-file ../.env.model-candidate.local -attempts 2
```

Probe output does not expose keys or execute returned actions. Use
`-pricing-file ../.env.pricing.local` for optional estimates. Hosted readback:

```sh
python3 scripts/read-traces.py TRACE_ID --run-file tmp/path/to/run.json
```

Add `--expected-status error` for traces containing errors, even when the
application recovered. Hosted verification and checkout outcome are separate.

## Public repository contents

```sh
python3 scripts/test-publication-audit.py
python3 scripts/audit-publication.py
```

The audit uses ripgrep ignore rules, reads tracked path metadata without Git
commands, and scans source/evidence for known and high-confidence secret values.
It does not audit staged blob contents or historical commits. User-managed
commit review remains necessary. `.gitignore` excludes private env files,
browser profiles, raw evidence/runs, local coding-agent state, Entire private
metadata, caches, dependencies, test outputs and generated builds.

Curated sandbox PNGs and safe trace/run summaries in `submission/evidence/`
are intentionally public. Raw runs/prompts and marketplace/account screenshots
are not part of that package. The CI workflow runs publication checks, Go/PG/
Chrome fixture tests, frontend tests and builds without production API keys.
After a fixture-test-only fix for Ubuntu's Chrome sandbox restrictions,
[hosted CI passed](https://github.com/fadillaoegi/checkoutGuard/actions/runs/38133238996)
on commit `0e21026`: backend race/PG/Chrome tests and Go vet, frontend tests/build,
and publication audit. This verifies synthetic fixtures, not live marketplace
checkout or model-provider reliability. See [evaluation](docs/EVALUATION.md).

## neatHack status

neatlogs integration and runtime readback are proven. Entire capture prerequisites
pass. The public Entire checkpoint `01M4NKBNCP60MN7M5BXS143W16` links to
`a134e4a`; its metadata matches all 35 changed files and transcript blobs are
accessible. Selected decisions, patches and original tool outputs were inspected;
seven source files match published commit blobs. See
[checkpoint verification](submission/evidence/entire-checkpoint.json).
For the demo, follow the [Entire context review](submission/ENTIRE-REVIEW.md).
The user has published the implementation and
[hosted CI passes for `a134e4a`](https://github.com/fadillaoegi/checkoutGuard/actions/runs/38144551603):
backend, frontend and publication all succeeded. Later local runway/context
review changes still await user publication.
Actual [cfo.ai business-plan exports](submission/CFO-RESULTS.md) are included:
pricing, cost assumptions, three revenue scenarios and a runway method.
Numerical cash runway awaits actual dated cash inputs. Demo video and publication
links remain pending; the latest local changes await user publication.
This repository preparation is not a completed submission.

[Runway inputs and local calculator](submission/RUNWAY-INPUTS.md) ·
[cfo.ai revision prompt](submission/CFO-RUNWAY-FOLLOWUP.md).
Missing cash inputs stay unavailable; financial results remain in ignored `tmp`.

[Submission checklist](submission/SUBMISSION-CHECKLIST.md) ·
[Demo script](submission/DEMO-SCRIPT.md) ·
[cfo.ai inputs](submission/CFO-PROMPT.md) ·
[cfo.ai results](submission/CFO-RESULTS.md) ·
[Detailed evaluation](docs/EVALUATION.md) · [Build log](hackathon.md)

Git operations and publication remain user-managed. Do not put credentials,
private browser state or personal/payment information into the public evidence.
