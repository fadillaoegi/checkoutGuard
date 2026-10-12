# Provider cadangan untuk demo CheckoutGuard

Status: konfigurasi dan alat uji siap, tetapi provider alternatif **belum diuji
atau diaktifkan**. Saat pemeriksaan, key LiteLLM/provider cadangan belum ditemukan
di konfigurasi proyek dan environment. Public AI masih menjadi provider aktif.

Pada 12 Oktober, profil demo lokal beralih ke **model** Apertus v1.5 8B pada
Public AI yang sama; ini bukan aktivasi provider cadangan. Key yang sudah ada
digunakan dalam `.env.model.local` berizin 600 dan diabaikan Git. Konfigurasi
70B di direktori induk tetap tersedia. Tiga probe JSON-object dan satu checkout
fixed lengkap berhasil tanpa model error; request/response label 8B cocok.
Pengujian berikutnya tetap menemukan ringkasan action terlalu panjang,
pemilihan tool sebelum tersedia, serta HTTP 502 dari gateway. Jadi model 8B
**belum cukup andal untuk menjamin demo checkout live**. State verifikasi dan
daftar tool yang tersedia kini ikut diberikan pada model; guard tetap menolak
aksi yang tidak sah. Perubahan ini tidak menjamin gateway menjadi stabil.
Uji schema JSON tambahan tidak dipilih untuk demo karena menghasilkan keputusan
finish terlalu dini yang ditolak guard. Adapter akhir tetap JSON-object dan
tidak melakukan fallback tersembunyi. Ringkasan keputusan kini diminta maksimum
100 karakter dengan validasi action yang tetap berlaku. Hasil video dan retest
terbaru dicatat di [skrip demo](../submission/DEMO-SCRIPT.md).

## Akses yang diperlukan

Acara menyediakan akses LiteLLM untuk Nova/Claude, dengan budget bersama $20 per
tim. Ketua tim mengklaim lewat email check-in pada
[halaman model neatHack](https://neatlogs.com/hackathon/models), sesuai
[informasi resmi](https://neatlogs.com/hackathon). Ketersediaan alokasi tim dan
model yang diperbolehkan harus diperiksa saat klaim. Budget tersebut bukan bukti
akses tim sudah aktif. Akun dan verifikasi email ditangani user.

Diperlukan tiga nilai dari penyedia: base URL API OpenAI-compatible, API key,
dan served model ID. Jangan menebak URL atau model ID dari nama produk, dan
jangan kirim key di chat. Adapter saat ini menggunakan `/chat/completions`,
messages, max_tokens, temperature dan response_format JSON. Kompatibilitas
parameter Nova/Claude melalui gateway harus dibuktikan pada probe.

## Konfigurasi terpisah

Salin template [model-candidate.env.example](examples/model-candidate.env.example)
ke `.env.model-candidate.local` pada root CheckoutGuard. Isi nilainya secara
lokal dan batasi izin file menjadi 600. File `.env.*` lokal diabaikan oleh aturan
repository. Tidak perlu mengedit konfigurasi Public AI di direktori induk.

`MODEL_BASE_URL`, `MODEL_API_KEY`, dan `MODEL_NAME` harus berupa satu konfigurasi
lengkap. Ketika salah satu variabel MODEL tersebut ada, adapter tidak mengisi
field kosong dengan variabel APERTUS. Konfigurasi kandidat yang belum lengkap
berhenti sebelum request, sehingga key Public AI tidak terkirim ke endpoint
lain. Beri `MODEL_PROVIDER=litellm` atau nama provider yang benar untuk telemetry.
Tarif tetap kosong sampai model respons dan harga diverifikasi.

Environment proses mengambil prioritas atas file. Jika terminal sudah memiliki
override MODEL yang berbeda, gunakan terminal baru atau bersihkan override
tersebut sebelum pengujian; mengganti env-file tidak mengubah environment induk.

## Uji sebelum aktivasi

Dari direktori `checkoutguard_BE`:

```sh
go run ./cmd/modelprobe -env-file ../.env.model-candidate.local -catalog
go run ./cmd/modelprobe -env-file ../.env.model-candidate.local -attempts 2
```

Katalog belum membuktikan inference bekerja. Kedua probe harus mengembalikan
action JSON yang valid; periksa provider, requested_model, model pada usage,
latensi dan token. Ketidaksesuaian label model perlu dicatat, bukan disamarkan.
Jika endpoint tidak mendukung katalog, kegagalan katalog perlu dinilai terpisah
dari dukungan chat completion. Request probe memakai key/provider tersebut dan
dapat mengonsumsi budget; hasil action tidak dieksekusi.

Input awal checkout dari diagnosis sebelumnya dapat diuji tanpa browser:

```sh
go run ./cmd/modelprobe -env-file ../.env.model-candidate.local -messages-file ../tmp/model-reliability/28e6abfb-initial-messages.json -attempts 2
```

Jangan menganggap dua respons cepat sebagai bukti provider stabil. Sesudah
probe lulus, backend dapat dijalankan dengan konfigurasi kandidat:

```sh
MODEL_ENV_FILE="$PWD/.env.model-candidate.local" sh scripts/dev-backend.sh
```

Perintah terakhir dijalankan dari root. Backend lama harus dihentikan saat
idle sebelum menjalankan penggantinya; menghentikan backend juga menutup browser
marketplace khusus. Profil browser tersimpan. Tidak ada perpindahan provider
otomatis atau replay browser pada workflow ini.

## Bukti pemilihan provider demo

Dengan backend kandidat aktif dan tidak ada run lain, jalankan satu `fixed`
sandbox untuk membuktikan loop agent, bukan hanya probe:

```sh
node scripts/check-model-reliability.mjs fixed
```

Periksa total standard 12.400 IDR, tepat satu order, satu submit, seluruh PNG,
usage dan hosted trace. Jika lulus, lanjutkan scenario timeout dan total_bug,
lalu retest yang menampilkan baseline/corrected fixture. Ini membuat run dan
order sintetis, bukan pesanan marketplace. Kegagalan provider tetap dicatat
sebagai gagal. Hasil dan konfigurasi/model respons dicatat sebelum memilih
provider demo. Aktifkan kembali Public AI melalui `scripts/dev-backend.sh` tanpa
override MODEL_ENV_FILE ketika backend sudah idle dan override proses bersih.

Tes konfigurasi saat ini memakai fixture: legacy Public AI tetap valid,
profil kandidat lengkap memakai credential sendiri, dan profil tidak lengkap
tidak mengirim request. Fixture ini bukan bukti LiteLLM live berhasil.
