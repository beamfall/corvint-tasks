# Corvint Tasks contract (TCP-00)

**Owner:** Russell Lewis  
**Date:** 2026-09-06  
**Intent status:** accepted for this repository (decision 0001; owner instruction "use claude to build this")  
**Delivery status:** experimental TCP-01 and TCP-02 support through J1, J2a/J2b and
J3-01 passed parent canonical gates and fresh independent reviews. J2b's two review
findings were repaired and passed fresh review. See
[fixture foundation evidence](reviews/2026-09-06-fixture-foundation-integration.md).
Reviewed J3-02 native read support also passed its parent gate and fresh review
([evidence](reviews/2026-09-06-j3-02-integration.md)). Decision 0002 selects fixture-only
first delivery. J3-03/04 passed its combined gate and fresh review
([delivery evidence](reviews/2026-09-06-fixture-delivery.md)). Decisions 0003/0004 delivered
callable fixture initialization and all fourteen ticket mutations. The 2026-09-19 writer
repair is tracked in [repair evidence](reviews/2026-09-19-writer-repair.md). Reservations,
process runtime, real-queue admission and performance qualification remain unbuilt or unqualified.
**Binary:** `corvint-tasks`, native Go 1.27.1, standard library only
**Wire profiles frozen here:** `taskman-queue/0`, `taskman-policy/0`, `taskman-ticket/0`,
`taskman-mutation/0`, `taskman-outcome/0`, `taskman-receipt/0`, `taskman-journal-head/0`,
`taskman-init/0`, `taskman-attempt/0`, `taskman-reservation-set/0`, `taskman-effect/0`, `taskman-plan/0`,
`taskman-gate-result/0`, `taskman-claim-disposition/0`, `taskman-completion-manifest/0`,
`taskman-capability-profile/0`, `taskman-import-plan/0`, `taskman-import-map/0`,
`taskman-archive/0`, `taskman-barrier/0`, `taskman-command-result/0`, `taskman-perf-baseline/0`
**Repair round:** 2026-09-06 R1 (independent Gate A review, verdict REPAIR; findings F1..F6 and
dropped-obligation list repaired below; gates remain NOT_RUN); 2026-09-06 R2 (fresh review,
verdict REPAIR; F1 pending-effect ownership, F2 bootstrap windows, F3 `ADOPT_FILE` composition,
F4 receipt overwrite/chain bounds/`archive export` snapshot, LOW cutover A5 exemption; repaired
in §1, §3.1, §3.3, §3.4, §5.2, §5.3, §5.4, §6.2, §6.4, §9.2, §11; gates remain NOT_RUN);
2026-09-06 R3 (Gate A review of R2, no HIGH findings; one MED documentation repair applied with
TCP-01: §5.2 temp cleanup limited to `receipts/*.tmp-*` and `head.json.tmp`, §6.4 non-`EEXIST`
link failures exit 3 without exec, AGENTS invariant 3 aligned to the link-in commit; and the
TCP-01 digest preimages frozen in §3.1 as an implementation detail; gates remain NOT_RUN)

## Agent digest
- Claim: Freezes the on-disk schemas, numeric limits, lock/commit/recovery protocol, transition
  table, acceptance tables and adoption protocol that TCP-01..09 implement.
- Status: accepted intent; fixture support through J3-03/04 passed parent gates and fresh reviews; callable init and ticket mutations delivered under decisions 0003/0004. TCP-02 remains incomplete.
- Exists: this contract, `docs/ROADMAP.md`, decisions 0001..0004, experimental TCP-01/TCP-02 Go source/tests and callable fixture writer. The frozen source-copy paths were missing at the 2026-09-19 audit; their recovery remains required.
- Blocked on: TCP-01 and TCP-02 implementation; Corvint-side amendment record (§10 U1).
- Read next: §2 identity and encoding; §3 records; §5 storage protocol; §6 transitions; §7
  acceptance tables; §9 slices and ownership.

## 0. Authority, boundary and sources

Authority descends: owner instruction and owner-recorded decisions (decision 0001) > the
accepted, governing invariants of the Corvint Task Control Plane V0 (`ATCP-V0-001..022`), of
ATM revision 6 (`ATM-V0-*`) and of WQO V0 (`WQO-V0-*`, unchanged) > this contract, which is a
generated local interpretation that freezes mechanics (schemas, limits, lock and commit
protocol, transitions) within that boundary and whose §8 rows hold only where the owner's
decision accepts them > code and tests. A generated document never silently accepts an
upstream amendment and never qualifies a real cutover; only the records of §7.4 and the owner
do. Frozen copies of those sources are under
`/private/tmp/corvint-tasks-build-20260906/context`. Provenance: the ATCP document
(`corvint-task-control-plane-v0.md`) and the four-expert review were **uncommitted working-tree
files** alongside Corvint base `0a227a94df6ad02e35c01965d05c46be9b367905`, not members of that
commit; ATM revision 6, WQO V0, `AGENTS.md` and `SPEC-DRIVEN-DEVELOPMENT.md` are the versions
present at that base. The reference for each frozen file is its path under the context directory
plus its SHA-256; the digests are recorded in decision 0001 §"Provenance" (computed by the
coordinator in R2; never inferred, and never a claim of commit membership). This repository never modifies Corvint or Beamfall
source.

Interpretation rule: clauses in this contract that go beyond the accepted ATCP/ATM text (for
example the canonical-record rule of §3.1, the barrier scopes of §5.2, the spawn bootstrap of
§6.4 and the performance gate of TM-V0-027) are implementation interpretations made under the
owner's build instruction. They may add safeguards; they never lower an accepted ATCP/ATM/WQO
safeguard, and model-authored prose here creates no authority over the owner's sources. Where the
two conflict, §8 must carry an explicit row or the stricter accepted rule stands.

Boundary: Corvint is the task-management product and evidence/selection engine; `corvint-tasks` in this
repository is the operative executor and the owner of the local ticket store, journal,
reservations and receipts. WQO V0 stays non-operative and closed; the operative planner is the
separate profile `taskman-plan/0` (§4.3). No `corvint` mutation facade exists.

This contract authorizes local implementation only. It does not authorize publication, real
queue cutover, autonomous dispatch of real repository tasks, foreign-queue writes, or deletion of
history. Those need the records named in §7.4 and §9.

## 1. Frozen numeric limits

Every limit is validated before any effect. A value over its limit fails closed with
`LIMIT_EXCEEDED` and produces no partial write. Policy files that exceed a `max` column or set a
value below a `min` column are rejected at load.

| Area | Limit | Value |
|---|---|---|
| Identifier / label / local token / requestId | bytes | 1..128 / 1..64 / 1..64 / 1..64 |
| PathText (absolute filesystem path: `head.primaryWorktree`, archive `primaryWorktree`) | bytes | 1..4096 |
| Title / body / one acceptance criterion / one prose field | bytes | 1..512 / ≤64 KiB / ≤4 KiB / ≤64 KiB |
| Ticket record file / mutation envelope / outcome | bytes | ≤128 KiB / ≤256 KiB / ≤64 KiB |
| Per ticket: acceptanceCriteria, dependencies, requirementRefs, resources, approvals | count | ≤64 each |
| Per ticket: touchPaths / labels / holds / requiredGates / capabilities | count | ≤256 / ≤32 / ≤16 / ≤32 / ≤32 |
| Tickets per queue / intent tree bytes | | ≤10,000 / ≤256 MiB |
| queue.json / policy.json / import-map.json (≤10,000 entries) | bytes | ≤1 MiB / ≤256 KiB / ≤8 MiB |
| Receipt file (inline post entries only) / inline post entry / attempt record / journal head | bytes | ≤1 MiB / ≤64 KiB each, ≤8 per receipt / ≤64 KiB / ≤4 KiB |
| Reservation set: resources per entry / entries / file | | ≤4,096 / ≤`maxActiveAttempts` max (64) / ≤64 MiB |
| Import apply records per `IMPORT_APPLY` receipt | count | ≤100 |
| Barrier record / command result envelope (excluding paginated `items`) | bytes | ≤4 KiB / ≤64 KiB |
| Journal saturation (admission refused above) / warning | | 1,000,000 receipts or 4 GiB / 50% |
| Evidence blob / evidence store total | bytes | ≤64 MiB / ≤16 GiB |
| Runtime input packet / captured lane transcript+logs per attempt / gate stdout+stderr | bytes | ≤4 MiB / ≤256 MiB / ≤64 MiB |
| Pagination page default / max; list result | | 100 / 1,000; ≤16 MiB |
| Import source items / import plan / archive total | | ≤10,000 / ≤16 MiB / ≤64 GiB |
| Archive `files` entries / `manifest.json` bytes / state-dir entries scanned by `archive export` (B4 freeze, §3.5; independent hard aggregate budgets, unverified) | | ≤2,100,000 / ≤768 MiB / ≤2,104,096 (files + 4,096) |
| Directory listing chunk (names materialized per read while enumerating a store directory) | count | 256 |
| Lock wait / heartbeat interval / drain grace / SIGTERM→SIGKILL | seconds | 30 / 30 / 60 / 10 |
| Bootstrap waits (§6.4): leader `.ack` wait / supervisor `.boot` wait / wrapper exit wait after a fence | seconds | 90 (lock wait + drain grace) / 90 / 10 |
| Consecutive `FAILED/SPAWN` per attempt with no runtime executed (`attempt.spawnNoExecCount`, §6.4), not charged to admissions per revision; exceeding it refuses `RETRY_EXHAUSTED` | count | 3 |
| Lane wall clock default / max; gate timeout default / max | minutes | 90 / 240; 30 / 120 |
| Liveness re-check threshold (tool-owned lanes) | | 2 × largest configured wall-clock cap |
| Per-lane token caps (ATM-V0-006 defaults) | | input 2,000,000; cacheCreation 1,000,000; cacheRead 100,000,000; output 400,000; turns 400 |
| Aggregate per ticket / per wave | | 4 × per-lane caps / per-ticket × `maxActiveAttempts` |
| `maxActiveAttempts` default / max; `maxWorkersTotal` default / max | | 1 / 64; 4 / 256 |
| Admissions per ticket revision (retries) / repair rounds / malformed-review retry / automatic gate re-run after STALE / adapter reconcile attempts | count | 3 / 2 / 1 / 1 / 3 |
| Evidence retention after terminal attempt before prune-eligible | days | default 90, min 30 |
| Decoded JSON depth / array elements / aggregate nodes | | 24 / 10,000 / 250,000 (unchanged; the archive manifest alone parses its top-level `files` array under ≤2,100,000 elements and ≤8,400,014 nodes through the opt-in entry point of §3.5; depth and every other array keep these bounds) |

Unknown usage never creates headroom: a budget field with any `NOT_OBSERVED` contribution has
aggregate `NOT_OBSERVED`, its cap is reported `NOT_ENFORCED`, and it is never treated as zero.
`turns` and wall clock are always observed and always enforced.

Post-file sizing rule: a post entry whose canonical bytes exceed 64 KiB, or any post entry beyond
the eighth, is written to `evidence/<sha256>` (fsync) before the receipt and referenced by
`blobSha256`; redo reads the blob. The reservation set, `import-map.json` and every import batch
therefore never breach the receipt cap, and the caps above are mutually reachable: the largest
inline receipt is 8 × 64 KiB plus envelope, under 1 MiB. Any transaction whose post set cannot be
represented this way fails `LIMIT_EXCEEDED` before writing.

## 2. Identity and canonical encoding

Encoding is WQO §4.1 verbatim: closed objects, sorted UTF-8 keys, no insignificant whitespace,
exactly one trailing LF on disk and transport, `Count` as decimal string ≤ `2147483647`,
`Digest` as 64 lowercase hex, `Path` as WQO `Path`, hostile code points rejected, no Unicode
normalization. `Size` is a second numeric primitive: a decimal string without leading zeros in
`0..18446744073709551615` (uint64). `Size` is used for every byte size, token budget field and
its aggregate, `seq`/`lastSeq`/`headSeq`/`*Seq`, `generation`, `pid`, `startTime`, `pgid` and
`policyVersion`; `Count` is used only for cardinalities and small counters (retries, rounds, units,
`order`, `nextSerial`, `criterionIndex`, `revision`, `acceptanceRevision`, `exitCode`, minutes,
days, seconds). Arithmetic on `Size` that would exceed uint64 fails `LIMIT_EXCEEDED` before any
write; a `Size` where a `Count` is declared, or the reverse, is a validation failure. Tests cover
`0`, max, max+1, leading zero and sign for both primitives (AS-10). `PathText` is a third text
primitive (B2 resolution, 2026-09-06, reversible clarification under the TCP-01 freeze): an
absolute filesystem path of 1..4096 UTF-8 bytes, leading `/`, hostile code points and
TAB/LF/CR rejected, never normalized; it types `head.primaryWorktree` (§3.4) and the archive
manifest's `primaryWorktree` (§3.5) and nothing else. It shares no bound with `Identifier`; the
repository, queue and ticket identity limits are unchanged. Witnesses:
`TestTMV0002_AS10_PathTextBoundaries`, `TestTMV0002_AS01_ReceiptAndHeadSchemas`. Semantic arrays (listed per record) keep declared order; all other arrays are
canonical-byte sorted without duplicates. `null` is legal only where a field says `|null`.
Timestamps are `YYYY-MM-DDTHH:MM:SSZ` (UTC, seconds) and are advisory: no identity, expiry,
reclaim or ordering decision uses wall time.

Content identity is WQO §4.3: `SHA-256(kind || 0x00 || profile || 0x00 || canonicalBody)`.
Record identities are `<kind>:sha256:<Digest>`; a file's chain digest is SHA-256 of its raw bytes.

| ID | Grammar |
|---|---|
| repository / queue | WQO `repo:<a>` and `queue:<a>:<q>`; tokens `[A-Za-z0-9][A-Za-z0-9._-]*` |
| ticket | `ticket:<a>:<q>:<local>`; `<local>` 1..64 bytes of the token grammar, unique under ASCII case folding, never reused (tombstones keep it) |
| attempt | `attempt:<a>:<q>:<32 lowercase hex from crypto/rand>` |
| gate / capability / capacity / resource keys | `label` (1..64 bytes) unique within policy |
| receipt | file name `receipts/<seq as 12 decimal digits>.json`; `seq` starts at `1` |
| profile version | the `/0` suffix; any other version fails `UNSUPPORTED_VERSION` before use |

Native `CREATE` without an explicit local token assigns `<queue.prefix>-<nextSerial as 4+
digits>` and increments `nextSerial` in the same transaction; explicit tokens are allowed
(import mapping such as `AT-07`) and must not collide under case folding with any live or
tombstoned ticket.

## 3. Records

### 3.1 Intent store (Git-tracked, in the repository tree)

```
.taskman/queue.json          taskman-queue/0     queue manifest and authority policy
.taskman/policy.json         taskman-policy/0    gates, limits, roles, capacity, runtimes
.taskman/import-map.json     taskman-import-map/0 import identity map (absent until first import)
.taskman/tickets/<local>.json taskman-ticket/0   one current record per ticket, tombstones included
```

**Canonical record and projection.** For every intent file the canonical current record is the
latest journal receipt's `post` entry for that path (§3.4); the Git-tracked file is a projection
of it. The journal decides `revision`, `expectedRevision` checks and eligibility; the file never
does. A projection that disagrees with the journal is a divergence (C6), never an alternative
truth: `corvint-tasks` neither adopts the file silently nor overwrites a human edit silently. Divergence
freezes that ticket until `reconcile intent` records the operator's explicit choice as its own
receipt: `KEEP_JOURNAL` (rewrite the projection from the canonical record, `RECONCILE` receipt
with the discarded file bytes stored in `evidence/`) or `ADOPT_FILE` (a composition of
`OPERATOR`-permitted mutations built from the file's differing fields by the fixed rule in §3.3
"ADOPT_FILE composition", so the adoption itself becomes a journal revision; any difference the
rule does not cover refuses and leaves the ticket diverged). Every intent-file write happens only inside a transaction that names the file as a
post entry and passes the branch and divergence checks of TM-V0-007, whatever the receipt kind
(`MUTATION`, `MANIFEST` reducer, `SYSTEM` hold, `IMPORT_APPLY`, `AUTHORITY_SWITCH`, `RESTORE`,
`RECONCILE`).

**Primary worktree.** The primary worktree is the worktree whose `.git` is the common directory
itself (`git rev-parse --git-common-dir` equals `--git-dir`); linked worktrees are never the
primary. `corvint-tasks init` records its absolute path bytes as `primaryWorktree` in `head.json` and in
the `INIT` receipt; every later command resolves the common dir (§3.4), compares the recorded
path to the primary worktree found at that common dir, and refuses `UNSUPPORTED` if they differ
(a moved repository needs the planned, unimplemented-in-this-preview
`corvint-tasks init --relocate`, an `OWNER` receipt). All intent-file reads and
writes, branch checks, `intent tree digest` values and publication status use the primary
worktree only, regardless of the caller's cwd.

**Digest preimages (R3; TCP-01 implementation detail frozen under the owner's build
instruction, not an accepted ATCP/ATM clause; stable IDs unchanged).** `intentTreeSha256` is
SHA-256 over the concatenation, for every file of the intent store in byte-sorted
store-relative path order, of `path || 0x00 || lowercase-hex SHA-256 of the file's raw bytes ||
0x0A`; a store with no files digests the empty string. The store layout is flat and closed:
the only admissible entries are the regular files `queue.json`, `policy.json`, `import-map.json`
and the directory `tickets/` holding only regular files `<local>.json` whose `<local>` is a
valid local token. Any other entry, a nested directory, a symlink, an empty file, more than
10,000 ticket files, a file over its §1 bound or a tree over 256 MiB fails the read before any
byte is hashed. The digest and every decoded record are computed from one captured copy of the
bytes, and a reader that pinned a tree digest refuses `SNAPSHOT_MOVED` when the captured tree
differs from it. `primaryWorktreeSha256` is SHA-256 of the exact absolute path bytes recorded
as `head.primaryWorktree` (no normalization, no trailing separator). Witnesses:
`TestTMV0008_AS08_IntentTreeDigestPreimage`, `TestTMV0008_AS08_PrimaryWorktreeDigestPreimage`.

**Bounded enumeration.** Every store directory is listed through the root-confined descriptor
in chunks of 256 names (§1) and refused `LIMIT_EXCEEDED` as soon as its entry bound is passed:
the store root admits at most its four named entries and `tickets/` at most 10,000 entries.
Every entry counts, whatever its name, so an empty-file or stray-name flood is refused during
the listing, before any name is validated and before any entry is stat'ed, opened or read; the
listing never materializes more than the bound plus one chunk of names. Witnesses:
`TestTMV0002_AS10_ReadDirNamesBoundedEnumeration`, `TestTMV0002_AS10_TreeDigestTicketCountBound`,
`TestTMV0002_AS10_TreeDigestFlatLayout`.

`taskman-queue/0`: `profile, queueId:queue:*, repositoryAuthorityId:repo:*, prefix:label,
nextSerial:Count, schemaVersion:"0", canonicalWriter:"NATIVE"|"ROADMAP"|"FOREIGN",
foreignAdapterId:label|null, intentBranch:label, fixture:boolean,
executionCutover:{enabledBy:label, decisionRef:label, gateEvidence:[Digest]}|null,
importMapSha256:Digest|null, writeBarrier:{reason:"CUTOVER"|"EMERGENCY"|"NONE",
since:timestamp|null}`. `importMapSha256` is the chain digest of `import-map.json` (null when
absent) so that `queue.json` stays under 1 MiB and ordinary `CREATE` never rewrites the map.

`taskman-import-map/0` (closed): `profile, queueId, entries:[{sourceQueueId:Identifier,
sourceItemId:Identifier, ticketId:ticket:*, sourceRevisionSha256:Digest|null,
appliedSeq:Size}]`, sorted, no duplicate `(sourceQueueId, sourceItemId)`.

`taskman-policy/0`: `profile, policyVersion:Count, roles:{<role>:[operation]}, capacity:
{maxActiveAttempts:Count, maxWorkersTotal:Count, classes:[{id:label, units:Count}]}, budgets:
{lane:{inputTokens,cacheCreationTokens,cacheReadTokens,outputTokens:Size, turns:Count,
wallClockMinutes:Count}, ticketMultiplier:Count, requireEnforcedFields:[label]},
retries:{admissionsPerRevision, repairRounds, malformedReviewRetry, gateRerunOnStale,
reconcileAttempts:Count}, retention:{evidenceDays:Count}, gates:[GateDefinition],
serialFallback:"BLOCK"|"WHOLE_REPOSITORY", integrationRequiredKinds:[kind],
allowEmptyObligationsKinds:[kind], reviewLane:{required:boolean}, docsLane:{required:boolean},
cemRequired:boolean, ocmRequired:boolean, runtimes:[RuntimeEntry], environment:
{allowedEnvKeys:[Identifier]}`. The four token budget fields are `Size` per §2 ("every token
budget field"); `turns` and `wallClockMinutes` are `Count` (B1 resolution, 2026-09-06; witness
`TestTMV0002_AS10_LaneBudgetPrimitives`). Policy identity `policy:sha256:*` is pinned per attempt.
Workers cannot write policy: policy changes are `OWNER`/`OPERATOR` mutations and are themselves
tasks (ATM-V0-027).

`RuntimeEntry` (closed): `runtimeId:label, executable:{pathSha256:Digest, fileSha256:Digest,
mode:Identifier} (WQO ExecutableIdentity), argvPrefix:[Identifier] (literal, ≤16),
capabilityProfileSha256:Digest (a pinned taskman-capability-profile/0), observedBudgetFields:
[label] (subset of the lane budget field names), roles:["BUILDER"|"REVIEWER"|"VERIFIER"|
"REPAIR"|"DOCS"], maxWorkers:Count, enabled:boolean`. A runtime whose profile digest is not
present in `pinned/` is unavailable (`CAPABILITY_UNAVAILABLE`), never probed lazily.

`taskman-capability-profile/0` (closed): `profile, runtimeId, probedAt, probeReceiptSeq:Size,
executableFileSha256:Digest, axes:{read:Axis, write:Axis, network:Axis, process:Axis,
tools:Axis}, readOnlyMode:{supported:boolean, enforcementAxis:"ENFORCED"|"OBSERVED"|"UNKNOWN",
argv:[Identifier]}, budgetObservability:{inputTokens,cacheCreationTokens,cacheReadTokens,
outputTokens,turns:"OBSERVED"|"NOT_OBSERVED"}, evidence:[Digest]` where `Axis = {enforcement:
"ENFORCED"|"OBSERVED"|"UNKNOWN", mechanism:"OS_PERMISSION"|"TOOL_PERMISSION"|"RUNTIME_FLAG"|
"NONE", scope:prose, probe:Digest|null}`. `ENFORCED` requires a recorded probe in which the
runtime attempted the forbidden action and the OS or the tool refused it; `OBSERVED` means only
that a violation would be detected afterwards; `UNKNOWN` means neither. Process-group membership
is never a mechanism for any axis.

