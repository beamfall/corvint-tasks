# Corvint Tasks roadmap

Human-owned until `corvint-tasks init` (TCP-02) creates this repository's own native queue and the
`docs/SPEC.md` §5.4 authority switch is recorded; only then does this file become a generated
view. Status words here are delivery labels under `docs/SPEC.md` and decision 0001; none claims a
qualification. The experimental support candidate passed the final parent-supervised
`make verify` and fresh independent R4 review. See [integration evidence](reviews/2026-09-06-experimental-support-integration.md).
Experimental J2b/J3-01 and J3-02 also passed their parent gates and fresh reviews.
J3-03/04 fixture integration passed its combined gate and fresh review; see
[fixture delivery](reviews/2026-09-06-fixture-delivery.md).
G1..G6 and GP remain unqualified; earlier G0 review cleared bounded support work only.

Owner constraint (2026-09-06, exact): "I want you to ensure that this functionality doesn't slow
down the main go binary speed" → SPEC TM-V0-027 and gate GP (§9.4). No Corvint Go source is
changed by TCP-00, TCP-01 or TCP-02; TCP-03 and TCP-06 may not be promoted without GP.

| Ticket | Priority | Depends on | Scope (SPEC references) | Status |
|---|---|---|---|---|
| TCP-00 | P0 | none | Contract: schemas, limits, storage protocol, transition table, acceptance tables, adoption protocol, performance gate (`docs/SPEC.md`, decision 0001) | repair rounds R1, R2 and R3 applied (R3: Gate A cleared TCP-01 for build with no HIGH findings; its one MED documentation repair — §5.2 temp cleanup scope, §6.4 non-`EEXIST` link failures, AGENTS invariant 3 — is applied with TCP-01); G0 recorded as plan-cleared for TCP-01, pending the limited repair review of these process paths; no runtime qualification; G0 NOT_RUN |
| TCP-01 | P0 | TCP-00 | Native ticket schema (`revision`/`acceptanceRevision`), validation, intent projection and primary-worktree resolution, command-result envelope, read verbs, export/verify (TM-V0-002..008, 022; AS-01..AS-10, AS-35, AS-36) — `ADOPT_FILE` per SPEC §3.3 composition; `archive export` as a TM-V0-008 read emitting a stream. Complete without journal wiring | experimental support reviewed: combined verification and fresh independent R4 review passed; pure mutation library remains non-durable |
| TCP-02 | P0 | TCP-00 | Authority resolution, filesystem qualification, lock, journal commit/redo/chain (link-in commit, chain bounds), barrier scopes, reservations, liveness, `lane-leader` bootstrap, init/admit/cancel/retry/resume/reconcile/pause/unpause/drain, `archive restore` (TM-V0-001, 009..013, 016, 017, 022, 024; AS-11..AS-19, AS-27..AS-29, AS-33, AS-34, AS-36); **prerequisite before any durable producer is promoted: the SPEC §3.5 aggregate budget accounting (B4) with saturation and restore+unpause regressions** | partial, authority support reviewed: A1..A5 repaired; final combined gate and independent review passed; experimental J1 ledger read/audit passed the parent full gate and fresh independent review (2026-09-06-j1-integration.md); J2a pure archive encoding-size mechanics passed the combined parent gate and fresh independent review (2026-09-06-j2a-integration.md); J2b pure transaction/capacity/recovery and J3-01 private fixture primitives passed the parent combined gate and fresh separate reviews; J3-02 native staging/empty-D reads passed their parent gate and fresh review; J3-03/04 test-owned transaction/recovery/native-read integration passed the combined parent gate and fresh review of all fixture changes (SPEC §5.7); callable writer, local-operator issuer and `atm init` delivered under decision 0003 (2026-09-07): `internal/store` applies a plan in the §5.2 order with the redo rule, and an initialized store answers reads instead of `UNINITIALIZED` (`TestTMV0008_AS11_InitMakesTheStoreReadable`); hostile-editor CAS, reservations/proc, staging-descriptor recovery, restore and every mutation verb other than `init` held; TCP-02 incomplete; GP NOT_RUN |
| TCP-02b | P0 | TCP-01, TCP-02 | Serial wiring of ticket mutations to the journal; AS-02/03/05 rerun against the real journal | delivered 2026-09-07 (decision 0004): one `MUTATE` transaction operation carries any `taskman-mutation/0` envelope, so all fourteen `ticket` mutation verbs commit through the §5.2 writer and no longer answer `NOT_RUN` (SPEC §5.7.2). AS-02 (`TestTMV0005_AS02_CreateCommitsATicket`, `TestTMV0005_AS02_RefineChainsFromTheCommittedRevision`, `TestTMV0005_AS02_StaleExpectedRevisionLeavesTheStoreByteIdentical`), AS-03 (`TestTMV0006_AS03_IdenticalRetryReplays`, `TestTMV0006_AS03_SameRequestIDDifferentBytesConflicts`) and AS-05 (`TestTMV0004_AS05_HoldsArchiveAndRestore`) rerun against the real journal. Also delivered because a second transaction needs them: the §5.2 redo rule's pre-state case, post-commit redo of a pending receipt (crash point C2), and CAS replacement of intent projections. Attempt liveness still uses the zero-attempt oracle; `reconcile intent`, the administrative verbs, `archive restore`, reservations/proc and the runtime remain held; GP NOT_RUN |
| TCP-03 | P0 | TCP-01, TCP-02 | Corvint native queue adapter and `taskman-plan/0` priority-first planner; WQO/0 unchanged; frozen perf baseline before the first Corvint edit (TM-V0-015, 027; AS-21, AS-31, AS-32) | not-started; GP NOT_RUN |
| TCP-04 | P0 | TCP-01, TCP-02 | Runtime registry, capability profiles with probe evidence, supervisor and ack writer, budgets, cancellation (TM-V0-014, 018; AS-16, AS-20, AS-32) | not-started |
| TCP-05 | P0 | TCP-04 | Gate, review, docs and manifest lanes; reducer derives obligations (TM-V0-019..021; AS-22..AS-25) | not-started |
| TCP-06 | P0 | TCP-01, TCP-03, TCP-05 | Batched import dry-run/apply, authority switch, roadmap projection, one real ticket in explicit-reviewer serial mode; GP rerun (TM-V0-023, 027; AS-26, AS-27, AS-31) | not-started; GP NOT_RUN |
| TCP-07 | P1 | TCP-03, TCP-06 | Bounded fanout: `maxActiveAttempts` > 1; contention envelope measured (AS-30, AS-32) | not-started |
| TCP-08 | P1 | TCP-05, TCP-06 | Expert routing evaluation under ATM §7 thresholds, automatic reviewer selection | not-started |
| TCP-09 | P2 | TCP-00; operative use TCP-05, TCP-07 | Beamfall adapter and recorder under its own owner; foreign execution disabled until qualified | not-started; no owner (SPEC §10 U2) |

