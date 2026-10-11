# CheckoutGuard

Panduan setup publik dan status terbaru tersedia di [README English](README.md).
Paket screenshot/trace yang dapat ditinjau ada di [bukti reviewer](submission/EVIDENCE.md).
Versi ini mempertahankan penjelasan Indonesia dan catatan awal pengembangan.

Agent pemeriksa checkout untuk neatHack, tim FLdevBoo. Pengguna memberi tujuan;
model mengamati browser, memilih aksi berschema, lalu memeriksa total dan order
yang tersimpan. Backend Go, frontend React/TypeScript, PostgreSQL Docker.
Demo menggunakan satu toko dan pembayaran sintetis.

Backend dan frontend memakai clean architecture. React memakai Redux Toolkit
untuk state aplikasi, async thunk untuk use case, dan typed hooks/selectors.
Struktur layer, dependency injection, dan pemeriksaan import ada di
[dokumentasi arsitektur](docs/ARCHITECTURE.md).

## Hasil yang sudah berjalan

- Loop model dan Chrome nyata, screenshot setiap langkah, hasil tool dan alasan
  keputusan yang dapat dibuka dari dashboard.
- Total salah ditemukan melalui perhitungan integer.
- Timeout HTTP 504 sesudah order tersimpan dipulihkan melalui pemeriksaan status
  dengan ID korelasi; tidak ada submit kedua atau order ganda pada evaluasi.
- Retest memakai tujuan/basket yang sama pada fixture toko yang sudah benar,
  dengan perbandingan baseline dan hasil baru.
- PostgreSQL menyimpan run, langkah, state, findings, usage, order dan referensi
  screenshot. State run ini terpisah dari checkpoint Entire.
- Lima trace evaluasi aplikasi berhasil dibaca kembali dari neatlogs, lengkap
  dengan hierarchy, input/output dan token yang cocok dengan data lokal.
- Entire mendeteksi sesi coding aktif dan file implementasi yang berubah.

Lihat [evaluasi aktual](docs/EVALUATION.md), [status capture](docs/ENTIRE-CAPTURE-NEXT.md)
dan [log acara](hackathon.md). Ini MVP lokal; kesiapan produksi belum dievaluasi.

## Menjalankan aplikasi

Prasyarat: Go 1.26.3, Node, pnpm, Docker dan Google Chrome/Chromium. Pengujian
lokal memakai pnpm 11.1.3 dan PostgreSQL 18. Chrome dapat ditentukan melalui
`CHROME_PATH` jika tidak ditemukan otomatis.

Dari root repository:

```sh
docker compose up -d db
sh scripts/dev-backend.sh
```

Terminal kedua:

```sh
cd checkoutguard_FE
pnpm install --frozen-lockfile
pnpm dev
```

Buka http://127.0.0.1:5173. Backend berada di http://127.0.0.1:8080 dan database
khusus proyek di `127.0.0.1:55433`. Browser demo mendapatkan URL sandbox unik
untuk setiap run; URL pada form adalah `http://127.0.0.1:8080/sandbox`.

Pilih **Incorrect total**, jalankan tujuan default, lalu pilih
**Retest corrected store** untuk membandingkan hasil. Pilih **Payment timeout**
untuk melihat pemulihan tanpa duplikasi. Tujuan seperti membeli dua tote dan
basket **Alternate prices** juga telah diuji dengan model live.

Retest mengaktifkan fixture dengan rumus benar; MVP belum mengedit source code
toko atau memperbaiki toko produksi secara otomatis.

Dashboard memakai cobalt blue/navy dengan teks berkontras tinggi, ringkasan
total/order/durasi, serta ilustrasi perisai checkout 3D transparan. Asset dan
prompt imagegen berada di `checkoutguard_FE/public/images/`. Badge tetap memakai
label teks untuk membedakan passed, findings dan failed. Ilustrasi dekoratif
terpisah dari screenshot evidence yang diambil browser agent.

