# Repair R2 linked-worktree staging closure — 2026-09-06

Status: bounded repair and focused archive regressions PASS; parent canonical gate and fresh
independent review PENDING.

`internal/archive/export.go` now resolves the proposed staging directory's enclosing Git
authority before pinning or creating a staging file. `UNINITIALIZED` is the only accepted
resolution error. A resolved authority sharing the queue repository's `CommonDir`, including a
linked worktree, is refused; malformed, unreadable or otherwise uncertain metadata is refused.
The check performs no worktree enumeration, project scan, write, lock or subprocess.

`TestTMV0022_AS36_N4fStagingRefused` constructs the existing resolver's linked-worktree metadata
grammar and proves refusal with zero stdout and byte-identical primary and linked trees. It also
proves malformed authority refuses unchanged and ordinary outside plus default OS-temp staging
still succeed and clean up their immediately unlinked stage file.

Focused command (Go 1.27.0 local toolchain, shared bounded cache, two-way package parallelism):

```sh
GOTOOLCHAIN=local GOCACHE=/private/tmp/corvint-tasks-build-20260906/go-cache GOMAXPROCS=2 GOFLAGS=-p=2 \
  go test -count=1 ./internal/archive \
  -run 'TestTMV0022_AS36_N4fStagingRefused|TestTMV0022_AS09_ExportThenVerify'
```

Result: PASS (`ok corvint-tasks/internal/archive 0.265s`). No canonical gate was run. Prior review
and evidence reports remain unchanged. The complete package also passed with the same environment:
`go test -count=1 ./internal/archive` (`ok corvint-tasks/internal/archive 12.296s`). This repair does
not enumerate or pre-discover worktrees,
prevent an OS owner relocating a retained staging directory, qualify native Linux/RSS/durability,
or authorize real queues.
