# Build log

## 2026-09-15 — Corvint Tasks identity candidate

Decision 0006 renames the current product, executable, command source directory, and Go module to
Corvint Tasks, `corvint-tasks`, and `github.com/Beamfall/corvint-tasks`. The frozen `.taskman` state
layout and `taskman-*/0` protocol identities remain unchanged. `CORVINT_TASKS_ACTOR` is primary,
`ATM_ACTOR` remains a fallback, and a conflict is rejected before mutation-command payload or
repository reads and before any write.

Current source selects Go 1.27.1 with `GOTOOLCHAIN=local`. Historical Go 1.27.0 evidence below is
preserved and was not rerun under the new toolchain. Full canonical verification is pending the
independent root gate; focused binary, actor, wire-hash, and fixture-path evidence accompanies the
candidate commit. Runtime, production queue, Linux-native, crash/power-loss, and performance
qualification remain NOT_RUN or NOT_OBSERVED.

## 2026-09-07 — Tickets can be created: one MUTATE operation, and the redo rule finished

TCP-02b is delivered. All fourteen `ticket` mutation verbs commit through the §5.2 writer, so a
queue can hold real tickets for the first time.

The transaction model had five closed operations and none of them could create a ticket.
`ADOPT_FILE` is the near miss and is explicitly not it: §3.3 requires the offered file to carry a
`ticketId` equal to an existing canonical record's, and forbids it from completing, archiving,
reopening, restoring, approving or re-sourcing. A ticket mutation had no journal representation.

Decision 0004 adds exactly one operation, `MUTATE`, carrying one `taskman-mutation/0` envelope,
rather than one per §3.3 verb. `mutation.Apply` already validates and computes the post record for
all fourteen, and the staged artifact shape is identical for every one — ticket projection,
request entry, receipt, head, plus `intent/queue.json` when a CREATE allocated a serial. Fourteen
operations would have multiplied the closed preimage, template and slot tables by fourteen and
added no check the envelope codec does not already make. The receipt keeps the §3.1 vocabulary
that already anticipated this: `ARCHIVE` and `RESTORE` receipts carry those kinds, everything else
is `MUTATION`.

**Three defects the genesis transaction could not reach.** TCP-02 wrote only to absent
destinations, so each of these was unreachable until a second transaction ran, and each one
wedged the store when it fired.

The redo rule was missing its middle case. §5.2 admits a destination at the post digest "or the
`pre` digest of a receipt that is still redo-pending"; the delivered rule handled the post digest
and absence and refused everything else, so an ordinary CREATE — whose `queue.json` legitimately
holds the pre state — was refused `INTENT_DIVERGED`. The pre digest now comes from the committed
receipt's own `pre` entries, the same evidence a crash recovery reads.

Post-commit redo did not exist. The receipt links in before any post file, so a post-phase failure
leaves a committed receipt, unwritten projections and a stale head (crash point C2). That is not a
hypothetical: the first defect produced exactly this state, and the store then refused every
subsequent transaction because the receipt count disagreed with the head. `internal/store/redo.go`
completes a pending receipt before anything new is modelled.

Intent projections could not be replaced at all. The reviewed applier admitted replacement for the
head, barrier and reservations only, and its `expected == nil` path degrades to an exclusive
create — so `head.json` could be written exactly once, and no ticket or queue projection could
ever be rewritten. The applier now admits the intent roles too, and admits them *only* under an
explicit expected digest: a projection is replaced by a real compare-and-swap on the bytes the
transaction read, which is stricter than what the three original roles are held to. That is a
13-line change to the 893 reviewed lines.

**The CLI shape is a decision, not a spec reading.** §3.3 fixes the envelope wire format and says
nothing about flags. The verb names the operation, `--payload` carries the closed payload as
canonical JSON, `--request-id` is the idempotency key. Because the TM-V0-006 digest is the digest
of the envelope bytes and `issuedAt` is one of them, a retry replays only by reproducing the same
envelope — so `--issued-at` exists, and without it every invocation is a new request. A CREATE
that allocates a serial reports the id it allocated; that is the only place a caller can learn it.

Verified: `make verify` exits 0, including the 145 s authority suite over the edited applier and
the reviewed transaction and snapshot suites. Against a real repository, `ticket list` goes from
zero rows to the created ticket, an identical retry replays without a second receipt, a stale
`expectedRevision` refuses with the store byte-identical, hold/release/archive/restore each carry
their own receipt kind, staging is left empty, and `archive export | archive verify` round-trips.