`GateDefinition`: `gateId:label, kind:"COMMAND"|"REVIEW"|"DOCS"|"MANUAL"|"EXTERNAL",
argv:[Identifier] (literal, 1..16, never ticket text), cwd:"WORKTREE"|"CANDIDATE",
env:[Identifier] (names drawn from allowedEnvKeys), timeoutSeconds:Count,
expected:{exitCode:Count|null, reducer:label|null}, evidence:[label],
inputs:[Path], sharedResource:{class,key}|null, reusable:boolean, required:boolean`.
`WORKTREE` runs in the lane worktree and requires, immediately before spawn, an empty porcelain
status and `HEAD^{tree}` equal to the attempt's `candidateTreeOid`, else the result is `BLOCKED/
DIRTY_WORKTREE`. `CANDIDATE` runs in a fresh detached worktree checked out at exactly
`candidateTreeOid` (a `WORKTREE_ADD` effect under the attempt's generation) and removed after
the result is recorded. There is no gate cwd in the primary checkout.

`taskman-ticket/0` (closed; every key present):

| Field | Type | Notes |
|---|---|---|
| `profile` | `"taskman-ticket/0"` | |
| `ticketId` | `ticket:*` | stable, never reused |
| `revision` | `Count` | starts `"1"`, +1 per committed mutation of any kind (chain identity) |
| `acceptanceRevision` | `Count` | starts `"1"`; +1 only when an acceptance-relevant field changes: `kind, acceptanceCriteria, requirementRefs, dependencies, requiredGates, effects, capabilities, executionClass, supersedes, source`, and on `REOPEN`, `RESTORE`, `COMPLETE_MANUAL`. Attempts, plans, gate results, reviews, manifests and dependency obligations bind this value; `revision` still bumps and still chains |
| `previousRecordSha256` | `Digest|null` | chain digest of the prior record file; null at revision 1 |
| `status` | `DRAFT|OPEN|HELD|COMPLETED|ARCHIVED` | §3.2 |
| `archivedFrom` | `DRAFT|OPEN|HELD|COMPLETED|null` | status at `ARCHIVE`; `RESTORE` returns to it; null unless `ARCHIVED` |
| `title` / `body` | prose / prose`|null` | untrusted data, never a command |
| `kind` | `FEATURE|BUG|CHORE|SPIKE|DOC|MANUAL|EXTERNAL` | |
| `owner` / `milestone` | label`|null` | |
| `priority` | `P0|P1|P2|P3` | P0 highest; default `P2` |
| `order` | `Count` | tie-break inside a priority; ties then by ticketId bytes |
| `labels` | `[label]` | sorted |
| `dependencies` | `[{ticketId, obligation:"COMPLETED"|"GATE_PASSED", gateId:label|null}]` | semantic order; `gateId` non-null iff `GATE_PASSED`. `COMPLETED` is satisfied by status `COMPLETED`, or `ARCHIVED` with `archivedFrom:"COMPLETED"`; `GATE_PASSED` is satisfied only by a `PASSED`, non-`STALE` result for `gateId` whose `ticketRevision` equals the dependency's current `acceptanceRevision`, taken from the dependency's most recent attempt at that revision; `STALE`, `NOT_RUN` or a result from an earlier revision does not satisfy it |
| `acceptanceCriteria` | `[prose]` | semantic order; index is the claim key |
| `requirementRefs` | `[Identifier]` | sorted, e.g. `ATCP-V0-003` |
| `source` | `{kind:"NATIVE"|"IMPORT", sourceQueueId, sourceItemId:Identifier|null, sourceRevisionSha256:Digest|null}` | |
| `effects` | `{coverage:"QUALIFIED"|"INCOMPLETE"|"UNKNOWN", touchPaths:[Path], resources:[Resource], externalUnbounded:boolean}` | `Resource = {class:"PATH"|"SHARED_GATE"|"SCHEMA"|"GENERATED_OUTPUT"|"PORT"|"DATABASE"|"WHOLE_REPOSITORY"|"OTHER", key:Identifier}` |
| `capabilities` | `[Identifier]` | opaque required runtime capabilities |
| `requiredGates` | `[label]` | must exist in policy; workers cannot edit |
| `holds` | `[{holdId:label, actor:label, reason:prose, placedAt}]` | semantic order |
| `executionClass` | `AUTONOMOUS|APPROVAL_REQUIRED|MANUAL|EXTERNAL|NEVER` | |
| `approvals` | `[{grantId:label, actor:label, operation:"RUN"|"COMPLETE"|"INTEGRATE"|"ADJUDICATE", targetRevision:Count, scope:[Identifier], revoked:boolean, grantedAt}]` | `targetRevision` names `acceptanceRevision`; a grant is valid only while that value is current; revocation keeps the entry. Approval, hold and completion checks use `acceptanceRevision`; `expectedRevision` optimistic checks use `revision` |
| `completion` | `{kind:"VERIFIED"|"MANUAL", actor, reason:prose|null, evidence:[Digest], manifestSha256:Digest|null, recordedAt}|null` | `MANUAL` never carries a manifest. Non-null iff the effective status (`archivedFrom` for a tombstone, else `status`) is `COMPLETED`; `REOPEN` writes `null`, and the prior completion is retained only in the chained prior record and its receipts, so readers key on status and a non-completed record can never present a completion (TCP-01 mutation repair, 2026-09-06; witnesses `TestTMV0004_AS05_ArchiveTombstoneAndReopen`, `TestTMV0004_AS05_HoldsArchiveRestoreReopen`) |
| `dueDate` / `estimateMinutes` | `YYYY-MM-DD|null` / `Count|null` | advisory only |
| `supersedes` / `supersededBy` | `ticket:*|null` | |
| `shadowOverlay` | `boolean` | true when native edits sit over an imported record pre-cutover |
| `createdAt` / `updatedAt` / `updatedBy` | timestamp / timestamp / label | actor attribution |

### 3.2 Ticket status, eligibility and roles

Status transitions: `DRAFT→OPEN` (refine with non-empty acceptanceCriteria for autonomous kinds),
`OPEN↔HELD` (hold/release-hold; `HELD` iff `holds` non-empty), `OPEN→COMPLETED` (reducer
`COMPLETED`, or `COMPLETE_MANUAL`), `COMPLETED→OPEN` (reopen, new revision; the new record
carries `completion:null` and the prior completion is kept in history, that is in the chained
prior record and its receipts, never in the active record), `{DRAFT,OPEN,HELD,COMPLETED}→ARCHIVED`
(tombstone; record retained), `ARCHIVED→` previous status (restore). Nothing else.

Effective-state mutation rules (TCP-01 mutation clarifications, 2026-09-06; mechanics frozen
under the experimental scope, no governing invariant changed): `CREATE` yields status `OPEN`
when `acceptanceCriteria` is non-empty or `kind` is `MANUAL` or `EXTERNAL`, else `DRAFT`, and a
`REFINE` that leaves a `DRAFT` with non-empty criteria (or one of those kinds) opens it. An
`ARCHIVED` record accepts only `RESTORE`; every other operation, and every `ADOPT_FILE`, is
refused `BLOCKED/TICKET_STATE`. A `COMPLETED` record accepts routine mutations but refuses any
change to an acceptance-relevant field (§3.1) `BLOCKED/TICKET_STATE` until `REOPEN`.
`SET_DEPENDENCIES`, `SET_GATES` and `SET_EFFECTS` require `DRAFT`, `OPEN` or `HELD`; `HOLD`
requires `OPEN` or `HELD`; `RELEASE_HOLD` requires `HELD`; `COMPLETE_MANUAL` requires `OPEN`;
`REOPEN` requires `COMPLETED`. `GRANT_APPROVAL` requires the payload `actor` to equal the
invoking (trusted) actor and `targetRevision` to equal the current `acceptanceRevision`, else
`UNAUTHORIZED` / `VALIDATION_FAILED`. Attempt liveness `NOT_OBSERVED` (no oracle, or an oracle
that cannot answer) is treated as live: an acceptance-relevant mutation is then refused
`BLOCKED/ATTEMPT_LIVE` and a routine one proceeds.

Eligibility is derived, never stored: `OPEN` and not held; every dependency obligation satisfied
at the dependency's current `acceptanceRevision` (§3.1 `dependencies`); `executionClass` is
`AUTONOMOUS`, or `APPROVAL_REQUIRED` with a non-revoked `RUN` grant whose `targetRevision` equals
the current `acceptanceRevision`; no live attempt; `effects.externalUnbounded` false; coverage
`QUALIFIED`, or policy `serialFallback` is `WHOLE_REPOSITORY`; `source.kind` `IMPORT` only when
`canonicalWriter` is `NATIVE`.

Mutations while an attempt is live (any non-terminal phase, including `BLOCKED_RECOVERY`): a
mutation that would change `acceptanceRevision` is refused `BLOCKED/ATTEMPT_LIVE` (cancel or let
the attempt finish first); routine mutations (`REFINE` of `title, body, owner, milestone, labels,
dueDate, estimateMinutes`, `PRIORITIZE`, `HOLD`, `RELEASE_HOLD`, `GRANT_APPROVAL`,
`REVOKE_APPROVAL`) are accepted and bump `revision` only, so a live attempt stays completable and
priority editing (ATCP-V0-009) never invalidates work. This is the explicit decision replacing
the earlier text under which every mutation invalidated the live attempt; it does not weaken
ATCP-V0-011 because grants, holds and revocations are still rechecked at every effect and
completion.

Role matrix (policy `roles` may remove rows, never add): `OWNER` all operations; `OPERATOR`
all except `GRANT_APPROVAL` for `ADJUDICATE`; `IMPORTER` only `CREATE`/`REFINE` with
`source.kind` `IMPORT`; `WORKER` only `REFINE` of `body`; `SYSTEM` only `HOLD` with
`holdId:"ESCALATED"` (placed by the tool on review escalation, §6.2); `REVIEWER` no ticket
mutations. A worker can never change priority, dependencies, gates, effects, approvals, holds,
status or completion.

### 3.3 Mutation and outcome envelopes

`taskman-mutation/0`: `profile, requestId:Identifier (1..64 UTF-8 bytes), actor:{id:label,
role:"OWNER"|"OPERATOR"|"WORKER"|"REVIEWER"|"IMPORTER"|"SYSTEM"}, queueId:queue:*,
targetId:ticket:*|null, expectedRevision:Count|null, operation, payload, issuedAt`.
`targetId`/`expectedRevision` are null only for `CREATE`.

| Operation | Payload (closed) |
|---|---|
| `CREATE` | full record minus `ticketId` (optional `localToken`), `revision`, `acceptanceRevision`, `previousRecordSha256`, `status`, `archivedFrom`, `completion`, `holds`, `approvals`, `createdAt`, `updatedAt`, `updatedBy`, `shadowOverlay` |
| `REFINE` | any subset of `title, body, kind, owner, milestone, labels, acceptanceCriteria, requirementRefs, dueDate, estimateMinutes, supersedes` |
| `PRIORITIZE` | `{priority, order}` |
| `SET_DEPENDENCIES` / `SET_GATES` / `SET_EFFECTS` | full replacement of `dependencies` / `requiredGates` / `{effects, capabilities, executionClass}` |
| `HOLD` / `RELEASE_HOLD` | `{holdId, reason}` / `{holdId}` |
| `REOPEN` / `ARCHIVE` / `RESTORE` | `{reason:prose}` |
| `COMPLETE_MANUAL` | `{reason:prose, evidence:[Digest]}` |
| `GRANT_APPROVAL` / `REVOKE_APPROVAL` | grant minus `revoked`,`grantedAt` / `{grantId, reason}` |

The `CREATE` payload is read literally from the table: it carries `supersededBy` (a reference
validated against the queue) and never `profile`; `localToken` is optional and, when absent,
the serial is allocated from `queue.nextSerial` skipping every serial whose token
`<prefix>-<serial>` is occupied, exactly or under ASCII case folding, by a live or tombstoned
ticket, and the manifest is advanced to one past the allocated serial. The scan is bounded by
the inventory cardinality and the `Count` range: a serial whose successor would leave the
range, or a queue at the §1 ticket bound, is refused `LIMIT_EXCEEDED` with the manifest
unchanged (witness `TestTMV0005_AS02_SerialAllocationSkipsOccupied`). The `actor` of every
envelope is an untrusted claim; a mutation is applied only under a trusted actor binding
supplied by the authority layer, and a claim that differs from it in id or role is
`UNAUTHORIZED`.

**ADOPT_FILE composition (R2 F3).** `reconcile intent --adopt-file` is issued by an `OWNER` or
`OPERATOR` actor with its own `requestId`. The file must parse as a complete `taskman-ticket/0`
record whose `profile` and `ticketId` equal the canonical record's; otherwise `VALIDATION_FAILED`.
The tool computes the per-field difference between the file and the canonical record and maps
every differing field to exactly one `OPERATOR`-permitted operation, in this fixed order:

| Differing field(s) | Composed operation |
|---|---|
| `title, body, kind, owner, milestone, labels, acceptanceCriteria, requirementRefs, dueDate, estimateMinutes, supersedes` | one `REFINE` carrying exactly those fields |
| `priority`, `order` | one `PRIORITIZE` (both values taken from the file) |
| `dependencies` | one `SET_DEPENDENCIES` |
| `requiredGates` | one `SET_GATES` |
| `effects`, `capabilities`, `executionClass` | one `SET_EFFECTS` (all three values taken from the file) |
| `holds` | per `holdId` present only in the file: `HOLD {holdId, reason}` with the reconciling actor (the file's `actor`/`placedAt` are ignored); per `holdId` present only in the canonical record: `RELEASE_HOLD {holdId}` |
| `updatedAt` | ignored; the adoption receipt sets it |
| any of `status, archivedFrom, completion, approvals, revision, acceptanceRevision, previousRecordSha256, source, shadowOverlay, supersededBy, createdAt, updatedBy` | refuse `VALIDATION_FAILED/ADOPT_UNSUPPORTED_FIELD`; nothing is written; the ticket stays `INTENT_DIVERGED` |

Each composed operation is validated by the same §3.3 validator, role matrix and §3.2
`ATTEMPT_LIVE` rule as if the reconciling actor had issued it (so an acceptance-relevant
difference while an attempt is live refuses `BLOCKED/ATTEMPT_LIVE` and the ticket stays
diverged). The composition is all-or-nothing: it is applied in one `RECONCILE` receipt whose
`post` entry is the resulting record with `revision` +1 exactly once, `acceptanceRevision` +1
iff any acceptance-relevant field (§3.1) changed, and `previousRecordSha256` chained from the
canonical record, never from the file. `ADOPT_FILE` can therefore never complete, archive,
reopen, restore, approve or re-source a ticket; those remain their own explicit operations, and
a file edit that attempts them is refused rather than partially adopted.

**ADOPT_FILE status, holds and request digest (TCP-01 mutation repair, 2026-09-06).** `status`
is derived, never adopted: a difference is tolerated only between `DRAFT`, `OPEN` and `HELD`,
and the composed operations decide the post status (a `REFINE` that supplies acceptance
criteria opens a `DRAFT` per §3.2, `HOLD`/`RELEASE_HOLD` derive `OPEN`↔`HELD`). The file's
status must equal the derived one, with one carve-out: a file that adds criteria to a `DRAFT`
without rewriting `status` is adopted as `OPEN`. A hand-opened `DRAFT` whose criteria stay
empty, a `DRAFT` written over an `OPEN` record, and transitions into or out of `COMPLETED` or `ARCHIVED`
are refused `VALIDATION_FAILED/ADOPT_UNSUPPORTED_FIELD`; an `ARCHIVED` canonical
record refuses `BLOCKED/TICKET_STATE` before anything is composed, so not even an identical or
`updatedAt`-only file mints a revision on a tombstone. A record remaining `COMPLETED` may
adopt otherwise permitted edits; its protected completion record is never editable. A hold present on both sides must be
byte-identical (`actor`, `reason`, `placedAt`) and the relative order of common holds unchanged;
an edited or reordered existing hold has no composed operation and is refused
`ADOPT_UNSUPPORTED_FIELD`. The TM-V0-006 request digest of an adoption is the SHA-256 of the
canonical encoding, trailing LF included, of the closed object `{actor:{id,role}, fileSha256,
operation:"ADOPT_FILE", queueId, requestId, targetId}` where `actor` is the trusted binding
of the reconciling principal, `fileSha256` the SHA-256 of the exact file bytes offered,
`queueId` the queue and `targetId` the canonical record's `ticketId`; no timestamp and no
current revision enters the preimage, so an identical retry after commit replays while the same
`requestId` reused by another actor or role, for another target or queue, or with other bytes
is `REQUEST_ID_CONFLICT`. Witnesses `TestTMV0007_AS35_DraftOpensByAdoption`,
`TestTMV0007_AS35_HoldEditsRefused`, `TestTMV0007_AS35_TombstoneRefusesAdoption`,
`TestTMV0007_AS35_RequestDigestBindsActorAndTarget`.

`taskman-outcome/0`: `profile, requestId, outcome, replayed:boolean, resultingRevision:Count|null,
resultingAcceptanceRevision:Count|null, receiptSeq:Size|null, codes:[Code]`. `outcome` is
exactly `COMPLETED | REVISION_CONFLICT | VALIDATION_FAILED | UNAUTHORIZED | BLOCKED | UNSUPPORTED
| CAPACITY_EXHAUSTED | UNCERTAIN_EFFECT | REQUEST_ID_CONFLICT | STORAGE_FAILED`. Codes are the
closed set in §11. Experimental J1 clarification (parent steering, 2026-09-06): a
recorded refusal may carry `receiptSeq`, but both resulting revisions remain null.
Uncommitted Apply/Adopt refusals still carry `receiptSeq:null`. A sequence in a decoded
outcome is only an assertion until bound to actual receipt/post evidence; it never proves
commitment. COMPLETED permits both resulting revisions or neither (mixed nullness is malformed).
The receipt binder requires exact target-afterimage revision/acceptanceRevision values for
indexed ticket successes, and null revisions for no-ticket administrative successes.
Standalone decoded fields remain assertions until that binding is checked.
Witnesses `TestTMV0002_AS01_RecordedRefusalSequence`,
`TestTMV0006_AS03_LookupErrorRefusesApplyAndAdopt`,
`TestTMV0006_AS03_RequestBindingReplayAndErrors`.

`taskman-command-result/0` (closed; the stdout envelope of every `corvint-tasks` command, one canonical
JSON document with trailing LF; non-zero exit iff `outcome` is not `OK`): `profile, command:
[Identifier] (verb path, literal), outcome:"OK"|"REFUSED"|"ERROR"|"NOT_RUN", codes:[Code],
snapshot:{headSeq:Size|null, headReceiptSha256:Digest|null, intentTreeSha256:Digest|null,
primaryWorktreeSha256:Digest|null, pendingRedo:boolean, barrier:{scope,reason}|null}|null,
mutation:taskman-outcome/0|null, items:[object] (record kind per verb, semantic order),
page:{offset:Count, limit:Count, total:Count|null, truncated:boolean}|null,
untrusted:["UNTRUSTED_QUEUE_DATA"]|[], warnings:[prose]`. One exception: `archive export`
(TM-V0-022) writes the archive stream to stdout and this envelope to stderr. Validation,
snapshot and staging failures before delivery emit zero stdout bytes. Once delivery begins,
a destination error or interruption may leave a partial archive; failure must be reported and
the consumer must discard the prefix. Stdout cannot be rolled back. `snapshot` is null for `UNINITIALIZED` and for the commands that
probe no store: `help`, `version`, a usage error refused before any read, and `archive verify`
of a stream, which is store-independent (B3 resolution, 2026-09-06; witnesses
`TestTMV0008_AS07_HelpAndVersion`, `TestTMV0008_AS07_ReadsLeaveStoreByteIdentical`). Every
other outcome carries the snapshot the command observed, partial on a failed probe. A read
whose snapshot changed while it ran reports `outcome:"NOT_RUN"` with `REDO_PENDING` or
`SNAPSHOT_MOVED` rather than mixed states (TM-V0-008). Prose inside `items` appears only as
escaped JSON strings.

**Read verb inputs and items (TCP-01 freeze, 2026-09-06; every verb below is a pure TM-V0-008
read over the intent store and the state-dir snapshot).** Paging flags `--offset N` (Count,
default `0`) and `--limit N` (Count, `1..1000`, default `100`) are validated before any read;
`page.total` is the number of rows the verb selected and `page.truncated` is true iff rows
remain after the ones returned (the next offset is `offset + len(items)`). Facts that live in
the journal (attempts, gate results, publication) are reported `NOT_OBSERVED`, never omitted,
assumed or invented; no verb below renders attempt, receipt, config, plan or replay data.

| Verb | Inputs (closed) | Items (semantic order) | `untrusted` |
|---|---|---|---|
| `ticket list` | paging | one compact ticket view per ticket in §4.3 order (`record:null`) | `["UNTRUSTED_QUEUE_DATA"]` iff items non-empty |
| `ticket search` | paging plus at least one of `--status` (§3.2 status), `--kind`, `--priority` (`P0..P3`), `--owner` label, `--milestone` label, `--label` label (membership), `--text` prose 1..512 bytes without TAB/LF/CR (case-folded substring of `title` or `body`); every given filter must match | exactly the `ticket list` items over the matching tickets | as `ticket list` |
| `ticket show <id>` / `ticket blockers <id>` | one ticket ID or local token | one view (`record` embedded / null) | `["UNTRUSTED_QUEUE_DATA"]` |
| `ticket export` | paging | per ticket in §4.3 order `{ticketId, path:"tickets/<local>.json", sha256 (file chain digest), bytes:Size, record}` where `record` is the full `taskman-ticket/0` and re-encodes to the file bytes; a page holds at most `limit` records and at most 16 MiB of encoded items, and the byte bound ends the page early with `truncated:true`; no record is ever cut | as `ticket list` |
| `queue status` | none | one object (queue, policy digest, counts, head, barrier; `attempts` and `publication` `NOT_OBSERVED`) | `[]` |
| `roadmap` | paging | one row per ticket grouped by milestone (milestones in byte order, tickets without one last, §4.3 order inside a group): `{milestone:label\|null, ticketId, title, status, kind, priority, order, owner, requiredGates:[label], gateResults:"NOT_OBSERVED", eligibility, nextAction}` | `["UNTRUSTED_QUEUE_DATA"]` iff items non-empty |
| `gate list` | none | every policy `GateDefinition` in policy order, with the §3.1 field names | `[]` |
| `gate show <gateId>` | one label | the definition plus `requiredBy:[ticketId]` (sorted) and `results:"NOT_OBSERVED"`; an undefined gate is `REFUSED` with `GATE_UNKNOWN` and no item | `[]` |

Witnesses: `TestTMV0008_AS08_Pagination`, `TestTMV0008_AS08_TicketSearch`,
`TestTMV0008_AS08_TicketExport`, `TestTMV0008_AS08_RoadmapAndGates`,
`TestTMV0008_AS07_ReadsLeaveStoreByteIdentical` (NOT_RUN until a verify run classifies them).
`config show`, `plan preview`, `receipt show|replay` and `attempt show` answer `NOT_RUN`
until the journal, pinned documents and plans exist (TCP-02+). The fourteen `ticket` mutation
verbs answered `NOT_RUN` until TCP-02b wired them to the journal (§5.7.2); the administrative
verbs still do except the delivered fixture `init`, reconciliation, `pause` and `unpause`.

Read/archive repair witnesses: `TestTMV0008_AS36_ShowRetriesDiscardPriorAttempt` proves
callback-owned output is reset per snapshot attempt. `TestTMV0022_AS09_ContentAddressedNames`
requires exact lowercase Digest filenames for `evidence/<digest>` and
`pinned/<digest>.json`, matching hashes of the consumed bytes, independently of manifest
claims. `TestTMV0022_AS09_TerminalReadError` requires actual EOF after the end marker;
`TestTMV0022_AS09_PartialDelivery` covers an accepted prefix followed by a destination error.
`TestTMV0022_AS10_EncodedStreamLimit` and `TestTMV0022_AS10_EncodedLimitIncludesPAX`
prove the complete encoded-stream bound with reduced internal test limits.

**Native receipt audit (2026-09-19).** `receipt audit` accepts no arguments. It validates the
complete native journal through J1 under the TM-V0-008 outer snapshot protocol. Its successful
`items` contains exactly one object with `headSeq:Size`, `lastReceiptSha256:Digest`,
`structuralConsistency`, `projectionAgreement`, `semanticCoverage`, `historicalAcceptance`,
`actorAuthentication`, `liveness`, `runtimeQualification` and `stagingPresent:boolean`.
String axes are the native audit's reported values: CONSISTENT, AGREES, KNOWN_CODECS or UNKNOWN
for codec coverage, and NOT_OBSERVED for historical acceptance, actor authentication, liveness
and runtime qualification. These last axes are never inferred from structural consistency.
The item is emitted only after both audit and outer snapshot succeed with equal head and intent
identities; `untrusted:[]`, `mutation:null`, `page:null`. Pending, malformed, forked, diverged or
moved observations produce no item and preserve the existing failure code/outcome mapping.
Read-only audit never repairs or replays. Witnesses `TestTMV0008_AS07_ReceiptAuditReportsBoundedEvidence`,
`TestTMV0008_AS11_ReceiptAuditRefusesWithoutRecovery`,
`TestTMV0008_AS07_ReceiptAuditUsageAndUninitialized`, and
`TestTMV0008_AS36_ReceiptAuditDiscardsMovedResults`.

**Explicit fixture reconciliation (decision 0007, 2026-09-19).**
`reconcile inspect <ticket-id>` completes one bounded native journal audit permitting divergence
only of that ticket projection, including present-empty or malformed bytes. All other projections
remain strict. Successful `items` contains one object: `ticketId`, `canonicalRecordSha256:Digest`,
`canonicalRecord` (the canonical ticket JSON), `projectionAgreement:"ALL_EXCEPT_TARGET_AGREE"`,
`historicalAcceptance`, `actorAuthentication`, `liveness`, `runtimeQualification`. The last four
axes remain NOT_OBSERVED. `untrusted:["UNTRUSTED_QUEUE_DATA"]`, `mutation:null`, `page:null`;
snapshot head, intent identity and barrier derive from the same complete audit. Errors emit no
item. This read takes no lock and never repairs, initializes or redoes state. A terminal receipt
digest is not a canonical ticket-record digest.

The closed write syntax is `reconcile intent --target ID --request-id ID --file PATH` plus exactly
one of `--keep-journal --canonical-sha256 DIGEST` or `--adopt-file`, optionally
`--role OWNER|OPERATOR` (default OWNER). `--file -` reads bounded stdin; filesystem input is
bounded and no-follow. Duplicate, unknown and inapplicable flags refuse. File paths are explicit
operator arguments, never derived from ticket prose. Before choosing, inspect the canonical
record/hash and preserve an immutable copy of the edited projection outside `.taskman`; provide
that copy as `--file`, retaining it and the KEEP digest for exact retries. The supplied original
bytes must match the current projection for a fresh transaction. A different choice under the
same request ID conflicts; exact retries can replay after later intent changes and retain the
original ticket ID. Existing §3.3/§5.5 acceptance and digest rules are unchanged.

Only settled fixture queues are supported. Full request lookup precedes fresh planning; fresh
planning uses journal-derived canonical records and a second complete audit binds head, inventory
and intent immediately before apply. Pending receipts refuse without recovery. Active staging
refuses fresh reconciliation while validated read-only replay remains available. ALL and ADMISSION permit this explicit reconciliation; branch, primary identity,
VERSION and RESTORE_INCOMPLETE guards still apply. Fresh ticket mutation and pending redo now
also refuse active staging instead of treating free-looking slots as a recovery grant.
KEEP publishes discarded evidence before receipt link-in even when the model represents it as
an evidence POST, preserving immutable destination semantics and the existing capacity accounting.

Witnesses: `TestTMV0007_AS35_ReconciliationCanonicalReadIsTargetScoped`,
`TestTMV0007_AS35_ReconciliationReadRefusesUnprovedHistory`,
`TestTMV0007_AS35_NativeKeepThenMutation`, `TestTMV0007_AS35_NativeAdoptThenMutation`,
`TestTMV0007_AS35_NativeReconcileRefusalsPreserveInput`,
`TestTMV0009_AS11_NativeReconcileRefusesPending`, `TestTMV0009_AS11_ActiveDescriptorBlocksFreshWriters`,
`TestTMV0009_AS11_KeepEvidencePrecedesCommitAndSurvivesReturnedFault`,
`TestTMV0016_AS27_ReconciliationUnderBarriersAndNoChange`,
`TestTMV0007_AS35_CLIInspectReconcileRetryWorkflow`,
`TestTMV0007_AS35_CLIReconcileInputFailuresWriteNothing`,
`TestTMV0007_AS35_InapplicableReconcileFlagRefusesBeforeInput`,
`TestTMV0008_AS35_CLIReconcileInspectFailureHasNoItem`,
`TestTMV0016_AS27_ReconcileInspectReportsJournalledBarrier`.
Returned-fault ordering/cleanup evidence is not process-crash or power-loss qualification.

**Settled fixture barriers (decision 0008, 2026-09-19).**
`pause --request-id ID [--role OWNER|OPERATOR]` and
`unpause --request-id ID [--role OWNER|OPERATOR]` expose only the existing §5.5 templates;
default role OWNER. Flags are closed; missing, duplicate, unknown and invalid-role arguments
refuse before store access. Queue identity is read from the head and revalidated under lock.
PAUSE creates ADMISSION/OPERATOR; it accepts no arbitrary scope or reason. Existing matching
pause and absent unpause return unrecorded NoChange. Exact recorded retries replay their original
outcome; reusing a request ID for another actor, role or operation conflicts. Administrative
results have no ticket or resulting revisions. All qualification coverage remains NOT_OBSERVED.

Fresh PAUSE requires complete projection agreement. Fresh UNPAUSE uses a complete settled
`Reader.BarrierRemoval` audit that reports `TICKETS_NOT_COMPARED`: canonical ticket records are
fully validated and each must have a present regular physical projection, while unrelated ticket
bytes may differ, including present-empty/malformed D. Extra/missing tickets, queue/policy drift
and private corruption still refuse. Neither this value-copy exception nor canonical records
change ordinary audit, reconciliation or request lookup behavior. A second complete audit binds
head, inventory and intent immediately before effects. ALL exempts only UNPAUSE (and the already
specified reconciliation operations); branch equality is not required for these non-intent posts.

Barrier deletion is the receipt's existing paired non-null pre/null post: publish receipt, ordinary
posts, exact-pre-digest barrier unlink plus parent sync, then head. The narrow session API can
remove only barrier.json. A deletion failure stops before head and preserves a third-value barrier.
Pending receipts refuse without redo; active staging refuses fresh operations. Postcommit UNPAUSE
failure remains pending and unsupported by **all current callable recovery paths**, including
ordinary ticket mutation redo, which refuses null posts. This preview does not qualify recovery,
process interruption, power loss, hostile editors, liveness or runtime execution.

Witnesses: `TestTMV0016_AS27_NativeBarrierCycleAndReplay`,
`TestTMV0016_AS27_UnpausePreservesMultipleDivergentTickets`,
`TestTMV0016_AS27_UnpauseALLAndBranchIndependence`,
`TestTMV0009_AS11_BarrierPublicationReturnedFaults`,
`TestTMV0009_AS11_ChangedBarrierAfterReceiptStopsBeforeHead`,
`TestTMV0009_AS11_BarrierRefusalsPreserveStore`,
`TestTMV0016_AS27_BarrierRemovalAuditRetainsStrictBoundaries`,
`TestTMV0016_AS27_CLIBarrierWorkflow`, and `TestTMV0016_AS27_CLIBarrierClosedArguments`.

### 3.4 Private state (journal), outside the repository tree

```
<git-common-dir>/taskman/           state dir; exactly one per repository, shared by every worktree
  VERSION                            "taskman-state/0\n", written once by `corvint-tasks init`
  barrier.json                       taskman-barrier/0; presence = a barrier is in force (scope inside)
  head.json                          taskman-journal-head/0
  receipts/<seq>.json                taskman-receipt/0, immutable once linked in (§5.2); never overwritten
  requests/<xx>/<sha256>.json        requestId index → {requestId, seq:Size, mutationSha256, outcome}
  attempts/<attemptId>.json          taskman-attempt/0 (current record)
  reservations.json                  taskman-reservation-set/0
  effects/<key>.json                 taskman-effect/0 pending/complete external effects
  effects/<key>.boot                 spawn bootstrap record (§6.4), link-created by exactly one writer, outside the journal
  effects/<key>.ack                  spawn acknowledgement (§6.4), link-created by exactly one writer; supervisor ack, leader no-exec, or recovery fence
  pinned/<sha256>.json               policy, config, plans, capability profiles (content-addressed)
  evidence/<sha256>                  gate outputs, transcripts, reports, manifests, blobs (content-addressed)
  staging/                           J3-02 reviewed: closed read-only transient layout (§5.6)
    active.json, active.json.tmp, a00..a10  regular children only; omitted from archive payload
  worktrees/<attemptId>-<generation>/ lane worktrees (git worktree add --detach), one per generation
