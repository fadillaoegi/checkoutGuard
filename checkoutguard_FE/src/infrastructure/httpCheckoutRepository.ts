import type { CheckoutRepository } from '../application/CheckoutRepository'
import type {
  DashboardSnapshot,
  Environment,
  Health,
  MarketConnection,
  MarketSession,
  Run,
  StartRunInput,
} from '../domain/checkout'

export class HttpCheckoutRepository implements CheckoutRepository {
  constructor(
    private readonly request: typeof fetch = (...args) => globalThis.fetch(...args),
    private readonly baseURL = '/api',
  ) {}
  private async api<T>(path: string, body?: unknown, signal?: AbortSignal): Promise<T> {
    // Native browser fetch must keep its Window receiver, not the repository.
    const request = this.request
    const response = await request(this.baseURL + path, {
      signal,
      ...(body === undefined
        ? {}
        : {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
          }),
    })
    const result = await response.json()
    if (!response.ok) throw new Error(result.error || 'Request failed')
    return result
  }
  async snapshot(signal?: AbortSignal): Promise<DashboardSnapshot> {
    const [health, runs, marketSession] = await Promise.all([
      this.api<Health>('/health', undefined, signal),
      this.api<Run[]>('/runs', undefined, signal),
      this.marketSession(signal),
    ])
    return { health, runs, marketSession }
  }
  start(input: StartRunInput, signal?: AbortSignal) {
    return this.api<Run>('/runs', input, signal)
  }
  retest(id: string, signal?: AbortSignal) {
    return this.api<Run>('/runs/' + encodeURIComponent(id) + '/retest', {}, signal)
  }
  async stop(id: string, signal?: AbortSignal) {
    await this.api('/runs/' + encodeURIComponent(id) + '/stop', {}, signal)
  }
  connect(environment: Environment, signal?: AbortSignal) {
    return this.api<MarketConnection>('/marketplaces/' + environment + '/connect', {}, signal)
  }
  marketSession(signal?: AbortSignal) {
    return this.api<MarketSession>('/marketplaces/session', undefined, signal)
  }
}
