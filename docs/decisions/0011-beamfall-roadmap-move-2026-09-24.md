# Decision 0011: Beamfall's roadmap moves to corvint-tasks after Corvint 1.0

- Date: 2026-09-24
- Status: accepted (owner direction; the rows it adds are not started)
- Owner: Russell Lewis

## Owner direction (2026-09-24, exact)

- "we should be looking to move the beamfall roadmap to the convint-task manager"
- "we should make sure we add them then" (the verbs corvint-tasks lacks for that move)
- "once we hit 1.0 we are going to go full into dogfooding corvint on beamfall and using it to its
  fullest"

## Context

Beamfall keeps 2501 tickets as checklist items in `docs/plans/roadmap/*.md` (44 shards; 2137
done, 323 open, 42 held in place, 16 retired). `script/roadmap_lib.py` parses them, and
`script/roadmap.sh` runs claim, lease, worktree, gate, review, rebase, fast-forward merge under a
per-repository lock, and a checkoff commit that rewrites the shard. Dashboard, fanout planning,
fleet dispatch (Beamfall ADR-0204) and several Beamfall gates read the shards through
`roadmap_lib.parse`.

corvint-tasks today stores tickets, dependencies, holds, approvals, policy and releases through
the §5.2 writer. `TM-V0-026` excludes a Beamfall write adapter and cross-repository scheduling,
and TCP-02's administrative verbs, TCP-04's runtime and TCP-06's `import` are not built. TCP-09
had no owner (SPEC §10 U2).

## Decision

1. Beamfall's roadmap moves to corvint-tasks: first the ticket data, then the claim-to-merge
   loop. The owner owns TCP-09, which closes SPEC §10 U2.
2. The move starts after Corvint 1.0 stable is accepted (Corvint V1-0021). Until then these rows
   stay not-started, and nothing here competes with the Corvint 1.0 critical path or its
   quiet-host gate runs.
3. `TM-V0-026` no longer excludes, as a matter of scope, a Beamfall adapter or scheduling across
   the repositories of one Beamfall workspace. Both stay unavailable, and are reported so, until
   the rows below are delivered and qualified. The other `TM-V0-026` exclusions are unchanged.
4. Rows TCP-10 to TCP-14 add what the existing rows do not cover. The order is TCP-02 (rest),
   TCP-04, TCP-05, TCP-03, TCP-10, TCP-06 (Beamfall import dry run), TCP-11, TCP-12, TCP-14,
   TCP-07, TCP-13, then TCP-09 (cutover).
5. Beamfall's Markdown shards stay the single source of truth until the TCP-06 authority switch
   runs for Beamfall with an owner cutover record. There are no dual writes. After the switch the
   shards are a generated projection that Beamfall readers can still parse (TCP-13).
6. GP (`TM-V0-027`) binds every row. Nothing here changes Corvint or Beamfall source; Beamfall's
   own side of the cutover (its ADR-0204 dispatcher and `roadmap.sh`) needs a Beamfall record
   under the same owner.

## Rows added

- TCP-10: Beamfall ticket mapping and import at scale. Map every shard field (`Depends`,
  `Primary Repo`/`Repo`, `Runtime`, `Fanout Lane`, `Size`, `Work Points`, `Review`, `Priority`,
  `Context Class`, `Complexity`, `Launch-blocking`, `Allowed Paths`, `Forbidden Paths`, `Do`,
  `Files`, `Verify`, `ADR`, `Mockup`, `AtomicCrossRepo`, `Held`) to the native schema or a
  recorded extension, and each status (`[ ]` OPEN, `[x]` COMPLETED with provenance, `[~]` HELD,
  `[r]` ARCHIVED). A field with no mapping is a refusal, never a silent drop. Store reads and
  writes at 2501 tickets are measured against GP.
- TCP-11: Integration step. Rebase onto the target branch, rerun the gate on the result,
  fast-forward merge under a per-repository lock, and record completion in the store, matching
  `roadmap.sh closeout --merge`.
- TCP-12: Scheduling across the repositories of one workspace: per-repository leases and merge
  locks, and `AtomicCrossRepo` tickets that land together or not at all.
- TCP-13: Projection for existing Beamfall readers: shards that `roadmap_lib.parse` accepts
  byte-for-byte on the same ticket set, so the dashboard, fanout planning, fleet dispatch and
  Beamfall gates keep working unchanged through cutover.
- TCP-14: Per-stage authority (claim, review, merge; Beamfall ADR-0204 D6) as policy, not flag
  files.

## Rollback

Before cutover, delete the rows and this decision; nothing else changed. After cutover, the
shards in Git history at the cutover commit are the recovery point: revert the authority switch,
and Beamfall's `roadmap.sh` resumes from them.

## Residual risks

- The ATCP source document is still missing (SPEC §0); this decision rests on owner authority, not
  on an amended ATCP text.
- Beamfall's `[~]` state mixes held and in-flight tickets; TCP-10 must tell them apart from the
  `.roadmap/claims` leases, or refuse.
- Beamfall's ticket bodies are free text written by agents; the import must secret-screen them.
