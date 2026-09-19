# J3-02 — bounded read-only staging layout

Status: proposed implementation plan, not implementation or gate evidence. No candidate source changes, tests or delegation were performed. Owner decision `0002-fixture-only-first-delivery.md` accepts the fixture-only delivery boundary; it does not supply writer, actor, cleanup, filesystem or runtime authority.

Baseline inspected: `tcp02-j2a-final/internal/journal/{audit,source,records}.go`, `internal/archive/{layout,export}.go` and `internal/snapshot/probe.go`. Consume J2b's final reviewed shared descriptor codec only after its freeze; names below describe responsibilities rather than unfinished symbol spellings. Do not import transaction from archive: transaction already depends on archive sizing.

## Scope and ownership

Claim new `internal/snapshot/stage_observation.go` and `stage_observation_test.go` for a small pure staging-byte/layout validator using the shared descriptor codec. It owns no filesystem traversal, full-plan classification, capacity or cleanup decision. Do not duplicate or weaken the shared descriptor decoder.

Claim `internal/archive/layout.go`, `export.go`, new `stage_read.go`, `stage_read_test.go`, and the relevant existing archive tests. Claim `internal/journal/audit.go`, `records.go`, new `stage_read.go`, `stage_read_test.go`, and relevant existing audit tests. Claim `journal/source.go` and focused source/decorator tests for the concrete directory observation and retained native lifetime amendment below. `snapshot/probe.go` and its public Reader API stay unchanged. SPEC, ROADMAP and BUILD-LOG record the closed layout amendment, named evidence and J3-02 partial delivery only.

No authority/transaction builder files, writer operations, CLI mutation wiring, restore, new actor grants, resource admission, arbitrary intent overwrite, real-queue execution or historical semantic archive integration are included. No changes to the archive manifest schema, standard wire limits or existing legacy-temp rules.

## 1. Closed observable layout and bounds

Only the exact state-root directory `staging/` is added. Its only possible children are regular files named `active.json`, `active.json.tmp`, and `a00` through `a10`: at most thirteen children, no nested directory, alias, traversal, suffix variant, symlink, FIFO or other special file. An empty persistent staging directory is inactive; its existence alone is not an active transaction or REDO_PENDING.

Count the staging directory and every child against the existing cumulative state scan budget. Bound the child listing by min(remaining scan budget, thirteen), refusing excess during enumeration before opening/stat'ing excess entries. Unexpected names still consume the enumeration budget before refusal. The intent-tree budget remains separate as in the accepted journal capture. Export counts actual payload files against its file cap before opening bodies; staging children contribute no payload files.

Use retained no-follow parents and the existing nonblocking/no-follow file readers. Collect the complete bounded names/metadata before reading staging bytes. Bound active/temp descriptor reads by the shared maximum (currently 2,422 bytes); then decode the active descriptor before opening its payload slots. A slot is read only under the descriptor's declared size, which has already passed its per-operation cap. Recheck actual consumed size; stat is not the byte proof. Preserve present-empty versus absent semantics. Aggregate staging bytes are bounded by the finite operation artifact/descriptor caps, not by archive-body limits or caller-provided arbitrary sizes.

The pure staging validator receives the finite observed names and bytes, uses the single shared descriptor decoder, and reports structural observations only:

| Observed staging | Read interpretation |
|---|---|
| absent or empty directory | inactive; no preparation claim |
| only bounded `active.json.tmp` | incomplete descriptor preparation; never an active descriptor or ownership proof |
| valid active descriptor; assigned slots absent or shorter than declared | incomplete payload preparation |
| assigned slot exactly declared length and matching digest | complete slot bytes, not a complete authorized plan |
| empty slot with declared size zero and empty-content digest | valid complete empty raw bytes |
| malformed/oversized active descriptor, undeclared slot, overlong slot, or full-length digest mismatch | precise MALFORMED/LIMIT_EXCEEDED/JOURNAL_FORKED refusal |
| wrong queue/base/request binding | foreign or inconsistent staging refusal; no silent skip |

The accepted J2b temp-only rule intentionally permits bounded zero-length, partial or malformed `active.json.tmp` bytes. Therefore “malformed staging refuses” applies to active descriptors and invalid assigned data/layout, not to unpublished temp bytes. With active present, read active first, bound temp by the actual active byte length, and require byte equality if temp reaches that exact length; excess is LIMIT_EXCEEDED and full unequal bytes are JOURNAL_FORKED. A shorter temp is interrupted preparation, never replacement authority. Only without active may arbitrary temp bytes use the global 2422-byte bound. Do not import transaction's abortable/FullPlanPrepared classifications into readers or authorize cleanup from these observations.

## 2. Bind staging to the observed store without claiming a transaction is accepted

