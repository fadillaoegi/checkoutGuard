import { languageLocale } from '../../domain/language'
import { useAppSelector } from '../state/hooks'
import { localizeMessage, translate, type TranslationKey, type TranslationParams } from './messages'

export function useTranslation() {
  const language = useAppSelector((s) => s.checkout.language)
  return {
    language,
    locale: languageLocale(language),
    t: (key: TranslationKey, params?: TranslationParams) => translate(language, key, params),
    message: (value: string) => localizeMessage(language, value),
  }
}
