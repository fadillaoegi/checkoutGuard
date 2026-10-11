package browser

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"

	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"checkoutguard/internal/application"
	"checkoutguard/internal/domain"
)

// A visible, dedicated browser profile lets the user authenticate directly.
// Login credentials and cookies never pass through the API or model prompt.
type marketSession struct {
	environment string
	ctx         context.Context
	cancel      context.CancelFunc
	allocCancel context.CancelFunc
	pageState   string
	pageCheck   string
	pageMu      sync.RWMutex
}

func (s *marketSession) setPageState(state string) {
	s.pageMu.Lock()
	defer s.pageMu.Unlock()
	s.pageState = state
	s.pageCheck = "checked"
}
func (s *marketSession) getPageState() string {
	s.pageMu.RLock()
	defer s.pageMu.RUnlock()
	return s.pageState
}
func (s *marketSession) pageStatus() (string, string) {
	s.pageMu.RLock()
	defer s.pageMu.RUnlock()
	return s.pageState, s.pageCheck
}
func (s *marketSession) markPageCheckUnavailable() {
	s.pageMu.Lock()
	defer s.pageMu.Unlock()
	s.pageCheck = "unavailable"
	if s.pageState == "" {
		s.pageState = "unknown"
	}
}

func (s *marketSession) isOpen() bool {
	if s == nil || s.ctx.Err() != nil {
		return false
	}
	c := chromedp.FromContext(s.ctx)
	if c == nil || c.Browser == nil {
		return false
	}
	select {
	case <-c.Browser.LostConnection:
		return false
	default:
		return true
	}
}

func marketResourceAllowed(environment, raw, method string, document bool) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" {
		return false
	}
	if document {
		return domain.MarketPageAllowed(environment, raw) && (method == "GET" || method == "HEAD")
	}
	roots := []string{"tokopedia.com", "tokopedia.net"}
	if environment == "shopee" {
		roots = []string{"shopee.co.id", "shopeemobile.com", "susercontent.com"}
	} else if environment != "tokopedia" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	allowed := false
	for _, root := range roots {
		if host == root || strings.HasSuffix(host, "."+root) {
			allowed = true
		}
	}
	if !allowed {
		return false
	}
	path := strings.ToLower(u.Path)
	for _, term := range []string{"place_order", "placeorder", "create_order", "createorder", "/order/create", "/order/submit", "/payment", "/pay/", "/wallet", "/transfer"} {
		if strings.Contains(path, term) {
			return false
		}
	}
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return true
	}
	// Only the known Shopee cart/check-out preview routes may mutate during
	// preparation. Unknown mutations (including opaque GraphQL) fail closed.
	return environment == "shopee" && method == "POST" && strings.HasSuffix(host, "shopee.co.id") && (strings.HasPrefix(path, "/api/v4/cart/") || path == "/api/v4/checkout/get")
}

// Tokopedia GraphQL is permitted only for read queries or explicit cart root
// fields. Unknown mutations remain blocked, even if their operation is renamed.
func marketGraphQLAllowed(raw, body string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" || !(u.Hostname() == "tokopedia.com" || strings.HasSuffix(u.Hostname(), ".tokopedia.com")) || !strings.HasPrefix(u.Path, "/graphql") || len(body) > 65536 {
		return false
	}
	type operation struct {
		Query string `json:"query"`
	}
	var operations []operation
	if strings.HasPrefix(strings.TrimSpace(body), "[") {
		if json.Unmarshal([]byte(body), &operations) != nil {
			return false
		}
	} else {
		var op operation
		if json.Unmarshal([]byte(body), &op) != nil {
			return false
		}
		operations = []operation{op}
	}
	if len(operations) == 0 || len(operations) > 20 {
		return false
	}
	for _, op := range operations {
		query := strings.TrimSpace(op.Query)
		if strings.HasPrefix(query, "query ") || strings.HasPrefix(query, "query(") || strings.HasPrefix(query, "query{") || strings.HasPrefix(query, "{") {
			if regexp.MustCompile(`\b(mutation|subscription)\b`).MatchString(query) {
				return false
			}
			continue
		}
		// A single cart mutation, without aliases, directives, variables containing
		// object defaults, or additional root fields; unsupported forms stop safely.
		root := regexp.MustCompile(`(?s)^mutation(?:\s+[A-Za-z_][A-Za-z0-9_]*)?\s*(?:\([^{}]*\))?\s*\{\s*(addToCart|addCart|cartAdd|updateCart|cartUpdate|updateCartQuantity)\s*\([^{}]*\)\s*\{[^{}]*\}\s*\}$`)
		if !root.MatchString(query) {
			return false
		}
	}
	return true
}