An active descriptor must name the observed queue and match a recognized head/receipt relationship. Existing nonfixture reads without staging remain unchanged; this increment's nonempty native staging recognition is limited to the fixture contract, with fixture content remaining a restriction rather than ownership proof.

- **Precommit:** a valid initialized head matches the descriptor's base sequence and last-receipt digest. No next/extra linked receipt exists. Archive may export the unchanged committed files and all linked evidence/pins, while omitting the transient staging children. A proposed head/receipt slot never replaces the actual head or journal in a read.
- **Missing head:** archive's existing snapshot probe stays UNINITIALIZED; neither a descriptor nor a complete slot initializes anything. Journal preserves its stronger existing distinction: absent head and no linked receipt1 is UNINITIALIZED; a fully validated linked genesis is REDO_PENDING. Staging alone never upgrades either case. A malformed/foreign stage may refuse, but cannot create a success or initialized-state claim.
- **Linked next receipt:** existing pending/fork precedence remains; export returns REDO_PENDING and no archive, journal validates and reports the linked receipt without redo. More than one next receipt/gaps remain JOURNAL_FORKED. A stage descriptor never excuses an invalid receipt inventory.
- **Completed head with leftover stage:** recognize only when the actual head and its actual last receipt match the descriptor's HEAD and RECEIPT artifact hashes; the receipt sequence/prev matches the descriptor base; queue/request ID and receipt-bound request digest match the descriptor. Reuse existing receipt/request codecs and bounded post/blob identity checks for this narrow linkage. This is a current-byte relationship, not execution, authorization, completion-of-cleanup or full-plan proof. Missing consumed slots can remain absent; present complete slots must still match. A stale/foreign descriptor that cannot meet these relationships refuses. The archive exports the actual current files; never staged afterimages.

Put the shared relationship checks in the pure observer only where both callers supply the same bounded bytes. Keep native byte acquisition in each package. No new generic I/O/capability framework or archive-to-journal/transaction dependency is required. If a required bound byte is unavailable, return explicit refusal; do not call a hash assertion proof. Existing full journal audit/request validation remains intact.

## 3. Actual state inventory and staging tuple

Each observation contains the actual state-root and visited-parent identities plus every bounded enumerated entry: exact relative name, type/mode, size, modification time and file identity. Compare identities with os.SameFile; do not invent a portable inode string. Root/directory identities participate even if empty. Worktrees retain the accepted policy of counting their root without traversing lane checkouts.

Also hash the exact consumed bytes of every present staging child, including partial slots and descriptor temp. Metadata alone is insufficient: same-length edits with restored mtime must change the staging tuple. Keep bounded descriptor/slot content hashes alongside metadata; do not retain body-sized data beyond validation. A deterministic internal tuple need not become a new wire profile or authority token.

Every linked orphan under evidence/pinned remains in the actual inventory and archive payload. No descriptor reference filter is allowed. New orphan/link/slot/directory names without head movement change the observation. Existing two-pass hashing of archived source bodies remains; inventory metadata is not a substitute for digesting consumed payload bytes.

Capture gathering and staging semantic validation are separate enough that a stable semantic error can be rechecked. Once a complete bounded observation exists, validate the stage and run the body, then capture again even if validation/body failed. Same tuple returns the original stable error; a changed tuple discards that attempt. Hard bounds, unsupported file types and genuine I/O failures still fail closed when no complete comparison is possible. Only positively observed movement, including disappearance/replacement of a previously listed entry, earns a movement retry; do not convert arbitrary I/O failure into absence.

## 4. One bounded retry owner per command

### Archive export

Use an archive-local coordinator with exactly four attempts. Inside each attempt invoke the existing snapshot.Reader with retries disabled (`Retries:-1`, currently one total attempt). This one-shot invocation still probes head/receipt slots/barrier/intent before and after the callback. Do not leave its default four attempts enabled inside the archive loop.

Within that callback: collect the bounded layout/staging observation, validate it and copy only the observed committed source-file list, then re-observe the actual state/staging tuple on both success and body-error paths. On tuple movement return SNAPSHOT_MOVED. The outer coordinator retries only movement, up to four total body attempts; existing header movement from the one-shot Reader follows the same route. Stable body errors remain verbatim. A re-probe that exposes REDO_PENDING/JOURNAL_FORKED retains the existing behavior instead of pretending the store is stable.

This solves the current corner case: snapshot.Reader.Same compares only head, barrier and intent; it returns a callback SNAPSHOT_MOVED immediately if those remain unchanged. The local coordinator supplies the necessary retry without a new generic callback API, changing every read command, or nesting sixteen attempts.

