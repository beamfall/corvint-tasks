# 0007 — Explicit local fixture intent reconciliation

Status: delivery scope authorized by the owner's 2026-09-19 “get it there” instruction following
the functionality/orchestration audit. Implementation interpretation of SPEC §3.1, §3.3 and
§5.5; extends the fixture writer of decisions 0003/0004 without changing their runtime boundary.

Deliver settled-fixture KEEP_JOURNAL and ADOPT_FILE through the existing pure transaction model
and receipt-first writer. A read-only `reconcile inspect <ticket-id>` supplies the validated
canonical ticket record and its record digest even when that ticket's projection is malformed
or empty. It completes the journal audit; only the named ticket projection may differ. Other
intent and private state remain strict. Partial audit results never authorize mutation.

An operator preserves the original edited bytes outside `.taskman` and supplies that copy with
an explicit choice and request ID. KEEP_JOURNAL also requires the canonical ticket-record digest
returned by inspection. Inputs remain the original choice on retry; they are never silently
refreshed from a newer projection. Exact retries return their original ticket ID without writes.
The original KEEP bytes are immutable evidence published before receipt link-in, including when
represented by a blob-backed evidence POST. ADOPT retains the existing protected-field and
revision/acceptance rules. No new reconciliation or mutation digest is introduced.

Reconciliation is permitted under ADMISSION and ALL as SPEC requires. Branch, VERSION, primary
identity, restore marker, scope, capacity and local-operator binding still apply. Pending receipts
are refused without redo or cleanup. Active staging refuses fresh reconciliation; validated
read-only replay remains available. Fresh ticket writes
and pending redo also refuse active staging until descriptor recovery is implemented.

This is not a grant to execute real tasks, authenticate shared agents, qualify durability, repair
hostile-editor races, delete history or write foreign queues. Reservations, process liveness,
runtime supervision, pending/staging reconciliation and owner cutover remain held. The six
missing governing source copies still require digest-verified recovery before upstream contract
ambiguities or qualification can be resolved; their absence does not prevent this already
specified local fixture behavior.

Evidence and named tests: [writer repair](../reviews/2026-09-19-writer-repair.md).
