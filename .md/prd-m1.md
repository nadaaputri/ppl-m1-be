# PRD Modul M1 --- Manajemen Pengguna & Hak Akses

**Nama file:** `prd-m1.md`\
**Proyek:** Araya Business Hub\
**Modul:** M1 --- Manajemen Pengguna & Hak Akses\
**Acuan API:** *Panduan API Bersama v1.0 (draf), 28 Oktober 2026*

## 1. Masalah

Araya Business Hub membutuhkan pengelolaan akun dan autentikasi terpusat
agar pengguna dapat masuk sesuai status akun dan kewenangannya.
Aktivitas perubahan data juga perlu dicatat supaya dapat ditelusuri. API
M1 harus mengikuti kontrak bersama agar dapat digunakan oleh modul lain
tanpa perbedaan format.

## 2. Pengguna

-   **Administrator:** mengelola akun, memeriksa audit log, dan
    menjalankan fungsi administrasi sesuai kewenangannya.
-   **Pengguna terdaftar:** login, mengakses fungsi sesuai role/divisi,
    dan mengubah kata sandi sendiri.
-   **Modul lain:** menggunakan token yang diterbitkan M1 untuk
    autentikasi pengguna dan membaca profil melalui kontrak bersama.

## 3. Tujuan dan ukuran keberhasilan

-   Pengguna dapat login dengan kredensial yang valid; login ditolak
    bila kredensial salah atau akun tidak aktif.
-   Token yang diterbitkan dapat dipakai lintas modul dan memiliki masa
    berlaku **8 jam**, sesuai Panduan API Bersama.
-   Akses dapat dibatasi berdasarkan informasi `role` dan
    `unit_usaha_id` pada token; setiap modul memeriksa token yang sama.
-   Administrator dapat menambah dan menonaktifkan akun serta melihat
    audit log aktivitas.
-   Pengguna dapat mengganti kata sandi setelah memasukkan kata sandi
    lama dengan benar.
-   Endpoint M1 memiliki kontrak OpenAPI, format respons, kode status,
    dan aturan penamaan yang konsisten dengan panduan bersama.

## 4. User story dan kriteria penerimaan

### M1-US1 --- Login dan hak akses per divisi

**User story:** Sebagai Administrator, saya ingin login dan mengatur hak
akses per divisi, supaya setiap pengguna hanya bisa mengakses data
sesuai kewenangannya.

-   Kredensial valid: pengguna berhasil login dan akses dibatasi sesuai
    kewenangan.
-   Kredensial tidak valid: akses ditolak.
-   Token pengguna diterbitkan oleh M1 dan dapat digunakan oleh modul
    lain melalui `Authorization: Bearer <token>`.
-   Kontrak bersama mencantumkan `GET /api/m1/v1/me` untuk token login
    dan profil pengguna (K-01).
-   Token tidak ada, salah, atau kedaluwarsa menghasilkan HTTP `401`;
    token valid tetapi role tidak berhak menghasilkan `403`.

### M1-US2 --- Audit log aktivitas pengguna

**User story:** Sebagai Administrator, saya ingin sistem mencatat audit
log aktivitas pengguna, supaya setiap perubahan data dapat ditelusuri.

-   Setelah perubahan data berhasil disimpan, sistem mencatat siapa yang
    melakukan perubahan, perubahan apa yang dilakukan, dan kapan
    waktunya.
-   Administrator dapat melihat riwayat audit log yang diurutkan
    berdasarkan waktu.
-   Audit log yang perlu dipersistenkan dan ditampilkan harus disepakati
    dalam kontrak M1; pencatatan ke log aplikasi saja belum membuktikan
    tersedianya riwayat yang dapat ditelusuri melalui API.
-   Jangan mencatat kata sandi, token, atau rahasia autentikasi ke audit
    log.

### M1-US3 --- Menambah dan menonaktifkan akun

**User story:** Sebagai Administrator, saya ingin menambah dan
menonaktifkan akun pengguna, supaya hanya orang yang berwenang yang
dapat masuk.

-   Data akun valid yang disimpan oleh administrator menghasilkan akun
    aktif dan dapat login sesuai kewenangannya.
-   Akun yang dinonaktifkan tidak dapat login.
-   Data historis tetap tersimpan; akun yang memiliki riwayat tidak
    dihapus secara fisik.

### M1-US4 --- Mengubah kata sandi