## TCP-01 experimental support inventory (review passed, 2026-09-06)

The current execution results and remaining holds are in the integration evidence above.
The following inventory records authored scope, not delivery. Owned files: `internal/wire`,
`internal/ticket`, `internal/intent`, `internal/snapshot`, `internal/archive`,
`internal/cli`, `internal/fixture`, `cmd/corvint-tasks`.

Authored scope: TM-V0-002 (AS-01, AS-10 wire and limits), TM-V0-003
(record chain and relationship rules), TM-V0-004 (status/eligibility separation, MANUAL
completion, AS-06), TM-V0-005 structural checks (missing dependency, cycles direct/transitive/via
gate, case-fold duplicates, invalid priority, unsupported version; AS-04), TM-V0-007 publication
and divergence as pure functions, TM-V0-008 read protocol, read verbs `ticket
list|search|show|blockers|export`, `queue status`, `roadmap`, `gate list|show` with the exact
per-verb inputs and items frozen in SPEC §3.3 (AS-08; journal facts reported `NOT_OBSERVED`),
pagination and explicit truncation (AS-08), byte/mode/mtime immutability of every read on
uninitialized, pending-redo and forked stores (AS-07, TCP-01 half), TM-V0-022 `archive export`
and `archive verify` (AS-09 export→verify half; AS-36 N4b, N4d, N4e, N4f) including the B4
archive capacity freeze (AS-10: parser boundary, >10,000-member export and verify).

