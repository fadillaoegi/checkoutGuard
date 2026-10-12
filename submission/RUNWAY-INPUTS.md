# Melengkapi runway numerik

Hasil cfo.ai asli sudah memuat metode runway, tetapi belum memiliki saldo kas
bertanggal. Nilai kas dan cash flow berikut diperlukan sebelum angka runway
dapat dihitung. Template tetap kosong; input yang belum diketahui tidak dianggap
nol. [Contoh input](../docs/examples/runway-input.json.example) tidak memuat
data keuangan pribadi.

## Input yang diperlukan

| Input | Yang diisi |
|---|---|
| Dasar kas | `actual` untuk saldo yang benar-benar tersedia, atau `simulation` untuk budget simulasi yang ditentukan user |
| Tanggal saldo | `YYYY-MM-DD`, sesuai tanggal saldo tersebut |
| Mata uang | Misalnya `IDR` atau `USD`; seluruh input harus dalam mata uang yang sama |
| Kas bebas digunakan | Dana likuid yang dialokasikan untuk CheckoutGuard; sudah mengecualikan dana yang tidak boleh dipakai |
| Buffer kas minimum | Batas saldo yang ingin dipertahankan; isi 0 hanya jika memang tidak ada buffer |
| Penerimaan kas bulanan | Uang yang diperkirakan benar-benar diterima, bukan revenue yang belum tertagih |
| Pembayaran operasi bulanan | Hosting, provider, DB, storage, telemetry, support dan biaya operasi lain; mengecualikan kategori terpisah berikut |
| Pembayaran founder | Pengambilan gaji/dana founder yang benar-benar direncanakan |
| Pajak, capex, debt service, refund | Masing-masing diisi terpisah; tidak dihitung lagi dalam pembayaran operasi |
| Pembiayaan masa depan | Konfirmasi jika tidak ada pembiayaan yang diasumsikan; bila ada pembiayaan committed bertanggal, gunakan jadwal kas cfo.ai |

Tidak perlu mengirim rekening, kartu, credential atau screenshot banking.
Cukup berikan jumlah, mata uang, tanggal dan asumsi pembayaran.

Asumsi asli cfo.ai adalah **US$6.800 biaya tetap per bulan**, termasuk
**US$6.000 kompensasi founder** dan US$800 biaya tetap lain. Ini angka rencana,
bukan pembayaran wajib atau invoice. Jika founder tidak mengambil gaji,
pembayaran founder dapat diisi 0 berdasarkan keputusan user. Biaya variabel,
support langganan serta payment fee juga perlu dimasukkan ke cash payments
saat benar-benar dibayarkan; tidak cukup memakai biaya tetap saja jika ada run
dan pelanggan. Jangan mengubah semua proyeksi USD menjadi IDR tanpa kurs yang
ditentukan beserta tanggalnya.

## Perhitungan lokal

Salinan input privat sudah disiapkan di `tmp/runway-input.json`. Isi field
`null`, termasuk jumlah 0 yang benar-benar dipilih. Jalankan dari root proyek:

```sh
python3 scripts/calculate-runway.py --input tmp/runway-input.json
```

Hasil privat tersimpan di `tmp/runway-result.json`. Rumus yang dipakai:

```text
Net cash burn bulanan = pembayaran kas bulanan − penerimaan kas bulanan
Runway sampai buffer = max(0, kas bebas − buffer minimum) / net cash burn positif
Kas akhir periode = kas awal + penerimaan − pembayaran
```

Ini model cash flow bulanan konstan tanpa pembiayaan masa depan. Periode dimulai
pada tanggal saldo dan berakhir pada tanggal ulang bulan, dengan penyesuaian
hari akhir bulan. Hasil akhir periode yang menyentuh buffer tidak mengklaim
tanggal persis kehabisan kas di dalam periode. Jika burn nol atau negatif dan
kas masih di atas buffer, runway numerik menjadi `null` dengan status kas tidak
berkurang di bawah asumsi konstan; bukan klaim bisnis memiliki runway tak terbatas.
Jika kas sudah menyentuh atau di bawah buffer, runway 0.

Input kosong menghasilkan `needs_input`, bukan runway palsu. Hasil kalkulator
lokal tidak dilabeli keluaran cfo.ai, dan tidak otomatis masuk paket publik.
Metode kas dibagi burn dijelaskan dalam
[referensi cash runway CFI](https://corporatefinanceinstitute.com/resources/valuation/cash-runway-explained/).

## Memasukkan hasil ke cfo.ai

Setelah input terisi, kirim [prompt tindak lanjut](CFO-RUNWAY-FOLLOWUP.md) di sesi
cfo.ai yang sama. Lampirkan `EVIDENCE.md`, `evidence/cfo-validation.json`, dan
hasil runway privat yang sudah dipilih untuk dibagikan. Mintalah jadwal kas
bertanggal dan ekspor revisi PDF/XLSX asli dari layanan tersebut. Jika actual
data belum ada, gunakan budget yang ditentukan user dengan label **simulasi**.

Saldo pribadi hanya dipublikasikan jika user memilihnya. Simpan revisi dengan
nama berbeda; jangan menimpa ekspor awal atau menganggap supplement lokal
sebagai rencana yang sudah diperbarui oleh cfo.ai. Defisit operasional kumulatif
US$59.175 pada ekspor awal tetap bukan saldo kas atau runway.
