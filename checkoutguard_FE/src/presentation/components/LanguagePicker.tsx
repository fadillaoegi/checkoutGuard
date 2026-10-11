import { isLanguage } from '../../domain/language'
import { actions } from '../state/checkoutSlice'
import { useAppDispatch } from '../state/hooks'
import { useTranslation } from '../i18n/useTranslation'

export function LanguagePicker() {
  const dispatch = useAppDispatch()
  const { t, language } = useTranslation()
  return (
    <label className="language-picker">
      <span>{t('Language')}</span>
      <select
        id="language"
        value={language}
        onChange={(event) => {
          const value = event.target.value
          if (isLanguage(value)) dispatch(actions.setLanguage(value))
        }}
      >
        <option value="en" lang="en">
          English
        </option>
        <option value="id" lang="id">
          Bahasa Indonesia
        </option>
      </select>
    </label>
  )
}
