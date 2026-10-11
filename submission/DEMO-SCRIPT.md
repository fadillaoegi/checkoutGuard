# Three-minute demo script — draft

Record the completed workflows before assembling the video. Do not fill the
Entire or cfo.ai segments with simulated evidence; those segments wait for
actual checkpoint/plan access. Label any cuts in execution wait time.

| Time | Screen | Narration / proof |
|---|---|---|
| 0:00–0:15 | CheckoutGuard dashboard | “CheckoutGuard checks checkout flows before release. A goal drives real browser actions, and results are verified against stored orders.” |
| 0:15–0:45 | Pricing-defect run | Show quantity one, expected 9,800 versus displayed 11,050 IDR, verified finding and screenshot. “The discount was not applied to the total.” |
| 0:45–1:05 | Linked corrected retest | Show same goal/basket, 9,800 IDR and one saved order. “Retest uses the corrected synthetic fixture; the agent does not edit merchant source code.” |
| 1:05–1:25 | Payment-timeout run | Show original timeout message, order_status and one submit. “The order already exists, so the agent checks status and avoids a second submission.” |
| 1:25–1:55 | neatlogs actual project dashboard | Show workflow → LLM/tool tree, latency/tokens and a failure trace. Explain how traces guided diagnosis of model timeout. Distinguish estimates from billed costs. |
| 1:55–2:15 | Entire actual checkpoint | Show committed feature checkpoint/context and repository link. Explain capture during implementation. Record only after checkpoint is verified. |
| 2:15–2:40 | Actual cfo.ai plan | Show pricing, cost assumptions, revenue scenarios and runway. Say customer/revenue figures are assumptions. Record after real output is available. |
| 2:40–2:55 | Scope and result | “Sandbox is verified; Shopee/Tokopedia preparation is experimental and stops before final ordering. Model provider reliability remains a limitation.” |
| 2:55–3:00 | Public repository | Team FLdevBoo, public code and supporting evidence. |

Optional marketplace material should fit within the limit by replacing a short
app segment, not the required partner evidence. Show `ready_for_review` only if
it actually occurs. A CAPTCHA/blocked result is not a successful checkout.

## Building-in-public draft — not published

> Building CheckoutGuard for #neatHack with FLdevBoo: a goal-driven browser agent
> that checks checkout totals, saves evidence, and confirms orders. Our synthetic
> demo found an incorrect total, passed the corrected retest, and recovered a
> payment timeout without a duplicate order. We are documenting provider timeouts
> and using neatlogs traces to investigate them. Code and verified evidence:
> https://github.com/fadillaoegi/checkoutGuard

Do not describe cfo.ai or a published Entire checkpoint as complete until their
actual outputs are available. Add progress links to SUBMISSION-CHECKLIST.md
after the user publishes them. Quote-submission draft waits for the official
neatlogs post and completed demo/partner evidence; opening a draft is not submission.