<git-common-dir>/taskman.lock        flock file, outside the state dir (ATM-V0-028)
```

The pure mutation library of TCP-01 (`internal/mutation`) takes the request index as a
required explicit input alongside the queue manifest, policy, inventory and trusted actor
binding; a context without an index, including a typed-nil implementation of `RequestIndex`, is refused `VALIDATION_FAILED/MALFORMED`, never treated as
an empty index. TCP-02b supplies the `requests/` index of this directory as the transaction
model's mandatory replay observation, read under the same lock as every other input (decision
0004 §2): the model decides replay before it computes anything, binding the receipt sequence and
the stored outcome shape as well as the digest. The library keeps its own requirement
(witnesses `TestTMV0005_ContextInputsRequired`, `TestTMV0006_AS03_TypedNilIndexRequired`).

`taskman-barrier/0` (closed): `profile, queueId, scope:"ADMISSION"|"ALL", reason:"CUTOVER"|
"EMERGENCY"|"DRAIN"|"OPERATOR", actor:label, sinceSeq:Size, since:timestamp`. `ADMISSION`
(written by `pause`) refuses only `admit`, `retry`, `resume`, `scope-expand`, `import apply` and
`AUTHORITY_SWITCH`; transitions, effects, `cancel`, `reconcile`, gate/review/manifest receipts,
ticket mutations and `archive` proceed. Under an `ADMISSION` barrier whose `reason` is
`CUTOVER`, `import apply` (A4) and `AUTHORITY_SWITCH` (A5) are exempt only when invoked by
`cutover` (§5.4). `ALL` (written by `drain` and by `archive restore`)
refuses every transaction except `drain`, `reconcile`, `cancel`, `unpause` and the cutover steps
that §5.4 names. A barrier is written and removed only inside a `PAUSE`/`DRAIN`/`UNPAUSE` receipt.

`<git-common-dir>` resolution: walk up from cwd to the first `.git`; if a directory, that is the
common dir; if a file, parse `gitdir:`, then read `<gitdir>/commondir` and resolve it. Symlinks in
the resolved path fail `UNSUPPORTED_FILESYSTEM`. No environment variable or flag relocates the
state dir. Windows is unsupported.

`taskman-journal-head/0`: `profile, queueId, lastSeq:Size, lastReceiptSha256:Digest|null,
generation:Size, initSha256:Digest, primaryWorktree:PathText, versionSha256:Digest`.
`generation` is the queue-wide counter; it is written in every head and is therefore always
derivable from the head, never from the kind of the last receipt.

`taskman-receipt/0`: `profile, seq:Size, prev:Digest|null, kind, requestId:Identifier (1..64 UTF-8 bytes)|null,
actor, ticketId|null, attemptId|null, generation:Size|null, expectedRevision:Count|null,
headGeneration:Size, pre:[{path:Identifier, sha256:Digest|null}], post:[{path:Identifier,
sha256:Digest|null, record:object|null, blobSha256:Digest|null}], outcome, codes:[Code],
recordedAt`. `kind` is exactly `INIT | MUTATION | ADMIT | TRANSITION | EFFECT_INTENT |
EFFECT_OUTCOME | GATE_RESULT | REVIEW | MANIFEST | RELEASE | PAUSE | UNPAUSE | DRAIN |
IMPORT_PLAN | IMPORT_APPLY | AUTHORITY_SWITCH | ARCHIVE | RESTORE | PRUNE | RECONCILE |
ADJUDICATE | CONFIG_PIN | QUALIFICATION`. `outcome` uses exactly the §3.3 outcome enum. A `MUTATE`
transaction records the kind its operation names: `ARCHIVE` for `ARCHIVE`, `RESTORE` for
`RESTORE`, and `MUTATION` for every other §3.3 mutation. The stage operation name is `MUTATE` for
all of them, because the staged artifact shape does not vary with the verb.
Except for the explicit UNPAUSE deletion below, every post entry carries the digest of the complete
new file bytes and exactly one of `record` (inline, §1 sizing rule) or `blobSha256` (bytes in
`evidence/`, written and fsynced before the receipt), so a receipt is a self-contained redo
record. `headGeneration` is the head generation after this receipt, so replay from receipts
alone re-derives the counter. Evidence larger than a record goes to `evidence/` and is
referenced by digest.

**Experimental J1 receipt and ledger clarification (2026-09-06, reviewed plan revision 3
plus parent schema steering).** This amends unpublished local /0 mechanics only; no public
compatibility, migration, upstream acceptance or cutover is implied. Old experimental stores
are not silently upgraded. `internal/snapshot` owns the shared strict codecs;
`internal/journal` adds an independent native read library, without changing existing CLI
probe semantics or wiring `receipt audit` to a command.

Receipt paths are clean relative slash-only Identifier128 strings in the archive namespace:
`intent/queue.json`, `intent/policy.json`, `intent/import-map.json`,
`intent/tickets/<local>.json`, `VERSION`, `barrier.json`, `reservations.json`,
`requests/<first2>/<digest>.json`, `attempts/<attemptId>.json`, `effects/<Digest>.json`,
`pinned/<Digest>.json`, `evidence/<Digest>`. Exact token/attempt/digest grammars apply;
head, receipt, worktree, bootstrap/acknowledgement, temp and symlink destinations are
forbidden. Paired pre/post path sets are identical, duplicate-free and in byte order. Each
pre digest/absence is the physical redo precondition, not the canonical revision predecessor.
It is type/path/pairing checked without imposing historical pre == prior canonical post.
RECONCILE may record edited projection D as pre while an adopted ticket chains from canonical
C through previousRecordSha256. Historical reconciliation authorization, discarded-byte
retention and ticket revision/predecessor-chain semantics are NOT_OBSERVED in J1. Inline
objects encode to complete canonical file bytes including LF (≤64 KiB, ≤8); other records
and raw VERSION use retained blobs. A retained raw evidence blob may contain zero bytes;
its presence and SHA-256 remain distinct from an absent blob or deletion afterimage.
Consumed bytes must equal post/blob digests and any
content-addressed filename; manifest claims alone are insufficient. Seq starts at 1, prev is
null iff seq1, receipt1 is INIT at headGeneration0, later INIT is forbidden, and generation
never decreases. Scope identities belong to the selected queue.

The only deletion afterimage is `{sha256:null,record:null,blobSha256:null}` in an UNPAUSE
receipt for `barrier.json`, paired with a non-null pre digest. The pure `journal.DeleteRedo`
classifier returns ALREADY_APPLIED for absent, DELETE_ELIGIBLE for exact pre, and
JOURNAL_FORKED for a third digest; it never unlinks, and callers must propagate read errors
before classifying absence. PRUNE is explicitly UNSUPPORTED in J1; receipts and referenced
bytes must remain available in a contiguous 1..lastSeq chain.

Noncircular genesis: `head.initSha256` is SHA-256 of the complete canonical receipt1 bytes.
`taskman-init/0` is the closed object `{profile,queueId,primaryWorktree:PathText,
versionSha256:Digest}`. Its unique `pinned/<SHA256(complete descriptor bytes)>.json`
afterimage occurs in receipt1 only. Per explicit parent steering, minimal J1 INIT posts
**exactly** VERSION, initial queue, initial policy, empty `taskman-reservation-set/0`, and
one INIT descriptor, plus the request-index entry iff requestId is non-null. Every pre is
null. Initial tickets/import-map, attempts/effects, other extra posts and head.json are
forbidden. VERSION requires an actual retained blob containing exactly `taskman-state/0\n`;
its bytes/digest agree with descriptor and head. Decode the required genesis records and
require the selected queue/primary identity. Later heads retain that identity; relocation
is outside J1. This is not an existing-queue adoption protocol.

Every non-null requestId binds exactly one closed request post
`{requestId,seq:Size,mutationSha256:Digest,outcome:<complete taskman-outcome/0>}` at
`requests/<first2>/<SHA256(exact UTF-8 requestId bytes)>.json`; there is no receipt hash in
that entry. Entry/outcome/receipt IDs and sequences agree, original outcome has replayed:false,
and outcome/codes agree with the receipt. For COMPLETED and receipt.ticketId non-null, both
resulting revisions must equal the decoded matching target-ticket afterimage's values; a
missing target or fabricated 99/99 against posted 2/1 refuses. For receipt.ticketId null,
both revisions are null. Recorded refusals retain null revisions and a matching receiptSeq. Null requestId forbids request posts; a request
occurs once in retained history. Mutation bytes cannot be reconstructed from their digest.
The internal `mutation.RequestIndex.Lookup` contract is now `(IndexEntry,bool,error)`;
Apply/Adopt refuse STORAGE_FAILED on any lookup error and emit no fresh post/plan.
`journal.RequestIndex` audits a fresh snapshot on every lookup and exposes that observation's
head/inventory/intent identity. Only complete successful audit plus private projection
agreement can return verified absence; an earlier successful lookup does not cache away
later I/O failure. Per parent Gate A PASS, the request adapter skips only intent/* equality
and untracked-intent checks in both projection loops, allowing stable bounded divergence D
without hiding corrupt/missing request projections. Full receipt/blob/request binding,
non-intent projection checks, bounded intent reads and tuple rechecks remain mandatory;
intent projection agreement is explicitly NOT_OBSERVED. Public Audit stays strict and never
catches-and-ignores INTENT_DIVERGED. Lookup mints no branch/actor/overwrite/reconciliation grant.
The adapter is for serial use and mints neither trusted actor bindings nor a cross-snapshot
atomic transaction. No-ticket successful indexed INIT/admin outcomes carry null revisions; no fabricated ticket
revision is needed. Future reconciliation still needs a validated canonical-record retrieval
entrypoint during divergence; strict Audit(selected) refuses and its partial result is never
authority. That wiring and broader writer/request semantics remain J3 work.

The native source opens read-only through existing safeopen, lists in 256-name chunks, and
counts state metadata against the established 2,104,096 state scan budget; intent has its
independent four-root-entry and 10,000-ticket bounds. Bounded retained metadata includes
both domains, without subtracting intent entries from the supported state scan capacity. It enumerates every receipt filename, including
slots beyond head+2, before streaming receipt bytes. Aliases, malformed names, duplicates,
gaps, head-ahead state, bad links and late extras refuse. One next receipt is REDO_PENDING
only after its chain, required records, bytes and bindings validate; current destinations
must be its pre or post state. A headless valid linked genesis is pending, never empty;
missing genesis VERSION/blob fails before pending classification. With no head or linked
receipt, reads report UNINITIALIZED and staging presence without validating or cleaning
remnants. No writer safety is inferred for those remnants; stable INIT retry and artifact
charging remain held under J3.

Head, metadata inventory (including file identity/mode/mtime/size) and the streamed primary
intent digest are captured/rechecked around the body, including body errors, with at most
three retries; persistent movement is SNAPSHOT_MOVED. The latest canonical post controls
projection agreement, including every request projection; extra unjournalled intent/private
projections refuse. Unreferenced evidence/pinned bytes and temp remnants are not historical
facts and are not promoted by audit. Reads retain bounded digest/path metadata, not all
receipt bytes. Selected afterimages alone may be returned, bounded by the existing 256 MiB
aggregate live intent-state byte budget; a larger selection refuses LIMIT_EXCEEDED and must
be split. The streamed request adapter imposes no arbitrary small request-history limit.
No aggregate durable producer/capacity/escrow proof is supplied by these read bounds.

Results separate structural consistency and projection agreement from semantics. Known
codecs validate current record structure; later profiles without full J1 decoders yield
SemanticCoverage UNKNOWN. Historical acceptance, actor authentication, liveness and runtime
qualification are always NOT_OBSERVED. Archive verification continues its existing
chain/digest checks, **not full historical semantic validation**: retained archive data
currently omit required blob bytes from its semantic input. Archive historical-semantic
integration remains NOT_OBSERVED; J1 does not infer success from unavailable bytes. The
changed positive archive fixture supplies actual genesis/request proofs; existing malformed
archive negatives remain negatives.

J1 witnesses (TM-V0-002/006/007/008/009/016; AS-01/03/07/10/11/27/36):
`TestTMV0007_AS08_CanonicalLatest`, `TestTMV0009_AS11_GenesisProofCounterexamples`,
`TestTMV0009_AS11_ChainInventoryAndPending`, `TestTMV0009_AS11_QueueScopeAndGenerationDecrease`,
`TestTMV0002_AS01_ClosedReceiptPathsPairsAndDeletion`,
`TestTMV0002_AS10_ReceiptRequestIDBytesAndInitDescriptor`,
`TestTMV0002_AS01_ReceiptAndHeadSchemas` (including nine-inline refusal),
`TestTMV0006_AS03_RequestCounterexamples`, `TestTMV0006_AS03_RequestBindingReplayAndErrors`,
`TestTMV0006_AS03_StreamedLedgerExactReplayConflictAndLateIO`,
`TestTMV0016_AS27_DeletionRedoAndAudit`, `TestTMV0008_AS07_LedgerReadPurity`,
`TestTMV0008_AS36_LedgerMovementAndBodyFailure`, `TestTMV0002_AS10_LedgerBoundsAndStreaming`,
`TestTMV0002_AS10_ReceiptBytesAreStreamed`,
`TestTMV0002_AS10_NativeEnumerationStopsBeforeInspectingFlood`,
`TestTMV0002_AS01_UnknownSemanticCoverage`,
`TestTMV0002_AS10_PostBlobSizeAndNativeSymlinkRefusal`,
`TestTMV0002_AS10_EmptyRawEvidencePreservesPresence`,
`TestTMV0007_AS35_PhysicalPreIsNotCanonicalPredecessor`,
`TestTMV0006_AS03_IndexedInitAndExactTicketRevisions`,
`TestTMV0006_AS35_IndexDuringStableIntentDivergence`,
`TestTMV0002_AS10_IntentDoesNotConsumeStateScanBudget`. These do not qualify syscall durability,
full historical authorization, production RSS or performance. Parent owns the full gate and
fresh independent review; TCP-02 stays incomplete and J2/J3 held.

`taskman-attempt/0`: `profile, attemptId, ticketId, ticketRevision:Count (the ticket's
acceptanceRevision at admission), ticketRecordSha256, generation:Size, phase,
phaseSinceSeq:Size, cause:Code|null, mode:"DEVELOPMENT"|"QUALIFIED", planSha256:Digest|null,
policySha256, configSha256, runtimeId:label, capabilityProfileSha256, baseCommit:OID,
branch:label, worktreePath:Identifier|null, candidateTreeOid:OID|null, supervisor:{pid:Size,
startTime:Size}|null, lane:{pgid:Size, leaderPid:Size, leaderStartTime:Size,
spawnEffectKey:Digest}|null, quiescence:"UNPROVED"|"PROVED"|"SURVIVORS",
noExec:"PROVED"|null (non-null only for `FAILED/SPAWN`, §6.4), spawnNoExecCount:Count,
pendingEffects:[Digest], retryCount:Count, repairRound:Count, budget:BudgetUsage,
gateResults:[Digest], reviews:[Digest], manifestSha256:Digest|null,
scopeCheck:"WITHIN"|"OUT_OF_SCOPE"|"UNKNOWN", priorGenerations:[{generation:Size,
quiescence:"PROVED"|"FENCED", provedSeq:Size}]`. `BudgetUsage` carries each field as
`{value:Size|null, state:"OBSERVED"|"NOT_OBSERVED"}`.

**Supervisor identity (R2 F1).** `supervisor` is written, from the committing process's own
`pid`/`startTime`, in the transaction that records the generation's first `EFFECT_INTENT`, or in
the `ADMIT`/`retry`/`resume` transaction when the admitting process itself supervises; it is
therefore durable before any child is forked. It is set to null only by a `reconcile`, `drain` or
`cancel` transaction that has proved the recorded supervisor dead by §6.3, and it is replaced
only by the next generation's first transaction. Every supervisor-originated `TRANSITION`,
`EFFECT_INTENT` and `EFFECT_OUTCOME` carries `supervisor:{pid,startTime}`; a value that differs
from the attempt record (including null) is rejected `REVISION_CONFLICT/FENCED` and recorded,
exactly as a generation mismatch (TM-V0-011).

`taskman-reservation-set/0`: `profile, queueId, entries:[{attemptId, generation:Size, ticketId,
ticketRevision:Count, resources:[Resource] (≤4,096 after §4.2 normalization),
capacityUses:[{classId:label, units:Count}], workers:Count, state:"ACTIVE"|"QUIESCING"|
"BLOCKED_RECOVERY", createdSeq:Size, coverage:"QUALIFIED"|"WHOLE_REPOSITORY"}]`. Exactly one
entry per live attempt ID; `retry` replaces the entry for its attempt ID in the same transaction.

`taskman-effect/0`: `profile, key:Digest, attemptId, generation, kind:"WORKTREE_ADD"|
"WORKTREE_REMOVE"|"PROCESS_SPAWN"|"PROCESS_SIGNAL"|"ADAPTER_CALL", adapterId:label|null,
args:object, reconcileClass:"IDEMPOTENT"|"RECONCILABLE"|"NONE", state:"PENDING"|"DONE"|"FAILED"|
"UNKNOWN"|"RESOLVED", intentSeq:Size, outcomeSeq:Size|null, resolvedBy:label|null,
reconcileAttempts:Count`. `key` is `SHA-256("taskman-effect" || 0x00 || attemptId || 0x00 ||
generation || 0x00 || kind || 0x00 || canonical(args))`. The reconcile class is pinned per kind
(never chosen by the caller):

| Kind | Class | Key-based reconcile query at recovery |
|---|---|---|
| `WORKTREE_ADD` | `RECONCILABLE` | `git worktree list --porcelain` names `worktrees/<attemptId>-<generation>` with `HEAD` = `args.baseCommit` and the directory exists → `DONE`; neither exists → `FAILED` (retry allowed); partial → `git worktree remove --force` then `FAILED` |
| `WORKTREE_REMOVE` | `RECONCILABLE` | directory absent and not listed → `DONE`; else re-run removal only after §6.3 proves the lane dead |
| `PROCESS_SPAWN` | `RECONCILABLE` | Only `reconcile`, `drain` or `cancel`, and only after `attempt.supervisor` is dead by §6.3 (a live supervisor owns the effect: any other command leaves it untouched, `EFFECT_OWNED`). Then the §6.4 fenced-record rule: `.boot` absent → fence-create `{key, fenced:true}` (`BOOT_FENCED`); `.ack` absent → fence-create `{key, fenced:true}`; read both; no-exec proof and wrapper quiescence per §6.4 → `FAILED/SPAWN` with quiescence `PROVED` (boot pid dead) or `FENCED` (no pid recorded); `.ack` is `ack:true` → the runtime may have executed: apply §6.3 to the recorded lane; dead → `FAILED/BUILD`, live → `BLOCKED_RECOVERY/SUPERVISOR_LOST` with `lane` from `.boot` |
| `PROCESS_SIGNAL` | `IDEMPOTENT` | re-send after identity revalidation (§6.3); `SIGNAL_REFUSED_IDENTITY` ends retries |
| `ADAPTER_CALL` | declared by the adapter's pinned profile; `NONE` unless it provides a query-by-key or fence operation | adapter query by `key` |

A pending `PROCESS_SPAWN` is never retried as if idempotent, and no new spawn for an attempt is
admitted until every prior generation of that attempt is `PROVED` or `FENCED` in
`priorGenerations`.

## Requirements

- `TM-V0-001`: Exactly one repository authority exists per Git common directory (§3.4). Every
  worktree, coordinator and command of this repository resolves the same state dir and lock. A
  second state directory, a relocated state directory, a network or unqualified filesystem, or a
  symlinked path is refused before any effect. No distributed lock, remote claim or
  machine-independent ownership exists; ownership is host-local and recorded in the journal.
  (ATCP-V0-001, 007, 018; ATM-V0-005a, 028)
- `TM-V0-002`: All records implement §2 and §3 exactly: unknown keys, missing keys, duplicate
  keys, wrong types, out-of-range values, unsupported profile versions, hostile code points and
  §1 over-limit inputs fail before identity use and before any write. Identities are rederived
  before use. (ATCP-V0-002, 018; WQO-V0-001/002)
- `TM-V0-003`: A ticket record carries every §3.1 field with the stated required/optional/null
  semantics; stable IDs are never reused; each revision chains to its predecessor with actor and
  time; `acceptanceRevision` changes exactly when §3.1 says and never otherwise; supersession is
  a reference, not a rewrite; due date and estimate never affect eligibility or permission.
  (ATCP-V0-002)
- `TM-V0-004`: Task status (§3.2), attempt phase (§6.1) and gate state (§6.1) are separate
  fields with separate enums. Eligibility is derived by the §3.2 rule. No process exit, build,
  review or gate-incomplete run changes a task status; only the §7.3 reducer or an explicit
  `COMPLETE_MANUAL` does, and a manual completion is labelled `MANUAL` forever. (ATCP-V0-004/005)
- `TM-V0-005`: Mutations use the §3.3 envelope and role matrix. `expectedRevision` must equal the
  current revision or the outcome is `REVISION_CONFLICT` with no write. Validation rejects
  missing dependencies, dependency cycles (including self and via `GATE_PASSED`), duplicate or
  case-folded duplicate local tokens, unknown gates, invalid priority or order, unknown
  operations, and unauthorized roles; `VALIDATION_FAILED`/`UNAUTHORIZED` never write. Holds,
  approvals and completion cannot be changed by `WORKER`. (ATCP-V0-003, 011)
- `TM-V0-006`: A request ID is recorded with the SHA-256 of the canonical mutation bytes in the
  same transaction as its receipt. An identical retry returns the original outcome with
  `replayed:true` and writes nothing; the same request ID with different bytes fails
  `REQUEST_ID_CONFLICT`. The request index is pruned only by explicit archive-verified prune and
  never while any referenced attempt is unresolved. (ATCP-V0-003)
- `TM-V0-007`: The journal's latest post record is the canonical ticket record and the intent
  file is its projection (§3.1). Every transaction that emits an intent-file post entry, of any
  receipt kind (mutation, §7.3 reducer, `SYSTEM` `ESCALATED` hold, import apply, authority
  switch, restore, reconcile), writes that file in the primary worktree and applies the same two
  guards: the primary worktree's HEAD is on `queue.intentBranch`, else
  `BLOCKED/INTENT_BRANCH_MISMATCH`; and the file's current digest equals the journal's latest
  post digest for that path, or equals the `pre` digest of a receipt that is still redo-pending,
  else `INTENT_DIVERGED`. Divergence refuses every transaction naming that ticket, including
  admission and completion, until `reconcile intent` records `KEEP_JOURNAL` or `ADOPT_FILE`
  (§3.1); `corvint-tasks` never overwrites a diverged file on its own and never adopts it silently. Git
  commit of intent is the operator's explicit act; `corvint-tasks` never runs `git commit`. `publication`
  is reported from a read-only query as `PUBLISHED` (HEAD blob digest = file digest = journal
  post digest), `UNPUBLISHED` (file = journal, HEAD differs) or `DIVERGED` (file ≠ journal);
  a non-fixture queue admits a ticket only when its publication is `PUBLISHED` or `UNPUBLISHED`,
  never `DIVERGED`. (ATCP-V0-001, 010, 018)
- `TM-V0-008`: Read commands (`ticket list|search|show|blockers|export`, `queue status`, `roadmap`,
  `gate show|list`, `plan preview`, `receipt show|audit|replay`, `config show`, `attempt show`,
  `archive verify`, `archive export`) open no file for writing under the state dir or the
  repository, take no lock, create no directory there, never redo a
  pending receipt, never migrate `VERSION`, and report `UNINITIALIZED`, `REDO_PENDING`,
  `JOURNAL_FORKED` or `UNSUPPORTED_FILESYSTEM` verbatim. Read snapshot protocol: read
  `head.json`; check that `receipts/<lastSeq+1>.json` is absent (else `REDO_PENDING`) and that
  `receipts/<lastSeq+2>.json` is absent (else `JOURNAL_FORKED`); read the
  files; re-read `head.json` and re-check both receipt slots; if anything changed, retry at most
  three times, then report `NOT_RUN/SNAPSHOT_MOVED`. A failure while reading the files is
  re-checked the same way before it is reported: when the re-check shows the same snapshot the
  store is stable and the failure is reported verbatim (`JOURNAL_FORKED`, `MALFORMED`, ...);
  when it shows a changed snapshot the failure was a read across a commit or an intent edit
  and counts as one moved attempt, never as a fact about the store (witness
  `TestTMV0022_AS36_N4dSnapshotMovedDuringExport`). `archive export` follows exactly this
  protocol over the whole store (TM-V0-022). No output ever mixes two snapshots. Every
  item shows why it is tracked, eligible or blocked, owner, priority, current attempt, stale or
  missing gates, holds and next permitted action; `roadmap` renders milestones and per-ticket gate
  state, `gate show` one gate's definition and every result. Pagination pins `(head.lastSeq,
  intent tree digest of the primary worktree)` and reports truncation. `receipt replay <seq>`
  re-derives a `PLAN`, `REVIEW` aggregation or `MANIFEST` reducer result byte-for-byte from the
  recorded inputs and reports `IDENTICAL|DIFFERS` (ATM-V0-019). Output is the
  `taskman-command-result/0` envelope (§3.3); prose is rendered as escaped JSON labelled
  `UNTRUSTED_QUEUE_DATA`. (ATCP-V0-001, 016; ATM-V0-019/020)
- `TM-V0-009`: Every mutating command follows §5.2: exclusive `flock` on `taskman.lock`, barrier
  check, one-link chain check, chain-bounds check (§5.2: at most one receipt beyond the head,
  and it must link to the head), redo of a pending receipt, then its own single transaction whose
  commit point is the atomic, non-overwriting link-in of `receipts/<seq>.json` after fsync
  (§5.2); an existing destination is `JOURNAL_FORKED` and is never replaced. A receipt is
  observable if and only if its transition is accepted; state files never precede their receipt;
  `head.json` never runs ahead of the receipt files; a head ahead of the receipts, a broken link,
  or a receipt more than one beyond the head is `JOURNAL_FORKED` and freezes all mutation. The
  lock is never held across an external effect.
  (ATCP-V0-010, 018; ATM-V0-019a, 028)
- `TM-V0-010`: `corvint-tasks init` and every mutating command qualify the filesystem (§5.1): a working
  exclusive `flock`, local filesystem type in the allowed set, directory fsync success, and no
  symlink in the state path. Failure is `UNSUPPORTED_FILESYSTEM` and refuses the command; a read
  reports it. (ATCP-V0-018)
- `TM-V0-011`: Every attempt has an immutable ID and a `generation` taken from the queue-wide
  monotonic counter in `head.json` inside the admission transaction; `retry` and `resume` assign a
  new, larger generation to the same attempt only after every prior generation is `PROVED` or
  `FENCED` (§6.2). Every transition request, supervisor callback, gate result, review, manifest
  and adapter effect carries `(attemptId, generation, ticketRevision = acceptanceRevision,
  policySha256, configSha256)`; any mismatch is rejected `REVISION_CONFLICT/FENCED` and recorded.
  (ATCP-V0-010, 012)
- `TM-V0-012`: Admission is one transaction that checks eligibility at `expectedRevision`,
  approvals and holds, capability availability, §4.2 resource collision, capacity and budget
  headroom, plan freshness (§4.3), journal saturation and the queue's execution cutover record,
  and then writes attempt, reservation, generation and receipt together or nothing. Resources
  are normalized before comparison: `PATH` keys are canonical `Path` bytes and a `dir/` prefix
  covers every path beneath it; `WHOLE_REPOSITORY` collides with every other reservation; other
  classes collide on exact `(class,key)`. Reservations remain through review, repair, blocked
  recovery and uncertain shutdown and are compared against every live entry, not only ready
  tickets. (ATCP-V0-006/007/008; ATM-V0-005a/007)
- `TM-V0-013`: Coverage is `QUALIFIED` only when the ticket declares `QUALIFIED` and, when a
  plan is used, the plan's closure for that ticket is complete. Otherwise the attempt is blocked
  `COVERAGE_UNKNOWN`, or, when policy `serialFallback` is `WHOLE_REPOSITORY`, admitted only with
  a `WHOLE_REPOSITORY` reservation while zero other reservations are live. `externalUnbounded`
  is always blocked. Writes outside the reserved closure detected at lane exit or in the shared
  checkout quarantine the attempt (§6.2). Scope expansion is a new reservation under a new
  generation after the lane is stopped. (ATCP-V0-008)
- `TM-V0-014`: Policy caps (§1) are validated at load; every coordinator, builder, reviewer,
  verifier, repair, failed and cancelled lane's usage is accounted to its ticket and wave; caps
  on fields the runtime entry does not mark observable are `NOT_ENFORCED`; aggregate headroom is
  reserved at admission for every field in `requireEnforcedFields` (default `turns`,
  `wallClockMinutes`); a field with `NOT_OBSERVED` contributions never yields headroom; saturation
  refuses admission with `CAPACITY_EXHAUSTED` and discards nothing. (ATCP-V0-015; ATM-V0-006)
- `TM-V0-015`: Operative plans use `taskman-plan/0` and the priority-first profile in §4.3.
  `admit --from-plan` admits only entries present in the approved plan bytes and checks, inside
  the admission transaction, that the plan's `queueId`, `policySha256`, `intentTreeSha256` and
  `reservationSetSha256` equal the live values and that every entry's `ticketRevision`
  (`acceptanceRevision`) and resource set equal the live journal. Any drift refuses the **whole
  plan** with `PLAN_STALE` and admits nothing; a fresh plan is required (ATCP-V0-006 "any
  snapshot drift requires a new proposal", preserved unnarrowed). It never adds, drops or
  reorders entries. WQO/0 profiles, `MAXIMUM`/`GREEDY` labels and `queueSourceId` semantics are
  unchanged; a WQO proposal is never admitted directly. (ATCP-V0-006, 009; ATM-V0-002/002a;
  WQO-V0-044/045)
- `TM-V0-016`: Attempts and gates follow the §6 transition and recovery table exactly; an
  unenumerated transition is rejected. `pause` writes an `ADMISSION` barrier and stops admission
  only (§3.4 lists exactly what it refuses); `hold` targets a ticket; `cancel` requests
  quiescence of one attempt. No reservation, lease or worktree is released or deleted until
  quiescence is `PROVED` by the §6.3 liveness test with no pending effect, or fencing is proved
  because every adapter effect for the old generation is declared refused or, for a pending
  `PROCESS_SPAWN`, both §6.4 bootstrap records are fences. `retry`, `resume` and
  `scope-expand` require that proof first and never reuse a worktree of an earlier generation.
  Supervisor loss, survivors and uncertain effects keep ownership in `BLOCKED_RECOVERY`.
  Process groups are cleanup, never a permission boundary. (ATCP-V0-012; ATM-V0-005a/005c/028)
- `TM-V0-017`: Before any external effect the transaction records an `EFFECT_INTENT` with the
  §3.4 idempotency key; after it, an `EFFECT_OUTCOME`. A pending effect is owned by the
  supervisor recorded in `attempt.supervisor` (§3.4) while that supervisor is live by §6.3; only
  `reconcile`, `drain` and `cancel` may act on a pending effect, and only after proving that
  supervisor dead (`cancel`/`drain` against a live supervisor use the §6.2 `STOPPING` request
  instead and never touch the effect files). Any other transaction leaves a pending effect
  untouched and never fences, signals or declares it lost. A pending effect found at recovery is
  reconciled by the class the §3.4 table pins for its kind: `IDEMPOTENT` is retried at most
  `reconcileAttempts` times, `RECONCILABLE` runs the key-based query in that table (for
  `PROCESS_SPAWN` the §6.4 bootstrap protocol), `NONE` is never admitted for autonomous use. An
  adapter without a reconcile or fencing operation for a mutation is not admitted for that
  mutation. No retry of any effect for an attempt starts until the previous generation is
  `PROVED` or `FENCED`. No exactly-once claim about an external effect is made without adapter
  support. (ATCP-V0-010; ATM-V0-001c)
- `TM-V0-018`: Each runtime entry pins a `taskman-capability-profile/0` (§3.1) with per-axis
  `read|write|network|process|tools` enforcement `ENFORCED|OBSERVED|UNKNOWN` and the mechanism
  that enforces it. Qualification: `QUALIFIED` mode (the only mode that can reach `COMPLETED`)
  requires `write`, `process` and `tools` `ENFORCED` for builder and repair lanes, and `write`
  and `tools` `ENFORCED` in the profile's `readOnlyMode` for reviewer, verifier and docs lanes;
  `OBSERVED` or `UNKNOWN` on any of those axes limits the runtime to `DEVELOPMENT` mode on
  fixture queues, and a runtime without a pinned, probed profile is unavailable. After-the-fact
  detection is never enforcement: the before/after worktree digest check (`CONTAMINATED`), the
  shared-checkout dirty check and process-group cleanup are monitors that apply in every mode and
  qualify nothing. Reviewers and verifiers run in a separate detached worktree of the candidate
  tree; because a detached worktree's `.git` file necessarily points into the common dir, "no
  path to the state dir" is satisfied only by an `ENFORCED` `read`/`write` axis whose scope
  excludes `<git-common-dir>/taskman` and `taskman.lock`, recorded in the profile probe.
  Approval grants bind actor, operation, target, `acceptanceRevision` and scope; holds and
  revocations are rechecked at admission, spawn, each gate/review lane start, repair, and
  completion; an `acceptanceRevision` change invalidates grants naming the old value. Repository
  content, ticket prose and transcripts are untrusted data and never become argv, environment,
  paths or policy. Credentials never enter tickets, receipts, evidence, archives or repository
  configuration. (ATCP-V0-011; ATM-V0-008/008a/015/024a)
- `TM-V0-019`: Gate definitions live in policy and pin literal argv, environment names, timeout,
  expected predicate and required evidence; results use `taskman-gate-result/0` (§7.1) and
  record actual start, exit, signal, timeout or unknown. A result is `STALE` when candidate tree,
  definition digest, inputs digest, environment digest or policy digest differ. A `PASSED` result
  with an identical tuple may be reused with `reusedFrom` recorded. Shared canonical gates take
  their `sharedResource` reservation. Workers cannot edit definitions or required sets.
  (ATCP-V0-013)
- `TM-V0-020`: Reviews use `taskman-claim-disposition/0` and the §7.2 truth table: exactly one
  disposition per required claim, stable finding IDs with the ATM-V0-014 anchors, closed claim
  classes; missing, duplicate, malformed or fabricated dispositions never reduce the obligation
  set; empty `ACCEPT` is valid only for an explicitly empty obligation set; parser retries,
  verifiers, repair and re-review consume the aggregate budget; unresolved findings and
  insufficient independence block completion until an `ADJUDICATE` grant records the decision.
  The reviewer report bytes themselves are stored in `evidence/<reportSha256>` before the
  `REVIEW` receipt, and the claim record is derived from those bytes by the tool's deterministic
  parser, so `receipt replay` can re-derive it (ATM-V0-013a). (ATCP-V0-014, 021;
  ATM-V0-012/013/013a/014/014a)
- `TM-V0-021`: A `taskman-completion-manifest/0` (§7.3) is reduced by the deterministic §7.3
  table inside the acceptance transaction. The reducer derives every obligation itself from the
  pinned policy and the canonical ticket record (required gate set, required claim set,
  unresolved findings, CEM/OCM requirement) and validates each cited digest by loading the
  evidence, checking its profile, and checking that it binds this `attemptId`, `generation`,
  `acceptanceRevision`, `policySha256`, `configSha256` and the tree it exercised; caller-supplied
  lists never shorten an obligation and a digest that merely exists never satisfies one. It fails
  closed on missing or stale mandatory evidence; `READY_FOR_INTEGRATION` waits for an
  `INTEGRATE` grant plus recorded integration evidence before `COMPLETED`; `DEVELOPMENT` mode
  never produces `COMPLETED`. (ATCP-V0-005, 022; ATM-V0-016..019)
- `TM-V0-022`: `archive export` is a non-mutating read (TM-V0-008) that emits a
  `taskman-archive/0` stream (§3.5) on stdout holding
  `VERSION`, `head.json`, `barrier.json` if present, queue, policy, import map, every ticket
  record, every receipt, the complete `requests/` index, every attempt, the reservation set,
  every effect record and bootstrap/acknowledgement file, every pinned document and evidence
  blob, with a digest manifest. It takes no lock and creates no file, directory or lock under
  the state dir or the repository; it captures exactly one snapshot under the TM-V0-008 protocol
  (head, absent `<lastSeq+1>` and `<lastSeq+2>` receipts, and the primary-worktree intent tree
  digest, all re-checked after the copy; any change is `NOT_RUN/SNAPSHOT_MOVED` with zero bytes
  on stdout); it stages the stream in `--staging` (default the OS temp dir; refused if inside
  the state dir or the repository, or on a symlinked path), verifies the staged stream by the
  same procedure as `archive verify`, and only then copies it to stdout. Before creating a
  staging file, the exporter resolves the staging directory's enclosing Git authority with the
  bounded read-only §3.4 resolver. It refuses authority whose common directory has the same
  filesystem identity as this repository's common directory, including a linked worktree or
  case alias. Both common paths are bounded by PathText and opened as no-follow directories;
  their pinned descriptor identities are compared and inspection failures refuse. `UNINITIALIZED`
  means no enclosing Git authority and is allowed, while any other resolution failure is not proof of being outside
  and is refused. The validated outside staging directory is pinned through the no-follow safe-opening boundary; exclusive
  creation and immediate unlink use that retained directory handle. The same unlinked file
  descriptor is retained through write, seek, verification, seek and delivery, then closed;
  no staging pathname is reopened or removed during delivery. Cleanup and seek errors are
  reported, including alongside a primary error. Unlink and close failures remain sticky across
  snapshot retries; any cleanup failure observed before delivery refuses stdout even if the
  source moved. A failed unlink can leave an abandoned entry in the outside staging directory;
  it is reported, never hidden by a successful retry. Errors observed after delivery cannot
  retract bytes already delivered. This closes ordinary pathname replacement,
  not an OS owner's relocation of the containing directory. Witnesses:
  `TestTMV0022_AS36_StagingDirectorySwap`, `TestTMV0022_AS36_StagingFileReplacement`,
  `TestTMV0022_AS36_StagingFailuresBeforeDelivery`,
  `TestTMV0022_AS36_StagingCommonDirectoryCaseAlias`,
  `TestTMV0022_AS36_StagingIdentityInspectionFailsClosed`,
  `TestTMV0022_AS36_StagingUnlinkFailureSurvivesMovement`,
  `TestTMV0022_AS36_N4fStagingRefused`. Its 64 GiB limit counts the exact encoded tar stream, including manifest, headers,
PAX extensions, padding and end markers; the exporter refuses before a write would exceed
that staging limit, and verification bounds consumed stream bytes independently of payload
sums. The existing command-result item `bytes` remains the sum of `files[].bytes` (payload
only, excluding manifest and tar framing); it is not the archive limit meter. Delivery errors
report the number of stdout bytes accepted before failure; interruption may leave any prefix.
`archive verify` retains the first non-EOF source error independently of the tar decoder,
including errors returned together with complete body or trailer bytes, and refuses before
success (`TestTMV0022_AS09_BytesAndSourceError`, `TestTMV0022_AS09_TerminalReadError`).
Normal bytes-plus-EOF and the exact encoded cap remain valid.
`archive verify` reads a stream from a file or stdin, recomputes every digest, the
  receipt chain and the head generation from `headGeneration`, and fails if the stream is
  truncated, if `receiptCount` differs from `head.lastSeq`, or if any receipt beyond
  `head.lastSeq` is present (`JOURNAL_FORKED`: an archive never carries a pending or forked
  receipt); `archive restore` writes only into an empty state dir, restores every
  file above byte-identically, sets `executionCutover:null`, leaves an `ALL` barrier in place
  until the operator runs `unpause`, moves every restored non-terminal attempt to
  `BLOCKED_RECOVERY/RESTORED` with its reservation retained (no phase is resumed without
  `reconcile`), restores no credentials, and records `RESTORE`. After restore an identical
  mutation retry replays and a conflicting request ID is detected, because the index is part of
  the archive. Export/restore is lossless for identity, history and replay inputs.
  (ATCP-V0-017, 018)
- `TM-V0-023`: `import plan` inventories every source item into `taskman-import-plan/0` (§3.5)
  with a frozen `sourceDigest`, per-item mapping `NATIVE_MANAGED|FOREIGN_ADAPTER|MANUAL|HELD|
  COMPLETED|EXCLUDED` with reason, and blocking `reconciliation` entries for unknown or
  conflicting fields; `import apply` refuses on any reconciliation entry or source drift, applies
  the plan in batches of ≤100 records per `IMPORT_APPLY` receipt in plan order, resumable and
  idempotent by `import-map.json` (a record already mapped is skipped byte-identically), creates
  no duplicate, marks records `source.kind:"IMPORT"` and, before cutover, `shadowOverlay`.
  Shadow records are ineligible until the single publication boundary, the `AUTHORITY_SWITCH`
  receipt (§5.4 A5), so a partially applied import is never visible as an authoritative queue.
  Authority switches follow §5.4; a foreign adapter never acquires ownership through import.
  (ATCP-V0-001, 017, 019)
- `TM-V0-024`: Prune and drain deletion require a verified archive containing every record
  needed for replay, never remove a live, blocked or unresolved attempt, never remove the only
  copy of history, and record `PRUNE`. Retention floors in §1 apply. Migration of `VERSION` is an
  explicit resumable command that first exports and verifies an archive; reads never migrate.
  (ATCP-V0-018)
- `TM-V0-025`: Promotion of any slice requires its §9 scenarios to pass as named tests, the
  gate rows in §9.2, fresh independent review, and recorded subsequent use; a real (non-fixture)
  queue admits work only after a `QUALIFICATION` receipt names passing G2 and G3 scenarios and
  the owner's cutover record exists. Comparative and world-class claims are not delivery labels.
  (ATCP-V0-020)
- `TM-V0-026`: Not provided and reported as unavailable rather than emulated: Beamfall or any
  foreign write adapter, automatic reviewer routing, vendor SDKs, a daemon, database service, UI,
  network fetch, cross-repository scheduling, a sandbox claim for process groups, exactly-once
  external effects, and any vendor capability profile that has not been probed and pinned.
  (ATCP-V0-011, 014; ATM §5; WQO §11)
- `TM-V0-027` (owner steering 2026-09-06, exact: "I want you to ensure that this functionality
  doesn't slow down the main go binary speed"): The task control plane never slows an existing
  `corvint` command. Storage, journal, execution and planning live in `corvint-tasks` and in this
  repository. No existing `corvint` command, including `help`, `version` and startup, loads a
  native task store, initializes or executes runner code, scans a queue, loads runtime
  configuration, records scheduling events, rebuilds an index, adds a default dependency, changes
  its wire output, or starts background work; Corvint-facing operative analysis is an explicit
  opt-in verb entered lazily, and the `taskman-plan/0` planner stays in Corvint's explicit
  evidence/planning surface. Any Corvint-side change (TCP-03, TCP-06) is gated by GP (§9.2):
  a preregistered `taskman-perf-baseline/0` (§3.5) frozen before the first Corvint Go edit, then a
  directly interleaved old-Go/new-Go comparison under the §9.4 decision rule. Until GP has run,
  every performance statement is `NOT_RUN`, the unchanged baseline code path is preserved, and no
  "no slowdown", "unchanged" or "negligible" claim is made. Active `corvint-tasks` on the same host is
  covered by §9.4's contention rows; process separation alone is not evidence of absent
  contention. No existing Corvint latency budget is relaxed. (owner steering; ATCP-V0-020)
- `TM-V0-028`: Fixture queues may model ordered releases under decision 0009. The public release
  read contract exposes the complete definition, canonical candidate, attestation and promotion
  digests, and their bindings so compatible evidence and successor releases can be driven without
  reading `.taskman` projections. Mutation remains fixture-only; native gate execution,
  publication, deployment and real-queue cutover remain `NOT_RUN` and unauthorized. (owner
  decision 0009; AS-38)

### 3.5 Remaining closed schemas

`taskman-plan/0`: §4.3. `taskman-gate-result/0`, `taskman-claim-disposition/0`,
`taskman-completion-manifest/0`: §7.

`taskman-import-plan/0` (closed): `profile, queueId, sourceKind:"ROADMAP"|"FOREIGN"|"WORKLIST",
sourceRef:Identifier, sourceDigest:Digest, sourceRevision:Identifier|null, plannedAtSeq:Size,
policySha256, items:[{sourceItemId:Identifier, mapping:"NATIVE_MANAGED"|"FOREIGN_ADAPTER"|
"MANUAL"|"HELD"|"COMPLETED"|"EXCLUDED", reason:prose, localToken:label|null,
record:object|null (the CREATE payload, null unless NATIVE_MANAGED/MANUAL/HELD/COMPLETED),
sourceRevisionSha256:Digest|null, batch:Count}] (semantic order),
reconciliation:[{sourceItemId, field:Identifier, sourceValue:prose, reason:Code}],
counts:{total,nativeManaged,foreignAdapter,manual,held,completed,excluded:Count}`. `counts.total`
must equal `len(items)` and equal the source inventory; `import apply` refuses unless
`reconciliation` is empty and the source digest is unchanged.

`taskman-archive/0` (a stream: POSIX pax tar written by the standard library, whose first
entry is `manifest.json`, followed by every file in `files` order, then the tar end-of-archive
marker; `manifest.json` closed): `profile, queueId, exportedAtSeq:Size,
headSha256:Digest, versionSha256:Digest, barrierSha256:Digest|null, primaryWorktree:PathText,
intentTreeSha256:Digest,
files:[{path:Identifier, sha256:Digest, bytes:Size}] (sorted; every state-dir file plus
`intent/queue.json`, `intent/policy.json`, `intent/import-map.json`, `intent/tickets/*`),
receiptCount:Size, lastReceiptSha256:Digest, headGeneration:Size, complete:boolean`. A stream
without a leading `manifest.json`, with `complete:false`, with a missing end-of-archive marker,
or whose entries differ from `files` in order, path, size or digest never verifies.
`exportedAtSeq` equals `head.lastSeq` and `receiptCount` of the snapshot.

**Archive capacity (B4 resolution, 2026-09-06: an unverified implementation freeze under the
TCP-01 experimental freeze; the flat closed `taskman-archive/0` schema is unchanged).** `files`
is bounded by the §1 archive entry bound of 2,100,000 entries: the journal saturates at
1,000,000 receipts, so the ordinary 10,000 decoded-array bound could never hold a saturated
store. `manifest.json` is parsed only through the named opt-in entry point
`archive.ParseManifestDocument` (`wire.ParseWith`), which widens exactly the top-level `files`
array (≤2,100,000 elements) and the node cap (≤4 × 2,100,000 + 14 = 8,400,014: the manifest
object, its thirteen members and four nodes per entry) and nothing else: depth stays 24, every
other array in the manifest and every other document keep 10,000 / 250,000, and `wire.Parse`
is unchanged. The manifest byte cap is 768 MiB, checked on the tar header and on the byte
length before any parse. Arithmetic: the largest valid entry is `{"bytes":"` 9 + Size 20 +
`","path":"` 10 + encoded path ≤256 (128 Identifier bytes, each `"` escaped to two bytes;
backslash and controls are refused by `Path`) + `","sha256":"` 12 + Digest 64 + `"}` 2 + `,` 1
= 374 bytes; 374 × 2,100,000 = 785,400,000; the envelope outside `files` (thirteen keys and
punctuation, five Digests, three Sizes, the profile, `queueId` ≤128, `primaryWorktree` ≤8,192
encoded, brackets and LF) is under 9,100 bytes; 785,409,100 < 805,306,368 (768 MiB), whereas
the first proposed 512 MiB (536,870,912) would not hold it, so 768 MiB is the cap. Element and
node caps fail `LIMIT_EXCEEDED` while parsing, before the document is materialized. Duplicate,
unknown and missing keys, Unicode, framing, `Path`, sortedness and order validation are exactly
those of every other document; nothing is truncated. `archive export` enumerates the state dir
under a separate scan bound of 2,104,096 entries (files + 4,096): it lists every directory in
chunks of 256 names, counts every entry it sees (directories, skipped `*.tmp-*` and
`head.json.tmp` entries, `worktrees/` and unexpected names included) against that scan bound
and refuses `LIMIT_EXCEEDED` as soon as it is passed, without validating, stat'ing or opening
the entries beyond it; the exported files (state-dir files plus every intent file) are then
counted against the 2,100,000 entry bound before any byte is read. A store beyond either bound
is refused whole, never truncated or rewritten; every store within them, including one already
holding more receipts than a future writer cap admits, stays exportable. Witnesses:
`TestTMV0022_AS10_ParseWithWidensOnlyTheNamedArray`,
`TestTMV0022_AS10_ExportEntryBoundDuringListing`,
`TestTMV0022_AS10_ExportBeyondTenThousandMembers`,
`TestTMV0022_AS10_VerifyBeyondTenThousandReceipts` (small realistic fixtures; no
multi-million-entry fixture is built). Worst-case flat-manifest memory (a 768 MiB manifest is
held as bytes, as a parsed value and as its re-encoding during the canonical check, and
`archive verify` retains every receipt and intent file for the chain check) is NOT_RUN: no
low-overhead claim is made, and a segmented manifest would need a later versioned profile only
if measurement shows flat parsing exceeds the qualified memory budget.

**J2a pure archive encoding-size mechanics (2026-09-06; experimental partial).**
`archive.MeasureManifestEncoding` measures a supplied typed `Manifest` without materializing
the complete manifest or any member body. It returns `Files = len(files)` (excluding
`manifest.json` and PAX extension entries), `PayloadBytes = sum(files.bytes)`,
`ManifestBytes = M` including LF, and the exact `TarBytes = header(manifest.json,M) +
round512(M) + sum(header(path,bytes) + round512(bytes)) + 1024`. `M` is the production encoding
of the same metadata with `files:[]`, plus each production `entryValue` encoding and the
inter-entry commas. Header bytes come from a fresh standard-library tar writer using the exact
production header constructor, so PAX headers and their padding are included; an unfinished
header-only writer is not closed. All addition, rounding and int64 conversion are checked, and
the existing file-count, manifest and exact-tar caps apply. The typed metadata and every member
must independently satisfy the existing archive scalar, Identifier, Path, digest, Size,
per-path byte, reserved-name and unique-path rules. An empty list is measurable arithmetic but
does not verify as a complete store. This API proves encoding size only: inventory completeness,
digest/content agreement, scan and transient entries, evidence and journal totals, admission,
terminal/recovery/admin reserves, filesystem durability, restore, runtime and performance are
NOT_OBSERVED. J2b aggregate inventory/prefix/escrow mechanics are an experimental reviewed model
under §5.5; every J3 durable writer remains HELD and TCP-02 is incomplete and GP (TM-V0-027) is NOT_RUN. Witnesses:
`TestTMV0022_AS22_MeasureManifestEncodingParity`,
`TestTMV0022_AS22_MeasureMetadataWidthAndNullParity`,
`TestTMV0022_AS22_MeasureRejectsTypedAliasesAndMembers`,
`TestTMV0022_AS22_MeasureCapsAndCheckedArithmetic`,
`TestTMV0022_AS22_MeasureDoesNotAllocateClaimedBodies`,
`TestTMV0022_AS22_ArchiveHeaderFactoringPreservesBytes`.

**TCP-02 writer prerequisite (B4; mandatory before any durable producer is promoted).** The
archive entry bound, the manifest byte cap, the archive total (64 GiB), the journal saturation
thresholds (1,000,000 receipts or 4 GiB, unchanged, and still the ADMISSION refusal of
TM-V0-012) and the evidence store total are independent hard aggregate budgets. Before EVERY
durable producer (receipt link-in, post file, evidence blob, pinned document, bootstrap
`.boot`, acknowledgement `.ack`, `archive restore`) the writer must prove that cost(proposed
store) + the reserved remaining terminal/recovery costs of every live attempt + the reserved
admin costs of `archive restore`, `unpause` and `prune` fit the file-count bound, the manifest
byte cap, the exact tar size (per-entry 512-byte header, the PAX extended header when one is
needed, body padding to 512 and the 1,024-byte end marker, all counted), the receipt
count/bytes thresholds and the evidence cap. Unique retained paths and content are counted
once each; the reservation is taken before the work starts and released only when its output
is durably consumed (linked in, archived or pruned). The terminal/recovery headroom R
(receipts) and B (bytes) are derived from the bounded recovery paths of §6.2 and §6.4 (the
maximum number and size of `TRANSITION`, `EFFECT_OUTCOME`, `RELEASE`, `RECONCILE`, fence and
acknowledgement records one live attempt can still emit), never invented; the same derivation
covers the admin path (`RESTORE`, `UNPAUSE`, `PRUNE` and their barrier records). A store that
already exceeds a proposed writer cap is never silently truncated or rewritten: the writer
refuses new work (`JOURNAL_SATURATED` / `LIMIT_EXCEEDED`) and the store remains exportable.
Every future durable-writer promotion is held until this accounting exists in code with
saturation and restore+unpause regressions (AS-09, AS-10, AS-27) that drive a store to each
budget and prove refusal, recovery and export. Pure TCP-01 code writes no native state and is
complete with this prerequisite recorded for the next slice.

`taskman-perf-baseline/0` (closed; frozen by `OWNER`/`OPERATOR` before any Corvint Go edit,
pinned by `CONFIG_PIN`): `profile, frozenAt, corvintCommit:OID, corvintTreeOid:OID,
binarySha256:Digest, goVersion:Identifier, os:Identifier, arch:Identifier, host:{cpuModel:prose,
cpuCount:Count, memoryBytes:Size, loadCondition:"IDLE"|"CONTENDED"|"UNKNOWN", concurrentGates:
Count|null}, corpora:[{name:label, repositoryCommit:OID, indexStateSha256:Digest, taskStateSha256:
Digest}], cacheStates:["COLD"|"WARM"], commands:[{name:label, argv:[Identifier], corpus:label,
cacheState, expectedStdoutSha256:Digest|null, expectedExit:Count}], samplesPerRow:Count (≥100),
warmupsPerRow:Count (≥5), order:"ALTERNATING", witnesses:["WALL_P50","WALL_P95","CPU_USER",
"CPU_SYS","RSS_MAX","ALLOC_BYTES","IO_READ_BYTES","IO_WRITE_BYTES"], budgets:[{command:label,
witness, absoluteMax:Size}] (the existing Corvint budgets, never relaxed), atmConditions:
["STOPPED","IDLE_INITIALIZED","MAX_ADMITTED"], atmBudgets:{coordinatorCpuPercent:Count,
coordinatorRssBytes:Size, childRssBytes:Size, ioBytesPerSecond:Size, processes:Count,
pollingIntervalSeconds:Count}|null, decisionRule:"SPEC-9.4"`. `atmBudgets` is null until
measured; a null budget is `NOT_RUN`, never a pass.

## 4. Admission, reservations and planning

### 4.1 Admission order (single transaction)

1. No barrier of any scope; journal not saturated; filesystem qualified; chain link valid;
   primary worktree matches `head.primaryWorktree`.
2. Ticket exists at `expectedRevision`; §3.2 eligibility at its current `acceptanceRevision`;
   publication not `DIVERGED`; `executionCutover` present unless `queue.fixture` is true (then
   `mode` is `DEVELOPMENT`).
3. Runtime entry exists with a pinned capability profile; required `capabilities` ⊆ profile;
   TM-V0-018 qualification for the requested `mode`; for `retry`, every prior generation of the
   attempt is `PROVED` or `FENCED` and no effect is `PENDING`.
4. Budget: for each `requireEnforcedFields` field, `Σ live reserved + lane cap ≤ wave cap`, and
   ticket aggregate `Σ prior attempts observed + lane cap ≤ ticket cap`; any `NOT_OBSERVED`
   contribution to a required field blocks `BUDGET_UNKNOWN`.
5. Capacity: `len(live) + 1 ≤ maxActiveAttempts`, workers `+ lane workers ≤ maxWorkersTotal`,
   per-class units fit.
6. Resources: normalized set intersects no live entry; `WHOLE_REPOSITORY` needs zero live entries.
7. Plan (if any) is fresh per TM-V0-015 (whole-plan check).
8. Write attempt (`ADMITTED`), reservation (for `retry`: replace the entry with the same
   `attemptId` atomically), `generation+1`, receipt `ADMIT`.

Any failed step yields the first applicable outcome in this order: `STORAGE_FAILED`,
`UNSUPPORTED`, `REVISION_CONFLICT`, `BLOCKED` (with code), `CAPACITY_EXHAUSTED`. Nothing is
reserved on failure.

### 4.2 Collision normalization

`PATH` keys are compared after WQO `Path` validation: equal bytes collide; `a/` collides with
`a/` and any `a/...`; a file and its own prefix directory collide. A ticket's resource set is
`touchPaths ∪ resources ∪ (plan closure for that ticket when present)`. Unparsed, untracked or
unindexed paths contribute only themselves (WQO-V0-042). Two attempts never share a resource
while both are live, including `QUIESCING` and `BLOCKED_RECOVERY`.

### 4.3 Operative planning profile `taskman-plan/0` (TCP-03)

`profile, planningProfile:"taskman-priority-first/0", queueId, policySha256, headSeq:Size,
intentTreeSha256:Digest, reservationSetSha256:Digest, capacity:{maxActiveAttempts,
availableWorkers}, entries:[{ticketId, ticketRevision (acceptanceRevision), resources:[Resource],
closureComplete:boolean, state:"SELECTED"|"DEFERRED"|"BLOCKED", reason:Code,
deferredSinceSeq:Size|null, blockers:[Code|ticket:*]}], mutationAuthority:false`.

Algorithm: order eligible tickets by `(priority asc, order asc, ticketId bytes asc)`. The first
ticket that is runnable (eligible, complete coverage or serial fallback, no collision with live
reservations, fits capacity) is `SELECTED` first and is never displaced. Continue in the same
order, selecting each ticket that collides with neither a live reservation nor an earlier
selection and fits remaining capacity; otherwise `DEFERRED` with reason. Tickets that are not
eligible are `BLOCKED` with their blockers. `deferredSinceSeq` is the `headSeq` of the earliest
pinned plan (any plan recorded by `plan record`, a `CONFIG_PIN` receipt, or by an admission) that
deferred the same `(ticketId, acceptanceRevision)` without an intervening plan that selected it;
`plan preview` computes it from the pinned plans it can read, so deferral age is reported even
for a queue that has only been previewed and recorded, and is null only when no pinned plan has
deferred the entry. No automatic priority promotion exists; a hold is never bypassed. This
profile is distinct from WQO-V0-044 and carries no `waveOptimality`.

## 5. Storage protocol

### 5.1 Filesystem qualification

Allowed filesystem types: Darwin `apfs`, `hfs`; Linux `ext4`, `xfs`, `btrfs`, `tmpfs`. Anything
else, including network and FUSE mounts, fails. `flock(LOCK_EX|LOCK_NB)` then unlock must
succeed; `ENOLCK`/`EOPNOTSUPP`/`ENOTSUP` fail. File durability: `fsync` on Darwin uses
explicit `F_FULLFSYNC` under a pinned descriptor, retrying only `EINTR`, and never falling
back to `fsync` on `ENOTSUP`. Go 1.27.0 `File.Sync` has such a fallback and is not used for
Darwin file durability. Directory sync is a separate `fsync` after every publication.
Linux magic `0xEF53` is shared by ext2/ext3/ext4: it is reported ambiguous and refused; the
allowlist is unchanged, but no ext4 qualification is claimed without descriptor-bound
disambiguation evidence. Unsupported platforms refuse before any lock or temp creation.

The opening boundary walks every absolute and relative directory component using atomic
`openat` with `O_NOFOLLOW|O_DIRECTORY`, pins the resulting descriptor, and retains identity
for publication and qualification. Final file opens also use `O_NOFOLLOW|O_NONBLOCK` so a
stat-to-FIFO swap cannot block before type validation. No later open re-traverses the original
absolute path. The internal descriptor-only `/dev/fd/<fd>` conversion to `os.Root` is accepted
only while that source descriptor is pinned and only after same-directory identity checks;
missing or non-equivalent bridges refuse. Generic `os.Root` traversal alone does not establish
no-follow behavior. Lock acquisition checks cancellation/deadline on every flock retry,
after acquisition and before returning a held lock. Local tokens remain case-fold unique
because the default Darwin filesystem is case-insensitive.

Repair witnesses (executed results recorded in the repair report):
`TestTMV0010_AS10_DarwinStrictFullSyncBoundary`, `TestTMV0010_AS10_AmbiguousExtMagic`,
`TestTMV0010_AS10_UnsupportedPlatformBeforeEffects`,
`TestTMV0001_AS10_AuthorityDirectoryReplacement`,
`TestTMV0001_AS10_DirectorySwapNoRedirectOrBlock`,
`TestTMV0001_AS10_PinnedRootSurvivesPathReplacement`,
`TestTMV0001_AS10_DescriptorBridgeFailsClosed`, `TestTMV0001_AS10_DescriptorLifetime`,
`TestTMV0009_AS10_LockInterruptedDeadlineAndCancellation`,
`TestTMV0008_AS07_ReadBoundedFIFORegression`, `TestTMV0008_AS07_ReadInRootFIFORegression`.

### 5.1.1 Experimental J3-01 private fixture mechanism (2026-09-06)

This Gate-A-cleared increment adds only unexported `fixtureSession` filesystem primitives
under `internal/authority/fixture_session*.go`. Only tests in package `authority` construct
them, on disposable repositories created and retained by the harness. The harness is the
sole writer and keeps directory/mount topology stable. Raw paths, a resolved repository,
`fixture:true`, same UID and flock confer no ownership, actor, administrative or transaction
grant. These primitives write test bytes, including receipt-shaped filenames; they produce
no real receipt or trusted transition. No production caller, exported operation, CLI, journal
sequencer, runtime or J2b dependency is added. The public authority prerequisite statements
continue to hold; this private mechanism is not their callable-writer implementation.

The constructor re-resolves the exact primary/common/state/lock tuple, requires a matching
live acquired `Lock`, and retains its file identity. It pins primary, common and every
present admitted parent through existing `safeopen.Root`, `SubRoot` and `InRoot`. Absent
state/intent children remain absent until an explicit primitive creates them; a constructor
failure closes its handles and leaves lock release with the acquiring caller. Successful
construction adopts responsibility for explicit session closure. Every operation locks
**session mutex, then `Lock.mu`**, retaining both through all native calls, identity checks,
sync, cleanup and temporary-handle closes. Session close takes that same order indirectly
through `Lock.Close`; it waits for an operation, closes all retained handles and releases the
lock, joining failures. Its first result is sticky on repeated close. Public `Lock.Close`
keeps its first-call reporting/idempotence semantics: if an external close wins, its caller
receives that release error. No background finalizer or cleanup process stands in for close.
Zero, wrong, closed or observed-replaced identities refuse.

Directory authority is a closed internal key, never a caller-provided directory FD or an
absolute destination path. The admitted hierarchy is common `taskman/` (state); state
`staging/`, `receipts/`, `evidence/`, `pinned/`, `requests/`; request shards exactly lowercase
`00`..`ff`; primary `.taskman/` (intent) and its `tickets/`. The constructor retains existing
parents, including existing shards; `mkdir` creates only an absent admitted child and requires
both child and parent directory sync. It never creates/replaces the Git common directory,
adopts an entry that appeared later, recursively cleans, or removes a directory. Unknown
entries do not become layout acceptance. Initial inventory/capacity and actual staging-reader
recognition remain separate J2b/J3-02 work.

The only stage names are `active.json.tmp`, `active.json` and `a00`..`a10`.
`active.json.tmp` prepares nonempty descriptor bytes at most 2,422 bytes; `active.json` can
only be linked exclusively from that prepared temp. It is not a direct write slot. Payload
slots admit closed roles with the existing per-file limits: receipt/queue 1 MiB; evidence/
reservations 64 MiB; pin 16 MiB; request 64 KiB; VERSION exactly `taskman-state/0\n`; head/
barrier 4 KiB; policy 256 KiB; import map 8 MiB; ticket 128 KiB. Only raw evidence permits
zero bytes, retaining its empty SHA-256 identity independently of absence. Limits and roles
are checked before create/write or payload allocation by this mechanism. Preparation creates
exclusively with no-follow, confirms a full write and exact length/digest/metadata, performs
strict file sync plus staging-directory sync, and requires successful close before returning
a session-owned token. Partial, short, sync-failed, close-failed or identity-moved preparation
never returns a publishable token. Failure cleanup verifies the known complete or partial
bytes and inode before unlink and sync; detected foreign entries are retained.

Publication consumes a session-owned token. Targets are canonical receipt sequence names,
exact content-addressed evidence/pin names, canonical request digest names in their matching
shards, VERSION, head, barrier, reservations, and the closed intent names (queue, policy,
import-map and valid local-token ticket names). Content-addressed filenames must match the
staged bytes. Descriptor bytes cannot become another role. A bounded payload may also be
linked as digest-matching evidence; this proves byte identity only. Native cross-parent
`linkat` is exclusive for every existing entry type and keeps the source slot. Only private
head/barrier/reservations may be checked-replaced: exact expected pre permits rename, equal
post is already applied (strict file and directory sync still required), a third value is
preserved. Null pre uses exclusive link. Rename requires destination **and source** directory
sync, attempting both even if one fails. Controlled removal covers only an exact-pre barrier
(or already-absent barrier with directory sync) and session-owned stage names. Receipt,
evidence, pin, request and VERSION identities never replace/remove. Existing intent overwrite
and removal are forbidden, including same-byte entries; absent intent may link exclusively.
This is no atomic content-CAS against unrelated editors and does not implement C6.

Every native call pins both parent descriptors and its open source (and existing destination
for replacement) with `safeopen.Control`. Observed parent/name/inode/mode/mtime/size/digest
drift refuses before effects and cleanup, including required directory-sync boundaries.
Already-applied paths recheck the named post or absence before directory sync. Filesystem checks use retained-descriptor device,
filesystem **and mount identity**, never `st_dev` alone. Darwin uses `Fstat` plus `Fstatfs`
(type, nonzero fsid, local flag and mounted-on/from identity); Linux uses `Fstat`/`Fstatfs`
plus a bounded, unique nonzero `mnt_id` from the pinned descriptor's `/proc/self/fdinfo`
observation. Unknown/ambiguous observation or mismatch refuses. Linux ext-family ambiguous
magic remains refused. Missing procfs/descriptor bridge or unsupported platforms refuse;
`EXDEV` has no copy/unlink fallback. File sync preserves the existing strict Darwin
`F_FULLFSYNC`/Linux fsync boundary; directory durability is separate.

Private `fixtureEffectError` reports a performed create/link/rename/unlink/mkdir whose
required completion failed, joining primary, cleanup, sync and close errors. A successful
link followed by failed destination-directory sync is visibly published with retained source,
not “nothing happened”. Rename/unlink/mkdir errors after effects preserve the analogous
on-disk distinction. This type says nothing about transaction commitment. In-process hooks
exercise finite failures; no interruption child is introduced. These are syscall-result and
identity witnesses, not whole-store recovery, physical power-loss durability, historical
acceptance or production performance qualification.

Named finite witnesses (actual results in the J3-01 builder report):

| Gate A category | Named TM/AS tests |
|---|---|
| Session identity/lifetime | `TestTMV0001_AS10_FixtureSessionIdentity`, `TestTMV0009_AS10_FixtureLifetimeAndClose`, `TestTMV0009_AS10_FixtureTransientHandlesAndCleanupRetention` |
| Stage roles, limits, foreign entries and preparation failures | `TestTMV0002_AS10_FixtureStageRolesAndBounds`, `TestTMV0009_AS10_FixtureStageForeignAndFaults` |
| Exclusive links, admitted names and mount/source preservation | `TestTMV0009_AS11_FixtureExclusiveLink`, `TestTMV0002_AS10_FixtureDestinationRoles`, `TestTMV0010_AS10_FixtureMountBoundaries`, `TestTMV0001_AS10_FixtureReplacementPreservation` |
| Private replace/remove and directory creation | `TestTMV0009_AS11_FixtureMutableReplace`, `TestTMV0009_AS11_FixtureRemove`, `TestTMV0010_AS10_FixtureMkdir`, `TestTMV0007_AS10_FixtureClosedDestinations` |
| Before/after effects, joined errors and required durability | `TestTMV0009_AS10_FixtureJoinedFailures`, `TestTMV0010_AS10_FixtureActualSyncFaultsAndAlreadyApplied`, `TestTMV0009_AS11_FixtureSyncBoundaryPreservesForeignNames`, plus the primitive fault tables above |
| Compatibility/platform | existing `TestTMV0009_LinkInPublishesCompleteFilesAndValidatesInput`, `TestTMV0009_AS36_N4c_LinkInNeverReplacesDestination`, `TestTMV0009_LinkInCleansUpOnInjectedSyncFailure`; new `TestTMV0010_AS10_FixtureAbsentHierarchyAndUnsupported`, Linux-only `TestTMV0010_AS10_FixtureLinuxMountObservation`, other-platform `TestTMV0010_AS10_FixtureUnsupportedPlatform` |

TCP-02 remains incomplete. J3-02 read support is reviewed (§5.6); test-owned J3-03/04
sequencing/recovery is authored (§5.7), pending combined verification and fresh review.
J3-05 is partly delivered under decision 0003 (2026-09-07): the local-operator issuer, the
`OBSERVED_LOCAL_OPERATOR` premise, the callable writer of `internal/store` and `atm init`.
Existing-intent CAS/C6 against a hostile concurrent editor, capacity qualification, runtime and
real (non-fixture) queues remain held. GP is NOT_RUN; unsupported
host tests are NOT_RUN rather than durability passes. No Corvint source is edited.

### 5.7.1 Callable writer and genesis (decision 0003, 2026-09-07)

`internal/store` applies one `transaction.Plan` to a real repository authority in the §5.2 order:
evidence and blob-backed bytes, then the receipt link-in (the commit point), then the post files,
then the head. Each artifact is staged into a `staging/` slot and fully synced before publication,
and every slot a transaction occupies is cleared before the call returns, including on a failure
part-way through: a slot left occupied with no live descriptor is an unassigned slot that every
later reader refuses (witness `TestTMV0009_AS11_InitLeavesNoUnassignedStageSlot`).

The §5.2 redo rule is applied per post destination before publishing it: a destination already
holding the post bytes is settled and is not written again; an absent destination is written; a
destination holding the `pre` digest the receipt records for that path is written; a destination
at any third value is never overwritten and is reported `INTENT_DIVERGED` for an intent
projection or `JOURNAL_FORKED` for a state file (witnesses
`TestTMV0009_AS35_InitDoesNotOverwriteAnEditedProjection`,
`TestTMV0009_AS35_RedoDoesNotOverwriteAnEditedProjection`). The pre digest is read from the
committed receipt, which is the evidence a crash recovery reads, so a first write and a redo of
it decide identically. An absent destination is link-created; a destination holding the pre bytes
is replaced by renaming the staged file over it under a compare-and-swap on exactly those bytes,
so a file that changed between the read and the write is refused rather than overwritten. The
request index shard `requests/<xx>`, and the intent store's `tickets` directory in a queue that
has never held a ticket, are created the first time a transaction publishes into them.

A receipt that was linked in but whose post files or head were not written (crash point C2) is
completed before any new transaction is modelled, under the same lock: the receipt is a
self-contained redo record, so its post entries supply the bytes inline or name a blob already
published in `evidence/`. Redo advances the head to the receipt it completed and never
renumbers or rewrites it (witness `TestTMV0009_AS11_RedoCompletesAPendingReceipt`).

`corvint-tasks init` qualifies the filesystem, takes the exclusive lock, creates the §3.4 directories
`receipts`, `evidence`, `pinned`, `requests` and `staging`, and commits the genesis `INIT`
receipt. It reads `.taskman/queue.json` and `.taskman/policy.json` as the operator's own authored
input and refuses when they are absent, before creating anything. An existing state dir is
reported as the model's `BLOCKED` refusal with nothing written, never as a filesystem failure.
The lane directories `attempts`, `effects` and `worktrees` are not created: the runtime slice has
not been built.

The recorded local-operator actor comes first from `CORVINT_TASKS_ACTOR`, with `ATM_ACTOR` as a
compatibility fallback. Equal nonempty values agree. Different nonempty values are `MALFORMED`
before repository resolution, payload reads, locking, or mutation. When both are absent, the OS
user lookup and label validation remain the decision 0003 behavior.

### 5.7.2 Ticket mutations reach the journal (decision 0004, 2026-09-07)

Every §3.3 ticket mutation commits through one `MUTATE` transaction carrying one
`taskman-mutation/0` envelope, delivering TCP-02b. `internal/store.Mutate` reads the queue,
policy, complete canonical ticket inventory, the state-dir inventory and this request's index
entry under the exclusive lock; `transaction.Model` decides replay and computes the plan through
`mutation.Apply`; the §5.2 writer publishes it. The fourteen `ticket` mutation verbs are wired to
it and no longer answer `NOT_RUN`.

The CLI composes the envelope so a caller never hand-writes the profile, queue id or timestamp:
the verb decides the operation, `--payload` (or `--payload-stdin`) carries the closed §3.3
payload as canonical JSON, `--request-id` is the idempotency key, and `--target` with
`--expected-revision` names the record for every verb except `CREATE`. Because the request digest
is the digest of the envelope bytes, `--issued-at` lets a retry reproduce the original envelope
and replay; without it each invocation issues a new request. A `CREATE` that allocates a serial
reports the ticket id it allocated; an exact replay returns that same ID from its validated original receipt, with no new receipt.

Witnesses: `TestTMV0005_AS02_CreateCommitsATicket`,
`TestTMV0005_AS02_RefineChainsFromTheCommittedRevision`,
`TestTMV0005_AS02_StaleExpectedRevisionLeavesTheStoreByteIdentical`,
`TestTMV0006_AS03_IdenticalRetryReplays`,
`TestTMV0006_AS03_SameRequestIDDifferentBytesConflicts`,
`TestTMV0004_AS05_HoldsArchiveAndRestore`,
`TestTMV0008_AS02_TicketCreateIsVisibleToTheReadVerbs`,
`TestTMV0006_AS03_TicketCreateRetryReplays`,
`TestTMV0008_AS07_MutationsRefuseOnAnUninitializedStore`.

Held: attempt liveness is answered by the zero-attempt oracle, which is sound only while no
runtime exists to start an attempt; the remaining administrative verbs and `archive restore`
remain unwired. Settled fixture intent reconciliation is specified below under decision 0007.

Not delivered here: crash recovery through the staging descriptor (`active.json`), which belongs
with `reconcile`; administrative mutations other than `init`; and any non-fixture queue, which still
requires the §7.4 cutover record and a `QUALIFICATION` receipt.

#### Writer repair boundary (2026-09-19)

Before recovery, the callable ticket writer binds the envelope actor ID/role to the invoking
local-operator binding and refuses unsupported VERSION, any RESTORE_INCOMPLETE marker,
primary-worktree mismatch and an ALL barrier. ADMISSION permits ticket mutations. These
checks do not authenticate the local operator or qualify a runtime.

A pending receipt authorizes recovery only after the complete native journal audit returns
its terminal REDO_PENDING result with CONSISTENT structure and PRE_OR_POST projections.
The writer binds the audited head, intent tree and terminal receipt digest before redo;
earlier-history corruption, foreign scope, decreasing generation and broken request bindings
cause no recovery writes. Other partial audit results never authorize a write or replay.

Settled request lookup validates the journal and private projections while permitting stable
intent divergence. An exact replay or request-ID conflict precedes fresh branch/projection
requirements. Original ticket identity comes from that same validated receipt walk. Fresh
requests obtain canonical queue, policy and ticket bytes from a strict journal audit, then
bind them to the model inventory. This conservatively refuses the whole fixture queue on
any intent divergence until reconciliation is wired. Request and intent selections have
separate bounded reads; no existing aggregate limit is raised.

INIT validates its model under the exclusive lock before creating state directories, and
checks for an intervening initializer. INIT and fresh mutations observe the primary
common-directory HEAD through a bounded no-follow read. A detached, missing, unreadable or
wrong primary branch refuses; a linked caller's own HEAD cannot substitute. The writer
rechecks the branch and captured intent/head before effects. These cooperative observations
do not qualify atomic CAS against hostile concurrent editors. Interrupted genesis and
staging-descriptor recovery remain held.

Named regressions (the test identifiers carry TM-V0 and AS traceability):
`TestTMV0007_AS35_WriterRejectsProjectionDrift`,
`TestTMV0007_AS35_WriterObservesPrimaryBranch`,
`TestTMV0007_AS29_LinkedCallerUsesPrimaryHEAD`,
`TestTMV0009_AS11_WriterGuardsBeforeRedo`,
`TestTMV0009_AS11_RedoRequiresCompleteJournalProof`,
`TestTMV0009_AS11_InvalidInitLeavesCorrectableInput`,
`TestTMV0009_AS11_InvalidInitRoleThenRetry`,
`TestTMV0006_AS03_ReplayPreservesIdentityAndAllowsStableDivergence`,
`TestTMV0006_AS03_ActorBindingBeforeReplayAndRedo`,
`TestTMV0006_AS03_ForgedRequestNeverReplays`,
`TestTMV0016_AS27_AdmissionBarrierAllowsNativeMutation`, and strengthened
`TestTMV0016_AS27_BarrierTemplatesExemptionsAndNoChange` /
`TestTMV0006_AS03_TicketCreateRetryReplays`.

### 5.2 Transaction and commit points

```
acquire taskman.lock exclusive (≤30 s, else STORAGE_FAILED/LOCK_TIMEOUT)
verify head.primaryWorktree == resolved primary worktree (else UNSUPPORTED)
read barrier.json; refuse per its scope (§3.4): ADMISSION refuses admit/retry/resume/
  scope-expand/import apply/AUTHORITY_SWITCH with BLOCKED/PAUSED; ALL refuses everything except
  drain, reconcile, cancel, unpause and the §5.4 steps that name themselves
verify receipts/<lastSeq>.json digest == head.lastReceiptSha256 (else JOURNAL_FORKED)
verify receipts/<lastSeq+2>.json is absent (else JOURNAL_FORKED: chain bounds)
if receipts/<lastSeq+1>.json exists: verify its seq == lastSeq+1 and prev == head.lastReceiptSha256
  (else JOURNAL_FORKED); REDO it (see redo rule below), advance head
delete stray receipts/*.tmp-* files and head.json.tmp (R3: never effects/*.tmp-*, see below)
validate inputs; for every intent post file apply the TM-V0-007 branch and divergence guards;
  compute post records; write evidence blobs and blob-referenced post bytes (fsync) first
write receipts/<seq>.tmp-<pid>; fsync; link(2) tmp -> receipts/<seq>.json (EEXIST => JOURNAL_FORKED,
  nothing committed, tmp removed); unlink tmp; fsync dir                        <- COMMIT POINT
for each post file: write .tmp, fsync, rename, fsync dir   (intent file in primary worktree too)
write head.json.tmp {lastSeq:seq, lastReceiptSha256, generation, ...}; fsync; rename; fsync dir
release lock
```

Temp cleanup (R3): the recovery step above removes only `receipts/*.tmp-*` and `head.json.tmp`.
An `effects/*.tmp-*` file is removed only by its own writer, or by recovery (`reconcile`,
`drain`, `cancel`) after the owning supervisor is proved dead by §6.3 and that effect key has
been fenced per §6.4; no other transaction touches it.

Receipt destinations are immutable: the commit uses `link(2)` rather than `rename(2)` so that a
pre-existing `receipts/<seq>.json` (for example one restored from an archive taken after this
head) can never be replaced; the caller sees `JOURNAL_FORKED` and the store is unchanged. The
same link-in rule (write temp, fsync, `link`, unlink temp, fsync dir) is the single exclusive
creation primitive used for `.boot` and `.ack` in §6.4, so a linked-in file is always complete.

Redo rule (C2..C4, C6): for each post entry, read the destination's current digest first. Equal
to the entry's `sha256` → nothing to do. Equal to the receipt's `pre` digest for that path (or
the path is absent and `pre` is null) → write the post bytes. Any third value → for an intent
file, mark the ticket `INTENT_DIVERGED` and skip that file (the receipt stays committed, the
projection waits for `reconcile intent`); for a state-dir file, which only `corvint-tasks` writes, stop
and report `JOURNAL_FORKED`. Redo never overwrites a file it has not proved to be at the pre or
post state, so a concurrent or manual edit is never destroyed.

External effects run strictly between transactions: `EFFECT_INTENT` transaction, effect,
`EFFECT_OUTCOME` transaction. A supervisor whose transition write is refused (`PAUSED`,
`FENCED`, `JOURNAL_FORKED`, `LOCK_TIMEOUT`) records nothing, leaves the lane in its current OS
state, retries the write on the next heartbeat, and after the lane's wall-clock cap reports
`SUPERVISOR_LOST` semantics by exiting; the attempt is then recovered by `reconcile`, never by a
transition invented from process exit status.

### 5.3 Crash matrix

| Point | On-disk state | Recovery (next mutating command, under lock) | Read reports |
|---|---|---|---|
| C1 before non-overwriting receipt link | owned staging may exist; linked orphan evidence may remain | remove only receipt/head temps permitted by §5.2; fixture staging cleanup requires §5.5 item 5 ownership and durability proofs; retain immutable evidence; no receipt committed | existing committed view, or `UNINITIALIZED` without a head |
| C2 after non-overwriting receipt link, before post files | receipt N+1, head N (or no head for INIT) | redo by the exact PRE/POST/third-value rules of §5.2, then advance head | `REDO_PENDING`; archive remains `UNINITIALIZED` without a head |
| C3 some post files renamed | same | leave exact POST bytes; apply POST only to exact PRE, or absence when PRE is null; preserve third intent values and refuse third private values under §5.2 | `REDO_PENDING` |
| C4 post done, head not renamed | same | advance head | `REDO_PENDING` |
| C5 after head rename | committed | none | normal |
| C6 intent file digest differs from the journal's latest post digest, or from the pre digest while that receipt is redo-pending (covers a Git checkout/reset of the file, a manual edit, and a read from any worktree other than the primary) | conflict | mark ticket `INTENT_DIVERGED`; refuse every transaction naming it (mutation, admission, reducer, hold, import, switch) until `reconcile intent` records `KEEP_JOURNAL` or `ADOPT_FILE` | `INTENT_DIVERGED` |
| C6a supervisor/reducer/`SYSTEM` transaction finds the primary worktree off `intentBranch` | no write | refuse `INTENT_BRANCH_MISMATCH`; attempt stays in its phase; supervisor retries per §5.2 | `INTENT_BRANCH_MISMATCH` |
| C7 operator never commits intent to Git | working tree ahead of HEAD | none; publication is the operator's act | `UNPUBLISHED` |
| C8 crash between EFFECT_INTENT and EFFECT_OUTCOME | effect `PENDING` | only `reconcile`/`drain`/`cancel`, and only after §6.3 proves `attempt.supervisor` dead: attempt `BLOCKED_RECOVERY/UNCERTAIN_EFFECT`; §TM-V0-017 reconcile by the §3.4 class table. While the supervisor is live every other command leaves the effect untouched (`EFFECT_OWNED`) | `UNCERTAIN_EFFECT` only with the supervisor dead; otherwise the attempt's phase with `pendingEffects` listed |
| C8a supervisor dies after `PROCESS_SPAWN` intent and after the leader wrapper was created, before the `RUNNING` receipt | `.boot` present or absent, no `.ack` | §6.4 fenced-record rule by `reconcile`/`drain`/`cancel` after the supervisor is proved dead: fence whichever of `.boot`/`.ack` is absent, so the wrapper exits without exec; `FAILED/SPAWN` with `PROVED` (boot pid dead) or `FENCED` (no pid recorded); reservation retained until then | `UNCERTAIN_EFFECT` |
| C8b supervisor dies after the `RUNNING` receipt, before linking `.ack` | lane recorded, wrapper waiting, `.ack` absent | reconcile fences `.ack` (`{key, fenced:true}`); wrapper exits without exec; lane dead by §6.3 → `FAILED/SPAWN` (never `FAILED/BUILD`: no `ack:true` exists); retry needs a new generation | `SUPERVISOR_LOST` |
| C9 authority switch step | §5.4 row | resume from last receipt | `CUTOVER_IN_PROGRESS` |
| C10 archive export interrupted | truncated stdout stream or abandoned staging files; nothing under the state dir or repository | `archive verify` of the stream fails (no end-of-archive marker or manifest mismatch); source untouched; rerun export | n/a |
| C10a archive taken while a transaction commits | export's re-check sees a changed head, a pending receipt or a changed intent tree | `NOT_RUN/SNAPSHOT_MOVED`, zero bytes on stdout; a torn archive is never produced | n/a |
| C10b restore of an archive whose receipts exceed its head | never verifies (`receiptCount` ≠ `lastSeq`, or `<lastSeq+1>` present) | `archive restore` refuses before writing; after any restore, a later commit at an occupied `seq` is `JOURNAL_FORKED`, never an overwrite | `JOURNAL_FORKED` |
| C11 restore interrupted | `RESTORE_INCOMPLETE` marker | only `archive restore` rerun or manual delete of the target | `RESTORE_INCOMPLETE` |
| C12 lock holder dies | OS releases flock | next holder runs recovery | normal |
| C13 head ahead of receipts, link mismatch, pending receipt whose `seq`/`prev` do not continue the head, or `receipts/<lastSeq+2>.json` present | fork | refuse all mutation `JOURNAL_FORKED`; restore from archive | `JOURNAL_FORKED` |

### 5.4 Authority switch and interruption matrix

Cutover from `ROADMAP` (or `FOREIGN`) to `NATIVE`, each step one receipt:

| Step | Action | Interrupted here → state | Resume | Revert |
|---|---|---|---|---|
| A1 | `pause --reason CUTOVER` writes an `ADMISSION` barrier | barrier present, writer unchanged | rerun `cutover` | `unpause` |
| A2 | reach quiescence under the `ADMISSION` barrier: `cancel` each live attempt or wait for it to reach a terminal phase (transitions, callbacks, cancel and reconcile all proceed under `ADMISSION`); `cutover` stops with `QUIESCENCE_UNPROVED` until every reservation entry is `PROVED` or released | as A1 | rerun | `unpause` |
| A3 | `archive export` (read, no lock, TM-V0-022) + `verify` | as A1, truncated stream or `SNAPSHOT_MOVED`; nothing written under the state dir | rerun export | `unpause` |
| A4 | `import plan` frozen digest; `import apply` in ≤100-record batches, each its own `IMPORT_APPLY` receipt (exempt from the `ADMISSION` refusal only when invoked by `cutover` and the barrier `reason` is `CUTOVER`) | some batches written as shadow, writer unchanged | rerun apply (idempotent by `import-map.json`) | records stay as shadow, ineligible |
| A5 | `AUTHORITY_SWITCH` receipt with post `queue.json` (`canonicalWriter:NATIVE`, `writeBarrier` on the old source); the single publication boundary for every shadow record (exempt from the `ADMISSION` refusal on the same terms as A4; a stand-alone `AUTHORITY_SWITCH` under any barrier is refused `PAUSED`) | C2..C4 rows | redo | `AUTHORITY_SWITCH` back with archive of post-cutover history |
| A6 | `unpause` removes the barrier | complete | none | n/a |

Revert performs A1..A3 and then A5 in reverse, keeps every native ticket and receipt created
after cutover as history, and never deletes records. `EMERGENCY` barrier: any external action on
the queue source freezes admission until `reconcile source` records the difference.

### 5.5 Experimental J2b pure transaction/capacity model (2026-09-06)

`internal/transaction` proposes local mechanics for a **hypothetical no-runtime fixture**
only. This is not upstream acceptance, a filesystem writer, CLI wiring, actual admission,
archive-layout acceptance, a real queue, restore or qualification. J1 (including present-empty
raw evidence) and J2a are the accepted base; their existing Go sources and read/export behavior
are not changed. Parent steering and its independent Gate A layering PASS place the single
shared descriptor codec in new `internal/snapshot/stage.go`; transaction aliases/reuses it.
This avoids a future archive -> transaction -> archive cycle without activating layout reads. J2b remains experimental; its repaired model passed the parent gate and fresh review. TCP-02 is incomplete.

1. **Explicit inputs and scope (TM-V0-002/005/008/014/018; AS-03/07/10/35).** `NewInventory`
   privately copies supplied regular-file path/digest/length metadata and state directories;
   it verifies bounded shape and internal agreement, never physical existence. Intent files
   count in archive/intent budgets, not state scans. Runtime/attempt/effect/bootstrap/worktree
   files and legacy temps are unsupported by this model, even where existing reads support
   them. Queue/policy/head/barrier/empty reservations and the complete canonical ticket set
   are supplied as captured bytes. Only `fixture:true`, no configured runtimes, no cutover or
   import map, and explicit `HYPOTHETICAL_FIXTURE_NO_RUNTIME` premises enter fresh templates.
   Missing, malformed, unknown or incompatible facts refuse. Inputs are not modified;
   inventory and plan accessors return copies. There are no filesystem, clock, environment,
   network reads, I/O callbacks, arbitrary liveness or index implementations in the API.
   `ActorAuthentication`, `AdministrativeAuthorization`, `InventoryObservation`, `Durability`
   and `RuntimeQualification` always report `NOT_OBSERVED`, including replay and NoChange.
   OWNER/OPERATOR is a hypothetical shape restriction, not an administrative policy grant.
   No administrative policy row, authorization codec or default-allow boolean is introduced.

2. **Closed original requests (TM-V0-002/006/007; AS-03/35).** Digest is SHA-256 of the
   canonical closed object below, including LF; actor is `{id,role}`. Inapplicable request
   fields refuse. All identifiers and captured bytes retain their existing bounds. INIT's
   queue/policy bytes are complete canonical desired inputs, not regenerated from a clock.
   VERSION is exactly `taskman-state/0\n`. RecordedAt and current head/revisions are excluded.

   | Operation | Exact preimage members (in addition to actor, operation, queueId, requestId) |
   |---|---|
   | INIT | policySha256, primaryWorktree, queueSha256, versionSha256 |
   | PAUSE | reason:"OPERATOR", scope:"ADMISSION" |
   | UNPAUSE | none |
   | KEEP_JOURNAL | canonicalSha256, fileSha256, targetId |
   | ADOPT_FILE | fileSha256, targetId; digest is `mutation.AdoptDigest` exactly |
   | MUTATE | the request is the `taskman-mutation/0` envelope itself; the digest is the SHA-256 of its canonical bytes exactly (TM-V0-006), `issuedAt` included, so a retry replays only by reproducing the same envelope |

   KEEP binds original canonical C and original physical D; a retry never silently refreshes
   either choice. Replay observation is mandatory and explicitly ABSENT or FOUND (original
   canonical request-index bytes). Missing/unknown/error refuses STORAGE_FAILED; it never
   means absent. FOUND first passes strict source decoding, canonical encoding,
   request/outcome identity, receipt-sequence binding and the 1,000,000 sequence bound. Once
   those generic source checks pass, a different original digest is REQUEST_ID_CONFLICT before
   any incoming-operation outcome-shape check. With the same digest, the stored COMPLETED
   outcome must still carry null revisions for INIT/PAUSE/UNPAUSE and paired revisions for
   reconciliation, with no codes. Exact replay returns the full original outcome with
   replayed:true before fresh-transition checks, with no receipt/index/post/artifact.
   Supplied replay evidence is still hypothetical. J3 must obtain qualified fresh J1 index
   observations and a trusted actor/administrative binding before consuming these artifacts.

3. **Five templates (TM-V0-003/005/006/007/009/016; AS-01/03/27/35).** Fresh success has
   non-null requestId, exactly one request-index post, COMPLETED, codes=[], recordedAt from
   the frozen input and seq=base.lastSeq+1. AttemptId/generation are null, head generation is
   preserved (genesis zero). No caller-authored fresh successful outcome is accepted.
   INIT/PAUSE/UNPAUSE have null ticket/expectedRevision and null resulting revisions.
   Reconciliation names its target, expectedRevision=C.revision and both posted revisions.
   NoChange is unrecorded COMPLETED, replayed:false, receiptSeq/revisions null and codes=[];
   replay, refusal and NoChange produce zero fresh receipt/index/staging growth.

   | Template | Fresh posts and decision |
   |---|---|
   | INIT | Exactly VERSION, queue, policy, empty reservations, pinned taskman-init/0 and request; all pre null. Head.initSha256 hashes receipt1. Existing new-request initialization is BLOCKED with no code. Retained precommit orphan evidence is preserved and deduplicated; genesis destinations must be empty. |
   | PAUSE | Create ADMISSION/OPERATOR barrier and request. Existing ADMISSION/OPERATOR is NoChange, preserving actor/time/seq. Any other barrier is BLOCKED/PAUSED; never replaced. |
   | UNPAUSE | Paired non-null physical barrier pre and null deletion post, plus request. Absent barrier is NoChange. |
   | KEEP_JOURNAL | Restore byte-identical C; pre hashes present regular D. Retain D as evidence POST with sha256=blobSha256 and record:null, including empty/malformed D bounded to 128 KiB. Already equal C/D is NoChange after original-choice checks. |
   | ADOPT_FILE | Existing `mutation.Adopt` supplies digest, validation, policy composition, finalization and C-prime. Revision increments once, acceptanceRevision follows existing rules, predecessor hashes C. Empty accepted composition still increments revision. |
   | MUTATE | Existing `mutation.Apply` supplies validation, the role matrix and the post record for every §3.3 verb. Posts are the ticket projection and the request entry, plus `intent/queue.json` when a CREATE allocated a serial. A CREATE has no predecessor: its receipt carries `expectedRevision:null` and a `pre` entry with a null digest. Every other verb chains from the record it read and names that revision. |

   Reconciliation is one target with full bounded ticket inventory validation. It checks
   the supplied primary branch. UNPAUSE alone does not require unrelated physical ticket
   projections to equal their canonical records: an abandoned reconciliation must preserve C,
   D and linked orphans while still permitting the reserved barrier removal. This exception
   does not admit an absent physical projection or weaken canonical decode, queue scope,
   queue/policy/head/barrier/inventory binding, complete ticket inventory or any other
   operation's projection-equality checks. The private in-memory absent index and zero-attempt fixture
   oracle are used only inside this declared hypothetical subset. Every existing ADMISSION
   and ALL reconciliation/UNPAUSE exemption is preserved; this does not change cutover,
   drain, cancellation or any other operation outside the subset.

4. **Proposed taskman-stage/0 (TM-V0-002/009; AS-01/10/11).** The closed descriptor is
   `{profile,queueId,operation,requestId,requestSha256,recordedAt,base,artifacts}`. Base is
   null only for INIT, otherwise `{lastSeq:Size,lastReceiptSha256:Digest}`, lastSeq 1..999999.
   Artifacts are semantic slot order, each exactly `{slot,role,target,sha256,bytes:Size}`.
   Sort by (role,target) bytes, then assign consecutive a00..a10. Roles are EVIDENCE, HEAD,
   POST, RECEIPT. Targets are exactly the operation's receipt, derived head, non-deletion
   posts and required new blobs; duplicate paths/roles/slots or arbitrary destinations refuse.
   An evidence POST serves its own blob once and publishes before receipt. Other new blobs
   require promised bytes and dedup against the retained inventory; equal bytes at different
   paths otherwise cost separately. Full plans bind the actual receipt, preimage, outcome,
   head and afterimage encoders. Structural descriptors alone are not complete transactions.

   | Operation | Slots at most | Actual encoder structural maximum, bytes including LF |
   |---|---:|---:|
   | INIT | 11 | 2422 |
   | PAUSE | 4 | 1243 |
   | UNPAUSE | 3 | 1098 |
   | KEEP_JOURNAL | 6 | 1676 |
   | ADOPT_FILE | 5 | 1467 |
   | MUTATE | 6 | 1615 |

   These are structural upper witnesses, not claims that every maximal payload coexists in a
   reachable transaction. Ticket/queue widths obey queueBytes+localBytes<=126. INIT's closed
   post bounds are VERSION16, queue1048576, policy262144, reservations194, pinned INIT8465,
   request551; administrative request563 and reconciliation request579 apply in the other
   rows. Receipt1048576/head4096 and existing ticket/blob bounds apply; minimal UNPAUSE
   receipt1683 is tighter. Actual head encoding must still fit 4096. The only proposed flat
   staging names are active.json, active.json.tmp and a00..a10: at most thirteen files plus
   staging/ (fourteen scanned entries). They would be excluded from F/M/T/Jn/Jb/Eb and charged
   to D/Pb. **Reviewed J3-02 read recognition is specified in §5.6.**

5. **Preparation and cleanup (TM-V0-008/009/016; AS-11/27).** Publish/fsync a descriptor before
   payload production. Missing/partial assigned regular slots are PreparationPending when
   lengths do not exceed their promises; complete wrong hashes are corruption. A matching
   full frozen plan plus complete bytes is FullPlanPrepared. An owned descriptor temp is
   bounded by that frozen descriptor's actual length and, if complete, must match it. A
   malformed active.json is never authority. Complete retained receipt metadata and head
   binding are required. Matching linked next receipt is RedoPending; verified post/head
   agreement is CompletedCleanup. These are conditional model classifications, not read CLI
   statuses or permission to repair a physical path.

   Selected Gate A exception: no active.json, no payload/unknown entries, stable unchanged
   base, complete no-next/no-extra receipt inventory and explicit no-runtime fixture premises
   permit only regular active.json.tmp of length 0..2422, without parse/full-hash requirements,
   as TempPreparationAbortable. Empty staging is EmptyPreparation. Unlink this fixed temp,
   then fsync staging before a fresh descriptor. Cap+1, symlink/nonregular, stray slot, changed
   base or next/extra receipt refuses. With a valid active descriptor, remove payload/temp
   names, fsync staging, unlink active.json last, fsync staging again. Names and logical bytes
   remain charged until their fsync; promised afterimages discharge only at durable descriptor
   removal. No new ownership era starts early. Old projections, all receipt history and all
   linked orphans remain. No age cleanup, abort receipt, history deletion or new blob occurs.
   J3 must prove sole-writer/type/identity/ownership and durability with retained handles.

6. **Capacity closure (TM-V0-002/009/014/022; AS-09/10/27).** The old hard-extension proposal
   is superseded: hard limits remain 1,000,000 receipts and 4 GiB. A state retaining/creating
   a barrier must reserve one minimal UNPAUSE *inside* them: Jn+1<=1000000 and
   Jb+1683<=4294967296. Already saturated barrier stores without reserve are ineligible for
   this model; their existing readability is not changed and nothing is truncated.
   Other hard bounds stay F2100000, D2104096, M768MiB, T64GiB, Eb16GiB, intent256MiB and all
   individual file limits. All arithmetic is checked; no unknown contribution is zero.

   Membership changes use explicit paths and immutable inventory transforms. Exact retained
   F/payload/M/T use production J2a `MeasureManifestEncoding`; directories/transients count
   separately in D. Journal/evidence/intent bytes are accumulated from complete supplied
   metadata. Pending-head envelopes measure promised bounded redo, not exportable snapshots.
   Every ordered evidence-publication prefix is checked for abort cleanup followed by fresh
   UNPAUSE, keeping old D and all linked orphan blobs. New request shard and head decimal
   growth are included. Growing post publications precede shrinking ones in the conservative
   size order. All required posts, including barrier deletion, precede head publication.
   A deletion/shrink cannot finance an earlier producer. The separately labeled UNPAUSE peak
   retains the old barrier with the new receipt/index and head replacement budget; it is a
   conservative envelope, not an invalid head-before-deletion durable publication point.

   Maximal staging-name/byte reservations dominate all bounded preparation and cleanup
   subsets; retained membership at each listed prefix is measured exactly. Pb is logical
   temporary charging, not measured physical blocks or preallocation. Outstanding artifact
   afterimages are reported separately as promised bytes and held through cleanup fsync.
   Actual minimal UNPAUSE codecs reproduce receipt/index/outcome 1683/563/308. Its three
   artifact caps sum 6342=1683+563+4096 and two descriptors give Pb8538=6342+2*1098, with five
   stage files plus staging/. No constant was expanded. Matching pending redo reuses its
   linked receipt and checks pre/post/third-value guards (including `journal.DeleteRedo`);
   it allocates no additional receipt or request identity and preserves retained history.
   Focused abort witnesses invoke the actual fresh UNPAUSE model after durable cleanup for
   both KEEP and ADOPT prefixes, preserve present-empty or valid divergent D plus multiple
   unrelated divergent projections and every linked orphan, apply exact capacity, then prove
   identical replay and a later barrier-absent NoChange add no receipt.

7. **Evidence limits and holds (TM-V0-010/018/022/027; AS-09/10/27/36).** Focused new-package
   tests cover actual descriptor/UNPAUSE encoder parity, immutable small inventory transforms,
   production arithmetic thresholds without million-receipt/4-GiB allocations, private scaled
   capacity limits, every one of 2^11 structural slot subsets, cleanup durability, corrupt and
   temp-only preparation, abort/orphan/no-shrink and receipt-preserving redo. This is modest
   and scaled pure evidence, not production RSS/timing, actual inventory observation, crash
   durability or real qualification. Parent owns one final full gate and a fresh independent
   reviewer. J3/writer promotion remains held for qualified fresh J1 observations, trusted
   actor/admin binding, filesystem ownership/qualification, same-filesystem state staging
   and primary intent rename proof (or separately accepted alternative), read/export staging
   amendment, publication/recovery/cleanup fault injection and saturation/export evidence.
   Runtime, live-attempt escrow, restore/prune, real cutover and GP remain held/NOT_RUN.
   Owner decision 0002 selects fixture-only first delivery; it does not qualify any issuer,
   filesystem durability or runtime, and this J2b increment remains pure.
   Repair witnesses `TestTMV0006_AS03_ReplayConflictPrecedesOperationShape`,
   `TestTMV0007_AS35_UnpauseDivergenceValidationBoundary` and
   `TestTMV0014_AS27_AbortCleanupPermitsFreshUnpauseAcrossDivergence` cover the two confirmed
   independent-review findings. See `docs/reviews/2026-09-06-j2b-builder.md` and
   `docs/reviews/2026-09-06-j2b-repair.md` for exact test/file evidence and failed hypotheses.

### 5.6 Experimental J3-02 bounded read-only staging observation (2026-09-06)

Accepted experimental implementation under decision 0002 and the amended J3-02 Gate A plan.
The parent source freeze, canonical gate and fresh independent review passed; see
[final J3-02 integration](reviews/2026-09-06-j3-02-integration.md). This acceptance excludes
J3-03/04 fixture integration, separately accepted in §5.7. TCP-02 is incomplete. The shared J2b descriptor codec,
transaction model, J3-01 authority primitives, safeopen, snapshot Reader retry contract,
archive manifest and wire limits retain their accepted source bytes.

The closed state layout additionally admits `staging/`, with at most thirteen regular
children: `active.json`, `active.json.tmp`, `a00`..`a10`. No aliases, nested paths, symlinks
or special children are admitted. Every enumerated name, including unexpected names,
consumes D; enumeration refuses before inspecting children beyond min(remaining D,13).
The staging directory itself consumes one D; its directory identity is not a second charge.
Staging children consume no archive file entries or payload. Every linked evidence/pinned
orphan remains in the inventory and export; descriptor references never filter history.

One shared pure observer uses the existing strict descriptor codec. Absent/empty staging
is inactive. A sole descriptor temp may contain arbitrary bytes, including empty or malformed
bytes, up to 2422. With active present, temp is bounded by active's actual byte length;
shorter temp is incomplete preparation, equal length requires exact equality (otherwise
JOURNAL_FORKED), and excess is LIMIT_EXCEEDED. Assigned payload reads use decoded per-artifact
bounds; absent/shorter slots are incomplete, exact length requires the declared digest, and
zero-length raw evidence remains present with the empty digest. Unknown/unassigned slots
and malformed active descriptors refuse. Only digests and the bounded descriptor survive
structural validation; captured staging payloads are not retained as archive afterimages.

Nonempty staging is recognized only for a decoded fixture queue. This is a content restriction,
not ownership or actor authority. Precommit requires the actual head's base sequence and
receipt digest. Completed leftovers require actual HEAD/RECEIPT bytes matching the descriptor,
actual receipt sequence/prev/generation/request binding, descriptor/receipt timestamp equality,
and operation/kind equality (INIT→INIT, PAUSE→PAUSE, UNPAUSE→UNPAUSE, both reconciliation
labels→RECONCILE), plus actual receipt-bound request
bytes/digest/length and original request digest. Inline and bounded retained-blob request
proof use existing strict codecs. Missing proof or stale/foreign relationships refuse.
Neither relation proves a full accepted operation, actor, cleanup, admission or runtime.
Missing head remains UNINITIALIZED for archive; journal validates a linked genesis before
REDO_PENDING. Headless staging binding runs only after the complete existing genesis chain/blob/request/
projection audit succeeds. On actual physical-queue absence it uses that audit’s bounded
queue afterimage; present invalid bytes or read errors never select that branch. It does not
excuse missing genesis blobs. A linked next receipt and forks retain existing command-specific checks.

Journal Source.List now requires explicit directory FileInfo from the listing handle plus
bounded child entries. No missing-metadata verification fallback exists, including for empty
roots. Direct Native audits retain and reuse state/intent roots and visited no-follow parents
through each capture/body/recheck; explicit synthetic Sources supply their own directory
observations without a native filesystem claim. Archive retains the analogous local lifetime.
Named directory bindings are compared with retained identities using os.SameFile. Tuples
include roots/parents, names, type/mode/size/mtime and file identities, plus hashes of exact
consumed staging bytes, including partial/temp/empty bytes. Restoring mtime does not hide a
staging edit. A positively observed disappearance/replacement is SNAPSHOT_MOVED; initial
optional-root absence and unrelated I/O failures remain distinct. Handles close on all exits;
close failures are joined and cannot be erased by a moved attempt.

Journal extends its single four-attempt loop. Archive owns a local four-attempt loop with
one-shot `snapshot.Reader{Retries:-1}` calls. Complete observations are recaptured after
validation/body errors: changed tuples retry, stable errors remain errors, and failed
recaptures refuse. Pending/fork header re-probes retain their existing precedence. Archive
resets all attempt output/manifest/file-list/stream state, preserves sticky staging cleanup
and stream-cap errors across retries, verifies its unlinked external staging file before
delivery, and emits no pre-delivery stdout on failure. Existing partial-delivery wording stays.

The narrow J1 empty-D repair accepts present-empty physical **ticket** capture/hash in the
journal inventory; missing files remain distinct. Full Audit returns INTENT_DIVERGED against
canonical C. Request-only lookup still verifies all private/chain/blob/request content and
can return FOUND/ABSENT with intent agreement NOT_OBSERVED. Actual read errors remain fatal;
physical queue/policy and canonical codecs stay strict. Ordinary snapshot/intent-tree ticket
views keep their existing empty-file refusal and gain no staging/accounting claim.

Named finite evidence (TM-V0-002/006/007/008/009/022; AS-07/10/11/35/36):
`TestTMV0002_AS10_StageObservationFiniteShapes`,
`TestTMV0008_AS07_StageObservationTempAndClosedLayout`,
`TestTMV0002_AS10_StageObservationEmptyRaw`,
`TestTMV0002_AS10_StageObservationDelegatesStrictCodec`,
`TestTMV0002_AS10_CompletedStageOperationAndTimestamp`,
`TestTMV0008_AS07_JournalCompletedStageLabels`,
`TestTMV0022_AS07_ArchiveCompletedStageLabels`,
`TestTMV0006_AS35_EmptyTicketRequestLookup`,
`TestTMV0006_AS35_EmptyTicketStrictReadBoundary`,
`TestTMV0008_AS07_JournalStageReadBindings`,
`TestTMV0008_AS36_JournalStageMovementFourAttempts`,
`TestTMV0008_AS36_NativeDirectoryObservationLifetime`,
`TestTMV0008_AS36_JournalStageReadErrorRecheck`,
`TestTMV0008_AS07_JournalReadHandlesCloseOnEveryExit`,
`TestTMV0008_AS07_JournalFileClosureErrorsSurviveMovement`,
`TestTMV0002_AS10_JournalStageScanBounds`,
`TestTMV0009_AS11_JournalStageGenesisAndNext`,
`TestTMV0009_AS11_StageLinkedGenesisBeforeIntentProjection`,
`TestTMV0009_AS11_StageReceiptOnlyGenesisValidation`,
`TestTMV0008_AS07_JournalStageFixtureRestriction`,
`TestTMV0002_AS10_JournalStageEmptyRaw`,
`TestTMV0022_AS07_StageReadStableBindingsAndPurity`,
`TestTMV0022_AS36_StageAndOrphanMovementFourAttempts`,
`TestTMV0022_AS36_StageBodyErrorRecheck`,
`TestTMV0022_AS10_StageScanAndPayloadBounds`,
`TestTMV0022_AS36_StageStreamAndCleanupSticky`,
`TestTMV0022_AS36_ArchiveCoordinatorOneShotAndReprobe`,
`TestTMV0022_AS36_ArchiveStageValidationAndFailedAfterCapture`,
`TestTMV0022_AS07_ArchiveNativeHandlesAndFIFO`,
`TestTMV0022_AS07_ArchiveFileClosureErrorsSurviveMovement`,
`TestTMV0022_AS07_ArchiveStageFixtureRestriction`,
`TestTMV0022_AS10_ArchiveStageEmptyRaw`,
`TestTMV0008_AS07_JournalCompletedStageRequestBlob`,
`TestTMV0022_AS07_ArchiveCompletedStageRequestBlob`,
`TestTMV0022_AS10_StageExcludedFromFileCapBeforeBodies`.
The builder report records exact focused commands and limitations. Writer/issuer/CAS,
capacity qualification, restore, native Linux execution, physical million-file/power-loss
witnesses, runtime and GP remain held or NOT_RUN. Archive remains chain/digest verification,
not historical semantic acceptance. No CLI mutation, real queue or publication is enabled.

### 5.7 J3-03/04 test-owned fixture integration (2026-09-06)

Authored against reviewed J3-02 `bc3b3dd7df3d383f0d546cc9e6ffe7e46a3d0d28` and the frozen
three-file foundation. This changes test-owned mechanisms only, not a wire profile or
production API. Combined parent canonical gate and fresh independent review of **all**
fixture changes, including the foundation, PASS; see [delivery evidence](reviews/2026-09-06-fixture-delivery.md). Exact
focused results and source preservation are in the
[builder report](reviews/2026-09-06-j3-fixture-integration-builder.md).

The disposable sole-writer harness uses one finite publication sequence for empty INIT and
journal-backed ticket histories containing PAUSE, UNPAUSE, KEEP_JOURNAL and ADOPT_FILE.
The ticket seed is validated fixture construction, not an ADD executor. Each modeled opaque
Plan is associated with a copied immutable inventory from that actual observation; publication
compares it before issuing an era or creating staging. Stable preflight, linked-next physical
PRE/POST mixtures and completed cleanup remain separate. Fresh recovery consumes actual
receipt/blob/request/genesis bytes without the original Plan, Model or cached afterimages.
The shared `snapshot.ObserveStage`/`StageObservation.Bind` validates structure and actual
linkage at observation/reopen; retained no-follow inode, mode, mtime, size and digest checks
still mint session-local tokens. Partial/temp observations remain cleanup-only. No observer
result supplies cleanup, execution, actor or administrative authority. Headless queue binding
comes from completely validated actual genesis bytes when the physical queue is absent.
Completed linkage includes operation-kind mapping and the original timestamp.

Native strict Audit and Export/Verify consume completed flows with persistent empty staging
and matching completed-head leftovers. Precommit archives carry actual current payload and
all linked orphan evidence, excluding staging bodies. Valid divergent D stays D; malformed
D is never replaced with C to obtain an archive. Nonempty ticket syntax remains the ordinary
intent loader's responsibility; archive verification does not claim ticket semantics. Empty-D
archive refusal and malformed-D ordinary intent refusal remain intact. Strict Audit preserves INTENT_DIVERGED;
empty-D request-only lookup verifies FOUND/ABSENT and all private indexes with intent
agreement NOT_OBSERVED. Headless linked INIT requires full genesis proof for journal
REDO_PENDING; archive remains UNINITIALIZED without a physical head. Linked-next and fork
refusals emit zero pre-delivery archive bytes. Every native read checks source, lock,
directory, mode and mtime purity.

Finite acceptance mapping (test names carry the indicated TM/AS IDs):

| Deferred group | Implemented named evidence / remaining gate |
|---|---|
| 1 — shared observer/reopen, slots and actual binding | `TestTMV0008_AS27_FixtureNativeStageParity`, `TestTMV0008_AS11_FixtureNativeStageBindings`, `TestTMV0002_AS27_FixtureNativeAssignedEmptyEvidence`; retained `TestTMV0001_AS10_FixtureFoundationSourceLifetime` and new `TestTMV0009_AS11_FixtureFoundationStaleInventory` |
| 2 — actual completed/precommit/pending native reads | `TestTMV0022_AS11_FixtureNativeCompletedHistories`, `TestTMV0009_AS11_FixtureNativeInterruptedStatuses`, `TestTMV0022_AS35_FixtureNativePrecommitIntent`, `TestTMV0007_AS35_FixtureNativeOrdinaryIntentRefusal` |
| 3 — native empty-D index and strict intent boundary | `TestTMV0006_AS35_FixtureNativeEmptyIntentIndex` (FOUND/ABSENT, original outcome, NOT_OBSERVED, corrupt/missing private index, strict INTENT_DIVERGED) |
| 4 — combined history accounting and read boundary | `TestTMV0008_AS10_FixtureNativeAccountingAndSequentialReads`: staging directory/children in D, excluded from F; exact archive payload/inventory; linked orphans; unchanged-head sequential reads including restored-mtime staging bytes |
| 5 — combined integration/evidence closure | This section, ROADMAP, BUILD-LOG and builder report rebase onto final J3-02; protected production/other-test hashes checked. Parent combined canonical gate and fresh independent review PASS; experimental fixture-only delivery, no TCP-02 or production qualification |

Group 4's exhaustive **in-read** movement/error matrices remain the existing final J3-02
unit evidence in §5.6: `JournalStageMovementFourAttempts`, `NativeDirectoryObservationLifetime`,
`JournalStageReadErrorRecheck`, `StageAndOrphanMovementFourAttempts`, `StageBodyErrorRecheck`,
`StageStreamAndCleanupSticky`, `ArchiveCoordinatorOneShotAndReprobe`,
`ArchiveStageValidationAndFailedAfterCapture`, `JournalFileClosureErrorsSurviveMovement`,
`ArchiveFileClosureErrorsSurviveMovement`, `StageScanAndPayloadBounds` and
`StageExcludedFromFileCapBeforeBodies` (full TM/AS names above). Those tests own private hooks
for in-read byte edits/restored mtime, directory replacement/disappearance, exactly four
attempts, stable errors, failed recapture and sticky close/unlink/stream-cap refusal. Combined
fixture tests exercise representative native reads, not duplicate private-hook matrices.

The 17 foundation suites retain finite publication/fault/prefix coverage, exact final
inventory/cost equality, and componentwise conservative transient domination retaining
unsynced name/shrink debt. Returned faults execute Go cleanup; no crash child was used.
Process-interruption and physical power-loss evidence are NOT_PRODUCED. Twelve-receipt,
128-entry-per-directory, 4-MiB fixture bounds are not general journal qualification. Callable
writer qualification, authenticated issuer/admin binding, hostile-editor CAS, live runtime/escrow,
real queues, restore and GP remain held. Decisions 0003/0004 subsequently released the callable
fixture writer and ticket CLI within their narrower local-operator premise. Native Linux execution and performance are NOT_RUN. TCP-02
remains incomplete; read support and fixture transactions grant no durable-writer authority.

## 6. Transition and recovery table

### 6.1 Enumerations

Attempt phase: `ADMITTED | RUNNING | BUILT | CHECKING | REVIEWING | REPAIRING | STOPPING |
QUARANTINED | BLOCKED_RECOVERY | FAILED | CANCELLED | READY_FOR_INTEGRATION | COMPLETED`.
Terminal: `FAILED`, `CANCELLED`, `COMPLETED`. Gate state: `PENDING | RUNNING | PASSED | FAILED |
BLOCKED | NOT_RUN | STALE`. Mapping from ATM-V0-005 names: `claimed→ADMITTED`,
`built→BUILT`, `build_failed→FAILED/BUILD`, `gated→CHECKING`, `gate_failed→FAILED/GATE`,
`reviewed→REVIEWING`, `repaired→REPAIRING`, `accepting→` the §7.3 reducer transaction,
`accepted→READY_FOR_INTEGRATION|COMPLETED`, `rejected→FAILED/REVIEW_REJECTED`,
`blocked→BLOCKED_RECOVERY` or ticket `HELD`, `TAKEN→REVISION_CONFLICT/FENCED`.

### 6.2 Transitions

| From | Event | Guard | To | Reservation / worktree |
|---|---|---|---|---|
| eligible ticket | `admit` | §4.1 | `ADMITTED` | reserved; no worktree yet |
| `ADMITTED` | worktree add `DONE`, spawn bootstrap record read (§6.4 step 4) | supervisor matches `attempt.supervisor`; effects recorded; `.boot` identity matches key and `.boot.pid` equals the pid the supervisor forked; phase still `ADMITTED` (not `STOPPING`) | `RUNNING` | worktree `<attemptId>-<generation>` created; `lane` recorded from `.boot`; `.ack` linked after commit |
| `ADMITTED` | spawn `FAILED`, or bootstrap fenced/timed out (§6.4 `BOOT_FENCED`/`BOOT_TIMEOUT`/`NOEXEC`) | §6.4 no-exec proof; recorded by the owning supervisor, or by `reconcile`/`drain`/`cancel` after that supervisor is dead | `FAILED/SPAWN` | wrapper quiescence `PROVED` → released; `FENCED` (no wrapper pid known) → released as fenced; otherwise retained; not charged to admissions per revision (§1) |
| `ADMITTED` (pending `PROCESS_SPAWN`, supervisor live) | `cancel` / `drain` | generation matches; effect files untouched | `STOPPING` | retained; the owning supervisor observes `STOPPING` at §6.4 step 4, fences `.ack`, and records `EFFECT_OUTCOME` `FAILED` + `CANCELLED` |
| `RUNNING` | runner exit 0, lane dead | liveness dead; `.ack` is `ack:true` | `BUILT` | retained |
| `RUNNING` | exit≠0 / crash / cap hit | lane dead; `.ack` is `ack:true` | `FAILED/BUILD` | retained until `PROVED`, then released |
| `RUNNING` | leader exited without exec (`.ack` is `noexec` or `fenced`, or absent with the leader dead by §6.3) | recorded by the owning supervisor, or by `reconcile` after it is dead | `FAILED/SPAWN` | as the `ADMITTED` `FAILED/SPAWN` row |
| `RUNNING` | exit, survivors live | liveness live | `BLOCKED_RECOVERY/SURVIVORS` | retained |
| `RUNNING` | out-of-scope write detected | dirty outside closure | `QUARANTINED` | retained; never integrates |
| `RUNNING`/`REPAIRING`/`CHECKING`/`REVIEWING` | `cancel` | generation matches | `STOPPING` | retained |
| `STOPPING` | SIGTERM, then SIGKILL after 10 s; lane dead | identity revalidated before signal | `CANCELLED` | released after `PROVED` |
| `STOPPING` | survivors after SIGKILL or `SIGNAL_REFUSED_IDENTITY` | | `BLOCKED_RECOVERY/SURVIVORS` | retained |
| any non-terminal | supervisor dead, lane live | `attempt.supervisor` dead by §6.3 (ATM-V0-005a precedence); recorded only by `reconcile`/`drain`/`cancel` | `BLOCKED_RECOVERY/SUPERVISOR_LOST` | retained; next `reconcile` signals by pgid |
| any non-terminal | pending effect without outcome | `attempt.supervisor` dead by §6.3; actor is `reconcile`, `drain` or `cancel`. With a live supervisor no command takes this edge (`EFFECT_OWNED`) | `BLOCKED_RECOVERY/UNCERTAIN_EFFECT` | retained |
| `BUILT` | gates start | candidate tree recorded | `CHECKING` | retained |
| `CHECKING` | all required pre-review gates `PASSED` | | `REVIEWING` (or reducer when review lane not required) | retained |
| `CHECKING` | a required gate `FAILED`/`BLOCKED`/timeout | | `FAILED/GATE` | retained until `PROVED` |
| `CHECKING` | gate `STALE` | rerun budget left | rerun once, else `FAILED/GATE_STALE` | retained |
| `REVIEWING` | verdict `ACCEPT` and §7.2 pass | | reducer → `READY_FOR_INTEGRATION`/`COMPLETED`/`FAILED/MANIFEST` | retained until terminal or integration |
| `REVIEWING` | `REPAIR` with confirmed findings, rounds left | | `REPAIRING` | retained |
| `REVIEWING` | `REJECT`, rounds exhausted, or non-converging | | `FAILED/REVIEW_REJECTED`; ticket hold `ESCALATED` placed | retained until `PROVED` |
| `REVIEWING` | malformed/unparsed round | retry left | rerun with fresh reviewer, else `BLOCKED_RECOVERY/ADJUDICATION` | retained |
| `REVIEWING` | independence `UNVERIFIED`/`CONTRADICTED` or `UNVERIFIABLE` claim | | `BLOCKED_RECOVERY/ADJUDICATION` | retained |
| `REPAIRING` | repair exit 0, lane dead | | `BUILT` (prior gate results `STALE`) | retained |
| `REPAIRING` | exit≠0 / cap | | `FAILED/REPAIR` | retained until `PROVED` |
| `READY_FOR_INTEGRATION` | `integrate` with `INTEGRATE` grant and commit containing candidate tree | | `COMPLETED` | released after `PROVED`; worktree removable |
| `BLOCKED_RECOVERY/*` | `reconcile` resolves cause | actor + reason recorded; quiescence `PROVED` (or `FENCED` for adapter-only causes and for a spawn whose bootstrap records are both fences, §6.4); no `PENDING` effect; old generation appended to `priorGenerations` | `ADMITTED`, `BUILT` or `READY_FOR_INTEGRATION` (phases with no lane process), else `FAILED`/`CANCELLED`; never `RUNNING`, `REPAIRING`, `CHECKING` or `REVIEWING` | retained |
| `BLOCKED_RECOVERY/*`, `FAILED`, `CANCELLED` | `retry` | ticket OPEN, not held, `retryCount<3`, `acceptanceRevision` matches, every prior generation `PROVED`/`FENCED`, no `PENDING` effect, prior worktree removed (`WORKTREE_REMOVE` `DONE`) or never created | new generation; `ADMITTED`; retry re-runs §4.1 in full and replaces the reservation entry atomically | new worktree `<attemptId>-<generation>` |
| `FAILED`, `CANCELLED`, `BLOCKED_RECOVERY/*` | `resume` | as `retry`, and the target is a phase with no lane process (`BUILT` with the candidate tree still reachable, or `READY_FOR_INTEGRATION`); otherwise `UNSUPPORTED` (use `retry`) | new generation; that phase | retained; gate results `STALE` |
| `QUARANTINED` | `cancel` or `scope-expand` | expand = stop lane, prove `PROVED`, new reservation, new generation | `CANCELLED` / `ADMITTED` (respawn under the new generation) | expand fails all-or-nothing |
| any | `pause` | writes `ADMISSION` barrier | unchanged; only `admit`, `retry`, `resume`, `scope-expand`, `import apply`, `AUTHORITY_SWITCH` refused `PAUSED`; everything else proceeds | unchanged |
| any | ticket `HOLD` | | unchanged; completion transition refused `TICKET_HELD` until release | unchanged |
| any non-terminal | `drain` | writes `ALL` barrier, exclusive lock, grace, cancel sweep | each lane as `cancel`; unproved lanes stay `BLOCKED_RECOVERY` | never deleted while unresolved |

### 6.3 Liveness, quiescence and release

Liveness is ATM-V0-005a verbatim: a lane is live when a live process carries the recorded pgid
and either the leader is absent or its start time matches the recorded one; start time comes
from `kern.proc.pid` on Darwin and `/proc/<pid>/stat` plus `btime` on Linux; before any
`kill(-pgid)` the leader start time is re-read and a mismatch records `SIGNAL_REFUSED_IDENTITY`.
Supervisor liveness uses the same test on the supervisor pid. Quiescence is `PROVED` when the
lane is dead by that test, no effect is `PENDING`, and the shared checkout shows no writes
outside the closure. Only `PROVED`, or `FENCED` exactly as TM-V0-016 (adapter effects all
declared refused) and §6.4 (both bootstrap records are fences, so no leader identity ever
existed and no runtime executed) define it, releases a reservation. Worktree deletion additionally
requires a terminal phase, or `COMPLETED`, and that the candidate tree is reachable from the
attempt branch. The residual pid-reuse race named in ATM-V0-005a remains a non-goal. Process
groups are a cleanup handle for these signals only; they qualify no capability axis (TM-V0-018).

### 6.4 Spawn bootstrap and acknowledgement (`PROCESS_SPAWN`)

The lane runtime is never exec'd directly. Two fenced records exist per `key`, `effects/<key>.boot`
and `effects/<key>.ack`; each is created at most once, by exactly one writer, with the §5.2
link-in primitive (write `<name>.tmp-<pid>`, fsync, `link(2)` to the final name, unlink temp,
fsync dir), so a record is either absent or complete and canonical, never partial. `EEXIST`
means another party won and the loser reads the winner's record. A `link(2)` failure other than
`EEXIST` (R3; for example `EIO`, `ENOSPC`, `EACCES`) is not a race outcome and never authorizes
anything: the writer removes its temp file, treats the record as unwritten, and a leader that
meets one exits 3 without exec; the owning supervisor then records `FAILED/SPAWN` by the step 4
rule. There is no other fence authority: no unlink, no rewrite, no rename over either file.
Record contents are closed:

| Record | Writer | Content |
|---|---|---|
| `.boot` | leader (step 2) | `{key, attemptId, generation, pid, startTime, pgid}` |
| `.boot` | recovery fence (§3.4 row, after the supervisor is proved dead; or the owning supervisor on `BOOT_TIMEOUT`) | `{key, fenced:true}` |
| `.ack` | owning supervisor (step 4) | `{key, attemptId, generation, pid, startTime, ack:true}` |
| `.ack` | leader on its no-exec path (step 3) | `{key, attemptId, generation, pid, startTime, noexec:true}` |
| `.ack` | recovery fence, or the owning supervisor on cancel/timeout | `{key, fenced:true}` |

Steps, with `key` the effect key of §3.4:

1. `EFFECT_INTENT` transaction records the `PROCESS_SPAWN` effect with `args` = `{attemptId,
   generation, runtimeId, worktreePath, argvSha256}` and, if not already set for this
   generation, `attempt.supervisor` = the committing process's `{pid, startTime}` (§3.4). The
   supervisor identity is durable before the fork in step 2.
2. The supervisor forks and spawns `corvint-tasks lane-leader --key <key> --attempt <attemptId>
   --generation <generation> --state <state dir>` (a `corvint-tasks` process, not the runtime) and
   remembers the forked pid. The leader calls `setsid`, then links in `.boot` with its own
   `{key, attemptId, generation, pid, startTime, pgid}`. On `EEXIST` the leader exits 3
   immediately; it has executed nothing.
3. The leader polls for `.ack` for at most the leader ack wait (90 s, §1). A parseable `.ack`
   whose `{key, attemptId, generation, pid, startTime}` equal its own and carries `ack:true`
   authorizes step 5. Any other parseable content (`fenced`, or an identity mismatch) ends the
   wait: the leader exits 3 without exec. At the deadline with no `.ack`, the leader links in
   `.ack` with its identity and `noexec:true` (on `EEXIST` it re-reads the winner and applies
   this same rule once more) and exits 3 without exec. Because the leader is single-threaded
   and execs only after reading a matching `ack:true` record, exactly one of "runtime exec'd"
   and "a `noexec`/`fenced` `.ack` exists or `.ack` is absent" holds for every key.
4. The supervisor waits for `.boot` at most the supervisor boot wait (90 s, §1) or until the
   forked child exits, whichever is first. If the child exited or the deadline passed with no
   `.boot`, it links in a fence `.boot` (`{key, fenced:true}`); on `EEXIST` it proceeds with the
   leader's record. Then, under lock, it: checks the transaction carries its own
   `supervisor` identity (§3.4) and the attempt is still `ADMITTED` under the same generation
   (a `STOPPING` phase from `cancel`/`drain` is honoured: it fences `.ack`, records
   `EFFECT_OUTCOME` `FAILED` and the `CANCELLED` transition, and stops); reads `.boot`; if it is
   a fence, or its `pid` differs from the forked pid, or its `pid`/`startTime` are dead by §6.3,
   it records `EFFECT_OUTCOME` `FAILED` and `FAILED/SPAWN` with cause `BOOT_FENCED` or
   `BOOT_TIMEOUT` (no-exec proof: no `ack:true` was ever written for this key); otherwise it
   commits `EFFECT_OUTCOME` `DONE` + `TRANSITION` to `RUNNING` with `lane` taken from `.boot`.
   Grants, holds and the barrier are rechecked in that transaction. Only after that receipt is
   linked in does it re-check the leader's liveness by §6.3 and link in `.ack` with `ack:true`.
   If the leader is dead, or the link fails `EEXIST` and the existing `.ack` is `noexec`, it
   records in a further transaction `RUNNING → FAILED/SPAWN` with cause `NOEXEC` (§6.2 row).
5. The leader, on a matching `ack:true`, execs the runtime with the pinned argv and environment.

Two separate proofs are recorded for every `FAILED/SPAWN`:

- **No runtime executed** (`noExec:"PROVED"`, required for `FAILED/SPAWN`): `.boot` is a fence;
  or `.ack` is `noexec` or `fenced`; or `.ack` is absent and the leader named by `.boot` is dead
  by §6.3. It is never inferred from exit codes, time elapsed, or an unparsable record.
- **Wrapper quiescence** (the attempt's `quiescence`): `PROVED` when the leader pid/startTime
  from `.boot`, or the forked pid reaped by the owning supervisor, is dead by §6.3; `FENCED`
  when no leader identity was ever recorded (`.boot` is a fence) and `.ack` is fenced, so a
  late wrapper can only read the fences and exit; `UNPROVED` otherwise, and then the
  reservation is retained: the recovering command fences whichever of `.boot`/`.ack` is absent,
  waits at most the wrapper exit wait (10 s, §1) for the wrapper to exit, and re-tests §6.3;
  if still live it signals the wrapper's pgid by the §6.3 protocol (SIGTERM, SIGKILL after 10 s),
  and if survivors remain records `BLOCKED_RECOVERY/SURVIVORS`.

A `.boot` or `.ack` that is present but unparsable cannot arise from this protocol (link-in is
all-or-nothing) and is treated as corruption: no-exec is `PROVED` only via a fence or `noexec`
record already present or via an absent `.ack`; the recovering command fences the other record,
and quiescence stays `UNPROVED` until the leader pid is known dead (from the forked pid, when
the owning supervisor recovers) or an operator `reconcile` records `SURVIVORS` handling.

Consequences: a runtime executes only when its `(attemptId, generation, pid, startTime)` is
durable in the journal and authorized in that same transaction; a crash at any point before
step 5 leaves at most a `corvint-tasks` leader that exits on its own within the leader ack wait;
ordinary lock contention (up to the 30 s lock wait plus redo and fsync work) fits inside the
90 s leader wait, and a leader that nevertheless times out produces `FAILED/SPAWN` with
`NOEXEC`, never `FAILED/BUILD`, and does not consume an admission (§1). Recovery of a pending
spawn is reserved to `reconcile`/`drain`/`cancel` after `attempt.supervisor` is dead (§3.4,
TM-V0-017); an unrelated transaction never fences a live supervisor's key, and `cancel`/`drain`
against a live supervisor request `STOPPING` and let that supervisor fence and record. `.boot`
and `.ack` are removed only when the effect reaches `RESOLVED` and the generation is `PROVED` or
`FENCED`. An uncertain spawn keeps its reservation. None of this is a permission boundary; it
is identity and ordering only.

## 7. Gate, review and completion schemas and acceptance tables

### 7.1 `taskman-gate-result/0`

`profile, gateId, attemptId, generation, ticketRevision:Count, definitionSha256,
candidateTreeOid, executedTreeOid:OID|null, executedCwd:"WORKTREE"|"CANDIDATE"|null,
porcelainClean:boolean|null, inputsSha256, environmentSha256, policySha256, configSha256,
startedAt, endedAt, state, outcomeClass:"EXIT"|"TIMEOUT"|"SIGNAL"|"SPAWN_FAILED"|"OUTPUT_LIMIT"|
"UNKNOWN", exitCode:Count|null, signal:Identifier|null, evidence:[{label, sha256, bytes:Size}],
reusedFrom:Digest|null`. `executedTreeOid` is `HEAD^{tree}` of the directory the command ran in,
read immediately before spawn together with the porcelain status; `PASSED` requires
`outcomeClass:"EXIT"`, the expected predicate satisfied, every required evidence label present,
`porcelainClean:true` and `executedTreeOid == candidateTreeOid`. A result that fails the last two
is `BLOCKED/DIRTY_WORKTREE` and is never reusable.

### 7.2 `taskman-claim-disposition/0` and review truth table

`profile, reviewId, attemptId, generation, ticketRevision:Count, policySha256, configSha256,
candidateTreeOid:OID, reviewWorktreeTreeBefore:OID, reviewWorktreeTreeAfter:OID, round:Count,
reviewer:{runtimeId, model:label, effort:label, catalogId:label|null, charterSha256,
conversationId:Identifier|null},
independence:"VERIFIED"|"UNVERIFIED"|"CONTRADICTED", claims:[{claimId, class:
"ACCEPTANCE_CRITERION"|"REQUIREMENT"|"DOC_CONSISTENCY"|"SCOPE"|"GATES"|"SAFETY",
disposition:"ACCEPT"|"REJECT"|"UNVERIFIABLE", findingIds:[label]}], findings:[{findingId,
severity:"HIGH"|"MED"|"LOW", path:Path, range:Identifier, claim:prose, bytesAnchor:Digest,
textAnchor:Digest, boundaryCrossed:boolean, verification:"CONFIRMED"|"REFUTED"|"UNVERIFIED"}],
verdict:"ACCEPT"|"REPAIR"|"REJECT", reportSha256:Digest`. `reportSha256` names the verbatim
reviewer report bytes in `evidence/`, present before the `REVIEW` receipt; the record is
derived from them by the tool's deterministic parser (TM-V0-020). `reviewWorktreeTreeBefore`
and `After` must equal `candidateTreeOid`, else the round is `CONTAMINATED` (TM-V0-018).

Required claim set = one `ACCEPTANCE_CRITERION` per criterion index, one `REQUIREMENT` per
`requirementRefs` entry, plus `DOC_CONSISTENCY`, `SCOPE` and `GATES`, derived from the canonical
ticket record at the review's `ticketRevision`. Rows apply in order:

| # | Condition | Result |
|---|---|---|
| 1 | report unparsable, a required claim missing, duplicated, or with an unknown class/ID | round `MALFORMED`: no repair round consumed; retry once fresh; then `BLOCKED_RECOVERY/ADJUDICATION` |
| 2 | verdict `REPAIR`/`REJECT` with zero findings | `MALFORMED` as row 1 |
| 3 | independence `CONTRADICTED` | round void; rerun fresh; counts against malformed retry |
| 4 | any claim `UNVERIFIABLE` | verdict recorded; completion blocked `ADJUDICATION` |
| 5 | verdict `ACCEPT` but any claim `REJECT` or any `CONFIRMED` HIGH/MED finding open | treated as `REPAIR` (obligation never reduced) |
| 6 | verdict `ACCEPT`, all claims `ACCEPT`, no open `CONFIRMED` finding, independence `VERIFIED` | review passes |
| 7 | as 6 but independence `UNVERIFIED` | review recorded; completion blocked until `ADJUDICATE` grant |
| 8 | verdict `REPAIR`, ≥1 `CONFIRMED` finding, `repairRound < repairRounds` | `REPAIRING` |
| 9 | verdict `REPAIR` whose 3-tuples all existed in the prior round | non-converging: escalate `FAILED/REVIEW_REJECTED` + `ESCALATED` hold |
| 10 | verdict `REJECT` or rounds exhausted | `FAILED/REVIEW_REJECTED` + `ESCALATED` hold |
| 11 | zero required claims | valid only when `kind ∈ allowEmptyObligationsKinds`; otherwise `MALFORMED` |

Finding verification follows ATM-V0-013 (resolve, quote, machine byte check); `UNVERIFIED`
findings never drive repair and never count as resolved.

### 7.3 `taskman-completion-manifest/0` and reducer

`profile, attemptId, generation, ticketId, ticketRevision (acceptanceRevision),
ticketRecordSha256, policySha256, configSha256, baseCommit, candidateTreeOid,
acceptanceMap:[{criterionIndex:Count, evidence:[Digest]}], gateResults:[Digest],
reviews:[Digest], docsResult:Digest|null, cem:Digest|null, ocm:Digest|null,
unresolvedFindings:[label], budget:BudgetUsage, budgetComplete:boolean, scopeCheck, mode`.

The manifest is an input the reducer re-derives, not a claim it trusts: every obligation set is
recomputed from the pinned policy and the canonical ticket record, every digest is loaded from
`evidence/` and validated against its profile, and caller-supplied lists (including an empty
`unresolvedFindings`) never shorten an obligation. Reducer rows apply in order; the first failure
is the result code and no status changes:

| # | Check | Failure code |
|---|---|---|
| 1 | `ticketRevision` equals the current `acceptanceRevision` and `ticketRecordSha256` equals the digest of the canonical record at the revision that last changed it; publication not `DIVERGED` | `STALE_TICKET` / `INTENT_DIVERGED` |
| 2 | ticket `OPEN` (not held, archived, completed) | `TICKET_STATE` |
| 3 | `policySha256` and `configSha256` equal the attempt's pinned values and the current policy | `STALE_POLICY` |
| 4 | `candidateTreeOid` equals the attempt record and the lane worktree `HEAD^{tree}`, with empty porcelain status | `STALE_TREE` / `DIRTY_WORKTREE` |
| 5 | `scopeCheck == WITHIN` | `OUT_OF_SCOPE` |
| 6 | required gate set = policy `required:true` gates ∪ ticket `requiredGates`; for each, a result in `evidence/` with profile `taskman-gate-result/0`, state `PASSED`, not `STALE`, `attemptId`/`generation`/`ticketRevision`/`policySha256`/`configSha256` equal to the manifest, `definitionSha256` equal to the current policy's definition, and `executedTreeOid == candidateTreeOid` | `MISSING_GATE` / `STALE_GATE` / `GATE_FAILED` |
| 7 | every acceptance criterion index of the canonical record maps to ≥1 digest that resolves to an evidence blob referenced by a gate result or review of this attempt/generation | `MISSING_EVIDENCE` |
| 8 | the latest `REVIEW` receipt for this `attemptId`/`generation` resolves to a claim record with `candidateTreeOid`, `ticketRevision`, `policySha256`, `configSha256` equal to the manifest; its required claim set recomputed from the canonical record equals its claims; it passes §7.2 row 6, or row 7 with an `ADJUDICATE` grant at this `acceptanceRevision`; the set of open `CONFIRMED` findings across the attempt's reviews, recomputed, is empty | `REVIEW_INCOMPLETE` / `UNRESOLVED_FINDING` / `INDEPENDENCE_UNVERIFIED` |
| 9 | docs lane result (`taskman-gate-result/0`, kind `DOCS`) `PASSED` and bound as row 6 when `docsLane.required` | `DOCS_MISSING` |
| 10 | when `cemRequired`/`ocmRequired`: `cem`/`ocm` resolve to evidence blobs whose recorded profile is the Corvint CEM / OCM record, whose tree/commit binding equals `candidateTreeOid`/`baseCommit`, and whose recorded status is complete (`cem-status:ready`; OCM aggregate closed with `outcome` recorded); an unknown profile, a mismatched tree or a non-ready status fails | `CEM_MISSING` / `OCM_MISSING` |
| 11 | `budgetComplete` and no required field over cap, recomputed from the attempt's recorded usage | `BUDGET_UNKNOWN` / `BUDGET_EXCEEDED` |
| 12 | `APPROVAL_REQUIRED` tickets hold a non-revoked `COMPLETE` grant at this `acceptanceRevision` | `APPROVAL_MISSING` |
| 13 | `mode == QUALIFIED` | `DEVELOPMENT_MODE` (recorded `DEVELOPMENT_ACCEPT`; ticket unchanged) |

Named negative fixtures (AS-25): an evidence digest that exists but is a different profile; a
gate result for the same gate from a different attempt or generation; a gate result whose
`executedTreeOid` differs from `candidateTreeOid`; a claim record whose claims omit one criterion
of the current record; a manifest with empty `unresolvedFindings` while a `CONFIRMED` HIGH
finding is open; a `cem` digest naming an unrelated blob; a `cem` blob with `cem-status:
not-ready`. Each must fail on the row named above.

Pass: `kind ∈ integrationRequiredKinds` → `READY_FOR_INTEGRATION`; otherwise `COMPLETED` with
`completion.kind:"VERIFIED"`. Content identity proves identity only, never sufficiency.

### 7.4 Execution permission

A non-fixture queue admits nothing until `queue.executionCutover` names the owner's decision and
a `QUALIFICATION` receipt lists passing scenarios for G2 and G3 (§9.2). A fixture queue
(`fixture:true`) is the only place `DEVELOPMENT` attempts run, and they never write outside the
fixture repository. Restore clears `executionCutover`.

## 8. Amendments to accepted ATM revision 6 (IDs preserved, text retained upstream)

| # | ATM clause | Amendment frozen here |
|---|---|---|
| A1 | ATM-V0-001 / S1 | Native ticket ownership, CRUD and lifecycle per §3; plain-JSON queue adapter becomes the native store; `ROADMAP.md` adapter is read-only until §5.4 cutover; Beamfall writer boundary retained. |
| A2 | ATM-V0-002/002a, 007 | Stable `queueId` distinct from WQO `queueSourceId`; native admission uses `taskman-plan/0` with whole-plan freshness (TM-V0-015: any drift refuses the plan, preserving ATCP-V0-006 unnarrowed), never a WQO wave; reservations span running, review, repair and blocked attempts (§4.2). ATM-V0-002a stays for WQO-driven foreign waves only (TCP-09). |
| A3 | ATM-V0-005a/005c/019a/028, §8 | Attempt generation fencing (TM-V0-011); receipt is the commit point (§5.2); lock is exclusive per transaction rather than shared (stricter, same bounded drain wait); stale age never frees a lane, only the liveness test does; drain never deletes with unresolved attempts; archive contents per TM-V0-022. |
| A4 | ATM-V0-015/021 | Capability profiles (TM-V0-018) separate enforcement from prompt hygiene; process groups and quoted prose are never called a sandbox. |
| A5 | wave policy | ATCP-V0-009 priority-first profile accepted for native queues as `taskman-priority-first/0`; WQO `MAXIMUM`/`GREEDY` untouched. |
| A6 | ATM-V0-012/013/013a/014 | Closed per-claim dispositions, §7.2 truth table, §7.3 manifest; routing thresholds (ATM §7, §9) retained unchanged; aggregate budgets per §1. ATM-V0-008a's "NOT_OBSERVED never blocks a review" stands; completion (not review) is what waits for adjudication. ATM-V0-013a report bytes retained verbatim in `evidence/` (TM-V0-020). |
| A7 | ATM-V0-025, 019 | Configuration pin unit is the attempt: the resolved configuration document is written to `pinned/` at admission (`configSha256`) and is the only configuration that attempt's lanes read; `reconfigure` is an `OPERATOR` `CONFIG_PIN` receipt that applies to attempts admitted afterwards, never to a live one, and keeps the verifier tuple fixed within a review round. Routine ticket edits (§3.2) do not re-pin. ATM-V0-019 replay is `receipt replay` (TM-V0-008) over plan, review aggregation and reducer inputs; wave replay stays out of scope. |

## 9. Slices, ownership, scenarios and rollback

### 9.1 Package ownership (proposed paths; a worker claims exact files before starting)

| Ticket | Owns | Depends on |
|---|---|---|
| TCP-01 | `internal/wire` (canonical JSON, `Count`/`Size`, IDs, bounds, `taskman-command-result/0`), `internal/ticket` (records, validation, mutation computation, eligibility, `acceptanceRevision`), `internal/intent` (intent files, primary-worktree resolution, publication status, divergence detection), `cmd/corvint-tasks` read and ticket verbs, `archive export|verify` | TCP-00 only. Done without journal wiring: mutation validation and post-record computation are tested against an in-memory journal fake; the fake never counts as commit evidence |
| TCP-02 | `internal/authority` (common-dir resolution, filesystem qualification, lock), `internal/journal` (receipts, head, redo, chain, request index, barrier), `internal/reservation`, `internal/proc` (liveness, signals, `lane-leader` bootstrap §6.4), `init|admit|cancel|retry|resume|reconcile|pause|unpause|drain`, `archive restore` | TCP-00 |
| TCP-02b (serial, after TCP-01 and TCP-02) | wiring `ticket` mutations to `journal` commit; AS-02/AS-03/AS-05 rerun against the real journal | TCP-01, TCP-02 |
| TCP-03 | `internal/plan` and the Corvint-side native queue adapter/planner profile; frozen `taskman-perf-baseline/0` before the first Corvint edit; GP | TCP-01, TCP-02 |
| TCP-04 | `internal/runtime` (registry, capability probes, supervisor, budgets, ack writer) | TCP-01, TCP-02 |
| TCP-05 | `internal/gate`, `internal/review`, `internal/manifest` | TCP-04 |
| TCP-06 | `internal/importer`, cutover commands, roadmap projection; GP rerun for its Corvint-side changes | TCP-01, TCP-03, TCP-05 |

TCP-01 and TCP-02 proceed in parallel on disjoint packages. TCP-01 is complete without the
journal; the commit wiring is the separate serial step TCP-02b, checked off only when the real
journal exists. Nothing in TCP-00's G0 requires code from TCP-01 or TCP-02.

### 9.2 Acceptance scenarios (each becomes a named Go test; all NOT_RUN)

| ID | Scenario | Gate | Ticket |
|---|---|---|---|
| AS-01 | wire goldens: min/max, unknown/duplicate key, framing, hostile code points, every profile version, `Count` vs `Size` | G1 | 01 |
| AS-02 | create/refine; stale `expectedRevision` fails with byte-identical store | G1 | 01 |
| AS-03 | identical retry replays; same requestId different bytes conflicts | G1 | 01 |
| AS-04 | cycle (direct, transitive, via gate), missing dependency, case-fold duplicate, invalid priority, unsupported version all refused | G1 | 01 |
| AS-05 | hold/release, archive tombstone retained and restored, reopen keeps prior completion (in the chained prior record; the reopened record carries `completion:null`, §3.1) | G1 | 01 |
| AS-06 | manual completion labelled `MANUAL`, no manifest, never reported as verified | G1/G3 | 01 |
| AS-07 | every read on uninitialized, pending-redo and forked stores leaves state dir and intent tree byte-identical | G1 | 01/02 |
| AS-08 | pagination pins seq/tree digest and reports truncation | G1 | 01 |
| AS-09 | export → verify → restore into empty dir → byte-identical records; restore clears execution permission | G1/G4 | 01/02 |
| AS-10 | every §1 limit at boundary and boundary+1; receipt with 8 inline and 1 blob post entry; reservation set at 64 × 4,096 | G1 | 01/02 |
| AS-11 | crash matrix C1..C8b, C10..C13 by fault injection at each point; C6 with a Git checkout of the ticket file to its pre digest after commit, a manual edit, and a linked-worktree read; redo never rewrites a third-value file | G2 | 02 |
| AS-12 | lost receipt, reordered receipt, head ahead → `JOURNAL_FORKED`; audit detects | G2 | 02 |
| AS-13 | two processes from two worktrees race one ticket; exactly one admits; unsupported filesystem refuses | G2 | 02 |
| AS-14 | callback with old generation rejected `FENCED` | G2 | 02 |
| AS-15 | capacity fits but resource collides → nothing reserved; prefix/file collision normalization | G2 | 02 |
| AS-16 | supervisor death with live lane → `BLOCKED_RECOVERY`, reservation kept; survivors after SIGKILL never released; §6.4 windows: supervisor killed before `.boot`, after `.boot` before `RUNNING`, after `RUNNING` before `.ack`; in every window at most zero runtimes execute and a second spawn is refused until `PROVED`/`FENCED` | G2 | 02/04 |
| AS-17 | pending effect at restart reconciled by the §3.4 class table per kind; `PROCESS_SPAWN` never retried while unfenced; NONE refused at admission; retry with live prior lane refused; retry replaces the reservation entry atomically | G2 | 02 |
| AS-18 | interrupted drain resumes; drain never deletes with unresolved attempt; pause (`ADMISSION`) lets a running lane transition, cancel and reconcile, and cutover A2 completes | G2 | 02 |
| AS-19 | every §6.2 row and every unenumerated transition rejected; `resume` into `RUNNING`/`REPAIRING`/`CHECKING`/`REVIEWING` refused; routine mutation on a ticket with a live attempt keeps the attempt completable; acceptance-relevant mutation refused `ATTEMPT_LIVE` | G2 | 02/04 |
| AS-20 | `NOT_OBSERVED` never zero; aggregate headroom refuses admission; failed/cancelled lanes accounted | G2/G3 | 04 |
| AS-21 | any plan entry drift (revision, resources, intent tree, reservation set) refuses the whole plan and admits nothing; urgent runnable never displaced; `deferredSinceSeq` reported from a recorded preview; WQO conformance unchanged | G5 | 03 |
| AS-22 | gate reuse on identical tuple, `STALE` on each drift axis, worker edit of gate refused; gate in a dirty lane worktree → `DIRTY_WORKTREE`; `CANDIDATE` cwd runs at exactly `candidateTreeOid` | G3 | 05 |
| AS-23 | §7.2 every row incl. seeded false ACCEPT, empty ACCEPT, malformed, non-converging; report bytes retained and re-parsed identically by `receipt replay` | G3 | 05 |
| AS-24 | reviewer worktree mutated → `CONTAMINATED`; reviewer runtime with `OBSERVED` write axis limited to `DEVELOPMENT`; probe evidence required for `ENFORCED`; state dir write refused by the enforced axis, not by the dirty check | G3 | 05 |
| AS-25 | §7.3 every row incl. stale tree/policy, missing docs/CEM/OCM, development mode, and every named negative fixture in §7.3 | G3 | 05 |
| AS-26 | import dry-run inventories every item; drift refuses apply; repeat apply creates no duplicate; interrupted batch apply resumes; shadow records ineligible until A5 | G4 | 06 |
| AS-27 | §5.4 interruption at every step resumes; revert keeps post-cutover history; restore reinstates `requests/`, barrier and generation and parks non-terminal attempts in `BLOCKED_RECOVERY/RESTORED` | G4 | 06 |
| AS-28 | non-fixture queue refuses admission without cutover + qualification receipt, or with `DIVERGED` publication | G2/G4 | 02 |
| AS-29 | filesystem qualification on each allowed type and one refused type; primary worktree mismatch refused | G2 | 02 |
| AS-30 | two disjoint admissions, competing coordinator, active collision, scope expansion, cancelled worker with complete receipts | G5 | 07 |
| AS-31 | performance: frozen `taskman-perf-baseline/0`, interleaved old-Go/new-Go run over every preregistered row under §9.4; failed and noisy runs kept; output identity per row | GP | 03/06 |
| AS-32 | contention: existing Corvint command rows re-measured with `corvint-tasks` `STOPPED`, `IDLE_INITIALIZED` and `MAX_ADMITTED` (concurrent lanes, gates, worktree adds); no background process, poller or daemon exists after any `corvint-tasks` command exits | GP | 03/04/07 |
| AS-33 | R2 F1 pending-effect ownership: the named negative fixtures below (N1a..N1e) | G2 | 02/04 |
| AS-34 | R2 F2 bootstrap windows: the named negative fixtures below (N2a..N2f) | G2 | 02/04 |
| AS-35 | R2 F3 `ADOPT_FILE` composition: the named negative fixtures below (N3a..N3f) | G1 | 01 |
| AS-36 | R2 F4 receipt immutability, chain bounds and export snapshot: the named negative fixtures below (N4a..N4f) | G1/G2 | 01/02 |
| AS-37 | R2 LOW cutover exemption: N5a stand-alone `AUTHORITY_SWITCH` under the A1 `ADMISSION` barrier refused `PAUSED`; N5b the same receipt issued by `cutover` step A5 commits; N5c `cutover` under an `ADMISSION` barrier with reason `OPERATOR` (not `CUTOVER`) is refused at A4 and A5 | G4 | 06 |

Named negative fixtures (R2), each a required test with the stated observable result:

- N1a: supervisor S commits `EFFECT_INTENT(PROCESS_SPAWN)` and forks the leader; before the
  leader links `.boot`, an unrelated `ticket refine` on another ticket runs under the lock. It
  must create no `.boot`/`.ack`, write no `EFFECT_OUTCOME`, and leave the attempt `ADMITTED`
  with the effect `PENDING`; S then reaches `RUNNING` normally.
- N1b: as N1a but `.boot` exists and the leader is waiting; a concurrent `queue status` and a
  concurrent `ticket hold` both report/leave the attempt `ADMITTED` with `pendingEffects`
  listed, never `BLOCKED_RECOVERY/SUPERVISOR_LOST` or `UNCERTAIN_EFFECT`.
- N1c: `reconcile` against a pending spawn whose `attempt.supervisor` is live by §6.3 is refused
  `BLOCKED/EFFECT_OWNED` with a byte-identical state dir.
- N1d: S is SIGKILLed after `EFFECT_INTENT`; `reconcile` proves S dead, fences the absent
  record(s), and records `FAILED/SPAWN` with `noExec:PROVED`; a leader started afterwards gets
  `EEXIST` on `.boot` and exits 3 having exec'd nothing.
- N1e: a `TRANSITION` to `RUNNING` carrying a `supervisor` value different from
  `attempt.supervisor` (including one issued after `reconcile` nulled it) is rejected
  `REVISION_CONFLICT/FENCED` and recorded; the attempt phase is unchanged. A `.boot` whose `pid`
  differs from the forked pid yields `FAILED/SPAWN`, never `RUNNING`.
- N2a: the lock is held by another process for 30 s while the leader waits; S commits `RUNNING`
  within the 90 s leader wait and links `ack:true`; the runtime execs exactly once.
- N2b: S is delayed past the leader's 90 s wait (fault injection); the leader links a `noexec`
  `.ack` and exits 3; S's `.ack` link fails `EEXIST`, S records `RUNNING → FAILED/SPAWN` with
  cause `NOEXEC`, `noExec:PROVED`, `retryCount` unchanged and `spawnNoExecCount` +1; the phase is
  never `FAILED/BUILD`.
- N2c: the leader is SIGKILLed before linking `.boot`; S observes the child exit, fences `.boot`,
  records `FAILED/SPAWN` with `BOOT_TIMEOUT` and quiescence `PROVED` (forked pid reaped); no
  90 s wait elapses.
- N2d: a `.boot` truncated by fault injection (simulating corruption) with `.ack` absent:
  recovery fences `.ack`, records `noExec:PROVED`, keeps quiescence `UNPROVED` and the
  reservation until the forked pid is known dead; a second spawn for the attempt is refused.
- N2e: `cancel` of an `ADMITTED` attempt whose spawn is pending and whose supervisor is live
  writes `STOPPING` only; S at step 4 fences `.ack`, records `EFFECT_OUTCOME FAILED` and
  `CANCELLED`; the leader exits 3; `cancel` itself never touches `effects/`.
- N2f: three consecutive `FAILED/SPAWN` with `noExec:PROVED` leave `retryCount` unchanged; the
  fourth spawn attempt is refused `RETRY_EXHAUSTED` by the §1 spawn-failure cap.
- N3a: file sets `status:"COMPLETED"` → `VALIDATION_FAILED/ADOPT_UNSUPPORTED_FIELD`; the
  ticket stays `INTENT_DIVERGED`, has no `completion`, and the state dir is byte-identical.
- N3b: file adds an `approvals` entry → refused as N3a; no grant exists afterwards.
- N3c: file hand-bumps `revision` or `acceptanceRevision` → refused as N3a.
- N3d: file changes only `priority` and `order` → one `RECONCILE` receipt composing
  `PRIORITIZE`; `revision` +1, `acceptanceRevision` unchanged, `previousRecordSha256` = the
  canonical record's digest (not the file's).
- N3e: file changes `acceptanceCriteria` while an attempt is live → `BLOCKED/ATTEMPT_LIVE`;
  the ticket stays diverged.
- N3f: file changes `dependencies` to introduce a cycle → `VALIDATION_FAILED/CYCLE`; stays
  diverged.
- N3g (TCP-01 mutation repair): a `DRAFT` canonical record with empty criteria; the file adds
  one criterion and leaves `status:"DRAFT"` → one `REFINE` composed, the post record is `OPEN`
  with `acceptanceRevision` +1; the same file with `status:"OPEN"` adopts identically; a file
  that writes `OPEN` without criteria is refused `ADOPT_UNSUPPORTED_FIELD`.
- N3h (TCP-01 mutation repair): `OWNER` adopts file F under `adopt-1`; `OPERATOR` (or the same
  id under another role) issues `adopt-1` with F → `REQUEST_ID_CONFLICT`, never the owner's
  replay; the same owner under `adopt-1` for another target or other bytes → conflict; the
  owner's identical retry after commit → replay with `replayed:true` and no post record.
- N3i (TCP-01 mutation repair): an `ARCHIVED` canonical record with an identical, `updatedAt`-only
  or `title`-edited file → `BLOCKED/TICKET_STATE`, nothing composed, no revision minted.
- N4a: `receipts/<lastSeq+1>.json` pre-planted with a foreign `prev` → commit path reports
  `JOURNAL_FORKED`, no redo, no write.
- N4b: `receipts/<lastSeq+2>.json` present (with `<lastSeq+1>` present or absent) → both the
  mutating path and every read report `JOURNAL_FORKED`; nothing is redone or renamed.
- N4c: the commit's `link(2)` target already exists → `JOURNAL_FORKED`; the pre-existing file's
  bytes are unchanged and the temp file is removed.
- N4d: the R2 counterexample: export snapshot at N, commits N+1 and N+2 injected between the
  head read and the receipts copy → export reports `NOT_RUN/SNAPSHOT_MOVED` with zero bytes on
  stdout; after three retries the same; no torn archive exists.
- N4e: a hand-assembled archive stream with `receiptCount` = `lastSeq` + 2 → `archive verify`
  fails `JOURNAL_FORKED`; `archive restore` refuses before writing.
- N4f: `archive export` with `--staging` inside the state dir, inside any worktree resolving to
  the repository's common directory, under malformed or unreadable enclosing Git authority, or
  on a symlinked path, is refused; absence of enclosing Git authority is allowed. A successful
  export leaves the state dir, `taskman.lock` mtime and the repository byte-identical and creates
  no lock.

Gate rows G0..G6 are ATCP §"Acceptance and rollout" verbatim. G0 is the document/protocol freeze
plus this fresh independent review and needs no code; executable schema and limit tests are G1.
GP (performance, §9.4) is additional and owner-mandated: it must pass before any Corvint-facing
promotion (TCP-03, TCP-06) and is re-run for TCP-07's fanout envelope. M2 (one real ticket,
explicit reviewer, capacity one) requires TCP-01..06 with G1..G4 and GP, which includes G2 and
G3 and the ATM S1/S2/S3 equivalents. Automatic routing (TCP-08, G6) and foreign adapters
(TCP-09) are qualified separately and later.

