# Skrip video demo CheckoutGuard — target 2 menit 55 detik

## Video final untuk submit (v2, dengan segmen bukti) — 12 Oktober 2026

File: `submission/video/CheckoutGuard-final.mp4`, **2 menit 56,6 detik**, 1920 × 1080,
30 fps, H.264 + AAC. SHA-256 `183a24ee96d2830e0d2de4e9ad448095242e29879b750b2939492562cc688733`.
Dekode penuh lulus; −14,9 LUFS, puncak −1,5 dBFS.
**Folder `submission/video/` di-ignore Git** (±1,1 GB: video, voice-over, footage); semua file
video hanya ada di komputer lokal. Reviewer menonton lewat link video yang diunggah. Versi sebelumnya (2:43.7, tanpa segmen
bukti) diarsipkan di `footage/older-versions/CheckoutGuard-final-v1-no-proof-20261012.mp4`.

| Waktu | Bagian |
| --- | --- |
| 0:00 | Intro logo 3D |
| 0:07 | Animasi penjelasan: masalah (total 14.400 vs seharusnya 12.400 IDR), cara kerja 4 langkah |
| 0:25.4 | Demo dashboard + caption per bab + subtitle Inggris + voice-over 1–6 (mulai 0:25.8, 0:35.8, 0:55.9, 1:21.7, 1:34.7, 1:53.4). Bagian tanpa narasi dipercepat 3–4× dengan badge "⏩ sped up · no narration"; before & after di-zoom |
| 2:03.6 | **Segmen bukti** (45 detik), 4 bagian: |
| 2:03.6 | 1 · Failure recovery — screenshot asli `evidence/payment-timeout.png`, run `b74d8ccb`: timeout pembayaran, `order_status`, tepat 1 order / 1 submit |
| 2:13.8 | 2 · neatlogs — readback trace terverifikasi `000d9deb…` (12.429 token) → `7a97b78b…` (12.307 token), 15 span, hierarki valid, token cocok dengan catatan lokal |
| 2:27.8 | 3 · Entire — screenshot asli GitHub: commit `a134e4a` dengan trailer `Entire-Checkpoint: 01M4NKBNCP60MN7M5BXS143W16`, 35 file, dan `metadata.json` checkpoint |
| 2:38.4 | 4 · cfo.ai — halaman 1 PDF rencana bisnis: harga $29/$99/$299, biaya ≈ $0,019/run, break-even ≈ 96 langganan (semua asumsi) |
| 2:48.6 | Closing logo 3D, repo, partner, tim |

Catatan kejujuran:
- Panel neatlogs adalah tampilan data readback (bukan UI neatlogs) dan memakai trace run
  bukti 11 Oktober (keranjang 9.800 IDR), bukan dua run di demo (`8dd424fc…`/`7b4c0a2a…`,
  belum di-readback).
- Rekaman layar dashboard neatlogs asli dari akun user akan memperkuat bagian ini. Jika
  tersedia (`footage/partners/neatlogs.mov`), ganti panel kanan scene 2 lalu render ulang
  dengan `footage/final-edit/source/v2/build_v2.sh`.
- Halaman checkpoint GitHub di-crop tanpa banner peringatan umum "commit does not belong to any branch",
  karena checkpoint Entire disimpan di ref `refs/entire/checkpoints`, bukan di branch.

Musik orisinal (disintesis sendiri, `synth_v2.py`), diredam otomatis saat voice-over.
Subtitle: `submission/video/CheckoutGuard-final.en.srt` (25 cue, timestamp per kata) untuk YouTube.

### Voice-over opsional untuk segmen bukti (Inggris)

Segmen bukti saat ini tanpa narasi (teks di layar dan musik). Jika ingin merekam,
simpan sebagai `CheckoutGuard-voiceover-7..10.mp3` dan sesuaikan durasinya:

- **7 · Recovery (≤ 9 detik):** “When the payment response times out, the agent doesn't retry blindly.
  It checks the saved orders first, finds exactly one, and finishes without submitting again.”
- **8 · neatlogs (≤ 13 detik):** “Every run is traced in neatlogs. Each model call and browser tool is a
  span, and the token totals match our local records, before and after the fix.”
- **9 · Entire (≤ 9 detik):** “The code was built with a coding agent. Entire checkpoints that session
  and links it to the exact public commit.”
- **10 · cfo.ai (≤ 9 detik):** “With cfo.ai we priced and costed the product: plans from twenty-nine
  dollars, about two cents per run, and break-even near ninety-six subscriptions.”

## Video tanpa suara, UI bahasa Inggris — 12 Oktober 2026