Logo memakai mark SVG kantong belanja/perisai dengan centang, yang juga dipakai
favicon. PNG 16/32/48px, ICO dan Apple touch icon 180px dibuat dari satu sumber
`checkoutguard_FE/public/brand/checkoutguard-mark.svg`. Setelah mengubah SVG,
regenerasi ikon dari `checkoutguard_BE` dengan `go run ./cmd/brandassets` (Chrome
diperlukan), lalu build frontend.

Dashboard menyediakan **English** (default) dan **Bahasa Indonesia** melalui
pilihan Language/Bahasa di header. Preferensi disimpan di localStorage browser
dan dipulihkan saat reload. Label, status, form, pesan aplikasi yang dikenali,
serta format angka/waktu mengikuti pilihan bahasa. Contoh tujuan yang belum
diedit ikut diterjemahkan; tujuan yang ditulis pengguna, record run, keluaran
model/tool asli, serta isi screenshot tetap dipertahankan. Jika browser menolak
storage, pilihan bahasa tetap bekerja selama halaman terbuka.

## Konfigurasi backend

`scripts/dev-backend.sh` membaca konfigurasi privat dari file di direktori induk:

- `.neatlogs-api.env`: `NEATLOGS_API_KEY`, `NEATLOGS_ENDPOINT`.
- `.apertus-api.env`: `APERTUS_API_BASE_URL`, `APERTUS_API_KEY`, `APERTUS_MODEL`.

Lokasi dapat diubah dengan `NEATLOGS_ENV_FILE` dan `MODEL_ENV_FILE`. Loader membaca
assignment tanpa mengeksekusi shell. Environment proses mengambil prioritas.
Credential hanya dimuat backend dan tidak disalin ke frontend, trace, atau repo.

Untuk endpoint OpenAI-compatible lain, file model dapat memakai `MODEL_BASE_URL`,
`MODEL_API_KEY`, `MODEL_NAME`, dan `MODEL_PROVIDER`. Endpoint harus mendukung
`/chat/completions` dan keluaran JSON. Lihat [.env.example](.env.example) untuk
variabel lain; file contoh ini **tidak dimuat otomatis**.

Request evaluasi memakai `swiss-ai/apertus-v1.5-70b`, tetapi provider Public AI
mengembalikan label `aisingapore/Qwen-SEA-LION-v4-32B-IT`. Run baru mencatat model
yang diminta dan yang dikembalikan secara terpisah. Evaluasi ini tidak membuktikan
bahwa inference dilakukan oleh Apertus. Label provider juga bukan verifikasi
independen terhadap bobot model.

Token berasal respons provider. Biaya tetap `null`/Unconfigured sampai tarif
input dan output diketahui. `MODEL_INPUT_USD_PER_MILLION` serta
`MODEL_OUTPUT_USD_PER_MILLION` menghitung estimasi hanya jika model respons sama
dengan `MODEL_PRICED_MODEL` (default: model request).

Konfigurasi harga terpisah `.env.pricing.local` kini didukung. Tarif hanya
diterapkan jika provider dan model respons cocok; total parsial dilabeli minimal/
belum lengkap. Retest biaya terbaru sudah merekam estimasi USD 0,00318975,
bukan invoice. Histori lama tetap utuh. Contoh tarif dan setup portabel ada di
[README](README.md).

Request Public AI menyertakan `User-Agent` aplikasi. Deadline request mencakup
koneksi dan pembacaan response body: default 45 detik, dapat diatur melalui
`MODEL_REQUEST_TIMEOUT_SECONDS` (10–60). Deadline run tetap 180 detik.
Error autentikasi/request permanen langsung menghentikan run. Error sementara
menunggu minimal satu detik atau `Retry-After`, dengan batas dua kegagalan
berturut-turut atau tiga kegagalan total. Instruksi menunggu lebih dari 15 detik
menghentikan run agar tidak mengabaikan batas provider. Retry hanya meminta
keputusan model baru; tidak mengulang aksi browser atau submit order.
Error dan evidence disimpan, serta span LLM ditandai error. Usage request yang
timeout tidak diketahui, meskipun token respons yang diterima tetap dicatat.