### 9.4 GP: performance gate and decision rule (TM-V0-027)

Preregistration: the `taskman-perf-baseline/0` document (§3.5) is frozen and pinned before the
first Corvint Go edit. It reuses the existing Corvint `conformance/perf-v0` pinned Corvint and Beamfall
corpora, cache-state controls, five warmups, ≥100 fresh-process samples per row, alternating
order, invalidity and output-parity checks, bootstrap infrastructure and cleanup tests, as the
coordinator's source check reports them; it does **not** reuse that harness's Python-versus-Go
comparison, `--candidate-revision` or reduced-sample modes, which cannot certify this change.
Rows cover `help`, `version`, bare startup and every existing Corvint command, including the
commands the current migration manifest omits; the command corpus is enumerated from the Corvint
CLI at freeze time, not from the manifest. Witnesses: wall p50/p95, user and system CPU, max
RSS, allocations, I/O read/write bytes, stdout/stderr identity and exit code, per row.

Comparison: old-Go (baseline binary digest) and new-Go (candidate) are run directly
interleaved on the same host, same corpus state and same cache state, in one session, with the
load condition recorded. Runs on this shared host while other full gates are executing are
recorded `CONTENDED` and cannot establish the baseline; they are kept, never discarded.

Decision per row (command × corpus × cache state × platform × `corvint-tasks` condition): zero output
differences, errors or failures; one-sided 95% upper confidence bound (bootstrap) on candidate
minus baseline p95 ≤ 0; the same no-increase test on required CPU and RSS witnesses; every
absolute Corvint budget still met. Interval crossing zero → `INCONCLUSIVE`, baseline retained, no
promotion. Absent or invalid measurement → `NOT_RUN`. A pass is never stated as "unchanged" or
"no regression proved"; the recorded statement is the bound and the sample. No positive noise
allowance exists. `corvint-tasks` conditions `IDLE_INITIALIZED` and `MAX_ADMITTED` are required rows;
the coordinator/child CPU, RSS, I/O, process-count and polling budgets in `atmBudgets` are
measured under them before any value is written, and a worker-count cap by itself is not
evidence of absent contention. Nothing in this section is measured at this revision; GP is
`NOT_RUN`.

