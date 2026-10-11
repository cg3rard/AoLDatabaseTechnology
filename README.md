# AoL Database Technology – Property Management

## Struktur

```
.
├── backend/            # API Go (Fiber + GORM), port 8080
├── frontend/           # React + Vite + Tailwind
├── db-init/            # Script .sql yang dijalankan otomatis saat volume MySQL pertama kali dibuat
├── docker-compose.yml  # MySQL 8.0, port 3306
└── run.bat             # Menjalankan DB, backend, dan frontend sekaligus (Windows)
```

## Menjalankan

Prasyarat: Docker, Go, Node.js.

```bat
copy .env.example .env
cd frontend && npm install && cd ..
run.bat
```

Isi kredensial di `.env` dulu. File ini tidak di-commit.

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- MySQL: localhost:3306 (database `property_management`)

Script di `db-init/` hanya dijalankan sekali. Untuk menjalankan ulang, hapus volume: `docker-compose down -v`.
