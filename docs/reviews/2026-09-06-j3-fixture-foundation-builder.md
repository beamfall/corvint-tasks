# J3-03/04 independent fixture foundation — builder report

Date: 2026-09-06. Base: `27460b6caa8a760344d15fc160d6effbb9de7a54`.
Worktree: `/private/tmp/corvint-tasks-build-20260906/tcp02-j3-fixture-tests`.

This report covers only the independently scheduled fixture foundation. It does
not complete J3-03/04. Owner decision 0002, the integration plan including its
parent-selected final amendments, and Gate A including its bounded scheduling
clarification govern this work. AGENTS and its four referenced conventions were
read. The optional `tcp02-j3-03-04-steering.md` was checked before finalization and
was absent.

## Exact source claim

Only these three NEW worktree files were written:

| File | SHA-256 |
|---|---|
| `internal/authority/fixture_transaction_test.go` | `92be72c1a04a5741464bba2c910f5de8c272d6edd7f4b8d02e8c4c89133e5879` |
| `internal/authority/fixture_recovery_test.go` | `6e71754b0a5cb80001a6c31d1f69b8cf40ea440a93e61fba13eee4ce911d6abd` |
| `internal/authority/fixture_transaction_fault_test.go` | `9d70e2ea4103789a397c6469ab77ecef4cd36d47da60f8420f4bc2a8a083601d` |

No production source, public/private production API, codec, existing test,
SPEC, ROADMAP, BUILD-LOG, memory file, Corvint or Beamfall source was edited. No
Git mutation, delegation, canonical full gate, or independent review was run.
The parent owns the combined documentation, gate and fresh review. Scratch logs
and this report are outside the worktree.

## Focused execution evidence

Final command, exit 0:

```sh
GOTOOLCHAIN=local \
GOCACHE=/private/tmp/corvint-tasks-build-20260906/go-cache \
GOMAXPROCS=2 GOFLAGS=-p=2 \
go test ./internal/authority -run FixtureFoundation -count=1 -json
```

Environment: `go version go1.27.0 darwin/arm64`, Darwin arm64. Package elapsed
145.444 seconds. All 17 top-level tests passed, with 333 named test/subtest PASS
events, no FAIL or SKIP events. Final JSON evidence:
`/private/tmp/corvint-tasks-build-20260906/tcp02-j3-foundation-focused.jsonl`.
`gofmt -l` on exactly the three claimed files returned no paths. Linux execution
was NOT_RUN. The Darwin/Linux build constraints reflect the accepted native
primitive platforms; they do not disable deferred integration assertions.

All new top-level tests carry their TM-V0/AS IDs:

