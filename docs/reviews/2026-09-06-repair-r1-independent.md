**REPAIR.** This exact candidate should not yet be applied as experimental support. Repair the findings below, complete parent-owned verification, then obtain another independent review.

| Item | Disposition |
|---|---|
| A1 | PASS — explicit pinned F_FULLFSYNC, EINTR-only retry, no ENOTSUP fallback; directory fsync separate. |
| A2 | PASS — shared ext magic remains ambiguous/refused; ext4 policy does not imply qualification. |
| A3 | PASS — Darwin syscall 463 argument layout, descriptor lifetime, component no-follow opens, bridge identity/refusal and retained publication roots close the stated counterexamples, assuming trusted OS `/dev`. |
| A4 | PASS — cancellation/deadline checked across flock retries and acquisition; refusal closes the descriptor. |
| A5 | PASS — unsupported mutation entries refuse before creation, including stale LinkIn handles. |
| M1 | PASS — typed-nil indexes refused before Apply/Adopt replay or planning. |
| R1 | PASS — callback output and absence state reset on every attempt. |
| R2 | PASS — exact lowercase digest grammar and consumed-content identity enforced. |
| R3 | REPAIR — zero-byte terminal EIO handled; bytes-plus-EIO remains lossy below tar decoding. |
| R4 | PASS — partial delivery truthfully reports accepted bytes; pre-delivery zero-output guarantee correctly scoped. |
| R5 | PASS — encoded-stream cap includes manifest, PAX, padding and trailer; public Bytes remains payload-only. |
| T1 | PASS — valid COMPLETED protected-field witness; malformed OPEN coverage retained. |
| FIFO | REPAIR — production blocking race closed; controlled legacy witness is meaningful, but setup-error synchronization can hang. |

Three actionable findings, derived from source inspection; these counterexamples were not executed during this review:

1. **HIGH — staging validation is separated from pathname reuse.** [export.go:213](/private/tmp/corvint-tasks-build-20260906/codex-repair-r1-review-snapshot/internal/archive/export.go:213). After `checkStaging` accepts an outside directory, another ordinary process renames it and replaces its pathname with a symlink to an existing repository directory. `CreateTemp` follows that replacement and writes inside the repository, violating **TM-V0-008 / TM-V0-022 / AS-36 N4f**. Cleanup can erase the evidence afterward. This predates R1 but remains present. Smallest correction: retain the validated directory handle and create/remove staging entries through it; retain the staging file descriptor through verification and delivery. Add a deterministic post-validation swap witness.

2. **MED — an error accompanying a complete trailer block can disappear.** [stream_limit.go:56](/private/tmp/corvint-tasks-build-20260906/codex-repair-r1-review-snapshot/internal/archive/stream_limit.go:56), [verify.go:175](/private/tmp/corvint-tasks-build-20260906/codex-repair-r1-review-snapshot/internal/archive/verify.go:175). A valid-stream reader returns the final 512 trailer bytes together with EIO, then EOF. Tar’s full-buffer read can discard EIO; the wrapper remembers only bytes, and the terminal probe sees EOF, permitting success. This leaves **R3 / TM-V0-022** incomplete. Preserve non-EOF source errors independently of decoder returns and refuse before success. Add this bytes-plus-error case alongside the existing zero-byte EIO witness.

3. **MED — FIFO witness can wait forever before its timeout starts.** [fifo_test.go:44](/private/tmp/corvint-tasks-build-20260906/codex-repair-r1-review-snapshot/internal/intent/fifo_test.go:44). If `Mkfifo` returns ENOSPC, the worker sends an error to `done` without closing `started`; the test blocks indefinitely on `started`, and hook restoration is never reached. This violates the requested bounded, cleanup-safe witness requirement. Select readiness against completion, register hook restoration before launch, and ensure release/join cleanup covers every exit. Add a setup-failure witness.

Recorded Go 1.27/gofmt/tests/vet results support the inspected repairs; I reran nothing. Linux/Windows cross-compilation establishes compilation only. JS/wasm execution establishes the refusal witness, not runtime qualification. Private hooks expose no production input bypass; inspected repair tests launch no child processes.

Separate holds remain: journal/caller binding, aggregate capacity and recovery reserves, reservations/runtime/restore/wiring, worst-case RSS, GP, physical durability and owner-authorized real-queue qualification. This is not a complete runtime.
