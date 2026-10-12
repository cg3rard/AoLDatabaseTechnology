// Konfigurasi semua resource (tabel). Satu halaman generik (ResourcePage) dibangun dari sini.
// Daftar peran (read/write/residentCreate) HARUS sama dengan backend/routes/routes.go;
// backend tetap yang menentukan izin sebenarnya, ini hanya untuk menyembunyikan tombol/menu.
//
// Field:
//   name      nama kolom JSON
//   type      text | textarea | email | password | number | id | select | date | time | datetime | bool
//   nullable  boleh dikosongkan (dikirim null)
//   required  true, atau 'create' (wajib hanya saat membuat)
//   staffOnly disembunyikan dari form penghuni (backend mengisinya otomatis)
//   ref       ambil pilihan dari /api/<ref> (label dari refLabel), jatuh ke input angka bila gagal
//   money     tampil sebagai Rupiah
//   filter    tampil sebagai filter dropdown di atas tabel

import { humanize } from './format'

export const ROLE = {
  admin: 'admin',
  resident: 'resident',
  laundry: 'laundry_staff',
  security: 'security',
  technician: 'technician',
}

export const ROLE_LABEL = {
  admin: 'Admin',
  resident: 'Penghuni',
  laundry_staff: 'Staf Laundry',
  security: 'Keamanan',
  technician: 'Teknisi',
}

const { admin, resident, laundry, security, technician } = ROLE
const STAFF = [admin, security, laundry, technician]

const opts = (...values) => values.map((value) => ({ value, label: humanize(value) }))

const unitRef = { type: 'id', ref: 'units', refLabel: 'unit_number' }

