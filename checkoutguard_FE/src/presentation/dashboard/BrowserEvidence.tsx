import { useTranslation } from '../i18n/useTranslation'
import { actions } from '../state/checkoutSlice'
import { useAppDispatch, useAppSelector } from '../state/hooks'
import { selectRun, selectStep } from '../state/selectors'
import { isMarket } from '../../domain/checkoutRules'
import { Icon } from '../components/common'
import { cost, environmentNames } from './format'

export function BrowserEvidence() {
  const { t, language, locale } = useTranslation()
  const dispatch = useAppDispatch()
  const run = useAppSelector(selectRun),
    step = useAppSelector(selectStep)
  const showDialog = (url: string | null) => dispatch(actions.setDialog(url))
  return (
    <section className="card evidence-card" id="evidence">
      <div className="section-header">
        <div>
          <h2>{t('Browser evidence')}</h2>
          <p>
            {step
              ? t('Captured after step {number}', { number: step.number })
              : t('Real screenshots from the selected environment')}
          </p>
        </div>
        {step?.evidence && (
          <button
            className="icon-button"
            title={t('Open full screenshot')}
            aria-label={t('Open full screenshot')}
            onClick={() => showDialog(step.evidence || null)}
          >
            <Icon name="link" size={16} />
          </button>
        )}
      </div>
      <div className="browser-frame">
        <div className="browser-toolbar">
          <div className="browser-dots">
            <i />
            <i />
            <i />
          </div>
          <span>
            {run
              ? isMarket(run)
                ? environmentNames[run.environment!] + ' / checkout'
                : 'field-supply / checkout'
              : t('Waiting for a checkout run')}
          </span>
          <Icon name="shield" size={12} />
        </div>
        {step?.evidence ? (
          <button className="screenshot-button" onClick={() => showDialog(step.evidence || null)}>
            <img
              src={step.evidence}
              alt={t('Browser evidence after step {number}', { number: step.number })}
            />
          </button>
        ) : (
          <div className="browser-empty">
            <Icon name="shield" size={44} />
            <span>{t('Observe. Act. Verify.')}</span>
            <small>{t('Every step leaves a screenshot.')}</small>
          </div>
        )}
      </div>
      {run && (
        <div className="evidence-stats">
          <div>
            <span>{isMarket(run) ? t('Order verification') : t('Persisted orders')}</span>
            <strong>{isMarket(run) ? t('Unavailable') : run.order_count}</strong>
          </div>
          <div>
            <span>{t('Actual tokens')}</span>
            <strong>
              {(run.usage.input_tokens + run.usage.output_tokens).toLocaleString(locale)}
            </strong>
          </div>
          <div>
            <span>{t('Estimated API cost')}</span>
            <strong className="cost-value">{cost(run.usage, language)}</strong>
          </div>
        </div>
      )}
      {step && (
        <details className="tool-result">
          <summary>{t('Step {number} · tool result', { number: step.number })}</summary>
          <pre>{JSON.stringify(step.result, null, 2)}</pre>
        </details>
      )}
      {run?.requested_model && (
        <div className="model-record">
          {t('Returned model: {model}', { model: run.model || t('Waiting for provider') })}
        </div>
      )}
    </section>
  )
}
