# CLAUDE.md

Read `AGENTS.md` first; it is the repository contract. This file only adds Claude-specific notes.

- Before implementing anything, read `docs/SPEC.md` sections 1 through 7. Field names, enums,
  limits, lock order, commit points and transition rows are frozen there. Do not invent
  alternatives; propose an amendment in `docs/decisions/` if one is truly needed.
- Current truth: TCP-01 Go source and tests exist; only the root operator runs `make verify`,
  and every gate is NOT_RUN until a run is recorded. Never report a gate as green until a real
  `make verify` run shows it.
- Authority: user/owner decisions outrank the accepted Corvint/ATM/WQO invariants, which outrank
  `docs/SPEC.md` (generated local interpretation; it freezes mechanics inside that boundary),
  which outranks code. Never claim that a generated document accepts an upstream amendment or
  qualifies a real cutover.
- Pending-work memory lives in `docs/agent-memory/` as `bugs.md`, `fixes.md`, `tests.md`,
  `optimizations.md`, `ideas.md`, `questions.md`; record pending work only, never a completed
  backlog.
- Standard library only. No `go get`. No network. No daemon. No Corvint or Beamfall edits.
- Tests are named after the requirement and scenario they witness, for example
  `TestTMV0009_AS11_CrashMatrix`.
- When a task touches a wire record or a limit, update `docs/SPEC.md` in the same change and say
  so in the commit message.
- Real queues are off limits for execution until `docs/SPEC.md` §7.4 is satisfied; use a fixture
  queue (`queue.fixture: true`) for every development run.
- Performance (TM-V0-027): do not touch Corvint Go source from this repository. When a later ticket
  authorizes an Corvint-facing change, the frozen baseline and the §9.4 interleaved comparison come
  first; report GP as NOT_RUN until it has actually run, never as "no impact".
- `TM-V0-027` and gate GP are stable IDs; add new requirements as `TM-V0-028+`, never renumber.
