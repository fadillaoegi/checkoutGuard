package application

import (
	"checkoutguard/internal/domain"
	"context"
	"encoding/json"
	"time"
)

func (a *Service) executeMarketplace(ctx context.Context, run *Run) {
	// Admission keeps the session alive throughout this active run. Do not hold
	// the admission lock during navigation, so Stop can cancel it immediately.
	browser, err := a.Markets.Begin(ctx, *run, a.EvidenceDir)
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
	run.Steps = append(run.Steps, Step{Number: 0, Action: Action{Tool: "observe", Reason: "Observe marketplace product page"}, Result: obs, Evidence: obs.Evidence})
	if err := a.Store.Save(ctx, run); err != nil {
		run.Status = "failed"
		run.Reason = "Cannot persist marketplace evidence"
		return
	}
	blocked := func() bool {
		if obs.State == "needs_login" || obs.State == "blocked" {
			run.Status = obs.State
			if obs.State == "needs_login" {
				run.Reason = "Sign in manually in the marketplace browser, then start another run"
			} else {
				run.Reason = "Marketplace access challenge or unavailable page; resolve it manually or select Sandbox demo"
			}
			return true
		}
		return false
	}
	if blocked() {
		return
	}
	marshalObs := func() Message {
		return observationMessage(Observation{URL: obs.URL, Text: obs.Text, Controls: obs.Controls, Evidence: obs.Evidence})
	}
	messages := []Message{{Role: "system", Content: marketplacePrompt}, {Role: "user", Content: "Goal: " + run.Goal}, marshalObs()}
	modelErrors, consecutiveErrors := 0, 0
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
				zero := 0.0
				run.Usage.CostUSD = &zero
			}
			*run.Usage.CostUSD += *usage.CostUSD
		}
		step := Step{Number: number, Action: action, Usage: usage}
		if planErr != nil {
			modelErrors++
			consecutiveErrors++
			step.Action = Action{Tool: "model_error", Reason: "Model request failed"}
			step.Error = planErr.Error()
			step.Result = modelFailureResult(planErr)
		} else {
			_, span := a.Tracer.Tool(ctx, action.Tool)
			span.Input(action)
			switch action.Tool {
			case "observe":
			case "click", "fill":
				err = browser.Act(action, obs)
				if err != nil {
					step.Error = err.Error()
				}
			case "verify_total":
				step.Error = "Marketplace price breakdown is not independently verified; review the displayed total manually"
			case "order_status":
				step.Error = "No marketplace order was submitted; sandbox database cannot verify marketplace orders"
			case "finish":
				if obs.State == "review" {
					run.Status = "ready_for_review"
					run.Reason = "Checkout review page reached; final order confirmation remains manual"
					run.CheckoutReview = &CheckoutReview{URL: obs.URL, Summary: obs.Summary}
				} else {
					run.Status = "blocked"
					run.Reason = "Agent could not verify a checkout review page; no order submitted"
				}
			}
			step.Result = map[string]any{"status": run.Status, "error": step.Error, "order_submitted": false}
			span.Output(step.Result)
			if step.Error != "" {
				span.Error(step.Error)
			}
			span.End()
		}
		if run.Status == "running" {
			next, observeErr := browser.Observe(number)
			if observeErr != nil {
				step.Error = observeErr.Error()
				run.Status = "failed"
				run.Reason = observeErr.Error()
			} else {
				obs = next
			}
		}
		step.Evidence = obs.Evidence
		step.LatencyMS = time.Since(started).Milliseconds()
		if step.Result == nil {
			step.Result = map[string]any{"error": step.Error}
		}
		run.Steps = append(run.Steps, step)
		if err := a.Store.Save(ctx, run); err != nil {
			run.Status = "failed"
			run.Reason = "Cannot persist marketplace step"
			return
		}
		if run.Status != "running" || blocked() {
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
		messages = decisionMessages(messages, step, marshalObs())
	}
}

func actionJSON(value Action) string { b, _ := json.Marshal(value); return string(b) }
func resultJSON(value any) string    { b, _ := json.Marshal(value); return string(b) }