File: `submission/video/CheckoutGuard-demo-sandbox.mp4` (intro 3D 7 detik: `submission/video/CheckoutGuard-intro.mp4`;
footage mentah: `submission/video/footage/`).
Durasi **2 menit 6,7 detik**, 1920 × 1080, H.264, 30 fps, **nol audio stream**,
tanpa voice-over atau subtitle. Dekode penuh lulus; frame sampel diperiksa.
SHA-256 `1a829176f92b9ea93dd52564e7c428e1db6e468b2b5a3e7241f22692ea34b519`.
[Manifest](evidence/sandbox-video-en.json).

Dashboard tetap **bahasa Inggris** sepanjang video (tidak ada pergantian
bahasa). Gaya interaksi seperti tester manual: pointer menyorot pemilih bahasa,
kartu metrik dan pilihan Sandbox; goal diketik per karakter dengan jeda tidak
seragam dan jeda lebih panjang setelah titik; scroll halus dengan easing.

Ini **review dua run sandbox tersimpan** dengan goal bahasa Inggris, bukan
checkout baru. Saat perekaman, dua percobaan live (`5cb9de57`, `55a8b1bd`)
dihentikan recorder karena Public AI timeout pada langkah 1; `modelprobe`
juga timeout 45 detik pada tiga request. Kedua run itu terlihat berstatus
**Stopped** di tabel riwayat (sekitar 1:58). Rekam ulang secara live setelah
provider pulih (perintah di bawah).

| Bagian | Run | Hasil |
| --- | --- | --- |
| Total salah, keranjang Standard | `e7af36ee-00e5-4a2b-b324-62c5025992c7` | Seharusnya 12.400 IDR, ditampilkan 14.400 IDR; 10 langkah; satu order; nol step error |
| Retest fixture benar | `8f007408-71ca-4817-9ea2-20d9e2b9a217` | Lulus, 12.400 IDR; 9 langkah; satu order; nol step error |

Trace: `8dd424fc716c6f9b45b26e659e88a85e` dan `7b4c0a2af252a387f7160450deff6760`.
Hosted readback neatlogs **belum** dilakukan untuk kedua trace ini; jangan
klaim token cocok dengan neatlogs sebelum readback.

| Waktu mulai | Bagian di layar |
| --- | --- |
| 0:00 | Dashboard Inggris; pointer ke Language dan kartu metrik |
| 0:10 | Pilihan Sandbox, goal diketik, kondisi *Incorrect total*, keranjang Standard |
| 0:38 | Menu Evidence; 10 langkah baseline dipilih satu per satu beserta screenshot |
| 1:15 | Ringkasan: expected 12,400 vs displayed 14,400 IDR; temuan dan modal bukti |
| 1:27 | Run history; memilih retest tersimpan |
| 1:35 | Langkah retest 4, 6, 7 dan hasil tool |
| 1:56 | *Before & after*; tabel *Latest runs* |
| 2:03 | Ringkasan retest lulus, satu order |

### Voice-over Inggris dan Indonesia (direkam terpisah)

Pilih satu bahasa. Jika narasi selesai sebelum pergantian layar, beri jeda.

**0:00–0:10 — Pembukaan / Introduction**

- EN: “This is CheckoutGuard, an agent that tests a checkout and keeps the
  evidence. In this video, I’m reviewing two saved sandbox checks.”
- ID: “Ini CheckoutGuard, agent yang menguji alur checkout dan menyimpan
  buktinya. Di video ini, saya meninjau dua hasil pemeriksaan sandbox yang
  sudah tersimpan.”

**0:10–0:38 — Tujuan dan kondisi / Goal and condition**

- EN: “I choose the Sandbox and describe the purchase: buy one tote, verify the
  checkout total, and confirm that exactly one order is saved. The condition is
  Incorrect total, so this store has a deliberate pricing bug. Everything here
  is simulated, with no real payment.”
- ID: “Saya memilih Sandbox lalu menuliskan tujuannya: membeli satu tote,
  memverifikasi total checkout, dan memastikan tepat satu pesanan tersimpan.
  Kondisinya Incorrect total, jadi toko ini sengaja memiliki kesalahan harga.
  Semua transaksi di sini simulasi, tanpa pembayaran nyata.”

**0:38–1:15 — Langkah agent / Agent steps**

- EN: “Every agent decision is saved as a step with a screenshot. The agent
  observes the page, sets the quantity, and adds the tote to the bag. Then it
  fills in the demo email and verifies the total. Finally, it places the order
  and reads the saved orders. Selecting a step shows the browser exactly as it
  was at that moment.”
- ID: “Setiap keputusan agent tersimpan sebagai langkah, lengkap dengan
  screenshot. Agent mengamati halaman, mengisi jumlah barang, lalu menambahkan
  tote ke keranjang. Setelah itu, agent mengisi email demo dan memverifikasi
  total. Terakhir, agent membuat pesanan dan membaca pesanan yang tersimpan.
  Memilih satu langkah menampilkan kondisi browser persis pada saat itu.”