### 9.3 Rollback per slice

| Ticket | Rollback |
|---|---|
| TCP-00 | none; documents only |
| TCP-01 | disable native writer (`canonicalWriter` unchanged, mutations refused); intent files restored from Git; verified archive retained |
| TCP-02 | `pause`; keep unresolved reservations; no state deletion |
| TCP-03 | disable operative planning; shadow WQO observation remains; any Corvint-side edit reverts to the frozen baseline tree if GP is not passed |
| TCP-04 | development runs only, gate-incomplete; cancellation retains work |
| TCP-05 | no completion while any mandatory lane is missing |
| TCP-06 | §5.4 revert after freeze/drain/archive; history preserved |
| TCP-07 | `maxActiveAttempts` back to 1; reservations kept until `PROVED` |
| TCP-08 | explicit reviewer mode |
| TCP-09 | foreign execution disabled |

## 10. Unresolved items

- U1: The Corvint-side records (ATCP intent status, the ATM amendment file, `docs/specs/README.md`)
  cannot be written from this repository. Decision 0001 accepts §8 for this repository; the
  Corvint repository still needs its own amendment commit.
- U2: TCP-09 (Beamfall adapter and recorder) has no owner; it does not block Corvint-only serial use.
- U3: Go 1.27.1 `make verify` executes format, all tests and vet. The 2026-09-19 baseline passed; writer-repair validation is recorded separately. Passing this gate does not qualify G1..G6 or GP.
- U4: Darwin `kern.proc.pid` start-time extraction and Linux `btime` parsing are asserted
  stdlib-feasible; TCP-02 must prove it with a synthetic process-table fixture.
