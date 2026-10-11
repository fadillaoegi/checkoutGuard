import type { LanguagePreferenceRepository } from '../application/languagePreferences'
import type { Language } from '../domain/language'

export const languageStorageKey = 'checkoutguard.language'
export class BrowserLanguagePreferences implements LanguagePreferenceRepository {
  constructor(private readonly storage: () => Storage = () => window.localStorage) {}
  read(): unknown {
    try {
      return this.storage().getItem(languageStorageKey)
    } catch {
      return null
    }
  }
  write(language: Language) {
    try {
      this.storage().setItem(languageStorageKey, language)
    } catch {
      // The in-memory choice still works when browser storage is unavailable.
    }
  }
}