**1:15–1:27 — Temuan harga / Pricing finding**

- EN: “Unit price minus discount plus shipping should be twelve thousand four
  hundred rupiah, but the checkout shows fourteen thousand four hundred.
  CheckoutGuard records this as a verified finding, with the screenshot as
  evidence.”
- ID: “Harga satuan dikurangi diskon, ditambah ongkir, seharusnya dua belas ribu
  empat ratus rupiah. Namun checkout menampilkan empat belas ribu empat ratus.
  CheckoutGuard mencatatnya sebagai temuan terverifikasi, dengan screenshot
  sebagai bukti.”

**1:27–1:56 — Riwayat dan retest / History and retest**

- EN: “From the run history, I open the retest. It uses the same goal and
  basket against the corrected synthetic store. This time the total matches,
  and the order check confirms exactly one saved order. The tool results and
  token usage are recorded as well.”
- ID: “Dari riwayat, saya membuka hasil retest. Retest memakai tujuan dan
  keranjang yang sama pada toko simulasi yang sudah diperbaiki. Kali ini
  totalnya sesuai, dan pemeriksaan pesanan memastikan tepat satu order
  tersimpan. Hasil tool dan penggunaan token juga tercatat.”

**1:56–2:06 — Perbandingan dan penutup / Comparison and closing**

- EN: “Before and after: fourteen thousand four hundred becomes twelve thousand
  four hundred. CheckoutGuard makes checkout verification clear and
  traceable.”
- ID: “Sebelum dan sesudah: empat belas ribu empat ratus menjadi dua belas ribu
  empat ratus. CheckoutGuard membuat verifikasi checkout jelas dan dapat
  ditelusuri.”

Jangan menyebut run sedang dieksekusi live, provider sudah stabil, atau trace
sudah dibaca ulang dari neatlogs.

Rekam ulang, dari `checkoutguard_BE`:

```sh
# Live (butuh provider model merespons; membuat dua run baru, keranjang Alternate)
go run ./cmd/demorecord
# Review hasil tersimpan seperti video ini (tanpa inference baru)
go run ./cmd/demorecord -review e7af36ee-00e5-4a2b-b324-62c5025992c7 -review-retest 8f007408-71ca-4817-9ea2-20d9e2b9a217
```

Flag `-lang` default `en` (dashboard tetap Inggris); `-lang id` mengembalikan
alur lama yang beralih ke bahasa Indonesia. Render dengan
`python3 scripts/render-sandbox-demo.py tmp/video/DIREKTORI-BARU`.

## Video tanpa suara sebelumnya (beralih ke Indonesia) — 12 Oktober 2026

File lokal: `tmp/video/20261012-023423/CheckoutGuard-sandbox-review-silent.mp4`.
Durasi **2 menit 36,53 detik**, 1920 × 1080, H.264, 30 fps. Tidak ada audio,
musik, voice-over atau subtitle. Dekode seluruh video lulus, tujuh sampel frame
dari MP4 diperiksa, dan ffprobe memastikan **nol audio stream**.

Ini **review dua hasil sandbox tersimpan**, dengan interaksi dashboard asli,
ketikan bertahap, pointer klik dan scroll halus. Tidak ada checkout atau
inference baru. Jumlah run dan order sama sebelum/sesudah rekaman. Submit form
diblokir oleh recorder pada mode review. Durasi run yang terlihat di dashboard
tetap durasi eksekusi historis, bukan waktu eksekusi model dalam video baru.
Gateway masih pernah timeout/HTTP 502 pada uji live; video ini tidak membuktikan
masalah provider telah selesai. [Diagnosis](../docs/MODEL-CONNECTION-DIAGNOSIS.md).

Yang terlihat: pergantian Inggris → Indonesia, pengisian goal, keranjang harga
alternatif, sidebar aktif, delapan langkah baseline beserta screenshot, temuan
harga, bukti dalam modal, riwayat, verifikasi retest, hasil tool, token/estimasi
biaya, dan perbandingan sebelum/sesudah. Tidak ada tes Shopee/Tokopedia, run
timeout pembayaran, stop run aktif, atau demonstrasi portal partner pada video
ini. Video submission lengkap tetap perlu segmen bukti partner.

| Waktu mulai | Bagian |
| --- | --- |
| 0:00 | Dashboard Inggris dan penggantian bahasa |
| 0:10 | Ketikan goal dan pemilihan keranjang |
| 0:45 | Review langkah baseline dan screenshot |
| 1:31 | Ringkasan total salah, temuan dan bukti modal |
| 1:46 | Sidebar riwayat dan pemilihan retest tersimpan |
| 1:57 | Review verifikasi retest serta hasil tool |
| 2:25 | Perbandingan total 11.050 → 9.800 IDR |
| 2:32 | Ringkasan retest lulus, satu order |

