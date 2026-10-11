export type Language = 'en' | 'id'
export const defaultLanguage: Language = 'en'
export const isLanguage = (value: unknown): value is Language => value === 'en' || value === 'id'
export const languageLocale = (language: Language) => (language === 'id' ? 'id-ID' : 'en-US')
