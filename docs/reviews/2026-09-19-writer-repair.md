# Callable writer repair — 2026-09-19

Baseline: `5117b9238f9ccc7caf01bc9a5e0ab1877ba76208`. Stage 1 scope: the seven
[confirmed audit defects](2026-09-19-functionality-orchestration-audit.md), complete validation
before pending redo, and actor binding before replay/recovery. No Corvint or Beamfall source
changes, runtime, real-queue operation, new owner decision or qualification claim.

Independent Gate A: PASS after adding full pre-redo validation and separating bounded request
and intent reads. Model selection: Astra/high for invariant reasoning; running model/effort
settings cannot be inspected or changed from this task. One builder and one read-only reviewer.

## Behavior and regression evidence

| Audit finding | Repair | Named regression |
|---|---|---|
| F1 canonical authority | Strict journal-derived canonical queue/policy/tickets; queue-wide refusal on drift | `TestTMV0007_AS35_WriterRejectsProjectionDrift` |
| F2 actual branch | Bounded no-follow primary HEAD observation; recheck before effects | `TestTMV0007_AS35_WriterObservesPrimaryBranch`, `TestTMV0007_AS29_LinkedCallerUsesPrimaryHEAD` |
| F3 primary identity | Head authority checked before recovery | `TestTMV0009_AS11_WriterGuardsBeforeRedo` |
| F4 ALL barrier | Entry guard and pure-model defense; ADMISSION remains usable | `TestTMV0009_AS11_WriterGuardsBeforeRedo`, `TestTMV0016_AS27_AdmissionBarrierAllowsNativeMutation` |
| F5 invalid INIT | Model validation under lock before state-directory creation; concurrent initializer recheck | `TestTMV0009_AS11_InvalidInitLeavesCorrectableInput`, `TestTMV0009_AS11_InvalidInitRoleThenRetry` |
| F6 restore marker | Presence/type/error refuses before recovery | `TestTMV0009_AS11_WriterGuardsBeforeRedo` |
| F7 replay identity | Original receipt ticket ID returned by validated journal lookup | `TestTMV0006_AS03_ReplayPreservesIdentityAndAllowsStableDivergence`, `TestTMV0006_AS03_TicketCreateRetryReplays` |
| Recovery ordering | Full audit before redo; exact pending digest/head binding | `TestTMV0009_AS11_RedoRequiresCompleteJournalProof` |
| Replay trust | Actor binding precedes replay/redo; private request projections validated | `TestTMV0006_AS03_ActorBindingBeforeReplayAndRedo`, `TestTMV0006_AS03_ForgedRequestNeverReplays` |

Refusal witnesses compare complete state/intent file contents. Positive witnesses cover normal
mutation, valid pending redo, linked callers, exact replay despite stable drift, corrected INIT
and a journalled ADMISSION barrier. Focused tests passed during implementation.

## Validation

Canonical `make verify`: **PASS**, exit 0 in 159.10 seconds on Go 1.27.1 Darwin/arm64; format,
all tests and vet. [Log](2026-09-19-writer-evidence/verify.log),
[result](2026-09-19-writer-evidence/verify-result.json), and
[frozen source manifest](2026-09-19-writer-evidence/source-manifest.json). Source hashes and
HEAD were unchanged throughout. Independent Gate B: **PASS**, no actionable findings; reviewer independently verified the source hashes ([review](2026-09-19-writer-evidence/gate-b-review.md)). The gate supervisor uses signal/finally process-group cleanup with an
interruption regression; it does not operate a runtime or qualify process containment.

## Remaining work

TCP-02 remains incomplete: pending/staging reconciliation, remaining administrative commands, archive restore, interrupted
genesis/staging recovery, hostile-editor CAS, authenticated authority, reservations and liveness.
The CLI still uses the no-runtime oracle. G1..G6 and GP remain unqualified; no performance
claim is made. Strict mutation audit conservatively blocks unrelated work on any divergent
intent projection; exact validated retries remain available.

Corvint owns planning and evidence. Native planning/GP baseline, runtime handshake/usage budgets,
exact-candidate gates/review/reducer, owner cutover, fanout and routing remain later stages.
The six exact accepted source copies recorded in decision 0001 are missing from their temporary
paths; recovering and digest-verifying them remains a prerequisite to resolving upstream
contract questions. Current sibling source cannot silently replace accepted bytes.

## Read-only receipt inspection

`receipt audit` exposes the existing full native journal audit through the snapshot-protected
CLI. Success reports structural consistency, projection agreement, codec coverage and the
terminal receipt digest; historical acceptance, actor authentication, liveness and runtime
qualification remain explicitly NOT_OBSERVED. Errors emit no partial evidence item. Reads take
no lock, initialize nothing and never redo a pending receipt. `receipt show` and `receipt replay`
remain NOT_RUN; planning/review/reducer replay is not fabricated.

