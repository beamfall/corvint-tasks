package store

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Beamfall/corvint-tasks/internal/authority"
	"github.com/Beamfall/corvint-tasks/internal/intent"
	"github.com/Beamfall/corvint-tasks/internal/mutation"
	"github.com/Beamfall/corvint-tasks/internal/snapshot"
	"github.com/Beamfall/corvint-tasks/internal/transaction"
	"github.com/Beamfall/corvint-tasks/internal/wire"
)

// Mutate commits one §3.3 ticket mutation to the journal (TCP-02b).
//
// The envelope is the caller's; the trusted binding is the caller's; this
// package supplies only the facts a mutation depends on and the §5.2 ordering
// that makes the result durable. It never authenticates, never edits the
// envelope, and never decides an outcome the model did not reach.
func Mutate(ctx context.Context, repo *intent.Repository, actor mutation.Binding, envelope []byte, now wire.Timestamp) (*Report, error) {
	report := &Report{}
	if repo == nil {
		return report, wire.Errorf(wire.CodeMalformed, "", "no repository authority was resolved")
	}
	env, err := mutation.Decode(envelope)
	if err != nil {
		return report, err
	}
	if _, err = os.Stat(repo.StateDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return report, wire.Errorf(wire.CodeUninitialized, repo.StateDir, "no store: run `corvint-tasks init` first")
		}
		return report, err
	}
	if _, err = authority.Qualify(repo.CommonDir); err != nil {
		return report, err
	}
	lock, err := authority.AcquireLock(ctx, repo, authority.LockOptions{})
	if err != nil {
		return report, err
	}
	defer lock.Close()

	session, err := authority.NewSession(repo, lock)
	if err != nil {
		return report, err
	}
	defer session.Close()
	// §5.2: a receipt that was linked in but whose posts or head were not
	// written is completed before anything new is modelled, so the inventory
	// below describes a settled journal.
	redone, err := redoPending(repo, session)
	if err != nil {
		return report, err
	}
	report.Redone = redone

	// Every input is read under the lock, so the inventory, the intent store
	// and the replay observation all describe the same instant.
	store, err := intent.Load(repo.PrimaryWorktree)
	if err != nil {
		return report, err
	}
	inv, err := inventory(repo)
	if err != nil {
		return report, err
	}
	replay, err := observeReplay(repo, env.RequestID)
	if err != nil {
		return report, err
	}
	queue, policy, err := intentBytes(repo)
	if err != nil {
		return report, err
	}
	head, barrier, reservations, err := journalBytes(repo)
	if err != nil {
		return report, err
	}
	canonical, err := canonicalTickets(repo, store)
	if err != nil {
		return report, err
	}

	result := transaction.Model(
		transaction.Request{
			Operation: transaction.Mutate,
			QueueID:   store.Queue.QueueID.Raw,
			RequestID: env.RequestID,
			Actor:     actor,
			Envelope:  envelope,
		},
		transaction.Input{
			Inventory:        inv,
			Head:             head,
			Queue:            queue,
			Policy:           policy,
			Barrier:          barrier,
			Reservations:     reservations,
			CanonicalTickets: canonical,
			Premise:          transaction.LocalOperator,
			Branch:           store.Queue.IntentBranch,
			Replay:           replay,
			RecordedAt:       now,
		},
	)
	report.Outcome = result.Outcome
	report.Coverage = result.Coverage
	report.Detail = result.Detail
	report.Kind = result.Kind
	if result.Kind != "Transaction" || result.Plan == nil {
		return report, nil
	}

	report.Receipt, err = apply(repo, session, result.Plan)
	if err != nil {
		return report, err
	}
	receipt, err := snapshot.DecodeReceipt(result.Plan.Receipt())
	if err != nil {
		return report, err
	}
	if receipt.TicketID != nil {
		report.Ticket = receipt.TicketID.Raw
	}
	return report, nil
}

// observeReplay reads this request's index entry. Absence is an observation,
// never a default: a read error is returned so the model refuses
// STORAGE_FAILED rather than committing a second copy of a retried request.
func observeReplay(repo *intent.Repository, requestID string) (transaction.ReplayObservation, error) {
	path, err := snapshot.RequestPath(requestID)
	if err != nil {
		return transaction.ReplayObservation{}, err
	}
	raw, err := intent.ReadFile(filepath.Join(repo.StateDir, path), wire.MaxAttemptRecordBytes)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return transaction.ReplayObservation{State: "ABSENT"}, nil
		}
		return transaction.ReplayObservation{}, err
	}
	return transaction.ReplayObservation{State: "FOUND", Record: raw}, nil
}

func intentBytes(repo *intent.Repository) (queue, policy []byte, err error) {
	root := filepath.Join(repo.PrimaryWorktree, intent.Dir)
	queue, err = intent.ReadFile(filepath.Join(root, "queue.json"), wire.MaxQueueFileBytes)
	if err != nil {
		return nil, nil, err
	}
	policy, err = intent.ReadFile(filepath.Join(root, "policy.json"), wire.MaxPolicyFileBytes)
	if err != nil {
		return nil, nil, err
	}
	return queue, policy, nil
}

// journalBytes reads the three §5.2 mutable state files. Only the head is
// required; an absent barrier means no barrier is in force.
func journalBytes(repo *intent.Repository) (head, barrier, reservations []byte, err error) {
	head, err = intent.ReadFile(filepath.Join(repo.StateDir, "head.json"), 4096)
	if err != nil {
		return nil, nil, nil, err
	}
	barrier, err = optional(filepath.Join(repo.StateDir, "barrier.json"), 4096)
	if err != nil {
		return nil, nil, nil, err
	}
	reservations, err = optional(filepath.Join(repo.StateDir, "reservations.json"), wire.MaxReservationSetBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	return head, barrier, reservations, nil
}

func optional(path string, bound int) ([]byte, error) {
	raw, err := intent.ReadFile(path, bound)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	return raw, err
}

// canonicalTickets returns the exact bytes of every ticket projection, in the
// inventory's own order. The model binds these bytes, not a re-encoding: a
// projection that differs from its canonical record must be visible to it.
func canonicalTickets(repo *intent.Repository, store *intent.Store) ([][]byte, error) {
	root := filepath.Join(repo.PrimaryWorktree, intent.Dir, intent.TicketsDir)
	out := make([][]byte, 0, len(store.Tickets))
	for _, rec := range store.Tickets {
		raw, err := intent.ReadFile(filepath.Join(root, rec.TicketID.Local+".json"), wire.MaxTicketFileBytes)
		if err != nil {
			return nil, err
		}
		out = append(out, raw)
	}
	return out, nil
}
