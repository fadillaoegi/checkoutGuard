# Submission checklist

Repository: https://github.com/fadillaoegi/checkoutGuard · Team: FLdevBoo.
Checked against [official neatHack requirements](https://neatlogs.com/hackathon)
on 11 October 2026. This is a readiness record, not confirmation of submission.

| Requirement or evidence | Current state | Remaining work |
|---|---|---|
| Goal-driven agent with real tool use/state/verification | Proven for synthetic sandbox | Keep synthetic scope clear |
| Failure recovery | One order after payment timeout verified | Preserve original timeout evidence |
| Runtime neatlogs traces | Hosted readback verified | Add dashboard screenshots or accessible trace links |
| Entire use | Local capture prerequisites pass; user push output reports one checkpoint ref sent | Verify checkpoint contents linked to implementation commit and reviewer access |
| cfo.ai pricing/cost/revenue/runway plan | Actual PDF/XLSX/Markdown exports included; 75 forecast periods and 5 sensitivities checked; runway method included | Publish artifacts and show actual cfo.ai use in video; numerical cash runway awaits real dated cash inputs |
| Public code with commits during the hackathon | User push output confirms implementation commit `380ed7a` on public `main` | Verify commit timestamp and reviewer access |
| GitHub Actions | Verified success for all three jobs on commit `0e21026`; [hosted run](https://github.com/fadillaoegi/checkoutGuard/actions/runs/38133238996) | Preserve this evidence; synthetic CI does not verify marketplace access or provider reliability |
| Before/after improvement evidence | Store-fixture retest and exact-label cost observability verified | Preserve scope and source/checkpoint provenance; this does not establish provider reliability or autonomous source repair |
| Demo video ≤3 minutes covering all three partners | Indonesian recording script updated; target 2:55, sandbox focus and marketplace limitations | Record real partner workflows, export, verify duration and accessible link |
| Building in public | Draft prepared | Publish genuine progress and keep #neatHack links |
| Official quote post on X | Pending | Quote the official submission post with team, project, repo and demo, then publish |
| Deadline | 12 Oct 20:30 WIB opening; planned close 13 Oct 01:25 WIB | Follow confirmed cutoff in Discord |

Sandbox results and curated evidence: [EVIDENCE.md](EVIDENCE.md). Shopee and
Tokopedia remain experimental preparation adapters; logged-in checkout success
has not been demonstrated. No requirement to make a real marketplace purchase
was found in the official event page.

Recording sequence and ready-to-read narration: [DEMO-SCRIPT.md](DEMO-SCRIPT.md).
Keep unsuccessful marketplace tests as limitations; use sandbox results for
the complete demonstrated workflow. Switching environments does not recover
marketplace access.

Deadline rechecked against the official page: submission opens Monday
**12 October 2026 at 20:30 WIB**, with a planned close Tuesday
**13 October at 01:25 WIB** (12 October 19:00–23:55 IST).
The confirmed cutoff is announced on the official Discord; the website still
labels the closing time as planned. Prepare the video and cfo.ai evidence before
the opening window. The actual cfo.ai exports are now included in
[CFO-RESULTS.md](CFO-RESULTS.md).

## Checks before the user commits

```sh
python3 scripts/test-publication-audit.py
python3 scripts/review-cfo-model.py
python3 scripts/audit-publication.py
```

The audit checks repository ignore rules with fixtures, scans candidate source
and evidence for known/high-confidence secret values, and checks tracked path
names by reading index metadata. It also scans decompressed XLSX/ZIP members
within bounded archive sizes. It does not execute Git. It does not audit
staged blob contents or historical commits; the user still reviews their commit
selection. An ignore rule does not remove files that were already tracked.

Keep private env files, API keys, account screenshots, browser profiles, cookies,
raw prompts/runs, entire metadata/logs, dependencies, caches and generated builds
out of the public commit. Public examples intentionally contain no real keys.
Commit the curated `submission/evidence/` assets, source, lockfiles, examples,
docs and CI configuration. Local preflight output does not replace an actual
published Entire checkpoint.

Git remains user-managed under the original instruction. No commit, push,
branch/remote operation, account signup, billing action or X publication was
performed as part of preparing this package.