Gunakan pembukaan voice-over pribadi yang sesuai ruang lingkup, misalnya:
“Saya meninjau hasil pemeriksaan sandbox yang sudah tersimpan. CheckoutGuard
menyimpan keputusan agent, screenshot dan hasil verifikasi untuk setiap langkah.”
Jangan mengatakan checkout baru sedang dieksekusi atau provider sudah stabil.
Tidak ada file suara yang dibuat. [Manifest video](evidence/sandbox-video.json)
dan [readback sumber trace](evidence/sandbox-video-traces.json).

Dua pengambilan review sebelumnya dihentikan dan tidak diekspor: `022808`
karena menu bahasa terpotong setelah scroll awal, serta `022944` karena batas
durasi recorder. Posisi awal Ringkasan diperbaiki dan jeda review diringkas
sebelum pengambilan selesai ini. Rekaman akhir mempertahankan waktu interaksi
asli tanpa mempercepat footage atau mengganti hasil run.

Untuk merekam review yang sama, dari `checkoutguard_BE`:

```sh
go run ./cmd/demorecord -review 7efb48b0-250a-4907-9c6f-718342588a55 -review-retest 6d374209-36ea-4a39-9c96-38f49b0fa3e8
```

Kemudian render direktori baru yang dicetak recorder dengan
`python3 scripts/render-sandbox-demo.py tmp/video/DIREKTORI-REKAMAN-BARU`.
Mode review memvalidasi hasil dan order tersimpan serta menolak sumber yang
memiliki step error. MP4 dan frame mentah tetap di `tmp/` yang diabaikan Git;
video belum diunggah atau disubmit.

## Voice-over untuk video review saat ini — Indonesia dan Inggris

Naskah berikut sudah dicocokkan dengan manifest video serta readback trace
terbaru. Pilih satu bahasa untuk rekaman suara pribadi. Timestamp mengikuti
video review berdurasi 156,53 detik dan menjadi panduan waktu membaca.
Beri jeda jika narasi selesai sebelum pergantian layar. Video tetap tanpa audio;
pembaruan ini hanya teks. Tidak ada klaim checkout baru, gateway sudah stabil,
atau portal Entire/cfo.ai didemonstrasikan dalam segmen ini.

### 0:00–0:10 — Pembukaan / Introduction

**Indonesia:**

“Ini CheckoutGuard, aplikasi untuk memeriksa alur checkout. Dalam video ini,
saya meninjau dua hasil pemeriksaan sandbox yang sudah tersimpan.”

**English:**

“This is CheckoutGuard, an application for checking checkout flows. In this
video, I’m reviewing two previously saved sandbox test results.”

### 0:10–0:45 — Bahasa dan tujuan / Language and goal

**Indonesia:**

“Dashboard mendukung bahasa Inggris dan Indonesia. Saya beralih ke Indonesia,
lalu mengetik tujuan pemeriksaan: membeli satu tote, memverifikasi total
checkout, dan memastikan tepat satu pesanan tersimpan.

Saya memilih keranjang dengan harga alternatif. Semua transaksi pada sandbox
merupakan simulasi, tanpa pembayaran nyata. Dalam review ini, saya tidak
memulai checkout baru.”

**English:**

“The dashboard supports English and Indonesian. I switch to Indonesian and
enter a test goal: buy one tote, verify the checkout total, and confirm that
exactly one order is saved.

I select the alternate-price basket. All sandbox transactions are simulated,
with no real payment. This review does not start a new checkout.”

### 0:45–1:31 — Langkah dan screenshot / Steps and screenshots

**Indonesia:**

“Melalui menu Bukti, kita bisa meninjau keputusan agent dan screenshot dari
setiap langkah. Skenario sandbox ini sengaja memiliki total yang salah untuk
menguji deteksi kesalahan harga.

Agent mengamati halaman, mengisi jumlah barang, menambahkan produk ke
keranjang, lalu mengisi email demo dan melanjutkan checkout.

Hasil verifikasi total dan pemeriksaan pesanan juga tersimpan. Kita dapat
memilih setiap langkah untuk melihat kondisi browser saat itu, sehingga hasil
pemeriksaan dapat ditelusuri melalui bukti yang nyata.”

**English:**

“The Evidence menu lets us review the agent’s decisions and a screenshot from
each step. This sandbox scenario intentionally contains an incorrect total to
test pricing defect detection.

The agent observes the page, enters the quantity, adds the product to the
basket, then fills in the demo email and continues through checkout.

The total verification and order check results are also saved. Selecting a
step shows the browser state at that moment, making the result traceable
through actual evidence.”

### 1:31–1:46 — Temuan harga / Pricing finding

**Indonesia:**

