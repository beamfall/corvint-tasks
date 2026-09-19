# J3-03/04 combined fixture integration — builder evidence

Date: 2026-09-06. Worktree: `/private/tmp/corvint-tasks-build-20260906/tcp02-j3-final-integration`.
Base: reviewed J3-02 `bc3b3dd7df3d383f0d546cc9e6ffe7e46a3d0d28`, plus the three frozen
foundation test files and two supplied provenance reports. Read AGENTS, all four conventions,
decision 0002, the full accepted integration plan/Gate A including appended selections and
scheduling, foundation handoff, final J3-02 builder/integration/review and final source APIs.
Optional `tcp02-j3-final-integration-steering.md` was absent when checked before finalization.

Status: five deferred integration groups implemented to the builder-owned boundary below.
**Parent combined freeze, canonical gate, fresh independent review of every fixture change
(including the foundation), and final intent/evidence acceptance remain PENDING.** The foundation
had preliminary inspection, not final independent code review. TCP-02 remains incomplete.
No durable-writer, historical acceptance, authority or performance qualification is claimed.

## Exact files and protected source

Owned source (three supplied NEW foundation tests modified, fourth NEW integration file):

| File | Final SHA-256 |
|---|---|
| `internal/authority/fixture_transaction_test.go` | `030fd83da3b5320b2675b673966636ad7bf32eafae6018cb5019d4acdc3371e3` |
| `internal/authority/fixture_recovery_test.go` | `818adc11c9db2b4be4442bb97ad013b4bd1eb2853631e326c19147bbe501d22b` |
| `internal/authority/fixture_transaction_fault_test.go` | `81b4faa693beaca2d39db742c3d15e3bdf7a8c839aa02b5f59549347bdea2529` |
| `internal/authority/fixture_read_integration_test.go` | `1c25147e83c1ad8ae6a396c697f0c911e0db7b58f8220a2609ab75b6a3602eea` |

Owned documentation: `AGENTS.md`, `docs/SPEC.md`, `docs/ROADMAP.md`, `docs/BUILD-LOG.md`,
`docs/agent-memory/tests.md`, `docs/agent-memory/fixes.md`, and this report. The two supplied
foundation/preliminary provenance reports were preserved. Other pending memory files were
read by title and left intact. No Corvint/Beamfall source, production Go file, other existing
test, codec/model/primitive/native-reader API or private hook changed.

All **99 protected Go files** match the pre-edit SHA-256 capture, including final reviewed
J3-02 and every pre-existing test. Capture: `j3-final-protected.json` under the scratch root
above. All 17 original foundation suite names remain; existing finite fault/prefix tests and
final inventory/cost equality remain present. No Git mutation, delegation, full canonical gate,
independent reviewer, parallel test coordinator, unsafe/linkname, race seam, test skip or
placeholder was introduced. Test commands ran sequentially without a PTY.

## Selected harness repairs

1. `ftHarness.model` associates each successfully modeled opaque Plan with a separately
   copied immutable inventory observed at that exact model boundary. `publish` compares full
   files/digests/lengths/directories before era issuance, staging creation or writes; head and
   affected preimage checks remain. It never reconstructs the expected inventory from the
   later observation. `FixtureFoundationStaleInventory` removes a linked orphan without head
   movement and requires refusal, no era/staging and exact source/metadata purity; copied
   accessor mutations do not change the saved base. Model associations are discarded at era
   issuance and restart, so recovery cannot obtain original Plans through the harness.
2. `FixtureFoundationSourceLifetime` selects the actual artifact targeting `ftTicket`, retains
   its valid token through receipt link, and reads exact physical D as preimage. Main-branch
   replacement succeeds; wrong-branch replacement refuses with unchanged bytes/metadata.
   The other four original lifetime cases remain. The old HEAD-slot branch case was replaced,
   not credited as branch coverage.

Both repairs have causal-red evidence from scratch **test-file-only** overlays:
`j3-final-branch-causal-red.log` fails `ticket-branch-other` with “invalid actual bytes accepted”
when only branch validation is disabled; the main control passes. `j3-final-inventory-causal-red.log`
fails with “stale inventory started an era” when only original-inventory comparison is disabled.
Each command exits 1 as expected. Production/worktree source was never replaced by an overlay;
subsequent green commands use no overlay.

## Five-group TM/AS acceptance mapping

The exact full test names are also recorded in SPEC §5.7. Names below are in package authority.

