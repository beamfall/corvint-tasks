# Corvint 1.0 fixture release chain

Create `v0-9`, then create `v1-0` with `v0-9` as its predecessor. Each release scopes explicit tickets and acceptance criteria. Complete the scoped tickets, capture a clean candidate, record compatible external or manual attestations for every derived gate and criterion, and promote `v0-9`. The `v1-0` candidate then binds the immutable `v0-9` promotion digest.

`release readiness` reports `BLOCKED`, `UNKNOWN`, or `READY_ATTESTED`. `.taskman` journal/projection writes do not change the captured source identity; any other tracked or untracked source change invalidates it. `release promote` records local readiness only. Native gate execution, tagging, publication, deployment, and real-queue cutover remain `NOT_RUN`.

The following sequence assumes an initialized disposable fixture, two newly created ticket IDs in `TICKET_09` and `TICKET_10`, and `jq`. It obtains every workflow revision and binding digest from command output; it never reads `.taskman` projections. Keep generated evidence outside the repository.

```sh
ATM=./corvint-tasks
EVIDENCE=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

$ATM release create --request-id create-v0-9 --target v0-9 --payload \
  "{\"acceptanceCriteria\":[\"Corvint 0.9 acceptance\"],\"predecessorReleaseIds\":[],\"requiredGates\":[\"verify\"],\"ticketIds\":[\"$TICKET_09\"],\"title\":\"Corvint 0.9\",\"version\":\"0-9\"}"
$ATM release create --request-id create-v1-0 --target v1-0 --payload \
  "{\"acceptanceCriteria\":[\"Corvint 1.0 acceptance\"],\"predecessorReleaseIds\":[\"v0-9\"],\"requiredGates\":[\"verify\"],\"ticketIds\":[\"$TICKET_10\"],\"title\":\"Corvint 1.0\",\"version\":\"1-0\"}"

$ATM ticket complete-manual --request-id complete-v0-9 --target "$TICKET_09" \
  --expected-revision 1 --payload '{"evidence":[],"reason":"0.9 accepted"}'
$ATM release candidate --request-id candidate-v0-9 --target v0-9 --expected-revision 1
SHOW_09=$($ATM release show v0-9)
REV_09=$(printf '%s' "$SHOW_09" | jq -r '.items[0].revision')
CANDIDATE_09=$(printf '%s' "$SHOW_09" | jq -r '.items[0].candidateSha256')
$ATM release record-gate --request-id attest-v0-9 --target v0-9 --expected-revision "$REV_09" --payload \
  "{\"attestation\":{\"actor\":\"external-ci\",\"attestationId\":\"verify-v0-9\",\"candidateSha256\":\"$CANDIDATE_09\",\"criteria\":[\"0\"],\"evidence\":[\"$EVIDENCE\"],\"gateId\":\"verify\",\"profile\":\"taskman-release-attestation/0\",\"provenance\":\"EXTERNAL_ATTESTATION\",\"recordedAt\":\"2026-09-20T12:04:00Z\",\"result\":\"PASS\",\"sourceIdentity\":\"external-ci\"}}"
REV_09=$($ATM release readiness v0-9 | jq -r '.items[0].revision')
$ATM release promote --request-id promote-v0-9 --target v0-9 --expected-revision "$REV_09"
PROMOTION_09=$($ATM release show v0-9 | jq -r '.items[0].promotionSha256')

$ATM ticket complete-manual --request-id complete-v1-0 --target "$TICKET_10" \
  --expected-revision 1 --payload '{"evidence":[],"reason":"1.0 accepted"}'
$ATM release candidate --request-id candidate-v1-0 --target v1-0 --expected-revision 1
SHOW_10=$($ATM release show v1-0)
test "$(printf '%s' "$SHOW_10" | jq -r '.items[0].candidateBinding.predecessors[0].promotionSha256')" = "$PROMOTION_09"
REV_10=$(printf '%s' "$SHOW_10" | jq -r '.items[0].revision')
CANDIDATE_10=$(printf '%s' "$SHOW_10" | jq -r '.items[0].candidateSha256')
$ATM release record-gate --request-id attest-v1-0 --target v1-0 --expected-revision "$REV_10" --payload \
  "{\"attestation\":{\"actor\":\"external-ci\",\"attestationId\":\"verify-v1-0\",\"candidateSha256\":\"$CANDIDATE_10\",\"criteria\":[\"0\"],\"evidence\":[\"$EVIDENCE\"],\"gateId\":\"verify\",\"profile\":\"taskman-release-attestation/0\",\"provenance\":\"EXTERNAL_ATTESTATION\",\"recordedAt\":\"2026-09-20T12:06:00Z\",\"result\":\"PASS\",\"sourceIdentity\":\"external-ci\"}}"
REV_10=$($ATM release readiness v1-0 | jq -r '.items[0].revision')
$ATM release promote --request-id promote-v1-0 --target v1-0 --expected-revision "$REV_10"
$ATM release show v1-0
```
