# Input rencana bisnis untuk cfo.ai — draft

Dokumen ini adalah draft input historis, bukan keluaran cfo.ai. User telah
menyerahkan ekspor asli pada 11 Oktober 2026; lihat
[hasil dan pemeriksaan](../submission/CFO-RESULTS.md). Sampel lima run di bawah
berbeda dari tiga run terbaru pada [prompt yang digunakan](../submission/CFO-PROMPT.md)
dan workbook asli. Pertahankan konteks sampel saat membandingkan rata-rata.

## Produk dan pelanggan

CheckoutGuard membantu developer memeriksa checkout sebelum merilis perubahan.
MVP menjalankan browser berbasis tujuan, memverifikasi total/order, menyimpan
evidence dan membandingkan retest. Calon pengguna: developer toko dan agensi
yang memelihara checkout. Kebutuhan, willingness to pay dan pelanggan berbayar
belum divalidasi. Tidak ada klaim revenue aktual.

## Input operasional yang sudah terukur

Lima evaluasi sintetis memakai total 62.406 token, rata-rata 12.481,2 token/run,
masing-masing tujuh panggilan model. Satu browser aktif pada satu waktu; durasi
21,79–54,81 detik/run dalam sampel. Setiap run menyimpan delapan screenshot.
Ukuran penyimpanan, kapasitas concurrent produksi dan support belum diukur.

Tarif katalog Public AI untuk label respons SEA-LION diperiksa pada 11 Oktober:
USD 0,25 per juta input dan 0,50 per juta output token. Ini dasar estimasi,
bukan invoice terverifikasi. Retest baru 8b2d9a3c melaporkan 12.361 input + 199
output, dengan estimasi aplikasi USD 0,00318975. Histori sebelumnya tetap null.
Request timeout bisa tidak melaporkan usage; biaya yang tidak diketahui jangan
dianggap nol. Tarif compute hosting, database, object storage dan support masih
belum diketahui. Jangan memakai tarif Apertus untuk respons berlabel SEA-LION.
Sumber: [katalog resmi](https://platform.publicai.co/models).

## Isi yang perlu dihasilkan melalui cfo.ai

1. Pilihan harga berbasis workspace/run dan kuota, dengan batas pemakaian jelas.
2. Model biaya per run: token input/output × tarif model respons, ditambah browser,
   database, storage evidence dan telemetry. Pisahkan estimasi dari tagihan aktual.
3. Skenario konservatif/base/tinggi untuk jumlah workspace, run/workspace, churn
   dan support. Semua angka pelanggan/revenue baru harus berlabel asumsi.
4. Margin kontribusi, biaya tetap bulanan, break-even dan runway berdasarkan
   saldo/budget yang diberikan user; jangan mengisi saldo atau revenue fiktif.
5. Catatan asumsi, sensitivitas kenaikan penggunaan, dan revisi setelah data
   pelanggan serta tarif aktual tersedia.

Hasil rencana, asumsi dan ekspor asli kini disertakan melalui CFO-RESULTS.md.
Rekam penggunaan nyata cfo.ai dalam demo. Aktivasi akun/billing dikelola user;
coding agent tidak melakukan operasi akun atau pembayaran.
