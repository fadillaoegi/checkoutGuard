# Diagnosis timeout model CheckoutGuard

Pemeriksaan awal 11 Oktober 2026, setelah user melaporkan tiga timeout pada demo
sandbox. Probe awal hanya melakukan inference. Perbaikan observabilitas dan
retest sandbox berikutnya dicatat terpisah di bawah; provider, key, dan default
model tetap sama. Tidak ada operasi Git.

## Provider dan konfigurasi aktual

| Bagian | Nilai yang diperiksa |
|---|---|
| Provider/gateway | Public AI (`publicai`) |
| API host | `api.publicai.co`, base `/v1` |
| Model yang diminta | `swiss-ai/apertus-v1.5-70b` |
| Model pada respons sukses | `aisingapore/Qwen-SEA-LION-v4-32B-IT` |
| Deadline request | 45 detik, termasuk pembacaan response body |
| API key | Ada; nilainya tidak dicetak |
| Proxy environment pada proses probe | HTTP_PROXY/HTTPS_PROXY/ALL_PROXY dan versi lowercase tidak dikonfigurasi |

Public AI menyediakan gateway OpenAI-compatible yang merutekan request ke
inference partner serta menangani load balancing/failover menurut
[FAQ resmi](https://platform.publicai.co/inference-partners-faq). Respons model
yang berbeda dari request terbukti pada run dan probe; belum diketahui apakah
ini failover lintas model, alias, atau kesalahan metadata/routing. Label respons
bukan verifikasi independen terhadap bobot model. CheckoutGuard tidak mengubah
nama model secara diam-diam.

## Run yang dilaporkan user

Run `28e6abfb-2197-4bb5-9ae9-e3e41e5054c0`: sandbox, total salah, basket harga
alternatif. Hasil akhir failed karena budget recovery habis.

| Keputusan | Hasil | Latensi termasuk evidence |
|---|---|---|
| 1 | Model timeout, tidak ada aksi browser | 45,083 detik |
| 2 | Model memilih fill quantity 1 | 15,261 detik |
| 3 | Model timeout, tidak ada aksi browser | 45,084 detik |
| 4 | Model timeout, tidak ada aksi browser | 45,222 detik |

Run tidak membuat order. Usage respons yang diterima 866 input + 37 output;
usage tiga request timeout tidak diketahui dan biaya tetap null. Snapshot privat
ada di `tmp/model-reliability/28e6abfb-2197-4bb5-9ae9-e3e41e5054c0.json`.
Trace `6bb72be9ea60a2c2a5854f3c89fd377f` diverifikasi dengan expected error:
enam span (satu root, empat LLM, satu tool), hierarchy valid, tanpa duplikasi,
input/output tersedia, dan 903 token cocok dengan snapshot lokal. Readback
`tmp/model-reliability/28e6abfb-trace.json`; verification bukan checkout passed.

## Probe nyata tanpa browser atau order

Probe memakai adapter model, timeout, header aplikasi dan schema action yang
sama. Hasil action hanya dicetak; tidak pernah dieksekusi. Pesan checkout awal
direkonstruksi dari system prompt aplikasi, goal asli, dan initial observation
yang tersimpan, lalu disimpan privat untuk reproduksi.

| Input | DNS/koneksi/TLS selesai | Byte respons pertama | Request selesai | Label respons |
|---|---|---|---|---|
| JSON action kecil, request 1 | 0,543 detik | 2,240 detik | 2,240 detik | SEA-LION |
| JSON action kecil, request 2 | Koneksi digunakan ulang | 0,640 detik | 0,640 detik | SEA-LION |
| Input awal checkout, request 1 | 0,098 detik | 2,372 detik | 2,372 detik | SEA-LION |
| Input awal checkout, request 2 | Koneksi digunakan ulang | 0,519 detik | 0,519 detik | SEA-LION |

Semua action valid. Probe kecil melaporkan 45 input + 11 output token per
respons; input checkout 595 + 35 per respons. Katalog autentikasi juga berhasil
melalui HTTP/2 dan mencantumkan model Apertus yang diminta. Tidak ada HTTP
401/403/429 pada uji ini. Ini tidak membuktikan saldo billing atau keandalan
provider pada semua request; respons bisa dipengaruhi load atau caching.

Pemeriksaan sebelumnya pada provider yang sama mengukur DNS/TCP/TLS selesai
kurang dari 200 ms, request terkirim, tetapi tidak ada byte respons sampai
deadline 45 detik. Hasilnya dicatat di [EVALUATION.md](EVALUATION.md).

## Kesimpulan dan batas diagnosis

Tidak ditemukan bukti bahwa koneksi internet user putus, credential ditolak,
atau seluruh request provider gagal. Request kecil dan input checkout bisa
berhasil; run live dan probe sebelumnya mengalami timeout. Tidak ditemukan
error schema sebagai penyebab pada sampel tersebut. Bukti mengarah ke waktu
tunggu respons yang berubah-ubah pada jalur API/gateway/inference. Antrean,
routing atau layanan upstream yang lambat adalah dugaan, bukan akar masalah
server yang sudah terkonfirmasi; akses metrik/log provider diperlukan untuk
menentukannya. Gangguan jalur internet khusus endpoint juga belum sepenuhnya
dapat dikesampingkan.

Timeout 45 detik merupakan batas yang dikonfigurasi aplikasi. Memperpanjangnya
akan memberi provider waktu lebih lama, tetapi tidak menjamin stabilitas dan
mengurangi jumlah keputusan yang bisa selesai dalam budget run 180 detik.
Provider lain perlu diuji dengan goal/tool loop sebenarnya sebelum dipilih.
Mengganti model di gateway yang sama belum tentu mengatasi keterlambatan;
dua probe SEA-LION sebelumnya juga pernah timeout.

Reproduksi inference saja, dari direktori backend:

```sh
go run ./cmd/modelprobe -catalog
go run ./cmd/modelprobe -attempts 2
go run ./cmd/modelprobe -messages-file ../tmp/model-reliability/28e6abfb-initial-messages.json -attempts 2
```

File pesan diagnosis berada di tmp dan tidak dipublikasikan. `-messages-file`
hanya dipakai probe, tidak mengubah server atau mengeksekusi returned action.
`go vet ./cmd/modelprobe` lulus; probe baru berhasil memakai file pesan tersebut.

## Pemisahan label permintaan dan respons

Implementasi sebelumnya mengisi label span dengan model yang diminta, lalu
menggantinya dengan label model dari respons sukses. Timeout tanpa respons
mempertahankan label Apertus, sementara respons sukses ditampilkan sebagai
SEA-LION. Pola ini merupakan artefak instrumentasi; ia tidak membuktikan
A/B routing atau kegagalan eksklusif model tertentu.

Adapter sekarang mempertahankan `neatlogs.llm.model_name` sebagai model yang
diminta sepanjang request. Label dari respons disimpan terpisah dalam custom
metadata SDK `neatlogs.metadata`, objek `checkoutguard`, dan pada setiap
`step.usage.model_call` yang tersimpan di aplikasi:

| Field | Makna |
|---|---|
| `requested_model` | Model pada body request |
| `returned_model` | Label dari respons JSON, hanya jika tersedia |
| `returned_model_known` | Apakah label respons benar-benar diketahui |
| `model_label_mismatch` | Perbandingan label, hanya jika keduanya diketahui |
| `timeout_ms` | Batas waktu HTTP client; bukan lama pemrosesan provider |
| `error_kind`, `retryable` | Kategori kegagalan dan kelayakan retry planning |
| `http_response_received`, `http_status_code` | Header/status HTTP telah diterima |
| `response_decoded` | JSON respons sudah berhasil dibaca |

Ini membedakan timeout sebelum header, timeout saat membaca body, HTTP error,
dan respons yang sudah diterima tetapi action tidak valid. Usage dan tarif
masih memakai label respons sesuai aturan harga; request tanpa usage tetap
unknown/incomplete. Tidak ada tambahan endpoint, header otorisasi atau body
HTTP error dalam metadata.

SDK Go memetakan custom metadata ke `neatlogs.metadata`; SDK Python resmi juga
menggunakan atribut metadata JSON untuk custom metadata. Lihat
[implementasi SDK](https://github.com/NeatLogs/neatlogs/blob/main/neatlogs/decorators/_base.py).
Namun API hosted v3 tidak menampilkan custom metadata LLM pada uji aktual.
Workflow output karena itu juga memuat `model_calls`, berupa JSON dengan marker
`checkoutguard.model_calls.v1=`. Setiap record memiliki nomor langkah dan ID
span model. Ini observasi client, bukan parameter tambahan yang dikirim ke
provider atau bagian dari jawaban model.

Adapter telemetry mengisi input/output sesuai jenis span (`neatlogs.workflow.*`
dan `neatlogs.tool.*`), serta declared trace output hanya untuk workflow.
API v3 mengubah output workflow menjadi teks dengan nama field yang diformat;
marker JSON mempertahankan struktur diagnostik melalui pemformatan tersebut.
`scripts/read-traces.py` membaca field itu, mengaitkannya ke span LLM, dan
memeriksa label request, timeout serta keberadaan label respons. Pembacaan
hanya mengekspor field diagnostik yang diizinkan, tanpa prompt mentah.

Label trace lama tidak ditulis ulang. Analisis 14 hari harus memisahkan
instrumentasi lama dari trace baru dan tidak memperlakukan label respons
sebagai bukti independen bobot model yang melayani request.

## Validasi perbaikan observabilitas

Retest akhir `9689d682-97ac-4982-bfd6-538aff957d4e` passed dalam 35,15 detik:
total 9.800 IDR, satu order, satu submit, delapan PNG valid, tanpa model timeout.
Usage 12.306 input + 197 output; estimasi USD 0,003175 dari tarif respons yang
cocok. Trace `869bdb1c0d98cd794c3a895292fe23c7` dibaca kembali dari neatlogs:
15 span, satu root, tujuh LLM, tanpa duplikasi, token cocok dengan run lokal.
Ketujuh request mempertahankan label Apertus; diagnostik terpisah menunjukkan
respons SEA-LION, HTTP 200, `error_kind=none`, `timeout_ms=45000`, dan label
mismatch diketahui. Readback memverifikasi semua field melalui ID span.
[Ringkasan aman](../submission/evidence/model-diagnostics-trace.json).

Retest sebelumnya pada sesi perbaikan, `a1be070b-6bf8-4232-8e60-f98d0d9954e0`,
passed dengan satu model timeout 45,17 detik, lalu pulih menjadi satu pesanan
9800 IDR dengan satu submit. Itu membuktikan bounded planning recovery dapat
melanjutkan run, tetapi tidak membuktikan timeout provider telah hilang.
Metadata hosted tambahan belum terbaca pada versi antara tersebut.

Go race tests termasuk Chrome/PostgreSQL, go vet, delapan kasus SDK tracing
(matching/mismatched/missing label, timeout sebelum/sesudah header, HTTP 503,
action invalid, cancellation), scope output workflow/tool, serta tiga tes
reader format hosted lulus. Histori dan bukti lama tetap dipertahankan.

Uji transport timeout hosted juga terverifikasi pada trace terpisah
`7e961856126283f3804024c6caa95069`, bernama
`checkoutguard.synthetic-model-diagnostics`. HTTP fixture lokal memberi satu
respons sukses dengan label berbeda, lalu satu timeout terkontrol 50 ms.
Neatlogs mempertahankan label request untuk keduanya; request timeout mencatat
`returned_model_known=false`, `http_response_received=false`,
`response_decoded=false`, `error_kind=timeout`, tanpa mengarang model respons.
Tiga span dan kedua record diagnostik verified. Ini fixture transport tanpa
inference provider, browser, atau order; bukan bukti performa provider dan
tidak dipakai sebagai bukti checkout kompetisi.
