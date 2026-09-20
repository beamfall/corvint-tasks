# Corvint 1.0 fixture release chain

Create `v0-9`, then create `v1-0` with `v0-9` as its predecessor. Each release scopes explicit tickets and acceptance criteria. Complete the scoped tickets, capture a clean candidate, record compatible external or manual attestations for every derived gate and criterion, and promote `v0-9`. The `v1-0` candidate then binds the immutable `v0-9` promotion digest.

`release readiness` reports `BLOCKED`, `UNKNOWN`, or `READY_ATTESTED`. `.taskman` journal/projection writes do not change the captured source identity; any other tracked or untracked source change invalidates it. `release promote` records local readiness only. Native gate execution, tagging, publication, deployment, and real-queue cutover remain `NOT_RUN`.