**User story:** Sebagai Pengguna, saya ingin mengubah kata sandi saya,
supaya akun tetap aman.

-   Jika kata sandi lama benar dan kata sandi baru memenuhi aturan, kata
    sandi berhasil diubah.
-   Jika kata sandi lama salah, permintaan ditolak dan kata sandi tidak
    berubah.
-   Kata sandi tidak dikembalikan dalam respons API maupun ditulis ke
    log.

## 5. Fitur dan aturan API bersama

### 5.1 Pola endpoint dan penamaan

Semua endpoint M1 menggunakan pola `/api/m1/v1/{resource}`. Nama
resource berupa kata benda jamak dalam bahasa Inggris, huruf kecil, dan
menggunakan tanda hubung bila diperlukan. Nama field JSON memakai
`snake_case`.

Kontrak yang secara eksplisit disebut Panduan API Bersama: -
`GET /api/m1/v1/me` --- token login dan profil pengguna, kontrak K-01.

Endpoint untuk login, pengelolaan akun, audit log, dan perubahan kata
sandi harus dirinci di `docs/openapi.yaml`. Nama endpoint berikut
**usulan awal untuk dibahas saat Sprint Planning**, bukan kontrak
final: - `POST /api/m1/v1/auth/sessions` --- login/menerbitkan token. -
`GET /api/m1/v1/users` dan `POST /api/m1/v1/users` --- melihat daftar
dan menambah akun. - `POST /api/m1/v1/users/{id}/deactivate` ---
menonaktifkan akun. - `GET /api/m1/v1/audit-logs` --- melihat audit
log. - `PATCH /api/m1/v1/users/me/password` --- mengubah kata sandi
sendiri.

Endpoint usulan tidak boleh dianggap final sebelum dikonfirmasi dan
dicatat pada OpenAPI/kontrak bersama. Jangan menggunakan pola seperti
`/getUsers` atau `/createUser`.

### 5.2 Autentikasi dan hak akses

-   Login hanya disediakan oleh M1; modul lain tidak membuat login
    sendiri.
-   Permintaan pengguna yang memerlukan autentikasi mengirim
    `Authorization: Bearer <token>`.
-   Token JWT dari M1 berlaku 8 jam dan memuat informasi seperti
    `user_id`, `name`, `role`, `division`, `unit_usaha_id`, dan `exp`,
    sesuai bentuk token dalam panduan.
-   Modul membaca `role` dan `unit_usaha_id` dari token untuk
    pemeriksaan kewenangan, tanpa menanyakan ulang ke M1 pada setiap
    permintaan.
-   Panggilan antarmodul tanpa pengguna memakai `X-Service-Key`; kunci
    disimpan di `.env`, bukan ditulis di kode atau dokumentasi publik.
-   Endpoint publik, jika ada, harus ditandai sebagai publik di OpenAPI.

### 5.3 Amplop respons

Respons sukses menggunakan bentuk bersama: - Satu data:
`{"success": true, "data": {...}}` - Daftar data:
`{"success": true, "data": [...], "meta": {"page": 1, "per_page": 20, "total": 0, "total_pages": 0}}` -
Galat:
`{"success": false, "error": {"code": "ERROR_CODE", "message": "Pesan dalam bahasa Indonesia"}}` -
Untuk galat validasi `422`, `error.fields` berisi nama field dan pesan
validasinya.

Kode `error.code` ditulis dengan `HURUF_BESAR_DENGAN_GARIS_BAWAH` dan
tidak berubah. Pesan `error.message` harus sama persis dengan teks yang
ditetapkan pada Acceptance Criteria. Modul pemakai membaca `error.code`,
bukan mengandalkan teks `message`.

### 5.4 Kode status HTTP

  -----------------------------------------------------------------------
  Status                              Penggunaan di M1
  ----------------------------------- -----------------------------------
  `200`                               Permintaan berhasil, termasuk
                                      daftar audit log kosong.

  `201`                               Akun baru berhasil dibuat.

  `400`                               Parameter permintaan tidak masuk
                                      akal atau salah format.

  `401`                               Kredensial tidak valid, token tidak
                                      ada/tidak sah, atau token
                                      kedaluwarsa.

  `403`                               Pengguna sudah terautentikasi
                                      tetapi tidak memiliki kewenangan.

  `404`                               Akun atau data yang diminta tidak
                                      ditemukan.

  `409`                               Konflik data, misalnya email sudah
                                      digunakan atau perubahan status
                                      tidak sesuai.

  `422`                               Isian gagal validasi.

  `500`                               Kesalahan internal server; tidak
                                      digunakan untuk menyatakan
                                      kesalahan input pengguna.
  -----------------------------------------------------------------------

