# 0001 — Build the Corvint task control plane as `atm` in this repository

**Date:** 2026-09-06  
**Status:** accepted (owner instruction), delivery not-started

## Owner instruction (exact)

After reviewing the complete 22-requirement, 10-ticket plan in Corvint
`docs/specs/corvint-task-control-plane-v0.md`, the owner instructed:

> use claude to build this

That instruction accepts building the proposed design. It is recorded verbatim; everything below
is the interpretation applied.

On 2026-09-06 the owner added, exactly:

> I want you to ensure that this functionality doesn't slow down the main go binary speed

That is a hard acceptance constraint, recorded as SPEC `TM-V0-027` and gate GP (SPEC §9.4).

## Provenance of the sources

The ATCP document was an **uncommitted working-tree file** in the Corvint checkout whose HEAD was
base `0a227a94df6ad02e35c01965d05c46be9b367905`; it was not a member of that commit, and the
document itself says only "Inspection base ... five pre-existing dirty paths were present". The
four-expert review (`corvint-task-runner-review-2026-09-06.md`) was likewise uncommitted. The
reference for each frozen source is therefore its path under
`/private/tmp/corvint-tasks-build-20260906/context/` plus its SHA-256, not commit membership:

| File | Status at base `0a227a94` | SHA-256 |
|---|---|---|
| `corvint-task-control-plane-v0.md` | uncommitted | `13a7448e86deaee8b0ba72a6f4c26f6ccba4033d75fc8eb5c9498b5607ceefb1` |
| `corvint-task-runner-review-2026-09-06.md` | uncommitted | `29a5ddd671f46678d87d614c45cdb6fc4c8dca6bad42a4983f8f1142fec8b8a3` |
| `AGENT-TASK-MANAGER-2026-09-04.md` (ATM revision 6) | as present at base | `7cca837eda3cd20fa06cf34ed3c1eb735afd79023004033791ac28912753d20a` |
| `work-queue-observation-v0.md` (WQO V0) | as present at base | `c1c05a308262cb43037f7782d633aed0a12843e4247c36fbaf29e31e9947e7d7` |
| `AGENTS.md` | as present at base | `8b536b7bded38572cdbdffe95009cffaf6d40246a7d73835b6aedcaee6d1d460` |
| `SPEC-DRIVEN-DEVELOPMENT.md` | as present at base | `0908bfdba891abc401decea9d0871c2de649a8bfeed88fe2d28fb6ae6824a054` |

The digests were computed by the coordinator on 2026-09-06 (R2) with `hashlib.sha256` over the
exact bytes of the six frozen copies under the context directory; this resolves SPEC §10 U5.
They identify the frozen copies only and imply nothing about membership of any file in an Corvint
commit; the "uncommitted" rows remain uncommitted at base. Authority stays with the
owner's sources: prose generated in this repository, including SPEC clauses that extend accepted
intent, is an implementation interpretation under the build instruction and never outranks an
unamended accepted invariant (SPEC §0).

## Decision

1. **Boundary.** The first interpretation in the ATCP proposal is selected: Corvint remains the
   task-management product, evidence and selection engine; the native-Go `atm` executor and
   store in this repository (decision 0052's accepted sibling path) is the operative side and
   owns the local ticket store, journal, reservations and receipts. No `corvint` mutation
   facade. WQO V0 stays non-operative and closed; the operative planner is the distinct profile
   `taskman-plan/0` / `taskman-priority-first/0`.
2. **Contract before code.** `docs/SPEC.md` is the executable contract (`TM-V0-001..027`) with
   traceability to `ATCP-V0-*`, `ATM-V0-*` and `WQO-V0-*`. It freezes schemas, numeric limits,
   the storage commit protocol, the transition and recovery table, acceptance tables, the
   adoption protocol and the performance gate. Implementation follows it; it is not
   implementation discretion.
3. **ATM amendments.** The amendment areas required by ATCP are frozen in SPEC §8 (A1..A7),
   preserving every ATM ID and its upstream text. They are accepted for this repository by this
   decision; the Corvint repository still needs its own amendment record (SPEC §10 U1). Where the
   SPEC does not carry an explicit row, the stricter accepted rule stands (for example
   ATCP-V0-006 whole-plan drift refusal).
5. **Performance.** The owner's 2026-09-06 constraint is binding: no existing `corvint`
   command may be slowed; storage and execution stay in `atm`; any Corvint-side change needs the
   frozen baseline and interleaved comparison of SPEC §9.4 before promotion. Until measured,
   every performance statement is `NOT_RUN` and the unchanged baseline path is preserved.
4. **Technology.** Go 1.27.0, standard library only, single `atm` binary, no daemon, no database
   service, no UI, no vendor SDK, no network fetch. Private state lives at
   `<git-common-dir>/taskman/`; intent lives in Git-tracked `.taskman/`.

## Explicitly not authorized by this decision

Publication, real-queue authority or execution cutover, autonomous dispatch of real repository
tasks, foreign-queue (Beamfall) writes, changes to Corvint or Beamfall source, and deletion or
pruning of history. Each needs its own record named in SPEC §5.4, §7.4 and §9.

## Alternatives rejected

- Operative storage and dispatch inside `corvint`: violates Corvint invariant 7 and WQO's
  authority boundary; would need an amendment of decision 0052.
- Extending Beamfall `roadmap.sh`: Beamfall-bound shell; kept only as a future adapter (TCP-09).
- Deferring schema and limit freezes to implementation: ATCP requires them in TCP-00.

## Consequences

`docs/ROADMAP.md` tracks TCP-00..09. `AGENTS.md` and `CLAUDE.md` carry the repository rules.
`go.mod` and `Makefile` bind the Go 1.27.0 gate; `make verify` fails closed until source exists,
so no gate is reported green at this revision. Repair round R1 (2026-09-06) applied the
independent Gate A review findings F1..F6 and the dropped-obligation list to `docs/SPEC.md`
without renumbering any `TM-V0` ID. Repair round R2 (2026-09-06, fresh review, verdict REPAIR)
applied findings F1..F4 and the LOW cutover item to `docs/SPEC.md` (pending-effect ownership and
supervisor identity, bootstrap windows and the single fenced-record protocol, `ADOPT_FILE`
composition, non-overwriting receipt commit with chain bounds and `archive export` as a
snapshot read, A5 cutover exemption), again without renumbering. G0 remains NOT_RUN until a
fresh reviewer re-reads it.
