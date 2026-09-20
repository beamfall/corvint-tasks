---
name: fixes
description: Repairs that are decided but not yet applied, or applied but not yet verified by root
updated: 2026-09-19
---

# Fixes

Pending work only: an entry is a repair someone still has to apply, or one applied in the
tree that a recorded `make verify` run has not yet classified. It leaves when the run shows
the witness passing or when the owner withdraws it.

### Frozen governing source copies are not durably available
- SPEC §0 and decision 0001 cite six exact accepted files under
  `/private/tmp/corvint-tasks-build-20260906/context`; that directory was absent during the audit.
- Repair: recover exact bytes into durable storage and verify all recorded SHA-256 values.
  Current sibling Corvint documents must not silently substitute for the accepted versions.
- Opened: 2026-09-19.

### Interrupted genesis and staging recovery
- Invalid INIT inputs now refuse before state-directory creation; crashes after directory
  creation remain a separate recovery gap. Do not classify an arbitrary existing partial
  state as initialized or delete it without the contracted journal/staging ownership proof.
- Implement receipt-bound genesis redo and active-descriptor recovery with interruption
  evidence before promoting the writer beyond its present fixture premise.
- Owner: TCP-02. Opened: 2026-09-19.

### Callable writer prerequisites
- J2b/J3-01 and J3-02 read support passed their parent gates and fresh reviews. J3-03/04
  test-owned transaction/recovery/native-read integration passed its combined gate/review.
  A callable writer still needs independently qualified issuer/actor/admin binding, observed
  full-store capacity/reserves, hostile-editor intent CAS and production durability evidence.
  Test-owned inventory/cleanup tokens and native reads grant none of those authorities.
- Owner: TCP-02 J3-05/TCP-02b.

### Archive historical semantic byte source
- J1 native ledger audit requires every referenced blob. Archive verify currently discards
  those bytes from its retained semantic input, so it remains chain/digest verification only.
  Design a bounded archive byte source before wiring full ledger semantic validation; never
  infer full historical validation from a digest or unavailable blob. Owner: TCP-02/TCP-01.

<!--
### <short title>
- Applies to: <path(s)>
- Repair: <decided text | commit>
- Opened: <YYYY-MM-DD>
-->
