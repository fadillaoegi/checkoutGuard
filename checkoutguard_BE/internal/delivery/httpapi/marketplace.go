package httpapi

import (
	"net/http"
	"strings"
)

func (a *Handler) marketplaceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/marketplaces/session", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, a.Service.MarketState()) })
	mux.HandleFunc("POST /api/marketplaces/{environment}/connect", func(w http.ResponseWriter, r *http.Request) {
		var input struct{}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") || decodeJSON(w, r, &input) != nil {
			writeError(w, 400, "expected an empty JSON object")
			return
		}
		result, err := a.Service.ConnectMarket(r.PathValue("environment"))
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 200, result)
	})
	mux.HandleFunc("POST /api/marketplaces/session/close", func(w http.ResponseWriter, r *http.Request) {
		var input struct{}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") || decodeJSON(w, r, &input) != nil {
			writeError(w, 400, "expected an empty JSON object")
			return
		}
		if err := a.Service.CloseMarket(); err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 200, map[string]string{"status": "closed"})
	})
}
