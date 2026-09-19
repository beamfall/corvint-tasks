# 0004 — One MUTATE transaction operation, and the redo rule completed

Status: accepted delivery scope, 2026-09-07. Extends decision 0003; TCP-02b.

The owner instructed "build TCP-02b so tickets can be created". TCP-02b is defined in
`docs/ROADMAP.md` as "serial wiring of ticket mutations to the journal; AS-02/AS-03/AS-05 rerun
against the real journal", and `internal/mutation/apply.go` names its counterparty exactly: a
`Plan` "says what the post record would be; the transaction writer of TCP-02b decides whether it
becomes real."

## 1. One operation, not fourteen

`transaction.Model` had five closed operations — INIT, PAUSE, UNPAUSE, KEEP_JOURNAL, ADOPT_FILE —
and none of them can create a ticket. ADOPT_FILE is the closest, and it is explicitly not it: the
offered file must carry a `ticketId` equal to an existing canonical record's, and §3.3 forbids it
from completing, archiving, reopening, restoring, approving or re-sourcing. A ticket mutation had
no journal representation at all.

Selected: **one `MUTATE` operation carrying one `taskman-mutation/0` envelope**, rather than a
transaction operation per §3.3 verb.

`mutation.Apply` already validates and computes the post record for all fourteen verbs, and the
staged artifact shape is identical for every one of them: the ticket projection, the request
index entry, the receipt, the head, and — for a CREATE that allocated a serial — the queue
manifest. Fourteen operations would multiply the §5.2 closed preimage, template and slot tables
by fourteen without adding a single check that the envelope codec does not already make.

The receipt keeps its own §3.1 vocabulary, which already anticipated this: `ARCHIVE` and
`RESTORE` receipts carry those kinds, and every other mutation is `MUTATION`. The stage operation
name is `MUTATE` for all of them because the staged shape does not vary.

The TM-V0-006 request digest of a mutation is the SHA-256 of the envelope bytes exactly, as
TM-V0-006 already specifies. `issuedAt` is one of those bytes, so a retry replays only when it
reproduces the same envelope; `atm` exposes `--issued-at` so a caller can. This is the same rule
ADOPT_FILE states, applied to the envelope that is itself the request.

## 2. The library index stays absent inside the model

SPEC §3.4 says TCP-02b "must supply the `requests/` index of this directory when it wires the
library to the journal". It is supplied — as the model's `ReplayObservation`, read from
`requests/<xx>/<sha256>.json` under the same lock as every other input.

It is not also handed to `mutation.Apply` as a second `RequestIndex`. The model decides replay
before it computes anything, binding the receipt sequence and the stored outcome shape as well as
the digest, which is strictly more than the library's own lookup checks. A second consultation
inside the same transaction could only agree redundantly or disagree, and disagreement has no
defined resolution. The reviewed ADOPT_FILE path already passes `absentIndex{}` for this reason;
`MUTATE` is consistent with it. `TestTMV0005_ContextInputsRequired` and
`TestTMV0006_AS03_TypedNilIndexRequired` continue to hold the library's own contract.

## 3. Three defects a real mutation exposed

The genesis transaction of TCP-02 wrote only to destinations that were absent, so three gaps in
the writer were unreachable until a second transaction ran.

**The redo rule was missing its middle case.** §5.2 admits a destination holding the post digest
"or the `pre` digest of a receipt that is still redo-pending". The delivered rule handled the
post digest and absence, and refused everything else — so an ordinary CREATE, whose `queue.json`
destination legitimately holds the pre state, was refused `INTENT_DIVERGED`. The pre digest is
now read from the committed receipt's own `pre` entries, which is the same evidence a crash
recovery reads, so a first write and a redo of it decide identically.

**Post-commit redo did not exist.** §5.2 places the receipt link-in before any post file, so a
failure in the post phase leaves a committed receipt, unwritten projections and a stale head
(crash point C2). Without redo the store is wedged: the next transaction sees a receipt count
that disagrees with the head and refuses. `internal/store/redo.go` completes a pending receipt
before any new transaction is modelled, using the receipt as the self-contained redo record §3.4
says it is.

**Intent projections could not be replaced.** The reviewed applier admitted replacement for the
head, barrier and reservations only, and its `expected == nil` path degrades to an exclusive
create — so `head.json` could be written exactly once, and a ticket or queue projection could
never be rewritten. The applier now also admits the intent roles, and admits them only under an
explicit expected digest, so a projection is replaced by a real compare-and-swap on the bytes the
transaction read and never by an unconditional overwrite. That is strictly stricter than what the
three original roles are held to.

## 4. What this does not decide

Held exactly as decision 0003 left them: the runtime, reservations and process liveness, the
staging-descriptor crash recovery of §5.6, `archive restore`, `reconcile intent`, the
administrative verbs (`admit`, `cancel`, `retry`, `resume`, `pause`, `unpause`, `drain`,
`cutover`, `import`), and GP.

A non-fixture queue still admits nothing. The queue this writer serves declares `fixture: true`,
and a real queue continues to require the §7.4 cutover record and a `QUALIFICATION` receipt.
Every Coverage axis remains `NOT_OBSERVED`: nothing here measures durability, authenticates an
actor, or qualifies a runtime.

Attempt liveness is answered by the zero-attempt oracle, exactly as the reviewed ADOPT_FILE path
answers it. That is sound only while no runtime exists and therefore no attempt can be live; the
runtime slice must replace it with a real oracle before it starts one.
