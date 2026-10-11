import { useTranslation } from '../i18n/useTranslation'
import type React from 'react'

export function Icon({ name, size = 20 }: { name: string; size?: number }) {
  const paths: Record<string, React.ReactNode> = {
    shield: (
      <>
        <path d="M12 3 4 6v6c0 5 8 9 8 9s8-4 8-9V6l-8-3Z" />
        <path d="m8 12 3 3 5-6" />
      </>
    ),
    grid: (
      <>
        <rect x="3" y="3" width="7" height="7" rx="1.5" />
        <rect x="14" y="3" width="7" height="7" rx="1.5" />
        <rect x="3" y="14" width="7" height="7" rx="1.5" />
        <rect x="14" y="14" width="7" height="7" rx="1.5" />
      </>
    ),
    play: <path d="m8 5 12 7-12 7V5Z" />,
    arrow: (
      <>
        <path d="M5 12h14M13 6l6 6-6 6" />
      </>
    ),
    check: <path d="m5 12 4 4L19 6" />,
    clock: (
      <>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v5l3 2" />
      </>
    ),
    layers: (
      <>
        <path d="m12 3 10 5-10 5L2 8l10-5Z" />
        <path d="m2 12 10 5 10-5M2 16l10 5 10-5" />
      </>
    ),
    retry: (
      <>
        <path d="M3 10a9 9 0 1 1 2 9M3 3v7h7" />
      </>
    ),
    link: (
      <>
        <path d="M14 3h7v7M21 3l-9 9" />
        <path d="M10 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5" />
      </>
    ),
    bolt: <path d="m13 2-9 12h7l-1 8 10-13h-7l0-7Z" />,
    alert: (
      <>
        <path d="m12 3 10 18H2L12 3Z" />
        <path d="M12 9v5m0 3v1" />
      </>
    ),
  }
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.7"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      {paths[name] || paths.layers}
    </svg>
  )
}
export function Status({ value }: { value: string }) {
  const { t } = useTranslation()
  return (
    <span className={'status ' + value}>
      <span />
      {(
        {
          ready_for_review: t('Ready for review'),
          needs_login: t('Sign-in needed'),
          blocked: t('Blocked'),
          passed: t('Passed'),
          findings: t('Defect found'),
          running: t('Running'),
          failed: t('Failed'),
          stopped: t('Stopped'),
        } as Record<string, string>
      )[value] || value}
    </span>
  )
}
export function Metric({
  name,
  value,
  note,
  icon,
}: {
  name: string
  value: string
  note: string
  icon: string
}) {
  return (
    <div className="metric">
      <div className="metric-top">
        <span>{name}</span>
        <Icon name={icon} size={18} />
      </div>
      <strong>{value}</strong>
      <small>{note}</small>
    </div>
  )
}
