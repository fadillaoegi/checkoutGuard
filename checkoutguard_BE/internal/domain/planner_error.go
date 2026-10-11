package domain

import "time"

// PlannerError is a safe, provider-neutral failure contract for orchestration.
// Message must never contain credentials, raw response bodies, or request URLs.
type PlannerError struct {
	Kind       string
	Message    string
	Retryable  bool
	RetryAfter time.Duration
}

func (e *PlannerError) Error() string { return e.Message }
