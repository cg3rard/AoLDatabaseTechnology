import { useEffect, useId, useState } from 'react'
import { api } from '../lib/api'
import { defaultInput, fromInput, toInput } from '../lib/format'
import { visibleFields } from '../lib/resources'
import Modal from './Modal'
import { btnGhost, btnPrimary, errorCls, inputCls } from './ui'

// Dropdown yang pilihannya diambil dari API (mis. daftar unit).
// Kalau gagal diambil (mis. tidak punya izin) atau kosong, jatuh ke input angka biasa.
function RefSelect({ field, id, value, onChange, required }) {
  const [options, setOptions] = useState(null)

  useEffect(() => {
    let cancelled = false
    api(`/api/${field.ref}`, { params: { limit: 500 } })
      .then((rows) => {
        if (cancelled) return
        setOptions(rows.map((r) => ({ value: String(r.id), label: `${r[field.refLabel]} (#${r.id})` })))
      })
      .catch(() => !cancelled && setOptions([]))
    return () => {
      cancelled = true
    }
  }, [field.ref, field.refLabel])

  if (options === null) return <input id={id} className={inputCls} disabled placeholder="Memuat…" />
  if (options.length === 0) {
    return <input id={id} className={inputCls} type="number" min="1" required={required} value={value} onChange={(e) => onChange(e.target.value)} />
  }
  return (
    <select id={id} className={inputCls} required={required} value={value} onChange={(e) => onChange(e.target.value)}>
      <option value="">{field.nullable ? '— kosong —' : 'Pilih…'}</option>
      {options.map((o) => (
        <option key={o.value} value={o.value}>
          {o.label}
        </option>
      ))}
    </select>
  )
}

function FieldInput({ field, id, value, onChange, required }) {
  const common = { id, className: inputCls, required, value }
  const set = (e) => onChange(e.target.value)

  if (field.ref) return <RefSelect field={field} id={id} value={value} onChange={onChange} required={required} />

  switch (field.type) {
    case 'bool':
      return (
        <input id={id} type="checkbox" className="h-5 w-5 rounded border-slate-300" checked={value} onChange={(e) => onChange(e.target.checked)} />
      )
    case 'textarea':
      return <textarea {...common} rows={3} onChange={set} />
    case 'select':
      return (
        <select {...common} onChange={set}>
          {field.nullable && <option value="">— kosong —</option>}
          {field.options.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      )
    case 'number':
      return <input {...common} type="number" step="any" onChange={set} />
    case 'id':
      return <input {...common} type="number" min="1" step="1" onChange={set} />
    case 'password':
      return <input {...common} type="password" autoComplete="new-password" minLength={8} onChange={set} />
    case 'datetime':
      return <input {...common} type="datetime-local" onChange={set} />
    case 'email':
    case 'date':
    case 'time':
      return <input {...common} type={field.type} onChange={set} />
    default:
      return <input {...common} type="text" onChange={set} />
  }
}

export default function FormModal({ resource, user, item, onClose, onSaved }) {
  const formId = useId()
  const editing = Boolean(item)
  const fields = visibleFields(resource, user.role)

  const [values, setValues] = useState(() =>
    Object.fromEntries(fields.map((f) => [f.name, editing ? toInput(f, item[f.name]) : defaultInput(f)])),
  )
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  async function submit(e) {
    e.preventDefault()
    setSaving(true)
    setError('')

    const payload = {}
    for (const f of fields) {
      if (f.type === 'password' && editing && !values[f.name]) continue // kosong = tidak diganti
      payload[f.name] = fromInput(f, values[f.name])
    }

    try {
      await api(editing ? `${resource.path}/${item.id}` : resource.path, {
        method: editing ? 'PUT' : 'POST',
        body: payload,
      })
      onSaved()
    } catch (err) {
      setError(err.message)
      setSaving(false)
    }
  }

  return (
    <Modal title={`${editing ? 'Ubah' : 'Tambah'} ${resource.label}`} onClose={onClose}>
      <form onSubmit={submit} className="space-y-4">
        {fields.map((f) => {
          const required = f.required === true || (f.required === 'create' && !editing)
          const id = `${formId}-${f.name}`
          return (
            <div key={f.name} className={f.type === 'bool' ? 'flex items-center gap-3' : ''}>
              <label htmlFor={id} className="mb-1 block text-sm font-medium text-slate-700">
                {f.label}
                {required && <span className="text-red-600"> *</span>}
              </label>
              <FieldInput
                field={f}
                id={id}
                value={values[f.name]}
                required={required}
                onChange={(v) => setValues((prev) => ({ ...prev, [f.name]: v }))}
              />
              {f.hint && <p className="mt-1 text-xs text-slate-500">{f.hint}</p>}
            </div>
          )
        })}

        {error && (
          <p role="alert" className={errorCls}>
            {error}
          </p>
        )}

        <div className="flex justify-end gap-2 pt-2">
          <button type="button" className={btnGhost} onClick={onClose}>
            Batal
          </button>
          <button type="submit" className={btnPrimary} disabled={saving}>
            {saving ? 'Menyimpan…' : 'Simpan'}
          </button>
        </div>
      </form>
    </Modal>
  )
}
