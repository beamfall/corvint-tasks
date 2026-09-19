# 0002 — Fixture-only first delivery

Status: accepted delivery scope, 2026-09-06.

The owner selected: “Keep this delivery fixture-only” in response to the concrete
operator-invocation decision packet in the J3 fixture-writer plan, section 7.

This delivery completes experimental support, pure transaction/capacity modeling,
and test-owned filesystem/transaction/recovery mechanisms against disposable fixtures.
A fixture bit, supplied role, same UID, held lock or successful test is not an actor
or administrative grant. The existing independently qualified issuer requirement and
concurrent human-edit preservation requirement remain in force. Callable real-queue
writers, operator integration, runtime execution, cutover and performance promotion
remain held behind their existing contracts and evidence gates.

Record test evidence and limitations explicitly. This decision accepts the delivery
boundary, not an unimplemented issuer, filesystem power-loss guarantee, GP result,
upstream policy amendment, or default-branch merge.
