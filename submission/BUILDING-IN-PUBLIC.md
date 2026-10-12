# Building in public — CheckoutGuard (#neatHack)

Kriteria penilaian **Building in public = 15%**. Juri meminta "links to public posts
tagged #neatHack". Caption di bawah siap disalin ke X (tiap posting ≤ 280 karakter menurut hitungan X:
link = 23, emoji = 2).
Semua klaim sesuai bukti di repo; jangan menambah angka yang tidak ada di
[EVIDENCE.md](EVIDENCE.md).

**Template paling sederhana (2 posting + formulir, tinggal isi): [X-POST-TEMPLATE.md](X-POST-TEMPLATE.md).**

**Aturan baru dari email neatlogs (12 Oktober):** quote-tweet submission **wajib menandai
@neatlogs, @EntireHQ dan @arithecfo**, dan **formulir submission juga wajib**. Tanpa salah
satunya, submission tidak terdeteksi.

Status: **belum dipublikasikan** (12 Oktober 2026). Setelah memposting, tempel
link-nya di tabel "Link yang sudah dipublikasikan" di bawah dan di
[SUBMISSION-CHECKLIST.md](SUBMISSION-CHECKLIST.md).

## Versi ringkas (direkomendasikan jika waktu terbatas): 2 posting

Yang **wajib** hanya posting submission (posting 4). Building in public dinilai dari
posting progres, jadi minimal satu posting progres sebelum submit. Posting 1–3 di bawah
bisa diganti satu posting gabungan ini (239/280 karakter menurut hitungan X):

> Building CheckoutGuard for #neatHack 🛒
>
> An AI agent that tests checkouts and keeps the evidence:
> 🐛 found a wrong total in a test store
> ⏱️ payment timeout: no double order
> 🛠️ built with neatlogs, Entire and cfo.ai
>
> github.com/fadillaoegi/checkoutGuard

