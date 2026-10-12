// demorecord records real sandbox UI runs or reviews verified saved results.
// Review mode never starts checkout or inference. A step error stops recording and
// cancels only the run started by this recorder; no automatic retry is made.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"checkoutguard/internal/domain"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

type frame struct {
	File    string  `json:"file"`
	Seconds float64 `json:"seconds"`
}
type chapter struct {
	Name    string  `json:"name"`
	Seconds float64 `json:"seconds"`
}
type recording struct {
	Status   string       `json:"status"`
	Error    string       `json:"error,omitempty"`
	Started  time.Time    `json:"started_at"`
	Duration float64      `json:"duration_seconds"`
	Frames   []frame      `json:"frames"`
	Chapters []chapter    `json:"chapters"`
	Runs     []domain.Run `json:"runs"`
	Style    string       `json:"style"`
	Audio    bool         `json:"audio"`
	Mode     string       `json:"mode"`
}

var preview = flag.Bool("preview", false, "Record interaction preview without starting a checkout run")
var geometry = flag.Bool("geometry", false, "Capture viewport geometry only, without interactions or checkout")
var reviewBaseline = flag.String("review", "", "Review a verified saved baseline instead of starting checkout")
var reviewRetest = flag.String("review-retest", "", "Verified saved retest to review with the baseline")
var language = flag.String("lang", "en", "Dashboard language shown in the recording: en keeps English, id switches to Indonesian")

func main() {
	flag.Parse()
	if err := record(); err != nil {
		fmt.Fprintln(os.Stderr, "Demo stopped:", err)
		os.Exit(1)
	}
}