Held: attempt liveness is still the zero-attempt oracle, which is sound only while no runtime can
start an attempt. `reconcile intent`, the administrative verbs, `archive restore`,
staging-descriptor crash recovery, reservations and the runtime remain unbuilt. Every Coverage
axis stays `NOT_OBSERVED`; a non-fixture queue still needs the §7.4 cutover record and a
`QUALIFICATION` receipt. GP is NOT_RUN.

## 2026-09-07 — The store works: callable writer, local-operator issuer, `atm init`

`atm init` was not missing wiring. Three locks held the durable path shut, all deliberate:
`transaction.Model` refused every premise except `HYPOTHETICAL_FIXTURE_NO_RUNTIME`;
`mutation.Binding` documents that deciding who may furnish a binding is "the runtime enforcement
profile's job", and no such profile exists; and decision 0002 held callable real-queue writers,
having already rejected same-UID, held-lock and supplied-role as grants. SPEC §5.7 named the
missing slice: J3-05 real issuer/authorization, existing-intent CAS/C6, callable writer
sequencing/redo.

The owner answered both open questions on 2026-09-07 (decision 0003): the local operator is the
authority for their own local store, with the binding recorded unauthenticated; and a second
premise `OBSERVED_LOCAL_OPERATOR` is admitted beside the fixture one rather than reworking the
gate into per-axis coverage. Coverage stays `NOT_OBSERVED` on every axis — admitting the premise
grants nothing, and `TestTMV0001_AS29_InitRecordsAnUnauthenticatedBinding` fails if that slips.

What was already built did most of the work. `internal/authority/fixture_session.go` is the real
applier — its roles and directories are the actual §3.4 layout, and only its name and visibility
were fixture-scoped — so it was wrapped by an exported `Session` in a new file rather than
rewriting 893 reviewed lines. `internal/store` is the new §5.2 driver: evidence, then the receipt
link-in as the commit point, then post files, then the head.

Three defects surfaced only by running it against a real repository, none visible to the fixture
tests. Publishing blindly hit `file exists` on `intent/queue.json`, because the operator's own
authored file is already at the post state: the §5.2 redo rule was missing, and a destination at
a third value must be refused `INTENT_DIVERGED`/`JOURNAL_FORKED` rather than overwritten. The
request index needed its `requests/<xx>` shard created on demand. And leaving staging slots
occupied broke `archive export` with `a00: unassigned stage slot` — every slot a transaction uses
is now cleared, including on a failure part-way through.

Evidence: `make verify` exits 0 (Go 1.27.0, gofmt, `go test ./...`, `go vet ./...`). On a fresh
repository `queue status`, `ticket list` and `gate list` go from `REFUSED`/`UNINITIALIZED` to
`OK`; `archive export | archive verify` round-trips; a second `init` refuses `BLOCKED` and writes
nothing. Held: every mutation verb other than `init`, staging-descriptor crash recovery,
hostile-editor CAS, reservations, the runtime, and any non-fixture queue, which still needs the
§7.4 cutover record and a `QUALIFICATION` receipt. GP is NOT_RUN.

Current status: experimental support through J3-02 passed parent gates and fresh reviews.
J3-03/04 fixture integration is authored; its combined parent gate and fresh review remain
PENDING. TCP-02 and callable writer/runtime/real-queue/restore/GP qualification remain held.

## 2026-09-06 — TCP-02 J3-02 bounded read-only staging candidate

Implemented the accepted finite J3-02 plan and both selected Gate A corrections, including
empty-D request lookup and the later selected completed operation/timestamp binding correction.
Both completed-linkage and receipt-only genesis corrections had causal-red/focused-green
witnesses. Headless binding follows the full existing audit and uses its canonical queue only
on actual physical absence. Owned sources are snapshot/stage_observation.go and tests;
archive/layout.go, export.go, stage_read.go and tests; journal/audit.go, records.go, source.go,
stage_read.go and directly relevant tests. No Source.List adapter needed a signature repair:
existing decorators embed Source and inherit its concrete Listing return. The relevant old
scan-budget test now excludes the separately observed root identities from D.

