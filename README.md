# Corvint Tasks

`corvint-tasks` is the native Go fixture ticket store and CLI for the Corvint task-control-plane
contract. This source currently supports a local **fixture queue**: initialization, fourteen
journal-backed ticket mutations, ticket and queue reads, roadmap and gate projections, and archive
export/verification, receipt audit, explicit intent reconciliation, and fixture pause/unpause.
It does not include runtime execution, admission or dispatch, reservations, remaining administrative verbs,
performance qualification, or authority to operate a real queue.

## Build and inspect

Go 1.27.1 is required.

```sh
make verify
go build -o corvint-tasks ./cmd/corvint-tasks
./corvint-tasks help
```

To install the command from a checkout:

```sh
go install ./cmd/corvint-tasks
corvint-tasks --version
```

The source module is `github.com/Beamfall/corvint-tasks`. A Corvint companion bundle records the
commit and tree it built `corvint-tasks` from in its `MANIFEST.json` and ships the matching source
archive. Repository work-queue observation and wave proposals are Corvint commands
(`corvint work init`, `observe`, `propose-wave`); see Corvint's `docs/INSTALL.md`. This store does
not yet feed them.

Commands emit canonical JSON envelopes. `corvint-tasks help` is the authoritative inventory of implemented
and omitted verbs.

## Supported ticket and roadmap path

The executable fixture path is covered end to end by the CLI tests: a valid fixture queue and
policy are written under `.taskman/`, `corvint-tasks init` commits the genesis receipt,
`corvint-tasks ticket create` commits a ticket, and `corvint-tasks ticket list` and
`corvint-tasks roadmap` read the resulting journal-backed projection. Run the focused witnesses
with:

```sh
go test -count=1 ./internal/cli -run 'TestTMV0008_AS(02_TicketCreateIsVisibleToTheReadVerbs|08_RoadmapAndGates|11_InitMakesTheStoreReadable)$'
```

The record schemas and canonical encoding are intentionally strict. Use the fixtures in
`internal/fixture` and the payloads in `internal/cli/*_test.go` as executable examples, and read
`docs/SPEC.md` before constructing or changing queue, policy, ticket, or mutation records.

Only fixture queues (`fixture: true`) are in the delivered scope. A non-fixture queue requires the
owner's cutover record and a `QUALIFICATION` receipt; this source contains neither. Do not point it
at a real queue or treat its local-operator binding as authentication.

## Try a disposable fixture

From this checkout, build the binary and prepare a fresh temporary Git repository:

```sh
project="$PWD"
demo=$(mktemp -d)
GOTOOLCHAIN=local go build -o "$demo/corvint-tasks" ./cmd/corvint-tasks
git init -b main "$demo"
mkdir "$demo/.taskman"
cp examples/fixture/queue.json examples/fixture/policy.json "$demo/.taskman/"
cd "$demo"
export CORVINT_TASKS_ACTOR=fixture-operator
./corvint-tasks init
./corvint-tasks ticket create --request-id example-create \
  --issued-at 2026-09-19T12:00:00Z --payload-stdin < "$project/examples/fixture/create.json"
./corvint-tasks receipt audit
```

Retain the original request ID, timestamp and payload when retrying a ticket command. Exact
retries return the original ticket ID. `receipt audit` verifies structural consistency; it leaves
historical acceptance, authentication, liveness and runtime qualification as NOT_OBSERVED.

A manually edited ticket freezes fresh mutations until an explicit reconciliation. Preserve the
edited bytes outside `.taskman`, then inspect the canonical ticket:

```sh
cp .taskman/tickets/AT-0002.json original-edit.json
./corvint-tasks reconcile inspect ticket:acme:main:AT-0002
```

To discard the edit, use `reconcile intent --target ticket:acme:main:AT-0002 --request-id keep-1
--file original-edit.json --keep-journal --canonical-sha256 DIGEST`, replacing DIGEST with the
inspector's `canonicalRecordSha256`. To adopt permitted fields, use `--adopt-file` instead of
`--keep-journal --canonical-sha256 DIGEST`. Retain the original copy and digest for retries.
KEEP stores discarded bytes as immutable evidence; ADOPT validates protected fields and revisions.
Pending receipts and fresh reconciliation with active staging require the still-unbuilt recovery
path and are refused. Validated read-only replay remains available with active staging.
These examples enable no runtime or real-queue execution.

## Fixture release control

`release create|update|candidate|record-gate|promote|list|show|readiness` operates only on initialized fixture queues. Candidate capture requires a clean primary worktree outside `.taskman`. Recorded manual/external attestations are provenance-preserving non-native evidence; the CLI reports native gate execution as `NOT_RUN`. Promotion is a local immutable journal fact, not a tag, publication, deployment, or cutover.

## Fixture admission barrier

```sh
./corvint-tasks pause --request-id example-pause
./corvint-tasks queue status
./corvint-tasks unpause --request-id example-unpause
```

Pause records an ADMISSION barrier; ordinary ticket mutations remain available. Unpause preserves
unrelated manual ticket edits. Recorded retries are stable even after a later opposite operation.
An already matching pause or absent unpause returns NoChange without recording the request.
Pending/staging recovery remains incomplete: a failure after an UNPAUSE receipt commits can
leave a pending deletion that no current callable command can finish. This remains a disposable
fixture workflow, without runtime or production qualification.

## License

Copyright (C) 2026 Russell Lewis. Licensed under the GNU Affero General Public License, version 3.0
or later. See `LICENSE` and `PROVENANCE.md`.
