import type { Usage } from '../../domain/checkout'
import { languageLocale, type Language } from '../../domain/language'
import { translate, type TranslationKey } from '../i18n/messages'

export const toolLabels: Record<string, TranslationKey> = {
  observe: 'Observe browser',
  click: 'Click page control',
  fill: 'Fill checkout field',
  verify_total: 'Verify checkout total',
  order_status: 'Check saved orders',
  finish: 'Confirm result',
  model_error: 'Model request issue',
}
export const money = (n: number, language: Language = 'en') =>
  new Intl.NumberFormat(languageLocale(language)).format(n) + ' IDR'
export const cost = (u?: Usage, language: Language = 'en') => {
  if (u?.cost_usd == null) return translate(language, 'Cost unavailable')
  const amount = new Intl.NumberFormat(languageLocale(language), {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 5,
    maximumFractionDigits: 5,
  }).format(u.cost_usd)
  return u.cost_incomplete ? translate(language, 'At least {amount} (partial)', { amount }) : amount
}

import type { Run } from '../../domain/checkout'
import { isMarket } from '../../domain/checkoutRules'
export const environmentNames = { sandbox: 'Sandbox', shopee: 'Shopee', tokopedia: 'Tokopedia' }
export const runLabel = (r: Run, language: Language = 'en') =>
  isMarket(r)
    ? translate(language, '{marketplace} checkout preparation', {
        marketplace: environmentNames[r.environment!],
      })
    : labels[r.scenario]
      ? translate(language, labels[r.scenario])
      : r.scenario
export const labels: Record<string, TranslationKey> = {
  fixed: 'Correct checkout',
  total_bug: 'Incorrect total',
  timeout: 'Payment timeout',
}
