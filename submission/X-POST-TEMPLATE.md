# Template submission — CheckoutGuard (#neatHack)

Sumber aturan: email resmi neatlogs "neatHack submissions open at 7 PM IST today"
(diterima 12 Oktober 2026). Submission dibuka **12 Oktober 19:00–23:55 IST =
20:30–01:25 WIB**. Pantau Discord untuk perubahan cutoff.

## Dua langkah submit — KEDUANYA WAJIB

1. **Quote-tweet posting peluncuran resmi @neatlogs** di X (link dibagikan pukul 20:30 WIB),
   dan quote-tweet itu **wajib menandai @neatlogs, @EntireHQ, dan @arithecfo**.
   Sistem mereka mendeteksi submission dari quote-tweet dan ketiga tag ini.
2. **Isi formulir submission** (link juga dibagikan pukul 20:30 WIB).

Hanya mengisi formulir, atau lupa salah satu tag, artinya **submission tidak terdeteksi**.
neatlogs akan berinteraksi dengan 50 quote-tweet pertama, jadi kirim segera setelah dibuka.

## Langkah 0 — siapkan data dulu

| Data | Isi Anda | Contoh |
| --- | --- | --- |
| Nama anggota tim | `[NAMA]` (maks. ±40 huruf) | Budi Santoso |
| Link video demo (bisa dibuka tanpa login) | `[LINK VIDEO]` | https://youtu.be/abc123 |
| Link posting peluncuran resmi neatlogs | dari Discord/X pukul 20:30 WIB | — |
| Link formulir submission | dari Discord/X pukul 20:30 WIB | — |

Unggah `submission/video/CheckoutGuard-final.mp4` ke YouTube (**Unlisted** atau
**Public**), lalu salin link-nya. Unggah juga `CheckoutGuard-final.en.srt` sebagai subtitle.

## Posting 1 — progres (building in public) · kirim SEKARANG

Langkah: buka x.com → **Post** → tempel teks → lampirkan
`submission/video/social/post1-intro-explainer.mp4` → **Post**. Tidak ada yang perlu diganti.

```
Building CheckoutGuard for #neatHack 🛒

An AI agent that tests checkouts and keeps the evidence:
🐛 found a wrong total in a test store
⏱️ payment timeout: no double order
🛠️ built with @neatlogs @EntireHQ @arithecfo

github.com/fadillaoegi/checkoutGuard
```

Panjang 242/280 menurut hitungan X (link = 23, emoji = 2). Aman untuk akun gratis.
Link cepat: [buka X dengan teks terisi](https://x.com/intent/tweet?text=Building%20CheckoutGuard%20for%20%23neatHack%20%F0%9F%9B%92%0A%0AAn%20AI%20agent%20that%20tests%20checkouts%20and%20keeps%20the%20evidence%3A%0A%F0%9F%90%9B%20found%20a%20wrong%20total%20in%20a%20test%20store%0A%E2%8F%B1%EF%B8%8F%20payment%20timeout%3A%20no%20double%20order%0A%F0%9F%9B%A0%EF%B8%8F%20built%20with%20%40neatlogs%20%40EntireHQ%20%40arithecfo%0A%0Agithub.com/fadillaoegi/checkoutGuard)

## Posting 2 — SUBMISSION (quote-tweet) · mulai 20:30 WIB

Langkah:
1. Buka link posting peluncuran resmi @neatlogs.
2. Klik **Repost** (panah melingkar) → **Quote**.
3. Tempel teks di bawah, ganti `[NAMA]` dan `[LINK VIDEO]`.
   **Jangan hapus baris "Built with"**: ketiga tag itu wajib.
4. Pastikan ketiga tag berubah biru (dikenali X), lalu **Post**.
5. Salin link quote-tweet Anda.

```
🛒 CheckoutGuard: an AI agent that tests your checkout and keeps the evidence.

Team: FLdevBoo
Name: [NAMA]
Project: CheckoutGuard
Repo: github.com/fadillaoegi/checkoutGuard
Demo: [LINK VIDEO]

Built with @neatlogs @EntireHQ @arithecfo
#neatHack
```

Dengan nama ±25 huruf panjangnya ±262/280. Jika nama lebih panjang dan X menolak,
hapus kata "your" di baris pertama, tetapi **jangan menghapus tag**.

## Langkah 3 — formulir submission (wajib)

Kolom formulir belum diketahui. Draf jawaban umum (bahasa Inggris), siap disalin:

| Kolom kemungkinan | Jawaban |
| --- | --- |
| Team name | FLdevBoo |
| Project name | CheckoutGuard |
| One-line description | An AI agent that tests e-commerce checkouts, verifies the total and saved orders, and keeps the evidence. |
| GitHub repo | https://github.com/fadillaoegi/checkoutGuard |
| Demo video | `[LINK VIDEO]` |
| Quote-tweet link | `[LINK QUOTE-TWEET]` |

**Project description (pendek):**

```
CheckoutGuard is a goal-driven browser agent for checkout QA. You describe a purchase in plain English; the agent drives Chrome step by step, checks the total with integer arithmetic, confirms exactly one saved order in PostgreSQL, and stores a screenshot for every step. In our synthetic sandbox it caught a wrong total (14,400 vs 12,400 IDR), passed the retest on the corrected store, and recovered from a payment timeout without creating a duplicate order.
```

**How did you use neatlogs, Entire and cfo.ai?**

```
neatlogs: every run is traced (workflow, model calls and browser tool calls); we read the hosted traces back to check the hierarchy and that token totals match our local usage records, and used them to debug model calls.
Entire: the coding-agent sessions are checkpointed; checkpoint 01M4NKBNCP60MN7M5BXS143W16 is linked to public commit a134e4a (35 files).
cfo.ai: we built the business plan with it: pricing ($29/$99/$299 per month), modeled delivery cost (~$0.019 per run) and a break-even of ~96 subscriptions. All figures are assumptions; the product is pre-revenue.
```

**Limitations (jika ditanya):**

```
Demo runs use a synthetic store with simulated payments. Shopee/Tokopedia adapters are experimental and have not completed a logged-in checkout. The video reviews saved runs because the model provider was timing out during recording.
```

## Cek sebelum menekan Post (posting 2)

- [ ] Ini **Quote** dari posting peluncuran resmi @neatlogs, bukan posting biasa.
- [ ] Ada **@neatlogs**, **@EntireHQ**, dan **@arithecfo** (ketiganya biru).
- [ ] Ada **#neatHack**, link repo, dan link video.
- [ ] Link video bisa dibuka di jendela incognito.
- [ ] Tidak ada `[` atau `]` yang tertinggal.
- [ ] Setelah posting: **isi formulir** submission.

## Setelah selesai

Catat link quote-tweet dan bukti formulir terkirim di
[BUILDING-IN-PUBLIC.md](BUILDING-IN-PUBLIC.md) dan [SUBMISSION-CHECKLIST.md](SUBMISSION-CHECKLIST.md).