| Named test suffix (`TestTMV..._AS..._FixtureFoundation...`) | Implemented evidence |
|---|---|
| `TMV0009_AS11_FreshRecovery` | Five operations; receipt-only, first-post mixtures and completed-head leftovers; fresh locks/sessions; actual-byte recovery and independent capacity observations. |
| `TMV0007_AS35_SeedRefusals` | Missing/corrupt canonical blob or private request, malformed/gapped ledger and naked ticket refuse before staging exists; source/metadata purity on refusal. |
| `TMV0006_AS03_NoArtifacts` | PAUSE/UNPAUSE/KEEP NoChange, actor/protected-ADOPT/stale-KEEP refusal, unchanged ADOPT revision bump, malformed and valid divergent KEEP. |
| `TMV0009_AS11_EveryPublicationPrefix` | Common table's actual finite boundaries for all five operations, every assigned slot, evidence links, receipt, each post/deletion, head and cleanup; reopen and abort/recover. |
| `TMV0009_AS11_CommittedSyncFailure` | Each operation's successful exclusive receipt link remains committed when receipt-directory sync returns an error. |
| `TMV0009_AS11_PendingRefusals` | Missing/corrupt blob, malformed ledger, unaffected/affected private third values, preserved third intent, foreign descriptor, wrong complete slot, extra receipt, completed post mismatch. |
| `TMV0009_AS27_PartialAndForeignStage` | Returned short write with retained partial slot, truncated malformed unpublished temp, malformed active, unknown/symlink/FIFO/directory entries and wrong complete temp; partial tokens cannot publish. |
| `TMV0009_AS27_OrphansAndMissingSlots` | Abort KEEP, retain exact D and linked orphan, reserved fresh UNPAUSE, replay/NoChange, later reuse of the same orphan inode; receipt-only recovery with all payload slots removed from the owned fixture subset. |
| `TMV0014_AS10_HeadDigitGrowth` | Actual finite history reaches sequence 10, including exact head growth from 9 and model-final inventory/cost equality. |
| `TMV0001_AS10_SourceLifetime` | Released lock, old-session token, replaced source inode/parent and wrong branch refuse while preserving foreign/source bytes and metadata. |
| `TMV0009_AS11_ReturnedFaults` | Each boundary name actually emitted by each operation is failed at its first occurrence, including primitive create/write/close/link/rename/unlink/mkdir/sync/cleanup boundaries; fresh-session recovery/abort. This is not every syscall occurrence cross-product. |
| `TMV0014_AS10_CapacityAndCleanupDebt` | Actual prepare byte costs, retained promised bytes, failed unlink-sync name/shrink debt and accepted Cleanup rejection of early descriptor unlink. |
| `TMV0006_AS03_ReplayOutcome` | Complete original recorded outcome and canonical request binding preserved; conflict remains conflict; model authority/runtime coverage stays NOT_OBSERVED. |
| `TMV0009_AS27_CleanupBoundaries` | Every staging-sync cut across the actual completed cleanup names for each operation, including failed active-unlink sync; old era held until fresh durable cleanup. |
| `TMV0009_AS11_StalePlan` | A changed physical D after Model refuses before staging preparation. |
| `TMV0009_AS11_RenameSyncAndRelease` | Returned failure after native rename at destination/source directory sync, plus real handle closes/flock release with sticky joined close errors and old-session refusal. |
| `TMV0009_AS11_FiveOperations` | Successful INIT/PAUSE/UNPAUSE/KEEP/ADOPT through the same table; exact replay/conflict; present-empty D retained; C/D/C-prime behavior; untouched fixture source bytes/mode/mtime. |

## Implemented mechanisms and boundaries

`fixture_transaction_test.go:110` (`ftSeed`) constructs a modest two-receipt J1
seed, with a blob-backed canonical ticket, request record, head, evidence and
physical projection. Native strict Audit, archive Export and archive Verify all
succeed before any staging directory or lock is created, with byte/mode/mtime
purity checks. This is fixture construction, not an ADD/import executor. Empty
INIT is a separate row and never supplies a ticket.

The constructor creates the disposable primary Git fixture and retains its root
and common-directory identities. Fresh sessions reacquire the real matching
Lock and accepted mount-qualified private handles. No arbitrary input path,
environment repository, supplied role, same UID, fixture bit or lock establishes
an issuer. The finite scan refuses runtime/attempt/effect/worktree/bootstrap
content and requires the fixture queue, empty reservations and zero configured
runtimes. Source text is data only.

`fixture_transaction_test.go:378` (`publish`) is the single operation table:
descriptor temp preparation/full sync, exclusive active link, durable temp
cleanup, bounded assigned slots, durable evidence (including evidence POST),
exclusive receipt link, growth/addition posts before shrink, barrier deletion
last among posts, head, then durable cleanup with active last. Existing evidence
is reused only after actual byte/identity/durability checks. Receipt link success
is remembered across later sync errors. Fresh preparation checks actual base and
physical preimages against the frozen plan. NoChange/replay/refusal never enters
publication. Persistent staging is retained, including after clean completion.

`fixture_transaction_test.go:292` (`replaceTicket`) is test-only and scoped to
the harness's selected AT-01 projection and explicit KEEP/ADOPT operation. It
checks the live branch, source token and exact physical pre digest inside the
session/lock operation lifetime, uses native rename, attempts both directory
syncs, and joins close/performed-effect errors. It does not widen production
roles and does not claim hostile-editor CAS.