## Batas dan penyimpanan

Dashboard memiliki tiga lingkungan: **Sandbox**, **Shopee**, dan **Tokopedia**.
Sandbox tetap menjadi demo sintetis yang teruji. Marketplace merupakan adapter
eksperimental untuk **menyiapkan checkout**, dengan konfirmasi order terakhir
dilakukan manual; belum ada evaluasi checkout marketplace menggunakan akun login.

Untuk marketplace, pilih lingkungan, tekan Open browser, lalu login langsung
di Chrome khusus CheckoutGuard. Masukkan URL produk HTTPS penuh (bukan shortlink):
Shopee `https://shopee.co.id/product/SHOP_ID/ITEM_ID` atau slug `-i.SHOP_ID.ITEM_ID`,
Tokopedia `https://www.tokopedia.com/SHOP/PRODUCT`. Parameter tracking dibuang
sebelum URL disimpan/dikirim ke model. Profil browser privat disimpan di
`checkoutguard_BE/tmp/marketplace-profiles/`, diabaikan dari Git. Jangan membagikan
profil tersebut: session login tersimpan di dalamnya. Agent tidak mengisi login,
OTP, kartu, atau alamat, dan tidak menyelesaikan access challenge.

Agent memilih aksi dari observasi nyata dengan budget yang sama. Kontrolnya
dibatasi pada varian, jumlah 1–5, tambah keranjang/beli sekarang, dan checkout
dari keranjang. Tombol order/payment terakhir tidak diekspos sebagai aksi.
Request mutasi dibatasi ke persiapan keranjang/preview Shopee, serta query baca
dan bentuk mutasi keranjang GraphQL yang dikenali untuk Tokopedia. Perubahan
API/DOM yang tidak dikenali dapat menghentikan alur; jangan menganggap adapter
ini sudah kompatibel dengan semua produk atau versi marketplace.

Status `ready_for_review` berarti halaman review checkout terlihat, **bukan
order berhasil dibuat**, bukan total terverifikasi, dan bukan pembuktian barang/
varian/anggaran sudah benar. Pengguna meninjau rincian di browser. `needs_login`
dan `blocked` menyimpan alasan dan evidence; data pesanan sandbox tidak digunakan
untuk membuktikan pesanan marketplace. Pada kegagalan navigasi sebelum observasi,
evidence halaman mungkin tidak tersedia. UI tidak menghitung review marketplace
sebagai checkout passed. Retest corrected-fixture hanya tersedia pada sandbox.

Tombol Use Sandbox demo hanya memilih form demo. Run marketplace tetap tersimpan
dengan hasil aslinya; fallback tidak berjalan otomatis dan tidak mengganti hasil.
Satu sesi browser marketplace aktif; pindah marketplace atau memulai sandbox
menutup browser sebelumnya, sementara profil login tetap tersimpan.

Endpoint tambahan: `GET /api/marketplaces/session`,
`POST /api/marketplaces/{shopee|tokopedia}/connect`, serta
`POST /api/marketplaces/session/close`. Endpoint POST menerima JSON `{}`.
`browser_open` menunjukkan browser aktif, bukan authentication terverifikasi.

Jika Shopee/Tokopedia menampilkan CAPTCHA, verifikasi akses, atau “Try again
later”, CheckoutGuard menandainya `blocked` berdasarkan URL dan teks halaman.
Deteksi juga mencakup halaman gagal verifikasi tanpa kata CAPTCHA di body.
Klik **Check browser / Periksa browser** untuk membaca ulang status halaman
browser yang sudah terbuka; tool ini tidak menavigasi ulang atau menyelesaikan
CAPTCHA. State terakhir tersedia sebagai `page_state` pada API koneksi/sesi.
Status ini diperbarui saat connect/check dan observasi run, bukan pemeriksaan
login terus-menerus. `browsing` tidak membuktikan akun sudah login.

`page_check: unavailable` berarti halaman tidak dapat dibaca pada pemeriksaan
terakhir. State akses sebelumnya dipertahankan dan checkout tetap ditahan,
termasuk bila state sebelumnya browsing. Tunggu browser merespons lalu tekan
Periksa browser lagi; kegagalan baca bukan bukti verifikasi sudah selesai.

