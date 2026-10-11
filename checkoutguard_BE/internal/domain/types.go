package domain

import (
	"time"
)

type Action struct {
	Tool     string `json:"tool"`
	Selector string `json:"selector,omitempty"`
	Value    string `json:"value,omitempty"`
	Reason   string `json:"reason"`
}
type Usage struct {
	Model          string   `json:"model,omitempty"`
	Input          int      `json:"input_tokens"`
	Output         int      `json:"output_tokens"`
	CostUSD        *float64 `json:"cost_usd"`
	CostIncomplete bool     `json:"cost_incomplete,omitempty"`
}
type Step struct {
	Number    int    `json:"number"`
	Action    Action `json:"action"`
	Result    any    `json:"result"`
	Evidence  string `json:"evidence,omitempty"`
	LatencyMS int64  `json:"latency_ms"`
	Usage     Usage  `json:"usage"`
	Error     string `json:"error,omitempty"`
}
type Finding struct {
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Evidence string `json:"evidence"`
	Verified bool   `json:"verified"`
}
type Run struct {
	Environment    string          `json:"environment,omitempty"`
	CheckoutReview *CheckoutReview `json:"checkout_review,omitempty"`
	ID             string          `json:"id"`
	ParentID       string          `json:"parent_id,omitempty"`
	Goal           string          `json:"goal"`
	URL            string          `json:"url"`
	Scenario       string          `json:"scenario"`
	Variant        string          `json:"variant"`
	Status         string          `json:"status"`
	Reason         string          `json:"reason,omitempty"`
	Model          string          `json:"model"`
	RequestedModel string          `json:"requested_model"`
	Provider       string          `json:"provider"`
	TraceID        string          `json:"trace_id,omitempty"`
	Correlation    string          `json:"correlation_id"`
	Started        time.Time       `json:"started_at"`
	Finished       *time.Time      `json:"finished_at,omitempty"`
	Steps          []Step          `json:"steps"`
	Findings       []Finding       `json:"findings"`
	Usage          Usage           `json:"usage"`
	OrderCount     int             `json:"order_count"`
	TotalVerified  bool            `json:"total_verified"`
	Expected       int64           `json:"expected_total_minor"`
	Displayed      int64           `json:"displayed_total_minor"`
}
type CheckoutReview struct {
	URL            string `json:"url"`
	Summary        string `json:"summary"`
	OrderSubmitted bool   `json:"order_submitted"`
}

func (r Run) IsMarketplace() bool { return r.Environment == "shopee" || r.Environment == "tokopedia" }

type Quote struct {
	Unit     int64  `json:"unit_minor"`
	Quantity int64  `json:"quantity"`
	Subtotal int64  `json:"subtotal_minor"`
	Discount int64  `json:"discount_minor"`
	Shipping int64  `json:"shipping_minor"`
	Total    int64  `json:"total_minor"`
	Currency string `json:"currency"`
}

func MakeQuote(variant string, quantity int64, broken bool) Quote {
	unit, shipping, discount := int64(12900), int64(1500), int64(2000)
	if variant == "alternate" {
		unit, shipping, discount = 8750, 2300, 1250
	}
	q := Quote{Unit: unit, Quantity: quantity, Subtotal: unit * quantity, Discount: discount, Shipping: shipping, Currency: "IDR"}
	q.Total = q.Subtotal - q.Discount + q.Shipping
	if broken {
		q.Total += q.Discount
	}
	return q
}
func (q Quote) Expected() int64 { return q.Unit*q.Quantity - q.Discount + q.Shipping }

type Order struct {
	ID          string `json:"id"`
	Correlation string `json:"correlation_id"`
	Total       int64  `json:"total_minor"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}

func (r *Run) AddFinding(f Finding) {
	for i, existing := range r.Findings {
		if existing.Kind == f.Kind {
			r.Findings[i] = f
			return
		}
	}
	r.Findings = append(r.Findings, f)
}
func (r *Run) HasDefect() bool {
	for _, f := range r.Findings {
		if f.Kind != "timeout_recovered" {
			return true
		}
	}
	return false
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Control struct {
	Selector string `json:"selector"`
	Label    string `json:"label"`
	Type     string `json:"type"`
}
type Observation struct {
	URL      string    `json:"url"`
	Text     string    `json:"text"`
	Controls []Control `json:"controls"`
	Quote    *Quote    `json:"quote,omitempty"`
	Evidence string    `json:"evidence,omitempty"`
}
type MarketObservation struct {
	URL      string    `json:"url"`
	Text     string    `json:"text"`
	Controls []Control `json:"controls"`
	Evidence string    `json:"evidence,omitempty"`
	State    string    `json:"state"`
	Summary  string    `json:"summary"`
}
type ModelInfo struct {
	Name       string
	Provider   string
	Configured bool
}
type MarketSessionState struct {
	Environment    string `json:"environment"`
	BrowserOpen    bool   `json:"browser_open"`
	Authentication string `json:"authentication"`
	Mode           string `json:"mode"`
	PageState      string `json:"page_state,omitempty"`
	PageCheck      string `json:"page_check,omitempty"`
}
type MarketConnection struct {
	Status         string `json:"status"`
	Authentication string `json:"authentication"`
	Navigation     string `json:"navigation,omitempty"`
	PageState      string `json:"page_state,omitempty"`
	PageCheck      string `json:"page_check,omitempty"`
}
