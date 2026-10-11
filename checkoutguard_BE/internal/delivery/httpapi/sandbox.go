package httpapi

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var sandboxTemplate = template.Must(template.New("store").Parse(`<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Field Supply — synthetic checkout</title><style>
*{box-sizing:border-box}body{font-family:system-ui;background:#f4f2ed;color:#252925;margin:0;padding:50px}header{display:flex;justify-content:space-between;border-bottom:1px solid #d9dcd2;padding-bottom:22px}.brand{font-size:25px;font-weight:750}.badge{font-size:12px;color:#64735a;background:#e7ebdf;border-radius:30px;padding:9px 15px}main{max-width:760px;margin:50px auto}h1{font-size:42px;letter-spacing:-2px}section{background:white;border:1px solid #e0e1d9;border-radius:18px;padding:28px;margin:20px 0}button{background:#344d3c;color:white;border:0;border-radius:10px;padding:15px 22px;font-weight:650;cursor:pointer;margin:8px 8px 0 0}button:disabled{opacity:.4}input{border:1px solid #ccc;border-radius:8px;padding:12px;margin:8px;max-width:300px}label{display:block}.product{font-size:68px;background:#e8ebdf;text-align:center;padding:32px;border-radius:12px}.row{display:flex;justify-content:space-between;padding:10px 0}.total{border-top:1px solid #ddd;font-weight:750;font-size:23px}.muted{color:#6b766b;font-size:13px}.notice{background:#fff4de;border:1px solid #eac98a;border-radius:8px;padding:16px}code{word-break:break-all;font-size:12px}</style></head><body><header><span class="brand">Field Supply.</span><span class="badge">Synthetic store · No real payment</span></header><main><p class="muted">GOOD THINGS, EVERY DAY</p><h1>The everyday carry.</h1><div id="product"><section><div class="product">◒</div><h2>Everyday canvas tote</h2><p>Made for the small adventures. One product, one simple checkout.</p><p>Unit price: {{.Unit}} IDR (integer units)</p><label>Quantity <input id="quantity" type="number" min="1" max="5" value="1" aria-label="Quantity"></label><button id="add">Add to bag</button></section></div><div id="checkout-panel" hidden><section><h2>Checkout</h2><div id="quote"></div><label>Demo email <input id="email" placeholder="demo@example.test" aria-label="Demo email"></label><p class="muted">Demo payment · Nothing will be charged</p><p class="muted">Order correlation <code id="correlation">{{.Correlation}}</code></p><button id="place">Place demo order</button><button id="status">Check order status</button><p id="message" aria-live="polite"></p></section></div><p class="muted">CheckoutGuard sandbox. All prices and orders are simulated.</p></main><script>
const base=location.pathname;let quantity=1;
const el=id=>document.getElementById(id);
const money=n=>n.toLocaleString('en-US')+' IDR';
async function request(path,options){document.body.dataset.busy='true';try{const response=await fetch(base+path,options);const data=await response.json();if(!response.ok)throw new Error(data.error||'Request failed');return data}finally{document.body.dataset.busy='false'}}
el('add').onclick=async()=>{try{quantity=Number(el('quantity').value);const q=await request('/quote?quantity='+quantity);el('product').hidden=true;el('checkout-panel').hidden=false;el('quote').setAttribute('data-quote',JSON.stringify(q));el('quote').innerHTML='<div class="row"><span>Unit price × quantity</span><span>'+money(q.unit_minor)+' × '+q.quantity+'</span></div><div class="row"><span>Subtotal</span><span>'+money(q.subtotal_minor)+'</span></div><div class="row"><span>Discount</span><span>− '+money(q.discount_minor)+'</span></div><div class="row"><span>Shipping</span><span>'+money(q.shipping_minor)+'</span></div><div class="row total"><span>Total</span><span>'+money(q.total_minor)+'</span></div>'}catch(e){alert(e.message)}};
el('place').onclick=async()=>{if(el('email').value!=='demo@example.test'){el('message').textContent='Enter demo@example.test to continue';return}el('place').disabled=true;try{const order=await request('/orders',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({quantity,correlation_id:el('correlation').textContent})});el('message').textContent='Order confirmed: '+order.id+' · '+money(order.total_minor)}catch(e){el('message').className='notice';el('message').textContent=e.message+' Correlation: '+el('correlation').textContent;el('place').disabled=false}};
el('status').onclick=async()=>{try{const orders=await request('/orders');el('message').className='';el('message').textContent=orders.length?'Order confirmed: '+orders[0].id+' · '+money(orders[0].total_minor)+' · '+orders.length+' order(s)':'No order found for this correlation';if(orders.length)el('place').disabled=true}catch(e){el('message').textContent=e.message}};
</script></body></html>`))

func (a *Handler) sandbox(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	run, err := a.Service.SandboxRun(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	resource := r.PathValue("resource")
	if resource == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; img-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
		quote, err := a.Service.SandboxQuote(r.Context(), id, 1)
		if err != nil {
			handleStartError(w, err)
			return
		}
		sandboxTemplate.Execute(w, map[string]any{"Unit": quote.Unit, "Correlation": run.Correlation})
		return
	}
	switch resource {
	case "quote":
		n, err := strconv.ParseInt(r.URL.Query().Get("quantity"), 10, 64)
		if err != nil {
			writeError(w, 400, "quantity must be 1..5")
			return
		}
		quote, err := a.Service.SandboxQuote(r.Context(), id, n)
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 200, quote)
	case "orders":
		if r.Method == "GET" {
			orders, err := a.Service.Orders(r.Context(), id)
			if err != nil {
				writeError(w, 500, "order lookup failed")
				return
			}
			writeJSON(w, 200, orders)
			return
		}
		if r.Method != "POST" {
			writeError(w, 405, "method not allowed")
			return
		}
		var input struct {
			Quantity    int64  `json:"quantity"`
			Correlation string `json:"correlation_id"`
		}
		if err := decodeJSON(w, r, &input); err != nil {
			writeError(w, 400, "invalid synthetic order")
			return
		}
		order, err := a.Service.PlaceSandboxOrder(r.Context(), id, input.Quantity, input.Correlation)
		if err != nil {
			handleStartError(w, err)
			return
		}
		writeJSON(w, 201, order)
	default:
		if !strings.Contains(resource, "favicon") {
			http.NotFound(w, r)
		} else {
			w.WriteHeader(204)
		}
	}
}