“Total seharusnya sembilan ribu delapan ratus rupiah, tetapi halaman menampilkan
sebelas ribu lima puluh. CheckoutGuard mencatat selisih ini sebagai temuan
terverifikasi, lengkap dengan screenshot dan bukti pesanan.”

**English:**

“The expected total is nine thousand eight hundred rupiah, but the page shows
eleven thousand fifty. CheckoutGuard records this difference as a verified
finding, supported by a screenshot and order evidence.”

### 1:46–1:57 — Riwayat dan retest / History and retest

**Indonesia:**

“Saya membuka riwayat dan memilih hasil retest yang sudah tersimpan. Retest
menggunakan tujuan dan keranjang yang sama pada toko simulasi yang telah
diperbaiki.”

**English:**

“I open the history and select the saved retest. It uses the same goal and
basket against the corrected synthetic store.”

### 1:57–2:25 — Verifikasi retest / Retest verification

**Indonesia:**

“Pada retest, total checkout sesuai dengan perhitungan: sembilan ribu delapan
ratus rupiah. Pemeriksaan pesanan juga mengonfirmasi tepat satu order
tersimpan pada run ini.

Dashboard menyediakan hasil tool, jumlah token, dan estimasi biaya API. Kedua
trace sumber telah diverifikasi melalui neatlogs, dengan jumlah token yang
cocok dengan catatan backend.”

**English:**

“In the retest, the checkout total matches the calculation: nine thousand
eight hundred rupiah. The order check also confirms exactly one saved order
for this run.

The dashboard provides tool results, token usage, and estimated API costs.
Both source traces have been verified through neatlogs, with token totals
matching the backend records.”

### 2:25–2:32 — Perbandingan / Comparison

**Indonesia:**

“Perbandingan ini memperlihatkan total salah pada pemeriksaan awal, lalu hasil
retest yang sudah sesuai.”

**English:**

“This comparison shows the incorrect initial total alongside the corrected
retest result.”

### 2:32–2:36,5 — Penutup / Closing

**Indonesia:**

“CheckoutGuard membantu memverifikasi checkout dengan bukti yang dapat
ditelusuri.”

**English:**

“CheckoutGuard makes checkout verification clear and traceable.”

## Draft sandbox sebelumnya — digantikan atas permintaan user

Sesuai pilihan user, rekaman pertama hanya menjalankan **sandbox**. Tidak ada
browser Shopee/Tokopedia yang dibuka atau run marketplace yang dibuat. Video
marketplace akan menjadi laporan terpisah tentang kendala akses, bukan klaim
checkout berhasil. Untuk video submission lengkap, sisakan segmen bukti penggunaan
neatlogs, Entire, dan cfo.ai; logo di sidebar bukan bukti penggunaan partner.

Draft pertama memakai dashboard asli dengan model aktif, tanpa replay hasil
lama. Draft itu memiliki anotasi bar bawah dan narasi sintetis, sehingga tidak
sesuai preferensi user terbaru. **Jangan gunakan draft ini sebagai hasil akhir.**
Urutan dan waktu tunggu kedua run dipertahankan sebagai bukti historis.

| Bagian | Run | Hasil |
| --- | --- | --- |
| Total salah | `7efb48b0-250a-4907-9c6f-718342588a55` | Total seharusnya 9.800 IDR, ditampilkan 11.050 IDR; satu order; nol step error |
| Retest fixture benar | `6d374209-36ea-4a39-9c96-38f49b0fa3e8` | Lulus, total 9.800 IDR; satu order; nol step error |

Trace masing-masing: `4102ea1e8e143672a7e41c813b14418a` dan
`aedd5ac65bb1a9895f7e2d894893af98`. Hosted readback keduanya kini terverifikasi:
15 span per trace, masing-masing tujuh LLM call, struktur valid dan token cocok
dengan run lokal. [Readback](evidence/sandbox-video-traces.json). Pengujian
sebelumnya tetap tersedia di [EVIDENCE.md](EVIDENCE.md).

File video dan manifest lokal berada di `tmp/video/20261011-215823/` (diabaikan
Git). Nama ekspor: `CheckoutGuard-sandbox-demo.mp4`; teks narasi:
`narration-id.txt`; subtitle tambahan: `subtitles-en.srt`.
Ekspor draft lama sudah diverifikasi: **1 menit 42,5 detik**, 1920 × 1080, 30 fps,
video H.264 dan audio AAC. Dekode seluruh video berhasil; sampel frame tiap
bagian sudah diperiksa. Audio memiliki sinyal, dengan peak −0,9 dB.
SHA-256 draft bersuara: `c953793cd252d9a27d57cb7129551d0a370c586c63f544f210d8fc645e57c74f`.
Manifest video saat ini menunjuk ekspor review tanpa suara, bukan draft ini.
File lokal belum merupakan link video publik atau submission.

