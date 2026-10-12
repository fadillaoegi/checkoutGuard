# Prompt tindak lanjut runway di cfo.ai

Gunakan di chat cfo.ai yang menghasilkan rencana CheckoutGuard sebelumnya,
setelah saldo bertanggal dan asumsi cash flow terisi. Lampirkan berkas yang
disebutkan. Coding agent belum mengirim prompt ini ke akun cfo.ai.

```text
Update the existing CheckoutGuard financial plan with a numerical cash runway.
Keep the original export intact and save this as a separately dated revision.

Attached:
- EVIDENCE.md and cfo-validation.json, connecting the three supplied synthetic
  run observations to curated evidence. Do not treat them as customer traction.
- runway-result.json with a supplied cash balance, balance date, currency,
  minimum buffer and explicitly stated constant monthly cash-flow assumptions.
  This attachment is a local calculation, not a previous cfo.ai export.

Validate the required cash inputs before calculating. If any are missing,
ask for the missing value; do not treat it as zero. Preserve whether the opening
cash is actual or a user-defined simulation, and label future cash flows as
assumptions. Do not invent a balance, funding, sales or an invoice.

Reconcile a dated cash schedule:
opening unrestricted cash + collected customer/other receipts + documented
committed financing - cash operating payments - founder payments - taxes
- capex - debt service - refunds = closing unrestricted cash.
Keep categories distinct to avoid double counting. Distinguish revenue from
collections and operating profit from cash. Use the actual supplied founder
payment assumption rather than automatically keeping the original $6,000 pay.

For a flat positive net cash burn, show (unrestricted cash - minimum buffer)
divided by monthly burn, floored at zero. Also show the first scheduled period
whose cash reaches or falls below the buffer, the lowest cash balance and the
cash bridge for that period. A month-end balance is not an exact within-month
exhaustion date. If cash burn is nonpositive, do not fabricate an infinite
runway; explain the result under the stated forecast horizon.

If cash flows vary or committed financing is included, use their dated schedule
instead of extrapolating one optimistic burn month. Explicitly label timing,
collections, taxes, liabilities, refunds and any unverified forecast assumptions.
Report customer growth scenarios only if their assumptions are supplied; do not
invent scenario inputs to fill the table. Do not equate the original $59,175
peak cumulative operating deficit with available cash or a verified funding need.

Keep a single currency per cash schedule. The original operating forecast uses
USD. If the cash inputs use another currency, retain a separate schedule unless
an explicit FX rate and date are supplied; do not infer a rate.

Return the revised business plan and actual PDF/XLSX exports with the cash
input date, calculation, assumptions and source notes. Prepare a public copy
using only financial inputs approved for sharing; remove account identities,
credentials, personal banking details and private links. Include the numerical
runway in the demo only after this revision is actually produced and checked.
```