Not implemented in this slice: `receipt show|audit|replay`, `config show`, `attempt show`,
`plan preview` (they need journal, pinned or plan data that no TCP-01 reader can observe) and
every mutation verb; each answers `NOT_RUN` with a warning until TCP-02/TCP-02b. The fourteen
`ticket` mutation verbs were wired by TCP-02b on 2026-09-07 and now commit; the administrative
verbs still answer `NOT_RUN`. §3.3 mutation
envelope decoding, role matrix, `expectedRevision` and post-record computation as a pure Plan
(AS-02), request-id replay/conflict against an explicit index (AS-03) and `ADOPT_FILE`
composition N3a..N3i (AS-35) are authored under `internal/mutation`. Status: a focused
`go test -count=1 ./internal/mutation` passed before the 2026-09-06 repair; the independent
review then found six MED items (ADOPT of a DRAFT gaining criteria, serial allocation wedged by
an explicit token, REOPEN leaving a completion on an OPEN record, an ADOPT_FILE digest over
file bytes only, a nil request index failing open, and adoption minting a revision on a
tombstone), all six repaired with counterexample regressions
(`TestTMV0007_AS35_DraftOpensByAdoption`, `TestTMV0005_AS02_SerialAllocationSkipsOccupied`,
`TestTMV0004_AS05_HoldsArchiveRestoreReopen`, `TestTMV0007_AS35_RequestDigestBindsActorAndTarget`,
`TestTMV0005_ContextInputsRequired`, `TestTMV0007_AS35_TombstoneRefusesAdoption`) and the
§3.1 completion rule tightened in `internal/ticket` (completion iff effective `COMPLETED`).
The repaired pure library passed the final combined verification and independent review; the package
plans and never commits, the mutation CLI stays absent, and no runtime qualification is claimed. Attempt, receipt replay and planner surfaces are later slices.

Contract blockers B1..B3 were resolved on 2026-09-06 as reversible implementation
clarifications under the TCP-01 experimental freeze (SPEC and witnesses changed in the same
repair; no repository, queue or ticket identity limit changed; no owner decision was needed):

- B1 (§2 vs §3.1): §3.1 now types the four token budget fields `Size` and `turns` /
  `wallClockMinutes` `Count`, as §2 already said; `intent.DecodePolicy` was already on that
  reading. Witness `TestTMV0002_AS10_LaneBudgetPrimitives`.
- B2: `head.primaryWorktree` and the archive manifest's `primaryWorktree` are `PathText` (§2:
  absolute path, 1..4096 UTF-8 bytes, controls and hostile code points refused), separate from
  `Identifier`. Witnesses `TestTMV0002_AS10_PathTextBoundaries`,
  `TestTMV0002_AS01_ReceiptAndHeadSchemas`.
- B3: §3.3 now says `snapshot` is null for `UNINITIALIZED` and for the commands that probe no
  store (`help`, `version`, usage errors before any read, `archive verify` of a stream).
  Witnesses `TestTMV0008_AS07_HelpAndVersion`, `TestTMV0008_AS07_ReadsLeaveStoreByteIdentical`.

