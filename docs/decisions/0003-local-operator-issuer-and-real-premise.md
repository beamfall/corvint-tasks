# 0003 — Local-operator issuer and an observed real premise

Status: accepted delivery scope, 2026-09-07. Amends decision 0002.

The owner instructed "build TCP-02 so the store works" and answered the two J3-05 questions
that `docs/SPEC.md` line 1367 left held.

## 1. Issuer (J3-05 "real issuer/authorization")

Selected: **the local operator is the authority for their own local store, and the binding is
recorded unauthenticated.**

`atm` writing this repository's own state dir is a single-operator local act, not admission of
work to a shared queue. The process that a person starts in their own terminal, inside a
repository they already have write access to, is that store's authority. `atm` performs no
authentication and therefore claims none: the receipt records the binding with
`ActorAuthentication: NOT_OBSERVED`, so an unproven fact stays visible as unproven (invariant 5)
instead of being laundered into a grant.

This amends decision 0002, which said a same-UID process is not an actor grant. That sentence
stands for **real multi-agent queues**, where a second agent under the same UID must not inherit
the operator's authority. It does not stand for a local store's own genesis and single-operator
mutation, which is what this decision releases. Nothing here qualifies an issuer for a shared or
foreign queue; `atm` mints no credential, stores no secret, and opens no network path
(invariant 1). Real-queue admission still requires the §7.4 cutover record and a `QUALIFICATION`
receipt (invariant 7), which this decision does not grant.

Rollback: revert this decision and the premise below; the fixture path is untouched by both.

## 2. Premise (the planner's hypothetical gate)

Selected: **add an observed real premise beside the fixture one.**

`internal/transaction` refused every premise except `HYPOTHETICAL_FIXTURE_NO_RUNTIME`, so no plan
it produced could be applied to a real store. A second premise, `OBSERVED_LOCAL_OPERATOR`, is
admitted alongside it. The fixture premise keeps its exact meaning and every reviewed fixture
test keeps passing unchanged, which is why this shape was chosen over reworking the gate into
per-axis coverage inputs.

A premise is a statement about the runtime the caller observed, never a grant. The `Coverage`
axes stay `NOT_OBSERVED` until each one is genuinely measured: admitting the real premise does
not make `Durability`, `RuntimeQualification`, `ActorAuthentication` or
`AdministrativeAuthorization` observed, and no code may read the premise as if it did.

## 3. What stays held

Existing-intent CAS/C6 against hostile concurrent editors, capacity qualification, runtime and
escrow, real foreign queues, `archive restore`, GP, and the §7.4 cutover record all remain held
behind their own contracts and evidence. This decision releases the local writer and `atm init`
only.
