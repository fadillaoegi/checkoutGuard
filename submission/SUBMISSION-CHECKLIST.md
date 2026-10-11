# Submission checklist

Repository: https://github.com/fadillaoegi/checkoutGuard · Team: FLdevBoo.
Checked against [official neatHack requirements](https://neatlogs.com/hackathon)
on 11 October 2026. This is a readiness record, not confirmation of submission.

| Requirement or evidence | Current state | Remaining work |
|---|---|---|
| Goal-driven agent with real tool use/state/verification | Proven for synthetic sandbox | Keep synthetic scope clear |
| Failure recovery | One order after payment timeout verified | Preserve original timeout evidence |
| Runtime neatlogs traces | Hosted readback verified | Add dashboard screenshots or accessible trace links |
| Entire use | Local capture prerequisites pass | User commits implementation; verify resulting checkpoint and reviewer access |
| cfo.ai pricing/cost/revenue/runway plan | Pending; inputs prepared | Generate actual plan through cfo.ai and include export/screenshots |
| Public code with commits during the hackathon | Repository visibility verified public; local source not committed | User publishes source and verifies event-dated implementation commits |
| Before/after improvement evidence | Store-fixture retest and exact-label cost observability verified | Preserve scope and source/checkpoint provenance; this does not establish provider reliability or autonomous source repair |
| Demo video ≤3 minutes covering all three partners | Script prepared | Record completed partner workflows, export, verify duration and accessible link |
| Building in public | Draft prepared | Publish genuine progress and keep #neatHack links |
| Official quote post on X | Pending | Quote the official submission post with team, project, repo and demo, then publish |
| Deadline | 12 Oct 20:30 WIB opening; planned close 13 Oct 01:25 WIB | Follow confirmed cutoff in Discord |

Sandbox results and curated evidence: [EVIDENCE.md](EVIDENCE.md). Shopee and
Tokopedia remain experimental preparation adapters; logged-in checkout success
has not been demonstrated. No requirement to make a real marketplace purchase
was found in the official event page.

## Checks before the user commits

```sh
python3 scripts/test-publication-audit.py
python3 scripts/audit-publication.py
```

The audit checks repository ignore rules with fixtures, scans candidate source
and evidence for known/high-confidence secret values, and checks tracked path
names by reading index metadata. It does not execute Git. It does not audit
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