Saat `blocked` atau `needs_login`, UI menahan Prepare checkout dan backend
menolak start sampai akses dinilai ulang. Ikuti instruksi marketplace secara
manual: tunggu bila diminta, lalu coba di browser yang sama dan tekan Periksa
browser. Pilihan Use Sandbox demo juga tersedia sebelum run dimulai. Pemilihan
Sandbox tidak otomatis menjalankan agent. Tidak ada perubahan untuk menyamarkan
automation atau melewati verifikasi; keberhasilan akses marketplace tetap belum
dijamin. Bila browser biasa juga gagal, masalah tidak terbatas pada sesi agent;
gunakan bantuan resmi marketplace jika kegagalan berlanjut.

Pengujian tambahan adapter memakai fixture browser sintetis (bukan run sukses
Shopee/Tokopedia):

```sh
cd checkoutguard_BE
CHECKOUTGUARD_BROWSER_TEST=1 CHECKOUTGUARD_TEST_DATABASE_URL='postgres://checkoutguard:checkoutguard_local@127.0.0.1:55433/checkoutguard?sslmode=disable' go test -race ./...
```

Satu run aktif dan satu browser agent, maksimum 18 aksi, 180 detik, serta satu
retry submit sandbox setelah pemeriksaan status. Request browser sandbox dibatasi ke origin
dan sandbox run sendiri. Hanya kontrol halaman yang terlihat dapat diklik/diisi;
email dibatasi `demo@example.test`, jumlah barang 1–5, tanpa JavaScript bebas dari
model. Instruksi dalam halaman diperlakukan sebagai data.

Order memakai constraint unik `(run_id, correlation_id)` dan upsert idempotent.
Restart backend menandai run yang belum selesai sebagai stopped; penghentian
normal menyimpan status akhir dan melakukan flush telemetry. PostgreSQL memakai
volume Docker terpisah; screenshot berada di `checkoutguard_BE/evidence/`.
Binary, evidence, laporan `tmp/`, credential, dependencies dan build frontend
diabaikan dari Git. Database dan file screenshot sama-sama diperlukan untuk
membuka ulang evidence lokal.

API utama: `GET /api/health`, `GET/POST /api/runs`, `GET /api/runs/{id}`,
`GET /api/runs/{id}/orders`, `POST /api/runs/{id}/retest`, dan
`POST /api/runs/{id}/stop`. Server bind ke loopback.

## Verifikasi

Diagnostik provider tanpa aksi browser/order, dari direktori backend:

```sh
cd checkoutguard_BE
go run ./cmd/modelprobe -catalog
go run ./cmd/modelprobe -attempts 2
```

Probe mencatat waktu DNS/koneksi/TLS/respons, model yang diminta dan label yang
dikembalikan, serta usage aktual. Credential tidak dicetak. `-model MODEL_ID`
hanya mengganti model untuk probe, bukan konfigurasi server; `-http1` menguji
protokol secara terpisah. Probe keluar dengan error jika semua percobaan gagal.
Katalog yang tersedia atau satu probe sukses bukan bukti checkout berhasil.

`-messages-file path.json` memuat array pesan privat untuk menguji input model
tertentu. Probe tetap tidak mengeksekusi action yang dikembalikan. Diagnosis
provider, timing aktual dan batas kesimpulannya ada di
[MODEL-CONNECTION-DIAGNOSIS.md](docs/MODEL-CONNECTION-DIAGNOSIS.md).

`-env-file path` memilih file konfigurasi provider untuk probe tanpa mengganti
server. Endpoint/key/model kandidat harus lengkap; tidak ada pencampuran dengan
credential APERTUS. Panduan akses, probe dan evaluasi sebelum aktivasi ada di
[PROVIDER-CANDIDATE.md](docs/PROVIDER-CANDIDATE.md). Provider cadangan masih
menunggu akses; belum ada failover otomatis atau klaim LiteLLM live lulus.

