package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool  *pgxpool.Pool
	lease *pgxpool.Conn
}

const schema = `
CREATE TABLE IF NOT EXISTS runs (
 id text PRIMARY KEY, status text NOT NULL, started_at timestamptz NOT NULL,
 parent_id text REFERENCES runs(id), payload jsonb NOT NULL
);
CREATE TABLE IF NOT EXISTS steps (run_id text REFERENCES runs(id), number integer, payload jsonb NOT NULL, PRIMARY KEY(run_id,number));
CREATE TABLE IF NOT EXISTS run_checkpoints (run_id text REFERENCES runs(id), number integer, state jsonb NOT NULL, PRIMARY KEY(run_id,number));
CREATE TABLE IF NOT EXISTS findings (run_id text REFERENCES runs(id), kind text, payload jsonb NOT NULL, PRIMARY KEY(run_id,kind));
CREATE TABLE IF NOT EXISTS evidence (run_id text REFERENCES runs(id), step integer, path text NOT NULL, PRIMARY KEY(run_id,step));
CREATE TABLE IF NOT EXISTS sandbox_orders (
 id text PRIMARY KEY, run_id text NOT NULL REFERENCES runs(id), correlation_id text NOT NULL,
 total_minor bigint NOT NULL, currency text NOT NULL, status text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now(), UNIQUE(run_id,correlation_id)
);
CREATE INDEX IF NOT EXISTS runs_started ON runs(started_at DESC);
`

func OpenStore(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	lease, err := pool.Acquire(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}
	var locked bool
	if err = lease.QueryRow(ctx, "SELECT pg_try_advisory_lock(7264012026)").Scan(&locked); err != nil || !locked {
		lease.Release()
		pool.Close()
		return nil, errors.New("another CheckoutGuard backend owns the database")
	}
	if _, err = pool.Exec(ctx, schema); err != nil {
		lease.Release()
		pool.Close()
		return nil, err
	}
	s := &Store{Pool: pool, lease: lease}
	// A process restart must not silently leave runs looking active.
	_, err = pool.Exec(ctx, `UPDATE runs SET status='stopped',payload=payload || jsonb_build_object('status','stopped','reason','Backend restarted; run interrupted','finished_at',now()) WHERE status='running'`)
	if err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() {
	if s.lease != nil {
		s.lease.Exec(context.Background(), "SELECT pg_advisory_unlock(7264012026)")
		s.lease.Release()
	}
	s.Pool.Close()
}

func (s *Store) Save(ctx context.Context, r *Run) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var parent any
	if r.ParentID != "" {
		parent = r.ParentID
	}
	_, err = tx.Exec(ctx, `INSERT INTO runs(id,status,started_at,parent_id,payload) VALUES($1,$2,$3,$4,$5) ON CONFLICT(id) DO UPDATE SET status=EXCLUDED.status,payload=EXCLUDED.payload`, r.ID, r.Status, r.Started, parent, data)
	if err != nil {
		return err
	}
	for _, step := range r.Steps {
		b, _ := json.Marshal(step)
		if _, err = tx.Exec(ctx, `INSERT INTO steps VALUES($1,$2,$3) ON CONFLICT(run_id,number) DO UPDATE SET payload=EXCLUDED.payload`, r.ID, step.Number, b); err != nil {
			return err
		}
		if step.Evidence != "" {
			if _, err = tx.Exec(ctx, `INSERT INTO evidence VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, r.ID, step.Number, step.Evidence); err != nil {
				return err
			}
		}
	}
	for _, f := range r.Findings {
		b, _ := json.Marshal(f)
		if _, err = tx.Exec(ctx, `INSERT INTO findings VALUES($1,$2,$3) ON CONFLICT(run_id,kind) DO UPDATE SET payload=EXCLUDED.payload`, r.ID, f.Kind, b); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO run_checkpoints VALUES($1,$2,$3) ON CONFLICT(run_id,number) DO UPDATE SET state=EXCLUDED.state`, r.ID, len(r.Steps), data)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (s *Store) Get(ctx context.Context, id string) (Run, error) {
	var r Run
	var b []byte
	err := s.Pool.QueryRow(ctx, "SELECT payload FROM runs WHERE id=$1", id).Scan(&b)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(b, &r)
	return r, err
}
func (s *Store) List(ctx context.Context) ([]Run, error) {
	rows, err := s.Pool.Query(ctx, "SELECT payload FROM runs ORDER BY started_at DESC LIMIT 60")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	runs := []Run{}
	for rows.Next() {
		var b []byte
		var r Run
		if err = rows.Scan(&b); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(b, &r); err != nil {
			return nil, err
		}
		runs = append(runs, r)
	}
	return runs, rows.Err()
}
func (s *Store) Place(ctx context.Context, r Run, q Quote) (Order, error) {
	o := Order{ID: uuid.NewString(), Correlation: r.Correlation, Total: q.Total, Currency: q.Currency, Status: "confirmed"}
	// A unique correlation key makes concurrent/retried submissions idempotent.
	err := s.Pool.QueryRow(ctx, `INSERT INTO sandbox_orders(id,run_id,correlation_id,total_minor,currency,status) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(run_id,correlation_id) DO UPDATE SET correlation_id=EXCLUDED.correlation_id RETURNING id,correlation_id,total_minor,currency,status`, o.ID, r.ID, o.Correlation, o.Total, o.Currency, o.Status).Scan(&o.ID, &o.Correlation, &o.Total, &o.Currency, &o.Status)
	return o, err
}
func (s *Store) Orders(ctx context.Context, id string) ([]Order, error) {
	rows, err := s.Pool.Query(ctx, `SELECT id,correlation_id,total_minor,currency,status FROM sandbox_orders WHERE run_id=$1 ORDER BY created_at`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Order{}
	for rows.Next() {
		var o Order
		if err = rows.Scan(&o.ID, &o.Correlation, &o.Total, &o.Currency, &o.Status); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
func (s *Store) EvidenceExists(ctx context.Context, id string, step int) bool {
	var found bool
	err := s.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM evidence WHERE run_id=$1 AND step=$2)`, id, step).Scan(&found)
	return err == nil && found
}
func (s *Store) PersistFinal(r *Run) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Save(ctx, r); err != nil {
		return fmt.Errorf("persist final run: %w", err)
	}
	return nil
}

var ErrNotFound = pgx.ErrNoRows

func (s *Store) Ping(ctx context.Context) error { return s.Pool.Ping(ctx) }