Focused regressions prove closed staging structure, consumed-byte hashes, native root/parent
lifetime and replacement, scan/payload separation, actual completed-head/request binding,
four attempts, body-error rechecks, sticky resource errors and read purity. Empty-D causal red
failed MALFORMED before capture; focused green permits verified request FOUND/ABSENT while
strict Audit remains INTENT_DIVERGED. See the builder report for commands, failures and limits.

Protected J2b/J3-01/safeopen/codec/Reader sources remain byte-identical to the supplied baseline.
No delegation, staging/commit/merge, canonical make verify or broad project gate was run.
Parent source freeze, canonical gate and fresh independent review remain pending. TCP-02
is incomplete; this is structural read observation and archive chain/digest evidence only.


## 2026-09-06 — TCP-02 J3-01 private fixture filesystem primitives

Base: accepted J2a `b67ed0d` (owner-supplied baseline). Read the full narrowed Gate A,
AGENTS instructions, SPEC §1..§7, and the larger J3 plan for retained assumptions/holds.
Gate A cleared this private primitive slice only; no transaction runner was built. Parent
steering confirms the owner's fixture-only selection; parent owns decision 0002.

Owned changes: eight new `internal/authority/fixture_session*.go` implementation/test files;
`docs/SPEC.md`, `docs/ROADMAP.md`, this log and the J3-01 builder report. Existing Go files
compare byte-for-byte with `/private/tmp/corvint-tasks-build-20260906/tcp02-j2a-final`.
No safeopen dependency or Lock lifecycle edit was needed. No Git operation, delegation,
canonical gate, server, database, runtime child or transaction-phase runner was invoked.

Verification by step:

1. Identity/lifetime: descriptor-pinned closed parents, actual Lock lifecycle mutex and
   fixed session-before-Lock order; wrong/closed/replaced identities refuse. Tests pause
   before native link and during destination durability with concurrent Close requests;
   exclusion stays held, handles close, subsequent acquisition succeeds, close errors join.
2. Primitives: deterministic tables cover exact staging roles/bounds, full/short/partial/empty
   writes, exclusive link, private rename, controlled unlink, admitted mkdir, strict file
   sync, every required directory sync and close faults. Fresh reads verify visibility,
   source retention and byte/mode/mtime preservation of foreign entries. Failed preparation
   never returns a token. Final self-check also added invalid-pre validation and file-name/absence
   rechecks at post-effect and already-applied directory-sync boundaries. Request bounds were checked against accepted J1 and fixed at 64 KiB.
3. Platform/compatibility: Darwin/arm64 native fixture tests execute; mount mismatch/unknown
   and late EXDEV refuse with no copy fallback. Existing LinkIn nonempty, empty, exclusive
   and injected-error tests run unchanged. Linux/amd64 and Windows/amd64 test compilation
   succeeds; native execution there is NOT_RUN.
4. Final focused package check: Go 1.27.0, `go test -race -count=1 -timeout=60s -json
   ./internal/authority` succeeds in 4.339 s. JSON records 40 top-level tests and 207 subtests
   passing, including 17 new fixture tests and 190 fixture subtests; zero skipped or failed
   tests. `gofmt -l internal/authority/fixture_session*.go` returns no files.

All Go commands use `GOTOOLCHAIN=local`,
`GOCACHE=/private/tmp/corvint-tasks-build-20260906/go-cache`, `GOMAXPROCS=2`, `GOFLAGS=-p=2`,
without a PTY. Exact commands, finite witness mapping and artifact paths are in
[the builder report](reviews/2026-09-06-j3-01-builder.md).

Parent combined canonical gate and fresh exact-diff review remain pending. TCP-02 is incomplete.
J3-02..05, staging-layout acceptance, J2b capacity/reserves integration, real issuer, existing-intent
CAS/C6, sequencer/redo, runtime and real queues remain held. These tests do not establish true
power-loss durability or accepted historical transitions. GP is NOT_RUN; no Corvint source edit
or performance claim is made.

## 2026-09-06 — accepted fixture foundation

Accepted experimental J2b pure transaction/capacity/recovery model and J3-01 private fixture filesystem primitives, on accepted J1/J2a b67ed0d. Owner decision 0002 keeps this delivery fixture-only.

