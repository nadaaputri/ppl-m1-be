# Aturan Pengembangan Backend --- Modul M1

**Proyek:** Araya Business Hub\
**Scope pekerjaan:** Modul M1 --- Manajemen Pengguna & Hak Akses, khusus
M1-US1 sampai M1-US4\
**Acuan API:** *Panduan API Bersama v1.0 (draf)*, disusun 28 Oktober
2026

Dokumen ini menjadi aturan kerja saat mempelajari, meninjau, atau
mengubah kode backend M1. Jika aturan atau keputusan belum jelas, jangan
menebak dan jangan langsung mengimplementasikan perubahan; tanyakan
terlebih dahulu.

## 1. Batasan pekerjaan dan instruksi

1.  Kerjakan hanya hal yang secara eksplisit diminta pengguna.
2.  Jangan menambah fitur, endpoint, tabel, kolom, dependensi,
    konfigurasi, file, atau refactor yang tidak diminta.
3.  Jika merasa perlu membuat sesuatu di luar instruksi, jelaskan alasan
    dan dampaknya, lalu minta persetujuan sebelum mengerjakannya.
4.  Jangan keluar dari scope Modul M1 dan M1-US1 sampai M1-US4.
5.  Jangan mengerjakan fitur modul lain (M2--M15), kecuali pengguna
    secara eksplisit meminta pekerjaan integrasi yang memang diperlukan.
6.  Jika acceptance criteria, kontrak API, atau aturan bisnis belum
    jelas, ajukan pertanyaan sebelum menetapkan keputusan.
7.  Jangan menganggap usulan sebagai keputusan final. Minta konfirmasi
    untuk hal yang belum disepakati tim.
8.  Sebelum mengubah kode, jelaskan secara singkat apa yang akan diubah
    bila perubahan berisiko atau cakupannya belum jelas.

## 2. Pelajari kode yang sudah ada terlebih dahulu

1.  Periksa struktur proyek, dokumentasi, konfigurasi, route,
    controller/handler, middleware, use case/service, repository/query,
    model, dan test yang relevan sebelum membuat perubahan.
2.  Telusuri alur fitur dari route sampai penyimpanan data atau respons
    API; jangan menilai hanya dari satu file.
3.  Cari di seluruh kode proyek apakah fungsi, endpoint, middleware,
    validasi, helper, atau service yang dibutuhkan sudah tersedia.
4.  Jika fungsi sudah ada, gunakan atau perbaiki implementasi yang ada
    sesuai instruksi. Jangan membuat versi kedua dengan tujuan yang
    sama.
5.  Jangan menimpa atau menghapus perilaku yang sudah berjalan tanpa
    memahami dampaknya.
6.  Perhatikan pola penamaan, struktur folder, dependensi, dan gaya kode
    yang sudah digunakan. Ikuti pola proyek yang konsisten, kecuali
    pengguna meminta perubahan.
7.  Jangan menganggap fitur belum ada hanya karena tidak terlihat di
    satu file. Periksa pemanggilan, registrasi route, konfigurasi, dan
    pengujiannya.
8.  Jika akses terhadap seluruh kode atau informasi yang diperlukan
    belum tersedia, nyatakan keterbatasannya dan minta bahan yang
    dibutuhkan. Jangan mengklaim sudah memeriksa seluruh kode jika
    belum.

## 3. Aturan penulisan kode

1.  **Maksimal 20 baris per fungsi.** Usahakan setiap fungsi memiliki
    satu tanggung jawab yang jelas.
2.  Jika fungsi melebihi 20 baris, pecah menjadi fungsi bantu yang
    memiliki tujuan jelas. Jangan membuat abstraksi atau fungsi tambahan
    yang tidak diperlukan hanya untuk memecah baris.
3.  Hitung isi fungsi secara konsisten; baris kosong dan komentar tidak
    digunakan untuk menyamarkan fungsi yang terlalu panjang.
4.  Hindari duplikasi kode. Gunakan fungsi/helper yang sudah ada bila
    sesuai.
5.  Gunakan nama variabel, fungsi, tipe, dan file yang jelas serta
    konsisten dengan konvensi proyek.
6.  Tangani error secara eksplisit. Jangan mengabaikan error atau
    menggantinya dengan respons sukses.
7.  Jangan menambahkan dependensi baru, mengubah versi library, atau
    mengubah arsitektur tanpa persetujuan pengguna.
8.  Jangan mengubah format file atau merapikan bagian kode yang tidak
    terkait dengan tugas.
9.  Jangan menaruh kredensial, token, API key, password, atau service
    key di source code maupun dokumentasi yang dibagikan.
