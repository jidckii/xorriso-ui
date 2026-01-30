import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'

function getInitialLocale() {
  // TODO: load from Wails GetSettings() when bindings are ready
  const saved = localStorage.getItem('xorriso-language')
  return saved || 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'en',
  messages: { en, ru },
})

export default i18n
