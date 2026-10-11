package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"checkoutguard/internal/application"
	"github.com/google/uuid"
)

type Handler struct {
	Service              *application.Service
	BaseURL, EvidenceDir string
}

func (a *Handler) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		info, err := a.Service.Health(r.Context())
		if err != nil {
			writeError(w, 503, "database unavailable")
			return
		}
		writeJSON(w, 200, map[string]any{"status": "ok", "model_configured": info.Configured, "model": info.Name, "provider": info.Provider, "sandbox_url": a.BaseURL + "/sandbox", "environments": []string{"sandbox", "shopee", "tokopedia"}, "marketplace_mode": "prepare_checkout", "telemetry": "instrumented; hosted workflow readback requires verification"})
	})
	a.marketplaceRoutes(mux)
	mux.HandleFunc("GET /api/runs", func(w http.ResponseWriter, r *http.Request) {
		runs, err := a.Service.List(r.Context())
		if err != nil {
			writeError(w, 500, "cannot read runs")
			return
		}
		writeJSON(w, 200, runs)
	})
	mux.HandleFunc("POST /api/runs", func(w http.ResponseWriter, r *http.Request) {
		var input application.StartInput
		if err := decodeJSON(w, r, &input); err != nil {
			writeError(w, 400, "invalid run request")
			return
		}
		run, err := a.Service.Start(r.Context(), input, "")
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 202, run)
	})
	mux.HandleFunc("GET /api/runs/{id}", func(w http.ResponseWriter, r *http.Request) {
		run, err := a.Service.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			writeError(w, 404, "run not found")
			return
		}
		writeJSON(w, 200, run)
	})
	mux.HandleFunc("GET /api/runs/{id}/orders", func(w http.ResponseWriter, r *http.Request) {
		orders, err := a.Service.Orders(r.Context(), r.PathValue("id"))
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 200, orders)
	})
	mux.HandleFunc("POST /api/runs/{id}/retest", func(w http.ResponseWriter, r *http.Request) {
		run, err := a.Service.Retest(r.Context(), r.PathValue("id"))
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 202, run)
	})
	mux.HandleFunc("POST /api/runs/{id}/stop", func(w http.ResponseWriter, r *http.Request) {
		if err := a.Service.Stop(r.PathValue("id")); err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 202, map[string]string{"status": "stopping"})
	})
	mux.HandleFunc("GET /sandbox/{id}", a.sandbox)
	mux.HandleFunc("/sandbox/{id}/{resource}", a.sandbox)
	mux.HandleFunc("GET /evidence/{id}/{file}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, err := uuid.Parse(id); err != nil {
			http.NotFound(w, r)
			return
		}
		file := r.PathValue("file")
		number, err := strconv.Atoi(strings.TrimSuffix(file, ".png"))
		if err != nil || file != strconv.Itoa(number)+".png" || number < 0 || !a.Service.EvidenceExists(r.Context(), id, number) {
			http.NotFound(w, r)
			return
		}
		serveScreenshot(w, r, filepath.Join(a.EvidenceDir, id, file))
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != "GET" && r.Method != "HEAD" && r.Header.Get("Origin") != "" {
			o, err := url.Parse(r.Header.Get("Origin"))
			if err != nil || (o.Host != "127.0.0.1:5173" && o.Host != "localhost:5173" && o.String() != a.BaseURL) {
				writeError(w, 403, "origin is outside local demo")
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}
func serveScreenshot(w http.ResponseWriter, r *http.Request, path string) {
	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	var header [512]byte
	n, err := file.Read(header[:])
	if err != nil && err != io.EOF {
		http.NotFound(w, r)
		return
	}
	kind := http.DetectContentType(header[:n])
	if kind != "image/png" && kind != "image/jpeg" {
		http.NotFound(w, r)
		return
	}
	info, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Historical captures used JPEG bytes under .png names; serve their actual format.
	w.Header().Set("Content-Type", kind)
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
func decodeJSON(w http.ResponseWriter, r *http.Request, value any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("expected one JSON object")
	}
	return nil
}
func handleStartError(w http.ResponseWriter, err error) {
	code := 400
	switch {
	case errors.Is(err, application.ErrBusy), errors.Is(err, application.ErrConflict):
		code = 409
	case errors.Is(err, application.ErrModelUnavailable), errors.Is(err, application.ErrUnavailable):
		code = 503
	case errors.Is(err, application.ErrNotFound):
		code = 404
	case errors.Is(err, application.ErrPaymentTimeout):
		code = 504
	case errors.Is(err, application.ErrInternal):
		code = 500
	}
	writeError(w, code, err.Error())
}
