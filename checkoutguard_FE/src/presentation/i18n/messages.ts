import type { Language } from '../../domain/language'

// English copy is also the typed key, so the default language stays explicit.
export const indonesian = {
  Language: 'Bahasa',
  'Cost unavailable': 'Biaya belum tersedia',
  'At least {amount} (partial)': 'Minimal {amount} (belum lengkap)',
  'Browser status unavailable': 'Status browser belum terbaca',
  'CheckoutGuard could not read the current page. Wait until the browser responds, then click Check browser again. Checkout remains paused until the page can be checked.':
    'CheckoutGuard belum dapat membaca halaman saat ini. Tunggu hingga browser merespons, lalu klik Periksa browser lagi. Checkout tetap ditahan sampai halaman bisa diperiksa.',
  'Marketplace verification required': 'Verifikasi marketplace diperlukan',
  'The marketplace is showing an access verification or unavailable page. If it says “Try again later”, wait a few minutes and try manually in the same browser. Then click Check browser. CheckoutGuard cannot complete this verification.':
    'Marketplace menampilkan verifikasi akses atau halaman tidak tersedia. Jika muncul “Try again later”, tunggu beberapa menit lalu coba secara manual di browser yang sama. Setelah itu klik Periksa browser. CheckoutGuard tidak dapat menyelesaikan verifikasi ini.',
  'Sign in manually in the open marketplace browser, then click Check browser. CheckoutGuard does not enter your login details.':
    'Login secara manual di browser marketplace yang terbuka, lalu klik Periksa browser. CheckoutGuard tidak mengisi data login kamu.',
  'After checking the browser, page access is reassessed. This does not verify your login or guarantee checkout will work.':
    'Setelah browser diperiksa, status akses halaman dinilai ulang. Ini belum memverifikasi login atau menjamin checkout berhasil.',
  'CheckoutGuard home': 'Beranda CheckoutGuard',
  'neatHack workspace': 'Ruang kerja neatHack',
  VERIFICATION: 'VERIFIKASI',
  Overview: 'Ringkasan',
  'Run history': 'Riwayat pemeriksaan',
  Evidence: 'Bukti',
  'BUILT WITH': 'DIBANGUN DENGAN',
  'Workflow traces': 'Trace alur kerja',
  'Coding capture': 'Rekaman proses coding',
  'Business plan': 'Rencana bisnis',
  'Open business plan': 'Buka rencana bisnis',
  'Local development': 'Pengembangan lokal',
  'Checkout verification': 'Verifikasi checkout',
  'Synthetic sandbox': 'Sandbox simulasi',
  'Prepare checkout': 'Siapkan checkout',
  'CHECKOUT QUALITY ASSURANCE': 'PEMERIKSAAN KUALITAS CHECKOUT',
  'Check your checkout.': 'Periksa checkout kamu.',
  'See the evidence.': 'Lihat buktinya.',
  'Verify totals, confirm orders, and investigate failures with every browser step saved.':
    'Verifikasi total, pastikan pesanan, dan telusuri kegagalan dengan bukti setiap langkah browser.',
  'Price verification': 'Verifikasi harga',
  'Order confirmation': 'Konfirmasi pesanan',
  'Safe recovery': 'Pemulihan aman',
  'neatHack · Local MVP': 'neatHack · MVP lokal',
  'Backend unavailable: {error}. Start the local backend to run checkout checks.':
    'Backend tidak tersedia: {error}. Jalankan backend lokal untuk memeriksa checkout.',
  'Saved checks': 'Pemeriksaan tersimpan',
  'All runs in this workspace': 'Semua pemeriksaan di ruang kerja ini',
  'Checks passed': 'Pemeriksaan lulus',
  'Correct total · One order confirmed': 'Total benar · Satu pesanan dikonfirmasi',
  'Verified defects': 'Bug terverifikasi',
  'Verified findings across checks': 'Temuan terverifikasi dari semua pemeriksaan',
  'Agent tool calls': 'Pemanggilan tool agent',
  'Browser actions and verification': 'Aksi browser dan verifikasi',
  'Selected check summary': 'Ringkasan pemeriksaan yang dipilih',
  'SELECTED CHECK': 'PEMERIKSAAN TERPILIH',
  'Expected total': 'Total seharusnya',
  'Not checked': 'Belum diperiksa',
  'Marketplace breakdown not verified': 'Rincian harga marketplace belum diverifikasi',
  'Calculated from the price breakdown': 'Dihitung dari rincian harga',
  'Displayed total': 'Total yang ditampilkan',
  'Matches the expected amount': 'Sesuai dengan total seharusnya',
  'Does not match the expected amount': 'Tidak sesuai dengan total seharusnya',
  'Awaiting total verification': 'Menunggu verifikasi total',
  'Persisted orders': 'Pesanan tersimpan',
  Manual: 'Manual',
  '{count} order(s)': '{count} pesanan',
  'No order submitted by the agent': 'Agent tidak membuat pesanan',
  'Exactly one order recorded': 'Tepat satu pesanan tercatat',
  'No order recorded': 'Tidak ada pesanan tercatat',
  'Multiple orders recorded': 'Lebih dari satu pesanan tercatat',
  'Run duration': 'Durasi pemeriksaan',
  sec: 'detik',
  'Stopped before completion': 'Berhenti sebelum selesai',
  'Finished check': 'Pemeriksaan selesai',
  'Check in progress': 'Pemeriksaan sedang berlangsung',
  'Recovered from {count} step issue(s). Final checkout checks passed.':
    'Berhasil pulih dari {count} kendala langkah. Pemeriksaan akhir checkout lulus.',
  'Review checkout in your browser': 'Tinjau checkout di browser kamu',
  'Keep the marketplace result, continue with the demo':
    'Simpan hasil marketplace, lanjutkan dengan demo',
  'Use Sandbox demo': 'Gunakan demo Sandbox',
  'This selects the demo form. It does not start a run or replace the marketplace result.':
    'Tombol ini memilih form demo. Pemeriksaan belum dimulai dan hasil marketplace tetap tersimpan.',
  'Agent activity': 'Aktivitas agent',
  'Manual final confirmation': 'Konfirmasi terakhir secara manual',
  '{basket} basket': 'Keranjang {basket}',
  'Your next checkout journey starts here': 'Mulai pemeriksaan checkout berikutnya di sini',
  'Your check starts with a goal.': 'Pemeriksaan dimulai dari sebuah tujuan.',
  'Start a run to see the agent observe, act, and verify the checkout.':
    'Mulai pemeriksaan untuk melihat agent mengamati, bertindak, dan memverifikasi checkout.',
  '{count} saved step(s)': '{count} langkah tersimpan',
  'Stop run': 'Hentikan pemeriksaan',
  'Load this goal again': 'Muat tujuan ini lagi',
  'Retest corrected store': 'Uji ulang toko yang diperbaiki',
  Initial: 'Awal',
  'Agent is choosing its next action…': 'Agent sedang memilih aksi berikutnya…',
  'What the agent found': 'Temuan agent',
  'Findings backed by amounts and persisted order data.':
    'Temuan didukung oleh rincian nilai dan data pesanan tersimpan.',
  '{count} finding(s)': '{count} temuan',
  'Deterministically verified': 'Terverifikasi melalui pemeriksaan data',
  'Needs verification': 'Perlu diverifikasi',
  'View evidence': 'Lihat bukti',
  'Before & after': 'Sebelum & sesudah',
  'The same goal and basket, against the corrected synthetic store.':
    'Tujuan dan keranjang yang sama, diuji pada toko simulasi yang diperbaiki.',
  RETEST: 'UJI ULANG',
  BASELINE: 'HASIL AWAL',
  '{defects} defect(s) · {orders} order(s)': '{defects} bug · {orders} pesanan',
  'Expected {amount}': 'Seharusnya {amount}',
  'Displayed {amount}': 'Ditampilkan {amount}',
  'Open final evidence': 'Buka bukti akhir',
  'neatlogs trace': 'Trace neatlogs',
  'Copy ID': 'Salin ID',
  'CheckoutGuard · Evidence before release.': 'CheckoutGuard · Bukti sebelum rilis.',
  'Sandbox demo · Marketplace checkout preparation is experimental':
    'Demo Sandbox · Persiapan checkout marketplace masih eksperimental',
  'Screenshot evidence': 'Bukti screenshot',
  'Close ×': 'Tutup ×',
  'Full browser screenshot evidence': 'Bukti screenshot browser lengkap',
  'Start a checkout check': 'Mulai pemeriksaan checkout',
  'Describe the purchase. The agent will act, verify, and save the evidence.':
    'Jelaskan pembelianmu. Agent akan bertindak, memverifikasi, dan menyimpan bukti.',
  'One active check at a time': 'Satu pemeriksaan aktif dalam satu waktu',
  'Checkout environment': 'Lingkungan checkout',
  'Tested demo · Simulated orders': 'Demo teruji · Pesanan simulasi',
  'Experimental · Prepare checkout': 'Eksperimental · Siapkan checkout',
  '{marketplace} browser is open': 'Browser {marketplace} terbuka',
  'Sign in using your own browser session': 'Login melalui sesi browser kamu sendiri',
  'Open the dedicated browser, sign in manually, then paste a product URL below. The agent stops before the final order is submitted.':
    'Buka browser khusus, login secara manual, lalu tempel URL produk di bawah. Agent berhenti sebelum pesanan dikirim.',
  'Browser open does not mean sign-in has been verified. Access challenges require your action.':
    'Browser yang terbuka belum membuktikan status login. Kendala akses perlu kamu tangani sendiri.',
  'Opening…': 'Membuka…',
  'Check browser': 'Periksa browser',
  'Open {marketplace} browser': 'Buka browser {marketplace}',
  'What should the agent do?': 'Apa yang harus dilakukan agent?',
  'Sandbox URL': 'URL Sandbox',
  '{marketplace} product URL': 'URL produk {marketplace}',
  'Sandbox condition': 'Kondisi Sandbox',
  Basket: 'Keranjang',
  Standard: 'Standar',
  'Alternate prices': 'Harga alternatif',
  'Starting…': 'Memulai…',
  'Agent running': 'Agent berjalan',
  'Run checkout': 'Jalankan checkout',
  'All orders and payments are simulated.': 'Semua pesanan dan pembayaran adalah simulasi.',
  'Final order confirmation is manual. Marketplace totals and orders are not independently verified.':
    'Konfirmasi pesanan terakhir dilakukan manual. Total dan pesanan marketplace belum diverifikasi secara independen.',
  'AI model configured': 'Model AI sudah dikonfigurasi',
  'AI model not configured': 'Model AI belum dikonfigurasi',
  'Configure model credentials in the backend to start a run.':
    'Atur kredensial model di backend untuk memulai pemeriksaan.',
  'Latest runs': 'Pemeriksaan terbaru',
  'A history of decisions, outcomes, and evidence.': 'Riwayat keputusan, hasil, dan bukti.',
  '{count} saved': '{count} tersimpan',
  'RUN / CONDITION': 'PEMERIKSAAN / KONDISI',
  STARTED: 'DIMULAI',
  STATUS: 'STATUS',
  ORDERS: 'PESANAN',
  STEPS: 'LANGKAH',
  COST: 'BIAYA',
  retest: 'uji ulang',
  'Inspect run {id}': 'Lihat pemeriksaan {id}',
  'No runs yet. Your first checkout check will appear here.':
    'Belum ada pemeriksaan. Pemeriksaan checkout pertama akan muncul di sini.',
  'Browser evidence': 'Bukti browser',
  'Captured after step {number}': 'Diambil setelah langkah {number}',
  'Real screenshots from the selected environment': 'Screenshot nyata dari lingkungan yang dipilih',
  'Open full screenshot': 'Buka screenshot lengkap',
  'Waiting for a checkout run': 'Menunggu pemeriksaan checkout',
  'Browser evidence after step {number}': 'Bukti browser setelah langkah {number}',
  'Observe. Act. Verify.': 'Amati. Bertindak. Verifikasi.',
  'Every step leaves a screenshot.': 'Setiap langkah menyimpan screenshot.',
  'Order verification': 'Verifikasi pesanan',
  Unavailable: 'Tidak tersedia',
  'Actual tokens': 'Token aktual',
  'Estimated API cost': 'Estimasi biaya API',
  'Step {number} · tool result': 'Langkah {number} · hasil tool',
  'Returned model: {model}': 'Model respons: {model}',
  'Waiting for provider': 'Menunggu provider',
  'Rates not set': 'Tarif belum diatur',
  'Correct checkout': 'Checkout benar',
  'Incorrect total': 'Total salah',
  'Payment timeout': 'Timeout pembayaran',
  '{marketplace} checkout preparation': 'Persiapan checkout {marketplace}',
  'Observe browser': 'Amati browser',
  'Click page control': 'Klik kontrol halaman',
  'Fill checkout field': 'Isi kolom checkout',
  'Verify checkout total': 'Verifikasi total checkout',
  'Check saved orders': 'Periksa pesanan tersimpan',
  'Confirm result': 'Konfirmasi hasil',
  'Model request issue': 'Kendala request model',
  'Ready for review': 'Siap ditinjau',
  'Sign-in needed': 'Perlu login',
  Blocked: 'Terhambat',
  Passed: 'Lulus',
  'Defect found': 'Bug ditemukan',
  Running: 'Berjalan',
  Failed: 'Gagal',
  Stopped: 'Dihentikan',
  'CheckoutGuard · Release with evidence': 'CheckoutGuard · Rilis dengan bukti',
} as const

