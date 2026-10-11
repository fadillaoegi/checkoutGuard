package application

import "context"

// Ports belong to the use cases. Infrastructure implements them at startup.
type Repository interface {
	Ping(context.Context) error
	Save(context.Context, *Run) error
	Get(context.Context, string) (Run, error)
	List(context.Context) ([]Run, error)
	Place(context.Context, Run, Quote) (Order, error)
	Orders(context.Context, string) ([]Order, error)
	EvidenceExists(context.Context, string, int) bool
	PersistFinal(*Run) error
}
type Planner interface {
	Next(context.Context, []Message) (Action, Usage, error)
	Info() ModelInfo
}
type SandboxBrowser interface {
	Observe(int) (Observation, error)
	Act(Action, Observation) error
	Close()
}
type BrowserFactory interface {
	Open(context.Context, string, string, string) (SandboxBrowser, error)
}
type MarketplaceBrowser interface {
	Observe(int) (MarketObservation, error)
	Act(Action, MarketObservation) error
	Close()
}
type MarketplaceGateway interface {
	State() MarketSessionState
	Connect(string, string) (MarketConnection, error)
	Close()
	Begin(context.Context, Run, string) (MarketplaceBrowser, error)
}
type Span interface {
	ID() string
	Input(any)
	Output(any)
	Error(string)
	End()
}
type Tracer interface {
	Workflow(context.Context, string) (context.Context, Span)
	Tool(context.Context, string) (context.Context, Span)
}