- B4 (§1 vs §3.5), closed on 2026-09-06 as an unverified implementation freeze (SPEC §1, §3.5;
  the flat closed `taskman-archive/0` schema is unchanged): archive `files` ≤2,100,000,
  `manifest.json` ≤768 MiB (the proposed 512 MiB does not hold 2,100,000 largest valid
  entries; arithmetic in §3.5 and `internal/wire/limits.go`), export scan bound files + 4,096
  counted separately from the exported-file count, and the archive parser node cap
  4 × files + 14 reached only through the named opt-in entry `archive.ParseManifestDocument`
  (`wire.ParseWith`); depth 24, ordinary arrays 10,000 and nodes 250,000 are unchanged, and
  `wire.Parse` is untouched. The 1,000,000-receipt / 4 GiB admission thresholds are preserved.
  Witnesses `TestTMV0022_AS10_ParseWithWidensOnlyTheNamedArray`,
  `TestTMV0022_AS10_ExportEntryBoundDuringListing`,
  `TestTMV0022_AS10_ExportBeyondTenThousandMembers`,
  `TestTMV0022_AS10_VerifyBeyondTenThousandReceipts`. The 768 MiB manifest ceiling was
  explicitly authorized by the prior packet as the arithmetic-adjusted cap for 2,100,000
  entries; it is retained, the entry bound stays at 2,100,000 (not lowered to 1,435,000), and
  the pending owner question is closed. The figure is a hard cap only, not a memory or
  performance qualification: worst-case flat-manifest RSS is NOT_RUN, no low-overhead claim is
  made, and every durable-writer capacity, terminal/recovery headroom and admin escrow gate
  below remains held.

TCP-02 WRITER prerequisite (explicit, SPEC §3.5 "TCP-02 writer prerequisite"): the archive
entry bound, manifest byte cap, archive total, journal saturation thresholds and evidence cap
are independent hard aggregate budgets. Before every durable producer (receipt, post file,
evidence, pinned, bootstrap, acknowledgement, restore) the writer proves cost(proposed store) +
reserved remaining terminal/recovery costs + reserved admin restore/unpause/prune costs fit the
file, manifest, exact tar (headers, PAX, padding, end marker), receipt count/bytes and evidence
caps; unique retained paths and content are counted once; reservations are taken before work
and released only when durably consumed; headroom R/B is derived from the bounded §6.2/§6.4
recovery paths, never invented; a store already over a writer cap is refused, never truncated
or rewritten, and stays exportable. Every durable-writer promotion is held until that
accounting and its saturation and restore+unpause regressions exist. TCP-01 code writes no
native state; the prerequisite remains open and does not establish completion.

Repair applied 2026-09-06 after the first `make verify` run (Go 1.27.0): the wire test source
carries no literal NUL or BOM (hostile bytes are constructed), `make fmt` propagates a gofmt
syntax failure, TM-V0-008 re-probes a body failure before reporting it (N4d), store directory
enumeration is chunked and bounded before materialization (§3.1, §3.5), and B1..B3 above. A
second focused run showed one failure, `TestTMV0002_AS01_EncodeRoundTrip`, whose positive
witness placed a non-ASCII key before `z` (UTF-8 byte order puts `z` first); the witness was
corrected and the unsorted form kept as a negative case. The later read-subset run passed; the combined candidate then failed T1 before vet, as
recorded in the candidate review.

Agent memory carriers live in `docs/agent-memory/` under the user's convention (`bugs.md`,
`fixes.md`, `tests.md`, `optimizations.md`, `ideas.md`, `questions.md`; pending work only).

## Milestones

| Milestone | Requires | Meaning |
|---|---|---|
| M0 | TCP-00 reviewed | accepted buildable contract |
| M1 | TCP-01, TCP-02, TCP-02b, G1 | durable searchable inventory, no dispatch |
| M2 | TCP-01..06, G1..G4, GP, owner cutover record, `QUALIFICATION` receipt | one real Corvint ticket, explicit reviewer, capacity one |
| M3 | TCP-07, G5, GP rerun under `MAX_ADMITTED` | bounded fanout |
| M4 | TCP-08, G6 | automatic reviewer routing |
| M5 | TCP-09 per foreign queue | individually qualified foreign queues |

## Working rules

- At most two builders (TCP-01 and TCP-02 in parallel on the SPEC §9.1 packages) and one fresh
  read-only reviewer per round. Shared wiring, CLI routing and policy are serial.
