import type {
  DashboardSnapshot,
  Environment,
  MarketConnection,
  MarketSession,
  Run,
  StartRunInput,
} from '../domain/checkout'

export interface CheckoutRepository {
  snapshot(signal?: AbortSignal): Promise<DashboardSnapshot>
  start(input: StartRunInput, signal?: AbortSignal): Promise<Run>
  retest(id: string, signal?: AbortSignal): Promise<Run>
  stop(id: string, signal?: AbortSignal): Promise<void>
  connect(environment: Environment, signal?: AbortSignal): Promise<MarketConnection>
  marketSession(signal?: AbortSignal): Promise<MarketSession>
}
