# 0006 — Corvint Tasks product identity

Status: accepted, 2026-09-15.

The product and executable are Corvint Tasks and `corvint-tasks`. The local source coordinate is
`github.com/Beamfall/corvint-tasks`; it records the planned module identity and does not claim that
a remote repository exists or is owned.

This identity change preserves `.taskman`, every `taskman-*/0` wire profile, canonical record and
archive bytes, numeric limits, journal behavior, history, licensing, and version
`0.0.0-tcp01-unverified`. `CORVINT_TASKS_ACTOR` is the primary local-operator identity variable.
`ATM_ACTOR` remains a compatibility fallback; equal nonempty values agree, while conflicting
nonempty values fail before repository resolution, payload reads, or mutation.

The selected current build toolchain is Go 1.27.1 with `GOTOOLCHAIN=local`. Earlier Go 1.27.0
qualification records remain historical facts and are not recharacterized as rerun evidence.

Rollback before publication: revert the product, command, module, current documentation, and actor
variable changes together. Publication remains a separate action.
