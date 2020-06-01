import readable from 'svelte/store'

const state = readable({}, set => {
  const es = new window.EventSource('/~~~/state')
  es.onerror = e => console.log('accountstore sse error', e.data)
  es.addEventListener('state', e => {
    let data = JSON.parse(e.data)
    set(data)
  })
})

export default state
