# Parent steering — completed-stage current-byte binding

2026-09-06. Independent bounded Gate A found and confirmed a MED linkage gap in the draft `StageObservation.Bind`. See the appended correction in `tcp02-j3-02-gate-a.md`. The parent selects the following narrow repair within existing ownership.

For completed-head leftovers, compare descriptor operation to actual receipt kind using the accepted model's mapping: INIT->INIT, PAUSE->PAUSE, UNPAUSE->UNPAUSE, KEEP_JOURNAL/ADOPT_FILE->RECONCILE. Also require descriptor RecordedAt to equal the actual receipt RecordedAt. Mismatch returns JOURNAL_FORKED. Retain existing actual head/receipt/hash/sequence/prev/request/outcome proof. This is current-byte field consistency, not full-plan reconstruction, actor authority or a claim to distinguish historical KEEP versus ADOPT acceptance.

Concrete wrong-label witness: construct a codec-valid three-slot UNPAUSE descriptor referencing an actual completed PAUSE receipt/head/request (receipt is below UNPAUSE cap), carrying that request's actual digest. Existing hash/sequence/request checks alone accept this mismatched operation. Add a negative for that case and a descriptor-only timestamp change, plus positive mapping witnesses including RECONCILE for both reconciliation labels. Do not weaken codec or enlarge bounds; no new API or ownership required.

## Second selected correction — headless INIT before physical queue publication

Independent Gate A confirmed a MED regression in the draft. `records.go` calls queue-dependent `validateStage` before the receipt walk; that helper requires physical `intent/queue.json`. After a real INIT receipt link, that file may still be absent. Existing J1 accepts its null PRE while fully validating the queue afterimage from receipt1. The current genesis witness merely removes head from a fully materialized fixture and misses this prefix.

Defer the queue-dependent stage binding for headless linked genesis until the existing complete genesis chain/blob/request/projection audit succeeds. If physical queue is absent, use the queue afterimage obtained from that actual fully validated linked genesis (bounded and content-bound); NEVER use a staged proposed afterimage, descriptor assertion or cached fixture bytes. Reuse the existing audit/post decoder and preserve all its checks. An actual non-not-exist read failure stays a refusal. No public canonical-observation API or Source/Reader contract widening is needed.

Add causal-red/green native journal coverage for a complete linked INIT receipt and required immutable evidence, active INIT descriptor, no head and no physical queue (prefer no materialized posts). It must report REDO_PENDING only after the complete existing audit passes. Cover missing/corrupt required genesis blob, invalid queue/genesis, foreign/non-fixture binding and read-error refusal. No linked receipt remains UNINITIALIZED; staging alone never initializes. Preserve byte/mode/mtime purity. Archive remains UNINITIALIZED when head is absent under its existing snapshot probe.

See the second appended J3-02 Gate A addendum for exact rationale. This is selected within your existing source ownership and supersedes any earlier report that only covered fully materialized genesis. Re-read this steering before finalizing.
