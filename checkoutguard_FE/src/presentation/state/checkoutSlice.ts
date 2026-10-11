import { createAsyncThunk, createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { CheckoutUseCases } from '../../application/checkoutUseCases'
import type { Environment, Health, MarketSession, Run } from '../../domain/checkout'
import { marketplaceGoals, sandboxGoals } from '../../domain/checkoutRules'
import { defaultLanguage, type Language } from '../../domain/language'

export type CheckoutState = {
  language: Language
  health: Health | null
  runs: Run[]
  selected: string
  stepNumber: number | null
  goal: string
  target: string
  environment: Environment
  marketSession: MarketSession | null
  connecting: boolean
  browserNotice: string
  scenario: string
  variant: string
  error: string
  connectionError: string
  pending: boolean
  dialog: string | null
  refreshing: boolean
  revision: number
}
type State = { checkout: CheckoutState }
export type Dependencies = { checkout: CheckoutUseCases }
const thunk = createAsyncThunk.withTypes<{
  state: State
  extra: Dependencies
  rejectValue: string
}>()
const message = (e: unknown) => (e instanceof Error ? e.message : 'Request failed')

export const refreshDashboard = thunk(
  'checkout/refresh',
  async (_, { extra, getState, signal, rejectWithValue }) => {
    const revision = getState().checkout.revision
    try {
      return { ...(await extra.checkout.refresh(signal)), revision }
    } catch (e) {
      return rejectWithValue(message(e))
    }
  },
  { condition: (_, { getState }) => !getState().checkout.refreshing },
)

export const launchRun = thunk(
  'checkout/launch',
  async (retest: string | undefined, { extra, getState, signal, rejectWithValue }) => {
    const { goal, target, environment, scenario, variant } = getState().checkout
    try {
      return retest
        ? await extra.checkout.retest(retest, signal)
        : await extra.checkout.start({ goal, url: target, environment, scenario, variant }, signal)
    } catch (e) {
      return rejectWithValue(message(e))
    }
  },
  {
    condition: (_, { getState }) => {
      const s = getState().checkout
      const pageState =
        s.marketSession?.environment === s.environment ? s.marketSession.page_state : undefined
      return (
        !s.pending &&
        !s.connecting &&
        !s.runs.some((r) => r.status === 'running') &&
        (s.environment === 'sandbox' ||
          (pageState !== 'blocked' &&
            pageState !== 'needs_login' &&
            !(
              s.marketSession?.environment === s.environment &&
              s.marketSession.page_check === 'unavailable'
            )))
      )
    },
  },
)

export const connectMarketplace = thunk(
  'checkout/connect',
  async (_, { extra, getState, signal, rejectWithValue }) => {
    try {
      return await extra.checkout.connect(getState().checkout.environment, signal)
    } catch (e) {
      return rejectWithValue(message(e))
    }
  },
  {
    condition: (_, { getState }) => {
      const s = getState().checkout
      return (
        s.environment !== 'sandbox' &&
        !s.connecting &&
        !s.pending &&
        !s.runs.some((r) => r.status === 'running')
      )
    },
  },
)

export const stopRun = thunk(
  'checkout/stop',
  async (id: string, { extra, signal, rejectWithValue }) => {
    try {
      await extra.checkout.stop(id, signal)
    } catch (e) {
      return rejectWithValue(message(e))
    }
  },
)

const initialState: CheckoutState = {
  language: defaultLanguage,
  health: null,
  runs: [],
  selected: '',
  stepNumber: null,
  goal: sandboxGoals.en,
  target: 'http://127.0.0.1:8080/sandbox',
  environment: 'sandbox',
  marketSession: null,
  connecting: false,
  browserNotice: '',
  scenario: 'total_bug',
  variant: 'standard',
  error: '',
  connectionError: '',
  pending: false,
  dialog: null,
  refreshing: false,
  revision: 0,
}
export const checkoutSlice = createSlice({
  name: 'checkout',
  initialState,
  reducers: {
    setLanguage: (s, a: PayloadAction<Language>) => {
      // Translate only untouched presets; user-authored goals remain verbatim.
      const presets = s.environment === 'sandbox' ? sandboxGoals : marketplaceGoals
      if (s.goal === presets[s.language]) s.goal = presets[a.payload]
      s.language = a.payload
    },
    setGoal: (s, a: PayloadAction<string>) => {
      s.goal = a.payload
    },
    setTarget: (s, a: PayloadAction<string>) => {
      s.target = a.payload
    },
    setScenario: (s, a: PayloadAction<string>) => {
      s.scenario = a.payload
    },
    setVariant: (s, a: PayloadAction<string>) => {
      s.variant = a.payload
    },
    selectRun: (s, a: PayloadAction<string>) => {
      if (s.selected !== a.payload) s.stepNumber = null
      s.selected = a.payload
    },
    selectStep: (s, a: PayloadAction<number>) => {
      s.stepNumber = a.payload
    },
    setDialog: (s, a: PayloadAction<string | null>) => {
      s.dialog = a.payload
    },
    selectEnvironment: (s, a: PayloadAction<Environment>) => {
      if (s.pending || s.connecting) return
      s.environment = a.payload
      s.error = ''
      s.browserNotice = ''
      s.goal = (a.payload === 'sandbox' ? sandboxGoals : marketplaceGoals)[s.language]
      s.target = a.payload === 'sandbox' ? s.health?.sandbox_url || initialState.target : ''
    },
    loadMarketplaceRun: (s, a: PayloadAction<Run>) => {
      if (s.pending || s.connecting) return
      s.environment = a.payload.environment || 'sandbox'
      s.target = a.payload.url
      s.goal = a.payload.goal
      s.error = ''
      s.browserNotice = ''
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(refreshDashboard.pending, (s) => {
        s.refreshing = true
      })
      .addCase(refreshDashboard.fulfilled, (s, a) => {
        s.refreshing = false
        // A poll begun before a mutation must not erase a newly started run/session.
        if (s.revision !== a.payload.revision) return
        s.health = a.payload.health
        s.runs = a.payload.runs
        s.marketSession = a.payload.marketSession
        s.connectionError = ''
        if (!s.selected || !s.runs.some((r) => r.id === s.selected)) {
          s.selected = s.runs[0]?.id || ''
          s.stepNumber = null
        }
      })
      .addCase(refreshDashboard.rejected, (s, a) => {
        s.refreshing = false
        if (!a.meta.aborted && !a.meta.condition)
          s.connectionError = a.payload || 'Backend unavailable'
      })
      .addCase(launchRun.pending, (s) => {
        s.pending = true
        s.error = ''
        s.revision++
      })
      .addCase(launchRun.fulfilled, (s, a) => {
        s.pending = false
        s.revision++
        s.runs = [a.payload, ...s.runs.filter((r) => r.id !== a.payload.id)]
        s.selected = a.payload.id
        s.stepNumber = null
      })
      .addCase(launchRun.rejected, (s, a) => {
        s.pending = false
        s.revision++
        s.error = a.payload || 'Could not start run'
      })
      .addCase(connectMarketplace.pending, (s) => {
        s.connecting = true
        s.error = ''
        s.browserNotice = ''
        s.revision++
      })
      .addCase(connectMarketplace.fulfilled, (s, a) => {
        s.connecting = false
        s.revision++
        s.marketSession = a.payload.session
        s.browserNotice =
          a.payload.connection.navigation && a.payload.connection.navigation !== 'loaded'
            ? 'The page did not finish loading. Retry it in the open browser before starting.'
            : ''
      })
      .addCase(connectMarketplace.rejected, (s, a) => {
        s.connecting = false
        s.revision++
        s.error = a.payload || 'Could not open browser'
      })
      .addCase(stopRun.pending, (s) => {
        s.error = ''
        s.revision++
      })
      .addCase(stopRun.fulfilled, (s) => {
        s.revision++
      })
      .addCase(stopRun.rejected, (s, a) => {
        s.revision++
        s.error = a.payload || 'Could not stop run'
      })
  },
})
export const actions = checkoutSlice.actions
