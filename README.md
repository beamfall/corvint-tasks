# Corvint Tasks

`corvint-tasks` is the native Go fixture ticket store and CLI for the Corvint task-control-plane
contract. This source currently supports a local **fixture queue**: initialization, fourteen
journal-backed ticket mutations, ticket and queue reads, roadmap and gate projections, and archive
export/verification.
It does not include runtime execution, admission or dispatch, reservations, administrative verbs,
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
and omitted verbs and their accepted arguments.

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

## License

Copyright (C) 2026 Russell Lewis. Licensed under the GNU Affero General Public License, version 3.0
or later. See `LICENSE` and `PROVENANCE.md`.
