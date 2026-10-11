import type { Run } from './checkout'
import type { Language } from './language'

export const sandboxGoal =
  'Buy one tote using demo@example.test. Verify the checkout total and confirm exactly one order. Recover safely from any timeout.'
export const marketplaceGoal =
  'Prepare checkout for one item from this product URL. Use the selected variant, then stop at the final order review so I can confirm it manually.'
export const sandboxGoals: Record<Language, string> = {
  en: sandboxGoal,
  id: 'Beli satu tote menggunakan demo@example.test. Verifikasi total checkout dan pastikan tepat satu pesanan tersimpan. Pulihkan alur dengan aman jika terjadi timeout.',
}
export const marketplaceGoals: Record<Language, string> = {
  en: marketplaceGoal,
  id: 'Siapkan checkout satu barang dari URL produk ini. Gunakan varian yang dipilih, lalu berhenti di halaman review pesanan agar saya dapat mengonfirmasinya secara manual.',
}
export const isMarket = (r: Run) => r.environment === 'shopee' || r.environment === 'tokopedia'
