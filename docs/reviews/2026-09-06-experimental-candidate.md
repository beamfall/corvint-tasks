# Experimental candidate review — 2026-09-06

Status: REPAIR REQUIRED. This is operator-recorded evidence, not an accepted upstream contract or a delivered runtime. Claude authored the implementation and repairs; Codex applied mechanical gofmt and performed independent review. Claude's next read review exited 1 because its monthly spending limit was reached. The owner has been asked whether to continue implementation with Codex or wait for Claude; no answer is recorded.

## Candidate and execution evidence

The complete candidate is frozen in `/private/tmp/corvint-tasks-build-20260906/combined-review-snapshot`; `SNAPSHOT.json` records each file SHA-256. Source bytes were checked against successful Claude Write/Edit output and recorded mechanical formatting before capture. The six generated superseded memory stubs were removed only after hash checks. No runtime source changed after capture.

The read subset's `make verify` passed Go 1.27.0, formatting, tests and vet (`tcp01-read-snapshot-verify.log` in the same build directory). It excluded mutation and authority. The combined candidate's `make verify` FAILED (`combined-verify.log`, exit 2): `TestTMV0007_AS35_ProtectedFieldsAndIdentity` expected ADOPT_UNSUPPORTED_FIELD, but its OPEN record with non-null completion is rejected earlier as MALFORMED. All other package test commands passed; combined vet was NOT_RUN because make stopped at the failure. A passing host test that accepts UNSUPPORTED_FILESYSTEM does not qualify host durability.

Fresh independent Codex reviewers inspected the frozen candidate. They ran no tests and made no edits. The following are confirmed source findings, not claims of executed race witnesses.

## Open repairs

| ID | Severity | Evidence | Counterexample and required repair |
|---|---|---|---|
| A1 | HIGH | internal/authority/fs_darwin.go:37; authority.go prerequisite | File.Sync silently falls back to fsync on ENOTSUP in pinned Go 1.27.0 internal/poll/fd_fsync_darwin.go:27. Issue explicit F_FULLFSYNC under a pinned descriptor, retry EINTR only, preserve unsupported failure; keep directory fsync separate. Test below the actual Darwin wrapper and forbid publication on failure. |
| A2 | HIGH | internal/authority/fs_linux.go:22,40 | ext2/ext3/ext4 share magic 0xEF53; naming all ext4 admits unapproved filesystems. Shared magic alone must remain ambiguous/refused unless descriptor-bound mount evidence distinguishes ext4. Test ext2/ext3 refusal and qualified ext4. |
| A3 | HIGH | internal/authority/publish.go:49,81; lock.go:92; qualify.go:68 | Separate Lstat/OpenRoot permits directory-to-symlink replacement and redirected publication; OpenRoot can block on a replacement FIFO before checking type. Use atomic no-follow, directory-only descriptor traversal and retained identity. Add controlled swap/no-redirect/no-block witnesses. |
| A4 | MED | internal/authority/lock.go:148 | EINTR immediately retries without context/deadline checks; cancellation during setup can still return a acquired lock. Check on every retry and after acquisition, closing on cancellation. Test interrupted deadlines and subsequent acquisition. |
| A5 | MED | internal/authority/fs_other.go:12; lock.go:112; publish.go:193 | Unsupported platforms create a lock or publication temp before refusal. Refuse at entry before filesystem effects and test identical before/after trees on an unsupported target. |
| M1 | MED | internal/mutation/apply.go:209; replay.go:37 | Requests = (*MemoryIndex)(nil) passes interface nil checks and Lookup returns absent, permitting a fresh plan without an index. Reject typed-nil implementations at both entry points; test Apply and Adopt. |
| R1 | MED | internal/cli/cli.go:746,811 | gateShow retains a removed gate from a failed snapshot attempt; ticketShow retains notFound after a ticket appears on retry. Reset callback-owned state at each invocation; add deterministic retry witnesses. |
| R2 | MED | internal/archive/verify.go:134; layout.go:148; archive_test.go:29 | Manifest digest is checked but evidence/<digest> and pinned/<digest>.json basename identities are not. Existing successful fixtures use unrelated names. Validate closed lowercase Digest names against consumed-byte hashes; correct positives and keep negative mismatches. |
| R3 | MED | internal/archive/verify.go:155 | A reader returning (0,EIO) after valid tar end blocks is accepted because the final read error is discarded. Require n=0 and io.EOF; propagate all other errors and test terminal failure. |
| R4 | MED | internal/archive/export.go:20,270; docs/SPEC.md:450 | A destination accepting a prefix then failing makes the unconditional zero-output-on-failure claim impossible. Scope zero output to pre-delivery validation/staging failure; document partial delivery/interruption and test a partial-write error. Do not try stdout rollback. |
| T1 | MED | internal/mutation/adopt_test.go:558 | Combined gate failure above. Separate invalid-record MALFORMED coverage from a structurally valid protected-completion adoption case; preserve the protected-field assertion rather than deleting coverage or weakening validation. |

Mutation AS-02 passed inspection; AS-03 fails M1; AS-35 otherwise passed inspected repaired cases subject to M1 and T1. Read AS-01/04/06/07 and inspected AS-10 bounds passed inspection; AS-08/09/36 are partial due to R1..R4. These labels are review judgments, not qualification.

## Hypotheses and future holds

Read-side FIFO replacement between stat and open in internal/intent/worktree.go:185 and store.go:301 still needs a focused deterministic witness. Do not promote it to an executed finding. ADOPT prose in SPEC:419-420 should distinguish changes entering/leaving COMPLETED from otherwise permitted unchanged COMPLETED records.

TCP-02 has authority primitives only, not a journal, aggregate capacity escrow, reservations, process runtime, restore or CLI mutation wiring. All durable producers remain held on the whole-store and terminal/recovery/admin reserve proof. Worst-case flat-manifest RSS, GP no-slowdown, cutover and runtime qualification are NOT_RUN. The owner no-slowdown requirement remains binding. Agent presence and targeted mailbox designs are proposed separately in Corvint branch codex/task-control-plane-proposal-20260906; they are not implemented here.

## Resume sequence

Resolve the builder choice; repair these bounded findings with witnesses; mechanically format; run focused regressions and one combined canonical gate; obtain a fresh independent review of repaired bytes. Then continue the journal/capacity slice and later runtime tickets subject to their explicit gates. Do not tick TCP-01/TCP-02 complete, run real queues, publish, or claim no slowdown from this report.
