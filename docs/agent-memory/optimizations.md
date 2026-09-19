---
name: optimizations
description: Performance or memory work that is required or proposed but not yet measured
updated: 2026-09-06
---

# Optimizations

Pending work only: an entry names a measurement or an improvement that has not happened. No
document may describe anything here as "unchanged", "negligible" or "low overhead" until the
named measurement exists (SPEC §9.4, TM-V0-027).

### Worst-case flat archive manifest memory
- Subject: a 768 MiB `manifest.json` held as bytes, parsed value and re-encoding during the
  canonical check; `archive verify` retaining every receipt and intent file for the chain check
- Required evidence: RSS measurement at the frozen caps (SPEC §3.5)
- State: NOT_RUN; a segmented manifest needs a later versioned profile only if this shows
  flat parsing exceeds the qualified memory budget
- Owner ticket: TCP-02 (before `archive restore`)

### GP performance gate
- Subject: no Corvint Go source may slow down (TM-V0-027)
- Required evidence: frozen `taskman-perf-baseline/0` and the §9.4 interleaved run
- State: NOT_RUN
- Owner ticket: TCP-03

<!--
### <subject>
- Required evidence: <measurement>
- State: <NOT_RUN|MEASURED(run <date>)>
- Owner ticket: <TCP-NN>
-->
