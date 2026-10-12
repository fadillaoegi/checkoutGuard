package application

import (
	"checkoutguard/internal/domain"
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Service struct {
	Store                Repository
	Planner              Planner
	Browsers             BrowserFactory
	Markets              MarketplaceGateway
	Tracer               Tracer
	NewID                func() string
	BaseURL, EvidenceDir string
	mu                   sync.Mutex
	activeID             string
	cancel               context.CancelFunc
	wg                   sync.WaitGroup
}
type StartInput struct {
	Environment string `json:"environment"`
	Goal        string `json:"goal"`
	URL         string `json:"url"`
	Scenario    string `json:"scenario"`
	Variant     string `json:"variant"`
}

var ErrBusy = errors.New("one run is already active")
var ErrModelUnavailable = errors.New("model credentials are not configured")
var ErrNotFound = errors.New("run not found")
var ErrConflict = errors.New("operation conflicts with the active run")
var ErrUnavailable = errors.New("service unavailable")
var ErrPaymentTimeout = errors.New("Payment response timed out. Check order status before retrying.")
var ErrInternal = errors.New("operation failed")

func (a *Service) Health(ctx context.Context) (ModelInfo, error) {
	return a.Planner.Info(), a.Store.Ping(ctx)
}

// Publish terminal state only after browser cleanup. Holding admission across
// final persistence prevents a caller observing "finished" and then ErrBusy.
// Trace export can finish independently without occupying the browser slot.
func (a *Service) persistFinal(run *Run) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	err := a.Store.PersistFinal(run)
	if a.activeID == run.ID {
		a.activeID, a.cancel = "", nil
	}
	return err
}

func (a *Service) persistProgress(ctx context.Context, run *Run) error {
	progress := *run
	progress.Status, progress.Reason, progress.Finished = "running", "", nil
	return a.Store.Save(ctx, &progress)
}
func (a *Service) List(ctx context.Context) ([]Run, error) { return a.Store.List(ctx) }
func (a *Service) Get(ctx context.Context, id string) (Run, error) {
	r, err := a.Store.Get(ctx, id)
	if err != nil {
		return r, ErrNotFound
	}
	return r, nil
}
func (a *Service) Orders(ctx context.Context, id string) ([]Order, error) {
	r, err := a.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.IsMarketplace() {
		return nil, errors.New("marketplace order verification is not available; sandbox orders cannot verify a marketplace checkout")
	}
	orders, err := a.Store.Orders(ctx, id)
	if err != nil {
		return nil, failure(ErrInternal, "cannot read orders")
	}
	return orders, nil
}
func (a *Service) Retest(ctx context.Context, id string) (Run, error) {
	parent, err := a.Get(ctx, id)
	if err != nil {
		return Run{}, err
	}
	if parent.Status == "running" {
		return Run{}, failure(ErrConflict, "wait for baseline to finish")
	}
	if parent.IsMarketplace() {
		return Run{}, errors.New("corrected-fixture retest is available only for sandbox runs; start another marketplace run instead")
	}
	return a.Start(ctx, StartInput{Goal: parent.Goal, URL: a.BaseURL + "/sandbox", Scenario: "fixed", Variant: parent.Variant}, parent.ID)
}
func (a *Service) Stop(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeID != id || a.cancel == nil {
		return failure(ErrConflict, "run is not active")
	}
	a.cancel()
	return nil
}
func (a *Service) MarketState() MarketSessionState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Markets.State()
}
func (a *Service) ConnectMarket(environment string) (MarketConnection, error) {
	if domain.MarketHome(environment) == "" {
		return MarketConnection{}, errors.New("unknown marketplace")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeID != "" {
		return MarketConnection{}, ErrBusy
	}
	state, err := a.Markets.Connect(environment, a.EvidenceDir)
	if err != nil {
		return state, failure(ErrUnavailable, err.Error())
	}
	return state, nil
}
func (a *Service) CloseMarket() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeID != "" {
		return ErrBusy
	}
	a.Markets.Close()
	return nil
}
func (a *Service) EvidenceExists(ctx context.Context, id string, number int) bool {
	return a.Store.EvidenceExists(ctx, id, number)
}
func (a *Service) Start(requestCtx context.Context, input StartInput, parent string) (Run, error) {
	var empty Run
	input.Goal = strings.TrimSpace(input.Goal)
	if input.Goal == "" || len(input.Goal) > 1500 {
		return empty, errors.New("goal must contain 1..1500 characters")
	}
	if input.Environment == "" {
		input.Environment = "sandbox"
	}
	if err := domain.ValidateTarget(input.Environment, input.URL, a.BaseURL); err != nil {
		return empty, err
	}
	if input.Environment == "sandbox" {
		switch input.Scenario {
		case "total_bug", "timeout", "fixed":
		default:
			return empty, errors.New("unknown scenario")
		}
		if input.Variant == "" {
			input.Variant = "standard"
		}
		if input.Variant != "standard" && input.Variant != "alternate" {
			return empty, errors.New("unknown basket variant")
		}
	} else {
		input.Scenario = "marketplace_checkout"
		input.Variant = "live"
	}
	if !a.Planner.Info().Configured {
		return empty, ErrModelUnavailable
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeID != "" {
		return empty, ErrBusy
	}
	if input.Environment != "sandbox" {
		state := a.Markets.State()
		if !state.BrowserOpen || state.Environment != input.Environment {
			return empty, errors.New("open the selected marketplace browser and sign in manually before starting")
		}
		if state.PageState == "blocked" || state.PageState == "needs_login" || state.PageCheck == "unavailable" {
			return empty, errors.New("resolve marketplace access manually, then check the browser before starting")
		}
	}
	id := a.NewID()
	run := Run{Environment: input.Environment, ID: id, ParentID: parent, Goal: input.Goal, Scenario: input.Scenario, Variant: input.Variant, Status: "running", RequestedModel: a.Planner.Info().Name, Provider: a.Planner.Info().Provider, Correlation: a.NewID(), Started: time.Now().UTC(), Steps: []Step{}, Findings: []Finding{}}
	if run.IsMarketplace() {
		u, _ := url.Parse(input.URL)
		u.RawQuery = ""
		run.URL = u.String()
	} else {
		run.URL = a.BaseURL + "/sandbox/" + id
	}
	if err := a.Store.Save(requestCtx, &run); err != nil {
		return empty, failure(ErrInternal, "cannot persist new run")
	}
	if !run.IsMarketplace() {
		a.Markets.Close()
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.activeID = id
	a.cancel = cancel
	a.wg.Add(1)
	go a.execute(ctx, run)
	return run, nil
}

func (a *Service) Close() {
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel()
	}
	a.mu.Unlock()
	a.wg.Wait()
	a.mu.Lock()
	a.Markets.Close()
	a.mu.Unlock()
}

type fault struct {
	kind    error
	message string
}

func (e fault) Error() string                  { return e.message }
func (e fault) Unwrap() error                  { return e.kind }
func failure(kind error, message string) error { return fault{kind: kind, message: message} }