The original J2b review found an unreachable fresh UNPAUSE after abort preserved divergent ticket bytes, and incorrect cross-operation request-conflict ordering. A fresh builder reproduced and repaired both; a fresh reviewer passed the exact repair. J3-01 passed its separate fresh review with no HIGH/MED findings. Reports, initial findings, Gate A selections and steering are retained beside this report.

One combined canonical `make -j1 verify` passed Go 1.27.0, formatting, all tests and vet. All 178 frozen files (93 Go) matched their hashes afterward. Both reviews used the same immutable combined snapshot; all component Go source was preserved. Integration changes only status prose and adds evidence. The primitive builder also passed authority race tests: 40 tests and 207 subtests, no skips/failures. Linux/Windows test binaries cross-compile; native execution there is NOT_RUN.

J2b carries five closed hypothetical templates, strict shared staging codec, exact archive encoding costs, in-cap UNPAUSE reservation, modeled abort/redo and retained orphan costs. J3-01 supplies unexported primitives called only by tests owning disposable repositories under sole-writer/stable-topology assumptions. Neither is a callable transaction writer or an actor/administrative authority.

J3-02 staging read integration and J3-03/04 test-only transaction/recovery integration are next. J3-05 issuer, hostile-editor intent CAS, mutation CLI, runtime, real queues, restore, power-loss and GP qualification remain held. TCP-02 is incomplete. Corvint Go code is unchanged; no main/default-branch merge or publication is included.

## 2026-09-06 — J3-02 parent acceptance

One canonical gate passed on all 197 frozen source/document files (99 Go), and the fresh read-only review passed with no HIGH/MED finding. The exact reviewed native-read increment is saved experimentally; [evidence and limits](reviews/2026-09-06-j3-02-integration.md). Fixture integration, durable writer/runtime and GP remain pending or held.

## 2026-09-06 — J3-03/04 combined fixture integration candidate

Rebased only the three frozen foundation test files onto reviewed J3-02 bc3b3dd, adding
`internal/authority/fixture_read_integration_test.go`. Shared structural observation and actual
queue/head/receipt/request binding now participate in observation/reopen, while retained
no-follow identity checks and cleanup-only partial tokens remain mandatory. Two parent-selected
harness repairs retain copied full model-time inventory and replace the vacuous branch witness
with an actual ticket token at a committed prefix plus a valid-branch control.

Native strict Audit and archive Export/Verify cover the five operation flows with persistent
empty staging and completed leftovers; precommit export preserves actual D and linked orphans,
and source syntax/intent refusals remain visible. Receipt-only INIT obtains its absent queue
binding from validated genesis; archive remains UNINITIALIZED. Empty-D FOUND/ABSENT lookup
reports NOT_OBSERVED intent agreement and refuses private-index corruption. Combined histories
measure staging D versus archive F and fresh sequential unchanged-head reads. The final J3-02
private-hook matrices retain exhaustive in-read movement/retry/sticky-error ownership.

Exact focused commands/results, causal-red overlays, protected-file hashes and the five-group
TM/AS mapping are in [the builder report](reviews/2026-09-06-j3-fixture-integration-builder.md)
and SPEC §5.7. Earlier builder-time pending records above remain historical; J3-02 acceptance
is recorded separately. The J3-03/04 combined parent gate and fresh review of all fixture changes
are PENDING. No Git mutation, delegation, canonical full gate, crash child or runtime was used.
Returned faults are not process death; interruption/power-loss evidence is NOT_PRODUCED.
TCP-02, callable writer/issuer/admin binding, hostile-editor CAS, CLI mutation, runtime, real
queues, restore and GP remain incomplete or held.

## 2026-09-06 — Fixture-only first delivery accepted experimentally

One canonical Go 1.27.0 formatting/all-tests/vet gate passed on the final 211-file candidate (103 Go). All hashes matched after the gate and fresh independent final review. The reviewer returned PASS with no HIGH/MED finding and accepted every fixture criterion and all five native integration groups. The original 99 Go files are unchanged from reviewed J3-02. [Delivery evidence](reviews/2026-09-06-fixture-delivery.md) records exact finite coverage and remaining holds. Parent status/evidence edits supersede earlier candidate-pending entries. Crash-matrix C1–C3 wording now cites the existing §5.2 receipt-link and exact PRE/POST rules; no behavior changed. TCP-02, real writer/runtime/restore and GP remain incomplete; no default-branch merge or publication.
