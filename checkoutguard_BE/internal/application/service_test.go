package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"checkoutguard/internal/domain"
)

// Ports here exercise orchestration without Chrome, a provider, or PostgreSQL.
// These are regression tests, not live model or marketplace evaluation evidence.
type memoryRepository struct {
	mu     sync.Mutex
	runs   map[string]Run
	orders map[string][]Order
}

func newRepository() *memoryRepository {
	return &memoryRepository{runs: map[string]Run{}, orders: map[string][]Order{}}
}
func (*memoryRepository) Ping(context.Context) error { return nil }
func (s *memoryRepository) Save(_ context.Context, r *Run) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *r
	copy.Steps = append([]Step{}, r.Steps...)
	copy.Findings = append([]Finding{}, r.Findings...)
	s.runs[r.ID] = copy
	return nil
}
func (s *memoryRepository) Get(_ context.Context, id string) (Run, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.runs[id]
	if !ok {
		return Run{}, ErrNotFound
	}
	return r, nil
}
func (*memoryRepository) List(context.Context) ([]Run, error) { return nil, nil }
func (s *memoryRepository) Place(_ context.Context, r Run, q Quote) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if orders := s.orders[r.ID]; len(orders) > 0 {
		return orders[0], nil
	}
	o := Order{ID: "order", Correlation: r.Correlation, Total: q.Total, Currency: q.Currency, Status: "confirmed"}
	s.orders[r.ID] = []Order{o}
	return o, nil
}
func (s *memoryRepository) Orders(_ context.Context, id string) ([]Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Order{}, s.orders[id]...), nil
}
func (*memoryRepository) EvidenceExists(context.Context, string, int) bool { return true }
func (s *memoryRepository) PersistFinal(r *Run) error                      { return s.Save(context.Background(), r) }

type testPlanner struct {
	next func(context.Context, []Message) (Action, Usage, error)
}

func (testPlanner) Info() ModelInfo {
	return ModelInfo{Name: "test", Provider: "test", Configured: true}
}
func (p testPlanner) Next(ctx context.Context, m []Message) (Action, Usage, error) {
	return p.next(ctx, m)
}

type noTrace struct{}
type noSpan struct{}

func (noTrace) Workflow(ctx context.Context, _ string) (context.Context, Span) { return ctx, noSpan{} }
func (noTrace) Tool(ctx context.Context, _ string) (context.Context, Span)     { return ctx, noSpan{} }
func (noSpan) ID() string                                                      { return "" }
func (noSpan) Input(any)                                                       {}
func (noSpan) Output(any)                                                      {}
func (noSpan) Error(string)                                                    {}
func (noSpan) End()                                                            {}

type marketGateway struct {
	state    MarketSessionState
	browser  MarketplaceBrowser
	connects int
	closes   int
}

func (g *marketGateway) State() MarketSessionState { return g.state }
func (g *marketGateway) Connect(string, string) (MarketConnection, error) {
	g.connects++
	return MarketConnection{}, nil
}
func (g *marketGateway) Close() { g.closes++ }
func (g *marketGateway) Begin(context.Context, Run, string) (MarketplaceBrowser, error) {
	return g.browser, nil
}

type marketFixture struct {
	state  string
	closed bool
}

func (b *marketFixture) Observe(int) (MarketObservation, error) {
	return MarketObservation{URL: "https://shopee.co.id/checkout", State: b.state, Summary: "Total 12400"}, nil
}
func (*marketFixture) Act(Action, MarketObservation) error { return errors.New("not expected") }
func (b *marketFixture) Close()                            { b.closed = true }

func newService(repository Repository, planner Planner, markets MarketplaceGateway) *Service {
	n := 0
	return &Service{Store: repository, Planner: planner, Markets: markets, Tracer: noTrace{}, BaseURL: "http://127.0.0.1:8080", NewID: func() string { n++; return fmt.Sprintf("id-%d", n) }}
}

func TestRunAdmissionAndCancellation(t *testing.T) {
	entered := make(chan struct{})
	markets := &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true}, browser: &marketFixture{state: "product"}}
	s := newService(newRepository(), testPlanner{next: func(ctx context.Context, _ []Message) (Action, Usage, error) {
		close(entered)
		<-ctx.Done()
		return Action{}, Usage{}, ctx.Err()
	}}, markets)
	r, err := s.Start(context.Background(), StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/123/456?tracking=secret", Goal: "Prepare one item"}, "")
	if err != nil {
		t.Fatal(err)
	}
	<-entered
	if strings.Contains(r.URL, "?") {
		t.Fatal("tracking persisted")
	}
	if _, err = s.Start(context.Background(), StartInput{Goal: "Buy", Scenario: "fixed"}, ""); !errors.Is(err, ErrBusy) {
		t.Fatalf("second run: %v", err)
	}
	if _, err = s.ConnectMarket("tokopedia"); !errors.Is(err, ErrBusy) || markets.connects != 0 {
		t.Fatal("active session was replaced")
	}
	if err = s.CloseMarket(); !errors.Is(err, ErrBusy) {
		t.Fatal("active session was closed")
	}
	if err = s.Stop(r.ID); err != nil {
		t.Fatal(err)
	}
	s.Close()
	final, _ := s.Get(context.Background(), r.ID)
	if final.Status != "stopped" || final.Finished == nil {
		t.Fatalf("cancellation was not saved: %+v", final)
	}
}

