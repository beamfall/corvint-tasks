# TCP-02 J1 empty-evidence repair — 2026-09-06

Status: the single MED static-review finding is repaired in this isolated candidate. Focused
journal tests pass. Parent `make verify` and fresh independent review remain PENDING.

## Scope

- `internal/journal/audit.go`: byte reads now carry explicit presence independently of the
  returned slice; required present empty bytes are normalized to a non-nil zero-length slice.
- `internal/journal/records.go`: current and pending projection equality hashes present empty
  bytes instead of treating them as absence.
- `internal/journal/source.go`: the Source contract explicitly reserves `os.ErrNotExist` for
  absence and admits `nil, nil` for a present empty file.
- `internal/journal/audit_test.go`: focused witnesses cover valid committed and pending empty
  raw evidence, request-adapter audit, selected empty versus deletion, missing retained bytes,
  and malformed empty structured blobs, including a Source returning `nil, nil`.
- `docs/SPEC.md`: records the already-authorized zero-byte raw-evidence semantics and names the
  witness.

No size bound, digest check, content-addressed identity check, structured decoder, snapshot
purity rule, runtime behavior, historical semantic claim, or qualification status was relaxed.
No non-journal source/test file, real queue, foreign repository, network, Git state, or runtime
was touched. No long-running process was spawned. GP remains NOT_RUN; historical acceptance,
actor authentication, liveness, runtime qualification, and archive historical semantic
integration remain NOT_OBSERVED. TCP-02 remains incomplete.

## Focused tests

All commands used `GOTOOLCHAIN=local`,
`GOCACHE=/private/tmp/corvint-tasks-build-20260906/go-cache`, `GOMAXPROCS=2`, and
`GOFLAGS=-p=2`, without a PTY.

- `go test -count=1 ./internal/journal -run 'TestTMV0002_AS10_EmptyRawEvidencePreservesPresence'`
  — PASS after the final witness changes (`ok`, 0.314s).
- `go test -count=1 ./internal/journal` — PASS (`ok`, 1.123s).
- `gofmt -l internal/journal/audit.go internal/journal/records.go internal/journal/source.go internal/journal/audit_test.go`
  — no output.

The first new-witness invocation failed only because its fixture used an unsupported receipt
kind and did not project the deletion. Those test inputs were corrected; the repaired focused
test then passed. No full gate or vet was run; the parent owns `make verify`.
