# Decision 0012: millisecond lookups at 10,000+ tickets, a Beamfall-sized workspace, and an enterprise seam

- Date: 2026-09-24
- Status: proposed. The owner direction below is accepted. The architecture choices A1 to A9 were
  made by an expert panel and become accepted when the owner merges this decision. Choices marked
  **owner** change a published contract and need that acceptance explicitly.
- Owner: Russell Lewis

## Owner direction (2026-09-24, exact)

- "I want you to ensure that corvint-tasks data structures are optimized and world class enough to
  support 10,000+ tickets and still have quick retrieval and lookup. I want all lookups to be
  milliseconds. you can create a new task to work on this and summon experts to ensure we have a
  solid and amazing plan"
- "we need to make sure that we can handle projects the size and multi repo nature of something like
  beamfall. and then do the same with even larger project, maybe even enterprise level"

## How this was decided

Three read-only reviews worked from one brief. None of them ran a build or benchmark:

- a performance and scalability engineer: cost model, budgets, benchmark method;
- a storage and migrations engineer: invariants, crash safety, migration, growth;
- an embedded on-disk index designer: file layouts, lookup algorithms, the multi-repo seam.

Their citations were spot-checked against the code. Every millisecond figure below is an estimate,
not a measurement. TCP-15 measures them before the design is locked.

## Where corvint-tasks stands (confirmed in code)

1. **Every read is O(whole store).**
   - `withStore` (`internal/cli/cli.go:289-323`) reads and hashes the whole intent tree three times
     per attempt (`intent.TreeDigest`, `internal/intent/store.go:180`):
     - in the snapshot probe (`cli.go:296`);
     - in `LoadExpecting` (`internal/intent/store.go:370`), which also decodes it once (`:379`);
     - in the re-probe (`internal/snapshot/probe.go:178`).
   - It then builds the full inventory and runs Tarjan over it (`internal/ticket/graph.go:26-72`).
   - `ticket show` does all of this to answer about one ticket. Estimate: about 1-2 s warm at 10k.
2. **Reads fail while a commit is in flight.**
   - Once a receipt is linked in and until `head.json` advances, the probe returns
     `REDO_PENDING` (`internal/snapshot/probe.go:127-142`). The read does not retry
     (`probe.go:173-181`); TM-V0-008 requires this.
   - With 5-20 writers this window recurs for the posts-and-head part of every commit.
3. **Every write is O(whole history) inside the lock.** `Mutate` takes `taskman.lock`
   (`internal/store/mutate.go:50`) and then makes three full passes:
   - the request lookup (`mutate.go:76`, `internal/journal/index.go:28`);
   - the capacity inventory (`mutate.go:88`), which reads and hashes every retained receipt,
     request, evidence and pinned file (`internal/store/scan.go:24`, `:182-194`);
   - the canonical audit (`mutate.go:98`).

   The two audits read every receipt (`internal/journal/records.go:78-83`). The lock is taken by
   polling every 10 ms, is unfair, and times out at 30 s (`internal/authority/authority.go:38-43`).
   Estimate: 20 agents claiming at once at 10k tickets pushes the last one past the timeout.
4. **Hard walls.**
   - Tickets are capped at 10,000 per queue, tombstones included (`internal/wire/limits.go:40`).
   - The intent tree is capped at 256 MiB.
   - The journal is capped at 1,000,000 receipts and 4 GiB (`internal/transaction/capacity.go:16-17`).
   - A single transaction may stage at most 11 artifacts (`internal/store/store.go:272`), so a
     100-record import batch cannot commit.
   - `import` is not built (`cli.go:61`).
   - `JOURNAL_SATURATED` is declared (`internal/wire/codes.go:40`) but no writer raises it.
   - Archive entries (at most 2,100,000) run out before the receipt cap does.
5. **One queue per store.** A dependency outside the queue is refused
   (`internal/ticket/record.go:377`).

