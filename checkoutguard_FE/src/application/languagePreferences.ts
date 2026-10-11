import { defaultLanguage, isLanguage, type Language } from '../domain/language'

export interface LanguagePreferenceRepository {
  read(): unknown
  write(language: Language): void
}

export class LanguagePreferences {
  constructor(private readonly repository: LanguagePreferenceRepository) {}
  current(): Language {
    const value = this.repository.read()
    return isLanguage(value) ? value : defaultLanguage
  }
  save(language: Language) {
    this.repository.write(language)
  }
}