Keep archive's staging resource lifetime outside the retry coordinator. Cleanup before the next attempt and final deferred cleanup preserve the existing sticky joined cleanup errors; no moved attempt can erase a failed close/unlink or turn it into success. Reset attempt-local manifest, file list, counts and stream state explicitly. Preserve the existing sticky stream-cap refusal, unlinked external staging descriptor, verification-before-delivery, and partial-delivery error wording. No source state or intent file is created, modified or removed.

### Journal reads and request lookup

Extend the existing capture's actual inventory tuple with staging content hashes; reuse the current single four-attempt audit loop. Add exact `staging` to directory recognition, validate only the thirteen names through the staging helper, and exempt only validated staging entries from canonical projection equality/untracked-post checks. Do not teach PostBound that stage files are receipt post destinations, and do not add a blanket `staging/*` temp predicate that skips validation.

Run stage semantic/base validation inside the rechecked audit body so a changing preparation cannot become a stable corruption assertion. Preserve full chain/blob/request binding and the distinction between full Audit and request-only intent projection coverage. Empty staging does not set an active-preparation indicator; existing general UNINITIALIZED-remnant reporting must not relabel it as an active transaction. Ordinary snapshot-only ticket views make no staging/accounting claim and need no new scan or API.

## 5. Required finite witnesses

Use modest disposable fixtures, direct finite byte observations and existing hooks. Name tests with the relevant TM-V0-008/022, TM-V0-002/010 and AS-07/10/36 IDs; journal pending-genesis/chain cases retain their existing IDs.

1. Shared helper: absent/empty directory; bounded empty/partial/malformed temp-only input; valid active with every assigned-slot absence/partial/complete combination needed by the five operation shapes; declared empty raw evidence; malformed active/version/closed fields; bad names/nesting/duplicates; per-operation bytes and cap+1; full-length bad digest; foreign queue/base/request. Exact codec goldens remain shared.
2. Archive: stable precommit exports only current committed records, includes unreferenced newly linked evidence/pins, omits every stage child, and verifies under existing archive checks. An empty persistent staging directory does not block export. Missing head, linked next and fork fixtures keep their codes and zero pre-delivery stdout. Matching completed leftovers use actual head/receipt/request proof; stale or missing-proof descriptors refuse.
3. Accounting: stage directory plus every child counts against scan at limit and limit+1, including malformed excess names; stage bodies do not count as archive files/payload. Flood refusal occurs before excess-entry reads. All source orphan files still count and export.
4. Movement: mutate staging or add a linked orphan during copy while head/barrier/intent stay byte-identical. A one-time change retries then returns the new coherent inventory; mutation on every attempt yields SNAPSHOT_MOVED after exactly four callbacks. Include same-size/restored-mtime staging edits and empty-directory/parent replacement. No stable accounting claim survives them.
5. Body errors: inject validation/read failure plus a changed tuple and prove retry; the same failure with a stable tuple returns the original error. Exercise failed after-capture and pending/fork re-probes. Test the archive-local coordinator directly to prove one-shot inner reads, and retain header-movement regressions.
6. Resource errors: combine movement with archive staging close/cleanup and stream-cap failures; those failures remain sticky, never success. No stale attempt-local outputs leak. Journal and archive retain byte/mode/mtime purity across every refusal and success; lock bytes/mtime and source directory contents are untouched.
7. Safety/platform: deterministic no-follow directory/file replacement and FIFO witnesses remain bounded/nonblocking, foreign bytes preserved. No source cleanup or automatic init on malformed/missing data. No child processes are required; if any are added, the repository cleanup/interruption contract applies.

Completion requires fresh independent review of the exact final diff and root's single canonical gate after focused regressions. Keep result language limited to structural layout, snapshot observation and archive chain/digest checks; NOT_OBSERVED actor/history/runtime/capacity/cleanup authority and existing writer/restore/performance holds remain. No TCP-02 completion tick or merge authorization follows from this plan.

## Selected Gate A amendments (parent, 2026-09-06)

Both MED findings in `tcp02-j3-02-gate-a.md` are selected and binding; no bounds or observation claims are weakened. The active/temp rule above uses the final J2b actual-length/equality rule. Add N-1/Nsame/Ndifferent/N+1 and no-active parity witnesses.

Change internal journal Source.List to return a concrete result carrying the listed directory's own os.FileInfo plus Entries, obtained from the same listing handle. Explicitly update all existing fixture/decorator implementations (including tests outside journal if compiler identifies callers); these mechanical adapters are in scope. Missing directory metadata refuses; never substitute a synthetic verified identity by default. The directory's identity does not add a child or extra scan charge.

Use a small private per-attempt native read lifetime in journal, pinning state and intent roots and retaining/reusing listed no-follow parents through capture/body/recheck, with joined close errors on every exit. Native file reads use those retained parents. Re-observe named root/path bindings against retained identities to detect detached/replaced roots; archive uses its analogous local lifetime. No new public Reader option, generic capability framework or extra retry owner. Explicit nonnative Sources keep their synthetic byte-source contract and supply explicit directory observations, without a native host claim.

