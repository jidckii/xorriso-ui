import FilePropertiesMedia from './FilePropertiesMedia.vue'

export default {
  title: 'Project/FilePropertiesMedia',
  component: FilePropertiesMedia,
  render: (args) => ({
    components: { FilePropertiesMedia },
    setup() {
      return { args }
    },
    template: '<div style="width: 400px;" class="p-4"><FilePropertiesMedia v-bind="args" /></div>',
  }),
}

export const Video = {
  args: {
    properties: {
      videoCodec: 'H.265',
      audioCodec: 'AAC',
      frameRate: '30fps',
      videoBitrate: '8 Mbps',
      sampleRate: '48000 Hz',
      audioBitrate: '256 kbps',
      audioChannels: 'Stereo',
    },
  },
}

export const Audio = {
  args: {
    properties: {
      audioCodec: 'FLAC',
      sampleRate: '96000 Hz',
      audioBitrate: '1411 kbps',
      audioChannels: '5.1',
    },
  },
}
