import { useTranslation } from '../i18n/useTranslation'
import { useEffect } from 'react'
import type { Environment, Run } from '../../domain/checkout'
import { isMarket } from '../../domain/checkoutRules'
import { actions, launchRun, refreshDashboard, stopRun } from '../state/checkoutSlice'
import { useAppDispatch, useAppSelector } from '../state/hooks'
import {
  defects,
  selectActive,
  selectBaseline,
  selectCheckout,
  selectMetrics,
  selectRun,
  selectStep,
} from '../state/selectors'
import { Icon, Metric, Status } from '../components/common'
import { money, toolLabels, environmentNames, runLabel } from './format'
import { CheckoutComposer } from './CheckoutComposer'
import { RunHistory } from './RunHistory'
import { BrowserEvidence } from './BrowserEvidence'
import { LanguagePicker } from '../components/LanguagePicker'

export function Dashboard() {
  const { t, message, language, locale } = useTranslation()
  const dispatch = useAppDispatch()
  const { runs, environment, connectionError, pending, dialog } = useAppSelector(selectCheckout)
  const run = useAppSelector(selectRun),
    step = useAppSelector(selectStep),
    baseline = useAppSelector(selectBaseline)
  const active = useAppSelector(selectActive)
  const { verified, totalTools } = useAppSelector(selectMetrics)
  useEffect(() => {
    document.documentElement.lang = language
    document.title = t('CheckoutGuard · Release with evidence')
  }, [language])
  useEffect(() => {
    let busy = false
    let task: ReturnType<ReturnType<typeof refreshDashboard>> | undefined
    const refresh = () => {
      if (busy) return
      busy = true
      const next = dispatch(refreshDashboard())
      // Retain the actual in-flight task so unmount can abort it.
      task = next
      void next.finally(() => {
        busy = false
      })
    }
    refresh()
    const timer = setInterval(refresh, 2000)
    return () => {
      clearInterval(timer)
      task?.abort()
    }
  }, [dispatch])
  useEffect(() => {
    if (!dialog) return
    const close = (event: KeyboardEvent) => {
      if (event.key === 'Escape') dispatch(actions.setDialog(null))
    }
    window.addEventListener('keydown', close)
    return () => window.removeEventListener('keydown', close)
  }, [dialog, dispatch])
  const showDialog = (url: string | null) => dispatch(actions.setDialog(url))
  const launch = (retest?: string) => dispatch(launchRun(retest))
  const selectEnvironment = (value: Environment) => dispatch(actions.selectEnvironment(value))
  const loadMarketplaceRun = (r: Run) => {
    dispatch(actions.loadMarketplaceRun(r))
    document.getElementById('goal')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
  const stopSelectedRun = () => {
    if (run) void dispatch(stopRun(run.id))
  }
  return (
    <div className="shell">
      <aside className="sidebar">
        <a href="#" className="brand" aria-label={t('CheckoutGuard home')}>
          <img
            className="brand-mark"
            src="/brand/checkoutguard-mark.svg"
            alt=""
            width="40"
            height="40"
          />
          <span className="brand-wordmark">
            <span>Checkout</span>
            <span className="brand-accent">Guard</span>
          </span>
        </a>
        <div className="workspace">
          <span className="workspace-avatar">F</span>
          <div>
            <strong>FLdevBoo</strong>
            <small>{t('neatHack workspace')}</small>
          </div>
          <span className="chevron">⌄</span>
        </div>
        <p className="nav-label">{t('VERIFICATION')}</p>
        <a className="nav-item current" href="#">
          <Icon name="grid" />
          {t('Overview')}
          <span className="nav-dot" />
        </a>
        <a className="nav-item" href="#run-history">
          <Icon name="layers" />
          {t('Run history')}
          <span className="count">{runs.length}</span>
        </a>
        <a className="nav-item" href="#evidence">
          <Icon name="shield" />
          {t('Evidence')}
        </a>
        <div className="sidebar-bottom">
          <div className="partner-title">{t('BUILT WITH')}</div>
          <div className="partner">
            neatlogs <span>{t('Workflow traces')}</span>
          </div>
          <div className="partner">
            Entire <span>{t('Coding capture')}</span>
          </div>
          <a
            className="partner partner-link"
            href="/business/CheckoutGuard-business-plan.pdf"
            target="_blank"
            rel="noreferrer"
            title={t('Open business plan')}
          >
            cfo.ai <span>{t('Business plan')} ↗</span>
          </a>
          <div className="profile">
            <span className="profile-avatar">FL</span>
            <div>
              <strong>FLdevBoo</strong>
              <small>{t('Local development')}</small>
            </div>
            <span className="status-dot" />
          </div>
        </div>
      </aside>
      <div className="main-shell">
        <header className="topbar">
          <div className="topbar-breadcrumb">
            {t('Checkout verification')}
            <span>/</span> <strong>{t('Overview')}</strong>
          </div>
          <div className="topbar-actions">
            <span className="environment">
              <span className="status-dot" />
              {environment === 'sandbox'
                ? t('Synthetic sandbox')
                : environmentNames[environment] + ' · ' + t('Prepare checkout')}
            </span>
            <LanguagePicker />
          </div>
        </header>
        <main>
          <section className="page-heading">
            <div className="hero-copy">
              <div className="eyebrow">
                <span className="hero-eyebrow-icon">
                  <Icon name="shield" size={14} />
                </span>
                {t('CHECKOUT QUALITY ASSURANCE')}
              </div>
              <h1>
                {t('Check your checkout.')}
                <br />
                {t('See the evidence.')}
              </h1>
              <p>
                {t(
                  'Verify totals, confirm orders, and investigate failures with every browser step saved.',
                )}
              </p>
              <div className="hero-capabilities">
                <span>
                  <Icon name="check" size={14} />
                  {t('Price verification')}
                </span>
                <span>
                  <Icon name="layers" size={14} />
                  {t('Order confirmation')}
                </span>
                <span>
                  <Icon name="retry" size={14} />
                  {t('Safe recovery')}
                </span>
              </div>
            </div>
            <div className="hero-art" aria-hidden="true">
              <div className="hero-orbit" />
              <img
                src="/images/checkoutguard-shield-3d.png"
                alt=""
                width="1536"
                height="1024"
                fetchPriority="high"
              />
            </div>
            <span className="local-tag">{t('neatHack · Local MVP')}</span>
          </section>
          {connectionError && (
            <div className="error-banner" role="alert">
              {t('Backend unavailable: {error}. Start the local backend to run checkout checks.', {
                error: message(connectionError),
              })}
            </div>
          )}
          <div className="metrics">
            <Metric
              name={t('Saved checks')}
              value={String(runs.length)}
              note={t('All runs in this workspace')}
              icon="layers"
            />
            <Metric
              name={t('Checks passed')}
              value={String(verified)}
              note={t('Correct total · One order confirmed')}
              icon="shield"
            />
            <Metric
              name={t('Verified defects')}
              value={String(runs.reduce((n, r) => n + defects(r), 0))}
              note={t('Verified findings across checks')}
              icon="alert"
            />
            <Metric
              name={t('Agent tool calls')}
              value={String(totalTools)}
              note={t('Browser actions and verification')}
              icon="bolt"
            />
          </div>
          <CheckoutComposer />
          {run && (
            <section className="card verification-summary" aria-label={t('Selected check summary')}>
              <div className="verification-heading">
                <div>
                  <span className="eyebrow">{t('SELECTED CHECK')}</span>
                  <h2>
                    {runLabel(run, language)} <span className="mono">{run.id.slice(0, 8)}</span>
                  </h2>
                  <p>{run.goal}</p>
                </div>
                <Status value={run.status} />
              </div>
              <div className="verification-values">
                <div>
                  <span>{t('Expected total')}</span>
                  <strong>
                    {run.total_verified
                      ? money(run.expected_total_minor, language)
                      : t('Not checked')}
                  </strong>
                  <small>
                    {isMarket(run)
                      ? t('Marketplace breakdown not verified')
                      : t('Calculated from the price breakdown')}
                  </small>
                </div>
                <div
                  className={
                    run.total_verified && run.displayed_total_minor !== run.expected_total_minor
                      ? 'amount-mismatch'
                      : ''
                  }
                >
                  <span>{t('Displayed total')}</span>
                  <strong>
                    {run.total_verified
                      ? money(run.displayed_total_minor, language)
                      : t('Not checked')}
                  </strong>
                  <small>
                    {run.total_verified
                      ? run.displayed_total_minor === run.expected_total_minor
                        ? t('Matches the expected amount')
                        : t('Does not match the expected amount')
                      : t('Awaiting total verification')}
                  </small>
                </div>
                <div>
                  <span>{isMarket(run) ? t('Order confirmation') : t('Persisted orders')}</span>
                  <strong>
                    {isMarket(run) ? (
                      t('Manual')
                    ) : (
                      <>{t('{count} order(s)', { count: run.order_count })}</>
                    )}
                  </strong>
                  <small>
                    {isMarket(run)
                      ? t('No order submitted by the agent')
                      : run.order_count === 1
                        ? t('Exactly one order recorded')
                        : run.order_count === 0
                          ? t('No order recorded')
                          : t('Multiple orders recorded')}
                  </small>
                </div>
                <div>
                  <span>{t('Run duration')}</span>
                  <strong>
                    {Math.max(
                      0,
                      (new Date(run.finished_at || Date.now()).getTime() -
                        new Date(run.started_at).getTime()) /
                        1000,
                    ).toLocaleString(locale, {
                      minimumFractionDigits: 1,
                      maximumFractionDigits: 1,
                    })}
                    <span className="value-unit"> {t('sec')}</span>
                  </strong>
                  <small>
                    {run.finished_at
                      ? run.status === 'failed' || run.status === 'stopped'
                        ? t('Stopped before completion')
                        : t('Finished check')
                      : t('Check in progress')}
                  </small>
                </div>
              </div>
              {run.status === 'passed' && run.steps.some((s) => s.error) && (
                <div className="recovery-note">
                  <Icon name="retry" size={15} />
                  {t('Recovered from {count} step issue(s). Final checkout checks passed.', {
                    count: run.steps.filter((s) => s.error).length,
                  })}
                </div>
              )}
            </section>
          )}
          {run && isMarket(run) && run.status !== 'running' && (
            <section className="card marketplace-result">
              <div>
                <h2>
                  {run.status === 'ready_for_review'
                    ? t('Review checkout in your browser')
                    : t('Keep the marketplace result, continue with the demo')}
                </h2>
                <p>{message(run.reason || '')}</p>
                {run.checkout_review && (
                  <p className="review-summary">{run.checkout_review.summary}</p>
                )}
              </div>
              <button
                className="secondary"
                type="button"
                disabled={active || pending}
                onClick={() => {
                  selectEnvironment('sandbox')
                  document
                    .getElementById('goal')
                    ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
                }}
              >
                {t('Use Sandbox demo')}
                <Icon name="arrow" size={15} />
              </button>
              <small>
                {t(
                  'This selects the demo form. It does not start a run or replace the marketplace result.',
                )}
              </small>
            </section>
          )}
          <div className="execution-grid">
            <section className="card journey">
              <div className="section-header">
                <div>
                  <h2>{t('Agent activity')}</h2>
                  <p>
                    {run
                      ? isMarket(run)
                        ? runLabel(run, language) + ' · ' + t('Manual final confirmation')
                        : runLabel(run, language) +
                          ' · ' +
                          t('{basket} basket', {
                            basket: t(
                              run.variant === 'alternate' ? 'Alternate prices' : 'Standard',
                            ),
                          })
                      : t('Your next checkout journey starts here')}
                  </p>
                </div>
                {run && <Status value={run.status} />}
              </div>
              {!run ? (
                <div className="empty">
                  <span className="empty-icon">
                    <Icon name="layers" size={28} />
                  </span>
                  <h3>{t('Your check starts with a goal.')}</h3>
                  <p>{t('Start a run to see the agent observe, act, and verify the checkout.')}</p>
                </div>
              ) : (
                <>
                  <div className="run-summary">
                    <span className="mono">
                      {t('{count} saved step(s)', { count: run.steps.length })}
                    </span>
                    <div>
                      {run.status === 'running' ? (
                        <button className="text-button" onClick={() => void stopSelectedRun()}>
                          {t('Stop run')}
                        </button>
                      ) : isMarket(run) ? (
                        <button
                          className="text-button"
                          disabled={active || pending}
                          onClick={() => loadMarketplaceRun(run)}
                        >
                          <Icon name="retry" size={14} />
                          {t('Load this goal again')}
                        </button>
                      ) : (
                        <button
                          className="text-button"
                          disabled={active || pending}
                          onClick={() => void launch(run.id)}
                        >
                          <Icon name="retry" size={14} />
                          {t('Retest corrected store')}
                        </button>
                      )}
                    </div>
                  </div>
                  <div className="timeline">
                    {run.steps.map((s) => (
                      <button
                        key={s.number}
                        className={'timeline-step ' + (step?.number === s.number ? 'selected' : '')}
                        aria-pressed={step?.number === s.number}
                        onClick={() => dispatch(actions.selectStep(s.number))}
                      >
                        <span className={'step-dot ' + (s.error ? 'error-dot' : '')}>
                          {s.error ? '!' : s.number}
                        </span>
                        <span className="step-body">
                          <span className="step-top">
                            <strong>
                              {toolLabels[s.action.tool]
                                ? t(toolLabels[s.action.tool])
                                : s.action.tool}
                            </strong>
                            <span>
                              {s.latency_ms > 0
                                ? (s.latency_ms / 1000).toLocaleString(locale, {
                                    minimumFractionDigits: 1,
                                    maximumFractionDigits: 1,
                                  }) +
                                  ' ' +
                                  t('sec')
                                : t('Initial')}
                            </span>
                          </span>
                          <span>{message(s.action.reason)}</span>
                          {s.action.selector && (
                            <code className="step-selector">{s.action.selector}</code>
                          )}
                          {s.error && <span className="step-error">{message(s.error)}</span>}
                        </span>
                      </button>
                    ))}
                    {run.status === 'running' && (
                      <div className="thinking">
                        <span className="pulse-dot" />
                        {t('Agent is choosing its next action…')}
                      </div>
                    )}
                  </div>
                  {run.reason && (
                    <div className={'outcome ' + run.status}>{message(run.reason)}</div>
                  )}
                </>
              )}
            </section>
            <BrowserEvidence />
          </div>
          {run && run.findings.length > 0 && (
            <section className="card findings">
              <div className="section-header">
                <div>
                  <h2>{t('What the agent found')}</h2>
                  <p>{t('Findings backed by amounts and persisted order data.')}</p>
                </div>
                <span className="quiet-badge">
                  {t('{count} finding(s)', { count: run.findings.length })}
                </span>
              </div>
              {run.findings.map((f) => (
                <div
                  className={'finding ' + (f.kind === 'timeout_recovered' ? 'recovery' : '')}
                  key={f.kind}
                >
                  <span className="finding-icon">
                    <Icon name={f.kind === 'timeout_recovered' ? 'retry' : 'alert'} />
                  </span>
                  <div>
                    <strong>{message(f.title)}</strong>
                    <p>{message(f.detail)}</p>
                    <small>
                      {f.verified ? t('Deterministically verified') : t('Needs verification')}
                    </small>
                  </div>
                  <button className="text-button" onClick={() => showDialog(f.evidence)}>
                    {t('View evidence')}
                    <Icon name="arrow" size={15} />
                  </button>
                </div>
              ))}
            </section>
          )}
          {baseline && run && (
            <section className="card comparison">
              <div className="section-header">
                <div>
                  <h2>{t('Before & after')}</h2>
                  <p>{t('The same goal and basket, against the corrected synthetic store.')}</p>
                </div>
                <Icon name="retry" />
              </div>
              <div className="comparison-grid">
                {[baseline, run].map((r, i) => (
                  <div key={r.id}>
                    <span className="eyebrow">{i ? t('RETEST') : t('BASELINE')}</span>
                    <div className="comparison-status">
                      <Status value={r.status} />
                      <span>
                        {t('{defects} defect(s) · {orders} order(s)', {
                          defects: defects(r),
                          orders: r.order_count,
                        })}
                      </span>
                    </div>
                    <p>
                      {t('Expected {amount}', { amount: money(r.expected_total_minor, language) })}
                      <br />
                      {t('Displayed {amount}', {
                        amount: money(r.displayed_total_minor, language),
                      })}
                    </p>
                    {r.steps.at(-1)?.evidence && (
                      <button
                        className="text-button"
                        onClick={() => showDialog(r.steps.at(-1)?.evidence || null)}
                      >
                        {t('Open final evidence')}
                        <Icon name="link" size={13} />
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </section>
          )}
          <RunHistory />
          {run?.trace_id && (
            <div className="trace-line">
              <Icon name="layers" size={14} />
              <span>{t('neatlogs trace')}</span>
              <code>{run.trace_id}</code>
              <button
                className="text-button"
                onClick={() => void navigator.clipboard.writeText(run.trace_id || '')}
              >
                {t('Copy ID')}
              </button>
            </div>
          )}
          <footer>
            <span>
              <Icon name="shield" size={14} />
              {t('CheckoutGuard · Evidence before release.')}
            </span>
            <span>{t('Sandbox demo · Marketplace checkout preparation is experimental')}</span>
          </footer>
        </main>
      </div>
      {dialog && (
        <div
          className="modal"
          role="dialog"
          aria-modal="true"
          aria-label={t('Screenshot evidence')}
          onClick={() => showDialog(null)}
        >
          <button className="modal-close" onClick={() => showDialog(null)}>
            {t('Close ×')}
          </button>
          <img
            src={dialog}
            alt={t('Full browser screenshot evidence')}
            onClick={(e) => e.stopPropagation()}
          />
        </div>
      )}
    </div>
  )
}