Daftar kosong bukan galat: kembalikan `200` dengan `data: []`.

### 5.5 Format data

-   Field JSON menggunakan `snake_case` dan bahasa Inggris.
-   Waktu menggunakan ISO 8601 dengan zona waktu WIB, misalnya
    `2026-10-28T09:15:00+07:00`.
-   Nilai kosong menggunakan `null`, bukan teks kosong.
-   Status menggunakan teks huruf kecil `snake_case` dalam bahasa
    Indonesia sesuai studi kasus.
-   ID adalah bilangan bulat dari modul pemilik data.
-   Untuk daftar, gunakan parameter bersama `page` dan `per_page`
    (bawaan 20, maksimal 100); gunakan `sort` dengan awalan tanda minus
    untuk urutan menurun bila didukung endpoint.

### 5.6 OpenAPI dan pengujian kontrak

-   Semua endpoint M1, termasuk endpoint internal, didokumentasikan di
    `docs/openapi.yaml`.
-   Dokumentasi API disediakan pada `/api/m1/v1/docs`.
-   Kontrak yang digunakan lintas modul dicatat di repositori kontrak
    bersama dan harus disetujui oleh semua pemakai terkait serta
    di-merge oleh Product Owner.
-   Setiap kontrak memiliki minimal satu pengujian sukses dan satu
    pengujian galat di koleksi Postman.
-   Perubahan yang menghapus/mengganti nama field, mengubah tipe/format,
    atau membuat field opsional menjadi wajib merupakan perubahan
    merusak dan memerlukan persetujuan pemakai atau versi baru.

## 6. Di luar cakupan

-   Implementasi fitur bisnis dari M2--M15.
-   Menentukan sendiri seluruh matriks role/divisi dan daftar role final
    tanpa kesepakatan tim.
-   Menambahkan fitur audit seperti ekspor, retensi khusus, analitik,
    atau pencarian lanjutan jika belum disepakati.
-   Menganggap pencatatan `slog` di aplikasi sebagai pengganti
    penyimpanan audit log yang dapat dibaca kembali melalui API.
-   Mengubah kontrak bersama yang sudah disetujui tanpa proses perubahan
    kontrak.

## 7. Hal yang harus diputuskan saat Sprint Planning

-   Nama dan bentuk final endpoint login, pengelolaan akun, audit log,
    dan ubah kata sandi.
-   Struktur respons profil pada `GET /api/m1/v1/me` serta field token
    final.
-   Daftar role/divisi dan aturan pembatasan akses yang disepakati
    bersama.
-   Detail audit log: jenis perubahan yang dicatat, identitas pelaku,
    waktu, objek yang berubah, dan apakah nilai sebelum/sesudah
    perubahan diperlukan.
-   Siapa yang boleh melihat audit log dan apakah endpoint daftar
    memerlukan pagination/filter tertentu.
-   Aturan validasi kata sandi serta pesan galat yang harus sama persis
    dengan Acceptance Criteria.
-   Nama `error.code` dan pesan bahasa Indonesia untuk setiap kondisi
    galat.
-   Penambahan Acceptance Criteria dari tim; AC tambahan bukan ketentuan
    final sebelum disepakati.

## 8. Catatan implementasi saat ini

Berdasarkan tinjauan kode yang tersedia sebelumnya, login, pemeriksaan
akun aktif, pembuatan JWT, dan pencatatan event login melalui logger
aplikasi sudah tampak tersedia. Namun, middleware autentikasi/hak akses
lintas endpoint, endpoint profil `GET /api/m1/v1/me`, penyimpanan serta
endpoint daftar audit log, endpoint tambah/nonaktif akun, dan endpoint
ubah kata sandi belum tampak lengkap pada bahan kode yang ditinjau.
Status ini perlu diverifikasi lagi pada kode terbaru sebelum dianggap
sebagai kondisi implementasi final.

## Referensi

-   *Panduan API Bersama v1.0 (draf)*, disusun 28 Oktober 2026,
    khususnya bagian 1--4, 7--8 dan kontrak K-01.
-   *Studi Kasus Capstone Project --- Kelas A*, user story M1-US1 sampai
    M1-US4.