Beamfall today (measured read-only on 2026-09-24):

- 28 Git repositories.
- One central roadmap of 2,501 tickets, covering about 25 of those repositories. The primary repo
  split is `beamfall` 971, `beamfall-apple` 181, `beamfall-android` 121, `beamfall-web` 89, and so on.
- 149 tickets name more than one repository, and 14 mention `AtomicCrossRepo`.

## Targets

- **"Lookup"** means `ticket show`, `ticket blockers`, `gate show`, `release show`, `queue status`,
  and one page (default limit 100) of `ticket list`, `search` or `roadmap`.
- **"Milliseconds"** means the wall time of a fresh process: exec, resolve, open, answer and write
  the output. Nothing is held in memory between calls.
- Budgets are reported both as a total and as the amount above the floor. The floor is a no-store
  verb on the same host.

| Operation | warm p50 / p99 | cold p99 |
|---|---|---|
| `show`, `blockers`, `gate show`, `release show`: 10k and 100k (independent of N) | 8 / 20 ms | 60 ms |
| One page of `list`/`search`/`roadmap`: 10k / 100k | 15 / 40 ms; 25 / 60 ms | 120 / 200 ms |
| `queue status` | 10 / 25 ms | 100 ms |
| One query across a Beamfall workspace (Tier B) | 20 / 50 ms | 200 ms |
| One write (create/refine/claim), unloaded, macOS; independent of history length | 60 / 150 ms | — |
| 20-writer claim burst | every writer commits, zero `LOCK_TIMEOUT`; the last one within 20 × measured single-write p50 + 250 ms (about 1.5 s at the estimate) | — |

Text search is the one lookup with a known worst case:

- A selective query stays in milliseconds.
- An unselective long query can cost time proportional to the total text in the store: about
  30-80 ms at 100k. It is reported, not hidden.

Release readiness hashes the repository's tracked files, not the store. It is not a lookup unless
its source digest is cached.

## Scale tiers

- **Tier S (single store):** 10,000 live tickets; 100,000 including completed and archived history.
- **Tier B (Beamfall workspace):** about 30 repositories, 10k live / 100k with history, 5-20
  concurrent agent writers, cross-repository dependencies and atomic cross-repository tickets.
- **Tier E (enterprise):** hundreds to thousands of repositories, 1,000,000+ tickets, hundreds of
  concurrent writers, many teams.

The local single-binary path officially supports up to:

- 100 stores in one workspace manifest;
- 100,000 tickets per store;
- about 20 writes per second per store. The macOS `F_FULLFSYNC` durability floor puts the real
  ceiling at about 10-50 commits/s.

Past those limits it refuses, and says so. Tier E's global queries, hosted leases and hundreds of
writers need their own accepted profile on Corvint's transport-neutral immutable index contract
(`docs/specs/deployment-neutral-index-platform-v0.md` in Corvint; Corvint invariant 7). They must
not silently widen the local path. This decision only fixes the seam, so that Tier S and Tier B do
not block Tier E.

## Decision

### A1. A derived, immutable index outside Git and outside the state directory

- **What it is:** `taskman-index/0`, stored at `<git-common-dir>/taskman-index/`, beside
  `taskman.lock`.
- **Not Git-tracked.** Every write would otherwise conflict on merge, and a checked-out index would
  describe another branch's history, while the journal is one per common directory.
- **Not in the state directory.** Older binaries' audits refuse unknown entries there
  (`internal/journal/audit.go:377-433`), and archive export would pick it up.
- **Not a journal post and not a `head.json` field.**
  - As a post, it would claim a cache as history and use up evidence and archive budget.
  - As a head field, it would change the frozen head schema (at most 4 KiB) and couple a cache to
    the commit point.
- **Built only from the journal**, meaning receipt posts plus evidence. It is never built from the
  `.taskman/` files, which humans and `git checkout` change without taking the lock.