- U5: resolved in R2. The coordinator computed SHA-256 over the six frozen context files and
  decision 0001 §"Provenance" records them; the reference is path plus digest, and this still
  says nothing about commit membership in Corvint.
- U6: `taskman-perf-baseline/0` names the Corvint `conformance/perf-v0` corpora and controls on
  the coordinator's source check; this repository has not opened that harness. The baseline
  freeze (TCP-03) must cite its exact path and digest.
- U7: The CEM/OCM profile names and status vocabulary used in §7.3 row 10 (`cem-status:ready`)
  are taken from the Corvint review record; TCP-05 must pin the exact Corvint profile identifiers
  before implementing row 10, without weakening the "status complete" requirement.

## 11. Closed detail codes

`ADJUDICATION, ADOPT_UNSUPPORTED_FIELD, APPROVAL_MISSING, APPROVAL_REVOKED, ATTEMPT_LIVE,
BOOT_FENCED, BOOT_TIMEOUT, BUDGET_EXCEEDED,
BUDGET_UNKNOWN, CAPABILITY_UNAVAILABLE, CEM_MISSING, CONTAMINATED, COVERAGE_UNKNOWN,
CUTOVER_IN_PROGRESS, CUTOVER_MISSING, CYCLE, DEPENDENCY_MISSING, DEPENDENCY_UNSATISFIED,
DEVELOPMENT_MODE, DIRTY_WORKTREE, DOCS_MISSING, DUPLICATE_ID, EFFECT_OWNED, EXTERNAL_UNBOUNDED, FENCED,
GATE_FAILED, GATE_STALE, GATE_UNKNOWN, INDEPENDENCE_UNVERIFIED, INTENT_BRANCH_MISMATCH,
INTENT_DIVERGED, INVALID_PRIORITY, JOURNAL_FORKED, JOURNAL_SATURATED, LIMIT_EXCEEDED,
LOCK_TIMEOUT, MALFORMED, MISSING_EVIDENCE, MISSING_GATE, NOEXEC, OCM_MISSING, OUT_OF_SCOPE, PAUSED,
PLAN_STALE, QUIESCENCE_UNPROVED, REDO_PENDING, REQUEST_ID_CONFLICT, RESOURCE_COLLISION,
RESTORED, RESTORE_INCOMPLETE, RETRY_EXHAUSTED, REVIEW_INCOMPLETE, REVIEW_REJECTED,
SIGNAL_REFUSED_IDENTITY, SNAPSHOT_MOVED, STALE_POLICY, STALE_TICKET, STALE_TREE,
SUPERVISOR_LOST, SURVIVORS, TICKET_HELD, TICKET_STATE, UNCERTAIN_EFFECT, UNINITIALIZED,
UNPUBLISHED, UNRESOLVED_FINDING, UNSUPPORTED, UNSUPPORTED_FILESYSTEM, UNSUPPORTED_VERSION`.