10. Jangan mencatat password, token, atau rahasia autentikasi ke log.
11. Perubahan harus sekecil mungkin dan fokus pada kebutuhan yang
    diminta.

## 4. Aturan API bersama

Semua endpoint M1 harus mengikuti *Panduan API Bersama v1.0 (draf)*.
Jika panduan diperbarui, tanyakan versi yang berlaku sebelum mengubah
kontrak.

### 4.1 Alamat dan penamaan

1.  Gunakan pola `/api/m1/v1/{resource}` untuk endpoint M1.
2.  Nama resource berupa kata benda jamak, bahasa Inggris, huruf kecil,
    dan tanda hubung bila diperlukan, misalnya `audit-logs`.
3.  Nama field JSON menggunakan `snake_case`, misalnya `user_id` dan
    `created_at`.
4.  Gunakan method HTTP sesuai tujuannya:
    -   `GET`: membaca data; tidak mengubah data.
    -   `POST`: membuat data atau menjalankan aksi/status.
    -   `PUT`: mengganti seluruh isi satu data.
    -   `PATCH`: mengubah sebagian field.
    -   `DELETE`: menonaktifkan atau menghapus sesuai aturan; data yang
        memiliki riwayat tidak dihapus fisik.
5.  Hindari kata kerja untuk operasi biasa di alamat endpoint, seperti
    `/getUsers` atau `/createUser`.
6.  Jangan menetapkan endpoint baru sebagai final tanpa konfirmasi jika
    endpoint tersebut belum ditentukan dalam kontrak bersama.

### 4.2 Autentikasi dan hak akses

1.  Login terpusat hanya disediakan oleh M1.
2.  Permintaan yang memerlukan autentikasi menggunakan header
    `Authorization: Bearer <token>`.
3.  Sesuai panduan versi ini, token JWT dari M1 berlaku **8 jam**.
4.  Token memuat informasi pengguna seperti `user_id`, `name`, `role`,
    `division`, `unit_usaha_id`, dan `exp`, mengikuti kontrak yang
    disepakati.
5.  Modul menggunakan `role` dan `unit_usaha_id` pada token untuk
    pemeriksaan kewenangan, tanpa meminta ulang informasi itu ke M1 pada
    setiap permintaan.
6.  Panggilan antarmodul tanpa pengguna menggunakan `X-Service-Key`;
    simpan nilainya di konfigurasi rahasia seperti `.env`, bukan di
    source code.
7.  Token tidak ada, tidak sah, atau kedaluwarsa menghasilkan `401`.
    Token sah tetapi role tidak berhak menghasilkan `403`.
8.  Jangan mengarang daftar role, matriks hak akses, atau format token
    final. Ikuti keputusan dan kontrak bersama yang sudah disetujui.

### 4.3 Bentuk respons

1.  Respons API menggunakan JSON dengan amplop yang konsisten.
2.  Respons sukses satu data: `{"success": true, "data": {...}}`.
3.  Respons sukses daftar:
    `{"success": true, "data": [...], "meta": {"page": 1, "per_page": 20, "total": 0, "total_pages": 0}}`.
4.  Respons galat:
    `{"success": false, "error": {"code": "ERROR_CODE", "message": "Pesan dalam bahasa Indonesia"}}`.
5.  Pada galat validasi `422`, sertakan `error.fields` yang memetakan
    nama field ke pesan validasinya.
6.  `error.code` memakai huruf besar dan garis bawah
    (`UPPER_SNAKE_CASE`) serta tidak berubah tanpa perubahan kontrak
    yang disetujui.
7.  `error.message` harus sama persis dengan pesan yang ditetapkan dalam
    Acceptance Criteria. Jangan mengganti teks pesan tanpa kesepakatan.
8.  Daftar kosong bukan galat: gunakan HTTP `200` dengan `data: []`.

### 4.4 Kode status HTTP

  -----------------------------------------------------------------------
  Status                              Arti/penggunaan
  ----------------------------------- -----------------------------------
  `200`                               Permintaan berhasil, termasuk
                                      daftar kosong.

  `201`                               Data baru berhasil dibuat.

  `400`                               Parameter permintaan salah atau
                                      tidak masuk akal.

  `401`                               Belum login, kredensial tidak
                                      valid, atau token tidak
                                      sah/kedaluwarsa.

  `403`                               Sudah terautentikasi tetapi tidak
                                      berwenang.

  `404`                               Data yang diminta tidak ditemukan.

  `409`                               Konflik dengan data atau status
                                      yang sudah ada.

  `422`                               Isian tidak lolos validasi.

  `500`                               Kesalahan internal server; jangan
                                      gunakan untuk kesalahan input
                                      pengguna.
  -----------------------------------------------------------------------

