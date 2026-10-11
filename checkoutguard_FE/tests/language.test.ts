import assert from 'node:assert/strict'
import test from 'node:test'
import { CheckoutUseCases } from '../src/application/checkoutUseCases'
import type { CheckoutRepository } from '../src/application/CheckoutRepository'
import { LanguagePreferences } from '../src/application/languagePreferences'
import { marketplaceGoals, sandboxGoals } from '../src/domain/checkoutRules'
import type { Language } from '../src/domain/language'
import {
  BrowserLanguagePreferences,
  languageStorageKey,
} from '../src/infrastructure/browserLanguagePreferences'
import { localizeMessage, translate } from '../src/presentation/i18n/messages'
import { cost, money, runLabel } from '../src/presentation/dashboard/format'
import { actions } from '../src/presentation/state/checkoutSlice'
import { createAppStore } from '../src/presentation/state/store'

const unexpected = () => {
  throw new Error('Changing language must not invoke checkout operations')
}
const repository: CheckoutRepository = {
  snapshot: unexpected,
  start: unexpected,
  retest: unexpected,
  stop: unexpected,
  connect: unexpected,
  marketSession: unexpected,
}
const checkout = new CheckoutUseCases(repository)

test('English is the default; a saved choice survives creating another store', () => {
  let saved: unknown = null
  let writes = 0
  const preferences = new LanguagePreferences({
    read: () => saved,
    write: (language) => {
      saved = language
      writes++
    },
  })
  const store = createAppStore({ checkout }, preferences)
  assert.equal(store.getState().checkout.language, 'en')
  store.dispatch(actions.setLanguage('id'))
  assert.equal(saved, 'id')
  assert.equal(store.getState().checkout.goal, sandboxGoals.id)
  store.dispatch(actions.setScenario('timeout'))
  assert.equal(writes, 1, 'Unrelated actions must not write the preference')
  const reloaded = createAppStore({ checkout }, preferences)
  assert.equal(reloaded.getState().checkout.language, 'id')
  assert.equal(reloaded.getState().checkout.goal, sandboxGoals.id)
  store.dispatch(actions.setLanguage('en'))
  assert.equal(saved, 'en')
})

test('language changes preserve user-authored goals, URLs and evidence selection', () => {
  const store = createAppStore({ checkout })
  store.dispatch(actions.selectEnvironment('shopee'))
  store.dispatch(actions.setGoal('Siapkan dua barang berwarna biru, jangan pesan dulu.'))
  store.dispatch(actions.setTarget('https://shopee.co.id/product/123/456'))
  store.dispatch(actions.selectStep(3))
  store.dispatch(actions.setDialog('/evidence/actual-run/3.png'))
  const before = store.getState().checkout
  store.dispatch(actions.setLanguage('id'))
  const after = store.getState().checkout
  assert.equal(after.goal, before.goal)
  assert.equal(after.target, before.target)
  assert.equal(after.stepNumber, 3)
  assert.equal(after.dialog, before.dialog)
  assert.deepEqual(after.runs, before.runs)
  store.dispatch(actions.selectEnvironment('tokopedia'))
  assert.equal(store.getState().checkout.goal, marketplaceGoals.id)
})

test('missing, invalid and inaccessible browser storage fall back to English', () => {
  for (const value of [null, 'fr', '{broken', '', 'ID']) {
    const preferences = new LanguagePreferences({ read: () => value, write: () => {} })
    assert.equal(preferences.current(), 'en')
  }
  const blocked = new BrowserLanguagePreferences(() => {
    throw new Error('Storage blocked')
  })
  assert.equal(new LanguagePreferences(blocked).current(), 'en')
  assert.doesNotThrow(() => blocked.write('id'))
  const entries = new Map<string, string>()
  const storage = {
    getItem: (key: string) => entries.get(key) ?? null,
    setItem: (key: string, value: string) => {
      entries.set(key, value)
    },
  } as Storage
  const allowed = new BrowserLanguagePreferences(() => storage)
  allowed.write('id')
  assert.equal(entries.get(languageStorageKey), 'id')
  assert.equal(allowed.read(), 'id')
})

test('display translation preserves unknown evidence and supports locale-specific amounts', () => {
  assert.equal(translate('en', 'Run checkout'), 'Run checkout')
  assert.equal(translate('en', '{count} order(s)', { count: 1 }), '1 order')
  assert.equal(translate('en', '{count} order(s)', { count: 2 }), '2 orders')
  assert.equal(
    translate('id', 'Open {marketplace} browser', { marketplace: 'Shopee' }),
    'Buka browser Shopee',
  )
  assert.equal(money(12400, 'id'), '12.400 IDR')
  assert.equal(money(12400, 'en'), '12,400 IDR')
  const raw = 'Product Blue Tote — actual provider response'
  for (const language of ['en', 'id'] as Language[])
    assert.equal(localizeMessage(language, raw), raw)
  assert.equal(
    localizeMessage('id', 'model response timed out'),
    'Respons model melewati batas waktu',
  )
  assert.match(
    localizeMessage('id', 'Expected 12400 IDR; checkout displays 14400 IDR'),
    /12400.*14400/,
  )
  assert.equal(
    runLabel({ environment: 'tokopedia' } as Parameters<typeof runLabel>[0], 'id'),
    'Persiapan checkout Tokopedia',
  )
})

test('unknown and partial model costs never look like a complete free request', () => {
  const usage = { input_tokens: 100, output_tokens: 10, cost_usd: null }
  assert.equal(cost(usage, 'en'), 'Cost unavailable')
  assert.equal(cost(usage, 'id'), 'Biaya belum tersedia')
  assert.match(
    cost({ ...usage, cost_usd: 0.00003, cost_incomplete: true }, 'en'),
    /^At least .*\(partial\)$/,
  )
  assert.match(
    cost({ ...usage, cost_usd: 0.00003, cost_incomplete: true }, 'id'),
    /^Minimal .*\(belum lengkap\)$/,
  )
  assert.doesNotMatch(cost({ ...usage, cost_usd: 0.00003 }, 'en'), /partial/)
})
