(() => {
  const path = location.pathname.toLowerCase();
  const text = (document.body?.innerText || '').toLowerCase();
  // Some challenge failure pages never display the word CAPTCHA in the body.
  if (/\/(?:verify(?:\/|$)|captcha(?:\/|$)|challenge(?:\/|$))/.test(path) ||
      /captcha|verify you are human|verifikasi keamanan|access denied|akses ditolak|unusual traffic|verification can['’]t be completed|halaman tidak ditemukan|page not found|produk tidak ditemukan/.test(text)) return 'blocked';
  const password = [...document.querySelectorAll('input[type="password"]')].some(e => e.getClientRects().length > 0);
  if (password || /\/login(?:\/|$)|\/buyer\/login/.test(path)) return 'needs_login';
  return 'browsing';
})()