export const RESOURCES = [
  // ---------------- Umum ----------------
  {
    key: 'announcements',
    label: 'Pengumuman',
    group: 'Umum',
    path: '/api/announcements',
    read: null,
    write: [admin],
    columns: ['id', 'title', 'category', 'is_published', 'created_at'],
    fields: [
      { name: 'title', label: 'Judul', type: 'text', required: true },
      { name: 'content', label: 'Isi', type: 'textarea', required: true },
      { name: 'category', label: 'Kategori', type: 'select', options: opts('general', 'maintenance', 'event', 'emergency'), filter: true },
      { name: 'is_published', label: 'Dipublikasikan', type: 'bool', default: true },
    ],
  },
  {
    key: 'lost-and-found',
    label: 'Barang Hilang & Temuan',
    group: 'Umum',
    path: '/api/lost-and-found',
    read: null,
    write: [admin, security],
    residentCreate: true,
    columns: ['id', 'type', 'item_name', 'category', 'location_found', 'status', 'created_at'],
    fields: [
      { name: 'type', label: 'Jenis', type: 'select', options: opts('lost', 'found'), filter: true },
      { name: 'item_name', label: 'Nama barang', type: 'text', required: true },
      { name: 'category', label: 'Kategori', type: 'select', options: opts('electronics', 'wallet_card', 'keys', 'accessories', 'others') },
      { name: 'location_found', label: 'Lokasi', type: 'text', required: true },
      { name: 'photo_url', label: 'URL foto', type: 'text', nullable: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('open', 'claimed', 'closed'), staffOnly: true, filter: true },
      { name: 'claimed_by', label: 'Diklaim oleh (ID user)', type: 'id', nullable: true, staffOnly: true },
    ],
  },
  {
    key: 'emergency-alerts',
    label: 'Darurat',
    group: 'Umum',
    path: '/api/emergency-alerts',
    read: [admin, security, resident],
    write: [admin, security],
    residentCreate: true,
    columns: ['id', 'unit_id', 'alert_type', 'status', 'triggered_by', 'created_at', 'resolved_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'alert_type', label: 'Jenis darurat', type: 'select', options: opts('medical', 'fire', 'security_threat') },
      { name: 'status', label: 'Status', type: 'select', options: opts('triggered', 'responding', 'resolved', 'false_alarm'), staffOnly: true, filter: true },
    ],
  },

  // ---------------- Hunian ----------------
  {
    key: 'units',
    label: 'Unit',
    group: 'Hunian',
    path: '/api/units',
    read: STAFF,
    write: [admin],
    columns: ['id', 'unit_number', 'tower_block', 'floor', 'status'],
    fields: [
      { name: 'unit_number', label: 'Nomor unit', type: 'text', required: true },
      { name: 'tower_block', label: 'Tower / blok', type: 'text', required: true },
      { name: 'floor', label: 'Lantai', type: 'number', required: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('vacant', 'occupied', 'under_renovation'), filter: true },
    ],
  },
  {
    key: 'users',
    label: 'Pengguna',
    group: 'Hunian',
    path: '/api/users',
    read: [admin],
    write: [admin],
    columns: ['id', 'full_name', 'email', 'role', 'unit_id', 'phone_number'],
    fields: [
      { name: 'full_name', label: 'Nama lengkap', type: 'text', required: true },
      { name: 'email', label: 'Email', type: 'email', required: true },
      { name: 'password', label: 'Password', type: 'password', required: 'create', hint: 'Minimal 8 karakter. Saat mengubah, kosongkan jika tidak diganti.' },
      { name: 'phone_number', label: 'No. telepon', type: 'text', required: true },
      {
        name: 'role',
        label: 'Peran',
        type: 'select',
        options: Object.entries(ROLE_LABEL).map(([value, label]) => ({ value, label })),
        filter: true,
      },
      { name: 'resident_type', label: 'Tipe penghuni', type: 'select', options: opts('owner', 'tenant'), nullable: true },
      { name: 'unit_id', label: 'Unit', ...unitRef, nullable: true, hint: 'Wajib untuk peran Penghuni.' },
    ],
  },
  {
    key: 'vehicles',
    label: 'Kendaraan',
    group: 'Hunian',
    path: '/api/vehicles',
    read: [admin, security, resident],
    write: [admin, security],
    residentCreate: true,
    columns: ['id', 'unit_id', 'plate_number', 'vehicle_type', 'brand_model', 'parking_slot_id', 'status'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'plate_number', label: 'Plat nomor', type: 'text', required: true },
      { name: 'vehicle_type', label: 'Jenis', type: 'select', options: opts('car', 'motorcycle'), filter: true },
      { name: 'brand_model', label: 'Merek / model', type: 'text', nullable: true },
      { name: 'parking_slot_id', label: 'Slot parkir', type: 'id', ref: 'parking-slots', refLabel: 'slot_code', nullable: true, staffOnly: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('active', 'inactive'), staffOnly: true, filter: true },
    ],
  },
  {
    key: 'parking-slots',
    label: 'Slot Parkir',
    group: 'Hunian',
    path: '/api/parking-slots',
    read: [admin, security, resident],
    write: [admin],
    columns: ['id', 'slot_code', 'zone', 'vehicle_type', 'is_occupied'],
    fields: [
      { name: 'slot_code', label: 'Kode slot', type: 'text', required: true },
      { name: 'zone', label: 'Zona', type: 'text', required: true },
      { name: 'vehicle_type', label: 'Jenis kendaraan', type: 'select', options: opts('car', 'motorcycle'), filter: true },
      { name: 'is_occupied', label: 'Terisi', type: 'bool', default: false },
    ],
  },
  {
    key: 'visitor-passes',
    label: 'Pas Tamu',
    group: 'Hunian',
    path: '/api/visitor-passes',
    read: [admin, security, resident],
    write: [admin, security],
    residentCreate: true,
    columns: ['id', 'unit_id', 'visitor_name', 'pass_code', 'valid_date', 'status'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'visitor_name', label: 'Nama tamu', type: 'text', required: true },
      { name: 'visitor_phone', label: 'No. telepon tamu', type: 'text', nullable: true },
      { name: 'vehicle_number', label: 'Plat kendaraan', type: 'text', nullable: true },
      { name: 'valid_date', label: 'Tanggal berlaku', type: 'date', required: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('active', 'checked_in', 'completed', 'expired'), staffOnly: true, filter: true },
      { name: 'check_in_time', label: 'Waktu masuk', type: 'datetime', nullable: true, staffOnly: true },
      { name: 'check_out_time', label: 'Waktu keluar', type: 'datetime', nullable: true, staffOnly: true },
    ],
  },

  // ---------------- Akses ----------------
  {
    key: 'access-cards',
    label: 'Kartu Akses',
    group: 'Akses',
    path: '/api/access-cards',
    read: [admin, security, resident],
    write: [admin, security],
    columns: ['id', 'unit_id', 'card_label', 'holder_name', 'rfid_uid', 'authorized_floor', 'status'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true },
      { name: 'rfid_uid', label: 'RFID UID', type: 'text', required: true },
      { name: 'card_label', label: 'Label kartu', type: 'text', required: true },
      { name: 'holder_name', label: 'Nama pemegang', type: 'text', required: true },
      { name: 'authorized_floor', label: 'Lantai yang diizinkan', type: 'number', required: true },
      { name: 'has_facility_access', label: 'Akses fasilitas', type: 'bool', default: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('active', 'blocked', 'lost', 'deactivated'), filter: true },
      { name: 'block_reason', label: 'Alasan blokir', type: 'textarea', nullable: true },
    ],
  },
  {
    key: 'card-requests',
    label: 'Permohonan Kartu',
    group: 'Akses',
    path: '/api/card-requests',
    read: [admin, security, resident],
    write: [admin, security],
    residentCreate: true,
    columns: ['id', 'unit_id', 'request_type', 'holder_name', 'fee_amount', 'status', 'created_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'request_type', label: 'Jenis permohonan', type: 'select', options: opts('new_additional', 'lost_replacement', 'damaged_replacement'), filter: true },
      { name: 'lost_card_id', label: 'ID kartu yang hilang/rusak', type: 'id', nullable: true },
      { name: 'holder_name', label: 'Nama pemegang', type: 'text', required: true },
      { name: 'reason', label: 'Alasan', type: 'textarea', nullable: true },
      { name: 'fee_amount', label: 'Biaya', type: 'number', money: true, default: 0, staffOnly: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('pending', 'approved', 'ready_for_pickup', 'completed', 'rejected'), staffOnly: true, filter: true },
      { name: 'processed_by', label: 'Diproses oleh (ID user)', type: 'id', nullable: true, staffOnly: true },
    ],
  },

  // ---------------- Keuangan ----------------
  {
    key: 'invoices',
    label: 'Invoice',
    group: 'Keuangan',
    path: '/api/invoices',
    read: [admin, resident],
    write: [admin],
    columns: ['id', 'invoice_number', 'unit_id', 'invoice_type', 'total_amount', 'payment_status', 'due_date'],
    fields: [
      { name: 'invoice_number', label: 'Nomor invoice', type: 'text', nullable: true, hint: 'Kosongkan untuk dibuat otomatis.' },
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true },
      { name: 'invoice_type', label: 'Jenis', type: 'select', options: opts('ipl', 'water', 'electricity_token', 'laundry', 'card_replacement', 'permit_deposit'), filter: true },
      { name: 'total_amount', label: 'Total', type: 'number', money: true, required: true },
      { name: 'payment_status', label: 'Status bayar', type: 'select', options: opts('unpaid', 'pending_verification', 'paid', 'cancelled'), filter: true },
      { name: 'payment_proof_url', label: 'URL bukti bayar', type: 'text', nullable: true },
      { name: 'due_date', label: 'Jatuh tempo', type: 'date', required: true },
    ],
  },
  {
    key: 'utility-tariffs',
    label: 'Tarif Utilitas',
    group: 'Keuangan',
    path: '/api/utility-tariffs',
    read: null,
    write: [admin],
    columns: ['id', 'utility_type', 'rate_per_unit', 'admin_fee', 'updated_at'],
    fields: [
      { name: 'utility_type', label: 'Jenis', type: 'select', options: opts('electricity_token', 'water', 'ipl') },
      { name: 'rate_per_unit', label: 'Tarif per unit', type: 'number', money: true, required: true },
      { name: 'admin_fee', label: 'Biaya admin', type: 'number', money: true, default: 0 },
    ],
  },
  {
    key: 'utility-records',
    label: 'Pemakaian Utilitas',
    group: 'Keuangan',
    path: '/api/utility-records',
    read: [admin, resident],
    write: [admin],
    columns: ['id', 'unit_id', 'utility_type', 'initial_meter', 'final_meter', 'usage_total', 'reading_date'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true },
      { name: 'utility_type', label: 'Jenis', type: 'select', options: opts('water', 'electricity'), filter: true },
      { name: 'initial_meter', label: 'Meter awal', type: 'number', required: true },
      { name: 'final_meter', label: 'Meter akhir', type: 'number', required: true },
      { name: 'reading_date', label: 'Tanggal baca', type: 'date', required: true },
      { name: 'invoice_id', label: 'Invoice', type: 'id', ref: 'invoices', refLabel: 'invoice_number', nullable: true },
    ],
  },
  {
    key: 'token-purchases',
    label: 'Pembelian Token',
    group: 'Keuangan',
    path: '/api/token-purchases',
    read: [admin, resident],
    write: [admin],
    columns: ['id', 'unit_id', 'invoice_id', 'nominal', 'token_code', 'status', 'created_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true },
      { name: 'invoice_id', label: 'Invoice', type: 'id', ref: 'invoices', refLabel: 'invoice_number', required: true },
      { name: 'nominal', label: 'Nominal', type: 'number', money: true, required: true },
      { name: 'token_code', label: 'Kode token', type: 'text', nullable: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('pending', 'success', 'failed'), filter: true },
    ],
  },

  // ---------------- Layanan ----------------
  {
    key: 'packages',
    label: 'Paket',
    group: 'Layanan',
    path: '/api/packages',
    read: [admin, security, resident],
    write: [admin, security],
    columns: ['id', 'unit_id', 'recipient_name', 'courier_name', 'receipt_number', 'status', 'received_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true },
      { name: 'recipient_name', label: 'Nama penerima', type: 'text', required: true },
      { name: 'courier_name', label: 'Kurir', type: 'text', required: true },
      { name: 'receipt_number', label: 'No. resi', type: 'text', nullable: true },
      { name: 'package_photo_url', label: 'URL foto paket', type: 'text', required: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('waiting_pickup', 'collected'), filter: true },
    ],
  },
  {
    key: 'laundry-orders',
    label: 'Laundry',
    group: 'Layanan',
    path: '/api/laundry-orders',
    read: [admin, laundry, resident],
    write: [admin, laundry],
    residentCreate: true,
    columns: ['id', 'unit_id', 'service_type', 'weight_kg', 'item_count', 'total_cost', 'status', 'created_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'service_type', label: 'Layanan', type: 'select', options: opts('kiloan_cuci_setrika', 'kiloan_cuci_kering', 'satuan_dry_clean') },
      { name: 'weight_kg', label: 'Berat (kg)', type: 'number', nullable: true },
      { name: 'item_count', label: 'Jumlah item', type: 'number', nullable: true },
      { name: 'notes', label: 'Catatan', type: 'textarea', nullable: true },
      { name: 'total_cost', label: 'Total biaya', type: 'number', money: true, default: 0, staffOnly: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('order_placed', 'received', 'washing', 'ironing', 'ready_for_pickup', 'completed'), staffOnly: true, filter: true },
      { name: 'staff_id', label: 'Staf (ID user)', type: 'id', nullable: true, staffOnly: true },
      { name: 'invoice_id', label: 'Invoice (ID)', type: 'id', nullable: true, staffOnly: true },
    ],
  },
  {
    key: 'facilities',
    label: 'Fasilitas',
    group: 'Layanan',
    path: '/api/facilities',
    read: null,
    write: [admin],
    columns: ['id', 'name', 'max_capacity', 'open_time', 'close_time', 'slot_duration_minutes', 'is_active'],
    fields: [
      { name: 'name', label: 'Nama', type: 'text', required: true },
      { name: 'max_capacity', label: 'Kapasitas per slot', type: 'number', required: true },
      { name: 'open_time', label: 'Jam buka', type: 'time', required: true },
      { name: 'close_time', label: 'Jam tutup', type: 'time', required: true },
      { name: 'slot_duration_minutes', label: 'Durasi slot (menit)', type: 'number', required: true },
      { name: 'is_active', label: 'Aktif', type: 'bool', default: true },
    ],
  },
  {
    key: 'facility-bookings',
    label: 'Booking Fasilitas',
    group: 'Layanan',
    path: '/api/facility-bookings',
    read: [admin, security, resident],
    write: [admin],
    residentCreate: true,
    columns: ['id', 'facility_id', 'user_id', 'booking_date', 'start_time', 'end_time', 'status'],
    fields: [
      { name: 'facility_id', label: 'Fasilitas', type: 'id', ref: 'facilities', refLabel: 'name', required: true },
      { name: 'booking_date', label: 'Tanggal', type: 'date', required: true },
      { name: 'start_time', label: 'Jam mulai', type: 'time', required: true },
      { name: 'end_time', label: 'Jam selesai', type: 'time', required: true, hint: 'Durasi harus sama dengan durasi slot fasilitas.' },
      { name: 'user_id', label: 'Atas nama (ID user)', type: 'id', nullable: true, staffOnly: true, hint: 'Kosongkan untuk akun Anda.' },
      { name: 'status', label: 'Status', type: 'select', options: opts('confirmed', 'cancelled', 'completed'), staffOnly: true, filter: true },
    ],
  },
  {
    key: 'maintenance-tickets',
    label: 'Tiket Perbaikan',
    group: 'Layanan',
    path: '/api/maintenance-tickets',
    read: [admin, technician, resident],
    write: [admin, technician],
    residentCreate: true,
    columns: ['id', 'unit_id', 'category', 'title', 'priority', 'status', 'created_at'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, nullable: true, staffOnly: true },
      { name: 'category', label: 'Kategori', type: 'select', options: opts('plumbing', 'electrical', 'air_conditioner', 'civil_structure'), filter: true },
      { name: 'title', label: 'Judul', type: 'text', required: true },
      { name: 'description', label: 'Deskripsi', type: 'textarea', required: true },
      { name: 'photo_url', label: 'URL foto', type: 'text', nullable: true },
      { name: 'priority', label: 'Prioritas', type: 'select', options: opts('low', 'medium', 'high', 'urgent'), default: 'medium', staffOnly: true, filter: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('open', 'assigned', 'in_progress', 'resolved', 'closed'), staffOnly: true, filter: true },
      { name: 'assigned_technician_id', label: 'Teknisi (ID user)', type: 'id', nullable: true, staffOnly: true },
      { name: 'resolution_notes', label: 'Catatan penyelesaian', type: 'textarea', nullable: true, staffOnly: true },
    ],
  },
  {
    key: 'permits',
    label: 'Izin Pindahan & Renovasi',
    group: 'Layanan',
    path: '/api/permits',
    read: [admin, security, resident],
    write: [admin],
    residentCreate: true,
    columns: ['id', 'unit_id', 'permit_type', 'start_date', 'end_date', 'worker_count', 'deposit_amount', 'status'],
    fields: [
      { name: 'unit_id', label: 'Unit', ...unitRef, required: true, staffOnly: true },
      { name: 'permit_type', label: 'Jenis izin', type: 'select', options: opts('moving_in', 'moving_out', 'renovation'), filter: true },
      { name: 'start_date', label: 'Tanggal mulai', type: 'date', required: true },
      { name: 'end_date', label: 'Tanggal selesai', type: 'date', required: true },
      { name: 'contractor_name', label: 'Kontraktor', type: 'text', nullable: true },
      { name: 'worker_count', label: 'Jumlah pekerja', type: 'number', default: 0 },
      { name: 'deposit_amount', label: 'Deposit', type: 'number', money: true, default: 0, staffOnly: true },
      { name: 'deposit_invoice_id', label: 'Invoice deposit', type: 'id', ref: 'invoices', refLabel: 'invoice_number', nullable: true, staffOnly: true },
      { name: 'status', label: 'Status', type: 'select', options: opts('submitted', 'approved', 'rejected', 'in_progress', 'completed', 'deposit_refunded'), staffOnly: true, filter: true },
      { name: 'notes', label: 'Catatan', type: 'textarea', nullable: true },
    ],
  },
]

export const RESOURCE_BY_KEY = Object.fromEntries(RESOURCES.map((r) => [r.key, r]))

export const canRead = (res, role) => res.read === null || res.read.includes(role)
export const canWrite = (res, role) => res.write.includes(role)
export const canCreate = (res, role) => canWrite(res, role) || (role === resident && Boolean(res.residentCreate))

// Field yang ditampilkan di form untuk peran tertentu.
export const visibleFields = (res, role) => res.fields.filter((f) => !(f.staffOnly && role === resident))

// Menu sidebar: kelompok -> resource yang boleh dibaca peran ini.
export function navGroups(role) {
  const groups = []
  for (const res of RESOURCES) {
    if (!canRead(res, role)) continue
    let group = groups.find((g) => g.name === res.group)
    if (!group) {
      group = { name: res.group, items: [] }
      groups.push(group)
    }
    group.items.push(res)
  }
  return groups
}