func TestMarketplaceStopsForManualReview(t *testing.T) {
	repository := newRepository()
	fixture := &marketFixture{state: "review"}
	markets := &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true}, browser: fixture}
	s := newService(repository, testPlanner{next: func(context.Context, []Message) (Action, Usage, error) {
		return Action{Tool: "finish", Reason: "Checkout review visible"}, Usage{Input: 10, Output: 2}, nil
	}}, markets)
	r, err := s.Start(context.Background(), StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/123/456", Goal: "Prepare"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s.wg.Wait()
	final, _ := s.Get(context.Background(), r.ID)
	if final.Status != "ready_for_review" || final.CheckoutReview == nil || final.CheckoutReview.OrderSubmitted || final.TotalVerified || final.OrderCount != 0 || !fixture.closed {
		t.Fatalf("manual review contract: %+v", final)
	}
	if _, err = s.Orders(context.Background(), r.ID); err == nil {
		t.Fatal("sandbox orders used for a marketplace")
	}
	if _, err = s.Retest(context.Background(), r.ID); err == nil {
		t.Fatal("marketplace accepted fixture retest")
	}
	if _, err = s.PlaceSandboxOrder(context.Background(), r.ID, 1, r.Correlation); !errors.Is(err, ErrNotFound) {
		t.Fatal("marketplace accepted a synthetic order")
	}
	s.Close()
}

func TestMarketplaceAdmissionWaitsForManualAccessCheck(t *testing.T) {
	for _, state := range []string{"blocked", "needs_login", "unavailable"} {
		t.Run(state, func(t *testing.T) {
			markets := &marketGateway{state: MarketSessionState{Environment: "shopee", BrowserOpen: true, PageState: state}}
			if state == "unavailable" {
				markets.state.PageState = "browsing"
				markets.state.PageCheck = "unavailable"
			}
			calls := 0
			s := newService(newRepository(), testPlanner{next: func(context.Context, []Message) (Action, Usage, error) { calls++; return Action{}, Usage{}, nil }}, markets)
			defer s.Close()
			run, err := s.Start(context.Background(), StartInput{Environment: "shopee", URL: "https://shopee.co.id/product/123/456", Goal: "Prepare"}, "")
			if err == nil || !strings.Contains(err.Error(), "check the browser") || calls != 0 || run.ID != "" {
				t.Fatalf("unresolved access started model/run: %+v %v calls=%d", run, err, calls)
			}
		})
	}
}

func TestSyntheticTimeoutPersistsExactlyOneOrder(t *testing.T) {
	repository := newRepository()
	s := &Service{Store: repository}
	r := Run{ID: "timeout", Correlation: "correlation", Status: "running", Scenario: "timeout", Variant: "standard"}
	repository.Save(context.Background(), &r)
	for attempt := 0; attempt < 2; attempt++ {
		_, err := s.PlaceSandboxOrder(context.Background(), r.ID, 1, r.Correlation)
		if !errors.Is(err, ErrPaymentTimeout) {
			t.Fatalf("timeout contract: %v", err)
		}
	}
	orders, err := s.Orders(context.Background(), r.ID)
	if err != nil || len(orders) != 1 || orders[0].Total != domain.MakeQuote("standard", 1, false).Expected() {
		t.Fatalf("timeout order persistence: %+v %v", orders, err)
	}
}

type fixtureBrowser struct {
	placed bool
	calls  []string
	q      Quote
}

func (b *fixtureBrowser) Observe(int) (Observation, error) {
	return Observation{Quote: &b.q, Text: "Payment response timed out", Evidence: "/evidence/test/0.png"}, nil
}
func (b *fixtureBrowser) Act(a Action, _ Observation) error {
	b.calls = append(b.calls, a.Selector)
	b.placed = true
	return nil
}
func (*fixtureBrowser) Close() {}

type fixtureFactory struct{ b *fixtureBrowser }

func (f fixtureFactory) Open(context.Context, string, string, string) (SandboxBrowser, error) {
	return f.b, nil
}

func TestAgentRejectsRetryUntilOrderStatusAndRejectsConfirmedOrder(t *testing.T) {
	repository := newRepository()
	b := &fixtureBrowser{q: domain.MakeQuote("standard", 1, false)}
	actionIndex := 0
	actions := []Action{{Tool: "verify_total"}, {Tool: "click", Selector: "#place"}, {Tool: "click", Selector: "#place"}, {Tool: "order_status"}, {Tool: "click", Selector: "#place"}, {Tool: "finish"}}
	s := newService(repository, testPlanner{next: func(_ context.Context, _ []Message) (Action, Usage, error) {
		if b.placed {
			repository.Place(context.Background(), Run{ID: "id-1", Correlation: "id-2"}, b.q)
		}
		a := actions[actionIndex]
		a.Reason = "Regression decision"
		actionIndex++
		return a, Usage{}, nil
	}}, &marketGateway{})
	s.Browsers = fixtureFactory{b: b}
	r, err := s.Start(context.Background(), StartInput{Goal: "Verify and recover", Scenario: "timeout"}, "")
	if err != nil {
		t.Fatal(err)
	}
	s.wg.Wait()
	final, _ := s.Get(context.Background(), r.ID)
	if final.Status != "passed" || len(b.calls) != 1 || !strings.Contains(final.Steps[3].Error, "inspect order_status") || !strings.Contains(final.Steps[5].Error, "already confirmed") {
		t.Fatalf("retry guards: status=%s calls=%v steps=%+v", final.Status, b.calls, final.Steps)
	}
	s.Close()
}