Pilih status berdasarkan kondisi sebenarnya dan kontrak yang disepakati;
jangan menggunakan satu status yang sama untuk semua galat.

### 4.5 Format data

1.  Tanggal menggunakan `YYYY-MM-DD`.
2.  Tanggal dan waktu menggunakan ISO 8601 dengan zona waktu WIB,
    misalnya `2026-10-28T09:15:00+07:00`.
3.  Uang, jika ada, dikirim sebagai bilangan bulat rupiah tanpa pemisah
    ribuan dan tanpa desimal.
4.  Nilai boolean menggunakan `true` atau `false`.
5.  Nilai kosong menggunakan `null`, bukan string kosong, jika field
    memang tidak memiliki nilai.
6.  Status memakai teks huruf kecil `snake_case` dalam bahasa Indonesia
    sesuai studi kasus.
7.  ID berupa bilangan bulat dari modul pemilik data, kecuali kontrak
    yang disepakati menyatakan lain.
8.  Untuk endpoint daftar yang memerlukan pagination, gunakan `page` dan
    `per_page`; nilai bawaan 20 dan maksimal 100 sesuai panduan. Gunakan
    filter dan `sort` hanya bila dibutuhkan oleh endpoint dan telah
    didefinisikan.

### 4.6 OpenAPI dan pengujian kontrak

1.  Dokumentasikan endpoint M1 di `docs/openapi.yaml`, termasuk endpoint
    yang hanya digunakan di dalam M1.
2.  Dokumentasi API disediakan pada `/api/m1/v1/docs`.
3.  Kontrak antarmodul harus dicatat di repositori kontrak bersama dan
    mengikuti proses persetujuan tim/PO.
4.  Setiap kontrak memiliki minimal satu pengujian sukses dan satu
    pengujian galat di koleksi Postman.
5.  Jalankan pengujian yang relevan setelah perubahan dan laporkan
    hasilnya. Jangan menyatakan test berhasil jika belum dijalankan.
6.  Jangan menghapus atau mengganti nama field, mengubah tipe/format,
    atau membuat field opsional menjadi wajib tanpa persetujuan pemakai
    kontrak atau proses versi baru.
7.  Jangan membuat endpoint, field, atau kode galat hanya demi
    melengkapi dokumentasi jika belum menjadi kebutuhan atau belum
    disetujui.

## 5. Batas scope Modul M1

Pekerjaan hanya mencakup kebutuhan berikut:

-   **M1-US1:** login dan hak akses per divisi/role sesuai kontrak yang
    disepakati.
-   **M1-US2:** pencatatan dan penelusuran audit log aktivitas pengguna.
-   **M1-US3:** menambah dan menonaktifkan akun pengguna.
-   **M1-US4:** mengubah kata sandi pengguna.

Jangan mengerjakan fitur bisnis M2--M15 atau menambah fitur di luar
empat user story tersebut tanpa instruksi eksplisit.

## 6. Prosedur sebelum dan sesudah perubahan

Sebelum mengubah kode: 1. Pahami permintaan dan tentukan US yang
terkait. 2. Periksa aturan API, PRD, Acceptance Criteria, serta
implementasi dan test yang sudah ada. 3. Cari implementasi serupa agar
tidak membuat ulang fungsi yang sudah tersedia. 4. Jika keputusan
penting belum jelas atau diperlukan perubahan di luar instruksi,
berhenti dan tanyakan kepada pengguna. 5. Sebutkan file dan perilaku
yang akan diubah jika ruang lingkupnya berpotensi meluas.

Setelah mengubah kode: 1. Tinjau diff agar hanya perubahan yang diminta
yang masuk. 2. Jalankan formatter dan test yang relevan jika tersedia.
3. Periksa ulang format API, kode status, validasi, autentikasi, dan
batas 20 baris per fungsi. 4. Laporkan file yang diubah, ringkasan
perubahan, test yang dijalankan beserta hasilnya, dan hal yang belum
dapat diverifikasi. 5. Jangan mengklaim pekerjaan selesai atau lolos
test tanpa bukti hasilnya.

## 7. Prinsip utama

**Pelajari sebelum membuat. Periksa sebelum mengulang. Kerjakan hanya
yang diperintahkan. Jika ingin menambah sesuatu atau ada keputusan yang
belum jelas, tanyakan terlebih dahulu. Jangan keluar dari scope M1-US1
sampai M1-US4.**

## Referensi

-   *Panduan API Bersama v1.0 (draf)*, disusun 28 Oktober 2026.
-   *Studi Kasus Capstone Project --- Kelas A*, M1-US1 sampai M1-US4.
-   `prd-m1.md`, PRD Modul M1 --- Manajemen Pengguna & Hak Akses.
