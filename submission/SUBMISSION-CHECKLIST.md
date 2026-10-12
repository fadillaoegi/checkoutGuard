# Submission checklist

Repository: https://github.com/fadillaoegi/checkoutGuard · Team: FLdevBoo.
Checked against [official neatHack requirements](https://neatlogs.com/hackathon)
on 11 October 2026. This is a readiness record, not confirmation of submission.

| Requirement or evidence | Current state | Remaining work |
|---|---|---|
| Goal-driven agent with real tool use/state/verification | Proven for synthetic sandbox | Keep synthetic scope clear |
| Failure recovery | One order after payment timeout verified | Preserve original timeout evidence |
| Runtime neatlogs traces | Hosted readback verified | Add dashboard screenshots or accessible trace links |
| Entire use | Public checkpoint `01M4NKBNCP60MN7M5BXS143W16` linked to `a134e4a`; selected decisions, patches and original outputs reviewed; seven source blobs match published files | Show the verified context in the demo; see ENTIRE-REVIEW.md |
| cfo.ai pricing/cost/revenue/runway plan | Actual PDF/XLSX/Markdown exports included; 75 forecast periods and 5 sensitivities checked; runway method included | Numerical runway awaits supplied dated cash and cash-flow inputs; calculator/template/follow-up prompt ready; revised cfo.ai export still required |
| Public code with commits during the hackathon | User push output confirms implementation commit `380ed7a` on public `main` | Verify commit timestamp and reviewer access |
| GitHub Actions | Verified success for backend, frontend and publication on commit `a134e4a`; [hosted run](https://github.com/fadillaoegi/checkoutGuard/actions/runs/38144551603) | Later local runway/context review changes await user publication; synthetic CI does not verify marketplace access or provider reliability |
| Before/after improvement evidence | Store-fixture retest and exact-label cost observability verified | Preserve scope and source/checkpoint provenance; this does not establish provider reliability or autonomous source repair |
| Demo video ≤3 minutes covering all three partners | `submission/video/CheckoutGuard-final.mp4` v2, 2:56.6, voice-over + English subtitles; 45 s proof segment shows failure recovery, neatlogs trace readback, Entire checkpoint on GitHub and the cfo.ai plan (12 Oct) | Optional: replace the neatlogs readback panel with a real dashboard recording; upload and verify the accessible link |
| Building in public | Ready-to-post X texts in [BUILDING-IN-PUBLIC.md](BUILDING-IN-PUBLIC.md) (3 progress posts + submission quote post); not yet published | Publish genuine progress with #neatHack and record the links there |
| Official quote post on X + submission form | Pending. Per neatlogs email (12 Oct): quote-tweet the official @neatlogs launch post **tagging @neatlogs, @EntireHQ and @arithecfo**, and fill out the submission form; both required, links shared at 19:00 IST (20:30 WIB). Text and form answers ready in [X-POST-TEMPLATE.md](X-POST-TEMPLATE.md) | Post the quote-tweet with all three tags, submit the form, record both links |
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

`submission/video/` (≈1.1 GB of MP4/MP3/footage) and all `*.mp4`, `*.mov`, `*.webm`, `*.mp3`,
`*.wav`, `*.flac` files are git-ignored; share the video through its uploaded link instead.

Keep private env files, API keys, account screenshots, browser profiles, cookies,
raw prompts/runs, entire metadata/logs, dependencies, caches and generated builds
out of the public commit. Public examples intentionally contain no real keys.
Commit the curated `submission/evidence/` assets, source, lockfiles, examples,
docs and CI configuration. Local preflight output does not replace an actual
published Entire checkpoint.

Git remains user-managed under the original instruction. No commit, push,
branch/remote operation, account signup, billing action or X publication was
performed as part of preparing this package.
