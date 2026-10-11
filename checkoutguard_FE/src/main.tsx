import React from 'react'
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'
import { CheckoutUseCases } from './application/checkoutUseCases'
import { HttpCheckoutRepository } from './infrastructure/httpCheckoutRepository'
import { LanguagePreferences } from './application/languagePreferences'
import { BrowserLanguagePreferences } from './infrastructure/browserLanguagePreferences'
import { createAppStore } from './presentation/state/store'
import { Dashboard } from './presentation/dashboard/Dashboard'
import './presentation/style.css'

// Composition root: this is the only place that chooses concrete adapters.
const store = createAppStore(
  { checkout: new CheckoutUseCases(new HttpCheckoutRepository()) },
  new LanguagePreferences(new BrowserLanguagePreferences()),
)
createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider store={store}>
      <Dashboard />
    </Provider>
  </React.StrictMode>,
)
