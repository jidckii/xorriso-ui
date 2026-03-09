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

export const Browser = {
  OpenURL(url) {},
}

export function Call(...args) {
  return Promise.resolve(null)
}

export function CancellablePromise(fn) {
  return new Promise(fn)
}

function passthrough(source) { return source }
passthrough.Nullable = function(creator) { return function(source) { return source == null ? null : creator(source) } }
passthrough.Array = function(creator) { return function(source) { return Array.isArray(source) ? source.map(creator) : [] } }
passthrough.Map = function(keyCreator, valueCreator) { return function(source) { return source || {} } }
passthrough.Any = passthrough

export function Create(obj) {
  return obj
}

Create.Nullable = passthrough.Nullable
Create.Array = passthrough.Array
Create.Map = passthrough.Map
Create.Any = passthrough
