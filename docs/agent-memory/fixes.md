---
name: fixes
description: Repairs that are decided but not yet applied, or applied but not yet verified by root
updated: 2026-09-06
---

# Fixes

Pending work only: an entry is a repair someone still has to apply, or one applied in the
tree that a recorded `make verify` run has not yet classified. It leaves when the run shows
the witness passing or when the owner withdraws it.

### Callable writer prerequisites
- J2b/J3-01 and J3-02 read support passed their parent gates and fresh reviews. J3-03/04
  test-owned transaction/recovery/native-read integration passed its combined gate/review.
  A callable writer still needs independently qualified issuer/actor/admin binding, observed
  full-store capacity/reserves, hostile-editor intent CAS and production durability evidence.
  Test-owned inventory/cleanup tokens and native reads grant none of those authorities.
- Owner: TCP-02 J3-05/TCP-02b.

### Validated canonical-record retrieval during intent divergence
- Strict J1 Audit(selected ticket) refuses on D != canonical C. The request adapter now
  permits stable intent divergence while validating private projections, but returns only
  request history. Before reconciliation wiring, add a validated canonical-record retrieval
  entrypoint; the J3-03/04 finite test-only actual-byte helper does not supply a production API.
  Never use arbitrary partial audit results as authority. Parent steering §5;
  outside J1, owner TCP-02 J3/TCP-02b.

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
