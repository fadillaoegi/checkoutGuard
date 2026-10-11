package application

import (
	"checkoutguard/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"log"
	"strings"
	"time"
)

func (a *Service) execute(ctx context.Context, run Run) {
	defer a.wg.Done()
	defer func() { a.mu.Lock(); a.activeID = ""; a.cancel = nil; a.mu.Unlock() }()
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()
	ctx, root := a.Tracer.Workflow(ctx, "checkoutguard.checkout")
	defer root.End()
	run.TraceID = root.ID()
	root.Input(map[string]any{"run_id": run.ID, "goal": run.Goal, "url": run.URL, "environment": run.Environment, "synthetic": !run.IsMarketplace(), "parent_run": run.ParentID})
	defer func() {
		now := time.Now().UTC()
		run.Finished = &now
		if ctx.Err() != nil && (run.Status == "running" || run.Status == "failed") {
			run.Status = "stopped"
			run.Reason = "Run cancelled or duration budget exhausted"
		}
		if run.Status == "running" {
			run.Status = "failed"
			run.Reason = "Action budget exhausted"
		}
		if run.Status == "failed" || run.Status == "stopped" || run.Status == "blocked" || run.Status == "needs_login" {
			root.Error(run.Reason)
		}
		verifyCtx, verifyCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if !run.IsMarketplace() {
			if orders, err := a.Store.Orders(verifyCtx, run.ID); err == nil {
				run.OrderCount = len(orders)
			}
		}
		verifyCancel()
		if err := a.Store.PersistFinal(&run); err != nil {
			log.Print("final run persistence failed; run cannot be confirmed complete")
		}
		output := map[string]any{"run_id": run.ID, "environment": run.Environment, "status": run.Status, "reason": run.Reason, "findings": run.Findings, "usage": run.Usage}
		calls := []map[string]any{}
		for _, step := range run.Steps {
			if step.Usage.ModelCall != nil {
				calls = append(calls, map[string]any{"step_number": step.Number, "diagnostics": step.Usage.ModelCall})
			}
		}
		// Hosted v3 humanizes structured output keys. A JSON string preserves
		// typed per-call observations and span IDs through that presentation.
		if encoded, err := json.Marshal(calls); err == nil {
			output["model_calls"] = "checkoutguard.model_calls.v1=" + string(encoded)
		}
		if run.IsMarketplace() {
			output["order_verification"] = "unavailable"
			output["checkout_review"] = run.CheckoutReview
		} else {
			output["order_count"] = run.OrderCount
		}
		root.Output(output)
	}()
	if err := a.Store.Save(ctx, &run); err != nil {
		run.Status = "failed"
		run.Reason = "Cannot persist run start"
		return
	}
	if run.IsMarketplace() {
		a.executeMarketplace(ctx, &run)
		return
	}
	browser, err := a.Browsers.Open(ctx, run.URL, a.EvidenceDir, run.ID)
	if err != nil {
		run.Status = "failed"
		run.Reason = err.Error()
		return
	}
	defer browser.Close()
	obs, err := browser.Observe(0)
	if err != nil {
		run.Status = "failed"
		run.Reason = err.Error()
		return
	}
	run.Steps = append(run.Steps, Step{Number: 0, Action: Action{Tool: "observe", Reason: "Initial browser observation"}, Result: obs, Evidence: obs.Evidence})
	if err = a.Store.Save(ctx, &run); err != nil {
		run.Status = "failed"
		run.Reason = "Cannot persist initial evidence"
		return
	}
	messages := []Message{{Role: "system", Content: systemPrompt}, {Role: "user", Content: "Goal: " + run.Goal}, observationMessage(obs)}
	placed, statusChecked, retries, modelErrors, consecutiveErrors := false, false, 0, 0, 0
	for number := 1; number <= 18; number++ {
		if ctx.Err() != nil {
			return
		}
		started := time.Now()
		action, usage, planErr := a.Planner.Next(ctx, messages)
		if planErr == nil {
			planErr = domain.ValidateAction(action)
		}
		if planErr == nil {
			consecutiveErrors = 0
		}
		run.Usage.Input += usage.Input
		run.Usage.Output += usage.Output
		run.Usage.CostIncomplete = run.Usage.CostIncomplete || usage.CostIncomplete
		if usage.Model != "" {
			run.Model = usage.Model
		}
		if usage.CostUSD != nil {
			if run.Usage.CostUSD == nil {
				v := 0.0
				run.Usage.CostUSD = &v
			}
			*run.Usage.CostUSD += *usage.CostUSD
		}
		step := Step{Number: number, Action: action, Usage: usage}
		var result any
		if planErr != nil {
			modelErrors++
			consecutiveErrors++
			step.Action = Action{Tool: "model_error", Reason: "Model response failed validation or transport"}
			step.Error = planErr.Error()
			result = modelFailureResult(planErr)
		}
		if planErr == nil {
			toolCtx, span := a.Tracer.Tool(ctx, action.Tool)
			span.Input(action)
			switch action.Tool {
			case "observe":
				result = map[string]any{"status": "observed"}
			case "click", "fill":
				if action.Selector == "#place" && placed {
					if !statusChecked {
						err = errors.New("inspect order_status before any retry")
					} else if run.OrderCount > 0 {
						err = errors.New("order already confirmed; do not resubmit")
					} else if retries >= 1 {
						err = errors.New("payment retry budget exhausted")
					} else {
						retries++
					}
				} else {
					err = nil
				}
				if err == nil {
					err = browser.Act(action, obs)
					if action.Selector == "#place" {
						placed = true
						statusChecked = false
					}
				}
				if err != nil {
					step.Error = err.Error()
				}
				result = map[string]any{"status": "action executed", "error": step.Error}
			case "verify_total":
				if obs.Quote == nil {
					step.Error = "checkout amounts are not visible"
					result = map[string]any{"error": step.Error}
				} else {
					q := *obs.Quote
					expected := q.Expected()
					run.Expected = expected
					run.Displayed = q.Total
					run.TotalVerified = true
					valid := expected == q.Total && q.Subtotal == q.Unit*q.Quantity
					result = map[string]any{"matches": valid, "expected_minor": expected, "displayed_minor": q.Total, "currency": q.Currency}
					if !valid {
						run.AddFinding(Finding{Kind: "incorrect_total", Title: "Checkout total does not match the breakdown", Detail: fmt.Sprintf("Expected %d %s; checkout displays %d %s", expected, q.Currency, q.Total, q.Currency), Evidence: obs.Evidence, Verified: true})
					}
				}
			case "order_status":
				orders, lookupErr := a.Store.Orders(toolCtx, run.ID)
				if lookupErr != nil {
					step.Error = "order status lookup failed"
					result = map[string]any{"error": step.Error}
				} else {
					run.OrderCount = len(orders)
					statusChecked = true
					result = map[string]any{"correlation_id": run.Correlation, "orders": orders, "count": len(orders)}
					if len(orders) > 1 {
						run.AddFinding(Finding{Kind: "duplicate_order", Title: "Duplicate order detected", Detail: "Multiple persisted orders for a single checkout", Evidence: obs.Evidence, Verified: true})
					}
					if placed && len(orders) == 1 && strings.Contains(strings.ToLower(obs.Text), "timed out") {
						run.AddFinding(Finding{Kind: "timeout_recovered", Title: "Payment timeout recovered without a duplicate", Detail: "Persisted order confirmed using the original correlation ID before retry", Evidence: obs.Evidence, Verified: true})
					}
				}
			case "finish":
				orders, lookupErr := a.Store.Orders(toolCtx, run.ID)
				if lookupErr != nil {
					step.Error = "final order verification failed"
				} else {
					run.OrderCount = len(orders)
					if !run.TotalVerified || !placed || !statusChecked || len(orders) != 1 {
						step.Error = "finish requires verified checkout amounts, order submission, and order_status confirming exactly one order"
					} else if orders[0].Total != run.Displayed {
						step.Error = "persisted order amount differs from the observed checkout"
					} else {
						run.Status = "passed"
						run.Reason = "Checkout verified; exactly one persisted synthetic order"
						if run.HasDefect() {
							run.Status = "findings"
							run.Reason = "Verified pricing defect; exactly one persisted synthetic order"
						}
					}
				}
				result = map[string]any{"status": run.Status, "error": step.Error, "order_count": run.OrderCount}
			}
			span.Output(result)
			if step.Error != "" {
				span.Error(step.Error)
			}
			span.End()
		}
		step.Result = result
		// Every step captures the resulting page; retain a last safe observation on errors.
		nextObs, observeErr := browser.Observe(number)
		if observeErr == nil {
			obs = nextObs
			step.Evidence = obs.Evidence
		} else {
			if step.Error == "" {
				step.Error = observeErr.Error()
			}
			run.Status = "failed"
			run.Reason = "Browser evidence became unavailable"
		}
		step.LatencyMS = time.Since(started).Milliseconds()
		run.Steps = append(run.Steps, step)
		if err = a.Store.Save(ctx, &run); err != nil {
			run.Status = "failed"
			run.Reason = "Cannot persist step and evidence"
			return
		}
		if run.Status != "running" {
			return
		}
		if ctx.Err() != nil {
			return
		}
		if planErr != nil {
			delay, reason := modelRecovery(planErr, consecutiveErrors, modelErrors)
			if reason != "" {
				run.Status = "failed"
				run.Reason = reason
				return
			}
			if err := waitModelRetry(ctx, delay); err != nil {
				return
			}
		}
		messages = decisionMessages(messages, step, observationMessage(obs))

	}
}
