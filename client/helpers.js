export function today() {
  return new Date(Date.parse(new Date().toISOString().split('T')[0]))
}

export function formatDate(d) {
  return d.toISOString().replace(/\.\d+/, '')
}

export function parseDate(s) {
  return new Date(Date.parse(s))
}

export function entryId(entry) {
  return formatDate(entry.time) + '~' + entry.pos
}
