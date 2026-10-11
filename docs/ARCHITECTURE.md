# Arsitektur CheckoutGuard

Backend dan frontend memisahkan aturan domain, use case, serta adapter. Dependency
mengarah ke layer dalam. Startup memilih implementasi konkret dan menyuntikkannya
ke use case; tidak ada akses database, Chrome, SDK tracing, HTTP, React, atau Redux
di domain. Use case frontend juga tidak mengimpor React atau Redux.

## Backend Go

```text
checkoutguard_BE/
  cmd/server/main.go             composition: config, adapter, service, HTTP
  internal/domain/              run, quote, finding, action, kebijakan URL
  internal/application/         port, lifecycle run, loop agent, sandbox
  internal/adapters/
    postgres/                   repository dan order idempotent
    model/                      provider OpenAI-compatible dan span LLM
    browser/                    Chrome, sesi marketplace, DOM/network policy
    telemetry/                  implementasi tracer neatlogs dan lifecycle SDK
  internal/delivery/httpapi/     routing, JSON, HTTP status, evidence, HTML demo
  internal/config/              loader environment untuk startup/adapter
  internal/architecture/        tes arah dependency
```

`application.Service` bergantung pada `Repository`, `Planner`, `BrowserFactory`,
`MarketplaceGateway`, dan `Tracer`. ID generator disuntikkan dari startup. Tipe
kontrak bersama berasal dari domain; alias tipe menjaga satu definisi dan JSON
API yang sama. Domain hanya mengimpor standard library. Application mengimpor
standard library dan domain, bukan pgx, chromedp, neatlogs, atau net/http.

HTTP memanggil service untuk start, stop, retest, pembacaan run/order, admission
sesi browser, quote dan order sintetis. Timeout sintetis ditentukan use case
sesudah persistensi order; delivery menerjemahkannya menjadi HTTP 504. Database
sandbox tidak dapat dipakai untuk membuat atau memverifikasi order marketplace.

Application tetap menjalankan keputusan model secara dinamis. Batas satu run,
18 aksi, 180 detik, pemeriksaan status sebelum retry, dan konfirmasi manual
marketplace tetap berlaku. Browser adapter menjaga scope, selector, screenshot,
dan larangan final submit. Tracer port mempertahankan workflow/span asli.

Adapter model menerjemahkan kegagalan transport/HTTP/validasi menjadi
`domain.PlannerError`: jenis error, pesan aman, sifat retry, dan waktu tunggu.
Application menentukan budget recovery, mempersist error/evidence, lalu meminta
keputusan berikutnya setelah jeda yang dapat dibatalkan. History error tidak
memuat aksi assistant palsu. Retry model tidak mengeksekusi ulang aksi browser;
verifikasi status dan batas retry order tetap ditangani use case sandbox.

## Frontend React + Redux Toolkit

```text
checkoutguard_FE/src/
  main.tsx                      composition: repository → use case → store
  domain/                       entitas dan aturan lingkungan checkout
  application/                  repository port dan use case
  infrastructure/               implementasi repository HTTP
  presentation/
    state/                      Redux slice, async thunk, selector, typed hook
    dashboard/                  page, composer, history, browser evidence
    components/                 komponen UI bersama
    style.css                   styling dashboard
```

State form, daftar/seleksi run, seleksi langkah, health, sesi marketplace,
loading/error, dan dialog evidence berada di Redux. Komponen mengirim action
atau thunk melalui typed hooks. Thunk menerima use case lewat `extraArgument`;
thunk dan komponen tidak memilih adapter HTTP sendiri.

Polling dapat dibatalkan dan tidak bertumpuk. Revision state menolak hasil poll
yang dimulai sebelum mutasi, sehingga snapshot lama tidak menghapus run atau
sesi baru. Launch ganda dicegah; memilih sandbox hanya mengubah form dan tetap
menyimpan hasil marketplace sebelumnya. Default transport memakai native fetch
dengan konteks browser yang benar. Store menyimpan data serializable, bukan
Promise, AbortController, sesi Chrome, atau objek repository.

Bahasa antarmuka (`en`/`id`, default `en`) berada di Redux. Domain mendefinisikan
bahasa yang valid; `LanguagePreferences` mengakses repository port dan adapter
browser menangani localStorage yang bisa ditolak. Composition root menyuntikkan
use case preferensi ke store; presentation tidak mengakses localStorage langsung.
Kamus typed pada presentation menerjemahkan label dan pesan aplikasi saat render.
Record run dan evidence asli tidak dimutasi. Pergantian bahasa hanya menerjemahkan
contoh tujuan yang belum diedit, serta memperbarui lang/title dokumen dan format
angka/waktu. Tujuan buatan pengguna tetap utuh.

## Menjaga batas dan menguji perubahan

```sh
# Backend: arah import, domain, use case dan adapter
cd checkoutguard_BE
go test -race ./...
go vet ./...

# Frontend: arah import, TypeScript, Redux/transport regression, build
cd checkoutguard_FE
pnpm check:architecture
pnpm test
pnpm build
```

Tes Go di `internal/architecture` melarang dependency luar pada domain dan
application serta import adapter konkret dari delivery. Pemeriksaan frontend
menggunakan AST TypeScript, termasuk import dinamis, dan berjalan saat build.

Tes dengan fake port memverifikasi orchestration; hasilnya bukan bukti live
model atau checkout marketplace. Tes PostgreSQL/Chrome opt-in dan evaluasi
provider nyata dijelaskan di [README](../README.md) dan
[evaluasi](EVALUATION.md). Screenshot QA dashboard tetap berada di `tmp/ui/`.

Tambahkan aturan bisnis ke domain/use case, implementasi layanan eksternal ke
adapter, dan presentasi ke UI. Fitur metode pembayaran marketplace belum
ditambahkan oleh refactor ini; agent masih berhenti sebelum submit order.