After successful bounded enumeration, disappearance/replacement during metadata or body acquisition is positively observed SNAPSHOT_MOVED, not absence of an optional root. Initially absent optional roots remain distinct. Unrelated I/O/type/limit errors fail closed. Add root replacement with identical moved children, empty-root/parent replacement, listed-parent-before-body replacement, missing directory metadata, initially absent versus disappeared, bounded FIFO safety and all-exit handle-closure witnesses. Keep exactly four attempts, source purity and sticky archive cleanup/stream-cap errors.

## Narrow existing reader repair required by empty-D fixtures

The parent confirmed an existing J1 request-only gap: `journal/audit.go` capture rejects all zero-byte physical intent at its size check before `requiredRead`; `journal/index.go` calls the same capture even with intent projection agreement disabled. Thus a valid ledger/index with a selected ticket truncated to present-empty D returns MALFORMED, although stable intent divergence and KEEP's empty raw D are supported contracts.

Within this claimed audit source, preserve present-empty physical **ticket** bytes in the bounded inventory/raw-intent observation (including empty-content digest), instead of rejecting them during capture. Missing and empty remain distinct. Full strict Audit still reports the physical projection disagreement as INTENT_DIVERGED against canonical C; the request-only audit continues full chain/blob/request/private projection checks and can return the verified recorded/absent index with intent agreement NOT_OBSERVED. Queue/policy and canonical ticket codecs stay strict; no empty canonical record or blanket ignored parse/read error is admitted. Add causal-red/green recorded+absent lookup, corrupt/missing request-index refusal, strict full Audit divergence, and bytes/mode/mtime purity witnesses. Keep actual read errors fatal. This is a narrow repair, not a new public canonical-observation mode.

## Selected completed-stage linkage correction

2026-09-06. Independent bounded Gate A found and confirmed a MED linkage gap in the draft `StageObservation.Bind`. See the appended correction in `tcp02-j3-02-gate-a.md`. The parent selects the following narrow repair within existing ownership.

For completed-head leftovers, compare descriptor operation to actual receipt kind using the accepted model's mapping: INIT->INIT, PAUSE->PAUSE, UNPAUSE->UNPAUSE, KEEP_JOURNAL/ADOPT_FILE->RECONCILE. Also require descriptor RecordedAt to equal the actual receipt RecordedAt. Mismatch returns JOURNAL_FORKED. Retain existing actual head/receipt/hash/sequence/prev/request/outcome proof. This is current-byte field consistency, not full-plan reconstruction, actor authority or a claim to distinguish historical KEEP versus ADOPT acceptance.

Concrete wrong-label witness: construct a codec-valid three-slot UNPAUSE descriptor referencing an actual completed PAUSE receipt/head/request (receipt is below UNPAUSE cap), carrying that request's actual digest. Existing hash/sequence/request checks alone accept this mismatched operation. Add a negative for that case and a descriptor-only timestamp change, plus positive mapping witnesses including RECONCILE for both reconciliation labels. Do not weaken codec or enlarge bounds; no new API or ownership required.

## Second selected correction — headless INIT before physical queue publication

Independent Gate A confirmed a MED regression in the draft. `records.go` calls queue-dependent `validateStage` before the receipt walk; that helper requires physical `intent/queue.json`. After a real INIT receipt link, that file may still be absent. Existing J1 accepts its null PRE while fully validating the queue afterimage from receipt1. The current genesis witness merely removes head from a fully materialized fixture and misses this prefix.

Defer the queue-dependent stage binding for headless linked genesis until the existing complete genesis chain/blob/request/projection audit succeeds. If physical queue is absent, use the queue afterimage obtained from that actual fully validated linked genesis (bounded and content-bound); NEVER use a staged proposed afterimage, descriptor assertion or cached fixture bytes. Reuse the existing audit/post decoder and preserve all its checks. An actual non-not-exist read failure stays a refusal. No public canonical-observation API or Source/Reader contract widening is needed.

Add causal-red/green native journal coverage for a complete linked INIT receipt and required immutable evidence, active INIT descriptor, no head and no physical queue (prefer no materialized posts). It must report REDO_PENDING only after the complete existing audit passes. Cover missing/corrupt required genesis blob, invalid queue/genesis, foreign/non-fixture binding and read-error refusal. No linked receipt remains UNINITIALIZED; staging alone never initializes. Preserve byte/mode/mtime purity. Archive remains UNINITIALIZED when head is absent under its existing snapshot probe.

See the second appended J3-02 Gate A addendum for exact rationale. This is selected within your existing source ownership and supersedes any earlier report that only covered fully materialized genesis. Re-read this steering before finalizing.
