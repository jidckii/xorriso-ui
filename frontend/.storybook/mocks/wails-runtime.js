export const Events = {
  On(event, callback) { return () => {} },
  Off(event) {},
  Emit(event, ...data) {},
}

export const Dialogs = {
  async OpenFile(options) { return '' },
  async SaveFile(options) { return '' },
  async Info(options) {},
  async Warning(options) {},
  async Error(options) {},
  async Question(options) { return 'Yes' },
}

export const Window = {
  SetTitle(title) {},
}

export function Call(...args) {
  return Promise.resolve(null)
}

export function CancellablePromise(fn) {
  return new Promise(fn)
}

export function Create(obj) {
  return obj
}