`fixture_recovery_test.go:330` (`observe`) reads actual no-follow retained bytes
twice with file/directory/HEAD identities, lengths, mode/mtime and byte agreement.
Its finite limits are 12 linked receipts, 128 entries per directory and 4 MiB
captured bytes, with existing per-file codecs/bounds and a narrower 1 MiB read
ceiling for this modest fixture. It validates actual receipt continuity,
genesis/head identity, post/blob content identity, exact request/outcome/revision
bindings and the finite seed/administrative/reconciliation shapes. Stable use
requires private projection agreement and permits only selected controlled D.
Pending use validates the committed base plus exactly one next receipt and its
descriptor bindings, allows affected physical PRE/POST mixtures, and retains
unaffected private agreement. KEEP uses physical D, including present-empty D;
its discarded evidence matches that preimage and the restored C remains exact.
UNPAUSE preserves the explicitly selected unaffected D after aborted KEEP.

`fixture_recovery_test.go:783` (`reopen`) qualifies fresh session-local tokens
from actual no-follow file identities and declared lengths/digests. It never
copies old tokens. Partial/unpublished observations are cleanup-only; wrong
complete slots refuse. It invalidates the previous token map only after all
replacement observations qualify. Precommit classification and cleanup reuse
`transaction.ClassifyStage` and `transaction.Cleanup`. Pending ClassifyStage's
Plan requirement is not bypassed or presented as actual-byte recovery.

`fixture_recovery_test.go:927` (`recover`) accepts no Plan, original invocation,
cached post bytes or Model output. It derives posts/missing slots and exact head
only from the actual linked receipt, retained blobs and genesis. Existing
destinations receive durability confirmation. Completed leftovers require
actual receipt/head/post agreement and durability confirmation before cleanup.
A failed sync after active removal keeps the ownership era held; fresh recovery
confirms staging sync before releasing that debt. Orphans and immutable history
are never removed by recovery/abort.

Capacity uses actual finite inventory bytes/names and accepted
`archive.MeasureManifestEncoding`, with staging charged separately. Successful
common-table prefixes are componentwise dominated by accepted CheckCapacity
prefixes. Prepare and fresh recovery have direct syscall-boundary byte witnesses;
recovery's external test observer retains only a separate expected/capacity
oracle and supplies no recovery bytes. Rename/name/shrink debt is retained across
the required sync boundary; Cleanup explicitly proves no early unlink credit.
Exact final Cost and complete Inventory equality are checked only after durable
cleanup. Transient manifest metadata is an encoding-cost oracle, not an archive
or native reader acceptance result. No accepted producer under-bound remained in
the exercised schedules; no model cap or production dependency was changed.

## Failures, corrections and evidence limits

Initial focused runs exposed these test-builder errors; the final run includes
their corrections:

- Used nonexistent `wire.Count.Uint64`; corrected to accepted `Count.Int`.
- Compared live pre-cleanup staging against the final zero-staging envelope;
  corrected the check to keep the preceding bound until the durable step ends.
  This was a test timing error, not a model-cap repair.
- Unchanged ADOPT lacked explicit selected-ticket ownership; selection now
  accompanies the original reconciliation request, even when D equals C.
- A foreign-descriptor fixture changed requestId without changing its declared
  request path, failing during setup; it now changes a valid recorded timestamp
  and reaches the intended retained ownership/binding refusal.
- Partial-write cleanup injection fired during descriptor-temp cleanup first;
  it now fires after the selected actual short write.
- UNPAUSE completed cleanup initially rejected preserved selected D after an
  aborted KEEP; selected unaffected intent is preserved without relaxing any
  unaffected private projection or affected-ticket PRE/POST guard.
- Prefix enumeration initially reused INIT's repository-specific pinned path in
  a new fixture. The stop selector now recognizes that pinned-post boundary while
  each publication still uses its own exact descriptor target.

Self-review additionally removed duplicate stale token registration, added stale
Plan/base/pre checks, strengthened finite receipt/KEEP bindings, confirmed
completed receipt/head/post durability before cleanup, and added the explicit
post-active-unlink sync-debt regression. Intermediate focused logs remain
scratch evidence; they are not substituted for the final successful run:
`tcp02-j3-foundation-faults.log`, `tcp02-j3-foundation-recovery.log`.

There were no changes to another owner's file and no remaining concrete source
dependency conflict for this bounded foundation. The external dependency is the
explicitly deferred, reviewed-final J3-02 integration below. No unfinished J3-02
source or guessed symbols were read/copied/imported.