- **Layout:**
  - one root file, `CURRENT`, replaced by rename;
  - content-addressed segments `seg/<sha256>.cti`, made of fixed-width little-endian tables in
    checksummed blocks under a Merkle root;
  - no timestamps and no map-order dependence, so a rebuild is byte-deterministic.
- **Segments:**

  | Segment | Contents |
  |---|---|
  | `ids` | sorted (queueId, folded local id) → slot, for binary search |
  | `slots` | canonical sha256, latest post seq, status, kind, priority, order and ranks |
  | `order` | the SPEC §4.3 order and the roadmap order as permutations with inverse ranks |
  | `postings` | roaring-style postings for status, kind, priority, owner, milestone, label, required gate and repository |
  | `graph` | forward and reverse dependency CSR plus SCC ids, so blockers and cycles are precomputed |
  | `rows` | the exact encoded compact view and roadmap row |
  | `text` | a lower-cased text column, trigram postings, and a dirty-slot set |

- **Estimated size:** about 30 MiB at 10k and about 300 MiB at 100k, mostly text.
- **Reads use `ReadAt`, not mmap.** A truncated file must produce a refusal, never SIGBUS.
- **Rejected:** a git-index-style stat cache over the projection. It still stats all N files on
  every read: about 15-30 ms at 10k and 150-300 ms at 100k. It is kept only for
  `index verify --projection`.
- **Rejected:** a minimal perfect hash. A 17-probe binary search already takes microseconds and is
  deterministic.

### A2 (owner). Lookups answer from the journal's canonical state at head

This is the only way a lock-free read can be both snapshot-consistent and O(1) to check. `head.json`
is the one file that changes atomically; the projection is N separately editable files with no
commit point.

- In index mode a read works like this:
  1. Probe as today: O(1).
  2. Open `CURRENT` once, and hold its descriptor and every segment descriptor for the rest of the
     read.
  3. Accept the root only if its magic, format version, view-algorithm version, Unicode version,
     `queueId`, `initSha256` and `primaryWorktreeSha256` match, and its `lastReceiptSha256` lies on
     the `prev` chain from the probed head, at most 16 receipts back. The reader applies those
     receipts' posts in memory. This keeps a read fast in the window where head has advanced and
     the index has not (A4).
  4. Answer the query, verifying every block it touches against its digest.
  5. Re-probe, and accept only if the head is byte-identical.
  6. Keep the three retries of TM-V0-008.
- **A commit in flight does not fail the read.** When exactly one receipt sits beyond head, an
  index-mode read answers at `head.lastSeq`, the settled state before that receipt, and the
  envelope names the pending sequence. Today this is `REDO_PENDING` (fact 2). A receipt two beyond
  head is still `JOURNAL_FORKED`.
