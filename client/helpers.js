export function formatDate(d) {
  return d.toISOString().replace(/\.\d+/, '')
}

export function parseDate(s) {
  return new Date(Date.parse(s))
}

export function entryId(entry) {
  return formatDate(entry.time) + '~' + entry.pos
}

export function findEntriesAt(entries, time) {
  let timeStr = time.toISOString()
  var matched = []
  for (let i = 0; i < entries.length; i++) {
    let entry = entries[i]
    if (entry.time.toISOString() === timeStr) {
      matched.push(entry)
    }
  }
  return matched
}