Backend (termasuk 12 submit bersamaan dengan satu order tersimpan):

```sh
cd checkoutguard_BE
CHECKOUTGUARD_TEST_DATABASE_URL='postgres://checkoutguard:checkoutguard_local@127.0.0.1:55433/checkoutguard?sslmode=disable' go test -race ./...
go vet ./...
```

Frontend:

```sh
cd checkoutguard_FE
pnpm test
pnpm build
```

Evaluasi end-to-end dari root, dengan backend berjalan dan tidak ada run aktif:

```sh
node scripts/smoke.mjs
```

Ini membuat empat run baru dengan panggilan API model dan telemetry nyata;
aksi browser dipilih model. Laporan ditulis ke `tmp/live-evaluation.json`.

Untuk memeriksa satu scenario saja:

```sh
node scripts/check-model-reliability.mjs timeout
```

Pilihan scenario: `fixed`, `timeout`, atau `total_bug`. Script memeriksa total,
tepat satu order persisten, satu klik submit, serta MIME/signature PNG semua
evidence. Hasil asli, termasuk kegagalan, disimpan di `tmp/model-reliability/`;
script keluar dengan error jika checkout tidak memenuhi assertion.

Retest run sandbox yang sudah memverifikasi expected total dapat dijalankan
dengan `node scripts/check-model-reliability.mjs --retest RUN_ID`. Script memakai
endpoint retest aplikasi, mempertahankan goal/basket/parent, dan memeriksa total
corrected fixture terhadap expected total baseline. Hasil baseline dan retest
disimpan terpisah; kondisi baseline tidak diubah.

Readback trace dapat dijalankan dengan ID dari laporan:

```sh
python3 scripts/read-traces.py TRACE_ID_1 TRACE_ID_2
```

Script memakai route read v3 autentikasi SDK dan menyimpan metadata struktur
aman di `tmp/trace-readback.json`. Hasil pass mencakup hierarchy, satu root,
tanpa span duplikat, input/output, token, serta kecocokan usage lokal jika ada.
Marker header diagnostik hanya untuk read; span aplikasi tidak diubah.
Gunakan `--run-file path/ke/run.json` untuk membandingkan token dengan snapshot
run tertentu atau laporan evaluasi yang dipindahkan.
`--expected-status error` memeriksa trace yang memiliki span error, termasuk run
yang pulih lalu passed. Hasil verified menyatakan telemetry lengkap sesuai
status yang diharapkan; status checkout harus diperiksa dari record run.

Dengan history baseline/retest tersedia dan agent sedang idle, QA dashboard:

```sh
cd checkoutguard_BE
go run ./cmd/visualcheck
```

Screenshot desktop, retest, dan mobile disimpan di `tmp/ui/`; tool memeriksa
runtime exception dan overflow horizontal halaman mobile.
Capture baru menghasilkan PNG asli. Screenshot evaluasi awal memiliki byte JPEG
di bawah nama `.png`; backend mendeteksi dan menyajikan MIME asli untuk evidence
lama tersebut tanpa mengubah rekamannya.

## Partner dan submission

neatlogs terintegrasi dengan workflow, span LLM/tool dan hosted readback.
Doctor local/hosted lulus. Rilis SDK `v0.1.12` masih mendefinisikan runtime
`Version = "0.1.9"`; ketidaksesuaian metadata upstream ini dicatat, tidak diubah.

Capture Entire aktif. Commit anchor kosong `34b821f` serta remote/upstream main
dibuat hanya atas izin khusus user. Commit/push fitur berikutnya dikelola user;
checkpoint fitur yang dipublikasikan dan link reviewer masih perlu diverifikasi.

cfo.ai masih **pending**: akses akun, rencana bisnis dari layanan tersebut, dan
bukti penggunaannya belum tersedia. [Input rencana bisnis](docs/CFO-PLAN-INPUTS.md)
disiapkan sebagai draft. Video, building in public, dan submission juga belum
selesai; MVP ini belum memenuhi seluruh keluaran submission neatHack.
