import { useTranslation } from '../i18n/useTranslation'
import type { Environment } from '../../domain/checkout'
import { actions, connectMarketplace, launchRun } from '../state/checkoutSlice'
import { useAppDispatch, useAppSelector } from '../state/hooks'
import {
  selectActive,
  selectCheckout,
  selectMarketBrowserOpen,
  selectMarketPageState,
} from '../state/selectors'
import { Icon } from '../components/common'
import { environmentNames, labels } from './format'

export function CheckoutComposer() {
  const { t, message } = useTranslation()
  const dispatch = useAppDispatch()
  const {
    health,
    goal,
    target,
    environment,
    connecting,
    browserNotice,
    scenario,
    variant,
    error,
    connectionError,
    marketSession,
    pending,
  } = useAppSelector(selectCheckout)
  const active = useAppSelector(selectActive),
    marketBrowserOpen = useAppSelector(selectMarketBrowserOpen)
  const pageState = useAppSelector(selectMarketPageState)
  const pageCheckFailed = marketBrowserOpen && marketSession?.page_check === 'unavailable'
  const needsManualAccess =
    pageState === 'blocked' || pageState === 'needs_login' || pageCheckFailed
  const launch = () => dispatch(launchRun(undefined))
  const selectEnvironment = (value: Environment) => dispatch(actions.selectEnvironment(value))
  const connectBrowser = () => dispatch(connectMarketplace())
  return (
    <section className="card composer">
      <div className="card-title">
        <div className="section-icon">
          <Icon name="play" />
        </div>
        <div>
          <h2>{t('Start a checkout check')}</h2>
          <p>{t('Describe the purchase. The agent will act, verify, and save the evidence.')}</p>
        </div>
        <span className="quiet-badge">{t('One active check at a time')}</span>
      </div>
      <div className="environment-picker" role="radiogroup" aria-label={t('Checkout environment')}>
        {(['sandbox', 'shopee', 'tokopedia'] as Environment[]).map((value) => (
          <button
            type="button"
            role="radio"
            aria-checked={environment === value}
            key={value}
            className={'environment-choice ' + (environment === value ? 'chosen' : '')}
            disabled={pending || connecting}
            onClick={() => selectEnvironment(value)}
          >
            <span className="environment-symbol">
              {value === 'sandbox' ? <Icon name="shield" /> : value === 'shopee' ? 'S' : 'T'}
            </span>
            <span>
              <strong>{environmentNames[value]}</strong>
              <small>
                {value === 'sandbox'
                  ? t('Tested demo · Simulated orders')
                  : t('Experimental · Prepare checkout')}
              </small>
            </span>
            {environment === value && <Icon name="check" size={17} />}
          </button>
        ))}
      </div>
      {environment !== 'sandbox' && (
        <div className="marketplace-connect">
          <div>
            <strong>
              {marketBrowserOpen
                ? t('{marketplace} browser is open', { marketplace: environmentNames[environment] })
                : t('Sign in using your own browser session')}
            </strong>
            <p>
              {t(
                'Open the dedicated browser, sign in manually, then paste a product URL below. The agent stops before the final order is submitted.',
              )}
            </p>
            <small>
              {t(
                'Browser open does not mean sign-in has been verified. Access challenges require your action.',
              )}
            </small>
            {browserNotice && !needsManualAccess && <p role="status">{message(browserNotice)}</p>}
          </div>
          <button
            type="button"
            className="secondary"
            disabled={active || connecting || pending || !!connectionError}
            onClick={() => void connectBrowser()}
          >
            {connecting
              ? t('Opening…')
              : marketBrowserOpen
                ? t('Check browser')
                : t('Open {marketplace} browser', {
                    marketplace: environmentNames[environment],
                  })}{' '}
            <Icon name="link" size={14} />
          </button>
        </div>
      )}
      {environment !== 'sandbox' && needsManualAccess && (
        <div className="marketplace-access-notice" role="status">
          <div>
            <strong>
              {t(
                pageCheckFailed
                  ? 'Browser status unavailable'
                  : pageState === 'blocked'
                    ? 'Marketplace verification required'
                    : 'Sign-in needed',
              )}
            </strong>
            <p>
              {t(
                pageCheckFailed
                  ? 'CheckoutGuard could not read the current page. Wait until the browser responds, then click Check browser again. Checkout remains paused until the page can be checked.'
                  : pageState === 'blocked'
                    ? 'The marketplace is showing an access verification or unavailable page. If it says “Try again later”, wait a few minutes and try manually in the same browser. Then click Check browser. CheckoutGuard cannot complete this verification.'
                    : 'Sign in manually in the open marketplace browser, then click Check browser. CheckoutGuard does not enter your login details.',
              )}
            </p>
            <small>
              {t(
                'After checking the browser, page access is reassessed. This does not verify your login or guarantee checkout will work.',
              )}
            </small>
          </div>
          <button
            type="button"
            className="secondary"
            disabled={active || pending || connecting}
            onClick={() => selectEnvironment('sandbox')}
          >
            {t('Use Sandbox demo')} <Icon name="arrow" size={14} />
          </button>
        </div>
      )}
      <form
        onSubmit={(event) => {
          event.preventDefault()
          void launch()
        }}
      >
        <div className="goal-input">
          <label htmlFor="goal">{t('What should the agent do?')}</label>
          <textarea
            id="goal"
            value={goal}
            onChange={(e) => dispatch(actions.setGoal(e.target.value))}
            maxLength={1500}
            rows={2}
            required
          />
        </div>
        <div className={'form-bottom ' + (environment !== 'sandbox' ? 'market-form' : '')}>
          <div className="field url-field">
            <label htmlFor="target">
              {environment === 'sandbox'
                ? t('Sandbox URL')
                : t('{marketplace} product URL', { marketplace: environmentNames[environment] })}
            </label>
            <input
              id="target"
              type="url"
              value={target}
              placeholder={
                environment === 'shopee'
                  ? 'https://shopee.co.id/product/123/456'
                  : environment === 'tokopedia'
                    ? 'https://www.tokopedia.com/shop/product'
                    : undefined
              }
              onChange={(e) => dispatch(actions.setTarget(e.target.value))}
              required
            />
          </div>
          {environment === 'sandbox' && (
            <>
              <div className="field">
                <label htmlFor="scenario">{t('Sandbox condition')}</label>
                <select
                  id="scenario"
                  value={scenario}
                  onChange={(e) => dispatch(actions.setScenario(e.target.value))}
                >
                  {Object.entries(labels).map(([key, value]) => (
                    <option value={key} key={key}>
                      {t(value)}
                    </option>
                  ))}
                </select>
              </div>
              <div className="field">
                <label htmlFor="variant">{t('Basket')}</label>
                <select
                  id="variant"
                  value={variant}
                  onChange={(e) => dispatch(actions.setVariant(e.target.value))}
                >
                  <option value="standard">{t('Standard')}</option>
                  <option value="alternate">{t('Alternate prices')}</option>
                </select>
              </div>
            </>
          )}
          <button
            className="primary"
            disabled={
              pending ||
              active ||
              connecting ||
              !health?.model_configured ||
              !!connectionError ||
              (environment !== 'sandbox' && needsManualAccess) ||
              (environment !== 'sandbox' && !marketBrowserOpen)
            }
          >
            <Icon name="play" size={16} />
            {pending
              ? t('Starting…')
              : active
                ? t('Agent running')
                : environment === 'sandbox'
                  ? t('Run checkout')
                  : t('Prepare checkout')}
          </button>
        </div>
        <div className="form-caption">
          <span>
            <Icon name="shield" size={14} />
            {environment === 'sandbox'
              ? t('All orders and payments are simulated.')
              : t(
                  'Final order confirmation is manual. Marketplace totals and orders are not independently verified.',
                )}
          </span>
          <span title={health?.model}>
            {health?.model_configured ? t('AI model configured') : t('AI model not configured')}
          </span>
        </div>
        {health && !health.model_configured && (
          <p className="form-error">
            {t('Configure model credentials in the backend to start a run.')}
          </p>
        )}
        {error && (
          <p className="form-error" role="alert">
            {message(error)}
          </p>
        )}
      </form>
    </section>
  )
}
