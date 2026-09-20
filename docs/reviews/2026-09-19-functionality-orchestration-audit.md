# Functionality and Corvint orchestration audit — 2026-09-19

**Verdict: useful experimental fixture ticket store; not ready for autonomous LLM task management.**
The existing verification gate passes, but seven reproduced defects affect the delivered writer.
Admission, runtime supervision, completion verification, and the Corvint planning connection remain
unbuilt. Feature declarations and passing fixture tests do not establish those capabilities.

## Scope and evidence

Audited clean baseline `5117b9238f9ccc7caf01bc9a5e0ab1877ba76208` on Darwin/arm64 with Go 1.27.1.
The owner requested equal attention to safe autonomous execution, smart Corvint orchestration,
and feature completeness/developer experience. This is an audit and completion plan; it changes no
implementation, accepted contract, owner decision, queue authority, or execution permission.

One primary auditor inspected functionality and integration; one independent read-only reviewer
examined storage/mutation correctness. Task-fit choice: Astra/high for contract and correctness
reasoning; live inherited model/effort settings were not inspectable or changeable. No nested delegation.
Weekly usage was 3% at admission and the checkpoint. Working checkpoint and original scratch evidence:
`/private/tmp/corvint-tasks-audit-20260919`.

| Verification | Result and limit |
|---|---|
| `make verify` on the frozen baseline | **PASS**, 156.65 seconds: exact toolchain, source existence, formatting, all package tests, vet. [Log](2026-09-19-audit-evidence/verify.log). No source changes or commits during the gate. |
| Built executable against a new temporary fixture Git repository | **PASS** for init → create → list/show/roadmap/gates/status → export → verify, and explicit `NOT_RUN` for admit/plan/receipt/cancel/restore. Eligibility remains `UNKNOWN`. [Observations](2026-09-19-audit-evidence/smoke-summary.json). |
| Independent counterexample tests through a Go overlay | **FAIL as expected on seven defects**; unsupported VERSION negative control passes. These tests are additional evidence, not members of the passing baseline suite. [Results](2026-09-19-audit-evidence/reproductions-results.txt), [test additions](2026-09-19-audit-evidence/reproductions.go.txt). |
| Retry through the built executable | Reproduces F7: first CREATE returns `ticket:acme:main:AT-0002`; identical retry succeeds with `replayed:true`, `ticketId:""`. |
| Qualification | G1–G6 and GP remain unqualified. Native Linux, power-loss/process-interruption, saturation and runtime evidence were not produced here. Another repository's gate was active; no performance conclusion is drawn. |

Corvint's narrow authority query returned AGENTS/SPEC/ROADMAP references with explicit source
identities; it is context evidence, not a behavioral witness. Additional read-only inspection of the
sibling Corvint checkout at HEAD `6a423ac091d848b8ac5b49e8002c61c252993ac3` examined its current
`docs/specs/work-queue-observation-v0.md` and `internal/workqueue/collision.go`. Those working-tree
sources inform integration recommendations, not amendments to this repository's accepted authority.
No Corvint or Beamfall source was changed.

## Confirmed defects, highest priority first

### F1 — P1: manual ticket edits are silently adopted

