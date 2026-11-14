# Araya Business Hub - Backend Modul M1

Backend REST API untuk **Modul M1: Manajemen Pengguna dan Hak Akses** pada proyek Araya Business Hub. Modul ini menyediakan login terpusat dan token JWT yang dipakai modul lain.

## Status

| User Story | Cakupan | Status |
|---|---|---|
| M1-US1 | Login dan hak akses per divisi | Login, JWT, `/me`, logout, dan middleware autentikasi/role tersedia. Permission dan pembatasan divisi belum terintegrasi lengkap. |
| M1-US2 | Audit log aktivitas pengguna | Technical logging tersedia. Audit log permanen di database dan API untuk melihat riwayat belum tersedia. |
| M1-US3 | Tambah dan nonaktifkan akun | Create account, list account, perubahan status, dan deaktivasi tersedia. Perubahan divisi pengguna dan pembatasan akses khusus Super Admin belum diselesaikan. |
| M1-US4 | Ubah kata sandi | Endpoint dan usecase Change Password tersedia. Validasi konfirmasi kata sandi belum tersedia. |

## Teknologi

- Go (versi mengikuti `go.mod`)
- Gin sebagai HTTP framework
- PostgreSQL di Supabase, diakses lewat `pgx` dan `database/sql`
- sqlc untuk membuat kode akses database dari query SQL
- JWT HS256 (`golang-jwt/jwt/v5`) dan bcrypt (`golang.org/x/crypto`)
- Deploy ke Vercel (Go runtime)

## Struktur Folder

```
.
├── main.go                       entry point (lokal dan Vercel)
├── cmd/app/main.go               entry point lokal alternatif
├── config/                       pembacaan environment variable
├── internal/
│   ├── app/                      koneksi database
│   ├── controller/http/v1/       router dan route per modul (auth_route)
│   └── usecase/
│       ├── auth_usecase/         logika login
│       ├── audit_usecase/        pencatat audit sementara (log aplikasi)
│       ├── repository/queries/   query SQL untuk sqlc
│       ├── repository/sqlc/      kode hasil sqlc generate (jangan diedit)
│       └── utils/                bcrypt dan JWT
├── middleware/                   Authenticate dan RequireRole
├── pkg/httpserver/               amplop respons dan kode galat
├── prisma/migrations/            file migrasi SQL, sumber skema untuk sqlc
├── sqlc.yaml                     konfigurasi sqlc
└── .md/                          aturan pengembangan dan PRD modul M1
```

Folder kosong seperti `audit_route`, `user_route`, dan `user_usecase` disiapkan untuk user story berikutnya.

## Menjalankan Secara Lokal

Prasyarat: Go, akses ke database PostgreSQL (Supabase), dan sqlc bila ingin mengubah query.

```bash
git clone https://github.com/nadaaputri/ppl-m1-be.git
cd ppl-m1-be
git checkout dev
cp .env.example .env      # Windows: copy .env.example .env
go mod download
go run .
```

Server berjalan di `http://localhost:8080` (sesuai `APP_PORT`). Cek dengan `GET /health` dan `GET /ready`.

## Konfigurasi

Semua nilai dibaca dari environment variable. Untuk lokal boleh memakai file `.env` (tidak masuk repo). Di Vercel isi lewat dashboard.

| Variabel | Keterangan |
|---|---|
| `APP_NAME` | Nama layanan, dipakai sebagai `iss` pada token |
| `APP_PORT` | Port lokal, bawaan 8080. Di Vercel dipakai `PORT` |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Koneksi PostgreSQL. Untuk Supabase pakai host pooler |
| `DB_SSLMODE` | `require` untuk Supabase. Bawaan `disable` |
| `JWT_SECRET` | Rahasia penandatangan token, minimal 32 byte |
| `JWT_EXPIRES_IN` | Masa berlaku token, `8h` sesuai Panduan API Bersama |
| `GIN_MODE` | Opsional, isi `release` untuk deploy |

Aplikasi berhenti saat start bila `JWT_SECRET` kurang dari 32 byte, `JWT_EXPIRES_IN` tidak valid, atau konfigurasi database tidak lengkap. Variabel `APP_ENV`, `LOG_LEVEL`, dan `DATABASE_URL` pada `.env.example` belum dibaca kode Go.

## Database dan sqlc

Skema dibuat dari file migrasi SQL di `prisma/migrations`. Tabel yang ada: `divisions`, `roles`, `users`, `permissions`, `role_permissions`, dan `audit_logs`. Semua ID bertipe integer.

Query ditulis di `internal/usecase/repository/queries`. Setelah mengubah query atau menambah migrasi, jalankan:

```bash
sqlc generate
```

Setiap file migrasi baru harus ditambahkan ke daftar `schema` di `sqlc.yaml`.

## Endpoint

Semua endpoint modul M1 menggunakan prefix
`/api/m1/v1` sesuai Panduan API Bersama.

