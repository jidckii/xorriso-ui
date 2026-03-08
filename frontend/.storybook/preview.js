import { setup } from '@storybook/vue3'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import en from '../src/locales/en.json'
import ru from '../src/locales/ru.json'
import '../src/assets/css/main.css'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en, ru },
})

setup((app) => {
  app.use(createPinia())
  app.use(i18n)
  document.documentElement.classList.add('dark')
  document.body.style.backgroundColor = '#111827'
})

/** @type { import('storybook').Preview } */
export default {
  parameters: {
    backgrounds: { disable: true },
    layout: 'centered',
  },
}
