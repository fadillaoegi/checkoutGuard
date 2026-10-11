package browser

import (
	"checkoutguard/internal/application"
	"context"

	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Browser struct {
	ctx         context.Context
	cancel      context.CancelFunc
	allocCancel context.CancelFunc
	allowed     *url.URL
	EvidenceDir string
	RunID       string
}

func OpenBrowser(ctx context.Context, target, evidenceDir, id string) (*Browser, error) {
	allowed, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	options := append([]chromedp.ExecAllocatorOption{}, chromedp.DefaultExecAllocatorOptions[:]...)
	options = append(options, chromedp.WindowSize(1280, 900), chromedp.Flag("disable-background-networking", true), chromedp.Flag("disable-component-update", true))
	if path := os.Getenv("CHROME_PATH"); path != "" {
		options = append(options, chromedp.ExecPath(path))
	}
	allocator, allocCancel := chromedp.NewExecAllocator(ctx, options...)
	browserCtx, cancel := chromedp.NewContext(allocator)
	b := &Browser{ctx: browserCtx, cancel: cancel, allocCancel: allocCancel, allowed: allowed, EvidenceDir: evidenceDir, RunID: id}
	if err = chromedp.Run(browserCtx); err != nil {
		b.Close()
		return nil, errors.New("cannot start Chrome; check CHROME_PATH and sandbox permissions")
	}
	// Only same-origin resources under this run's sandbox path are allowed.
	chromedp.ListenTarget(browserCtx, func(event any) {
		if e, ok := event.(*fetch.EventRequestPaused); ok {
			go func() {
				chromedp.Run(browserCtx, chromedp.ActionFunc(func(c context.Context) error {
					if b.Allowed(e.Request.URL) {
						return fetch.ContinueRequest(e.RequestID).Do(c)
					}
					return fetch.FailRequest(e.RequestID, network.ErrorReasonBlockedByClient).Do(c)
				}))
			}()
		}
	})
	err = chromedp.Run(browserCtx, fetch.Enable(), chromedp.Navigate(target))
	if err != nil {
		b.Close()
		return nil, errors.New("sandbox navigation failed")
	}
	return b, nil
}

func (b *Browser) Allowed(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return u.Scheme == b.allowed.Scheme && u.Host == b.allowed.Host && u.User == nil && (u.Path == b.allowed.Path || strings.HasPrefix(u.Path, b.allowed.Path+"/"))
}
func (b *Browser) Close() { b.cancel(); b.allocCancel() }
func (b *Browser) Observe(number int) (Observation, error) {
	var obs Observation
	ctx, cancel := context.WithTimeout(b.ctx, 10*time.Second)
	defer cancel()
	script := `(() => {const visible=e=>!!(e.offsetWidth||e.offsetHeight||e.getClientRects().length);const controls=[...document.querySelectorAll('button[id],input[id]')].filter(e=>visible(e)&&!e.disabled).map(e=>({selector:'#'+e.id,label:e.textContent||e.getAttribute('aria-label')||e.id,type:e.tagName.toLowerCase()}));const q=document.querySelector('#quote');return {url:location.href,text:document.body.innerText.slice(0,9000),controls,quote:q?JSON.parse(q.getAttribute('data-quote')):null}})()`
	if err := chromedp.Run(ctx, chromedp.Evaluate(script, &obs)); err != nil {
		return obs, errors.New("browser observation failed")
	}
	if !b.Allowed(obs.URL) {
		return obs, errors.New("browser left allowed sandbox")
	}
	var screenshot []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&screenshot, 100)); err != nil {
		return obs, errors.New("screenshot capture failed")
	}
	path, err := b.saveScreenshot(number, screenshot)
	obs.Evidence = path
	return obs, err
}
func (b *Browser) saveScreenshot(number int, screenshot []byte) (string, error) {
	directory := filepath.Join(b.EvidenceDir, b.RunID)
	if err := os.MkdirAll(directory, 0700); err != nil {
		return "", err
	}
	name := strconv.Itoa(number) + ".png"
	if err := os.WriteFile(filepath.Join(directory, name), screenshot, 0600); err != nil {
		return "", err
	}
	return fmt.Sprintf("/evidence/%s/%s", b.RunID, name), nil
}
func (b *Browser) Act(a Action, obs Observation) error {
	var control *Control
	for i := range obs.Controls {
		if obs.Controls[i].Selector == a.Selector {
			control = &obs.Controls[i]
			break
		}
	}
	if control == nil {
		return errors.New("selector is not a visible allowed control")
	}
	if a.Tool == "fill" {
		if control.Type != "input" {
			return errors.New("fill requires an input")
		}
		if a.Selector == "#email" && a.Value != "demo@example.test" {
			return errors.New("use the synthetic demo email")
		}
		if a.Selector == "#quantity" {
			n, err := strconv.Atoi(a.Value)
			if err != nil || n < 1 || n > 5 {
				return errors.New("quantity must be 1..5")
			}
		}
	}
	if a.Tool == "click" && control.Type != "button" {
		return errors.New("click requires a button")
	}
	ctx, cancel := context.WithTimeout(b.ctx, 12*time.Second)
	defer cancel()
	var action chromedp.Action
	if a.Tool == "click" {
		action = chromedp.Click(a.Selector, chromedp.ByQuery)
	} else {
		action = chromedp.SetValue(a.Selector, a.Value, chromedp.ByQuery)
	}
	if err := chromedp.Run(ctx, action); err != nil {
		return errors.New("browser action failed")
	}
	// The sandbox exposes busy state rather than relying on arbitrary sleeps.
	return chromedp.Run(ctx, chromedp.Poll(`document.body.dataset.busy !== 'true'`, nil, chromedp.WithPollingTimeout(8*time.Second)))
}

type Factory struct{}

func (Factory) Open(ctx context.Context, target, evidence, id string) (application.SandboxBrowser, error) {
	return OpenBrowser(ctx, target, evidence, id)
}
