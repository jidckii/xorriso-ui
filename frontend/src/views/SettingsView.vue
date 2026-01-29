<script setup>
import { ref, reactive, watch } from 'vue'
// TODO: import Wails service bindings when available
// import { GetSettings, SaveSettings } from '../../bindings/xorriso-ui/SettingsService'

// --- State ---
const saving = ref(false)
const saved = ref(false)

const settings = reactive({
  xorrisoBinaryPath: '/usr/bin/xorriso',
  devicePollInterval: 5,
  bdxlSafeMode: false,
  defaultBurnOptions: {
    speed: 'auto',
    dummyMode: false,
    verify: true,
    closeDisc: false,
    eject: true,
    burnMode: 'auto',
    streamRecording: false,
  },
  defaultIsoOptions: {
    rockRidge: true,
    joliet: true,
    md5: true,
    backupMode: false,
  },
})

// BDXL safe mode auto-toggle
watch(() => settings.bdxlSafeMode, (enabled) => {
  if (enabled) {
    settings.defaultIsoOptions.md5 = true
    settings.defaultBurnOptions.verify = true
    settings.defaultBurnOptions.streamRecording = true
  }
})

// --- Actions ---
async function loadSettings() {
  try {
    // TODO: replace with actual Wails call
    // const result = await GetSettings()
    // Object.assign(settings, result)
    console.log('Settings loaded (mock)')
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
}

async function saveSettings() {
  saving.value = true
  saved.value = false
  try {
    // TODO: replace with actual Wails call
    // await SaveSettings({ ...settings })
    console.log('Settings saved (mock):', JSON.parse(JSON.stringify(settings)))
    saved.value = true
    setTimeout(() => { saved.value = false }, 3000)
  } catch (error) {
    console.error('Failed to save settings:', error)
  } finally {
    saving.value = false
  }
}

// Load on mount
loadSettings()
</script>

<template>
  <div class="h-full overflow-y-auto">
    <div class="max-w-2xl mx-auto p-6 space-y-8">

      <div>
        <h1 class="text-xl font-semibold">Settings</h1>
        <p class="text-sm text-gray-500 mt-1">Configure xorriso-ui application preferences</p>
      </div>

      <!-- xorriso Binary Path -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wide">General</h2>

        <div>
          <label class="block text-sm text-gray-300 mb-1">xorriso binary path</label>
          <input
            v-model="settings.xorrisoBinaryPath"
            type="text"
            placeholder="/usr/bin/xorriso"
            class="w-full px-3 py-2 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 placeholder-gray-600 focus:outline-none focus:border-blue-500"
          />
          <p class="text-xs text-gray-600 mt-1">Full path to the xorriso executable</p>
        </div>

        <div>
          <label class="block text-sm text-gray-300 mb-1">Device poll interval (seconds)</label>
          <input
            v-model.number="settings.devicePollInterval"
            type="number"
            min="1"
            max="60"
            class="w-32 px-3 py-2 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 focus:outline-none focus:border-blue-500"
          />
          <p class="text-xs text-gray-600 mt-1">How often to check for device and media changes</p>
        </div>
      </section>

      <!-- Default ISO Options -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wide">Default ISO Options</h2>

        <div class="space-y-2">
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.rockRidge" class="accent-blue-500" />
            Rock Ridge extensions (Linux/Unix file attributes)
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.joliet" class="accent-blue-500" />
            Joliet extensions (Windows long filenames)
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.md5" class="accent-blue-500" />
            Record MD5 checksums
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.backupMode" class="accent-blue-500" />
            Backup mode (preserves ACL, xattr, hardlinks)
          </label>
        </div>
      </section>

      <!-- Default Burn Options -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wide">Default Burn Options</h2>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-gray-300 mb-1">Default speed</label>
            <select
              v-model="settings.defaultBurnOptions.speed"
              class="w-full px-2 py-1.5 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="auto">Auto</option>
              <option value="1">1x</option>
              <option value="2">2x</option>
              <option value="4">4x</option>
              <option value="8">8x</option>
              <option value="12">12x</option>
              <option value="16">16x</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-gray-300 mb-1">Default burn mode</label>
            <select
              v-model="settings.defaultBurnOptions.burnMode"
              class="w-full px-2 py-1.5 text-sm bg-gray-800 border border-gray-600 rounded text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="auto">Auto (DAO/SAO)</option>
              <option value="tao">TAO</option>
              <option value="sao">SAO/DAO</option>
            </select>
          </div>
        </div>

        <div class="space-y-2 mt-2">
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.verify" class="accent-blue-500" />
            Verify after burn
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.eject" class="accent-blue-500" />
            Eject disc when done
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.dummyMode" class="accent-yellow-500" />
            Simulation (dummy) mode
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.closeDisc" class="accent-blue-500" />
            Close disc (no multisession)
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.streamRecording" class="accent-blue-500" />
            Stream recording (for Blu-ray)
          </label>
        </div>
      </section>

      <!-- BDXL Safe Mode -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-400 uppercase tracking-wide">Blu-ray / BDXL</h2>

        <label class="flex items-start gap-3 text-sm text-gray-300">
          <input type="checkbox" v-model="settings.bdxlSafeMode" class="accent-cyan-500 mt-0.5" />
          <div>
            <span class="font-medium">BDXL safe mode</span>
            <p class="text-xs text-gray-500 mt-0.5">
              Automatically enables MD5 checksums, verify after burn, and stream recording
              for BD media. Recommended for large or archival burns.
            </p>
          </div>
        </label>
      </section>

      <!-- Save -->
      <div class="flex items-center gap-3 pt-4 border-t border-gray-700">
        <button
          @click="saveSettings"
          :disabled="saving"
          class="px-6 py-2 text-sm font-semibold rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-50 transition-colors"
        >
          {{ saving ? 'Saving...' : 'Save Settings' }}
        </button>
        <span v-if="saved" class="text-sm text-green-400">Settings saved successfully.</span>
      </div>
    </div>
  </div>
</template>
