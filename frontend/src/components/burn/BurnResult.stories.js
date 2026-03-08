import BurnResult from './BurnResult.vue'

export default {
  title: 'Burn/BurnResult',
  component: BurnResult,
  decorators: [
    () => ({
      template: '<div style="max-width: 600px;"><story /></div>',
    }),
  ],
}

const baseLogLines = [
  '[INFO] Starting burn process...',
  '[INFO] Speed: 8x',
  '[INFO] Writing track 1...',
  '[INFO] 100% complete',
]

export const Success = {
  args: {
    job: {
      result: {
        success: true,
        message: 'Burn completed successfully',
        verifyErrors: 0,
      },
    },
    logLines: [...baseLogLines, '[INFO] Burn completed successfully'],
  },
}

export const Failed = {
  args: {
    job: {
      result: {
        success: false,
        message: 'Write error at sector 12345',
      },
    },
    logLines: [
      ...baseLogLines,
      '[ERROR] Write error at sector 12345',
      '[FATAL] Burn failed',
    ],
  },
}

export const VerificationPassed = {
  args: {
    job: {
      result: {
        success: true,
        message: 'Burn and verification completed successfully',
        verifyErrors: 0,
        md5Match: true,
      },
    },
    logLines: [
      ...baseLogLines,
      '[INFO] Burn completed',
      '[INFO] Starting verification...',
      '[INFO] Verification passed, MD5 match',
    ],
  },
}

export const VerificationFailed = {
  args: {
    job: {
      result: {
        success: true,
        message: 'Burn completed, verification failed',
        verifyErrors: 3,
        md5Match: false,
      },
    },
    logLines: [
      ...baseLogLines,
      '[INFO] Burn completed',
      '[INFO] Starting verification...',
      '[ERROR] Verification failed: 3 errors found',
      '[ERROR] MD5 mismatch',
    ],
  },
}
