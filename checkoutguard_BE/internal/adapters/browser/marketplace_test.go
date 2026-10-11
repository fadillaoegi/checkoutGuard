package browser

import (
	"context"
	"encoding/json"
	"github.com/chromedp/chromedp"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMarketplaceNetworkBoundaries(t *testing.T) {
	for _, c := range []struct {
		environment, url, method string
		document, allowed        bool
	}{
		{"shopee", "https://shopee.co.id/checkout", "GET", true, true},
		{"shopee", "https://down-id.img.susercontent.com/file/product", "GET", false, true},
		{"shopee", "https://down-id.img.susercontent.com/file/product", "GET", true, false},
		{"shopee", "https://shopee.co.id/api/v4/cart/add_to_cart", "POST", false, true},
		{"shopee", "https://shopee.co.id/api/v4/checkout/get", "POST", false, true},
		{"shopee", "https://shopee.co.id/api/v4/checkout/place_order", "POST", false, false},
		{"shopee", "https://shopee.co.id/api/v4/account/update", "POST", false, false},
		{"shopee", "https://evilshopee.co.id/api/v4/cart/add", "POST", false, false},
		{"tokopedia", "https://www.tokopedia.com/order/create", "POST", true, false},
		{"tokopedia", "https://www.tokopedia.com/payment/charge", "GET", false, false},
		{"tokopedia", "http://127.0.0.1/api/runs", "GET", false, false},
	} {
		if got := marketResourceAllowed(c.environment, c.url, c.method, c.document); got != c.allowed {
			t.Errorf("%s accepted=%t", c.url, got)
		}
	}
	for _, c := range []struct {
		query   string
		allowed bool
	}{
		{`query Product($id: Int!) { product(id: $id) { name } }`, true},
		{`mutation Cart($id: Int!) { addToCart(id: $id) { success } }`, true},
		{`mutation Cart($id: Int!) { createOrder(id: $id) { success } }`, false},
		{`mutation Cart { addToCart(id: 1) { success } createOrder(id: 1) { success } }`, false},
		{`mutation Cart { addToCart: createOrder(id: 1) { success } }`, false},
		{``, false},
	} {
		body, _ := json.Marshal(map[string]string{"query": c.query, "operationName": "AddToCart"})
		if got := marketGraphQLAllowed("https://g.tokopedia.com/graphql", string(body)); got != c.allowed {
			t.Errorf("GraphQL %q accepted=%t", c.query, got)
		}
		if marketGraphQLAllowed("https://g.tokopedia.com.evil.test/graphql", string(body)) {
			t.Fatal("GraphQL host escape")
		}
	}
}

func TestFailedPageCheckPreservesLastAccessState(t *testing.T) {
	ctx, cancel := chromedp.NewContext(context.Background())
	cancel()
	session := &marketSession{ctx: ctx}
	session.setPageState("blocked")
	connection := inspectMarketSession(session)
	if connection.PageState != "blocked" || connection.PageCheck != "unavailable" {
		t.Fatalf("failed check erased blocked state: %+v", connection)
	}
	session.setPageState("browsing")
	state, check := session.pageStatus()
	if state != "browsing" || check != "checked" {
		t.Fatal("successful observation did not restore checked state")
	}
}

// Browser fixtures verify actual DOM execution and a changing final-order
// button. They are synthetic adapter tests, not evidence of marketplace success.
func TestMarketplaceDOMPreparation(t *testing.T) {
	if os.Getenv("CHECKOUTGUARD_BROWSER_TEST") != "1" {
		t.Skip("set CHECKOUTGUARD_BROWSER_TEST=1 for actual Chrome adapter tests")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Path == "/login" {
			w.Write([]byte(`<input type="password"><button>Login</button>`))
			return
		}
		if r.URL.Path == "/challenge" {
			w.Write([]byte(`<h1>Verify you are human</h1>`))
			return
		}
		if r.URL.Path == "/verify/captcha" || r.URL.Path == "/verification-error" {
			w.Write([]byte(`<h1>Please Try Again Later</h1><p>Verification can't be completed. Please try again after a few minutes.</p><button>Try Again</button><button>Beli sekarang</button>`))
			return
		}
		if r.URL.Path == "/verify/empty-challenge" {
			w.Write([]byte(`<h1>Please Try Again Later</h1><button>Beli sekarang</button>`))
			return
		}
		if r.URL.Path == "/product-retry" {
			w.Write([]byte(`<h1>Product</h1><p>Try again later for more stock</p><button>Beli sekarang</button>`))
			return
		}
		w.Write([]byte(`<h1>Checkout</h1><p>Total Rp 12.400</p><p>demo@example.test 081234567890</p><div class="shipping-address">Secret shipping address</div><div class="product-variant"><button id="variant" onclick="window.variantChosen=true">Blue</button></div><input id="quantity" name="quantity" aria-label="Quantity" value="1"><button id="buy">Beli sekarang</button><button id="final" onclick="window.orderCreated=true">Buat pesanan</button><button id="unknown">Change password</button>`))
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	browser, closeBrowser := chromedp.NewContext(ctx)
	defer closeBrowser()
	if err := chromedp.Run(browser, chromedp.Navigate(server.URL+"/checkout")); err != nil {
		t.Fatal(err)
	}
	var obs MarketObservation
	if err := chromedp.Run(browser, chromedp.Evaluate(marketObservationScript(), &obs)); err != nil {
		t.Fatal(err)
	}
	if obs.State != "review" || len(obs.Controls) != 3 || strings.Contains(obs.Text, "demo@example.test") || strings.Contains(obs.Text, "081234567890") || strings.Contains(obs.Text, "Secret shipping") {
		t.Fatalf("unexpected observation: %+v", obs)
	}
	var buy, qty, variant Control
	for _, c := range obs.Controls {
		switch c.Label {
		case "beli sekarang":
			buy = c
		case "quantity":
			qty = c
		case "blue":
			variant = c
		}
	}
	if err := performMarketAction(browser, Action{Tool: "fill", Selector: qty.Selector, Value: "2"}, obs); err != nil {
		t.Fatal(err)
	}
	if err := performMarketAction(browser, Action{Tool: "click", Selector: variant.Selector}, obs); err != nil {
		t.Fatal(err)
	}
	var picked bool
	if err := chromedp.Run(browser, chromedp.Evaluate(`window.variantChosen === true && document.querySelector('#quantity').value === '2'`, &picked)); err != nil || !picked {
		t.Fatal("variant/quantity action did not execute")
	}
	if err := performMarketAction(browser, Action{Tool: "fill", Selector: qty.Selector, Value: "99"}, obs); err == nil {
		t.Fatal("unbounded quantity accepted")
	}
	if err := performMarketAction(browser, Action{Tool: "click", Selector: "#final"}, obs); err == nil {
		t.Fatal("final order button accepted")
	}
	// A previously allowed control changed into a final-order button. Rechecking
	// the DOM must block it even when the old snapshot allows its selector.
	if err := chromedp.Run(browser, chromedp.Evaluate(`document.querySelector('#buy').textContent='Buat pesanan'; window.__checkoutGuardPolicy={allowed:()=>true}`, nil)); err != nil {
		t.Fatal(err)
	}
	if err := performMarketAction(browser, Action{Tool: "click", Selector: buy.Selector}, obs); err == nil {
		t.Fatal("changed final button or page-supplied policy accepted")
	}
	for _, c := range []struct{ path, state string }{{"/login", "needs_login"}, {"/challenge", "blocked"}, {"/verify/captcha?anti_bot_tracking_id=private-test-value", "blocked"}, {"/verify/empty-challenge", "blocked"}, {"/verification-error", "blocked"}, {"/product-retry", "browsing"}} {
		if err := chromedp.Run(browser, chromedp.Navigate(server.URL+c.path), chromedp.Evaluate(marketObservationScript(), &obs)); err != nil {
			t.Fatal(err)
		}
		if obs.State != c.state {
			t.Fatalf("%s state=%s", c.path, obs.State)
		}
		if (c.state == "blocked" || c.state == "needs_login") && len(obs.Controls) != 0 {
			t.Fatal("challenge/login page exposed agent controls")
		}
		session := &marketSession{ctx: browser}
		connection := inspectMarketSession(session)
		if connection.PageState != c.state || session.getPageState() != c.state || connection.Authentication != "sign in directly in the browser" {
			t.Fatalf("connection page state incorrect: %+v", connection)
		}
		encoded, _ := json.Marshal(connection)
		if strings.Contains(string(encoded), "private-test-value") || strings.Contains(string(encoded), "anti_bot_tracking_id") {
			t.Fatal("challenge token exposed by browser check")
		}
	}
}