Lampiran: `submission/video/social/post1-intro-explainer.mp4`.
Link cepat: [buka X dengan teks terisi](https://x.com/intent/tweet?text=Building%20CheckoutGuard%20for%20%23neatHack%20%F0%9F%9B%92%0A%0AAn%20AI%20agent%20that%20tests%20checkouts%20and%20keeps%20the%20evidence%3A%0A%F0%9F%90%9B%20found%20a%20wrong%20total%20in%20a%20test%20store%0A%E2%8F%B1%EF%B8%8F%20payment%20timeout%3A%20no%20double%20order%0A%F0%9F%9B%A0%EF%B8%8F%20built%20with%20neatlogs%2C%20Entire%20and%20cfo.ai%0A%0Agithub.com/fadillaoegi/checkoutGuard)

Lalu malam ini kirim posting 4 (quote post resmi). Total 2 posting.

## Cara posting

1. Posting dari akun X tim/anggota tim, **publik** (bukan akun terkunci).
2. Setiap posting memakai **#neatHack** dan link repo.
3. Tag **@neatlogs @EntireHQ @arithecfo** (handle resmi dari email neatlogs, 12 Oktober).
4. Lampirkan gambar atau klip pendek (lihat kolom "Lampiran").
5. Beri jeda antar-posting (misalnya 20–30 menit); jangan semuanya sekaligus.
6. Posting 1–3 bisa sekarang. Posting 4 (quote post resmi) adalah **cara submit**,
   dilakukan saat submission dibuka **12 Oktober 20:30 WIB** (tutup 13 Oktober 01:25 WIB).

## Link cepat (teks sudah terisi)

Klik, periksa teks, tambahkan lampiran, lalu tekan **Post**. Posting 4 tidak memakai link
cepat karena harus berupa quote dari posting resmi neatlogs.

- Posting 1: [buka X dengan teks terisi](https://x.com/intent/tweet?text=Ever%20shipped%20a%20checkout%20that%20looked%20fine%20but%20charged%20the%20wrong%20amount%3F%20%F0%9F%9B%92%0A%0AWe%27re%20building%20CheckoutGuard%20for%20%23neatHack%3A%20tell%20it%20the%20purchase%20in%20plain%20English%2C%20and%20an%20AI%20agent%20drives%20Chrome%2C%20checks%20the%20math%20and%20saves%20the%20evidence.%0A%0Agithub.com/fadillaoegi/checkoutGuard)
- Posting 2: [buka X dengan teks terisi](https://x.com/intent/tweet?text=Our%20agent%20caught%20a%20pricing%20bug%20in%20our%20test%20store%20%F0%9F%90%9B%0A%0A12%2C900%20%E2%88%92%202%2C000%20%2B%201%2C500%20%3D%2012%2C400%20IDR%0AThe%20checkout%20showed%3A%2014%2C400%20%E2%9D%8C%0A%0ARetest%20on%20the%20fixed%20store%3A%2012%2C400%20%E2%9C%85%2C%20exactly%20one%20order%2C%20every%20step%20with%20a%20screenshot.%0A%0A%23neatHack)
- Posting 3: [buka X dengan teks terisi](https://x.com/intent/tweet?text=What%20if%20the%20payment%20times%20out%20mid-checkout%3F%20%E2%8F%B1%EF%B8%8F%0A%0ACheckoutGuard%20doesn%27t%20retry%20blindly.%20It%20checks%20the%20saved%20orders%20first%2C%20finds%20exactly%20one%20and%20stops.%20No%20double%20order.%0A%0A%F0%9F%94%8D%20Traced%20with%20neatlogs%0A%F0%9F%A7%B1%20Built%20with%20Entire%20checkpoints%0A%F0%9F%93%8A%20Planned%20with%20cfo.ai%0A%0A%23neatHack)

File lampiran ada di `submission/video/social/` (lokal, di-ignore Git).

## Posting 1 — perkenalan proyek

> Ever shipped a checkout that looked fine but charged the wrong amount? 🛒
>
> We're building CheckoutGuard for #neatHack: tell it the purchase in plain English, and an AI agent drives Chrome, checks the math and saves the evidence.
>
> github.com/fadillaoegi/checkoutGuard

Lampiran: klip 25 detik `submission/video/social/post1-intro-explainer.mp4` (intro + animasi cara kerja, dengan musik).

## Posting 2 — hasil: bug harga dan retest

> Our agent caught a pricing bug in our test store 🐛
>
> 12,900 − 2,000 + 1,500 = 12,400 IDR
> The checkout showed: 14,400 ❌
>
> Retest on the fixed store: 12,400 ✅, exactly one order, every step with a screenshot.
>
> #neatHack

Lampiran: `submission/video/social/post2-before-after.png`.

## Posting 3 — kegagalan dan partner tools

> What if the payment times out mid-checkout? ⏱️
>
> CheckoutGuard doesn't retry blindly. It checks the saved orders first, finds exactly one and stops. No double order.
>
> 🔍 Traced with neatlogs
> 🧱 Built with Entire checkpoints
> 📊 Planned with cfo.ai
>
> #neatHack

Lampiran: `submission/video/social/post3-payment-timeout.png` dan `post3-entire-checkpoint.png`.

## Posting 4 — submission (quote post resmi neatlogs)

Quote posting submission resmi dari neatlogs, lalu isi:

> 🛒 CheckoutGuard, an AI agent that tests your checkout and keeps the evidence.
>
> Team: FLdevBoo
> Name: [NAMA ANGGOTA TIM]
> Project: CheckoutGuard
> Repo: github.com/fadillaoegi/checkoutGuard
> Demo: [LINK VIDEO]
>
> #neatHack

Ganti `[NAMA ANGGOTA TIM]` dan `[LINK VIDEO]` (nama sampai ±60 karakter masih muat). Halaman resmi meminta "team, name, project,
public repo and demo"; uji link video di jendela incognito (tanpa login).
neatlogs dan cfo.ai akan berinteraksi dengan 50 quote pertama, jadi posting segera
setelah submission dibuka.

## Opsional — LinkedIn

Gabungkan posting 1–3 menjadi satu posting LinkedIn dengan video final terlampir
dan #neatHack. Ini tambahan, bukan pengganti posting X.

## Jangan diklaim

- Checkout Shopee/Tokopedia berhasil (belum terbukti).
- Pelanggan berbayar atau pendapatan (rencana cfo.ai seluruhnya asumsi).
- Provider model sudah stabil, atau demo direkam live (video meninjau run tersimpan).

## Link yang sudah dipublikasikan

| Posting | Link | Tanggal |
| --- | --- | --- |
| 1 | — | — |
| 2 | — | — |
| 3 | — | — |
| 4 (submission) | — | — |