// Session changes are serialized by application admission. Begin runs while
// admission prevents replacing or closing the session.
type Gateway struct{ session *marketSession }

func (g *Gateway) State() MarketSessionState {
	state := MarketSessionState{Authentication: "not_checked", Mode: "prepare_checkout"}
	if g.session.isOpen() {
		state.Environment = g.session.environment
		state.BrowserOpen = true
		state.PageState, state.PageCheck = g.session.pageStatus()
	}
	return state
}
func (g *Gateway) Close() {
	if g.session != nil {
		g.session.cancel()
		g.session.allocCancel()
		g.session = nil
	}
}
func (g *Gateway) Connect(environment, evidenceDir string) (MarketConnection, error) {
	if s := g.session; s.isOpen() && s.environment == environment {
		return inspectMarketSession(s), nil
	}
	g.Close()
	profile := filepath.Join(filepath.Dir(evidenceDir), "tmp", "marketplace-profiles", environment)
	if err := os.MkdirAll(profile, 0700); err != nil {
		return MarketConnection{}, errors.New("cannot create private browser profile")
	}
	opts := append([]chromedp.ExecAllocatorOption{}, chromedp.DefaultExecAllocatorOptions[:]...)
	opts = append(opts, chromedp.Flag("headless", false), chromedp.UserDataDir(profile), chromedp.WindowSize(1280, 900))
	if path := os.Getenv("CHROME_PATH"); path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	}
	alloc, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(alloc)
	s := &marketSession{environment: environment, ctx: ctx, cancel: cancel, allocCancel: allocCancel}
	// The first Run owns Chrome's lifetime. A cancelled child context here
	// would kill the browser immediately after a successful connect.
	startupTimer := time.AfterFunc(8*time.Second, cancel)
	err := chromedp.Run(ctx)
	startupTimer.Stop()
	if err != nil {
		cancel()
		allocCancel()
		return MarketConnection{}, errors.New("cannot open marketplace browser; check Chrome, connectivity, and GUI permissions")
	}
	g.session = s
	navCtx, endNav := context.WithTimeout(ctx, 8*time.Second)
	navErr := chromedp.Run(navCtx, chromedp.Navigate(domain.MarketHome(environment)))
	endNav()
	navigation := "loaded"
	if navErr != nil {
		navigation = "page did not finish loading; retry the URL manually in the open browser"
	}
	connection := inspectMarketSession(s)
	connection.Navigation = navigation
	return connection, nil

}

