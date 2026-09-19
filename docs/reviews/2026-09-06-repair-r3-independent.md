**REPAIR — this exact candidate should not yet be applied as experimental support.**

| Finding | Disposition |
|---|---|
| F1 | **PASS.** Simultaneously retained no-follow roots establish common-directory identity; inspection failures refuse. The case-alias witness requires distinct spellings and actual identical filesystem identity. External/default staging remain covered. |
| F2 | **PASS.** Export-owned cleanup state preserves unlink and close failures across retries, preventing stdout. Intent/head movement witnesses inject EIO and independently EBUSY, checking both errors and zero delivery. |
| F3 | Parent interruption repair **PASS**; new outer fallback **needs repair**. Subscription precedes launch, cancellation reaches bounded kill/reap, and notification stops afterward. Parent-signal witnesses check child absence before rescue; outer-process signals exercise deferred cleanup. |

**MED — stale child PID can target an unrelated process.** [fifo_test.go:435](/private/tmp/corvint-tasks-build-20260906/codex-repair-r3-review-snapshot/internal/intent/fifo_test.go:435).

Concrete counterexample: the helper reaps its child; the outer process is descheduled; normal process churn reuses that child PID for another same-user process. On resumption, `childGone()` returns false and fallback sends that unrelated process SIGKILL. This can occur even after the earlier absence check at line 511, because absence is not retained. No root attacker or fixture tampering is required.

Requirement: F3 cleanup must target owned descendants. A published numeric PID ceases to establish ownership after reap. Smallest correction: keep destructive child cleanup in the actual helper while its child remains unreaped; make outer numeric-PID fallback diagnostic-only, using the bounded child lifetime for orphan termination. Add a stale-PID witness demonstrating that an unrelated process receives no signal.

The 15-second timer is confined to the child helper and stopped on return; expiry remains explicitly unexercised. Existing interruption witnesses are causal but do not cover PID reuse.

Parent evidence records the single canonical gate **PASS**, exit 0, cleanup complete and all 103 frozen hashes matching. Tests were not rerun; the defect above is source-derived.

Journal, capacity, runtime, GP, native-Linux qualification, durability and real-queue cutover remain separate holds, outside this bounded repair’s acceptance criteria.