Recorder `checkoutguard_BE/cmd/demorecord` berhenti pada step error, error UI,
exception browser, kegagalan backend, hasil verifikasi yang tidak sesuai, atau
batas durasi. Recorder membatalkan run miliknya yang masih aktif dan tidak
mengulang run secara otomatis. Temuan `incorrect_total` pada fixture rusak adalah
hasil QA yang diharapkan; skenario timeout tidak dijalankan dalam rekaman ini.
Encoder menolak manifest yang dihentikan atau memiliki step error.

Untuk rekaman baru yang memang ingin dijalankan lagi oleh user:

```sh
cd checkoutguard_BE
go run ./cmd/demorecord
```

Kemudian dari root project, gunakan direktori baru yang dicetak recorder:

```sh
python3 scripts/render-sandbox-demo.py tmp/video/DIREKTORI-REKAMAN-BARU
```

Perintah render hanya mengolah rekaman yang sudah ada. Membuat rekaman baru
membutuhkan backend/dashboard hidup, database siap, konfigurasi model aktif,
Chrome, serta tidak ada run lain yang sedang berjalan. Ekspor membutuhkan
ffmpeg dan ffprobe. Encoder baru menghasilkan `CheckoutGuard-sandbox-silent.mp4`
dengan 1920 × 1080, 30 fps, tanpa audio/voice-over/subtitle; hanya manifest dua
run lengkap tanpa step error yang dapat diekspor.

## Rencana video submission lengkap

