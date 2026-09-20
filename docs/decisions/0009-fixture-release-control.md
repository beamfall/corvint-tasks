# Decision 0009: fixture release control

- Date: 2026-09-20
- Status: accepted

Corvint Tasks may durably model ordered releases in fixture queues. A candidate binds its release definition, explicit ticket records and acceptance revisions, predecessor promotions, current policy, observed Git `HEAD` and tree, and a digest of repository content excluding `.taskman`.

Manual and external observations are `taskman-release-attestation/0`, never native `taskman-gate-result/0` execution. Promotion is an immutable local journal fact; it does not publish, tag, deploy, cut over, or authorize a real queue. Native gate execution and non-fixture release mutation remain `NOT_RUN` and refused. Release divergence freezes mutation; `KEEP_JOURNAL` retains discarded bytes and restores the canonical projection, while release `ADOPT_FILE` remains refused.
