package httpapi

import (
	"bytes"
	"checkoutguard/internal/application"
	"context"
	"net/http/httptest"
	"testing"
)

func TestMarketplaceAPIBoundaries(t *testing.T) {
	service := &application.Service{Planner: unconfiguredPlanner{}, BaseURL: "http://127.0.0.1:8080"}
	a := &Handler{Service: service, BaseURL: service.BaseURL}
	handler := a.Handler()
	for _, c := range []struct {
		path, body, origin, contentType string
		code                            int
	}{
		{"/api/marketplaces/shopee/connect", `{}`, "https://evil.test", "application/json", 403},
		{"/api/marketplaces/shopee/connect", `{}`, "", "application/x-www-form-urlencoded", 400},
		{"/api/marketplaces/other/connect", `{}`, "", "application/json", 400},
		{"/api/runs", `{"environment":"sandbox","goal":"Buy","url":"https://shopee.co.id/product/123/456"}`, "", "application/json", 400},
		{"/api/runs", `{"environment":"shopee","goal":"Buy","url":"https://www.tokopedia.com/shop/item"}`, "", "application/json", 400},
	} {
		r := httptest.NewRequest("POST", c.path, bytes.NewBufferString(c.body))
		r.Header.Set("Origin", c.origin)
		r.Header.Set("Content-Type", c.contentType)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != c.code {
			t.Fatalf("%s expected %d, got %d", c.path, c.code, w.Code)
		}
	}

}

type unconfiguredPlanner struct{}

func (unconfiguredPlanner) Info() application.ModelInfo { return application.ModelInfo{} }
func (unconfiguredPlanner) Next(context.Context, []application.Message) (application.Action, application.Usage, error) {
	panic("unexpected model call")
}