Team **FLdevBoo** · [Repository publik](https://github.com/fadillaoegi/checkoutGuard).
Narasi bahasa Indonesia; gunakan satu bahasa dashboard secara konsisten.

Fokus: agent menjalankan checkout QA pada sandbox, mendeteksi total salah,
memverifikasi retest, dan menangani timeout tanpa pesanan ganda. Shopee dan
Tokopedia tetap disebut sebagai fitur eksperimental yang belum berhasil.
Video maksimal tiga menit dan harus menunjukkan penggunaan nyata neatlogs,
Entire, serta cfo.ai. [Ketentuan resmi neatHack](https://neatlogs.com/hackathon).

## Siapkan sebelum merekam

- Buka hasil run sandbox dan bukti di [EVIDENCE.md](EVIDENCE.md). Rekam juga
  tujuan, klik mulai, dan aksi browser. Beri label **“Waktu tunggu dipotong”**
  pada bagian eksekusi yang disingkat; jangan mengubah klaim durasi run.
- Siapkan trace asli di neatlogs, checkpoint implementasi di Entire, dan
  rencana bisnis asli cfo.ai. Logo atau prompt saja belum membuktikan penggunaan.
- Push user mengirim commit `a134e4a` dan checkpoint
  `01M4NKBNCP60MN7M5BXS143W16`. Metadata publik cocok dengan 35 file commit;
  [hasil verifikasi](evidence/entire-checkpoint.json). Pilih keputusan/perubahan
  yang sudah diperiksa pada [ENTIRE-REVIEW.md](ENTIRE-REVIEW.md) saat merekam.
- Hasil cfo.ai sudah tersedia: [PDF dan model Excel asli](CFO-RESULTS.md).
  Buka sesi asli cfo.ai dan bagian business plan; siapkan PDF halaman 1–2 dan 4
  sebagai bukti ekspor. Dashboard juga menyediakan tautan rencana bisnis.
- Tutup notifikasi/tab pribadi dan sembunyikan credential, email akun, billing,
  alamat, serta data pelanggan dari rekaman.

## Urutan dan narasi siap dibaca

### 0:00–0:12 — Produk

**Layar:** dashboard, tujuan pemeriksaan, pilihan sandbox.

> “CheckoutGuard membantu developer memeriksa checkout sebelum rilis. Berikan
> tujuan; agent memilih aksi browser, memeriksa total, dan memverifikasi pesanan
> dengan bukti yang tersimpan.”

**Teks:** `Checkout QA · Goal → Browser → Verification`.

### 0:12–0:30 — Arsitektur dan eksekusi

**Layar:** mulai pemeriksaan dan aktivitas browser. Sisipkan diagram
`React + Redux → Go agent → Model / Chrome / PostgreSQL`, dengan neatlogs
mencatat trace.

> “Dashboard memakai React dan Redux, dengan backend Go berarsitektur bersih.
> Model memilih langkah berdasarkan observasi. Chrome menjalankan aksi;
> PostgreSQL memverifikasi pesanan, dan neatlogs merekam panggilan serta hasilnya.”

Tampilkan satu aksi dan hasil observasinya; labeli potongan waktu tunggu.

### 0:30–0:53 — Total salah ditemukan

**Layar:** temuan dan [screenshot pricing defect](evidence/pricing-defect.png).
Run `b87b1ae0-83a0-4355-b3ac-850d3fb6a4be`.

> “Harga dikurangi diskon, ditambah ongkir, seharusnya sembilan ribu delapan
> ratus rupiah. Halaman menampilkan sebelas ribu lima puluh. Agent membandingkan
> rincian harga dan mencatat kesalahan beserta screenshot sebagai bukti.”

**Teks besar:** `8,750 − 1,250 + 2,300 = 9,800 IDR`;
`Displayed: 11,050 IDR · Verified finding`.

### 0:53–1:10 — Retest setelah perbaikan

**Layar:** hubungan run awal dengan retest, total dan jumlah pesanan.
Run `6b111e4c-30ff-44e2-a6cc-0c533eb3af71`;
[screenshot](evidence/corrected-retest.png).

> “Retest memakai tujuan dan keranjang yang sama pada fixture yang sudah
> diperbaiki. Total kini sembilan ribu delapan ratus rupiah, dan tepat satu
> pesanan simulasi tersimpan. Agent memverifikasi perbaikan ini.”

**Teks:** `Same goal + basket · Corrected fixture · 1 saved order`.
Perbaikan tersedia melalui fixture berbeda; agent tidak mengedit source toko.

### 1:10–1:29 — Recovery tanpa pesanan ganda

**Layar:** pesan timeout asli, langkah `order_status`, satu order dan satu submit.
Run `b74d8ccb-9559-4a61-9c52-7eb3ae713815`;
[screenshot](evidence/payment-timeout.png).

> “Skenario berikutnya menampilkan timeout setelah pesanan tersimpan. Agent
> memeriksa status pesanan, menemukan satu pesanan, lalu selesai tanpa submit
> kedua. Ini membuktikan recovery tanpa pesanan ganda.”

**Teks:** `Timeout → order_status → 1 order / 1 submit`.
Screenshot tetap menunjukkan timeout; hasil tool/database membuktikan order.
Ini transaksi sintetis tanpa pembayaran uang nyata.

### 1:29–1:59 — neatlogs: diagnosis dan perbaikan

**Layar:** trace model timeout, pohon span workflow/model/tool, lalu trace retest
dan panel estimasi biaya aplikasi. Sorot perubahan adapter model bila terbaca.

> “Neatlogs menunjukkan request mencapai batas client 45 detik. Apertus adalah
> model yang diminta; respons sukses menyebut SEA-LION. Kami memisahkan keduanya
> dan merekam error tiap request, lalu membatasi retry. Retest menyelesaikan satu
> pesanan. Biaya yang tampil adalah estimasi dari token yang diketahui; provider
> masih bisa timeout.”

**Tab/bukti:**

- Timeout: `6bb72be9ea60a2c2a5854f3c89fd377f`.
- Sebelum konfigurasi tarif: `7a97b78bbc758fb1cbadf231ccb0f210`.
- Sesudah konfigurasi tarif: `db1fd6a1e4cb8d7159c6cefc1c0b3daa`;
  12.361 input + 199 output token, satu order, **US$0,00318975 estimasi**.
  [Readback](evidence/priced-retest-trace.json).

- Label/diagnostik sesudah perbaikan: `869bdb1c0d98cd794c3a895292fe23c7`;
  pada output workflow cari `Model Calls`, lalu sorot satu `requested_model`,
  `returned_model`, `timeout_ms` dan `error_kind` yang terkait ID span.
  [Readback terverifikasi](evidence/model-diagnostics-trace.json).

**Teks:** `Unknown cost → US$0.00318975 estimate · Not an invoice`.
Tampilkan biaya dari aplikasi jika UI neatlogs tidak memuat angka itu.
Jangan mengklaim biaya menurun atau provider pulih. Perbedaan permintaan
Apertus/label respons SEA-LION dijelaskan dalam EVIDENCE.md.

### 1:59–2:19 — Entire: konteks implementasi

**Layar:** commit implementasi dan checkpoint terkait. Buka konteks keputusan
serta satu perubahan file yang benar-benar tercatat.

**Narasi ketika konteks keputusan dan perubahan file sudah ditampilkan:**

> “Entire menangkap konteks kerja Codex dan menghubungkannya dengan commit
> implementasi. Reviewer bisa melihat tujuan perubahan, konteks keputusan,
> dan file yang berubah. Checkpoint ini ikut dipublikasikan bersama repository.”

Gunakan commit `a134e4a` dan checkpoint `01M4NKBNCP60MN7M5BXS143W16`.
[Metadata publik](https://github.com/fadillaoegi/checkoutGuard/blob/bef8f90aab69269462123a995948e574b08f6acc/metadata.json)
mengidentifikasi Codex dan 35 file yang sama dengan commit. Transkrip tersedia;
record keputusan 513 dan patch 519 sudah dicocokkan dengan arsip asli, termasuk
output terkait. Lihat [panduan inspeksi](ENTIRE-REVIEW.md) dan tampilkan perubahan
`model.go`/`tracing.go`. Arsip juga menyimpan identitas sesi; sembunyikan metadata
akun dan path pribadi saat merekam. Jangan memakai anchor kosong sebagai bukti.

### 2:19–2:44 — cfo.ai: rencana bisnis

**Layar:** hasil asli Ari di cfo.ai: pricing/kuota, costs, revenue scenarios,
dan metode runway. Sorot paket Team, biaya per run dan tabel sensitivitas yang
sudah dihasilkan; kemudian tunjukkan ekspor asli PDF atau Excel.

**Narasi:**

> “Cfo.ai menyarankan paket Team 99 dolar per bulan, dengan estimasi biaya
> variabel sekitar dua sen per pemeriksaan. Model mencakup tiga skenario
> pendapatan dan sensitivitas support. Semua angka bisnis adalah asumsi.
> Metode runway tersedia; perhitungannya menunggu data kas aktual.”

**Teks:** `Pricing · Costs · Revenue · Runway · Assumptions labelled`.
Tambahkan label angka biaya **US$0,0190856/run · Planning estimate**; support
bulanan per langganan, payment fee dan biaya tetap masih dihitung terpisah.
Tabel sensitivitas PDF halaman 2 menunjukkan support ×2 menurunkan hasil
operasional base Oktober 2028 dari sekitar US$4.209 menjadi US$2.117.
Tampilkan tabel hasil yang tersedia; jangan mengklaim kontrol simulator sudah
diuji. PDF halaman 4 membahas runway yang belum dihitung. Jangan menyebut
defisit kumulatif US$59.175 sebagai saldo, runway atau kebutuhan dana pasti.
Screenshot browser pribadi dari chat tidak masuk paket publik; rekam panel
rencana tanpa sidebar akun, bookmark dan tab pribadi.

### 2:44–2:51 — Batas marketplace

**Layar:** kartu status singkat, tanpa menunggu CAPTCHA/login live.

> “Shopee dan Tokopedia masih eksperimental; pengujian belum berhasil.
> Bukti alur lengkap berasal dari sandbox.”

**Teks:** `Shopee: verification blocked`;
`Tokopedia: checkout preparation not verified`;
`Sandbox: verified synthetic checkout QA`.

Jangan menyebut penyebab Tokopedia sebagai CAPTCHA tanpa bukti. Beralih ke
sandbox tidak memperbaiki akses marketplace. Adapter berhenti sebelum konfirmasi
pesanan akhir. Simpan detail kegagalan di laporan; screenshot hanya jika bukti
asli tersedia dan bebas informasi pribadi.

### 2:51–2:55 — Penutup

**Layar:** logo, team, URL repository.

> “CheckoutGuard, oleh FLdevBoo. Kode dan bukti pengujian tersedia di GitHub.”

**Teks:** `github.com/fadillaoegi/checkoutGuard · #neatHack`.

## Cek video hasil ekspor

- Maksimal **3:00**, termasuk intro/outro. Target 2:55 adalah rencana edit,
  bukan durasi video yang sudah diverifikasi. Ringkas narasi bila diperlukan.
- Angka, run, trace dan checkpoint sesuai layar; semua order bersifat sintetis.
- Ketiga partner ditunjukkan dengan hasil nyata; akses link video diuji tanpa login.
- Tidak ada klaim checkout marketplace berhasil, source toko diperbaiki otomatis,
  atau timeout provider sudah terselesaikan.
- Untuk pemirsa internasional, tambahkan subtitle Inggris jika memungkinkan.

## Building-in-public draft — belum dipublikasikan

> Building CheckoutGuard for #neatHack with FLdevBoo: a goal-driven browser agent
> for checkout QA. Our synthetic demo detects an incorrect total, verifies the
> corrected retest, and recovers from a timeout without a duplicate order.
> Shopee/Tokopedia preparation remains experimental and unsuccessful in our
> tests. We document limitations and use neatlogs to investigate them.
> Code and evidence: https://github.com/fadillaoegi/checkoutGuard

Tambahkan hasil cfo.ai: “The actual cfo.ai plan models capped pricing, variable
run costs and three revenue scenarios; cash runway awaits real cash inputs.”
Tambahkan link checkpoint setelah isi dan aksesnya diverifikasi. Simpan
link progress di [SUBMISSION-CHECKLIST.md](SUBMISSION-CHECKLIST.md) setelah user
mempublikasikannya. Submission final harus mengutip posting resmi neatlogs dan
memuat team, proyek, repo, serta link demo. Membuka draft belum berarti submit.
