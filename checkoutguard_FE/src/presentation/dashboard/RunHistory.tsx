import { useTranslation } from '../i18n/useTranslation'
import { actions } from '../state/checkoutSlice'
import { useAppDispatch, useAppSelector } from '../state/hooks'
import { selectCheckout } from '../state/selectors'
import { isMarket } from '../../domain/checkoutRules'
import { Icon, Status } from '../components/common'
import { cost, environmentNames, runLabel } from './format'

export function RunHistory() {
  const { t, language, locale } = useTranslation()
  const dispatch = useAppDispatch()
  const { runs, selected } = useAppSelector(selectCheckout)
  return (
    <section className="card history" id="run-history">
      <div className="section-header">
        <div>
          <h2>{t('Latest runs')}</h2>
          <p>{t('A history of decisions, outcomes, and evidence.')}</p>
        </div>
        <span className="quiet-badge">{t('{count} saved', { count: runs.length })}</span>
      </div>
      <div className="table-scroll">
        <table>
          <thead>
            <tr>
              <th>{t('RUN / CONDITION')}</th>
              <th>{t('STARTED')}</th>
              <th>{t('STATUS')}</th>
              <th>{t('ORDERS')}</th>
              <th>{t('STEPS')}</th>
              <th>{t('COST')}</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {runs.map((r) => (
              <tr
                key={r.id}
                className={selected === r.id ? 'active-row' : ''}
                onClick={() => dispatch(actions.selectRun(r.id))}
              >
                <td>
                  <strong>{runLabel(r, language)}</strong>
                  <small>
                    {r.id.slice(0, 8)}
                    {r.parent_id ? ' · ' + t('retest') : ''} ·{' '}
                    {environmentNames[r.environment || 'sandbox']}
                  </small>
                </td>
                <td>
                  {new Date(r.started_at).toLocaleTimeString(locale, {
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </td>
                <td>
                  <Status value={r.status} />
                </td>
                <td>{isMarket(r) ? '—' : r.order_count}</td>
                <td>{r.steps.length}</td>
                <td>{cost(r.usage, language)}</td>
                <td>
                  <button
                    className="icon-button"
                    aria-label={t('Inspect run {id}', { id: r.id.slice(0, 8) })}
                    onClick={() => dispatch(actions.selectRun(r.id))}
                  >
                    <Icon name="arrow" size={16} />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {runs.length === 0 && (
          <div className="table-empty">
            {t('No runs yet. Your first checkout check will appear here.')}
          </div>
        )}
      </div>
    </section>
  )
}
