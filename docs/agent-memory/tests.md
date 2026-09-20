---
name: tests
description: Tests that are unwritten, written but NOT_RUN, or failing, with the ticket that owns them
updated: 2026-09-19
---

# Tests

Pending work only: every `AS-NN` row of SPEC §9.2 and every `TM-V0-NNN` has a named Go test
(`TestTMV0NNN_ASNN_<Scenario>`) that a recorded `make verify` run shows passing, or an entry
here. A test that exists but has not executed under a recorded run is NOT_RUN, never
"passing"; a scenario is not removed by deleting its test.

### Native Linux authority and safe opening execution
- State: NOT_RUN on a native Linux host. Authority and safeopen tests cross-compile for
  linux/amd64, and the pure shared-ext-magic refusal witness passed on Darwin.
- Required before claiming Linux runtime qualification: execute the no-follow/FIFO/descriptor
  and authority suite on Linux and record filesystem observations. The shared ext magic
  remains refused; no ext4 mount is qualified by classification alone.
- Owner ticket: TCP-02.

### Process-interruption and power-loss qualification
- Test-owned finite returned-fault and publication-prefix evidence passed, but no crash child
  was used. Process-interruption/power-loss evidence remains NOT_PRODUCED; real saturated-store
  recovery and production RSS/timing remain unqualified.
- Owner ticket: TCP-02 production writer prerequisites.

### Saturation and restore+unpause regressions (B4 writer prerequisite)
- Tests: unwritten; must drive a store to the file, manifest, tar, receipt and evidence
  budgets and prove refusal, recovery and export (SPEC §3.5)
- State: UNWRITTEN; blocks every durable-writer promotion
- Owner ticket: TCP-02

### AS-09 restore half, AS-10 live reservations and production durable receipt boundaries
- Test-owned finite receipt boundaries are covered by J3-03/04; production writer/runtime
  and restore qualification tests remain unwritten.
- State: UNWRITTEN for production qualification
- Owner ticket: TCP-02

<!--
### <AS-NN> / <TM-V0-NNN> — <scenario>
- Test: <TestName or "unwritten">
- State: <UNWRITTEN|NOT_RUN|FAILING|PASSING(run <date>)>
- Owner ticket: <TCP-NN>
-->
