package application

import (
	"checkoutguard/internal/domain"
	"context"
	"errors"
)

func (a *Service) SandboxRun(ctx context.Context, id string) (Run, error) {
	r, err := a.Get(ctx, id)
	if err != nil {
		return r, err
	}
	if r.IsMarketplace() {
		return Run{}, ErrNotFound
	}
	return r, nil
}
func (a *Service) SandboxQuote(ctx context.Context, id string, quantity int64) (Quote, error) {
	r, err := a.SandboxRun(ctx, id)
	if err != nil {
		return Quote{}, err
	}
	if quantity < 1 || quantity > 5 {
		return Quote{}, errors.New("quantity must be 1..5")
	}
	return domain.MakeQuote(r.Variant, quantity, r.Scenario == "total_bug"), nil
}
func (a *Service) PlaceSandboxOrder(ctx context.Context, id string, quantity int64, correlation string) (Order, error) {
	r, err := a.SandboxRun(ctx, id)
	if err != nil {
		return Order{}, err
	}
	if r.Status != "running" {
		return Order{}, failure(ErrConflict, "run is not active")
	}
	if quantity < 1 || quantity > 5 || correlation != r.Correlation {
		return Order{}, errors.New("invalid synthetic order")
	}
	order, err := a.Store.Place(ctx, r, domain.MakeQuote(r.Variant, quantity, r.Scenario == "total_bug"))
	if err != nil {
		return Order{}, failure(ErrInternal, "order persistence failed")
	}
	if r.Scenario == "timeout" {
		return order, ErrPaymentTimeout
	}
	return order, nil
}
