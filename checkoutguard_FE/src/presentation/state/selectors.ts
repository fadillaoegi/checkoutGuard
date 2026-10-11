import { createSelector } from '@reduxjs/toolkit'
import type { RootState } from './store'
import type { Run } from '../../domain/checkout'

export const selectCheckout = (s: RootState) => s.checkout
export const selectRun = createSelector([selectCheckout], (s) =>
  s.runs.find((r) => r.id === s.selected),
)
export const selectActive = createSelector([selectCheckout], (s) =>
  s.runs.some((r) => r.status === 'running'),
)
export const selectStep = createSelector(
  [selectCheckout, selectRun],
  (s, r) => r?.steps.find((step) => step.number === s.stepNumber) || r?.steps.at(-1),
)
export const selectBaseline = createSelector([selectCheckout, selectRun], (s, r) =>
  r?.parent_id ? s.runs.find((run) => run.id === r.parent_id) : undefined,
)
export const defects = (r: Run) =>
  r.findings.filter((f) => f.verified && f.kind !== 'timeout_recovered').length
export const selectMetrics = createSelector([selectCheckout], (s) => ({
  verified: s.runs.filter((r) => r.status === 'passed').length,
  totalTools: s.runs.reduce(
    (n, r) =>
      n + r.steps.filter((step) => step.action.tool !== 'model_error' && step.number > 0).length,
    0,
  ),
}))
export const selectMarketBrowserOpen = createSelector(
  [selectCheckout],
  (s) => !!s.marketSession?.browser_open && s.marketSession.environment === s.environment,
)
export const selectMarketPageState = createSelector(
  [selectCheckout, selectMarketBrowserOpen],
  (s, open) => (open ? s.marketSession?.page_state || 'unknown' : 'unknown'),
)
