// Fungsi bantu untuk menampilkan data (tabel) dan mengubah nilai form <-> JSON API.

export const humanize = (v) => String(v).replace(/_/g, ' ').replace(/^./, (c) => c.toUpperCase())

const rupiah = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0,
})

export function optionLabel(field, value) {
  return field.options?.find((o) => o.value === value)?.label ?? humanize(value)
}

// Teks yang ditampilkan di sel tabel.
export function formatCell(name, field, value) {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'boolean') return value ? 'Ya' : 'Tidak'
  if (field?.money) return rupiah.format(value)
  if (field?.type === 'select') return optionLabel(field, value)
  if (field?.type === 'time') return String(value).slice(0, 5)
  if (field?.type === 'datetime' || /_at$/.test(name)) {
    const d = new Date(value)
    return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
  }
  return String(value)
}

const pad = (n) => String(n).padStart(2, '0')

// Nilai dari API -> nilai input form (selalu string, kecuali checkbox).
export function toInput(field, value) {
  if (field.type === 'bool') return Boolean(value)
  if (value === null || value === undefined) return ''
  if (field.type === 'date') return String(value).slice(0, 10)
  if (field.type === 'time') return String(value).slice(0, 5)
  if (field.type === 'datetime') {
    const d = new Date(value)
    if (Number.isNaN(d.getTime())) return ''
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
  }
  return String(value)
}

// Nilai awal saat membuat data baru.
export function defaultInput(field) {
  if (field.default !== undefined) return field.default
  if (field.type === 'bool') return false
  if (field.type === 'select' && !field.nullable) return field.options[0].value
  return ''
}

// Nilai input form -> nilai yang dikirim ke API.
export function fromInput(field, value) {
  switch (field.type) {
    case 'bool':
      return Boolean(value)
    case 'number':
    case 'id':
      if (value === '') return field.nullable ? null : 0
      return Number(value)
    case 'datetime':
      return value === '' ? null : new Date(value).toISOString()
    default: {
      const text = typeof value === 'string' ? value.trim() : value
      if (text === '') return field.nullable ? null : ''
      return text
    }
  }
}
