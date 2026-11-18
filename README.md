# Araya Business Hub --- Backend Modul 1 (M1)

Backend REST API untuk **Modul 1: Manajemen Pengguna dan Hak Akses**
pada platform Araya Business Hub.

M1 menyediakan layanan autentikasi terpusat, profil pengguna,
pengelolaan akun, role, katalog permission, division, dan audit log.
Modul lain menggunakan token yang diterbitkan M1 dan bertanggung jawab
memeriksa hak akses pada endpoint bisnis milik modul masing-masing.

> Dokumentasi ini disusun berdasarkan source code pada repository yang
> diperiksa. Keterangan status implementasi dibedakan dari kebutuhan PRD
> yang belum tersedia sebagai endpoint atau belum diterapkan sepenuhnya.

## Daftar Isi

-   [Ruang Lingkup](#ruang-lingkup)
-   [Status Implementasi](#status-implementasi)
-   [Teknologi](#teknologi)
-   [Struktur Proyek](#struktur-proyek)
-   [Persyaratan](#persyaratan)
-   [Menjalankan Secara Lokal](#menjalankan-secara-lokal)
-   [Konfigurasi Environment](#konfigurasi-environment)
-   [Database dan SQLC](#database-dan-sqlc)
-   [Arsitektur dan Alur Autentikasi](#arsitektur-dan-alur-autentikasi)
-   [Authorization](#authorization)
-   [Daftar Endpoint](#daftar-endpoint)
-   [Format Request dan Response](#format-request-dan-response)
-   [Audit Log](#audit-log)
-   [Pengujian](#pengujian)
-   [Dokumentasi API](#dokumentasi-api)
-   [Deployment](#deployment)
-   [Keterbatasan Implementasi Saat
    Ini](#keterbatasan-implementasi-saat-ini)

## Ruang Lingkup

Fokus backend M1:

1.  Login dan penerbitan JWT.
2.  Pengambilan profil pengguna yang sedang login.
3.  Perubahan kata sandi pengguna.
4.  Pembuatan, pembacaan, pembaruan, aktivasi, dan deaktivasi akun
    sesuai route yang tersedia.
5.  Pengelolaan role global.
6.  Pengambilan katalog permission.
7.  Penetapan permission ke role.
8.  Pengelolaan division global.
9.  Pencatatan dan pembacaan audit log.
10. Health check dan readiness check.

M1 merupakan pemilik autentikasi dan data akun platform. M1 bukan
pemilik utama data partner maupun unit usaha. `partner_id` dan
`unit_usaha_id` disimpan sebagai referensi yang berasal dari modul
pemilik datanya.

## Status Implementasi

Status berikut menggambarkan source code yang diperiksa, bukan klaim
bahwa seluruh acceptance criteria telah selesai.

  ----------------------------------------------------------------------------
  Area                    Implementasi yang tersedia   Catatan
  ----------------------- ---------------------------- -----------------------
  M1-US1 --- Login dan    Login JWT HS256, pemeriksaan Penerapan permission
  hak akses               akun aktif saat login, klaim dan division pada route
                          role/division/partner/unit   belum menyeluruh.
                          usaha, permission role dalam 
                          token, middleware            
                          autentikasi dan helper       
                          authorization                

  M1-US2 --- Audit log    Tabel audit log, pencatatan  Route audit partner
                          aktivitas penting, endpoint  memeriksa role partner
                          audit log partner dan        dan permission
                          platform, pagination, urutan `m1.audit.read`.
                          terbaru lebih dahulu

  M1-US3 --- Pengelolaan  Create, list, detail,        Route pengelolaan akun
  akun                    update, aktivasi, dan        saat ini dibatasi untuk
                          deaktivasi akun; soft        Super Admin.
                          deactivation; audit          
                          perubahan akun               

  M1-US4 --- Ubah kata    Verifikasi kata sandi lama,  Tersedia melalui
  sandi                   validasi kata sandi baru,    endpoint
                          konfirmasi, bcrypt, dan      change-password.
                          pencatatan audit             

  Role                    List/create role global,     Belum ada endpoint
                          melihat katalog permission,  update/delete/status
                          melihat dan mengganti        role.
                          permission role              

  Division                List dan create division     Belum ada endpoint
                          global                       update/delete/status
                                                       division.

  Integrasi M3--M1        Belum tersedia endpoint      Merupakan pekerjaan
                          provisioning Partner Admin   integrasi lanjutan,
                          atau penerimaan event status bukan endpoint aktif
                          partner                      dalam OpenAPI ini.
  ----------------------------------------------------------------------------

## Teknologi

-   **Go 1.26.2** --- versi yang ditetapkan pada `go.mod`.
-   **Gin** --- HTTP framework.
-   **PostgreSQL** --- database aplikasi.
-   **pgx** --- driver PostgreSQL melalui `database/sql`.
-   **sqlc** --- menghasilkan kode repository Go dari query SQL.
-   **JWT HS256** --- autentikasi token menggunakan `golang-jwt/jwt/v5`.
-   **bcrypt** --- hashing kata sandi.
-   **godotenv** --- membaca `.env` untuk pengembangan lokal.
-   **Prisma schema dan SQL migrations** --- definisi model dan riwayat
    perubahan database.
-   **Vercel Go Runtime** --- target deployment yang digunakan
    repository.

## Struktur Proyek

``` text
.
â”œâ”€â”€ main.go
â”œâ”€â”€ cmd/
â”‚   â””â”€â”€ app/
â”‚       â””â”€â”€ main.go
â”œâ”€â”€ config/
â”‚   â”œâ”€â”€ config.go
â”‚   â””â”€â”€ config_test.go
â”œâ”€â”€ docs/
â”‚   â”œâ”€â”€ embed.go
â”‚   â””â”€â”€ openapi.yaml
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ app/
â”‚   â”‚   â””â”€â”€ app.go
â”‚   â”œâ”€â”€ controller/
â”‚   â”‚   â””â”€â”€ http/
â”‚   â”‚       â””â”€â”€ v1/
â”‚   â”‚           â”œâ”€â”€ auth_route/
â”‚   â”‚           â”œâ”€â”€ user_route/
â”‚   â”‚           â”œâ”€â”€ role_route/
â”‚   â”‚           â”œâ”€â”€ division_route/
â”‚   â”‚           â”œâ”€â”€ audit_route/
â”‚   â”‚           â”œâ”€â”€ router.go
â”‚   â”‚           â””â”€â”€ mount_*.go
â”‚   â””â”€â”€ usecase/
â”‚       â”œâ”€â”€ auth_usecase/
â”‚       â”œâ”€â”€ user_usecase/
â”‚       â”œâ”€â”€ role_usecase/
â”‚       â”œâ”€â”€ division_usecase/
â”‚       â”œâ”€â”€ audit_usecase/
â”‚       â”œâ”€â”€ repository/
â”‚       â”‚   â”œâ”€â”€ queries/
â”‚       â”‚   â””â”€â”€ sqlc/
â”‚       â””â”€â”€ utils/
â”œâ”€â”€ middleware/
â”‚   â”œâ”€â”€ auth.go
â”‚   â”œâ”€â”€ cors.go
â”‚   â””â”€â”€ scope.go
â”œâ”€â”€ pkg/
â”‚   â””â”€â”€ httpserver/
â”œâ”€â”€ prisma/
â”‚   â”œâ”€â”€ migrations/
â”‚   â”œâ”€â”€ schema.prisma
â”‚   â”œâ”€â”€ seed.ts
â”‚   â””â”€â”€ seed-dummy.ts
â”œâ”€â”€ sqlc.yaml
â”œâ”€â”€ go.mod
â”œâ”€â”€ go.sum
â”œâ”€â”€ .env.example
â””â”€â”€ README.md
```

### Tanggung jawab folder

  ----------------------------------------------------------------------------
  Folder                                   Tanggung jawab
  ---------------------------------------- -----------------------------------
  `cmd/` dan `main.go`                     Entry point aplikasi.

  `config/`                                Membaca dan memvalidasi environment
                                           variable.

  `internal/app/`                          Membuat dan menguji koneksi
                                           database.

  `internal/controller/http/v1/`           Router, handler, DTO, dan mounting
                                           endpoint HTTP.

  `internal/usecase/`                      Logika bisnis autentikasi, akun,
                                           role, division, dan audit.

  `internal/usecase/repository/queries/`   Query SQL sumber untuk sqlc.

  `internal/usecase/repository/sqlc/`      Kode Go hasil generate sqlc; jangan
                                           diedit manual.

  `internal/usecase/utils/`                JWT, password, format, dan helper
                                           database.

  `middleware/`                            Autentikasi, role, permission,
                                           division, scope, dan CORS.

  `pkg/httpserver/`                        Envelope response, validasi
                                           request, pagination, dan query
                                           parameter.

  `prisma/migrations/`                     Riwayat migrasi database.

  `docs/`                                  Spesifikasi OpenAPI yang di-embed
                                           ke aplikasi.
  ----------------------------------------------------------------------------

## Persyaratan

Sebelum menjalankan aplikasi, siapkan:

-   Go versi **1.26.2** atau versi yang sesuai dengan `go.mod`.
-   PostgreSQL yang dapat diakses dari lingkungan aplikasi.
-   Kredensial database dengan hak akses yang diperlukan.
-   `sqlc` versi yang sesuai dengan konfigurasi proyek jika mengubah
    query.
-   Git untuk mengambil source code.

## Menjalankan Secara Lokal

### 1. Clone repository

``` bash
git clone https://github.com/nadaaputri/ppl-m1-be.git
cd ppl-m1-be
git checkout dev
```

### 2. Siapkan environment lokal

Linux/macOS:

``` bash
cp .env.example .env
```

Windows PowerShell:

``` powershell
Copy-Item .env.example .env
```

Isi `.env` dengan konfigurasi lokal yang benar. Jangan membagikan atau
melakukan commit file `.env`.

### 3. Unduh dependency Go

``` bash
go mod download
```

### 4. Jalankan aplikasi

``` bash
go run .
```

Secara default aplikasi mendengarkan pada port `8080`, kecuali
`APP_PORT` atau environment `PORT` menentukan port lain.

### 5. Periksa aplikasi

``` bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

Respons health check yang diharapkan:

``` json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

Respons readiness jika database dapat diakses:

``` json
{
  "success": true,
  "data": {
    "status": "ready",
    "database": "connected"
  }
}
```

## Konfigurasi Environment

`config/config.go` membaca konfigurasi berikut.

  -----------------------------------------------------------------------------------
  Variabel                                 Wajib Fungsi           Default/Aturan
  ------------------------ --------------------- ---------------- -------------------
  `APP_NAME`                    Tidak divalidasi Nama aplikasi    Tidak memiliki
                                           wajib dan issuer JWT   default eksplisit
                                                                  di loader.

  `APP_PORT`                               Tidak Port aplikasi    `8080`.
                                                 lokal            

  `PORT`                                   Tidak Port runtime     Jika terisi,
                                                 hosting; dibaca  mengesampingkan
                                                 oleh `main.go`   `APP_PORT`.

  `DB_HOST`                                   Ya Host PostgreSQL  Tidak ada default.

  `DB_PORT`                                   Ya Port PostgreSQL  Tidak ada default.

  `DB_USER`                                   Ya Username         Tidak ada default.
                                                 database         

  `DB_PASSWORD`                       Bergantung Password         Tidak ada default.
                            konfigurasi database database         

  `DB_NAME`                                   Ya Nama database    Tidak ada default.

  `DB_SSLMODE`                             Tidak Mode SSL koneksi Default `disable`;
                                                 PostgreSQL       sesuaikan dengan
                                                                  kebijakan
                                                                  database/hosting.

  `JWT_SECRET`                                Ya Secret           Minimal 32 byte.
                                                 penandatangan    
                                                 JWT              

  `JWT_EXPIRES_IN`                            Ya Masa berlaku JWT Harus tepat `8h`.

  `CORS_ALLOWED_ORIGINS`                   Tidak Daftar origin    Default `*`.
                                                 yang diizinkan,  
                                                 dipisahkan koma  

  `TRUSTED_PROXIES`                        Tidak IP/CIDR proxy    Kosong berarti
                                                 tepercaya,       tidak ada proxy
                                                 dipisahkan koma  eksplisit.
  -----------------------------------------------------------------------------------

Catatan konfigurasi:

-   Aplikasi memanggil `godotenv.Load()` untuk membaca `.env` lokal.
-   Koneksi database dibangun dari `DB_HOST`, `DB_PORT`, `DB_USER`,
    `DB_PASSWORD`, `DB_NAME`, dan `DB_SSLMODE`.
-   `DATABASE_URL` yang mungkin tercantum pada `.env.example` tidak
    digunakan oleh `config.Load()` dan `app.ConnectDB()` saat ini.
-   `APP_ENV`, `LOG_LEVEL`, serta `ADMIN_*` yang mungkin ada pada contoh
    environment tidak dibaca oleh konfigurasi Go yang diperiksa.
-   `JWT_EXPIRES_IN` selain `8h` akan menyebabkan aplikasi gagal memuat
    konfigurasi.
-   `JWT_SECRET` yang kurang dari 32 byte akan menyebabkan aplikasi
    gagal memuat konfigurasi.
-   Untuk production, gunakan secret yang acak dan kuat, batasi CORS ke
    origin resmi, dan atur SSL database serta trusted proxy sesuai
    infrastruktur.

Contoh format `.env` tanpa nilai rahasia:

``` dotenv
APP_NAME=araya-m1
APP_PORT=8080

DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=your-database-user
DB_PASSWORD=your-database-password
DB_NAME=your-database-name
DB_SSLMODE=require

JWT_SECRET=replace-with-a-random-secret-at-least-32-bytes
JWT_EXPIRES_IN=8h

CORS_ALLOWED_ORIGINS=http://localhost:3000
TRUSTED_PROXIES=
```

## Database dan SQLC

### Model data

Berdasarkan Prisma schema dan migration, entitas utama M1 mencakup:

  -----------------------------------------------------------------------
  Tabel                               Fungsi
  ----------------------------------- -----------------------------------
  `users`                             Identitas pengguna, password hash,
                                      role, division, partner, unit
                                      usaha, dan status akun.

  `roles`                             Master role global, berlaku untuk
                                      semua partner (tanpa `partner_id`).

  `divisions`                         Master division global, berlaku
                                      untuk semua partner (tanpa
                                      `partner_id`).

  `permissions`                       Katalog permission.

  `role_permissions`                  Relasi role dengan permission.

  `division_permissions`              Relasi division dengan permission.

  `audit_logs`                        Riwayat aktivitas dan perubahan
                                      data.

  `revoked_tokens`                    Daftar jti token yang dicabut saat
                                      logout, dibersihkan setelah
                                      kedaluwarsa.
  -----------------------------------------------------------------------

ID pada tabel inti menggunakan integer. `partner_id` dan `unit_usaha_id`
merupakan referensi data lintas modul; M1 tidak menjadi pemilik data
utama partner atau unit usaha.

### Query dan generate code

Query sumber berada di:

``` text
internal/usecase/repository/queries/
```

Kode Go hasil generate berada di:

``` text
internal/usecase/repository/sqlc/
```

Setelah mengubah query SQL, jalankan:

``` bash
sqlc generate
```

Jangan mengedit file hasil generate secara manual. Perubahan query harus
diikuti pemeriksaan hasil generate dan pengujian package terkait.

### Migration

Riwayat migration berada di `prisma/migrations/`. Terapkan migration
mengikuti urutan dan prosedur database proyek yang telah disetujui.
Jangan menjalankan ulang migration terhadap database yang sudah berisi
data tanpa meninjau dampaknya terlebih dahulu.

## Arsitektur dan Alur Autentikasi

### Login

1.  Klien mengirim email dan password ke `POST /api/m1/v1/auth/login`.
2.  M1 mencari akun berdasarkan email.
3.  M1 memeriksa password menggunakan bcrypt.
4.  M1 menolak akun nonaktif.
5.  M1 mengambil permission role dari `role_permissions`.
6.  M1 menerbitkan JWT dengan algoritma HS256.
7.  Klien mengirim JWT pada request terlindungi melalui header
    `Authorization`.

### Logout dan pencabutan token

-   Klien memanggil `POST /auth/sessions/logout` dengan Bearer token.
-   M1 menyimpan `jti` dan `exp` token di `revoked_tokens`, lalu
    mencatat audit `LOGOUT`.
-   Baris yang sudah kedaluwarsa dihapus otomatis pada setiap logout.
-   Request berikutnya dengan token yang sama dijawab 401.
-   Token tanpa `jti` (terbit sebelum fitur ini) ditolak.

### Klaim JWT

Token yang dibuat `internal/usecase/utils/token.go` memuat klaim:

  Klaim              Keterangan
  ------------------ -----------------------------------------------------
  `user_id`          ID pengguna M1.
  `name`             Nama pengguna.
  `role`             Kode role pengguna.
  `permissions`      Daftar permission role yang diambil saat login.
  `division`         Kode division atau `null`.
  `partner_id`       ID partner atau `null`.
  `unit_usaha_id`    ID unit usaha atau `null`.
  `is_super_admin`   Penanda Super Admin.
  `iss`              Issuer dari `APP_NAME`.
  `iat`              Waktu token diterbitkan.
  `exp`              Waktu kedaluwarsa, 8 jam setelah token diterbitkan.
  `jti`              ID unik token, dipakai untuk pencabutan saat logout.

Token ditandatangani dengan HS256. Modul lain harus menggunakan kontrak
secret/validasi token yang disepakati bersama.

### Endpoint profil

`GET /api/m1/v1/me` mengembalikan subset klaim profil: `user_id`,
`name`, `role`, `division`, `partner_id`, `unit_usaha_id`, dan
`is_super_admin`.

## Authorization

Middleware yang tersedia pada source code:

  -----------------------------------------------------------------------
  Middleware                          Fungsi
  ----------------------------------- -----------------------------------
  `Authenticate`                      Memvalidasi JWT tanpa membaca
                                      database. Tidak dipakai route M1.

  `AuthenticateCurrent`               Memvalidasi JWT, menolak token yang
                                      sudah dicabut, lalu membaca ulang
                                      status akun, role, division, dan
                                      permission dari database.

  `RequireRole`                       Membatasi route berdasarkan role;
                                      Super Admin dilewatkan oleh helper
                                      ini.

  `RequireSuperAdmin`                 Membatasi route khusus Super Admin.

  `RequirePermission`                 Memeriksa kode permission pada
                                      klaim JWT.

  `RequireDivision`                   Membatasi route berdasarkan kode
                                      division yang ditentukan.
  -----------------------------------------------------------------------

### Pembatasan route yang terpasang saat ini

  -----------------------------------------------------------------------
  Route                               Pembatasan pada router saat ini
  ----------------------------------- -----------------------------------
  `/auth/sessions/logout`             JWT valid dan belum dicabut

  `/me`                               JWT valid dan belum dicabut

  `/auth/change-password`             JWT valid dan belum dicabut

  `/users` dan turunannya             JWT valid dan Super Admin.

  `/roles` dan turunannya             JWT valid dan Super Admin.

  `/permissions`                      JWT valid dan Super Admin.

  `/divisions`                        JWT valid dan Super Admin.

  `/audit-logs`                       JWT valid, role partner, dan
                                      permission `m1.audit.read`.

  `/platform/audit-logs`              JWT valid dan Super Admin.
  -----------------------------------------------------------------------

### Catatan authorization yang penting

Keterangan berikut menggambarkan source code yang diperiksa:

-   `RequirePermission` dan `RequireDivision` sudah tersedia sebagai
    helper middleware, tetapi belum dipasang pada seluruh route yang
    memerlukan pemeriksaan tersebut.
-   Login saat ini mengambil permission dari `role_permissions`. Tabel
    `division_permissions` ada pada schema, tetapi belum digunakan dalam
    query login yang diperiksa.
-   Route M1 memakai `AuthenticateCurrent`, yang membaca ulang akun,
    role, division, dan permission pada setiap request. Akun atau role
    nonaktif ditolak walau token belum kedaluwarsa.
-   Perubahan role atau permission setelah token diterbitkan tidak
    otomatis memperbarui klaim token yang sudah dimiliki pengguna.
    Pengguna perlu login kembali untuk memperoleh token baru.
-   Pemeriksaan JWT pada M1 tidak otomatis mengamankan endpoint M2--M5.
    Setiap modul harus menerapkan validasi token dan authorization pada
    API miliknya.

## Daftar Endpoint

Base path seluruh resource M1:

``` text
/api/m1/v1
```

### Health dan dokumentasi

  -------------------------------------------------------------------------
  Method            Endpoint            Akses             Fungsi
  ----------------- ------------------- ----------------- -----------------
  GET               `/health`           Publik            Memeriksa
                                                          aplikasi
                                                          berjalan.

  GET               `/ready`            Publik            Memeriksa koneksi
                                                          database.

  GET               `/openapi.yaml`     Publik            Mengambil
                                                          spesifikasi
                                                          OpenAPI.

  GET               `/api/m1/v1/docs`   Publik            Membuka halaman
                                                          Scalar API
                                                          Reference.
  -------------------------------------------------------------------------

### Authentication

  -------------------------------------------------------------------------------
  Method            Endpoint                  Akses             Fungsi
  ----------------- ------------------------- ----------------- -----------------
  POST              `/auth/login`             Publik            Login dan
                                                                menerbitkan JWT.

  POST              `/auth/sessions/logout`   JWT valid dan       Mencabut JWT
                                                belum dicabut      yang dipakai
                                                                   di server dan
                                                                   mencatat audit
                                                                   LOGOUT.

  GET               `/me`                     Bearer            Mendapatkan
                                                                profil pengguna
                                                                sendiri.

  POST              `/auth/change-password`   Bearer            Mengubah kata
                                                                sandi sendiri.
  -------------------------------------------------------------------------------

### Users

Seluruh route `/users` saat ini memerlukan Bearer JWT dan Super Admin.

  -----------------------------------------------------------------------
  Method                  Endpoint                Fungsi
  ----------------------- ----------------------- -----------------------
  GET                     `/users`                Daftar akun dengan
                                                  pagination dan filter.

  POST                    `/users`                Membuat akun.

  GET                     `/users/{id}`           Detail akun.

  PATCH                   `/users/{id}`           Memperbarui nama, role,
                                                  division, atau unit
                                                  usaha.

  PATCH                   `/users/{id}/status`    Mengaktifkan atau
                                                  menonaktifkan akun.

  DELETE                  `/users/{id}`           Menonaktifkan akun
                                                  secara soft.
  -----------------------------------------------------------------------

Filter daftar akun:

  -----------------------------------------------------------------------
  Parameter               Tipe                    Keterangan
  ----------------------- ----------------------- -----------------------
  `page`                  integer                 Halaman, default 1.

  `per_page`              integer                 Item per halaman,
                                                  default 20, maksimal
                                                  100.

  `partner_id`            integer                 Filter partner untuk
                                                  Super Admin.

  `q`                     string                  Pencarian nama atau
                                                  email.

  `role`                  string                  Filter kode role.

  `division`              string                  Filter kode division.

  `unit_usaha_id`         integer                 Filter unit usaha.

  `is_active`             boolean                 Filter status akun.

  `sort`                  string                  `created_at`,
                                                  `-created_at`, `name`,
                                                  `-name`, `email`, atau
                                                  `-email`. Default
                                                  `-created_at`.
  -----------------------------------------------------------------------

Body pembuatan akun:

``` json
{
  "partner_id": 15,
  "name": "Staf Operasional",
  "email": "staf@partner.example",
  "password": "Password123",
  "role": "manajer",
  "division": "operasional",
  "unit_usaha_id": 10
}
```

`partner_id`, `name`, `email`, `password`, `role`, dan `division` wajib
untuk route pembuatan akun saat ini. `unit_usaha_id` bersifat opsional.

Body pembaruan akun:

``` json
{
  "name": "Nama Pengguna Baru",
  "division": "keuangan",
  "unit_usaha_id": 12
}
```

Minimal satu field harus dikirim. Email dan partner_id tidak dapat
diperbarui melalui endpoint ini.

Body perubahan status:

``` json
{
  "is_active": false
}
```

Deaktivasi bersifat soft: data akun dan riwayatnya tetap disimpan. Akun
sendiri tidak dapat dinonaktifkan oleh pemilik akun tersebut.

### Roles dan permissions

Route berikut hanya dapat diakses Super Admin.

  Method   Endpoint                    Fungsi
  -------- --------------------------- ------------------------------------
  GET      `/roles`                    Daftar role global.
  POST     `/roles`                    Membuat role global.
  GET      `/roles/{id}/permissions`   Melihat permission role.
  PUT      `/roles/{id}/permissions`   Mengganti seluruh permission role.
  GET      `/permissions`              Melihat katalog permission.

Filter daftar role:

-   `page`
-   `per_page`
-   `is_active`

Contoh membuat role:

``` json
{
  "code": "supervisor_gudang",
  "name": "Supervisor Gudang",
  "description": "Mengelola aktivitas gudang",
  "permission_codes": [
    "m1.users.read",
    "m1.audit.read"
  ]
}
```

Contoh mengganti permission role:

``` json
{
  "permission_codes": [
    "m1.users.read",
    "m1.audit.read"
  ]
}
```

Array `permission_codes` kosong diperbolehkan untuk menghapus seluruh
penetapan permission pada role.

Katalog permission hanya dapat dibaca melalui `GET /permissions`.
Penambahan permission baru bukan operasi yang disediakan endpoint ini.

### Divisions

Route division hanya dapat diakses Super Admin.

  Method   Endpoint       Fungsi
  -------- -------------- --------------------------
  GET      `/divisions`   Daftar division global.
  POST     `/divisions`   Membuat division global.

Filter daftar division:

-   `page`
-   `per_page`
-   `is_active`

Contoh membuat division:

``` json
{
  "code": "operasional",
  "name": "Operasional"
}
```

Kode harus menggunakan huruf kecil, angka, tanda hubung, atau garis
bawah.

### Audit Logs

  ------------------------------------------------------------------------------
  Method            Endpoint                 Akses             Fungsi
  ----------------- ------------------------ ----------------- -----------------
  GET               `/audit-logs`            Role partner      Riwayat audit
                                             dengan            untuk partner
                                             `partner_id`      terkait.

  GET               `/platform/audit-logs`   Super Admin       Riwayat audit
                                                               seluruh platform.
  ------------------------------------------------------------------------------

Kedua endpoint mendukung `page` dan `per_page`. Data diurutkan dari
waktu terbaru (`created_at DESC`) dan ID terbesar lebih dahulu untuk
timestamp yang sama.

Pada implementasi saat ini, route `/audit-logs` memeriksa role partner,
scope partner, dan permission `m1.audit.read`.

## Format Request dan Response

### Header autentikasi

Endpoint terlindungi menggunakan:

``` http
Authorization: Bearer <access_token>
Content-Type: application/json
```

Endpoint login tidak memerlukan Bearer token.

### Respons sukses satu data

``` json
{
  "success": true,
  "data": {
    "id": 25,
    "name": "Staf Operasional"
  }
}
```

### Respons sukses daftar

``` json
{
  "success": true,
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 0,
    "total_pages": 0
  }
}
```

Daftar kosong tetap menggunakan HTTP 200.

### Respons gagal

``` json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke fitur ini"
  }
}
```

### Respons validasi

``` json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Data yang dikirim tidak valid",
    "fields": {
      "email": "Format email tidak valid"
    }
  }
}
```

### Status HTTP dan kode error

  ------------------------------------------------------------------------------
                          HTTP Kode umum                   Penggunaan
  ---------------------------- --------------------------- ---------------------
                           200 ---                         Operasi berhasil,
                                                           termasuk daftar
                                                           kosong.

                           201 ---                         Resource berhasil
                                                           dibuat.

                           400 `BAD_REQUEST`               Format JSON atau body
                                                           tidak valid.

                           400 `INVALID_PARAMETER`         Parameter query atau
                                                           path tidak valid.

                           401 `INVALID_CREDENTIALS`       Login gagal.

                           401 `UNAUTHENTICATED`           Token tidak ada,
                                                           tidak sah, atau
                                                           kedaluwarsa.

                           403 `FORBIDDEN`                 Pengguna tidak
                                                           memiliki kewenangan.

                           404 `USER_NOT_FOUND` /          Akun atau role tidak
                               `NOT_FOUND`                 ditemukan.

                           409 `EMAIL_DUPLICATE`           Email sudah
                                                           terdaftar.

                           409 `ROLE_CODE_DUPLICATE`       Kode role sudah
                                                           digunakan.

                           409 `DIVISION_CODE_DUPLICATE`   Kode division sudah
                                                           digunakan.

                           409 `CANNOT_DEACTIVATE_SELF`    Akun sendiri tidak
                                                           boleh dinonaktifkan.

                           409 `USER_ALREADY_ACTIVE`       Akun sudah aktif.

                           409 `USER_ALREADY_INACTIVE`     Akun sudah nonaktif.

                           409 `USER_REFERENCE_INACTIVE`   Role atau division
                                                           akun tidak aktif saat
                                                           aktivasi.

                           422 `VALIDATION_ERROR`          Field request tidak
                                                           memenuhi validasi.

                           422 `OLD_PASSWORD_INVALID`      Kata sandi lama
                                                           salah.

                           500 `INTERNAL_ERROR`            Kesalahan internal
                                                           server.

                           503 `DATABASE_UNAVAILABLE`      Database tidak dapat
                                                           diakses pada
                                                           readiness check.
  ------------------------------------------------------------------------------

Pesan error mengikuti implementasi handler. Untuk kontrak final
antarmodul, perubahan kode atau pesan error harus diselaraskan pada
OpenAPI dan disepakati dengan pemakai API.

## Audit Log

Audit log disimpan di tabel `audit_logs`.

Informasi yang disediakan oleh model audit:

-   `id`
-   `partner_id`
-   `user_id`
-   `action`
-   `resource`
-   `resource_id`
-   `old_value`
-   `new_value`
-   `ip_address`
-   `actor_name`
-   `actor_email`
-   `created_at`

Aktivitas yang dicatat pada implementasi mencakup login berhasil/gagal,
perubahan kata sandi, pembuatan dan perubahan akun, perubahan status
akun, pembuatan role, penetapan permission role, pembuatan division,
dan logout.

Password, JWT, dan secret tidak boleh dimasukkan ke nilai audit.

## Pengujian

Jalankan perintah dari root repository.

### Formatting

``` bash
gofmt -w .
```

Jika tidak ingin memformat seluruh project, jalankan `gofmt -w` pada
file Go yang berubah.

### Unit test dan package test

``` bash
go test ./...
```

Untuk mengabaikan hasil cache:

``` bash
go test ./... -count=1
```

### Static analysis

``` bash
go vet ./...
```

### SQLC

Jika query SQL berubah:

``` bash
sqlc generate
```

Setelah generate, jalankan kembali:

``` bash
go test ./...
go vet ./...
```

Jangan menyatakan pengujian lulus sebelum perintah benar-benar
dijalankan pada lingkungan dengan versi Go yang sesuai `go.mod`.

## Dokumentasi API

Spesifikasi OpenAPI berada pada:

``` text
docs/openapi.yaml
```

File tersebut di-embed oleh `docs/embed.go` dan disajikan melalui:

``` text
GET /openapi.yaml
```

Halaman Scalar API Reference tersedia di:

``` text
GET /api/m1/v1/docs
```

Pada lingkungan lokal:

-   OpenAPI: `http://localhost:8080/openapi.yaml`
-   API Reference: `http://localhost:8080/api/m1/v1/docs`

Jika endpoint atau DTO berubah, perbarui `docs/openapi.yaml` agar
kontrak dokumentasi sesuai dengan route, handler, dan response aktual.

## Deployment

Repository menggunakan entry point `main.go` di root. Untuk deployment
Vercel:

1.  Hubungkan repository ke project Vercel.
2.  Gunakan runtime Go yang sesuai dengan `go.mod`.
3.  Isi environment variable yang diperlukan melalui dashboard hosting.
4.  Pastikan konfigurasi database dapat diakses dari lingkungan
    deployment.
5.  Pastikan `JWT_SECRET` dan `JWT_EXPIRES_IN=8h` valid.
6.  Setelah deploy, uji `/health` dan `/ready`.
7.  Periksa `/openapi.yaml` dan `/api/m1/v1/docs`.

`main.go` menggunakan environment `PORT` jika tersedia; jika tidak,
aplikasi menggunakan `APP_PORT`.

## Keterbatasan Implementasi Saat Ini

Bagian ini penting agar dokumentasi tidak dianggap menyatakan fitur yang
belum ada.

1.  **Pembaruan klaim:** perubahan role, division, atau permission tidak
    mengubah token yang sudah diterbitkan. Klaim baru diperoleh pada
    login berikutnya.
2.  **Permission route:** permission role dimasukkan ke JWT, dan helper
    `RequirePermission` tersedia. Namun pemasangan middleware permission
    pada seluruh endpoint belum selesai.
3.  **Permission division:** tabel `division_permissions` tersedia pada
    schema, tetapi query login saat ini hanya mengambil permission dari
    `role_permissions`.
4.  **Pembatasan division:** helper `RequireDivision` tersedia, tetapi
    belum dipasang pada seluruh route bisnis yang membutuhkan pembatasan
    division.
5.  **Role:** endpoint yang tersedia hanya daftar, pembuatan, melihat
    permission role, dan mengganti permission role. Update/delete/status
    role belum tersedia.
6.  **Division:** endpoint yang tersedia hanya daftar dan pembuatan
    division. Update/delete/status division belum tersedia.
7.  **Integrasi M3:** endpoint server-to-server untuk provisioning akun
    Partner Admin dan sinkronisasi status partner belum tersedia pada
    router saat ini.
8.  **Pencabutan token hanya berlaku di M1.** Modul lain yang memverifikasi
    JWT sendiri tetap menerima token yang sudah di-logout sampai `exp`,
    kecuali memanggil `GET /me`.

## Acuan

-   Product Requirements Document (PRD) Modul M1 --- revisi 10 November
    2026.
-   Aturan Pengembangan Backend Modul M1.
-   Panduan API Bersama v1.0 (draf), 28 Oktober 2026.
-   Source code, route, handler, usecase, query SQL, migration, dan
    konfigurasi repository BE Araya.