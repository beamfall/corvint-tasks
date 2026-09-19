```text
REPAIR — one MED static-review hypothesis; not execution-proven. No HIGH identified.

Acceptance dispositions:
PASS: closed codecs/byte identities; minimal noncircular genesis; complete receipt inventory (+3, gaps, pending/no-head); canonical latest/projection checks.
PASS: fresh full request verification, divergence-only intent exception, corrupt/missing request refusal even for absent IDs, lookup-error refusal without fresh Apply/Adopt plans.
PASS: exact outcome/revision binding, admin null revisions, uncommitted refusal null sequence; physical preconditions and structural RECONCILE; effect-free UNPAUSE deletion classifier.
PASS: bounded enumeration, streamed history, selected-output limits, snapshot/body-error retries and read purity.
REPAIR: zero-byte raw evidence handling.

MED — internal/journal/audit.go:141, reached from internal/journal/records.go:211
Input: valid genesis followed by a correctly linked receipt posting evidence/<h>, where h=SHA256(empty bytes), pre=null, sha256=blobSha256=h, record=null, and the retained evidence file exists with length zero. Audit—and consequently every request lookup—returns MALFORMED before digest verification.

Requirement: plan R3’s consumed-blob validation and SPEC §1/§3.4 admit raw evidence with an upper byte bound; no positive minimum is specified. Empty captured output is raw evidence, distinct from an empty JSON document.

Smallest correction: allow present zero-byte evidence while preserving an explicit distinction from missing files/deletions; retain empty structured-document refusal. Add focused witnesses for empty evidence success and missing evidence refusal.

Evidence accepted: parent make verify PASS (Go 1.27.0, format, all tests, vet), exit 0/cleanup, 127 files including 74 Go sources unchanged. No tests rerun or files changed.

Scope remains explicit: historical actor/acceptance/liveness/runtime semantics NOT_OBSERVED; unknown later semantic coverage UNKNOWN; archive historical semantic integration deferred for unavailable retained bytes; validated canonical retrieval during divergence remains a future reconciliation hold. These are not findings. The modest heap witness does not qualify production RSS; GP NOT_RUN.

This exact candidate should not be accepted/saved as experimental support until the MED is resolved and reviewed. No production, real-queue or TCP-02-complete qualification.
```