Faults are returned syscall errors or deterministic owned-fixture publication
prefixes/subsets. Go defers run on returned faults. No child crash test, killed
process, service or process supervisor was used; process-interruption/descendant
and physical power-loss evidence are NOT_PRODUCED. Missing-slot cases are
deliberately constructed owned subsets, not claims about killed-process states.
Logical name/file-byte accounting is not a physical-block, allocation,
power-loss, timing, RSS or million-receipt/GiB saturation qualification. Existing
pure arithmetic extremes were not rerun under this focused authorization.

## Exact deferred integration — all NOT_PRODUCED / pending

These are precisely the five Gate A scheduling groups, not disabled tests or
acceptance of current reader failures. No t.Skip or placeholder assertion was
added for them.

1. **Reviewed-final shared J3-02 observer parity.** Handoff at
   `fixture_recovery_test.go:783` (`reopen`), `:330` (`observe`), and
   `fixture_transaction_fault_test.go:283`
   (`TestTMV0009_AS27_FixtureFoundationPartialAndForeignStage`). Integrate the
   actual reviewed shared observer and assert absent/empty staging, bounded
   temp-only, active/temp exact-length/equality, absent/partial/complete slots
   including empty evidence, invalid/unknown/special entries, and queue/base/
   request/completed-head linkage parity. Current direct fixture checks are not
   shared/native staging acceptance.

2. **Completed/staged native Audit, archive Export/Verify and statuses.**
   Handoff at `fixture_transaction_test.go:608`
   (`TestTMV0009_AS11_FixtureFoundationFiveOperations`) after publication,
   `fixture_recovery_test.go:1031`
   (`TestTMV0009_AS11_FixtureFoundationFreshRecovery`) at completed-head
   leftovers, and `fixture_transaction_fault_test.go:173`
   (`TestTMV0009_AS11_FixtureFoundationEveryPublicationPrefix`) at interrupted
   prefixes. Assert strict Audit and export/verify with persistent empty staging,
   actual committed payload/all linked orphans without staging bodies, and exact
   missing-head/linked-next/fork outcomes with zero pre-delivery archive output
   on refusal. The existing seed audit at `ftSeed` does not cover these groups.

3. **Production present-empty-D request lookup.** Handoff beside
   `fixture_transaction_fault_test.go:605`
   (`TestTMV0006_AS03_FixtureFoundationReplayOutcome`) and the empty-D KEEP row
   at `fixture_transaction_test.go:608`. After the reviewed J3-02 capture repair,
   add actual native recorded/absent request lookup, explicit intent agreement
   NOT_OBSERVED, strict Audit INTENT_DIVERGED, corrupt/missing private index
   refusal and source purity. The finite canonical helper is not this evidence.

4. **Representative native accounting/retry/purity integration.** Handoff at
   `fixture_recovery_test.go:228` (`ftDisk.cost`) and
   `fixture_transaction_fault_test.go:536`
   (`TestTMV0014_AS10_FixtureFoundationCapacityAndCleanupDebt`), plus the prefix
   callback in `TestTMV0009_AS11_FixtureFoundationEveryPublicationPrefix`.
   Exercise native staging D versus archive F accounting and retained orphans;
   unchanged-head stage/orphan changes, same-size/restored-mtime edits, empty
   root/parent replacement and listed-entry disappearance; exact four-attempt
   retry/error behavior, stable errors, sticky archive cleanup/stream-cap
   failures, no stale output, and byte/mode/mtime/lock/directory purity. Cite
   J3-02's exhaustive unit evidence while adding representative combined-history
   native assertions here. Current fixture scans/costs do not substitute for it.

5. **Final reviewed integration, combined gate/review/intent/docs.** Parent must
   integrate/rebase on reviewed-final J3-02, close groups 1–4, obtain fresh exact
   combined-diff independent review, run its one combined canonical gate and
   intent/evidence check, and record final SPEC/ROADMAP/BUILD-LOG evidence. This
   report and the three uncommitted new files are the handoff, not that closure.

No callable writer, issuer/admin grant, hostile-editor CAS, runtime, real-queue
permission, cutover, restore, GP/performance result, merge, TCP-02 completion or
J3-03/04 completion is claimed.