- A ticket is done only when its AS scenarios exist as named Go tests, pass under
  `make verify`, and a fresh reviewer has read the diff. Until then its status stays
  `experimental`.
- Nothing here authorizes real-queue cutover, foreign writes, publication, or history deletion.

## TCP-02 J1 experimental read/audit increment (2026-09-06)

Plan-cleared revision 3 plus parent steering implemented serially in an isolated candidate.
Shared snapshot receipt/request/INIT codecs, streamed native ledger audit, canonical latest
projection checking, error-safe streamed RequestIndex and pure UNPAUSE deletion classification
are authored. Focused journal/snapshot/mutation and affected archive tests pass; see
[builder evidence](reviews/2026-09-06-j1-builder.md) for exact scope and results.
Parent canonical `make verify` and fresh independent review passed; see the J1 integration report. No CLI receipt
surface, writer, lock/redo, runtime, real queue, restore/prune, Corvint source edit, network or
qualification is added. Structural consistency is not historical acceptance. Archive
historical-semantic integration is NOT_OBSERVED because its retained semantic input omits
required blobs. J2b aggregate inventory/prefix/escrow mechanics and J3 writers remain HELD;
TCP-02 is not complete.

## TCP-02 J2a pure archive encoding-size increment (2026-09-06)

Gate-A-cleared encoding-only mechanics are authored in the isolated candidate. The exported
`archive.MeasureManifestEncoding` reports exact file count, payload, canonical manifest bytes
and complete tar bytes for a supplied typed manifest while retaining only bounded per-entry
temporaries and uniqueness state. Production and measurement share the unchanged fixed tar
header constructor; focused parity, padding, PAX, validation, cap, overflow and allocation
witnesses pass. See [builder evidence](reviews/2026-09-06-j2a-builder.md).

This is J2a partial status only. It does not establish a complete inventory, validate content
against digests, calculate scan/evidence/journal/transient costs, reserve terminal/recovery/admin
headroom, admit work, write or restore state, or qualify filesystem/runtime/performance behavior.
An empty files list is measurable but is not a verified store archive. Parent combination,
canonical gate and fresh independent review passed; see [integration evidence](reviews/2026-09-06-j2a-integration.md). J2b aggregate inventory,
publication-prefix and escrow mechanics and every J3 writer remain HELD; TCP-02 stays unchecked
and incomplete; GP is NOT_RUN.

## TCP-02 J2b experimental pure transaction/capacity model (2026-09-06)

Implementation candidate in `internal/transaction` plus the single shared descriptor codec
in new `internal/snapshot/stage.go`, based on accepted J1/J2a b67ed0d, both selected Gate A
corrections and the parent-approved shared-codec layering amendment. Owner decision 0002
selects fixture-only first delivery. Five administrative templates, proposed taskman-stage/0,
immutable inventory costs, abort-to-UNPAUSE closure and existing-receipt redo are modeled
without I/O. Hard caps stay 1,000,000 receipts/4 GiB with one UNPAUSE reserved inside them.
The original gate-green candidate failed fresh review on abort-to-UNPAUSE divergence and
cross-operation replay-conflict ordering. Both are repaired with causal-red and focused-green
regressions; parent hash verification, final combined gate and fresh review of the exact repair
passed. See [foundation evidence](reviews/2026-09-06-fixture-foundation-integration.md). J2b is experimental, **TCP-02 remains incomplete**. No archive-layout acceptance,
CLI/writer/runtime/admission, restore, real queue, upstream acceptance or qualification.
J3 and all existing runtime/restore/GP holds remain; SPEC §5.5 and the J2b builder report
record the exact mechanics and evidence limits.

## TCP-02 J3-01 private fixture filesystem increment (2026-09-06)

