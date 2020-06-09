export function formatDate(d) {
  return d.toISOString().replace(/\.\d+/, '')
}

export function parseDate(s) {
  return new Date(Date.parse(s))
}

export function entryId(entry) {
  return formatDate(entry.time) + '~' + entry.pos
}
