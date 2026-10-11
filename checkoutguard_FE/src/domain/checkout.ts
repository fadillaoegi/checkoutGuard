export type Usage = {
  input_tokens: number
  output_tokens: number
  cost_usd: number | null
  cost_incomplete?: boolean
}
export type Step = {
  number: number
  action: { tool: string; selector?: string; reason: string }
  result: unknown
  evidence?: string
  latency_ms: number
  usage: Usage
  error?: string
}
export type Finding = {
  kind: string
  title: string
  detail: string
  evidence: string
  verified: boolean
}
export type Environment = 'sandbox' | 'shopee' | 'tokopedia'
export type MarketPageState = 'unknown' | 'browsing' | 'needs_login' | 'blocked' | 'review'
export type MarketSession = {
  environment: string
  browser_open: boolean
  page_state?: MarketPageState
  page_check?: 'checked' | 'unavailable'
}
export type Run = {
  environment?: Environment
  checkout_review?: { url: string; summary: string; order_submitted: boolean }
  id: string
  parent_id?: string
  goal: string
  url: string
  scenario: string
  variant: string
  status: string
  reason?: string
  model: string
  requested_model?: string
  trace_id?: string
  started_at: string
  finished_at?: string
  steps: Step[]
  findings: Finding[]
  usage: Usage
  order_count: number
  total_verified: boolean
  expected_total_minor: number
  displayed_total_minor: number
}
export type Health = { model_configured: boolean; model: string; sandbox_url: string }

export type StartRunInput = {
  environment: Environment
  goal: string
  url: string
  scenario: string
  variant: string
}
export type DashboardSnapshot = { health: Health; runs: Run[]; marketSession: MarketSession }
export type MarketConnection = {
  navigation?: string
  page_state?: MarketPageState
  page_check?: 'checked' | 'unavailable'
}
