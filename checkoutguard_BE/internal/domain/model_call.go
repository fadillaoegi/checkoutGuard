package domain

// ModelCall records client observations, never inferred provider behavior.
type ModelCall struct {
	SpanID               string `json:"span_id"`
	RequestedModel       string `json:"requested_model"`
	ReturnedModel        string `json:"returned_model,omitempty"`
	ReturnedModelKnown   bool   `json:"returned_model_known"`
	ModelLabelMismatch   *bool  `json:"model_label_mismatch,omitempty"`
	TimeoutMS            int64  `json:"timeout_ms"`
	ErrorKind            string `json:"error_kind"`
	Retryable            bool   `json:"retryable"`
	HTTPResponseReceived bool   `json:"http_response_received"`
	HTTPStatusCode       int    `json:"http_status_code,omitempty"`
	ResponseDecoded      bool   `json:"response_decoded"`
}
