# Fullstack Engineer Assessment - Task Management

Implementasi lengkap berdasarkan brief assessment yang diberikan:

- Backend: Go + Gin + MySQL + Redis
- Frontend: React Native + TypeScript
- Fokus: filtering, pagination, update, soft delete, cache 60 detik, cache invalidation, loading state, edit modal, bug fixes, unit/component tests, migration, README.

> Catatan: saya membuat frontend dengan Expo/React Native karena ini paling sederhana untuk menjalankan React Native ketika masih belajar. Ini tetap React Native + TypeScript dan tidak mengubah kebutuhan fitur assessment.

## Struktur

```text
task-management-assessment/
├── backend/
│   ├── cmd/
│   │   ├── migrate/
│   │   ├── seed/
│   │   └── server/
│   ├── internal/
│   │   ├── cache/
│   │   ├── config/
│   │   ├── database/
│   │   ├── httputil/
│   │   ├── models/
│   │   └── task/
│   ├── migrations/
│   │   └── 001_init.sql
│   ├── docker-compose.yml
│   ├── .env.example
│   ├── go.mod
│   └── README.md
├── mobile/
│   ├── src/
│   │   ├── components/
│   │   ├── screens/
│   │   ├── services/
│   │   └── types/
│   ├── __tests__/
│   ├── .env.example
│   ├── App.tsx
│   ├── app.json
│   ├── package.json
│   └── README.md
├── CHECKLIST.md
└── README.md
```

# A. Install semua dari nol

Instruksi ini ditulis untuk macOS karena environment yang kamu gunakan adalah Mac. Untuk Windows/Linux konsepnya sama, tetapi command instalasi paket berbeda.

## 1. Install Homebrew

Cek dulu:

```\bash
brew --version
```

Kalau belum ada, install dari halaman resmi Homebrew.

## 2. Install Git

```bash
brew install git
```

Cek:

```bash
git --version
```

## 3. Install Go

Assessment ini menggunakan Go 1.27.x. Go 1.27.1 adalah stable release yang tercantum di halaman download resmi Go.

```bash
brew install go
```

Cek:

```bash
go version
```

Kalau `go version` masih menampilkan versi lama, cek:

```bash
which go
```

Lalu sesuaikan `PATH` sesuai instalasi Homebrew kamu.

## 4. Install Node.js

Frontend menggunakan Expo SDK 57. Dokumentasi Expo saat ini mencantumkan Node.js minimum 22.13.x untuk SDK 57.

Dengan Homebrew:

```bash
brew install node@22
```

Tambahkan ke PATH jika perlu:

```bash
echo 'export PATH="/opt/homebrew/opt/node@22/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Cek:

```bash
node -v
npm -v
```

## 5. Install Docker Desktop

Docker dipakai hanya untuk menjalankan MySQL dan Redis supaya kamu tidak perlu meng-install database tersebut secara manual.

Install Docker Desktop, buka aplikasinya, lalu cek:

```bash
docker --version
docker compose version
```

## 6. Siapkan project

Masuk ke folder hasil extract project:

```bash
cd task-management-assessment
```

---

# B. Jalankan Backend

Buka Terminal 1.

## 1. Masuk backend

```bash
cd backend
```

## 2. Buat `.env`

```bash
cp .env.example .env
```

Nilai default sudah cocok dengan Docker Compose.

## 3. Jalankan MySQL + Redis

```bash
docker compose up -d
```

Cek:

```bash
docker compose ps
```

Harus ada container:

```text
task-assessment-mysql
task-assessment-redis
```

## 4. Download dependency Go

```bash
go mod tidy
```

Atau:

```bash
go mod download
```

## 5. Jalankan migration

```bash
go run ./cmd/migrate
```

Output yang diharapkan:

```text
Database migration completed.
```

## 6. Masukkan data contoh

```bash
go run ./cmd/seed
```

Output:

```text
Seed data inserted.
```

## 7. Jalankan API

```bash
go run ./cmd/server
```

Output kurang lebih:

```text
API listening on http://localhost:8080
```

Jangan tutup Terminal 1.

---

# C. Tes Backend dari Terminal 2

Buka terminal baru.

## 1. Health check

```bash
curl http://localhost:8080/health
```

Expected:

```json
{"status":"ok"}
```

## 2. Search

```bash
curl "http://localhost:8080/api/tasks?page=1&limit=5&keyword=login"
```

## 3. Filter status

```bash
curl "http://localhost:8080/api/tasks?status=done"
```

## 4. Filter assignee

```bash
curl "http://localhost:8080/api/tasks?assignee=Betran"
```

## 5. Sorting

```bash
curl "http://localhost:8080/api/tasks?sort=title_asc"
```

Sort yang diterima:

```text
created_at_asc
created_at_desc
title_asc
title_desc
status_asc
status_desc
```

## 6. Cek cache

Request pertama:

```bash
curl -i "http://localhost:8080/api/tasks?page=1&limit=5"
```

Cari header:

```text
X-Cache: MISS
```

Request yang sama lagi:

```bash
curl -i "http://localhost:8080/api/tasks?page=1&limit=5"
```

Expected:

```text
X-Cache: HIT
```

Cache TTL adalah 60 detik.

## 7. Update

```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "title":"Fix login validation",
    "description":"Updated from assessment",
    "status":"done",
    "assignee":"Betran"
  }'
