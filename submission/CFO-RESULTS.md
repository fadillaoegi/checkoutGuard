# Hasil cfo.ai — CheckoutGuard

Ekspor asli diterima dari user pada **11 Oktober 2026**. Rencana ini membahas
pricing, biaya, proyeksi pendapatan dan metode runway untuk bisnis checkout QA.
Tidak ada pelanggan berbayar atau pendapatan aktual yang diklaim.

## Berkas untuk reviewer dan demo

- [Rencana bisnis PDF, 7 halaman](../checkoutguard_FE/public/business/CheckoutGuard-business-plan.pdf).
- [Model Excel asli, 11 sheet](../checkoutguard_FE/public/business/CheckoutGuard-cfoai-model.xlsx).
- [Rencana bisnis Markdown asli](business/CheckoutGuard-business-plan.md).
- [Manifest asli dari cfo.ai](business/cfoai-original-manifest.json) dan
  [provenance salinan repository](business/artifact-provenance.json).
- [Hasil pemeriksaan independen lokal](evidence/cfo-validation.json).

Dashboard menyediakan tautan **cfo.ai → Rencana bisnis / Business plan**.
Saat frontend berjalan, PDF tersedia di `/business/CheckoutGuard-business-plan.pdf`
dan Excel di `/business/CheckoutGuard-cfoai-model.xlsx`.

PDF, Excel, Markdown dan manifest asli disalin tanpa perubahan byte.
`Financial model (3 pages).xlsx` identik dengan `CheckoutGuard-cfoai-model.xlsx`,
sehingga cukup menyertakan satu workbook. ZIP sumber dan screenshot browser
pribadi tidak disertakan dalam paket repository. Tidak ada publikasi Git yang
dilakukan oleh coding agent.

## Temuan utama — seluruh angka bisnis adalah asumsi

| Bagian | Hasil rencana |
|---|---|
| Paket Developer | US$29/bulan; 1 workspace; 100 run |
| Paket Team yang disarankan | US$99/bulan; 3 workspace; 500 run gabungan; batas kuota dan persetujuan upgrade |
| Paket Agency | US$299/bulan; 10 workspace; 2.000 run gabungan |
| Rata-rata estimasi token dari tiga run sintetis | US$0,00315467/run berdasarkan tarif katalog; bukan tagihan |
| Biaya variabel lengkap | US$0,0190856/run selesai; mencakup allowance retry, browser, DB, evidence, telemetry dan support per run |
| Biaya tambahan per langganan | Support bulanan dan asumsi payment fee; belum termasuk dalam biaya variabel per run |
| Biaya tetap | US$6.800/bulan, termasuk asumsi kompensasi founder US$6.000 |
| Base break-even | 96 langganan pada komposisi asumsi; sekitar US$8.650,67 pendapatan bulanan |
| Bulan operasional base pertama nonnegatif | Maret 2028 |
| Sensitivitas support ×2 | Hasil operasional Oktober 2028 turun dari US$4.209 menjadi US$2.117 |

Pada pemakaian kuota 50%, margin kontribusi paket Team adalah sekitar **79,25%**.
Angka **79,16%** pada ringkasan adalah margin campuran pelanggan base. Keduanya
belum mengurangi biaya tetap; margin Team pada pemakaian kuota penuh sekitar
74,43%. Angka ini belum divalidasi lewat pelanggan maupun invoice.

| Skenario | Pendapatan bulanan Oktober 2027 | Pendapatan bulanan Oktober 2028 | Hasil operasional Oktober 2028 |
|---|---:|---:|---:|
| Conservative | US$1.864 | US$2.871 | −US$4.527 |
| Base | US$5.586 | US$13.907 | US$4.209 |
| Growth | US$12.342 | US$40.765 | US$25.464 |

Perbedaan skenario berasal dari asumsi conversion, pertumbuhan prospek dan
churn. Billing satu bulan penuh untuk jumlah langganan akhir bulan merupakan
asumsi optimistis, sehingga perlu diganti dengan tanggal langganan aktual.
Shopee/Tokopedia tidak menjadi dasar klaim checkout berhasil atau penawaran
berbayar awal; fokus rencana adalah staging yang dimiliki atau diizinkan.

## Runway dan batas bukti

Runway numerik **belum dihitung**. PDF halaman 4 memuat rumus kas tersedia dibagi
net cash burn bulanan positif, serta kebutuhan saldo bertanggal, penerimaan,
pembayaran, pajak, kewajiban dan pembiayaan yang benar-benar tersedia.
Puncak defisit operasional kumulatif base sekitar **US$59.175** bukan saldo kas,
runway atau kebutuhan pendanaan terverifikasi. Untuk menambahkan runway numerik,
berikan data kas tersebut ke cfo.ai dan simpan hasil revisi asli.

[Input dan kalkulator runway](RUNWAY-INPUTS.md) serta
[prompt revisi untuk cfo.ai](CFO-RUNWAY-FOLLOWUP.md) sudah disiapkan. Status
tetap menunggu input kas user; tidak ada angka saldo yang diisi otomatis.

cfo.ai menerima angka runtime yang ditempel ke prompt; `EVIDENCE.md` belum
dilampirkan di sesi tersebut. Pemeriksaan lokal berikut mencocokkan tiga baris
workbook dengan [run sintetis yang sudah dikurasi](evidence/runs.json), tanpa
mengubah klaim pemeriksaan yang dilakukan cfo.ai.

Pemeriksaan lokal membandingkan **409 angka**, termasuk **75 periode** dari tiga
skenario dan **5 sensitivitas**, dengan perhitungan independen. Tidak ditemukan
Excel error cell, external relationship, macro member, atau kecocokan pola
credential/akun yang diperiksa. Semua halaman PDF dibaca dan ditinjau secara
visual. Ini pemeriksaan cached export dan aritmetika, bukan recalculation Excel
native atau pengujian kontrol simulator cfo.ai. Screenshot Canvas light/dark
tidak tersedia sebagai file dalam ZIP sumber; render PDF lokal tidak menggantikan
bukti screenshot simulator tersebut.

Jalankan kembali pemeriksaan workbook tanpa jaringan atau Git:

```sh
python3 scripts/review-cfo-model.py
python3 scripts/test-publication-audit.py
python3 scripts/audit-publication.py
```

Ekspor ini menyediakan keluaran bisnis untuk bagian cfo.ai dalam
[ketentuan neatHack](https://neatlogs.com/hackathon). Tetap tunjukkan penggunaan
nyata layanan dalam video, verifikasi checkpoint Entire, lalu selesaikan
publikasi dan submission sesuai [checklist](SUBMISSION-CHECKLIST.md).