| Deferred group | Named acceptance evidence |
|---|---|
| 1. Shared observer at reopen/read | `TestTMV0008_AS27_FixtureNativeStageParity`: absent/empty, bounded arbitrary temp (empty/malformed/max/over), active/temp N−1/equal/different/N+1, assigned absent/empty/partial/full/wrong/long slots, invalid names/active/unassigned/symlink/FIFO/directory; shared observer versus fixture/native outcomes. `TestTMV0002_AS27_FixtureNativeAssignedEmptyEvidence`: present-empty declared evidence is complete. `TestTMV0008_AS11_FixtureNativeStageBindings`: actual queue/base/head/receipt/request linkage, request ID/digest, operation mapping and original timestamp. |
| 2. Native completed/precommit/interrupted histories | `TestTMV0022_AS11_FixtureNativeCompletedHistories`: strict Audit and actual tar Export/Verify for all five operation shapes with completed leftovers, then persistent empty staging; exact final native manifest/payload/encoding costs. `TestTMV0009_AS11_FixtureNativeInterruptedStatuses`: receipt-only prefixes for all five operations, headless temp, earliest INIT with no physical queue/head (also without staged queue), missing/corrupt actual genesis VERSION blob, linked-next/fork statuses and zero refusal output. `TestTMV0022_AS35_FixtureNativePrecommitIntent`: valid divergent D and empty D, retained actual linked orphan, abort purity; no D→C substitution. `TestTMV0007_AS35_FixtureNativeOrdinaryIntentRefusal`: malformed D refuses ordinary intent loading and strict Audit. |
| 3. Native empty-D request adapter | `TestTMV0006_AS35_FixtureNativeEmptyIntentIndex`: actual FOUND/ABSENT after committed PAUSE, complete original outcomes, NOT_OBSERVED intent agreement, strict Audit INTENT_DIVERGED with no claimed projection agreement, missing/corrupt private-index refusal and source/lock/directory purity. |
| 4. Combined accounting and read boundary | `TestTMV0008_AS10_FixtureNativeAccountingAndSequentialReads`: ticket-bearing seed→PAUSE→UNPAUSE→KEEP→ADOPT history plus two modeled precommit KEEP evidence-link prefixes; empty staging directory adds one D and zero F; children add D only; two actual linked orphans add F and persist after abort; exact tar bodies/full manifest equal actual inventory. Sequential native reads reflect staging/orphan changes and same-size restored-mtime staging bytes without head movement. Existing exhaustive in-read matrices are cited below, not duplicated or accessed via exposed hooks. |
| 5. Final rebase and evidence closure | Final J3-02 APIs consumed directly (`snapshot.ObserveStage`, `StageObservation.Bind`, `journal.Native`/Reader/RequestIndex and native archive Export/Verify). SPEC/ROADMAP/BUILD-LOG/current AGENTS and pending memory align delivered slices and holds. Protected-source checks and focused evidence below are complete. Combined canonical gate and fresh review remain parent-owned PENDING; this report does not self-approve them. |

`observe` retains stable preflight, pending physical PRE/POST mixtures and completed cleanup
as distinct uses. Actual receipt/blob/request/genesis proof precedes fallback to an absent
physical queue's validated genesis afterimage; staged proposed queue bytes never bind it.
Reopen uses shared structure **and** retained no-follow file identity/size/digest checks to
mint fresh session-local tokens; partial/temp observations remain cleanup-only. Shared
observations never grant cleanup or execution. Recovery takes no Plan, Model, original
invocation or cached afterimages. Existing conservative prefix domination and unsynced
name/shrink debt remain; exact cost/inventory equality occurs only after durable cleanup.

All native checks snapshot the entire disposable repository (including lock/source/directories)
before and after, preserving bytes, modes and mtimes. Successful exports compare actual tar
bodies and complete manifest entries to independently captured current inventory, preserving
all linked orphans and excluding staging bodies; Verify still means chain/digest integrity.

## Existing exhaustive native evidence (not rerun or duplicated here)

Accepted final J3-02 integration/review and canonical gate:
[2026-09-06-j3-02-integration.md](2026-09-06-j3-02-integration.md),
[final review](2026-09-06-j3-02-review.md), [builder evidence](2026-09-06-j3-02-builder.md).
The unchanged private-hook matrices own:

- `internal/journal/stage_read_test.go`: `TestTMV0008_AS36_JournalStageMovementFourAttempts`,
  `TestTMV0008_AS36_NativeDirectoryObservationLifetime`,
  `TestTMV0008_AS36_JournalStageReadErrorRecheck`,
  `TestTMV0008_AS07_JournalFileClosureErrorsSurviveMovement`.
- `internal/archive/stage_read_test.go`: `TestTMV0022_AS36_StageAndOrphanMovementFourAttempts`,
  `TestTMV0022_AS36_StageBodyErrorRecheck`, `TestTMV0022_AS10_StageScanAndPayloadBounds`,
  `TestTMV0022_AS36_StageStreamAndCleanupSticky`,
  `TestTMV0022_AS36_ArchiveCoordinatorOneShotAndReprobe`,
  `TestTMV0022_AS36_ArchiveStageValidationAndFailedAfterCapture`,
  `TestTMV0022_AS07_ArchiveFileClosureErrorsSurviveMovement`,
  `TestTMV0022_AS10_StageExcludedFromFileCapBeforeBodies`.