## 12. Traceability

| TM-V0 | ATCP-V0 | ATM-V0 | Ticket | Evidence |
|---|---|---|---|---|
| 001 | 001, 007, 018 | 005a, 028 | 02 | AS-13, AS-29 |
| 002 | 002, 018 | 019 | 01 | AS-01, AS-10 |
| 003 | 002 | — | 01 | AS-02, AS-05 |
| 004 | 004, 005 | 005 | 01, 05 | AS-06, AS-19, AS-25 |
| 005 | 003, 011 | — | 01 | AS-02, AS-04 |
| 006 | 003 | — | 01 | AS-03 |
| 007 | 001, 010, 018 | — | 01, 02 | AS-11 (C6, C6a, C7), AS-28, AS-35 |
| 008 | 001, 016 | 019, 020 | 01 | AS-07, AS-08, AS-21, AS-23, AS-36 |
| 009 | 010, 018 | 019a, 028 | 02 | AS-11, AS-12, AS-36 |
| 010 | 018 | 028 | 02 | AS-29 |
| 011 | 010, 012 | 005a | 02 | AS-14 |
| 012 | 006, 007, 008 | 005a, 007 | 02 | AS-13, AS-15 |
| 013 | 008 | — | 02, 03 | AS-15, AS-30 |
| 014 | 015 | 006 | 04 | AS-20 |
| 015 | 006, 009 | 002, 002a | 03 | AS-21 |
| 016 | 012 | 005a, 005c, 028 | 02, 04 | AS-16, AS-18, AS-19, AS-33, AS-34 |
| 017 | 010 | 001c | 02, 04 | AS-16, AS-17, AS-33, AS-34 |
| 018 | 011 | 008, 008a, 015, 024a | 04, 05 | AS-24 |
| 019 | 013 | — | 05 | AS-22 |
| 020 | 014, 021 | 012, 013, 013a, 014, 014a | 05 | AS-23 |
| 021 | 005, 022 | 016, 017, 018, 019 | 05 | AS-25 |
| 022 | 017, 018 | §8 | 01, 02 | AS-09, AS-27, AS-36 |
| 023 | 001, 017, 019 | 001 | 06 | AS-26, AS-27, AS-37 |
| 024 | 018 | — | 02, 06 | AS-18 |
| 025 | 020 | §7 | all | G0..G6, GP |
| 026 | 011, 014 | §5 | — | AS-17, AS-28 |
| 027 | 020 (evidence rule); owner steering 2026-09-06 | — | 03, 04, 06, 07 | AS-31, AS-32 (GP) |
| 028 | owner decision 0009 | ATM-V0-028 | 02 | AS-38 |
| 029 | owner decision 0009 (public read closure for 028) | ATM-V0-028 | 02 | AS-38 |

