/** @format */

import * as toast from './toast'
import App from './App.html'

fetch('/~/metadata')
  .then(r => r.json())
  .then(metadata => {
    for (let m in metadata.methods) {
      metadata.methods[m].splice(metadata.methods[m].indexOf('date'), 1)
    }
    window.metadata = metadata
  })
  .catch(err => {
    toast.error('failed to fetch metadata: ' + err.message)
  })
  .then(() => {
    new App({
      target: document.getElementById('app')
    })
  })
  .catch(err => {
    toast.error('failed to start app: ' + err.message)
  })