These are the exhaustive in-read staging/orphan movement, restored-mtime byte edits,
root/parent replacement/disappearance, exactly four attempts, stable errors/failed recapture,
one-shot inner probes, sticky close/unlink and stream-cap failures (including later smaller
streams), and D/F boundary evidence. The combined fixture test proves representative real
native reads/accounting and sequential freshness; it does not claim those private-hook
matrices were independently recreated in authority.

## Focused execution

Go 1.27.0, Darwin/arm64. Every Go test used:

```sh
GOTOOLCHAIN=local GOCACHE=/private/tmp/corvint-tasks-build-20260906/go-cache \
GOMAXPROCS=2 GOFLAGS=-p=2
```

Commands after that environment prefix (all logs under the scratch root):

| Command | Result / evidence |
|---|---|
| `go test ./internal/authority -run '^$'` | compile-only PASS; no tests claimed |
| `go test ./internal/authority -run 'FixtureNative\|FixtureFoundation(SourceLifetime\|StaleInventory)$' -count=1 -json` | PASS, 15.613 s, 11 suites / 66 test-subtest passes, no skips; `j3-final-native-focused-green.jsonl` |
| `go test -overlay .../j3-final-branch-overlay.json ./internal/authority -run 'FixtureFoundationSourceLifetime/ticket-branch' -count=1` | expected FAIL (exit 1), causal guard witness described above |
| `go test -overlay .../j3-final-inventory-overlay.json ./internal/authority -run 'FixtureFoundationStaleInventory$' -count=1` | expected FAIL (exit 1), causal inventory witness described above |
| `go test ./internal/authority -run 'Fixture(Foundation\|Native)' -count=1 -json` | PASS, 120.237 s, **27 suites / 393 test-subtest passes**, zero failures/skips; `j3-final-all-affected-focused.jsonl` |
| `go test ./internal/authority -run 'FixtureNativeStageBindings$' -count=1 -json` | final PASS, 3.408 s, 1 suite / 10 test-subtest passes; `j3-final-actual-binding-focused.jsonl` |
| `go test ./internal/authority -run 'FixtureNativeAccountingAndSequentialReads$' -count=1 -json` | final PASS, 3.354 s, 1 suite; `j3-final-accounting-focused.jsonl` |

The single affected-foundation rerun was warranted by the shared-observer and model-time
inventory changes on every publication/reopen path. It retains 17 original suites, adds the
new stale-inventory suite and nine native suites. After that run, the binding test was
strengthened to read the direct shared Bind receipt argument from the actual linked file
instead of the equal Plan accessor; its exact focused follow-up passed. The accounting setup was then tightened to obtain both linked orphans through modeled
precommit KEEP prefixes, replacing an intermediate unassigned-slot setup before native reads.
Its exact focused follow-up is recorded below. No foundation or full canonical gate was
repeated for these test-only strengthenings. Final owned-file gofmt
and `git diff --check` pass; all 99 protected Go hashes match.

Initial failures were test-builder assumptions, corrected without production changes: one
unused import; a directory child reports journal MALFORMED versus archive
UNSUPPORTED_FILESYSTEM; the first foreign-queue negative violated queue syntax and was
corrected to a syntactically valid foreign queue; nonempty malformed ticket bytes are not
rejected by archive's existing chain/digest-only verification. The final precommit export
fixtures therefore use valid divergent D and present-empty D, while a separate ordinary
intent-load test retains malformed-D refusal. No malformed D was replaced with C and no
ordinary intent validator was weakened. Failed logs remain `j3-final-focused-initial.log`
and `j3-final-native-focused.jsonl`; neither supplies PASS evidence.

## Remaining holds

No new underlying product contract is selected. Parent freeze/canonical gate/fresh review
and final intent/evidence acceptance remain PENDING, including final review of the foundation.
Callable durable writer, independently qualified issuer/actor/admin binding, hostile-editor
CAS, CLI mutations, live runtime/reservation escrow, real queues, restore and GP remain HELD.
The finite helper supports at most twelve receipts, 128 directory entries and 4 MiB captured
bytes; it is not a general production canonical-observation/recovery API. Archive verification
is not historical semantics, authorization or restore qualification. Unknown authority/runtime
coverage remains NOT_OBSERVED. Native Linux, physical saturation, RSS/timing and GP are NOT_RUN.
No crash child was used; returned errors run Go defers. Process-interruption and physical
power-loss evidence remain **NOT_PRODUCED**. No merge, promotion, real-queue cutover or
full-gate clearance is claimed; TCP-02 remains incomplete.