The table maps obligations to slices and acceptance scenarios. Executed support/fixture evidence is recorded in `docs/reviews/`; unimplemented runtime and qualification scenarios remain NOT_RUN. A passed package or repository gate does not fill a qualification cell.

### TM-V0-028 — fixture release control

Candidate ticket and predecessor bindings are semantic arrays in their release definition's ID order; promotion predecessor bindings retain the candidate order. Whole-object digest-byte sorting does not apply to these arrays. Their IDs must exactly match the sorted, unique definition, so duplicate, missing, extra, or reordered bindings remain malformed.

`taskman-release/0` records live at `.taskman/releases/<release-id>.json` under the release-count and file-byte limits. Definitions name ordered predecessors, scoped tickets, acceptance criteria, and required policy gates. Candidate capture, attestations, readiness, and promotion follow decision 0009. Required gates are derived from the current policy's required gates union the release's explicit gates; every release criterion must be covered by passing, compatible current-candidate evidence. Readiness is `BLOCKED`, `UNKNOWN`, or `READY_ATTESTED`. External/manual evidence is non-native, and native gate execution is `NOT_RUN`. Release mutation is fixture-only and journaled; promotion is local and grants no publication or real-queue authority.

Definitions reject unknown required gate IDs before commit. Attestations must supply the current candidate digest; the writer never rebinds stale evidence. Promotion binds only passing, compatible evidence for that candidate. Release CLI payloads accept insignificant JSON whitespace (including stdin) and derive authorization from parsed provenance. Source identity includes file mode and type plus file contents or symlink target. Every writer carries the complete canonical ticket and release inventory.

TM-V0-029 closes the public read surface: `release list` is a bounded summary. `release show` and `release readiness` additionally expose the complete definition, nullable `candidateSha256` and candidate bindings, attestations with evidence and canonical `attestationSha256`, and nullable promotion detail and canonical `promotionSha256`; readiness adds its state, missing conditions, and `nativeGateExecution`. These public canonical digests are sufficient to record compatible evidence and bind a successor without reading `.taskman` projections.

AS-38 proves two ordered releases, a public-output-only candidate/attestation/promotion workflow and predecessor binding, missing-gate blocking, compatible attestation readiness, promotion, invalidation by source/ticket/policy/predecessor change (including chmod), policy removal, rejection of unknown gates and stale attestations, manual payload routing, ticket/barrier/reconciliation interoperability, durable replay and pending-receipt redo, and refusal outside the fixture boundary. Active staging recovery remains NOT_RUN.
