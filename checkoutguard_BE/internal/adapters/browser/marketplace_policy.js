(() => {
  const visible = e => !!e.getClientRects().length && getComputedStyle(e).visibility !== 'hidden';
  const label = e => (e.tagName === 'INPUT' ? e.getAttribute('aria-label') || e.getAttribute('placeholder') || e.name || (/quantity/.test(e.className) ? 'quantity' : '') : e.innerText || e.getAttribute('aria-label') || '').trim().replace(/\s+/g, ' ').toLowerCase();
  const finalLabel = text => /buat pesanan|pesan sekarang|bayar|pay now|place order|submit order|confirm order|konfirmasi pesanan|order now|proses pembayaran/.test(text);
  const quantity = e => e.tagName === 'INPUT' && /quantity|kuantitas|jumlah/.test([e.name, e.className, e.getAttribute('aria-label'), e.getAttribute('placeholder')].join(' ').toLowerCase()) && !/password|email|tel/.test(e.type);
  const variant = e => !!e.closest('[class*="variation"],[class*="variant"],[class*="product-option"],[data-testid*="variant"],[data-testid*="Variant"]');
  const allowed = e => {
    const text = label(e);
    if (finalLabel(text) || e.disabled || !visible(e) || e.tagName === 'A') return false;
    if (e.tagName === 'INPUT') return quantity(e);
    if (!e.matches('button,[role="button"],[role="radio"]')) return false;
    if (variant(e) && text.length > 0 && text.length < 100) return true;
    if (/^(beli sekarang|buy now|add to cart|tambah ke keranjang|masukkan keranjang|masukkan ke keranjang)$/.test(text)) return true;
    return /\/cart(?:\/|$)|\/keranjang(?:\/|$)/.test(location.pathname) && /^checkout(?:\s*\([0-9]+\))?$/.test(text);
  };
  return {allowed, visible, label, finalLabel};
})()
