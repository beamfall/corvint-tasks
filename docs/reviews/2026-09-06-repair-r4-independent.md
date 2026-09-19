**PASS — exact R4 candidate may be applied as experimental support.** No new HIGH/MED flaws found.

- Child kill/reap remains confined to its actual parent while unreaped (`fifo_test.go:172`). Outer child-PID handling uses only signal 0 (`:352`, `:455`, `:477`).
- The 20-second observation bound covers the independent 15-second child lifetime (`:32`, `:477`).
- Expiry meaningfully exercises helper SIGKILL/reap followed by child disappearance without destructive child signaling (`:529`, `:572`), supported by the reported 15.01-second runs.
- The stale-PID witness verifies diagnostic-only signaling and the configured budget (`:355`). Minor report imprecision: it uses a single-call wait stub, not an injected clock or simulated 20-second progression.
- Actual helper/outer INT/TERM regressions and the pre-cleanup child-absence assertion remain (`:513`, `:521`, `:547`).

Parent canonical PASS and unchanged 108-file hashes accepted from supplied evidence; tests not rerun. No broader qualification implied.
