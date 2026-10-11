import assert from 'node:assert/strict'
import test from 'node:test'
import { CheckoutUseCases } from '../src/application/checkoutUseCases'
import type { CheckoutRepository } from '../src/application/CheckoutRepository'
import type { DashboardSnapshot, Run, StartRunInput } from '../src/domain/checkout'
import { HttpCheckoutRepository } from '../src/infrastructure/httpCheckoutRepository'
import {
  actions,
  connectMarketplace,
  launchRun,
  refreshDashboard,
  stopRun,
} from '../src/presentation/state/checkoutSlice'
import { createAppStore } from '../src/presentation/state/store'
import { selectMarketPageState } from '../src/presentation/state/selectors'

const snapshot: DashboardSnapshot = {
  health: { model_configured: true, model: 'test', sandbox_url: 'http://127.0.0.1:8080/sandbox' },
  runs: [],
  marketSession: { environment: '', browser_open: false },
}
const run = (id: string, status = 'running'): Run => ({
  id,
  status,
  environment: 'sandbox',
  goal: 'Verify',
  url: snapshot.health.sandbox_url,
  scenario: 'fixed',
  variant: 'standard',
  model: 'test',
  started_at: '2026-10-11T00:00:00Z',
  steps: [],
  findings: [],
  usage: { input_tokens: 0, output_tokens: 0, cost_usd: null },
  order_count: 0,
  total_verified: false,
  expected_total_minor: 0,
  displayed_total_minor: 0,
})
function repository(overrides: Partial<CheckoutRepository> = {}): CheckoutRepository {
  return {
    snapshot: async () => snapshot,
    start: async () => run('new'),
    retest: async () => run('retest'),
    stop: async () => {},
    connect: async () => ({ navigation: 'loaded' }),
    marketSession: async () => snapshot.marketSession,
    ...overrides,
  }
}
const makeStore = (repo: CheckoutRepository) =>
  createAppStore({ checkout: new CheckoutUseCases(repo) })
const deferred = <T>() => {
  let resolve!: (value: T) => void
  const promise = new Promise<T>((done) => {
    resolve = done
  })
  return { promise, resolve }
}

test('switching marketplaces clears product URL and selects sandbox without launching it', async () => {
  let starts = 0
  const store = makeStore(
    repository({
      start: async () => {
        starts++
        return run('unexpected')
      },
    }),
  )
  await store.dispatch(refreshDashboard())
  store.dispatch(actions.selectEnvironment('shopee'))
  store.dispatch(actions.setTarget('https://shopee.co.id/product/123/456'))
  store.dispatch(actions.selectEnvironment('tokopedia'))
  assert.equal(store.getState().checkout.target, '')
  store.dispatch(actions.selectEnvironment('sandbox'))
  assert.equal(store.getState().checkout.target, snapshot.health.sandbox_url)
  assert.equal(starts, 0)
})

test('a poll started before launch cannot erase the newly created run', async () => {
  const pending = deferred<DashboardSnapshot>()
  const store = makeStore(repository({ snapshot: () => pending.promise }))
  const poll = store.dispatch(refreshDashboard())
  await store.dispatch(launchRun(undefined))
  pending.resolve(snapshot)
  await poll
  assert.equal(store.getState().checkout.selected, 'new')
  assert.deepEqual(
    store.getState().checkout.runs.map((r) => r.id),
    ['new'],
  )
  assert.equal(store.getState().checkout.refreshing, false)
})

test('duplicate launch is suppressed and persisted failure remains in history', async () => {
  const pending = deferred<Run>()
  let starts = 0
  let received: StartRunInput | undefined
  const store = makeStore(
    repository({
      start: (input) => {
        starts++
        received = input
        return pending.promise
      },
    }),
  )
  store.dispatch(actions.selectEnvironment('shopee'))
  store.dispatch(actions.setTarget('https://shopee.co.id/product/123/456'))
  const first = store.dispatch(launchRun(undefined))
  await store.dispatch(launchRun(undefined))
  pending.resolve({
    ...run('failed-market', 'failed'),
    environment: 'shopee',
    reason: 'Access challenge',
  })
  await first
  assert.equal(starts, 1)
  assert.equal(received?.environment, 'shopee')
  assert.equal(received?.url, 'https://shopee.co.id/product/123/456')
  store.dispatch(actions.selectEnvironment('sandbox'))
  assert.equal(store.getState().checkout.runs[0].reason, 'Access challenge')
  assert.equal(store.getState().checkout.runs[0].environment, 'shopee')
})