| Method | Alamat | Akses | Fungsi |
|---|---|---|---|
| GET | `/health` | Publik | Cek server hidup |
| GET | `/ready` | Publik | Cek koneksi database |
| POST | `/api/m1/v1/auth/login` | Publik | Login pengguna |
| POST | `/api/m1/v1/auth/sessions/logout` | Bearer token | Logout |
| GET | `/api/m1/v1/me` | Bearer token | Profil pengguna |
| POST | `/api/m1/v1/auth/change-password` | Bearer token | Mengubah kata sandi |
| GET | `/api/m1/v1/users` | Bearer token | Daftar akun |
| POST | `/api/m1/v1/users` | Bearer token | Membuat akun |
| PATCH | `/api/m1/v1/users/{id}/status` | Bearer token | Mengubah status akun |
| DELETE | `/api/m1/v1/users/{id}` | Bearer token | Menonaktifkan akun |
| GET | `/api/m1/v1/roles` | Bearer token | Daftar role |
| POST | `/api/m1/v1/roles` | Bearer token | Membuat role |
| GET | `/api/m1/v1/divisions` | Bearer token | Daftar divisi |
| POST | `/api/m1/v1/divisions` | Bearer token | Membuat divisi |

Dokumentasi OpenAPI:
- Raw specification: `GET /openapi.yaml`
- Swagger UI: `/api/m1/v1/docs`

Catatan:
- Registrasi Customer tidak dicantumkan dalam dokumentasi fitur.
- Endpoint pengelolaan akun, role, dan divisi masih memerlukan penyesuaian authorization agar sesuai keputusan Super Admin.
- Role Global merupakan target arsitektur yang akan diterapkan pada tahap pengembangan Role dan Permission.

### Contoh login

```http
POST /api/m1/v1/auth/login
Content-Type: application/json

{
  "email": "kasir@contoh.id",
  "password": "sandi-pengguna"
}
```

Respons 200:

```json
{
  "success": true,
  "data": {
    "access_token": "<jwt>",
    "token_type": "Bearer",
    "expires_in": 28800
  }
}
```

Token dikirim pada request berikutnya dengan header `Authorization: Bearer <access_token>`.

### Contoh profil

```json
{
  "success": true,
  "data": {
    "user_id": 42,
    "name": "Rani Kusumawati",
    "role": "kasir",
    "division": "penjualan",
    "unit_usaha_id": 3
  }
}
```

`division` dan `unit_usaha_id` bernilai `null` bila akun tidak memilikinya.

## Format Respons

```json
{ "success": true, "data": { } }
```

```json
{ "success": false, "error": { "code": "KODE_GALAT", "message": "Pesan" } }
```

Galat validasi (422) menambah `error.fields` berisi nama field dan pesannya.

| Status | `error.code` | Kapan |
|---|---|---|
| 400 | `BAD_REQUEST` | Body kosong, JSON rusak, atau tipe data salah |
| 401 | `INVALID_CREDENTIALS` | Login gagal, apa pun alasannya |
| 401 | `UNAUTHENTICATED` | Token tidak ada, salah, atau kedaluwarsa |
| 403 | `FORBIDDEN` | Role tidak berhak |
| 422 | `VALIDATION_ERROR` | Isian tidak lolos validasi |
| 500 | `INTERNAL_ERROR` | Kesalahan server, detail hanya di log |

## Token JWT

Token berlaku 8 jam dan berisi `user_id`, `name`, `role`, `division`, `unit_usaha_id`, dan `exp`, ditambah `iss` dan `iat`. Nama role dan divisi ditulis huruf kecil. Modul lain membaca `role` dan `unit_usaha_id` dari token tanpa bertanya ulang ke M1.

## Deploy ke Vercel

1. Hubungkan repositori ke project Vercel.
2. Di Settings, Build and Deployment, pilih **Framework Preset: Go**.
3. Isi environment variable di atas, untuk Production dan Preview.
4. Pastikan branch yang berisi kode adalah branch produksi (merge `dev` ke `main` bila Production memakai `main`).
5. Setelah deploy, cek `GET /health` dan `GET /ready`.

Vercel mendeteksi `main.go` di root dan menjalankan server pada port dari `PORT`.

## Keterbatasan yang Diketahui

- Token JWT tidak disimpan di server. Logout hanya berlaku di sisi klien dan token tetap berlaku sampai masa kedaluwarsa.
- Pembatasan akses masih menggunakan role pada JWT. Permission dari database belum diintegrasikan ke middleware authorization.
- Middleware RequireDivision sudah tersedia, tetapi belum diterapkan pada seluruh endpoint yang memerlukan pembatasan divisi.
- Audit aktivitas masih menggunakan technical logging dan belum seluruhnya disimpan secara permanen ke tabel audit_logs.
- Belum tersedia fitur perubahan divisi pengguna.
- Change Password belum memiliki validasi konfirmasi kata sandi.
- Autentikasi server-to-server menggunakan X-Service-Key belum diimplementasikan.
- Provisioning akun Partner dari M3 belum tersedia.
- Pengujian otomatis belum tersedia secara lengkap.

Dokumentasi OpenAPI tersedia pada `docs/openapi.yaml`,
dapat diakses melalui `/openapi.yaml`,
dan ditampilkan melalui `/api/m1/v1/docs`.

## Acuan

- Panduan API Bersama v1.0 (draf), 28 Oktober 2026
- `.md/prd-m1.md` dan `.md/aturan.md`
- Studi Kasus Capstone Project Kelas A, M1-US1 sampai M1-US4