The narrowed Gate A mechanism is authored on accepted J2a `b67ed0d` in new
`internal/authority/fixture_session*.go` files only. Package-authority tests are its sole
constructors and own disposable repositories. Session/Lock closure coordinates through the
existing private mutex; no existing Go source, safeopen file, LinkIn behavior, runtime or
exported writer entry point is changed. SPEC §5.1.1 records exact roles, bounds, lifecycle,
mount observations, after-effect errors and six finite witness categories. Focused Darwin/
arm64 authority tests, including `-race`, pass. Linux/amd64 and Windows/amd64 test binaries
cross-compile; execution on those platforms is NOT_RUN. See
[builder evidence](reviews/2026-09-06-j3-01-builder.md) and [build log](BUILD-LOG.md).

This is fixture mechanism evidence only. Parent combined canonical gate and fresh exact-diff
review passed; see the foundation integration evidence. The owner-selected fixture-only delivery boundary is retained (parent
records decision 0002). J2b remains a separate parent-owned package increment; this slice
imports no unfinished symbols. J3-02 staging acceptance, J3-03 sequencing/recovery, J3-04
reconciliation and J3-05 issuer/callable writer remain held, as do existing-intent CAS/C6,
capacity/reserves, runtime, real queues and power-loss qualification. TCP-02 stays incomplete;
GP is NOT_RUN.

## Current accepted fixture foundation

J2b repair and J3-01 passed their combined parent gate and fresh separate reviews.
[Integration evidence](reviews/2026-09-06-fixture-foundation-integration.md) supersedes
earlier candidate-pending status wording for these slices. J3-02 is also accepted below;
J3-03/04 integration is authored pending its combined gate and fresh review.
Decision 0002 holds the callable writer/issuer/CAS/CLI/runtime boundary; TCP-02 is incomplete.

## TCP-02 J3-02 bounded read-only staging candidate (2026-09-06)

Implemented the selected Gate A amendments on the accepted J2b repair/J3-01 foundation:
shared strict structural stage observer, native retained directory observations, staging-byte
and orphan inventory rechecks, one four-attempt owner per journal/archive read, sticky export
resource errors, and narrow present-empty physical-ticket request lookup. Source.List's own
directory observation is mandatory; no fallback qualifies missing metadata. Stage files count
in scan D and never payload F. Matching completed leftovers consume actual receipt/request
proof and matching operation/kind plus timestamp; no full-plan or cleanup claim follows.

Focused snapshot/journal and selected archive regressions pass. See
[builder evidence](reviews/2026-09-06-j3-02-builder.md) and SPEC §5.6 for exact witnesses.
Parent owns source freeze, one canonical gate and a fresh independent exact-diff review;
these remain pending. No TCP-02 completion tick, Git commit/merge, CLI mutation, real queue,
writer/issuer/CAS, runtime, restore, physical power-loss or GP qualification is included.

## J3-02 accepted experimental read increment

Parent canonical gate and fresh independent review PASS; see [J3-02 integration](reviews/2026-09-06-j3-02-integration.md). This supersedes the preceding J3-02 candidate gate/review-pending text. J3-03/04 combined fixture evidence subsequently passed; TCP-02 remains incomplete.

## J3-03/04 accepted experimental fixture delivery

The three foundation test files now consume final reviewed J3-02 observation APIs, with a new
native-read integration test file. SPEC §5.7 maps all five deferred groups to exact TM/AS
evidence. The two selected test-harness repairs bind full model-time inventory before
preparation and exercise ticket branch refusal with a valid committed ticket token/control.
Native completed/precommit reads, empty-D index behavior, receipt-only INIT and combined
history staging/orphan accounting are authored. The 17 foundation suites remain intact.

See [focused builder evidence](reviews/2026-09-06-j3-fixture-integration-builder.md).
Parent combined canonical gate and fresh independent review of the complete fixture diff
PASS; see [delivery evidence](reviews/2026-09-06-fixture-delivery.md). This row does not tick
TCP-02 or qualify a callable durable writer. No
process-interruption or power-loss evidence was produced; issuer/admin binding, hostile-editor
CAS, CLI mutation, runtime, real queues, restore and GP remain held.