```

Setelah update, request list yang sebelumnya sudah `HIT` harus berubah menjadi `MISS` karena cache di-invalidate.

## 8. Duplicate title -> HTTP 409

Coba buat task dengan title yang sudah ada:

```bash
curl -i -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{
    "title":"Fix login validation",
    "description":"Duplicate test",
    "status":"todo",
    "assignee":"Betran"
  }'
```

Expected:

```text
HTTP/1.1 409 Conflict
```

Body:

```json
{
  "error": {
    "code": "DUPLICATE_TITLE",
    "message": "a task with this title already exists"
  }
}
```

## 9. Soft delete

```bash
curl -i -X DELETE http://localhost:8080/api/tasks/1
```

Task tidak dihapus secara fisik. Kolom `deleted_at` akan terisi.

Setelah itu:

```bash
curl "http://localhost:8080/api/tasks?keyword=Fix%20login%20validation"
```

Task tersebut tidak boleh muncul.

---

# D. Jalankan Backend Tests

Masih dari folder `backend`:

```bash
go test ./...
```

Test yang disediakan mencakup:

1. Update task
2. Search dengan keyword
3. Cache invalidation

---

# E. Jalankan Frontend React Native

Buka Terminal 3.

## 1. Masuk mobile

```bash
cd mobile
```

## 2. Install dependency

```bash
npm install
```

## 3. Konfigurasi API

Untuk iOS Simulator / Android Emulator, default:

```text
EXPO_PUBLIC_API_URL=http://localhost:8080/api
```

Buat file:

```bash
cp .env.example .env
```

### Penting untuk HP fisik

`localhost` pada HP menunjuk ke HP itu sendiri, bukan ke Mac.

Cari IP Mac:

```bash
ipconfig getifaddr en0
```

Contoh hasil:

```text
192.168.1.20
```

Ubah `.env`:

```text
EXPO_PUBLIC_API_URL=http://192.168.1.20:8080/api
```

Pastikan HP dan Mac berada di Wi-Fi/LAN yang sama.

## 4. Jalankan Expo

```bash
npm start
```

Untuk simulator:

```text
i -> iOS Simulator
```

atau:

```text
a -> Android Emulator
```

Karena Expo Go pada September 2026 memiliki perubahan login, kamu mungkin diminta login Expo saat memakai Expo Go. Alternatifnya adalah memakai simulator/emulator lokal.

## 5. Jalankan test frontend

```bash
npm test
```

Test minimal yang disediakan:

```text
TaskSearchBar.test.tsx
```

Test memverifikasi input pencarian meneruskan keyword yang diketik user ke callback.

## 6. TypeScript check

```bash
npm run typecheck
```

---

# F. Cara memahami project ini sebagai pemula

Urutan belajar yang saya sarankan:

## Backend

Mulai dari:

```text
cmd/server/main.go
        ↓
internal/task/handler.go
        ↓
internal/task/service.go
        ↓
internal/task/repository.go
        ↓
MySQL
```

Redis masuk di layer service:

```text
GET /api/tasks
       ↓
Service cek Redis
       ↓
HIT  -> langsung return
MISS -> Repository -> MySQL -> simpan Redis -> return
```

Saat update/create/delete:

```text
mutation
   ↓
MySQL berhasil
   ↓
Redis tasks:list:* dihapus
```

## Frontend

Alurnya:

```text
App.tsx
   ↓
TaskListScreen
   ├── TaskSearchBar
   ├── StatusFilter
   ├── TaskCard
   └── TaskEditModal
             ↓
        src/services/api.ts
             ↓
      Go API http://.../api
```

---

# G. Checklist assessment

Lihat juga file `CHECKLIST.md` untuk checklist lengkap dengan status.

Secara implementasi:

- Backend filtering: selesai
- PUT `/api/tasks/:id`: selesai
- Soft DELETE: selesai
- Consistent errors: selesai
- Redis 60 sec: selesai
- Cache key query parameters: selesai
- Cache invalidation create/update/delete: selesai
- Search input: selesai
- Status filter: selesai
- Pagination: selesai
- Edit modal: selesai
- Loading state: selesai
- Duplicate title 409: selesai
- Refresh after update: selesai
- Hide soft-deleted tasks: selesai
- Backend tests: selesai
- Frontend component test: selesai
- DB migration: selesai
- README: selesai

---

# H. Git repository untuk dikumpulkan

Dari root project:

```bash
git init
git add .
git commit -m "feat: complete fullstack assessment"
```

Buat repository GitHub lalu:

```bash
git branch -M main
git remote add origin https://github.com/USERNAME/task-management-assessment.git
git push -u origin main
```

Jangan commit file `.env` karena berisi credential lokal.

# I. Catatan implementasi

Untuk duplicate title, migration menggunakan generated column `active_title` yang hanya berisi title ketika `deleted_at IS NULL`, kemudian diberi unique index. Dengan begitu:

- task aktif dengan title yang sama -> ditolak
- task yang sudah soft-deleted -> tidak menghalangi title tersebut dipakai lagi

Untuk sort, SQL tidak menerima user input mentah. Nilai `sort` dipetakan melalui whitelist agar tidak menjadi SQL injection.
