import { mergeConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

/** @type { import('storybook').StorybookConfig } */
export default {
  stories: ['../src/**/*.stories.js'],
  addons: [
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
  ],
  framework: {
    name: '@storybook/vue3-vite',
    options: {},
  },
  async viteFinal(config) {
    return mergeConfig(config, {
      plugins: [tailwindcss()],
      resolve: {
        alias: {
          '@wailsio/runtime': path.resolve(__dirname, 'mocks/wails-runtime.js'),
        },
      },
    })
  },
}