Gate A: PASS. Focused CLI suite and built-executable fixture smoke: PASS. Final `make verify`: **PASS**, exit 0 in 142.72 seconds; frozen source and HEAD unchanged. Independent Gate B: PASS, including the completed gate and frozen-source verification. [Gate log](2026-09-19-writer-evidence/receipt-verify.log), [review](2026-09-19-writer-evidence/receipt-review.md), [smoke](2026-09-19-writer-evidence/smoke-summary.json).
Named witnesses: `TestTMV0008_AS07_ReceiptAuditReportsBoundedEvidence`,
`TestTMV0008_AS11_ReceiptAuditRefusesWithoutRecovery`,
`TestTMV0008_AS07_ReceiptAuditUsageAndUninitialized`, and
`TestTMV0008_AS36_ReceiptAuditDiscardsMovedResults`. These include full no-write checks,
corrupted early history, pending redo, marker, divergence, foreign primary, and bounded retries
for head/intent movement with no stale output.

Source recovery searched 206 candidate files under project checkouts, with zero digest matches
against decision 0001; the two uncommitted governing documents were also absent from searched
temporary and agent directories. [Search result](2026-09-19-writer-evidence/source-recovery-search.json).
The owner has been asked for a durable backup; current sibling documents are not substituted. This blocks upstream ambiguity resolution and qualification, not already specified local fixture work.

## Explicit fixture intent reconciliation

Decision 0007 releases the existing KEEP_JOURNAL/ADOPT_FILE mechanics for settled fixture queues.
`reconcile inspect` supplies the canonical ticket record/hash through a complete, target-scoped
audit. `reconcile intent` requires preserved original bytes and an explicit operator choice;
exact retries retain the original identity. KEEP retains malformed or empty discarded bytes
before receipt commit. ADOPT preserves protected fields and revision rules. Neither path
interprets ticket prose as executable input or changes a qualification axis.

Gate A, including returned-fault publication ordering: PASS. Focused store/journal/CLI workflows:
PASS. Final canonical `make verify`: **PASS**, exit 0 in 151.62 seconds; source hashes and HEAD
unchanged. Independent Gate B: **PASS**, including final source binding and executable evidence
([review](2026-09-19-writer-evidence/reconcile-review.md)). Named witnesses and exact CLI item
and input contracts are in SPEC §3.3. Pending receipts refuse without recovery; active staging
refuses fresh reconciliation but permits validated read-only replay;
ALL/ADMISSION reconciliation remains possible in a settled store. This slice does not implement
runtime supervision, reservations, real-queue admission, archive restore or hostile-editor CAS.


[Final reconciliation gate](2026-09-19-writer-evidence/reconcile-final-verify.log),
[result](2026-09-19-writer-evidence/reconcile-final-verify-result.json),
[source manifest](2026-09-19-writer-evidence/reconcile-final-source-manifest.json),
[built CLI workflow](2026-09-19-writer-evidence/reconcile-smoke-summary.json), and
[archive evidence retention](2026-09-19-writer-evidence/reconcile-archive-summary.json).
The one LOW review finding (empty inapplicable ADOPT digest flag) is repaired and covered by a
before-input regression. The earlier 137.62-second gate is retained as earlier-tree evidence,
not reused for the repaired tree. A concurrent gate belonged to `corvint-wt-issue20` and was
preserved; these results establish no timing/performance property.

The checked-in `examples/fixture` inputs and README walkthrough passed the built executable
path: init/create/retry, strict audit refusal of malformed input, canonical inspect, KEEP and
ADOPT, exact retries, subsequent mutation, receipt audit, export/verify retaining discarded
bytes, and explicit unavailable runtime commands. The resolved canonical-retrieval memory item
is closed; interrupted genesis/staging recovery, production writer prerequisites, archive
historical semantics, native Linux and qualification gaps remain open.


## Settled fixture pause/unpause

Decision 0008 exposes the existing PAUSE/UNPAUSE model. Pause creates ADMISSION/OPERATOR;
unpause removes a journalled barrier while retaining unrelated edited ticket bytes. Recorded
retries preserve the original outcome instead of reapplying an old pause after unpause. Unknown
qualification axes stay unknown. The complete BarrierRemoval audit permits only present regular
canonical ticket projections to differ; all queue/policy/private and extra/missing-file checks remain.

The callable applier now publishes the validated null barrier post after receipt and ordinary
posts, before head, through the existing pinned exact-pre unlink/sync primitive. Returned failures
before commit and before deletion preserve head/barrier; a third value at the deletion boundary
refuses without erasure or head advancement. Pending UNPAUSE recovery remains unsupported by all
current callable commands and is documented explicitly. This is not crash or power-loss proof.

Gate A PASS. Focused normal/refusal/replay/empty-D tests PASS. Canonical `make verify`: **PASS**,
exit 0 in 138.32 seconds; frozen source hashes and HEAD unchanged. Independent source/contract
review **PASS**, including final gate/hash/evidence binding
([review](2026-09-19-writer-evidence/barrier-review.md)). Exact syntax and named tests are in SPEC §3.3.


[Barrier gate](2026-09-19-writer-evidence/barrier-verify.log),
[result](2026-09-19-writer-evidence/barrier-verify-result.json),
[source manifest](2026-09-19-writer-evidence/barrier-source-manifest.json), and
[built workflow](2026-09-19-writer-evidence/barrier-smoke-summary.json).
The executable workflow includes pause/status, an unrecorded NoChange, ticket creation under
ADMISSION, unpause, exact retries after opposite operations, preservation of empty D, inspection,
KEEP and final receipt audit. These results do not promote the queue or qualify performance.
