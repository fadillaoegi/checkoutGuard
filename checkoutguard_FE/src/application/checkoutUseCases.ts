import type { CheckoutRepository } from './CheckoutRepository'
import type { Environment, StartRunInput } from '../domain/checkout'

// Framework independent use cases. HTTP and Redux are outside this layer.
export class CheckoutUseCases {
  constructor(private readonly repository: CheckoutRepository) {}
  refresh(signal?: AbortSignal) {
    return this.repository.snapshot(signal)
  }
  start(input: StartRunInput, signal?: AbortSignal) {
    return this.repository.start(input, signal)
  }
  retest(id: string, signal?: AbortSignal) {
    return this.repository.retest(id, signal)
  }
  stop(id: string, signal?: AbortSignal) {
    return this.repository.stop(id, signal)
  }
  async connect(environment: Environment, signal?: AbortSignal) {
    if (environment === 'sandbox') throw new Error('Choose a marketplace to open its browser')
    const connection = await this.repository.connect(environment, signal)
    return { connection, session: await this.repository.marketSession(signal) }
  }
}