[`canonicalTickets`](../../internal/store/mutate.go#L193) reads the physical projections and supplies
them as canonical records. The model then compares a projection against the same supplied bytes
([comparison](../../internal/transaction/model.go#L379)), rather than the last committed journal post.
Create a ticket, edit only its title on disk, then PRIORITIZE at revision 1: receipt 3 commits
revision 2 carrying the manual title. No reconciliation occurs.

This violates TM-V0-007 and the revision chain's authority. It is distinct from the intentionally
unbuilt reconcile command or hostile concurrent-editor CAS: a stable edit already bypasses the guard.
Repair by deriving validated canonical bytes from the journal and comparing projections before
planning. Regression: `TestAuditProjectionEditNotAdopted`; add title, acceptance-field, and queue/policy
divergence cases proving no receipt or projection changes on refusal.

### F2 — P1: intent branch restriction is never observed

[`Mutate`](../../internal/store/mutate.go#L106) passes `queue.IntentBranch` as the observed branch;
[`readIntent`](../../internal/store/store.go#L486) does the same for initialization. Setting primary
HEAD to `refs/heads/other` while the policy names `main` still permits CREATE and receipt 2.
The model's branch check therefore compares the configured value with itself.

Observe the primary worktree's actual branch under the transaction boundary; refuse detached,
unavailable, mismatched, or changed observations as required by TM-V0-007. Regression:
`TestAuditWrongBranchRefusesMutation`, extended to init, linked-worktree callers, and pending redo.

### F3 — P1: a changed primary-worktree identity permits writes

[`redoPending`](../../internal/store/redo.go#L24) decodes the head without comparing its recorded
primary worktree with the resolved repository. With `head.primaryWorktree` changed to another absolute
path, CREATE still commits. The read path already rejects this mismatch
([read guard](../../internal/cli/cli.go#L267)).

Enforce SPEC §5.2's identity guard before any redo or new effect. Regression:
`TestAuditUnsupportedStoreStatesRefuseMutation/primary-mismatch`; add a pending receipt case proving
that failed identity does not advance head or publish projections.

### F4 — P1: an ALL barrier permits ordinary mutations

The [model](../../internal/transaction/model.go#L328) inspects an existing barrier for PAUSE but does
not apply ALL scope to MUTATE. A valid seeded ALL/EMERGENCY barrier at sequence 1 still permits
CREATE and receipt 2. SPEC §5.2 requires ordinary mutations to refuse `BLOCKED/PAUSED`.

Apply barrier scope before redo/new effects, preserving the explicit recovery exemptions and the
narrower ADMISSION behavior. Regression: `TestAuditAllBarrierRefusesMutation`, plus the scope and
exemption matrix. Barrier creation is currently unwired: this is a seeded valid safety-state
counterexample, not a claim that today's pause CLI creates the barrier.

### F5 — P2: rejected initialization strands the store

[`Init`](../../internal/store/store.go#L418) creates state directories before the model validates the
actor role/request and complete inputs. `init --role INVALID` can return UNAUTHORIZED after creating
them; the next correct init is BLOCKED because [directory existence](../../internal/store/store.go#L391)
is treated as initialized state. No genesis receipt exists.

Validate rejectable inputs before creating state, and distinguish recoverable interrupted genesis
from committed initialization without deleting unknown history. Regression:
`TestAuditInvalidInitDoesNotStrandStore`; extend to malformed request IDs and interrupted genesis.

### F6 — P2: restore interruption marker is ignored

The [writer entry](../../internal/store/mutate.go#L33) checks directory existence and proceeds toward
redo; its [inventory](../../internal/store/scan.go#L19) does not include RESTORE_INCOMPLETE. A seeded
marker still allows CREATE, whereas [reads refuse it](../../internal/snapshot/probe.go#L61).

Refuse this state before redo or new effects, per SPEC C11. Regression:
`TestAuditUnsupportedStoreStatesRefuseMutation/restore-marker`. Archive restore itself is unbuilt;
the reproduction demonstrates failure to reject an explicitly unsupported interrupted state.

### F7 — P2: CREATE retry loses the allocated ticket identity

[`Mutate`](../../internal/store/mutate.go#L115) returns early for Replay and only sets `report.Ticket`
after a new transaction. Identical CREATE retries return an empty ticket ID. An orchestrator that
lost the first response cannot recover the allocated ID through the idempotent retry; receipt lookup
is also unwired. An empty *new receipt* field on replay is intentional and is not this defect.

Recover the ticket ID from the validated original receipt without writing. Regression:
`TestAuditCreateReplayPreservesTicketID`, including the CLI result and response-loss scenario.
Separately, reliable callers must retain `--issued-at`: its inclusion in the request digest is an
accepted decision, not a bug, but the automatic timestamp default makes naive retries conflict.

## Feature and functionality assessment

| Area | Delivered behavior | Missing before dependable orchestration |
|---|---|---|
| Task inventory | Strict ticket schema, dependencies/cycle checks, acceptance revisions, priorities, holds, approvals, archive/restore of individual tickets, manual completion; fourteen mutation verbs wired | Repair F1–F7. Callable transaction model permits OWNER/OPERATOR only; shared-agent identity/role enforcement is unqualified. |
| Discovery and status | List/search/show/blockers/export, roadmap and gate-definition projections, canonical JSON, pagination and explicit unknowns | Actual attempt/gate/publication evidence. `nextAction:admit` points to an omitted verb; no task is reported execution-qualified by these readers. |
| Durable store | Receipt-first link publication, request replay/conflict, post-commit redo, safe-opening and authority primitives | Reconciliation, complete startup/recovery guards, staging-descriptor crash recovery, hostile-editor CAS, production capacity/reserves and durability qualification. |
| Backup and lifecycle | Archive export and structural/digest verification; ticket tombstones retain history | Archive restore, restore/unpause recovery, administrative drain/cancel/retry/resume, receipt audit/replay CLI. Full historical archive semantic validation is a known pending item. |
| Scheduling | Ticket priority/order/dependencies/effect declarations and policy schemas | Atomic admission, durable reservations, resource collisions against live attempts, capacity/budget reservations, real liveness/generation fencing. |
| LLM execution | Runtime/capability/budget records are specified | Runtime registry execution, permission probes, process bootstrap/ack, supervised lanes, cancellation, interruption cleanup, attributable usage and unknown-budget refusal. |
| Verified completion | Manual completion is labelled MANUAL; gate/review/manifest contracts are specified | Candidate-bound gates, independent reviews, bounded repair, deterministic obligation reducer, CEM/OCM evidence and integration authorization. Process exit alone cannot satisfy completion. |
| Corvint connection | Corvint independently provides repository evidence and non-operative queue observations/proposals | Native queue adapter and `taskman-plan/0` producer/consumer path; current store does not feed those Corvint commands. |
| Model/expert routing | Capability declarations; TCP-08 describes later evaluation | Explicit model/effort selection records, qualified route alternatives, measured quality/cost evaluation and safe escalation. No smart runtime routing exists here today. |
| Operator experience | Helpful implemented/omitted command inventory and deterministic envelopes | Checked-in minimal fixture setup/examples, recoverable first-run errors, replay-complete responses, supported next actions and inspectable receipt/config/attempt state. |

The concrete unavailable commands are visible in
[`OmittedVerbs`](../../internal/cli/cli.go#L59); TCP-03 through TCP-08 are
[not started](../ROADMAP.md#L23). The synthetic runtime-free fixture mode cannot be repurposed for
real agents: decision 0004 expressly requires replacing the zero-attempt oracle first.

## How to use Corvint's strengths

1. **Supply bounded, source-bound context.** Pin Corvint query/feature/impact evidence to the task's
   source revision and accepted requirements. Keep omissions, unsupported language analysis,
   freshness and coverage visible. Use packets to guide source reads and test selection; they do
   not prove that a change works or grant execution authority.
2. **Expand declared scope into conservative collision evidence.** Corvint WQO expands touch paths
   through tracked directory prefixes and direct indexed imports/importers. That is useful conflict
   evidence, not transitive or semantic independence. Unindexed paths contribute only themselves;
   shared schemas, generators, ports and external resources still need explicit reservations.
   Carry incomplete coverage into refusal or the contracted whole-repository serial fallback.
3. **Preserve the two planning contracts.** WQO `propose-wave` analyzes collision-free wave size
   (`MAXIMUM` or bounded `GREEDY`) and remains non-operative. This project's §4.3 planner is
   priority-first: the highest-priority runnable ticket cannot be displaced by a larger cheap wave.
   Reuse evidence and collision primitives, not WQO's selection order or an implicit dispatch grant.
4. **Make the store the admission authority.** Pin queue/policy, acceptance revisions, head, intent
   tree and reservation digest in the plan; reject the entire plan on drift. Atomically recheck
   holds, approvals, qualification, collisions and budget before committing an attempt. Recheck
   before each effect, review and completion as the contract requires.
5. **Use measured model routing.** Start with an explicit builder and independent reviewer. Select
   only policy-approved, capability-probed runtime/model/effort alternatives. Record the routing
   reason, actual usage, gate results and confirmed review outcomes. Later compare quality, latency
   and cost before enabling automatic expert selection; ticket prose must never choose executable
   arguments or grant capabilities.
6. **Bind completion to evidence.** Use Corvint CEM/OCM where policy requires it, pinned to the exact
   candidate and profile versions. The deterministic reducer must derive obligations from the
   canonical ticket/policy and reject missing, stale or unverifiable evidence. Do not let a model's
   confident completion claim substitute for that reducer.

Corvint integration belongs to its own repository and authorization boundary. Freeze the required
performance baseline before any Corvint edit; run GP before promotion. No latency or cost benefit is
claimed by this audit, and process separation does not establish absence of contention.

## Completion sequence and measurable acceptance

| Order | Deliverable / owner | Smallest convincing proof |
|---|---|---|
| 1 | Repair delivered writer, TCP-02b | Promote the seven counterexamples to requirement-tagged regressions, include no-effect assertions, then one canonical gate and independent diff review. |
| 2 | Complete local store operations, TCP-02 | Inspect canonical receipt/attempt/config state; reconcile an edited intent file explicitly; recover interrupted genesis/transactions; export → restore → unpause without lost replay/history. Prove saturation and terminal/recovery reserves. |
| 3 | Connect Corvint planning, TCP-03 | Two fixture tickets plus a live reservation yield deterministic selected/deferred/blocked reasons; urgent runnable ticket wins; any plan drift refuses all admission. Read-only preview is an early usable increment. Corvint-owned record/profile alignment and GP are prerequisites to promotion. |
| 4 | Execute one bounded fixture task, TCP-02/TCP-04 | Admission → durable attempt/reservation → acknowledged spawn → controlled exit. Kill the supervisor during bootstrap, enforce turn/time caps, reject unknown required usage, cancel with a surviving child: no duplicate runtime, premature release or leaked descendant. |
| 5 | Verify one complete candidate, TCP-05 | Build → gates → independent review → at most two repairs → reducer. Seed false ACCEPT, stale gate, missing docs/CEM/OCM, dirty reviewer and revoked approval; none may complete. Pin the exact Corvint evidence profiles first. |
| 6 | Qualify one real serial ticket, TCP-06 | Passing required G2/G3 evidence, import/authority-switch qualification, GP and explicit owner cutover plus QUALIFICATION receipt. Then run capacity one with an explicit reviewer and recorded subsequent use. |
| 7 | Increase concurrency and routing, TCP-07/08 | Race admission across worktrees, include shared resources and blocked-recovery reservations, measure contention, then evaluate automatic routing against the accepted thresholds. Preserve explicit-reviewer fallback. |

Steps 3 and 4 can be developed against frozen interfaces once their storage prerequisites are
satisfied; integration remains serial. Do not add a daemon, database, UI, network provider SDK or
foreign scheduler to solve these gaps: those conflict with the chosen product boundary or require
a separate owner decision. Foreign queues remain deferred to TCP-09 and their own qualification.

## Documentation and qualification gaps

- `docs/SPEC.md` opening/digest, §10 U3, §12 ending and §5.7.2's final paragraph contain obsolete
  claims that the writer/mutations/tests do not exist. They contradict later accepted increments.
  The roadmap also retains historical “mutation CLI absent” statements. Correct current status
  without converting successful support tests into G1–G6 qualification.
- The six frozen governing sources cited by SPEC §0/decision 0001 live under a `/private/tmp`
  directory that is absent on this host. Their recorded hashes survive, but exact accepted bytes
  could not be reopened. Recover them into durable, digest-verified storage; current sibling
  documents cannot silently replace them.
- Existing [test backlog](../agent-memory/tests.md) correctly retains native Linux, interruption,
  saturation and restore gaps. Existing [fix backlog](../agent-memory/fixes.md) retains canonical
  retrieval and historical archive semantics. Preserve those holds. Large flat-manifest RSS and
  Corvint GP remain unmeasured; the passing local gate does not address them.

## Reproducing the independent counterexamples

The text file contains only the added tests; they use helpers in the baseline
`internal/store/mutate_test.go`. At the audited baseline, concatenate that source and
`2026-09-19-audit-evidence/reproductions.go.txt` into a temporary file, map the original absolute
path to it through Go's `{"Replace":{...}}` overlay, then run:

```sh
GOTOOLCHAIN=local go test -count=1 -overlay=/absolute/path/to/overlay.json ./internal/store -run '^TestAudit' -v
```

Expected baseline outcome: seven failing scenarios, one passing unsupported-VERSION control.
These tests prove the observed erroneous outcomes; they do not yet implement all required
no-write, crash, concurrent-editor or cross-platform repair regressions. The evidence additions
remain text so this audit does not alter the package suite or introduce deliberately failing tests
into the ordinary gate. All code, gates and authority decisions remain at the audited baseline.
