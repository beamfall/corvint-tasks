# Stage 3 independent Gate B — fixture reconciliation

Initial result: **PARTIAL**, one LOW CLI-contract defect; no HIGH or MED finding. Stage 1/2 review is retained. This review covers only the accepted reconciliation delta and accompanying decision/SPEC/roadmap/report/README/example documentation. No reviewer tests, source edits, nested delegation, or owned background processes. Selected Astra/high for invariant reasoning; live model/effort settings are not independently inspectable.

## Finding

- **LOW — an empty inapplicable flag is accepted.** `internal/cli/reconcile.go:119` tests the digest value rather than flag presence. With otherwise valid ADOPT arguments and input, `--adopt-file --canonical-sha256 ''` passes flag validation and can proceed to the transaction. The accepted plan and SPEC require inapplicable flags to refuse before input reads. Use `seen["--canonical-sha256"]` for ADOPT and extend `TestTMV0007_AS35_CLIReconcileInputFailuresWriteNothing` with an empty-value case. This is established by the opened parser branch; no extra reproduction was executed during the frozen gate.

## Acceptance

| Criterion | Result | Evidence |
|---|---|---|
| Complete canonical audit with exactly one target exception | PASS | `internal/journal/audit.go:86`, `internal/journal/records.go:68`, `internal/journal/records.go:420`; value-copy option, parsed queue scope, canonical target required, other/private projections strict, pending error never grants write authority. Named canonical-read and unproved-history tests cover empty/malformed input, other drift, corruption, scope, pending, and moving observations. |
| Pure inspect with canonical record digest and explicit unknowns | PASS | `internal/cli/reconcile.go:18`; one bounded native capture/walk/recapture audit, no item on error, untrusted record, same-proof barrier/snapshot and NOT_OBSERVED axes. `internal/cli/reconcile_test.go:15`, `:139`, `:168`. |
| Explicit input, stable digest, validated replay before fresh checks | PASS | `internal/store/reconcile.go:41`, `:54`, `:82`, `:88`, `:92`; validates original bytes/choice/actor, obtains original identity from validated request lookup, never silently refreshes original input. `internal/store/reconcile_test.go:49`, `:92`. |
| Canonical fresh inputs and complete pre-apply proof binding | PASS | `internal/store/reconcile.go:100`, `:110`, `:137`, `:144`; model inventory plus journal-derived records, second full audit, entire Identity equality, branch and startup guard rechecks. Conservative local-operator premise remains explicit; no hostile-editor CAS claim. |
| Guard scope and recovery refusal | PASS | `internal/store/guards.go:48`, `:76`; `internal/store/reconcile.go:114`; fresh mutation and redo staging checks. Reconciliation is the only ALL exemption. Private index/audit errors, pending receipts and fresh active staging cause no recovery. Refusal, pending, descriptor and journalled barrier tests in `internal/store/reconcile_test.go:119`, `:168`, `:216`, `:323`. |
| Discarded evidence before commit, immutable reuse | PASS | `internal/store/store.go:96`, `:211`, `:234`; evidence POST maps to RoleEvidence and phase 0, no added artifact or limit. `internal/store/reconcile_test.go:239` checks evidence visible before the returned fault, unchanged head/projection/requests/receipts, empty staging and successful orphan reuse. This proves finite returned-fault behavior, not process interruption or power loss. |
| Closed CLI/input bounds | PARTIAL | `internal/cli/reconcile.go:73`, `:125`; duplicate/unknown/choice/role checks and bounded no-follow file/stdin handling are present. Empty inapplicable digest presence remains the LOW finding above. |
| Fixture-only authority and end-to-end workflow | PASS | Decision 0007; SPEC explicit reconciliation section; README disposable fixture inputs and commands. Store/CLI named workflows cover inspect, KEEP malformed/empty, ADOPT, exact replay after newer state, and subsequent ordinary mutation. No runtime, real-queue, authentication, performance or qualification claim. |

## Documentation corrections

The repair report's opening `no new owner decision` and Remaining work's `reconciliation` refer to Stage 1 and need qualification now that decision 0007 and settled reconciliation are added (`docs/reviews/2026-09-19-writer-repair.md:6`, `:41`). Describe active-staging refusal as applying to fresh reconciliation: the accepted plan deliberately allows a validated read-only replay before the fresh staging check. Pending receipts still refuse replay. README/example inputs otherwise match the delivered fixture workflow.

## Gate evidence

First Stage 3 canonical `make verify`: **PASS**, exit 0 in **137.62s**, recorded in `reconcile-verify-result.json`. Reviewer independently recomputed all **126** hashes in `reconcile-source-manifest.json`: **zero mismatches**. HEAD remains `5117b9238f9ccc7caf01bc9a5e0ab1877ba76208`. No commit during this gate. A parser repair invalidates this final-tree gate; the builder will run one repair gate, with reviewer attention limited to the changed parser/test/docs and its frozen evidence.

Usage telemetry at review completion: 8% weekly consumed; no reset used. All review work remains read-only except this scratch evidence file.
