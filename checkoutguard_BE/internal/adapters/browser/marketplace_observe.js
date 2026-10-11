(() => {
  const {allowed, visible, label, finalLabel} = policy;
  document.querySelectorAll('[data-checkoutguard-control]').forEach(e => e.removeAttribute('data-checkoutguard-control'));
  const generation = crypto.randomUUID();
  const all = [...document.querySelectorAll('button,[role="button"],[role="radio"],input')].filter(visible);
  const controls = all.filter(allowed).slice(0, 60).map((e, i) => {
    const id = generation + '-' + i;
    e.setAttribute('data-checkoutguard-control', id);
    return {selector: '[data-checkoutguard-control="' + id + '"]', label: label(e), type: e.tagName === 'INPUT' ? 'input' : 'button'};
  });
  const redact = text => text.replace(/[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}/gi, '[email hidden]').replace(/(?:\+62|\b08)[0-9 ()-]{7,16}/g, '[phone hidden]');
  const clone = document.body.cloneNode(true);
  clone.querySelectorAll('script,style,input,textarea,address,[class*="address"],[data-testid*="Address"],[data-testid*="address"]').forEach(e => e.remove());
  const text = redact(clone.textContent || '').replace(/\s+/g, ' ').trim().slice(0, 12000);
  const bodyText = (document.body.innerText || '').toLowerCase();
  const finalVisible = all.some(e => finalLabel(label(e)));
  let state = pageState;
  if (state === 'browsing' && /\/checkout(?:\/|$)/.test(location.pathname) && finalVisible && /total|ringkasan belanja|rincian pembayaran/.test(bodyText)) state = 'review';
  if (!document.getElementById('checkoutguard-private-fields')) {
    const style = document.createElement('style'); style.id = 'checkoutguard-private-fields';
    style.textContent = 'input:not([data-checkoutguard-control]),textarea,address,[class*="address"],[data-testid*="Address"],[data-testid*="address"]{filter:blur(12px)!important}';
    document.head.appendChild(style);
  }
  return {url: location.href, text, controls: state === 'blocked' || state === 'needs_login' ? [] : controls, state, summary: state === 'review' ? text.slice(-1800) : ''};
})()
