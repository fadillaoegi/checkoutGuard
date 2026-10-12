import { configureStore } from '@reduxjs/toolkit'
import { checkoutSlice, type Dependencies } from './checkoutSlice'
import type { LanguagePreferences } from '../../application/languagePreferences'
import { sandboxGoals } from '../../domain/checkoutRules'
import { defaultLanguage } from '../../domain/language'
import { navigationSlice } from './navigationSlice'

export function createAppStore(dependencies: Dependencies, preferences?: LanguagePreferences) {
  let language = preferences?.current() ?? defaultLanguage
  const store = configureStore({
    reducer: { checkout: checkoutSlice.reducer, navigation: navigationSlice.reducer },
    preloadedState: {
      checkout: { ...checkoutSlice.getInitialState(), language, goal: sandboxGoals[language] },
    },
    middleware: (defaults) => defaults({ thunk: { extraArgument: dependencies } }),
  })
  if (preferences)
    store.subscribe(() => {
      const next = store.getState().checkout.language
      if (next !== language) {
        language = next
        preferences.save(next)
      }
    })
  return store
}
export type AppStore = ReturnType<typeof createAppStore>
export type RootState = ReturnType<AppStore['getState']>
export type AppDispatch = AppStore['dispatch']