func record() (failure error) {
	output := filepath.Join("..", "tmp", "video", time.Now().Format("20060102-150405"))
	if err := os.MkdirAll(filepath.Join(output, "frames"), 0700); err != nil {
		return err
	}
	r := recording{Status: "recording", Started: time.Now().UTC(), Frames: []frame{}, Chapters: []chapter{}, Runs: []domain.Run{}, Style: "continuous compositor capture; typed goal and pointer clicks; select values through DOM change events; smooth scrolling; no captions", Audio: false}
	r.Mode = "live"
	active := ""
	client := &http.Client{Timeout: 4 * time.Second}
	request := func(method, path string, into any) error {
		req, err := http.NewRequest(method, "http://127.0.0.1:8080"+path, nil)
		if err != nil {
			return err
		}
		response, err := client.Do(req)
		if err != nil {
			return errors.New("local backend request failed")
		}
		defer response.Body.Close()
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return fmt.Errorf("backend HTTP %d", response.StatusCode)
		}
		if into != nil {
			return json.NewDecoder(response.Body).Decode(into)
		}
		return nil
	}
	defer func() {
		if failure != nil {
			r.Status, r.Error = "stopped_on_error", failure.Error()
			if active != "" {
				if err := request("POST", "/api/runs/"+active+"/stop", nil); err != nil {
					fmt.Fprintln(os.Stderr, "Cancellation not confirmed:", err)
				}
			}
		}
		data, _ := json.MarshalIndent(r, "", "  ")
		if err := os.WriteFile(filepath.Join(output, "recording.json"), data, 0600); err != nil {
			fmt.Fprintln(os.Stderr, "Cannot save recording manifest")
		}
		fmt.Println("Recording:", output)
	}()
	var health struct {
		Status     string `json:"status"`
		Configured bool   `json:"model_configured"`
	}
	if err := request("GET", "/api/health", &health); err != nil {
		return err
	}
	if health.Status != "ok" || !health.Configured {
		return errors.New("backend/model preflight failed")
	}
	var initial []domain.Run
	if err := request("GET", "/api/runs", &initial); err != nil {
		return err
	}
	known := map[string]bool{}
	for _, run := range initial {
		known[run.ID] = true
		if run.Status == "running" {
			return errors.New("another run is active; recorder did not change it")
		}
	}
	if *language != "en" && *language != "id" {
		return errors.New("lang must be en or id")
	}
	if *reviewBaseline != "" || *reviewRetest != "" {
		if *preview || *geometry || len(*reviewBaseline) != 36 || len(*reviewRetest) != 36 {
			return errors.New("review requires both saved run IDs and cannot use preview/geometry")
		}
		for i, id := range []string{*reviewBaseline, *reviewRetest} {
			var saved domain.Run
			if err := request("GET", "/api/runs/"+id, &saved); err != nil {
				return err
			}
			status := "findings"
			if i == 1 {
				status = "passed"
			}
			if saved.ID != id || saved.Environment != "sandbox" || saved.Status != status || saved.Finished == nil || !saved.TotalVerified || saved.OrderCount != 1 {
				return errors.New("saved run is not verified sandbox evidence")
			}
			for _, step := range saved.Steps {
				if step.Error != "" {
					return errors.New("saved run contains a step error")
				}
			}
			var orders []domain.Order
			if err := request("GET", "/api/runs/"+id+"/orders", &orders); err != nil {
				return err
			}
			if len(orders) != 1 || orders[0].Total != saved.Displayed {
				return errors.New("saved order evidence mismatch")
			}
			r.Runs = append(r.Runs, saved)
		}
		if r.Runs[1].ParentID != r.Runs[0].ID || r.Runs[1].Goal != r.Runs[0].Goal || r.Runs[1].Variant != r.Runs[0].Variant {
			return errors.New("saved baseline/retest do not match")
		}
		r.Mode = "saved_run_review"
		r.Style += "; review of real saved results; no new checkout or inference"
	}
	ctx, stop := context.WithTimeout(context.Background(), 210*time.Second)
	defer stop()
	// This macOS Chrome build reserves 143px of window height for its shell;
	// keep the captured content area at 1920x1080 without rescaling the text.
	allocator, closeAllocator := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WindowSize(1920, 1223))...)
	defer closeAllocator()
	browser, closeBrowser := chromedp.NewContext(allocator)
	defer closeBrowser()
	var runtimeError atomic.Bool
	chromedp.ListenTarget(browser, func(event any) {
		if _, ok := event.(*runtime.EventExceptionThrown); ok {
			runtimeError.Store(true)
		}
	})
	runJS := func(script string) error { return chromedp.Run(browser, chromedp.Evaluate(script, nil)) }
	if err := chromedp.Run(browser, chromedp.EmulateViewport(1920, 1080), chromedp.Navigate("http://127.0.0.1:5173"), chromedp.WaitVisible("#scenario", chromedp.ByQuery)); err != nil {
		return errors.New("dashboard could not load")
	}
	if r.Mode == "saved_run_review" {
		selector := fmt.Sprintf(`button[aria-label*="%s"]`, r.Runs[0].ID[:8])
		encoded, _ := json.Marshal(selector)
		if err := chromedp.Run(browser, chromedp.WaitReady(selector, chromedp.ByQuery), chromedp.Evaluate(`document.querySelector(`+string(encoded)+`).click()`, nil), chromedp.WaitReady(`.finding .text-button`, chromedp.ByQuery)); err != nil {
			return errors.New("saved baseline could not load before recording")
		}
	}
	if err := runJS(`(() => {
	  const cursor=document.createElement('div'); cursor.id='demo-cursor';
	  cursor.style.cssText='position:fixed;left:0;top:0;width:25px;height:32px;z-index:2147483647;pointer-events:none;transform:translate(1800px,100px);filter:drop-shadow(0 2px 2px #0004)';
	  cursor.innerHTML='<svg viewBox="0 0 25 32" xmlns="http://www.w3.org/2000/svg"><path d="M2 2v24l6-6 5 10 4-2-5-10h9Z" fill="white" stroke="#19253a" stroke-width="1.6" stroke-linejoin="round"/></svg>';
	  document.body.append(cursor);
	})()`); err != nil {
		return err
	}
	if *preview || r.Mode == "saved_run_review" {
		if err := runJS(`document.querySelector('form').addEventListener('submit',event=>{event.preventDefault();event.stopImmediatePropagation();},true)`); err != nil {
			return err
		}
	}
	started := time.Now()
	captureCheck, captureFinish, err := startCapture(browser, output, started, &r)
	if err != nil {
		return err
	}
	defer func() {
		if err := captureFinish(); failure == nil && err != nil {
			failure = err
		}
	}()
	mark := func(name, caption string) error {
		r.Chapters = append(r.Chapters, chapter{Name: name, Seconds: time.Since(started).Seconds()})
		return nil
	}
	checkUI := func() error {
		if err := captureCheck(); err != nil {
			return err
		}
		if runtimeError.Load() {
			return errors.New("browser runtime exception")
		}
		var bad bool
		if err := chromedp.Run(browser, chromedp.Evaluate(`!!document.querySelector('.error-banner, .form-error')`, &bad)); err != nil {
			return errors.New("dashboard inspection failed")
		}
		if bad {
			return errors.New("dashboard reports an error")
		}
		if time.Since(started) > 175*time.Second {
			return errors.New("demo duration budget reached")
		}
		return nil
	}
	hold := func(duration time.Duration) error {
		until := time.Now().Add(duration)
		for time.Now().Before(until) {
			if err := checkUI(); err != nil {
				return err
			}
			time.Sleep(90 * time.Millisecond)
		}
		return nil
	}
	if *geometry {
		if err := hold(2 * time.Second); err != nil {
			return err
		}
		r.Status = "geometry-preview-no-run-started"
		return nil
	}
	cursorX, cursorY := 1800.0, 100.0
	move := func(selector string) error {
		var point struct {
			X, Y    float64
			Visible bool
		}
		sel, _ := json.Marshal(selector)
		if err := chromedp.Run(browser, chromedp.Evaluate(`(() => {const e=document.querySelector(`+string(sel)+`); if(!e)return {Visible:false}; const r=e.getBoundingClientRect();return {X:r.x+r.width/2,Y:r.y+r.height/2,Visible:r.y>=0&&r.bottom<=innerHeight};})()`, &point)); err != nil {
			return err
		}
		if !point.Visible {
			return fmt.Errorf("demo target not visible: %s", selector)
		}
		sx, sy := cursorX, cursorY
		for i := 1; i <= 24; i++ {
			t := float64(i) / 24
			eased := (1 - math.Cos(t*math.Pi)) / 2
			cursorX = sx + (point.X-sx)*eased
			cursorY = sy + (point.Y-sy)*eased
			if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error {
				return input.DispatchMouseEvent(input.MouseMoved, cursorX, cursorY).Do(ctx)
			}), chromedp.Evaluate(fmt.Sprintf(`document.querySelector('#demo-cursor').style.transform='translate(%.3fpx,%.3fpx)'`, cursorX, cursorY), nil)); err != nil {
				return err
			}
			time.Sleep(12 * time.Millisecond)
		}
		return hold(150 * time.Millisecond)
	}
	click := func(selector string) error {
		if err := move(selector); err != nil {
			return err
		}
		if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MousePressed, cursorX, cursorY).WithButton(input.Left).WithClickCount(1).Do(ctx)
		})); err != nil {
			return err
		}
		time.Sleep(90 * time.Millisecond)
		if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MouseReleased, cursorX, cursorY).WithButton(input.Left).WithClickCount(1).Do(ctx)
		})); err != nil {
			return err
		}
		return hold(350 * time.Millisecond)
	}
	scroll := func(selector string) error {
		sel, _ := json.Marshal(selector)
		var position struct{ Start, End float64 }
		if err := chromedp.Run(browser, chromedp.Evaluate(`(() => {const e=document.querySelector(`+string(sel)+`);if(!e)throw new Error('Missing scroll target');document.activeElement?.blur();return {Start:scrollY,End:Math.max(0,Math.min(e.getBoundingClientRect().top+scrollY-100,document.documentElement.scrollHeight-innerHeight))};})()`, &position)); err != nil {
			return err
		}
		// Give long journeys through the history time to remain readable.
		steps := max(40, min(150, int(math.Abs(position.End-position.Start)/35)))
		for i := 1; i <= steps; i++ {
			t := float64(i) / float64(steps)
			eased := (1 - math.Cos(t*math.Pi)) / 2
			if err := runJS(fmt.Sprintf(`window.scrollTo({top:%.3f,behavior:'instant'})`, position.Start+(position.End-position.Start)*eased)); err != nil {
				return err
			}
			time.Sleep(22 * time.Millisecond)
		}
		if err := hold(300 * time.Millisecond); err != nil {
			return err
		}
		var visible bool
		if err := chromedp.Run(browser, chromedp.Evaluate(`(() => {const r=document.querySelector(`+string(sel)+`).getBoundingClientRect();return r.top>=0&&r.top<innerHeight-100;})()`, &visible)); err != nil {
			return err
		}
		if !visible {
			return errors.New("smooth scroll target did not enter the viewport")
		}
		return nil
	}
	selectNext := func(selector string) error {
		if err := move(selector); err != nil {
			return err
		}
		value := "alternate"
		if r.Mode == "saved_run_review" {
			value = r.Runs[0].Variant
		}
		if selector == "#language" {
			value = *language
		}
		if err := chromedp.Run(browser, chromedp.Focus(selector, chromedp.ByQuery)); err != nil {
			return err
		}
		if err := hold(400 * time.Millisecond); err != nil {
			return err
		}
		sel, _ := json.Marshal(selector)
		val, _ := json.Marshal(value)
		if err := runJS(`(() => {const e=document.querySelector(` + string(sel) + `);e.value=` + string(val) + `;e.dispatchEvent(new Event('change',{bubbles:true}));})()`); err != nil {
			return err
		}
		return hold(700 * time.Millisecond)
	}
	// hover glides the pointer over an element and pauses, as a tester reading it would.
	hover := func(selector string, pause time.Duration) error {
		if err := move(selector); err != nil {
			return err
		}
		return hold(pause)
	}
	if err := mark("intro", "CheckoutGuard | Sandbox demo — all orders and payments are simulated."); err != nil {
		return err
	}
	if err := hold(2 * time.Second); err != nil {
		return err
	}
	if *language == "id" {
		if err := selectNext("#language"); err != nil {
			return err
		}
	} else {
		// English is the default; show the picker without changing it.
		if err := hover("#language", 900*time.Millisecond); err != nil {
			return err
		}
	}
	for _, selector := range []string{".metrics > :nth-child(1)", ".metrics > :nth-child(2)", ".metrics > :nth-child(3)"} {
		if err := hover(selector, 650*time.Millisecond); err != nil {
			return err
		}
	}
	if err := hold(500 * time.Millisecond); err != nil {
		return err
	}
	if err := scroll(".composer"); err != nil {
		return err
	}
	if err := mark("goal", "A goal drives the agent: verify the total and confirm exactly one saved order."); err != nil {
		return err
	}
	if err := hover(".environment-choice.chosen", 800*time.Millisecond); err != nil {
		return err
	}
	if err := click("#goal"); err != nil {
		return err
	}
	if err := runJS(`document.querySelector('#goal').select()`); err != nil {
		return err
	}
	if err := chromedp.Run(browser, chromedp.KeyEvent(kb.Backspace)); err != nil {
		return err
	}
	goal := "Buy one tote using demo@example.test. Verify the checkout total and confirm exactly one order is saved."
	if *language == "id" {
		goal = "Beli 1 tote menggunakan demo@example.test. Verifikasi total checkout dan pastikan tepat 1 pesanan tersimpan."
	}
	if r.Mode == "saved_run_review" {
		goal = r.Runs[0].Goal
	}
	// Human typing rhythm: uneven key intervals, short pauses between words and
	// a longer one after each sentence. Seeded so retakes look the same.
	rhythm := rand.New(rand.NewSource(7))
	for _, char := range goal {
		if err := chromedp.Run(browser, chromedp.ActionFunc(func(ctx context.Context) error { return input.InsertText(string(char)).Do(ctx) })); err != nil {
			return err
		}
		delay := time.Duration(45+rhythm.Intn(55)) * time.Millisecond
		switch {
		case char == '.':
			delay += time.Duration(260+rhythm.Intn(200)) * time.Millisecond
		case char == ' ':
			delay += time.Duration(30+rhythm.Intn(70)) * time.Millisecond
		case rhythm.Intn(18) == 0:
			delay += time.Duration(120+rhythm.Intn(160)) * time.Millisecond
		}
		time.Sleep(delay)
	}
	if err := hold(800 * time.Millisecond); err != nil {
		return err
	}
	if err := hover("#scenario", 900*time.Millisecond); err != nil {
		return err
	}
	if err := selectNext("#variant"); err != nil {
		return err
	}
	variant := "alternate"
	if r.Mode == "saved_run_review" {
		variant = r.Runs[0].Variant
	}
	var ready bool
	encodedGoal, _ := json.Marshal(goal)
	if err := chromedp.Run(browser, chromedp.Evaluate(`document.querySelector('#language').value===`+fmt.Sprintf("%q", *language)+` && document.querySelector('#variant').value===`+fmt.Sprintf("%q", variant)+` && document.querySelector('#scenario').value==='total_bug' && document.querySelector('#goal').value===`+string(encodedGoal), &ready)); err != nil {
		return err
	}
	if !ready {
		return errors.New("typed goal and form selections were not confirmed")
	}
	if r.Mode == "saved_run_review" {
		showStep := func(number int) error {
			selector := fmt.Sprintf(".timeline-step:nth-child(%d)", number+1)
			encoded, _ := json.Marshal(selector)
			if err := runJS(`(() => {const e=document.querySelector(` + string(encoded) + `),p=document.querySelector('.timeline');p.scrollTo({top:p.scrollTop+e.getBoundingClientRect().top-p.getBoundingClientRect().top-35,behavior:'smooth'});})()`); err != nil {
				return err
			}
			if err := hold(450 * time.Millisecond); err != nil {
				return err
			}
			if err := click(selector); err != nil {
				return err
			}
			var selected bool
			if err := chromedp.Run(browser, chromedp.Evaluate(`document.querySelector(`+string(encoded)+`).getAttribute('aria-pressed')==='true'`, &selected)); err != nil {
				return err
			}
			if !selected {
				return errors.New("saved step selection was not confirmed")
			}
			return hold(1300 * time.Millisecond)
		}
		if err := mark("baseline_review", ""); err != nil {
			return err
		}
		if err := click(`.nav-item[href="#evidence"]`); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := scroll(".execution-grid"); err != nil {
			return err
		}
		for i := range r.Runs[0].Steps {
			if err := showStep(i); err != nil {
				return err
			}
		}
		if err := scroll(".verification-summary"); err != nil {
			return err
		}
		if err := mark("finding", ""); err != nil {
			return err
		}
		if err := hold(4 * time.Second); err != nil {
			return err
		}
		if err := scroll("section.findings"); err != nil {
			return err
		}
		if err := click(".finding .text-button"); err != nil {
			return err
		}
		if err := hold(4 * time.Second); err != nil {
			return err
		}
		if err := click(".modal-close"); err != nil {
			return err
		}
		if err := mark("history", ""); err != nil {
			return err
		}
		if err := click(`.nav-item[href="#run-history"]`); err != nil {
			return err
		}
		if err := hold(1500 * time.Millisecond); err != nil {
			return err
		}
		selector := fmt.Sprintf(`button[aria-label*="%s"]`, r.Runs[1].ID[:8])
		if err := scroll(selector); err != nil {
			return err
		}
		if err := hold(1500 * time.Millisecond); err != nil {
			return err
		}
		if err := click(selector); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := mark("retest_review", ""); err != nil {
			return err
		}
		if err := click(`.nav-item[href="#evidence"]`); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := scroll(".execution-grid"); err != nil {
			return err
		}
		for _, i := range []int{4, 6, 7} {
			if err := showStep(i); err != nil {
				return err
			}
		}
		if err := click(".tool-result summary"); err != nil {
			return err
		}
		if err := hold(1800 * time.Millisecond); err != nil {
			return err
		}
		if err := click(".tool-result summary"); err != nil {
			return err
		}
		if err := scroll(".comparison"); err != nil {
			return err
		}
		if err := mark("comparison", ""); err != nil {
			return err
		}
		if err := hold(5 * time.Second); err != nil {
			return err
		}
		if err := scroll(".verification-summary"); err != nil {
			return err
		}
		if err := mark("outro", ""); err != nil {
			return err
		}
		if err := hold(4 * time.Second); err != nil {
			return err
		}
		r.Status = "verified"
		return nil
	}
	if *preview {
		if err := scroll("section.findings"); err != nil {
			return err
		}
		if err := click(".finding .text-button"); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := click(".modal-close"); err != nil {
			return err
		}
		if err := scroll(".execution-grid"); err != nil {
			return err
		}
		if err := click(".tool-result summary"); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := click(".tool-result summary"); err != nil {
			return err
		}
		if err := click(".screenshot-button"); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		if err := click(".modal-close"); err != nil {
			return err
		}
		if err := scroll("#run-history"); err != nil {
			return err
		}
		if err := click("tbody tr:first-child button"); err != nil {
			return err
		}
		if err := scroll(".verification-summary"); err != nil {
			return err
		}
		if err := hold(time.Second); err != nil {
			return err
		}
		r.Status = "preview-no-run-started"
		return nil
	}
	startRun := func(retest bool) error {
		var selector = "form .primary"
		if retest {
			selector = ".journey .text-button"
		}
		if err := click(selector); err != nil {
			return errors.New("could not click start/retest")
		}
		deadline := time.Now().Add(8 * time.Second)
		for time.Now().Before(deadline) {
			var runs []domain.Run
			if err := request("GET", "/api/runs", &runs); err != nil {
				return err
			}
			for _, run := range runs {
				if !known[run.ID] {
					if run.Environment != "sandbox" {
						return errors.New("unexpected environment")
					}
					active = run.ID
					known[run.ID] = true
					fmt.Println("Started sandbox run:", run.ID)
					return nil
				}
			}
			if err := hold(300 * time.Millisecond); err != nil {
				return err
			}
		}
		return errors.New("start was not confirmed")
	}
	waitRun := func(expectedStatus string) (domain.Run, error) {
		lastSteps := -1
		if err := scroll(".execution-grid"); err != nil {
			return domain.Run{}, err
		}
		for {
			var run domain.Run
			if err := request("GET", "/api/runs/"+active, &run); err != nil {
				return run, err
			}
			for _, step := range run.Steps {
				if step.Error != "" {
					r.Runs = append(r.Runs, run)
					return run, fmt.Errorf("step %d (%s) failed: %s", step.Number, step.Action.Tool, step.Error)
				}
			}
			if len(run.Steps) != lastSteps {
				lastSteps = len(run.Steps)
				fmt.Printf("%s: %s, %d steps\n", run.ID[:8], run.Status, len(run.Steps))
				if lastSteps > 3 {
					if err := runJS(`(() => {const e=document.querySelector('.timeline'); e.scrollTo({top:e.scrollHeight,behavior:'smooth'});})()`); err != nil {
						return run, err
					}
				}
			}
			if run.Status != "running" {
				active = ""
				r.Runs = append(r.Runs, run)
				if run.Status != expectedStatus || !run.TotalVerified || run.OrderCount != 1 {
					return run, errors.New("sandbox result did not pass verification")
				}
				var orders []domain.Order
				if err := request("GET", "/api/runs/"+run.ID+"/orders", &orders); err != nil {
					return run, err
				}
				if len(orders) != 1 || orders[0].Total != run.Displayed {
					return run, errors.New("saved order verification failed")
				}
				return run, nil
			}
			if err := hold(500 * time.Millisecond); err != nil {
				return run, err
			}
		}
	}
	if err := mark("baseline", "Live run: the model chooses browser actions; each completed step saves a screenshot."); err != nil {
		return err
	}
	if err := startRun(false); err != nil {
		return err
	}
	baseline, err := waitRun("findings")
	if err != nil {
		return err
	}
	if baseline.Expected != 9800 || baseline.Displayed != 11050 {
		return errors.New("baseline totals differ from this demo scenario")
	}
	if err := hold(2300 * time.Millisecond); err != nil {
		return err
	}
	if err := scroll(".verification-summary"); err != nil {
		return err
	}
	if err := mark("finding", "Verified finding: 8,750 − 1,250 + 2,300 = 9,800 IDR. The broken fixture displays 11,050 IDR."); err != nil {
		return err
	}
	if err := hold(4 * time.Second); err != nil {
		return err
	}
	if err := scroll("section.findings"); err != nil {
		return err
	}
	if err := click(".finding .text-button"); err != nil {
		return err
	}
	if err := hold(4 * time.Second); err != nil {
		return err
	}
	if err := click(".modal-close"); err != nil {
		return err
	}
	if err := scroll(".execution-grid"); err != nil {
		return err
	}
	if err := mark("retest", "Retest: same goal and basket, against a corrected sandbox fixture. The agent does not edit store source code."); err != nil {
		return err
	}
	if err := startRun(true); err != nil {
		return err
	}
	retest, err := waitRun("passed")
	if err != nil {
		return err
	}
	if retest.ParentID != baseline.ID || retest.Expected != 9800 || retest.Displayed != 9800 {
		return errors.New("retest does not match baseline and expected totals")
	}
	if err := hold(2300 * time.Millisecond); err != nil {
		return err
	}
	if err := scroll(".comparison"); err != nil {
		return err
	}
	if err := mark("comparison", "Before → after: 11,050 → 9,800 IDR. Exactly one order is saved in each synthetic checkout."); err != nil {
		return err
	}
	if err := hold(5 * time.Second); err != nil {
		return err
	}
	if err := mark("history", ""); err != nil {
		return err
	}
	if err := scroll("#run-history"); err != nil {
		return err
	}
	if err := hold(2200 * time.Millisecond); err != nil {
		return err
	}
	if err := click("tbody tr:first-child button"); err != nil {
		return err
	}
	if err := scroll(".execution-grid"); err != nil {
		return err
	}
	if err := click(".tool-result summary"); err != nil {
		return err
	}
	if err := hold(2800 * time.Millisecond); err != nil {
		return err
	}
	if err := click(".tool-result summary"); err != nil {
		return err
	}
	if err := click(".screenshot-button"); err != nil {
		return err
	}
	if err := hold(4 * time.Second); err != nil {
		return err
	}
	if err := click(".modal-close"); err != nil {
		return err
	}
	if err := scroll(".verification-summary"); err != nil {
		return err
	}
	if err := mark("outro", "CheckoutGuard · FLdevBoo | github.com/fadillaoegi/checkoutGuard | Sandbox segment; marketplace demos are separate."); err != nil {
		return err
	}
	if err := hold(4 * time.Second); err != nil {
		return err
	}
	r.Status = "verified"
	return nil
}
