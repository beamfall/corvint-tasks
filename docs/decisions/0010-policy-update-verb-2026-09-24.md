# Decision 0010: policy update verb

- Date: 2026-09-24
- Status: accepted
- Ticket: Corvint V1-0208

## Context

ATM-V0-027 accepts that policy changes are `OWNER`/`OPERATOR` mutations. Until now,
`corvint-tasks` wrote `intent/policy.json` only at `init`. Editing the projection by hand makes
the journal and the file disagree (`INTENT_DIVERGED`), and no verb could record a new policy.

Corvint's release checklist needs a policy gate that runs `--pre-promotion`. Adding that gate
means changing the policy of an initialized queue, which the store could not do.

## Decision

Add one journaled verb:

    corvint-tasks policy update --request-id ID --expected-policy-version N --file PATH [--role OWNER|OPERATOR]

- The role defaults to `OWNER`. `WORKER` and `REVIEWER` are refused `UNAUTHORIZED`.
- `PATH` must hold a canonical `taskman-policy/0` file of at most 256 KiB whose `policyVersion`
  is `N+1`. `N` must equal the current version; a stale `N` is `REVISION_CONFLICT`. The queue
  file and queue identity cannot change, and the policy has no queue field.
- The transaction is a new `POLICY_UPDATE` operation. The receipt kind and stage operation share
  that name. The stage slot and byte caps are 5 and 1470, the exact encoder maxima. The only
  post is `intent/policy.json`, written through the §5.2 writer in its usual order: evidence,
  then receipt link-in, then post, then head.
- The receipt's `pre` and `post` entries carry the old and new policy sha256. The CLI result
  repeats them as `oldPolicySha256` and `newPolicySha256`.
- Replay, `REQUEST_ID_CONFLICT`, pending-receipt redo, restore, `VERSION`, `ALL` barrier,
  intent-branch, divergence and staging guards behave as they do for the other fixture writers.
  The writer is fixture-only, like every other callable writer.
- Readers, `receipt audit` and the closed receipt-kind set accept `POLICY_UPDATE`.
- A release candidate pins the policy digest, so an update makes earlier candidates stale
  (`candidate-source-or-policy`). This is intended.

SPEC: TM-V0-030, §2, §3.4 receipt kinds, the §5.5 stage table, §5.7.3 and §12.

## Rollback

- **Policy:** run `policy update` again with the old policy bytes at the next version. For
  example, if version 2 was wrong, submit the version-1 content as version 3. History is never
  rewritten.
- **Binary:** journals are append-only, so a store that holds a `POLICY_UPDATE` receipt needs a
  binary that can read that kind. If the verb is reverted, keep the receipt-kind and stage
  decoders. Remove only the CLI verb and the `store.PolicyUpdate` entrypoint.

## Residual risks

- Roles are self-declared. `ActorAuthentication` remains `NOT_OBSERVED` (decision 0003), so
  `--role OWNER` is an operator claim, not an authenticated identity.
- Release attestations are not checked against the argv of the policy gate they claim. That is a
  separate follow-up.
- Returned-fault staging and redo are tested. Process crash, power loss and hostile-editor CAS
  qualification remain `NOT_RUN`.
