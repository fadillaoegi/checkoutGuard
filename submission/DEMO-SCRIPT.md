# Skrip video demo CheckoutGuard — target 2 menit 55 detik

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
- Push user sudah mengirim commit `380ed7a` dan satu checkpoint ref Entire.
  Periksa isi checkpoint yang terhubung sebelum merekam; gunakan ID aktual.
- Hasil cfo.ai masih pending. Gunakan [CFO-PROMPT.md](CFO-PROMPT.md), lalu simpan
  hasilnya. Jangan menerbitkan video final dengan segmen partner kosong.
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

> “Di neatlogs, workflow memuat panggilan model dan tool beserta waktu, input,
> dan hasil. Trace gagal menunjukkan model timeout. Kami menambahkan penanganan
> error dan batas retry; keandalan provider masih menjadi batasan. Perbaikan
> observabilitas yang terukur: biaya sebelumnya tidak tersedia, lalu retest
> melaporkan estimasi sekitar nol koma nol nol tiga dua dolar.”

**Tab/bukti:**

- Timeout: `6bb72be9ea60a2c2a5854f3c89fd377f`.
- Sebelum konfigurasi tarif: `7a97b78bbc758fb1cbadf231ccb0f210`.
- Sesudah konfigurasi tarif: `db1fd6a1e4cb8d7159c6cefc1c0b3daa`;
  12.361 input + 199 output token, satu order, **US$0,00318975 estimasi**.
  [Readback](evidence/priced-retest-trace.json).

**Teks:** `Unknown cost → US$0.00318975 estimate · Not an invoice`.
Tampilkan biaya dari aplikasi jika UI neatlogs tidak memuat angka itu.
Jangan mengklaim biaya menurun atau provider pulih. Perbedaan permintaan
Apertus/label respons SEA-LION dijelaskan dalam EVIDENCE.md.

### 1:59–2:19 — Entire: konteks implementasi

**Layar:** commit implementasi dan checkpoint terkait. Buka konteks keputusan
serta satu perubahan file yang benar-benar tercatat.

**Narasi setelah isi checkpoint diverifikasi:**

> “Entire menangkap konteks kerja Codex dan menghubungkannya dengan commit
> implementasi. Reviewer bisa melihat tujuan perubahan, konteks keputusan,
> dan file yang berubah. Checkpoint ini ikut dipublikasikan bersama repository.”

Gunakan commit `380ed7a` dan ID checkpoint aktual. Jangan memakai anchor kosong
sebagai bukti implementasi atau memperkirakan isi checkpoint.

### 2:19–2:44 — cfo.ai: rencana bisnis

**Layar:** hasil asli Ari di cfo.ai: pricing/kuota, costs, revenue scenarios,
runway; tunjukkan perubahan asumsi dan hasil.

**Narasi setelah hasil nyata tersedia:**

> “Kami memakai cfo.ai untuk merencanakan harga, biaya per pemeriksaan,
> pendapatan, dan runway. Pelanggan serta pendapatan proyeksi diberi label
> asumsi. Saat jumlah pemeriksaan atau harga token berubah, kami mengecek
> perubahan biaya dan margin. Biaya hosting dan support ikut diperhitungkan.”

**Teks:** `Pricing · Costs · Revenue · Runway · Assumptions labelled`.
Uji 1.000 versus 10.000 run/bulan, lalu dua kali harga token dengan input lain
tetap. Bila saldo belum diketahui, tampilkan rumus runway/input yang diperlukan.
Minta revisi jika biaya hosting/support belum tercakup sebelum memakai narasi.
Simpan ekspor/screenshot asli; prompt kita sendiri bukan keluaran cfo.ai.

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

Tambahkan hasil cfo.ai dan link checkpoint setelah bukti aktual siap. Simpan
link progress di [SUBMISSION-CHECKLIST.md](SUBMISSION-CHECKLIST.md) setelah user
mempublikasikannya. Submission final harus mengutip posting resmi neatlogs dan
memuat team, proyek, repo, serta link demo. Membuka draft belum berarti submit.