- **Projection files are checked, never trusted.** Every record a read returns is re-hashed against
  its slot.
  - If the projection file matches, it is served.
  - Otherwise the canonical bytes are served from the journal, and the item is marked `DIVERGED`.
  - A ticket file with no journal slot (for example, a teammate's new ticket after `git pull`) is
    not in any answer. A lookup of that ID answers `NOT_FOUND` and reports
    `UNTRACKED_PROJECTION`, found with one `lstat`. `index verify --projection` and
    `receipt audit` list all of them.
- **Inputs outside head are computed live** at read time and never stored in rows: attempt
  liveness, lease expiry by wall clock (TCP-12), Git `HEAD` for publication, and the barrier file.
- **Contract changes:**
  - Whole-tree projection agreement moves out of the lookup and into `receipt audit`.
  - `snapshot.intentTreeSha256` comes to mean the canonical tree digest at head, stored in the
    root. The pagination pin becomes (`lastSeq`, canonical tree digest).
  - This amends:
    - TM-V0-008;
    - §4.3;
    - the envelope digest meaning;
    - the digest preimage and the "flat and closed" layout (`docs/SPEC.md:203-210`), which A6's
      fan-out also changes.
  - TM-V0-015 plan staleness compares that digest (`docs/SPEC.md:978-983`). An edit to a
    projection file alone then no longer makes a plan stale. TM-V0-015 cites ATCP-V0-006 as
    "preserved unnarrowed", and `AGENTS.md` ranks ATCP above the SPEC, so this is an upstream
    amendment. The ATCP source text is still missing (SPEC §0).
- If the owner rejects A2, the fallback is the stat cache rejected in A1. That keeps projection
  semantics, and lookups stay O(N) stats.

### A3. A missing or stale index is slow, never wrong

- A missing, torn, foreign, ahead-of-head, too-far-behind or unknown-version index is treated as
  absent.
- The read then takes the **canonical slow path**. It computes the same journal-canonical answer by
  replaying the journal, from genesis or from a checkpoint (A6), without the index. It is O(history)
  and gives the same bytes. Like index mode, it re-hashes each returned record against its projection
  file and applies the same `DIVERGED` and `UNTRACKED_PROJECTION` markers, which costs
  O(returned records). Today's projection read is kept only inside
  `index verify --projection`.
- A block that fails its digest mid-query discards the partial answer and reruns the query on the
  slow path. It is never reported as an error.
- Readers never write, repair or lock the index. Invariant 2 and AS-07 byte identity extend to
  `taskman-index/`.
- A differential test is mandatory and runs on every change. Every read verb in index mode must be
  byte-equal to the canonical slow path across randomized mutation sequences, including
  projection edits and a lagging index. The markers are part of the compared bytes.

### A4 (owner, reader leases only). The writer updates the index after the head and never fails a commit on it

- The §5.2 order becomes: evidence → receipt link-in (commit point) → posts → head → index, all
  under the same `taskman.lock`.
- **Incremental update:**
  - It covers the posted slots, their reverse dependencies and any SCC changes.
  - It copies only the touched blocks under a new Merkle root: O(touched · log N), not a rewrite of
    whole tables.
  - It appends rows.
  - Text edits go to the dirty-slot set. The trigram segment is compacted when the dirty set
    exceeds max(1024, N/32).
  - It writes new segments, then `CURRENT.tmp`, fsyncs it, renames it and fsyncs the directory.
- **Full rebuild** happens on a policy change, a root that has fallen too far behind, or a broken
  `prev` chain from the root's `lastReceiptSha256`.
- **The view-algorithm version is derived from the build identity** (VCS revision and module sums
  from `runtime/debug.ReadBuildInfo`), not bumped by hand. A build without that identity never
  trusts an existing index; it rebuilds on its next write.
- **New crash points:**
  - C5a: after head, before the index rename. The index lags; reads apply up to 16 receipts
    (A2 step 3) or fall back, and the next writer catches up.
  - C5b: a stray index temp file. Only the index writer removes it, under the lock.
  - C5c: a rename without a data fsync. The root digest fails, and reads fall back.
  - C5d: complete segments that no root names. GC removes them under the lock.
- **Every path that advances head** must catch the index up or leave it detectably behind: redo
  (`internal/store/redo.go:79,189`), `init` (`internal/store/store.go:430`), and the release and
  policy writers.
- **GC and reader leases:**
  - GC keeps the generations named by `CURRENT` and `PREVIOUS`. A reader holds its descriptors, so
    an unlinked segment stays readable to it on APFS and ext4.
  - A reader that gets ENOENT before opening retries as if the store had moved. Those retries use
    up the three TM-V0-008 retries; keeping two generations means a reader has to lose two commits
    in a row to run out.
  - This is the local equivalent of a reader lease, and it takes no lock (**owner**: DNIP-LOC-010
    asks for explicit leases).
- **Index verbs:**
  - `index rebuild` takes the lock and writes no receipt.
  - `index verify` is a pure read that reports `FRESH`, `BEHIND`, `ABSENT`, `CORRUPT` or
    `MISMATCH`. `index verify --projection` also checks every tracked file.
  - The SPEC must say that a derived-cache write is not a mutation.

### A5 (owner). Writes stop replaying the whole history

- **Request replay becomes an O(1) read** of the content-addressed request file, checked against
  its receipt, with a receipt-digest table in the index.
- **Records a mutation touches are checked by hash through the index.** The full audit leaves the
  lock path.
- **Capacity is accounted incrementally.** Per-field totals are kept in the index root, and each
  commit adds its own delta. `receipt audit` recomputes them from a full inventory and reports any
  mismatch. With no usable index, the writer falls back to today's full inventory (fact 3).
- **Lock fairness is measured before it is built.**
  - `flock(2)` promises no first-come order on macOS or Linux. The standard library cannot bound a
    blocking `flock` wait without leaving a goroutine behind.
  - The candidate is a ticket queue: each waiter creates a sequence file under `taskman-lockq/`,
    waits only for its own turn, and a waiter whose process is gone is skipped by the same liveness
    check attempts use.
  - TCP-18 keeps today's unfair poll if the burst budget holds without the queue.
- **Cost:** tampering with an old receipt is no longer caught on every write, only by
  `receipt audit`, `archive export` and `index verify`. That weakens the "each lookup re-audits"
  comment (`internal/journal/index.go:8-9`), and the SPEC must say so.
- **Safeguard:** the randomized differential test against a full audit (Gates) runs on every
  change.

### A6 (owner). Raising the limits is a one-way store migration to `taskman-state/1`

- **Limits:**

  | Limit | Today | New |
  |---|---|---|
  | Tickets per queue | 10,000 | 100,000, with archived tickets excluded from the default load |
  | Intent tree | 256 MiB | 1 GiB, or archived tickets moved out of the tree budget |
  | Import map | 10,000 entries in one file | 100,000 entries, sharded per source |
  | Import source | 10,000 | 100,000 |
  | Artifacts per transaction | 11 | at least 210 (TCP-21; it lands with this migration only if older binaries cannot read such a transaction) |
  | `queue.json` | 1 MiB | unchanged; it does not grow with tickets |
  | JSON decode bounds | unchanged | unchanged; this is why the index is binary |

- **Layout:** `tickets/` fans out as `tickets/<xx>/`. A flat directory of 100k files makes a
  multi-megabyte Git tree object that is rewritten on every commit.
- **Retention:** the B4 saturation accounting is wired (`JOURNAL_SATURATED` raised before a store
  becomes impossible to export), and a `CHECKPOINT` receipt follows an archive-verified `PRUNE` so
  that audit can start from a checkpoint. `PRUNE` is unsupported today.
- **Migration:**
  - It is explicit and resumable, and exports and verifies an archive first (TM-V0-024). A receipt
    kind that exists only in `/1` records it.
  - Older binaries fail closed with `UNSUPPORTED_VERSION`.
  - Rollback means restoring the pre-migration archive, so in practice it is one-way.
- **The index alone needs no migration.** It is a cache: deleting `taskman-index/` reverts it, and
  older binaries ignore it.

### A7. Tier B: one workspace store, with the repository as a ticket attribute

This matches Beamfall's single central roadmap. Beamfall's workspace fits in one Tier S store.

- **Repository is a field on the ticket.** TCP-10 maps `Primary Repo`/`Repo` into the record or a
  recorded extension. The index gains a repository posting and per-repository counters.
- **Joining a workspace.** A member repository joins through a `workspace join` receipt and a small
  membership record in the member's common directory that names the workspace store.
  - The store's registry names the member back.
  - Readers verify both directions or refuse `UNSUPPORTED`. §3.4 forbids relocating the state
    directory, so this is a SPEC change.
- **Long work never holds the journal lock.**
  - Per-repository merge locks live at `taskman-merge/<repoId>.lock` (TCP-11).
  - Leases are journal receipts (TCP-12).
- **`AtomicCrossRepo` tickets:**
  - Take the merge locks in canonical `repoId` byte order, which is deadlock-free.
  - Run the pushes as `EFFECT_INTENT`/`EFFECT_OUTCOME` receipts.
  - Record the completion for every repository in one receipt, which is the atomic fact. The Git
    pushes are recovered effects.
- **Rejected for Tier B:** a store per repository. The 14 atomic tickets would then need two-phase
  commit across journals, since the `link(2)` commit point is per journal. The 149 multi-repository
  tickets would lose dependencies to `NOT_OBSERVED`, and every read would need a snapshot of about
  30 heads.

### A8. The Tier E seam is federation of independent stores

Tier S and Tier B do these now, so Tier E needs no rework later:

1. Index keys are fully qualified (queueId, local id), even with one queue.
2. The index root is a store root. A federation root is an immutable manifest of
   {queueId, headSha256, rootSha256} under a Merkle digest. A reader verifies it the same way as a
   store root and gets a snapshot of every store from one read. The command envelope reserves a
   `snapshots[]` form.
3. Per-store results merge. Each store's pages are already sorted by the §4.3 key, so a federated
   page is a k-way merge and a total is a sum of popcounts.
4. A cross-store dependency is a pinned observation:
   {storeId, ticketId, seenSeq, seenReceiptSha256, recordSha256, acceptanceRevision}, rechecked on
   every evaluation.
   - If the store is absent, the dependency is `NOT_OBSERVED` and blocks.
   - If the store is rewritten, the dependency gets a fork-class refusal.
   - Unknown is never counted as satisfied.
5. Atomic work across stores (presumed-abort two-phase commit with one coordinator decision receipt)
   is outside the local path. It belongs to the Tier E profile.
6. Blocks, Merkle roots, `ReaderAt` access and manifest-last publication follow DNIP-IDX-005, 006
   and 008, so the same bytes can later be served as authenticated ranges.

Moving from one central store to several means an import plus an authority switch per target store
(IDs change, and old IDs resolve through the import map). It is not a live re-shard.

### A9 (owner). `taskman-index/0` is a local cache, not a DNIP index profile

- DNIP-IDX-009 requires an immutable SQLite safety baseline before a custom pack is selected.
  corvint-tasks invariant 1 is standard library only.
- The local index is a derived cache of one single-writer store, so DNIP-IDX-009 does not bind it.
- The Tier E profile that serves these bytes remotely carries DNIP-IDX-009. Its SQLite baseline is
  measured in a harness outside the `corvint-tasks` binary.

## Rows

| Row | What | Depends on |
|---|---|---|
| TCP-15 | Scale benchmark harness and measured baseline. A deterministic generator calibrated from Beamfall: 85% completed, Zipf dependency fan-in and labels, log-normal bodies, 6% multi-repository and 0.6% atomic. Datasets S-2.5k, S-10k, S-100k with history R=N and R=10N, B-30 with 30 repositories scaled ×4, and a local-break-point E-1k. Journals are built by replaying writes. Fresh-process timing with p50/p95/p99 and bootstrap intervals, cold cache via detach/attach of an APFS disk image (macOS) and `posix_fadvise` (Linux), syscall and bytes-read counts. Writer simulators for 5/10/20 writers plus a burst. | none |
| TCP-16 | Read-path quick wins with no format change: capture and hash the tree once per attempt instead of three times, skip Tarjan for verbs that do not need it. | TCP-15 |
| TCP-17 | `taskman-index/0` (A1–A4): builder, verify/rebuild, index-mode reads, text search, differential and fault-injection tests. | TCP-15, owner A2 |
| TCP-18 | Incremental writer, incremental capacity accounting, lock fairness measured (A5). | TCP-17, owner A5 |
| TCP-19 | `taskman-state/1` (A6): limits, `tickets/<xx>/` fan-out, sharded import map, saturation accounting, checkpoint and prune. | TCP-18, TCP-21, TCP-02 B4, owner A6 |
| TCP-20 | Tier B workspace (A7, A8 items 1–4): repository attribute and postings, membership records, cross-repository reads within budget. | TCP-10, TCP-17 |
| TCP-21 | Per-transaction artifact cap raised to at least 210, so a 100-record `IMPORT_APPLY` batch (`docs/SPEC.md:1117`) commits. It decides whether older binaries can read such a transaction; if they cannot, it lands inside the TCP-19 migration and TCP-06 waits for it. | TCP-02 |

- TCP-06 (P0 import) now also depends on TCP-21, and TCP-10 (Beamfall import) on TCP-19 and
  TCP-21 through TCP-06, because a 100-record batch cannot commit today.
- TCP-12 now also depends on TCP-20, and TCP-20 on TCP-10 for the repository field. There is no
  cycle: TCP-12 → TCP-20 → TCP-10 → TCP-19 → TCP-18 → TCP-17 → TCP-15.
- The Tier E profile is not a row here. It is proposed on Corvint's platform contract when a
  workspace outgrows the local limits above.

## Gates

- **Budget gate.** A timed matrix runs nightly on a quiet host. It fails when the one-sided 95%
  upper bound of p99 exceeds its budget, or when the p95 change against a frozen baseline has an
  upper bound above +5%.
- **Growth gate.** It runs on every change at N = 10k and fails on any growth in syscall count or
  bytes read for a lookup at fixed N. These counts are deterministic, so they cheaply catch a
  regression back to O(N).
- **Acceptance tests** (proposed names; SPEC numbers are assigned when each row lands):
  - index missing, stale or torn ≡ slow path;
  - crash points C5a–C5d;
  - a read during an in-flight commit answers at head, and a read against an index up to 16
    receipts behind is byte-equal to the slow path;
  - a projection-only ticket file answers `NOT_FOUND` with `UNTRACKED_PROJECTION`;
  - an older binary commits, then catch-up, including a forked chain;
  - a restored or rewound store and a foreign index are refused;
  - a projection edit gives `DIVERGED`, not a wrong answer;
  - a read never writes the index;
  - an older reader ignores the sibling index;
  - archive-entry reserve refuses before a store becomes unexportable;
  - `/0`→`/1` migration is archive-first, and an older binary fails closed;
  - a randomized 10k-mutation differential test against a full audit.
- **GP (TM-V0-027) still binds.** The harness lives only in corvint-tasks, adds nothing to Corvint's
  main binary, and never runs on a host while a Corvint release gate is running.

## Scheduling

- This work is not on Corvint's 1.0 critical path.
- It may start before Corvint 1.0, but only on a quiet host.
- It must not change the corvint-tasks commit that a pending Corvint release candidate's companion
  gate ran on until that release is promoted.
- TCP-15 comes first. No format is locked until its baseline confirms the cost model, especially:
  - the per-file cost of the safe-open chain on APFS;
  - Git behaviour with 100k tracked `.taskman` files.

## Rollback

- **Before TCP-19:** delete `taskman-index/` and revert the reader and writer changes. Stores are
  untouched, and older binaries never see the index.
- **After TCP-19:** restore the archive exported before migration. Commits made after migration
  are lost unless exported first.
- **Before TCP-20 ships:** remove the membership records; the repository attribute stays as
  ordinary ticket data.

## Residual risks

- **Owner-level changes.** A2 changes what a read means. A5 moves tamper detection off the write
  path.
- **Builder drift.** If the index builder and `ticket.Inventory.View` drift apart, answers are
  silently wrong. The builder calls the same code and stores its exact bytes, and the differential
  test is mandatory.
- **Unmeasured costs.** Git cost with 100k tracked files, the text-search worst case, and the macOS
  fsync floor under a 20-writer burst are unmeasured. TCP-15 settles them, and the budgets may move
  with evidence.
