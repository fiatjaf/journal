import {readable, writable} from 'svelte/store'

import * as toast from './toast'
import {parseDate} from './helpers'

export const state = readable({}, set => {
  const es = new EventSource('/~~~/state')
  es.addEventListener('state', e => {
    set(JSON.parse(e.data))
  })
  es.addEventListener('error', e => {
    let data = JSON.parse(e.data)
    toast.error(data.error + ' error: ' + data.message)
  })

  return () => {
    es.close()
  }
})

export const entries = writable([], set => {
  fetch('/~/entries')
    .then(r => r.text())
    .then(text =>
      text
        .split('\n')
        .filter(x => x)
        .map(JSON.parse)
        .map(entry => {
          entry.time = parseDate(entry.time)
          return entry
        })
        .reduce((acc, entry) => acc.concat(entry), [])
    )
    .then(set)
    .catch(err => {
      toast.error('failed to fetch entries: ' + err.message)
    })

  return () => {}
})