test('connection warnings and errors release the pending state; stop targets the selected run', async () => {
  let stopped = ''
  const store = makeStore(
    repository({
      connect: async () => ({ navigation: 'page did not finish loading' }),
      marketSession: async () => ({ environment: 'shopee', browser_open: true }),
      stop: async (id) => {
        stopped = id
      },
      start: async () => {
        throw new Error('Model unavailable')
      },
    }),
  )
  store.dispatch(actions.selectEnvironment('shopee'))
  await store.dispatch(connectMarketplace())
  assert.equal(store.getState().checkout.marketSession?.browser_open, true)
  assert.match(store.getState().checkout.browserNotice, /did not finish loading/)
  assert.equal(store.getState().checkout.connecting, false)
  await store.dispatch(launchRun(undefined))
  assert.equal(store.getState().checkout.pending, false)
  assert.equal(store.getState().checkout.error, 'Model unavailable')
  await store.dispatch(stopRun('actual-run'))
  assert.equal(stopped, 'actual-run')
})

test('selecting another run resets the step; dialogs are serializable Redux state', async () => {
  const store = makeStore(
    repository({
      snapshot: async () => ({ ...snapshot, runs: [run('a', 'passed'), run('b', 'failed')] }),
    }),
  )
  await store.dispatch(refreshDashboard())
  store.dispatch(actions.selectStep(3))
  store.dispatch(actions.selectRun('b'))
  assert.equal(store.getState().checkout.stepNumber, null)
  store.dispatch(actions.setDialog('/evidence/b/0.png'))
  assert.equal(store.getState().checkout.dialog, '/evidence/b/0.png')
  store.dispatch(actions.setDialog(null))
  assert.equal(store.getState().checkout.dialog, null)
})

test('HTTP adapter encodes IDs, sends JSON and propagates backend errors', async () => {
  const requests: { path: string; init?: RequestInit }[] = []
  const request = (async (path: string, init?: RequestInit) => {
    requests.push({ path, init })
    if (path.endsWith('/stop'))
      return new Response(JSON.stringify({ error: 'run is not active' }), { status: 409 })
    return new Response(JSON.stringify(run('retest')), { status: 202 })
  }) as typeof fetch
  const repo = new HttpCheckoutRepository(request)
  await repo.retest('a/b')
  assert.equal(requests[0].path, '/api/runs/a%2Fb/retest')
  assert.equal(requests[0].init?.body, '{}')
  assert.equal(requests[0].init?.method, 'POST')
  await assert.rejects(repo.stop('x'), /run is not active/)
})

test('poll requests do not overlap, and abort clears refresh without a connection error', async () => {
  let reads = 0
  const repo = repository({
    snapshot: (signal) =>
      new Promise((_, reject) => {
        reads++
        signal?.addEventListener('abort', () => reject(new Error('aborted')), { once: true })
      }),
  })
  const store = makeStore(repo)
  const first = store.dispatch(refreshDashboard())
  await store.dispatch(refreshDashboard())
  assert.equal(reads, 1)
  first.abort()
  await first
  assert.equal(store.getState().checkout.refreshing, false)
  assert.equal(store.getState().checkout.connectionError, '')
})

test('injected request runs as a function instead of receiving the repository as this', async () => {
  const request = function (this: unknown) {
    assert.equal(this, undefined)
    return Promise.resolve(new Response(JSON.stringify(snapshot.marketSession)))
  } as typeof fetch
  const repo = new HttpCheckoutRepository(request)
  assert.deepEqual(await repo.marketSession(), snapshot.marketSession)
})

test('blocked or login pages cannot launch until a manual browser check updates access', async () => {
  for (const pageState of ['blocked', 'needs_login'] as const) {
    let starts = 0
    let checked = false
    const repo = repository({
      snapshot: async () => ({
        ...snapshot,
        marketSession: { environment: 'shopee', browser_open: true, page_state: pageState },
      }),
      connect: async () => {
        checked = true
        return { page_state: 'browsing' }
      },
      marketSession: async () => ({
        environment: 'shopee',
        browser_open: true,
        page_state: checked ? 'browsing' : pageState,
      }),
      start: async () => {
        starts++
        return { ...run('prepared'), environment: 'shopee' }
      },
    })
    const store = makeStore(repo)
    await store.dispatch(refreshDashboard())
    store.dispatch(actions.selectEnvironment('shopee'))
    store.dispatch(actions.setTarget('https://shopee.co.id/product/123/456'))
    await store.dispatch(launchRun(undefined))
    assert.equal(starts, 0)
    assert.equal(selectMarketPageState(store.getState()), pageState)
    await store.dispatch(connectMarketplace())
    assert.equal(selectMarketPageState(store.getState()), 'browsing')
    await store.dispatch(launchRun(undefined))
    assert.equal(starts, 1)
  }
})

test('a failed page check cannot enable checkout even if the last page was browsing', async () => {
  let starts = 0
  const store = makeStore(
    repository({
      snapshot: async () => ({
        ...snapshot,
        marketSession: {
          environment: 'shopee',
          browser_open: true,
          page_state: 'browsing',
          page_check: 'unavailable',
        },
      }),
      start: async () => {
        starts++
        return run('unexpected')
      },
    }),
  )
  await store.dispatch(refreshDashboard())
  store.dispatch(actions.selectEnvironment('shopee'))
  await store.dispatch(launchRun(undefined))
  assert.equal(starts, 0)
})
