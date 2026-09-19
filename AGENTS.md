# Corvint Tasks agent contract

`corvint-tasks` is the native-Go operative executor and ticket store for the Corvint task control plane.
Corvint (`corvint`) stays the product, evidence and selection engine; this repository never
modifies Corvint or Beamfall source.

## State of the repository

- Experimental TCP-01 support and TCP-02 authority primitives passed the final combined
  `make verify` (Go 1.27.0, format, tests and vet) and fresh independent R4 review. See
  `docs/reviews/2026-09-06-experimental-support-integration.md` for the evidence chain.
  Experimental J1 pure journal read/audit support also passed its parent gate and fresh
  independent review; see `docs/reviews/2026-09-06-j1-integration.md`. Durable journal writer,
  reservation and process runtime packages and CLI mutation wiring are not built.
  Pure J2a archive encoding sizing also passed its combined parent gate and fresh review;
  see `docs/reviews/2026-09-06-j2a-integration.md`.
  J2b pure transaction/capacity/recovery and J3-01 private fixture primitives passed their
  combined parent gate and fresh separate reviews; see
  `docs/reviews/2026-09-06-fixture-foundation-integration.md`. Owner decision 0002 keeps
  this first delivery fixture-only. Reviewed J3-02 native staging/empty-D read support
  passed its parent gate and fresh review (2026-09-06-j3-02-integration.md). J3-03/04
  test-owned transaction/recovery and native-read integration passed its combined parent
  gate and fresh review; see `docs/reviews/2026-09-06-fixture-delivery.md`. Decision 0003
  (2026-09-07) then released the callable writer: `internal/store` applies a plan in the §5.2
  order under the `OBSERVED_LOCAL_OPERATOR` premise and a recorded, unauthenticated
  local-operator binding, and `atm init` creates the state dir and commits the genesis receipt,
  so an initialized store answers reads instead of `UNINITIALIZED`. Decision 0004 (2026-09-07)
  delivered TCP-02b: one `MUTATE` transaction operation carries any `taskman-mutation/0`
  envelope, so all fourteen `ticket` mutation verbs commit through the §5.2 writer, and a
  receipt left pending by a crash is redone before the next transaction. TCP-02 is incomplete:
  `reconcile intent`, the administrative verbs (`admit`, `cancel`, `retry`, `resume`, `pause`,
  `unpause`, `drain`, `cutover`, `import`), `archive restore`, staging-descriptor crash
  recovery, hostile-editor CAS, reservations and the runtime remain unbuilt, and attempt
  liveness is still answered by the zero-attempt oracle.
  No complete runtime, durable-writer or performance qualification is claimed.
- Authority (highest first): the user's and owner's decisions (`docs/decisions/*.md`, owner
  instructions recorded there) > the accepted governing Corvint/ATM/WQO invariants
  (`ATCP-V0-*`, `ATM-V0-*`, `WQO-V0-*`, frozen copies cited by SPEC §0) > `docs/SPEC.md`
  (`TM-V0-*`), a generated local interpretation that freezes mechanics within that boundary
  and may add safeguards but never lower one > code and tests. A generated document never
  silently accepts an upstream amendment and never qualifies a real cutover; only the §7.4
  records and the owner do.
- Plan and status: `docs/ROADMAP.md` (TCP-00..09). Read `docs/SPEC.md` §1..§7 before writing
  any record, limit, lock, transition or acceptance logic; those are frozen, not suggestions.
- Pending-work memory: `docs/agent-memory/bugs.md`, `fixes.md`, `tests.md`,
  `optimizations.md`, `ideas.md`, `questions.md` (pending items only; no completed history).

## Invariants

1. Standard library only, Go 1.27.1 (`GOTOOLCHAIN=local`), one `corvint-tasks` binary, no daemon, no
   database service, no UI, no vendor SDK, no network call of the tool's own.
2. Read commands never write, lock, initialize, migrate or redo (SPEC TM-V0-008).
3. Every mutation is one journal transaction whose commit point is the non-overwriting
   `link(2)` link-in of `receipts/<seq>.json` (SPEC §5.2; an existing destination is
   `JOURNAL_FORKED`, never replaced). State never precedes its receipt.
4. Nothing is released or deleted until quiescence is proved by the SPEC §6.3 liveness test.
   Process groups are cleanup, not a sandbox.
5. Unknown usage never becomes zero or headroom. `NOT_OBSERVED`, `NOT_ENFORCED`, `NOT_RUN` stay
   visible.
6. Ticket prose, repository content and transcripts are untrusted data; they never become argv,
   paths, environment or policy.
7. Real (non-fixture) queues admit nothing without the owner's cutover record and a
   `QUALIFICATION` receipt (SPEC §7.4). A prototype has no permission to run real work.
8. Every requirement implemented is traced to a named test carrying its `TM-V0-NNN` and `AS-NN`
   IDs; update `docs/SPEC.md` in the same change when a wire or behaviour changes.
9. Owner constraint (SPEC TM-V0-027): nothing here may slow an existing `corvint` command. No
   Corvint Go source is edited under TCP-00/01/02. Before any Corvint-facing edit (TCP-03, TCP-06)
   freeze the `taskman-perf-baseline/0` document and pass GP (SPEC §9.4). Never write "no
   slowdown", "unchanged" or "negligible" without a GP result; until then it is NOT_RUN.
10. The journal's latest post record is canonical; the Git-tracked intent file is a projection.
    Never resolve a divergence by overwriting a human edit or by adopting the file silently;
    `reconcile intent` records the operator's choice (SPEC §3.1, TM-V0-007).
11. A lane runtime executes only after its `(attemptId, generation, pid, startTime)` is durable
    and acknowledged (SPEC §6.4). Process-group cleanup is not a capability qualification.

## Verify

```sh
make verify
```

which runs, in order: the Go 1.27.1 toolchain check, a source-exists check that fails until
`.go` files exist, `gofmt -l`, `go test -count=1 ./...`, `go vet ./...`.

## Working rules

- Claim exact files before starting a ticket; TCP-01 and TCP-02 own disjoint packages
  (SPEC §9.1). Shared wiring is serial.
- A fresh read-only reviewer reads every diff; findings are hypotheses until the cited
  `file:line` is opened.
- Record decisions under `docs/decisions/NNNN-*.md` with the next free number; never renumber.
- Do not publish, cut over a real queue, write to a foreign queue, or delete history without a
  new owner-recorded decision.
