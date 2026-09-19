# Fixture-only first delivery

Owner decision [0002](../decisions/0002-fixture-only-first-delivery.md) is delivered experimentally. TCP-01 support, TCP-02 authority/J1 reads, J2a/J2b sizing and transaction models, J3-01 private primitives, J3-02 native staging reads, and J3-03/04 test-owned transaction/recovery integration are preserved on the delivery branch. TCP-02 as a whole is incomplete.

## Verification and independent acceptance

The final J3-03/04 candidate is based on reviewed J3-02 `bc3b3dd7df3d383f0d546cc9e6ffe7e46a3d0d28`. One parent-owned `make -j1 verify` passed under Go 1.27.0: formatting, all tests and vet. All 211 pre-gate hashes (103 Go files) matched after the gate and fresh independent final review. All 99 pre-existing Go files are byte-identical; four new authority test files provide the fixture implementation. The authority package took 129.640 seconds. See [gate log](2026-09-06-j3-fixture-final-verify.log) and [frozen manifest](2026-09-06-j3-fixture-final-frozen-sha256.json).

The [fresh independent review](2026-09-06-j3-fixture-final-review.md) returned PASS with no HIGH/MED finding and scored every fixture acceptance criterion and all five deferred native integration groups PASS. The selected stale-inventory and branch-token repairs have isolated causal-red witnesses plus restored-source green tests. The focused run passed 27 suites / 393 named tests and subtests, with no failures or skips; final binding/accounting follow-ups passed before the canonical gate. [Builder evidence](2026-09-06-j3-fixture-integration-builder.md), [accepted plan](2026-09-06-j3-fixture-integration-plan.md), and [Gate A](2026-09-06-j3-fixture-gate-a.md) preserve the full mapping and limitations. Parent acceptance supersedes their historical pending-parent wording.

After review, the parent only added evidence, updated delivery-status text, and aligned crash-matrix C1–C3 terminology with the already-governing section5.2 non-overwriting receipt link and exact PRE/POST recovery rules. No Go, Makefile or module input changed after the gate/review. No new behavior or authority was accepted by that editorial alignment. A [fresh final editorial check](2026-09-06-j3-fixture-editorial-final.md) passed after tightening C1 cleanup ownership and C3 absent-PRE wording; all Go/build hashes still match.

## Exact boundary

Delivered tests use disposable finite fixtures: at most twelve receipts, 128 entries per directory and 4 MiB captured bytes. In INIT/PAUSE/UNPAUSE/KEEP/ADOPT order they cover 23/13/12/15/13 publication stops, 22/19/20/21/20 distinct first-occurrence boundary faults, 11/5/4/5/4 cleanup-sync cuts, plus fresh recovery, missing-slot, rename-sync and release cases. This is not every syscall occurrence or production-cap saturation.

Callable durable writer, independently qualified issuer/actor/admin binding, hostile-editor CAS, CLI mutations, runtime/reservation escrow, real queues, restore and GP remain held. Native Linux and physical saturation are NOT_RUN; process-interruption and power-loss evidence are NOT_PRODUCED. Archive verification remains chain/digest integrity, not historical semantic acceptance or restore qualification. Corvint Go source was not edited by this build; no measured no-slowdown claim is made and GP remains NOT_RUN.

No default-branch merge, push, deployment or real-queue cutover is included. The initial Corvint proposal is separately preserved at `802830a` on `codex/task-control-plane-proposal-20260906`, outside the Corvint main branch.
