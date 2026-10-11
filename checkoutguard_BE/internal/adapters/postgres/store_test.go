package postgres

import (
	"checkoutguard/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"sync"
	"testing"
	"time"
)

func TestConcurrentOrderIdempotency(t *testing.T) {
	dsn := os.Getenv("CHECKOUTGUARD_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set CHECKOUTGUARD_TEST_DATABASE_URL for PostgreSQL integration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := Store{Pool: pool}
	if _, err = pool.Exec(ctx, schema); err != nil {
		t.Fatal(err)
	}
	run := Run{ID: uuid.NewString(), Correlation: uuid.NewString(), Status: "test", Started: time.Now().UTC(), Steps: []Step{}, Findings: []Finding{}}
	if err = s.Save(ctx, &run); err != nil {
		t.Fatal(err)
	}
	defer func() {
		pool.Exec(ctx, "DELETE FROM sandbox_orders WHERE run_id=$1", run.ID)
		pool.Exec(ctx, "DELETE FROM run_checkpoints WHERE run_id=$1", run.ID)
		pool.Exec(ctx, "DELETE FROM runs WHERE id=$1", run.ID)
	}()
	q := domain.MakeQuote("standard", 2, false)
	var wg sync.WaitGroup
	ids := make(chan string, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o, e := s.Place(ctx, run, q)
			if e != nil {
				t.Error(e)
				return
			}
			ids <- o.ID
		}()
	}
	wg.Wait()
	close(ids)
	unique := map[string]bool{}
	for id := range ids {
		unique[id] = true
	}
	if len(unique) != 1 {
		t.Fatalf("concurrent submits produced %d order IDs", len(unique))
	}
	orders, err := s.Orders(ctx, run.ID)
	if err != nil || len(orders) != 1 || orders[0].Total != q.Expected() {
		t.Fatalf("unexpected persisted orders: %+v %v", orders, err)
	}
}
