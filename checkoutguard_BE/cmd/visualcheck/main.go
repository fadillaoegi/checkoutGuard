// visualcheck inspects the local dashboard. It never starts a checkout run.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, timeout := context.WithTimeout(context.Background(), 60*time.Second)
	defer timeout()
	allocator, stop := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WindowSize(1440, 1080))...)
	defer stop()
	browser, cancel := chromedp.NewContext(allocator)
	defer cancel()
	var errors []string
	var mu sync.Mutex
	chromedp.ListenTarget(browser, func(e any) {
		if _, ok := e.(*runtime.EventExceptionThrown); ok {
			mu.Lock()
			errors = append(errors, "browser runtime exception")
			mu.Unlock()
		}
	})
	out := filepath.Join("..", "tmp", "ui")
	if err := os.MkdirAll(out, 0755); err != nil {
		log.Fatal(err)
	}
	capture := func(name string) {
		var png []byte
		if err := chromedp.Run(browser, chromedp.FullScreenshot(&png, 100)); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(out, name), png, 0600); err != nil {
			log.Fatal(err)
		}
	}
	if err := chromedp.Run(browser, chromedp.Navigate("http://127.0.0.1:5173"), chromedp.WaitVisible("tbody tr", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	var idle bool
	if err := chromedp.Run(browser, chromedp.Evaluate(`fetch('/api/runs').then(r => r.json()).then(runs => !runs.some(r => r.status === 'running'))`, &idle, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) })); err != nil {
		log.Fatal(err)
	}
	if !idle {
		log.Fatal("Dashboard QA requires an idle agent; finish the active run before retrying")
	}
	check := func(script, description string) {
		var valid bool
		if err := chromedp.Run(browser, chromedp.Evaluate(script, &valid)); err != nil {
			log.Fatal(err)
		}
		if !valid {
			log.Fatal(description)
		}
	}
	check(`document.documentElement.lang === 'en' && document.querySelector('#language').value === 'en'`, "English must be the default language")
	var hero struct {
		Width       int `json:"width"`
		Height      int `json:"height"`
		CornerAlpha int `json:"cornerAlpha"`
	}
	if err := chromedp.Run(browser, chromedp.Poll(`document.querySelector('.hero-art img')?.complete && document.querySelector('.hero-art img').naturalWidth > 0`, nil), chromedp.Evaluate(`(() => { const img=document.querySelector('.hero-art img'); const canvas=document.createElement('canvas'); canvas.width=img.naturalWidth; canvas.height=img.naturalHeight; const ctx=canvas.getContext('2d'); ctx.drawImage(img,0,0); return {width:img.naturalWidth,height:img.naturalHeight,cornerAlpha:ctx.getImageData(0,0,1,1).data[3]}; })()`, &hero)); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hero asset loaded: %dx%d; corner alpha=%d\n", hero.Width, hero.Height, hero.CornerAlpha)
	var brandOK bool
	if err := chromedp.Run(browser, chromedp.Poll(`document.querySelector('.brand-mark')?.complete && document.querySelector('.brand-mark').naturalWidth === 64`, nil), chromedp.Evaluate(`(async () => { const links=[...document.querySelectorAll('link[rel="icon"],link[rel="apple-touch-icon"]')]; if(links.length!==4) return false; for(const link of links) { const response=await fetch(link.href); if(!response.ok) return false; const url=URL.createObjectURL(await response.blob()); try { const img=new Image(); img.src=url; await img.decode(); const expected=link.type==='image/svg+xml'?64:link.rel==='apple-touch-icon'?180:link.type==='image/png'?32:48; if(img.naturalWidth!==expected || img.naturalHeight!==expected) return false; } finally { URL.revokeObjectURL(url); } } return true; })()`, &brandOK, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) })); err != nil {
		log.Fatal(err)
	}
	if !brandOK {
		log.Fatal("Logo or favicon asset failed browser verification")
	}
	fmt.Println("Logo and SVG/PNG/ICO/Apple favicon assets loaded and decoded in Chrome")
	capture("dashboard-desktop.png")
	var overview []byte
	if err := chromedp.Run(browser, chromedp.CaptureScreenshot(&overview)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "dashboard-overview.png"), overview, 0600); err != nil {
		log.Fatal(err)
	}
	if err := chromedp.Run(browser, chromedp.Click(".timeline-step", chromedp.ByQuery), chromedp.Poll(`document.querySelector('#evidence').innerText.includes('Captured after step 0')`, nil), chromedp.Click(".screenshot-button", chromedp.ByQuery), chromedp.WaitVisible(".modal", chromedp.ByQuery), chromedp.KeyEvent("\u001b"), chromedp.WaitNotPresent(".modal", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	if err := chromedp.Run(browser, chromedp.Evaluate(`Array.from(document.querySelectorAll('tbody tr')).find(row => row.innerText.includes('· retest')).querySelector('button').click()`, nil), chromedp.Poll(`document.body.innerText.includes('Before & after')`, nil)); err != nil {
		log.Fatal(err)
	}
	capture("dashboard-retest.png")
	for _, environment := range []string{"shopee", "tokopedia"} {
		var valid bool
		selector := `.environment-choice:nth-child(2)`
		if environment == "tokopedia" {
			selector = `.environment-choice:nth-child(3)`
		}
		if err := chromedp.Run(browser, chromedp.Click(selector, chromedp.ByQuery), chromedp.WaitVisible(".marketplace-connect", chromedp.ByQuery), chromedp.Evaluate(`document.querySelector('#target').value === '' && !document.querySelector('#scenario') && (document.querySelector('form .primary').disabled || document.querySelector('.marketplace-connect').innerText.includes('browser is open')) && document.querySelector('.marketplace-connect').innerText.includes('sign in manually')`, &valid)); err != nil {
			log.Fatal(err)
		}
		if !valid {
			log.Fatal("Marketplace form did not switch correctly")
		}
		capture("dashboard-" + environment + ".png")
	}
	// Browser-state API fixtures exercise the access notice without opening a
	// real marketplace session or making a checkout request.
	if err := chromedp.Run(browser, chromedp.Click(".environment-choice:nth-child(2)", chromedp.ByQuery), chromedp.Evaluate(`(() => {
      const original = window.fetch;
      window.__cgRestoreFetch = () => { window.fetch = original; };
      window.__cgPageState = 'blocked'; window.__cgPageCheck = 'checked'; window.__cgCheckoutStarts = 0;
      window.fetch = async (input, init) => {
        const path = typeof input === 'string' ? input : input.url;
        if (path === '/api/runs' && init?.method === 'POST') { window.__cgCheckoutStarts++; throw new Error('QA cannot start checkout'); }
        if (path === '/api/marketplaces/shopee/connect') return new Response(JSON.stringify({status:'browser_open',page_state:window.__cgPageState}), {status:200,headers:{'Content-Type':'application/json'}});
        if (path === '/api/marketplaces/session') return new Response(JSON.stringify({environment:'shopee',browser_open:true,page_state:window.__cgPageState,page_check:window.__cgPageCheck}), {status:200,headers:{'Content-Type':'application/json'}});
        return original(input, init);
      };
    })()`, nil), chromedp.Click(".marketplace-connect button", chromedp.ByQuery), chromedp.WaitVisible(".marketplace-access-notice", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('form .primary').disabled && document.querySelector('.marketplace-access-notice').innerText.includes('Marketplace verification required')`, "Blocked marketplace did not stop checkout launch")
	capture("dashboard-marketplace-blocked-english.png")
	if err := chromedp.Run(browser, chromedp.Evaluate(`(() => { const select=document.querySelector('#language'); select.value='id'; select.dispatchEvent(new Event('change',{bubbles:true})); })()`, nil), chromedp.Poll(`document.documentElement.lang === 'id'`, nil)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('.marketplace-access-notice').innerText.includes('Verifikasi marketplace diperlukan') && document.querySelector('.marketplace-access-notice').innerText.includes('Gunakan demo Sandbox')`, "Blocked guidance did not translate")
	if err := chromedp.Run(browser, chromedp.EmulateViewport(390, 844)); err != nil {
		log.Fatal(err)
	}
	check(`document.documentElement.scrollWidth <= window.innerWidth`, "Blocked guidance overflows mobile")
	capture("dashboard-marketplace-blocked-indonesian-mobile.png")
	if err := chromedp.Run(browser, chromedp.Evaluate(`window.__cgPageState = 'browsing'; window.__cgPageCheck = 'unavailable'`, nil), chromedp.Click(".marketplace-connect button", chromedp.ByQuery), chromedp.Poll(`document.querySelector('.marketplace-access-notice')?.innerText.includes('Status browser belum terbaca')`, nil)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('form .primary').disabled && window.__cgCheckoutStarts === 0`, "Unreadable page check incorrectly enabled checkout")
	capture("dashboard-marketplace-check-unavailable.png")
	if err := chromedp.Run(browser, chromedp.Evaluate(`window.__cgPageCheck = 'checked'`, nil), chromedp.Click(".marketplace-connect button", chromedp.ByQuery), chromedp.WaitNotPresent(".marketplace-access-notice", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	check(`!document.querySelector('form .primary').disabled && window.__cgCheckoutStarts === 0`, "Manual browser recheck failed or started checkout")
	if err := chromedp.Run(browser, chromedp.Evaluate(`window.__cgPageState = 'blocked'`, nil), chromedp.Click(".marketplace-connect button", chromedp.ByQuery), chromedp.WaitVisible(".marketplace-access-notice", chromedp.ByQuery), chromedp.Click(".marketplace-access-notice button", chromedp.ByQuery), chromedp.WaitVisible("#scenario", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	check(`window.__cgCheckoutStarts === 0 && !document.querySelector('.marketplace-access-notice')`, "Sandbox selection launched a checkout")
	if err := chromedp.Run(browser, chromedp.Evaluate(`(() => { window.__cgRestoreFetch(); const select=document.querySelector('#language'); select.value='en'; select.dispatchEvent(new Event('change',{bubbles:true})); })()`, nil), chromedp.EmulateViewport(1440, 1080), chromedp.Poll(`document.documentElement.lang === 'en'`, nil)); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marketplace notice QA passed with synthetic API states: blocked launch, English/Indonesian guidance, manual recheck, and sandbox selection; no marketplace session or checkout created")
	if err := chromedp.Run(browser, chromedp.Click(".environment-choice:first-child", chromedp.ByQuery), chromedp.WaitVisible("#scenario", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	if err := chromedp.Run(browser, chromedp.Evaluate(`(() => { const goal=document.querySelector('#goal'); Object.getOwnPropertyDescriptor(HTMLTextAreaElement.prototype,'value').set.call(goal,'Siapkan dua tote. Saya akan meninjau pesanan.'); goal.dispatchEvent(new Event('input',{bubbles:true})); const select=document.querySelector('#language'); select.value='id'; select.dispatchEvent(new Event('change',{bubbles:true})); })()`, nil), chromedp.Poll(`document.documentElement.lang === 'id'`, nil)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('#goal').value === 'Siapkan dua tote. Saya akan meninjau pesanan.' && localStorage.getItem('checkoutguard.language') === 'id' && document.querySelector('form .primary').innerText.includes('Jalankan checkout') && document.querySelector('#scenario option').textContent.includes('Checkout benar')`, "Language switch failed or changed the user's goal")
	capture("dashboard-indonesian-desktop.png")
	if err := chromedp.Run(browser, chromedp.Reload(), chromedp.WaitVisible("tbody tr", chromedp.ByQuery), chromedp.Poll(`document.documentElement.lang === 'id'`, nil)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('#language').value === 'id' && document.title.includes('Rilis dengan bukti') && document.querySelector('#goal').value.startsWith('Beli satu tote')`, "Saved language did not survive reload")
	if err := chromedp.Run(browser, chromedp.Click(".timeline-step", chromedp.ByQuery), chromedp.Poll(`document.querySelector('#evidence').innerText.includes('Diambil setelah langkah 0')`, nil), chromedp.Click(".screenshot-button", chromedp.ByQuery), chromedp.WaitVisible(".modal", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('.modal').getAttribute('aria-label') === 'Bukti screenshot' && document.querySelector('.modal-close').innerText.includes('Tutup')`, "Evidence dialog did not translate")
	if err := chromedp.Run(browser, chromedp.KeyEvent("\u001b"), chromedp.WaitNotPresent(".modal", chromedp.ByQuery)); err != nil {
		log.Fatal(err)
	}
	for _, environment := range []string{"shopee", "tokopedia"} {
		selector := `.environment-choice:nth-child(2)`
		if environment == "tokopedia" {
			selector = `.environment-choice:nth-child(3)`
		}
		if err := chromedp.Run(browser, chromedp.Click(selector, chromedp.ByQuery), chromedp.WaitVisible(".marketplace-connect", chromedp.ByQuery)); err != nil {
			log.Fatal(err)
		}
		check(`document.querySelector('form .primary').innerText.includes('Siapkan checkout') && document.querySelector('.marketplace-connect').innerText.includes('login secara manual') && document.querySelector('#goal').value.startsWith('Siapkan checkout')`, "Marketplace form did not translate")
		capture("dashboard-" + environment + "-indonesian.png")
	}
	if err := chromedp.Run(browser, chromedp.Click(".environment-choice:first-child", chromedp.ByQuery), chromedp.EmulateViewport(390, 844)); err != nil {
		log.Fatal(err)
	}
	check(`document.documentElement.scrollWidth <= window.innerWidth`, "Indonesian mobile layout overflows")
	capture("dashboard-indonesian-mobile.png")
	if err := chromedp.Run(browser, chromedp.EmulateViewport(320, 844)); err != nil {
		log.Fatal(err)
	}
	check(`document.documentElement.scrollWidth <= window.innerWidth && document.querySelector('#language').getBoundingClientRect().right <= window.innerWidth`, "Small mobile language picker overflows")
	if err := chromedp.Run(browser, chromedp.Evaluate(`(() => { const select=document.querySelector('#language'); select.value='en'; select.dispatchEvent(new Event('change',{bubbles:true})); })()`, nil), chromedp.Poll(`document.documentElement.lang === 'en'`, nil)); err != nil {
		log.Fatal(err)
	}
	check(`document.querySelector('form .primary').innerText.includes('Run checkout') && localStorage.getItem('checkoutguard.language') === 'en'`, "Switching back to English failed")
	var horizontalOverflow bool
	if err := chromedp.Run(browser, chromedp.EmulateViewport(390, 844), chromedp.Evaluate(`document.documentElement.scrollWidth>window.innerWidth`, &horizontalOverflow)); err != nil {
		log.Fatal(err)
	}
	capture("dashboard-mobile.png")
	mu.Lock()
	defer mu.Unlock()
	if len(errors) > 0 || horizontalOverflow {
		log.Fatalf("UI verification failed: runtime errors=%d mobile overflow=%t", len(errors), horizontalOverflow)
	}
	fmt.Println("Dashboard QA passed: desktop/retest/mobile, English/Indonesian switching and reload persistence, preserved custom goal, translated evidence and marketplace forms; no runtime exceptions or mobile page overflow; no run started")
}
