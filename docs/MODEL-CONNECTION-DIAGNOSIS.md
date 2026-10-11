# Diagnosis timeout model CheckoutGuard

Pemeriksaan 11 Oktober 2026, setelah user melaporkan tiga timeout pada demo
sandbox. Tidak ada perubahan provider/key/default model, operasi Git, aksi
browser checkout, atau order pada diagnosis ini.

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
