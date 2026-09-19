**REPAIR — this exact candidate should not yet be applied as experimental support.**

| Area | Disposition |
|---|---|
| F1 | Descriptor retention and pathname-swap protection PASS; cleanup-error persistence needs repair. |
| F2 | PASS: sticky bytes-plus-EIO, valid bytes-plus-EOF, encoded cap and payload metrics preserved. |
| F3 | Setup/release/deadline paths bounded; parent-interruption cleanup needs repair. |
| Linked staging | Ordinary linked and malformed-metadata witnesses PASS; case-alias exclusion needs repair. |

1. **HIGH — linked-worktree exclusion compares path spellings, not directory identity.** [export.go:100](/private/tmp/corvint-tasks-build-20260906/codex-repair-r2-review-snapshot/internal/archive/export.go:100). On a case-insensitive Darwin filesystem, let the source common directory be `/…/Repo/.git` and the linked checkout’s `gitdir:` use `/…/repo/.git/worktrees/staging`, with `commondir` containing `../..`. `intent.Resolve` preserves that spelling; `filepath.Clean` does not normalize case. Both names identify the same common directory, but staging is accepted and written inside its linked worktree. Violates **TM-V0-008 / TM-V0-022 / AS-36 N4f**. Compare validated common-directory filesystem identities, refusing inspection failures; add this alias witness.

2. **MED — snapshot retries can erase failed staging cleanup.** [export.go:364](/private/tmp/corvint-tasks-build-20260906/codex-repair-r2-review-snapshot/internal/archive/export.go:364). Exclusive creation succeeds, immediate unlink returns EIO, and the source snapshot changes before re-probing. The unlink/close error becomes only the callback error; `snapshot.Reader` discards it on the moved attempt. A successful retry then delivers stdout and reports success despite the abandoned staging entry. Violates **F1’s truthful cleanup requirement / TM-V0-022**. Preserve unlink and associated close failures in export-level sticky cleanup state, independently of retryable snapshot errors, and refuse delivery. Add a one-shot unlink-error-plus-snapshot-move witness.

3. **MED — interrupting the test parent can orphan the blocked FIFO child.** [fifo_test.go:155](/private/tmp/corvint-tasks-build-20260906/codex-repair-r2-review-snapshot/internal/intent/fifo_test.go:155). After readiness, deliver SIGTERM or SIGINT only to the test-parent PID while its legacy child blocks in `os.Open`. Default signal termination bypasses deferred cleanup; the child has no independent deadline. The interruption test signals **the child**, leaving its parent alive to reap it. Violates **F3 / the supplied EXIT–INT–TERM cleanup requirement**. Install parent cancellation handling before launch, route it through bounded kill/reap, and add a parent-interruption witness proving no child survives.

Parent-supervised `make verify` **PASS**, exit code 0: all **98 frozen hashes**, including **65 Go sources**, match this snapshot’s [hash manifest](/private/tmp/corvint-tasks-build-20260906/codex-repair-r2-review-snapshot/docs/reviews/codex-repair-r2-linked-frozen-sha256.json). Counterexamples above are source-derived; no tests rerun.

Journal, capacity, runtime, platform qualification, durability, GP and real-queue work remain separate holds, not defects in this bounded repair.