// Check only page state: never return the challenge URL, token, credentials,
// cookies or page contents. Browsing is not proof that the user is signed in.
func inspectMarketSession(s *marketSession) MarketConnection {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	state := "unknown"
	if err := chromedp.Run(ctx, chromedp.Evaluate(marketPageStateScript, &state)); err != nil {
		s.markPageCheckUnavailable()
	} else {
		s.setPageState(state)
	}
	state, check := s.pageStatus()
	return MarketConnection{Status: "browser_open", Authentication: "sign in directly in the browser", PageState: state, PageCheck: check}
}
func (g *Gateway) Begin(ctx context.Context, run Run, evidenceDir string) (application.MarketplaceBrowser, error) {
	session := g.session
	if !session.isOpen() || session.environment != run.Environment {
		return nil, errors.New("Marketplace browser is not open")
	}
	bctx, cancel := context.WithCancel(session.ctx)
	stop := context.AfterFunc(ctx, cancel)

	cleanup := func() {
		clean, end := context.WithTimeout(session.ctx, 2*time.Second)
		_ = chromedp.Run(clean, fetch.Disable())
		end()
		stop()
		cancel()
	}
	chromedp.ListenTarget(bctx, func(event any) {
		if e, ok := event.(*fetch.EventRequestPaused); ok {
			go func() {
				_ = chromedp.Run(bctx, chromedp.ActionFunc(func(c context.Context) error {
					allowed := marketResourceAllowed(run.Environment, e.Request.URL, e.Request.Method, e.ResourceType == network.ResourceTypeDocument)
					if !allowed && run.Environment == "tokopedia" && e.Request.Method == "POST" && e.ResourceType != network.ResourceTypeDocument && e.NetworkID != "" {
						if body, err := network.GetRequestPostData(e.NetworkID).Do(c); err == nil {
							allowed = marketGraphQLAllowed(e.Request.URL, string(body))
						}
					}
					if allowed {
						return fetch.ContinueRequest(e.RequestID).Do(c)
					}
					return fetch.FailRequest(e.RequestID, network.ErrorReasonBlockedByClient).Do(c)
				}))
			}()
		}
	})
	navCtx, navCancel := context.WithTimeout(bctx, 20*time.Second)
	err := chromedp.Run(navCtx, fetch.Enable(), chromedp.Navigate(run.URL))
	navCancel()
	if err != nil {
		cleanup()
		return nil, errors.New("Marketplace navigation failed or timed out; no order submitted")
	}
	b := &Browser{ctx: bctx, EvidenceDir: evidenceDir, RunID: run.ID}
	observe := func(number int) (MarketObservation, error) {
		var obs MarketObservation
		readCtx, end := context.WithTimeout(bctx, 10*time.Second)
		defer end()
		if err := chromedp.Run(readCtx, chromedp.Evaluate(marketObservationScript(), &obs)); err != nil {
			return obs, errors.New("marketplace observation failed")
		}
		if !domain.MarketPageAllowed(run.Environment, obs.URL) {
			return obs, errors.New("marketplace redirected outside the selected site")
		}
		session.setPageState(obs.State)
		u, _ := url.Parse(obs.URL)
		u.RawQuery = ""
		u.Fragment = ""
		obs.URL = u.String()
		var screenshot []byte
		if err := chromedp.Run(readCtx, chromedp.FullScreenshot(&screenshot, 100)); err != nil {
			return obs, errors.New("marketplace screenshot failed")
		}
		path, err := b.saveScreenshot(number, screenshot)
		obs.Evidence = path
		return obs, err
	}
	return &preparedMarketBrowser{ctx: bctx, observe: observe, close: cleanup}, nil
}

type preparedMarketBrowser struct {
	ctx     context.Context
	observe func(int) (MarketObservation, error)
	close   func()
}

func (b *preparedMarketBrowser) Observe(n int) (MarketObservation, error) { return b.observe(n) }
func (b *preparedMarketBrowser) Act(action Action, obs MarketObservation) error {
	return performMarketAction(b.ctx, action, obs)
}
func (b *preparedMarketBrowser) Close() { b.close() }

//go:embed marketplace_observe.js
var marketObservationBody string

//go:embed marketplace_policy.js
var marketPolicyScript string

//go:embed marketplace_state.js
var marketPageStateScript string

func marketObservationScript() string {
	return "(() => { const policy = " + marketPolicyScript + "; const pageState = " + marketPageStateScript + "; return " + marketObservationBody + "; })()"
}

func performMarketAction(ctx context.Context, action Action, obs MarketObservation) error {
	known := false
	knownLabel := ""
	for _, c := range obs.Controls {
		if c.Selector == action.Selector && ((action.Tool == "fill" && c.Type == "input") || (action.Tool == "click" && c.Type == "button")) {
			known = true
			knownLabel = c.Label
		}
	}
	if !known {
		return errors.New("control is outside the approved marketplace preparation actions")
	}
	if action.Tool == "fill" && !regexp.MustCompile(`^[1-5]$`).MatchString(action.Value) {
		return errors.New("marketplace quantity must be 1..5")
	}
	encoded, _ := json.Marshal(action)
	encodedLabel, _ := json.Marshal(knownLabel)
	script := `(() => { const a = ` + string(encoded) + `; const e = document.querySelector(a.selector); const p = ` + marketPolicyScript + `; if (!e || p.label(e) !== ` + string(encodedLabel) + ` || !p.allowed(e) || e.disabled || !e.getClientRects().length) return false; if(a.tool === 'click') e.click(); else { const setter = Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, 'value').set; setter.call(e, a.value); e.dispatchEvent(new Event('input', {bubbles:true})); e.dispatchEvent(new Event('change', {bubbles:true})); } return true; })()`
	actCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var executed bool
	if err := chromedp.Run(actCtx, chromedp.Evaluate(script, &executed)); err != nil {
		return errors.New("marketplace action failed")
	}
	if !executed {
		return errors.New("page control changed or is not an approved preparation action")
	}
	return chromedp.Run(actCtx, chromedp.Poll(`document.readyState !== 'loading' && !document.querySelector('[aria-busy="true"]')`, nil, chromedp.WithPollingTimeout(5*time.Second)))
}
