# 0008 — Settled fixture pause and unpause

Status: delivery scope under the owner's 2026-09-19 “get it there” instruction. Extends the
callable fixture writer using the already specified PAUSE/UNPAUSE templates and digests of
SPEC §5.5. No new administrative policy or authenticated role is introduced.

Expose `pause|unpause --request-id ID [--role OWNER|OPERATOR]` for settled fixture queues. PAUSE
creates only ADMISSION/OPERATOR; it does not stop ordinary ticket mutations. UNPAUSE can remove
ADMISSION or ALL. Existing identical pause or absent unpause is the model's unrecorded NoChange.
Exact recorded retries return the original outcome before fresh transition checks, so retrying
an old PAUSE after UNPAUSE cannot pause again. Actor, operation and queue remain digest-bound.

UNPAUSE preserves unrelated edited ticket bytes, including empty or malformed D, while validating
the complete canonical ticket inventory, queue, policy and private journal. Every canonical ticket
still requires a present regular physical projection; extra or missing tickets refuse. Normal
audit and other mutations retain strict projection equality. No branch condition is added to the
existing PAUSE/UNPAUSE model, which changes no intent projection.

Barrier removal delegates to the existing pinned, exact-pre-digest unlink and directory sync
primitive. Its only destination is barrier.json. The validated UNPAUSE receipt commits first,
ordinary posts follow, barrier deletion follows them, and head publishes last. A deletion failure
never advances head or erases a changed barrier. No generic deletion API or history deletion is
added. The original capacity model, including reserved UNPAUSE capacity, remains unchanged.

The new entrypoint refuses pending receipts and fresh work with active staging. A postcommit
UNPAUSE failure can leave a pending receipt that **no current callable recovery path can finish**:
the existing ticket redo refuses null deletion posts. Evidence and projections remain retained;
operators must not manually erase them. Completing this recovery remains a later TCP-02 task.
Returned-fault ordering and boundary-edit tests do not qualify process crashes, power loss or
hostile-editor CAS.

No runtime, reservation release, real-queue operation, foreign write, performance or qualification
authority is granted. Evidence: [writer repair](../reviews/2026-09-19-writer-repair.md).
