package model

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"checkoutguard/internal/domain"
)

func newClient() *http.Client {
	timeout := 45 * time.Second
	if seconds, err := strconv.Atoi(os.Getenv("MODEL_REQUEST_TIMEOUT_SECONDS")); err == nil && seconds >= 10 && seconds <= 60 {
		timeout = time.Duration(seconds) * time.Second
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ResponseHeaderTimeout = timeout
	transport.MaxIdleConnsPerHost = 2
	return &http.Client{Timeout: timeout, Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

func requestError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return &domain.PlannerError{Kind: "canceled", Message: "model request cancelled or run duration exhausted"}
	}
	var dns *net.DNSError
	if errors.As(err, &dns) {
		return &domain.PlannerError{Kind: "network", Message: "model DNS lookup failed", Retryable: true}
	}
	var network net.Error
	if errors.As(err, &network) && network.Timeout() {
		return &domain.PlannerError{Kind: "timeout", Message: "model response timed out", Retryable: true}
	}
	return &domain.PlannerError{Kind: "network", Message: "model connection failed", Retryable: true}
}

func statusError(status int, header http.Header) error {
	retryable := status == 408 || status == 429 || status == 500 || status == 502 || status == 503 || status == 504
	kind := "request_rejected"
	if status == 401 || status == 403 {
		kind = "authentication"
	} else if status == 429 {
		kind = "rate_limited"
	} else if retryable {
		kind = "unavailable"
	}
	return &domain.PlannerError{Kind: kind, Message: "model HTTP " + strconv.Itoa(status), Retryable: retryable, RetryAfter: retryAfter(header.Get("Retry-After"), time.Now())}
}

func retryAfter(value string, now time.Time) time.Duration {
	if seconds, err := strconv.ParseInt(value, 10, 32); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if at, err := http.ParseTime(value); err == nil && at.After(now) {
		return at.Sub(now)
	}
	return 0
}

func invalidResponse(message string, retryable bool) error {
	return &domain.PlannerError{Kind: "invalid_response", Message: message, Retryable: retryable}
}