export type TranslationKey = keyof typeof indonesian
export type TranslationParams = Record<string, string | number>
export function translate(
  language: Language,
  key: TranslationKey,
  params: TranslationParams = {},
): string {
  let text: string = language === 'id' ? indonesian[key] : key
  if (language === 'en') {
    text = text.replace(
      /(\{(\w+)\}[^{}]*?)\(s\)/g,
      (_, prefix: string, name: string) => prefix + (params[name] === 1 ? '' : 's'),
    )
  }
  return text.replace(/\{(\w+)\}/g, (placeholder, name: string) =>
    String(params[name] ?? placeholder),
  )
}

// Translate known application messages at display time, preserving original records.
const applicationMessages: Record<string, string> = {
  'resolve marketplace access manually, then check the browser before starting':
    'Tangani akses marketplace secara manual, lalu periksa browser sebelum memulai',
  'Request failed': 'Request gagal',
  'Failed to fetch': 'Tidak dapat terhubung ke backend',
  'NetworkError when attempting to fetch resource.': 'Tidak dapat terhubung ke backend',
  'Backend unavailable': 'Backend tidak tersedia',
  'Could not start run': 'Tidak dapat memulai pemeriksaan',
  'Could not open browser': 'Tidak dapat membuka browser',
  'Could not stop run': 'Tidak dapat menghentikan pemeriksaan',
  'Choose a marketplace to open its browser': 'Pilih marketplace untuk membuka browsernya',
  'The page did not finish loading. Retry it in the open browser before starting.':
    'Halaman belum selesai dimuat. Coba muat ulang di browser yang terbuka sebelum memulai.',
  'model response timed out': 'Respons model melewati batas waktu',
  'model DNS lookup failed': 'Pencarian DNS model gagal',
  'model connection failed': 'Koneksi model gagal',
  'model request cancelled or run duration exhausted':
    'Request model dibatalkan atau batas durasi pemeriksaan habis',
  'Model recovery budget exhausted; provider or response remains unavailable':
    'Batas pemulihan model habis; provider atau respons masih tidak tersedia',
  'Provider requires a longer wait; start another run later':
    'Provider meminta waktu tunggu lebih lama; mulai pemeriksaan baru nanti',
  'Run cancelled or duration budget exhausted': 'Pemeriksaan dibatalkan atau batas durasi habis',
  'Action budget exhausted': 'Batas jumlah aksi habis',
  'Checkout verified; exactly one persisted synthetic order':
    'Checkout terverifikasi; tepat satu pesanan simulasi tersimpan',
  'Verified pricing defect; exactly one persisted synthetic order':
    'Bug harga terverifikasi; tepat satu pesanan simulasi tersimpan',
  'Checkout review page reached; final order confirmation remains manual':
    'Halaman review checkout tercapai; konfirmasi terakhir tetap manual',
  'Sign in manually in the marketplace browser, then start another run':
    'Login manual di browser marketplace, lalu mulai pemeriksaan baru',
  'Marketplace access challenge or unavailable page; resolve it manually or select Sandbox demo':
    'Ada kendala akses marketplace atau halaman tidak tersedia; tangani secara manual atau pilih demo Sandbox',
  'Agent could not verify a checkout review page; no order submitted':
    'Agent belum dapat memastikan halaman review checkout; tidak ada pesanan dikirim',
  'Initial browser observation': 'Observasi awal browser',
  'Observe marketplace product page': 'Amati halaman produk marketplace',
  'Model response failed validation or transport': 'Validasi respons atau koneksi model gagal',
  'Model request failed': 'Request model gagal',
  'Checkout total does not match the breakdown': 'Total checkout tidak sesuai dengan rincian harga',
  'Payment timeout recovered without a duplicate':
    'Timeout pembayaran dipulihkan tanpa pesanan ganda',
  'Persisted order confirmed using the original correlation ID before retry':
    'Pesanan tersimpan dikonfirmasi menggunakan ID korelasi awal sebelum retry',
  'Duplicate order detected': 'Pesanan ganda terdeteksi',
  'Multiple persisted orders for a single checkout':
    'Lebih dari satu pesanan tersimpan untuk satu checkout',
  'Browser evidence became unavailable': 'Bukti browser tidak lagi tersedia',
  'Cannot persist run start': 'Tidak dapat menyimpan awal pemeriksaan',
  'Cannot persist initial evidence': 'Tidak dapat menyimpan bukti awal',
  'Cannot persist step and evidence': 'Tidak dapat menyimpan langkah dan bukti',
  'run is not active': 'Pemeriksaan tidak aktif',
  'Another run is active': 'Pemeriksaan lain sedang aktif',
}
export function localizeMessage(language: Language, value: string): string {
  if (language === 'en') return value
  if (applicationMessages[value]) return applicationMessages[value]
  const mismatch = /^Expected (\d+) (\w+); checkout displays (\d+) (\w+)$/.exec(value)
  if (mismatch)
    return `Seharusnya ${mismatch[1]} ${mismatch[2]}; checkout menampilkan ${mismatch[3]} ${mismatch[4]}`
  const http = /^model HTTP (\d+)$/.exec(value)
  if (http) return `Model mengembalikan HTTP ${http[1]}`
  const fatal = 'Model request cannot be retried: '
  if (value.startsWith(fatal))
    return (
      'Request model tidak dapat diulang: ' + localizeMessage(language, value.slice(fatal.length))
    )
  return value
}
