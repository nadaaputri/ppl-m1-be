# Product Requirements Document (PRD)
## Modul M1 – Manajemen Pengguna dan Hak Akses
**Proyek:** Araya Business Hub  
**Versi:** Revisi – 10 November 2026

---

## 1. Masalah

Araya Business Hub terdiri dari beberapa modul yang dikembangkan oleh tim berbeda. Setiap modul memiliki kebutuhan pengguna dan fitur masing-masing, tetapi pengelolaan akun, autentikasi, serta pengaturan hak akses perlu dilakukan secara terpusat agar konsisten.

Permasalahan yang perlu ditangani:

1. Akun pengguna perlu dikelola secara terpusat agar pengguna dapat mengakses modul sesuai identitas dan hak aksesnya.
2. Setiap partner memiliki pengguna dengan kebutuhan akses yang berbeda berdasarkan divisi atau peran.
3. Partner yang dinonaktifkan oleh M3 tidak boleh lagi mengakses sistem melalui akun-akun yang terkait.
4. Perubahan akun dan hak akses perlu tercatat agar dapat ditelusuri.
5. Pengguna membutuhkan mekanisme perubahan kata sandi yang aman.

M1 berperan sebagai layanan pusat pengelolaan akun, autentikasi, dan penetapan hak akses. M1 bukan pemilik data utama partner, unit usaha, produk, pelanggan, maupun transaksi.

## 2. Pengguna

| Pengguna | Tanggung jawab dalam M1 |
|---|---|
| Super Admin | Mengelola akun dan hak akses pada tingkat platform serta melihat audit log platform sesuai kewenangan. |
| Partner Admin | Mengelola akun pengguna dalam lingkup partner, mengatur penetapan hak akses sesuai kewenangan, dan melihat audit log dalam lingkup partner. |
| Pengguna Partner | Login, mengakses fitur sesuai hak akses, melihat profil sendiri, dan mengubah kata sandi. |
| Modul lain (M2–M5) | Menggunakan token dari M1 dan menerapkan pemeriksaan hak akses pada API masing-masing. |
| Backend M3 | Mengirim data akun awal Partner Admin setelah partner disetujui serta mengirim perubahan status partner kepada M1. |

Catatan: Pembagian kewenangan Super Admin dan Partner Admin perlu mengikuti matriks hak akses yang disepakati tim.

## 3. Tujuan dan Ukuran Sukses

### Tujuan
1. Menyediakan pengelolaan akun pengguna secara terpusat.
2. Menyediakan autentikasi terpusat yang digunakan oleh modul-modul Araya Business Hub.
3. Memungkinkan penetapan hak akses berdasarkan role dan/atau divisi menggunakan katalog permission yang telah disediakan modul terkait.
4. Memastikan akun dari partner yang dinonaktifkan tidak dapat digunakan untuk mengakses sistem.
5. Menyediakan audit log untuk aktivitas pengelolaan akun dan hak akses di M1.

### Ukuran sukses
- Pengguna dengan kredensial valid dapat login dan memperoleh token autentikasi.
- Kredensial tidak valid, akun nonaktif, atau partner nonaktif ditolak saat autentikasi.
- Modul lain dapat memvalidasi token dan menerapkan otorisasi sesuai kontrak bersama.
- Pengguna hanya dapat mengakses fitur yang diberikan kepadanya.
- Aktivitas penting pengelolaan akun dan hak akses di M1 tercatat dan dapat ditelusuri.
- Pengguna dapat mengganti kata sandi dengan memasukkan kata sandi lama yang benar.

## 4. User Story

| ID | User Story | Sumber |
|---|---|---|
| M1-US1 | Sebagai administrator, saya ingin login dan mengatur hak akses berdasarkan divisi agar pengguna hanya dapat mengakses fitur sesuai kewenangannya. | Studi kasus |
| M1-US2 | Sebagai administrator, saya ingin aktivitas pengguna tercatat agar perubahan dan tindakan penting dapat ditelusuri. | Studi kasus |
| M1-US3 | Sebagai administrator, saya ingin menambahkan dan menonaktifkan akun pengguna agar akses pengguna dapat dikelola. | Studi kasus |
| M1-US4 | Sebagai pengguna, saya ingin mengubah kata sandi dengan memverifikasi kata sandi lama agar akun tetap aman. | Studi kasus |
| M1-US5 | Sebagai backend M3, saya ingin mengirim data akun awal Partner Admin setelah partner disetujui agar akun dapat dibuat di M1. | Integrasi antarmodul |
| M1-US6 | Sebagai backend M3, saya ingin mengirim perubahan status partner agar M1 dapat membatasi atau memulihkan akses akun terkait. | Integrasi antarmodul |
| M1-US7 | Sebagai Partner Admin, saya ingin mengelola akun pengguna dalam lingkup partner agar staf dapat menggunakan sistem sesuai kewenangannya. | Kebutuhan sistem |
| M1-US8 | Sebagai administrator yang berwenang, saya ingin menetapkan permission dari katalog yang tersedia kepada role/divisi agar hak akses dapat diatur. | Kebutuhan sistem |
| M1-US9 | Sebagai pengguna yang sudah login, saya ingin melihat profil akun sendiri agar dapat mengetahui identitas yang sedang digunakan. | Kontrak API bersama |

## 5. Fitur

### 5.1 Fokus M1

| Kode | Fitur | Deskripsi |
|---|---|---|
| F-M1-01 | Autentikasi pengguna | Login terpusat, penerbitan JWT, dan penolakan kredensial tidak valid. |
| F-M1-02 | Profil pengguna | Menyediakan informasi pengguna yang sedang login melalui endpoint profil. |
| F-M1-03 | Pengelolaan akun | Melihat, menambahkan, memperbarui, dan menonaktifkan akun sesuai kewenangan. |
| F-M1-04 | Pengelolaan role | Mengelola role sesuai aturan dan daftar role bersama yang disepakati. |
| F-M1-05 | Katalog permission | Menyimpan dan menampilkan daftar permission yang disampaikan secara manual oleh setiap modul. M1 tidak menyediakan fitur/API bagi modul lain untuk mendaftarkan permission secara mandiri. |
| F-M1-06 | Penetapan hak akses | Menetapkan permission dari katalog kepada role dan/atau divisi sesuai kewenangan. |
| F-M1-07 | Pengelolaan divisi | Menggunakan referensi divisi untuk pengelompokan pengguna dan hak akses. Kepemilikan master data divisi perlu dikonfirmasi dengan M2/PO sebelum CRUD diimplementasikan. |
| F-M1-08 | Perubahan kata sandi | Memverifikasi kata sandi lama, memvalidasi kata sandi baru, lalu menyimpan hash kata sandi. |
| F-M1-09 | Audit log M1 | Mencatat aktivitas penting M1, termasuk pelaku, tindakan, waktu, dan informasi pendukung yang relevan. |
| F-M1-10 | Pembuatan akun awal Partner Admin | Menerima data dari M3 setelah partner disetujui dan membuat akun awal di M1. |
| F-M1-11 | Sinkronisasi status partner | Menerima event perubahan status partner dari M3 dan membatasi atau memulihkan akses akun terkait. |
| F-M1-12 | Dokumentasi API | Menyediakan dokumentasi OpenAPI yang sesuai dengan endpoint aktual dan kontrak antarmodul. |

### 5.2 Tanggung jawab kelompok lain

| Modul | Tanggung jawab di luar M1 | Hubungan dengan M1 |
|---|---|---|
| M2 – Unit Usaha/Dashboard | Mengelola data unit usaha dan kebutuhan dashboard sesuai ruang lingkup M2. | Menjadi pemilik data `unit_usaha_id`; M1 hanya menyimpan referensinya bila diperlukan. |
| M3 – Partner | Pendaftaran, verifikasi, persetujuan, penolakan, penonaktifan, dan aktivasi kembali partner. | Mengirim data akun awal Partner Admin dan event status partner ke M1. M3 adalah sumber data utama partner. |
| M4 – Operasional | Mengelola produk/layanan, pelanggan, quotation, invoice, pembayaran, dan status pesanan. | Mengirim daftar permission M4 kepada M1 secara manual; memvalidasi hak akses pada API M4. |
| M5 – Katalog Publik | Mengelola katalog publik, permintaan layanan, keranjang, checkout, serta pencarian/filter sesuai ruang lingkup M5. | Mengirim daftar permission M5 kepada M1 secara manual; memvalidasi hak akses pada API M5. |

Modul lain bertanggung jawab melindungi endpoint API miliknya sendiri dengan memvalidasi JWT dan permission yang relevan. M1 menyediakan identitas, token, serta mekanisme penetapan hak akses; M1 tidak otomatis mengamankan endpoint milik modul lain.

### 5.3 Rancangan endpoint M1

Endpoint berikut mengacu pada OpenAPI M1 yang tersedia dan penyesuaian integrasi yang diperlukan. Endpoint tambahan merupakan rancangan kontrak yang perlu disepakati dengan modul terkait dan dicantumkan dalam OpenAPI.

| Method | Endpoint | Kegunaan |
|---|---|---|
| POST | `/api/m1/v1/auth/login` | Login pengguna. |
| GET | `/api/m1/v1/me` | Mendapatkan profil pengguna yang sedang login. |
| POST | `/api/m1/v1/auth/change-password` | Mengubah kata sandi. |
| GET | `/api/m1/v1/users` | Mendapatkan daftar pengguna. |
| POST | `/api/m1/v1/users` | Membuat akun pengguna oleh administrator berwenang. |
| GET | `/api/m1/v1/users/{id}` | Mendapatkan detail pengguna. |
| PATCH | `/api/m1/v1/users/{id}` | Memperbarui data pengguna. |
| DELETE | `/api/m1/v1/users/{id}` | Menghapus/menonaktifkan pengguna sesuai definisi kontrak; utamakan soft delete bila riwayat harus dipertahankan. |
| PATCH | `/api/m1/v1/users/{id}/status` | Mengubah status aktif pengguna. |
| GET/POST | `/api/m1/v1/roles` | Melihat dan membuat role. |
| GET/PATCH/DELETE | `/api/m1/v1/roles/{id}` | Melihat, memperbarui, dan menghapus role sesuai aturan. |
| PATCH | `/api/m1/v1/roles/{id}/status` | Mengubah status role. |
| GET | `/api/m1/v1/permissions` | Melihat katalog permission yang telah dicatat di M1. |
| GET/POST | `/api/m1/v1/divisions` | Melihat dan membuat divisi apabila kepemilikan CRUD divisi telah disepakati sebagai tanggung jawab M1. |
| GET/PATCH/DELETE | `/api/m1/v1/divisions/{id}` | Detail, perubahan, dan penghapusan divisi apabila disepakati. |
| PATCH | `/api/m1/v1/divisions/{id}/status` | Mengubah status divisi apabila disepakati. |
| GET | `/api/m1/v1/audit-logs` | Melihat audit log dalam lingkup kewenangan pengguna. |
| GET | `/api/m1/v1/platform/audit-logs` | Melihat audit log platform oleh Super Admin. |
| POST | `/api/m1/v1/service/partners/initial-account` | **Endpoint tambahan:** menerima data akun awal Partner Admin dari M3 menggunakan `X-Service-Key`. |
| POST | `/api/m1/v1/events` | **Endpoint tambahan:** menerima event status partner dari M3 menggunakan `X-Service-Key`. |

**Catatan integrasi M3–M1:**
- M3 mengirim data akun awal setelah partner disetujui. Data minimal yang disepakati: `name`, `email`, `password`, `role`, `partner_id`, dan `unit_usaha_id`.
- M1 melakukan validasi dan hash password sebelum menyimpan akun. Password tidak boleh disimpan dalam bentuk plaintext.
- Event yang diusulkan: `partner.deactivated` dan `partner.reactivated`, dengan `partner_id` sebagai identitas partner.
- Status akun pengguna dan status partner harus dibedakan. Partner yang aktif kembali tidak otomatis mengaktifkan akun yang sebelumnya dinonaktifkan secara individual.
- Pemrosesan event perlu idempotent agar pengiriman ulang event tidak menyebabkan perubahan ganda.
- Nama endpoint dan struktur payload final harus disepakati M1–M3.

### 5.4 Aturan pembuatan API

Aturan berikut mengikuti Panduan API Bersama dan perlu diterapkan pada endpoint M1:

1. **Base path:** `/api/{modul}/v1/{resource}`. Untuk M1 gunakan `/api/m1/v1/...`.
2. **HTTP method:** gunakan method REST sesuai operasi, seperti `GET`, `POST`, `PATCH`, dan `DELETE`.
3. **Format data:** request dan response menggunakan JSON; nama field menggunakan `snake_case`.
4. **Autentikasi pengguna:** gunakan header `Authorization: Bearer <token>`.
5. **Autentikasi antarmodul:** gunakan header `X-Service-Key` untuk komunikasi backend antarmodul.
6. **JWT:** M1 menjadi penerbit token; modul lain memvalidasi token yang sama. Klaim dan masa berlaku harus mengikuti kontrak bersama. Contoh klaim dalam panduan mencakup `user_id`, `name`, `role`, `division`, `unit_usaha_id`, dan `exp`; kebutuhan `partner_id` serta `permissions` harus diselaraskan pada kontrak final.
7. **Respons:** gunakan envelope konsisten `success`, `data`, dan `error`. Respons daftar dapat menyertakan `meta`.
8. **Status HTTP:** `401` untuk token/kredensial tidak valid atau kedaluwarsa; `403` untuk pengguna terautentikasi tetapi tidak memiliki hak akses; gunakan status lain sesuai hasil operasi.
9. **Referensi lintas modul:** gunakan field `{data}_id`, misalnya `partner_id` dan `unit_usaha_id`. ID lintas database merupakan referensi, bukan foreign key lintas database.
10. **Event:** gunakan format event bersama yang memuat `event_id`, `event_type`, `source_module`, `occurred_at`, `recipient_role`, dan `payload`. Penerima perlu mencegah pemrosesan duplikat berdasarkan `event_id`.
11. **Retry event:** pengirim melakukan percobaan ulang sesuai panduan bersama (1, 5, dan 30 menit) dan mencatat kegagalan.
12. **OpenAPI:** semua endpoint dan skema request/response harus didokumentasikan. Perubahan kontrak perlu melalui proses review bersama konsumen API.
13. **Contract test:** endpoint yang menjadi kontrak antarmodul perlu diuji kesesuaiannya dengan OpenAPI.
14. **Keamanan:** jangan mencatat password, token, atau secret ke audit log. Endpoint service harus dilindungi dengan `X-Service-Key` dan hanya diakses melalui koneksi aman HTTPS.

## 6. Di Luar Cakupan

1. Pendaftaran, verifikasi, persetujuan, penolakan, dan pengelolaan data utama partner — tanggung jawab M3.
2. Pengelolaan data utama unit usaha — tanggung jawab M2.
3. Pengelolaan produk, layanan, pelanggan, quotation, invoice, pembayaran, dan pesanan — tanggung jawab M4.
4. Pengelolaan katalog publik, permintaan layanan, keranjang, dan checkout — tanggung jawab M5.
5. Implementasi login terpisah pada setiap modul — autentikasi dipusatkan di M1.
6. Pembuatan atau pendaftaran permission baru secara mandiri melalui UI/API M1 oleh modul lain. Daftar permission disampaikan oleh masing-masing modul kepada tim M1 secara manual melalui dokumentasi/komunikasi yang disepakati.
7. Penerapan pemeriksaan permission pada endpoint bisnis milik M2–M5 — menjadi tanggung jawab backend masing-masing modul.
8. Pengelolaan akun customer guest untuk transaksi publik, kecuali ada keputusan baru yang memasukkannya ke ruang lingkup M1.
9. Pembuatan fitur bisnis yang tidak tercantum dalam user story dan kebutuhan M1.

## 7. Kriteria Penerimaan

### Autentikasi dan profil
- **AC-01:** Pengguna dapat login menggunakan kredensial yang valid.
- **AC-02:** Kredensial yang salah ditolak dan tidak menghasilkan token.
- **AC-03:** Respons login menghasilkan token sesuai kontrak API bersama.
- **AC-04:** Pengguna dapat mengambil profil sendiri melalui `GET /me` dengan token valid.
- **AC-05:** Request tanpa token, token tidak valid, atau token kedaluwarsa ditolak dengan `401`.

### Pengelolaan akun
- **AC-06:** Administrator berwenang dapat melihat daftar akun sesuai lingkup kewenangannya.
- **AC-07:** Administrator berwenang dapat membuat akun dengan data wajib yang valid.
- **AC-08:** Email akun harus unik sesuai aturan sistem.
- **AC-09:** Data akun dapat diperbarui oleh administrator berwenang.
- **AC-10:** Akun dapat dinonaktifkan tanpa menghilangkan riwayat yang diwajibkan untuk dipertahankan.
- **AC-11:** Akun nonaktif tidak dapat login.
- **AC-12:** Pengguna tanpa kewenangan tidak dapat mengelola akun di luar lingkupnya.

### Role dan permission
- **AC-13:** Pengguna berwenang dapat melihat daftar role.
- **AC-14:** Role dapat dibuat atau diperbarui sesuai aturan role yang disepakati.
- **AC-15:** M1 menampilkan katalog permission yang telah disampaikan oleh modul terkait.
- **AC-16:** M1 tidak menyediakan endpoint/UI pendaftaran permission baru oleh modul lain.
- **AC-17:** Administrator berwenang dapat menetapkan permission dari katalog kepada role dan/atau divisi.
- **AC-18:** Pengguna hanya memperoleh hak akses sesuai penetapan yang berlaku.
- **AC-19:** Modul pemilik fitur tetap memvalidasi permission pada API bisnisnya sendiri.

### Perubahan kata sandi
- **AC-20:** Pengguna dapat mengganti kata sandi dengan memasukkan kata sandi lama yang benar.
- **AC-21:** Kata sandi lama yang salah menyebabkan permintaan ditolak.
- **AC-22:** Kata sandi baru harus memenuhi kebijakan kata sandi yang disepakati.
- **AC-23:** Kata sandi disimpan dalam bentuk hash, bukan plaintext.

### Audit log
- **AC-24:** Aktivitas penting pengelolaan akun dan hak akses di M1 tercatat.
- **AC-25:** Catatan audit memuat identitas pelaku, tindakan, waktu, dan informasi relevan.
- **AC-26:** Pengguna hanya dapat melihat audit log sesuai lingkup kewenangannya.
- **AC-27:** Audit log dapat ditampilkan secara terurut dan mendukung pagination bila disediakan dalam kontrak.
- **AC-28:** Password, token, dan secret tidak dicatat ke audit log.

### Integrasi M3–M1
- **AC-29:** Endpoint pembuatan akun awal hanya dapat dipanggil menggunakan autentikasi service yang valid.
- **AC-30:** M1 menerima data akun awal Partner Admin dari M3 setelah partner disetujui.
- **AC-31:** M1 memvalidasi field wajib dan keunikan email sebelum membuat akun.
- **AC-32:** M1 menyimpan `partner_id` sebagai referensi dari M3, bukan foreign key ke database M3.
- **AC-33:** Password yang dikirim M3 diproses menjadi hash sebelum disimpan.
- **AC-34:** Pengiriman ulang permintaan yang sama tidak membuat akun duplikat.
- **AC-35:** Endpoint event hanya dapat dipanggil menggunakan autentikasi service yang valid.
- **AC-36:** Event `partner.deactivated` menyebabkan akun terkait partner tersebut tidak dapat mengakses sistem.
- **AC-37:** Event `partner.reactivated` memulihkan kelayakan akses akun yang masih aktif secara individual.
- **AC-38:** Akun yang dinonaktifkan secara individual tetap nonaktif meskipun partner diaktifkan kembali.
- **AC-39:** Event dengan `event_id` yang telah diproses tidak diterapkan ulang.
- **AC-40:** Kegagalan pemrosesan event dapat ditelusuri melalui log dan mengikuti mekanisme retry yang disepakati.

### API dan integrasi
- **AC-41:** Endpoint M1 menggunakan prefix `/api/m1/v1`.
- **AC-42:** Nama field JSON menggunakan `snake_case`.
- **AC-43:** Endpoint pengguna menggunakan Bearer token, sedangkan endpoint service antarmodul menggunakan `X-Service-Key`.
- **AC-44:** Respons API mengikuti envelope bersama.
- **AC-45:** Semua endpoint M1 yang tersedia tercantum dalam OpenAPI dengan request, response, dan status code.
- **AC-46:** Perubahan kontrak antarmodul dikomunikasikan dan disetujui konsumen API terkait.
- **AC-47:** Endpoint integrasi M3–M1 diuji berdasarkan kontrak yang telah disepakati.
- **AC-48:** Pemeriksaan kesiapan layanan (`/ready`) dapat memastikan koneksi database tersedia.

## Catatan Keputusan yang Harus Dikonfirmasi

Sebelum implementasi final, sepakati hal-hal berikut bersama PO/dosen dan modul terkait:

1. **Istilah partner:** seluruh penggunaan `organization` pada OpenAPI lama diganti menjadi `partner` karena keduanya merujuk entitas yang sama dalam proyek ini. Nama field dan path API harus konsisten, misalnya `partner_id` dan `/partners`.
2. **Kepemilikan divisi:** M2 memiliki kebutuhan pengelolaan struktur unit usaha/divisi, sedangkan OpenAPI M1 menyediakan endpoint divisi. Tentukan satu pemilik master data agar tidak terjadi duplikasi.
3. **Role global atau partner:** panduan API menyebut daftar role bersama dikelola M1, sedangkan kebutuhan platform dapat mengarah pada role yang dapat dibuat partner. Tentukan apakah role bersifat global, partner-specific, atau kombinasi.
4. **Klaim JWT:** selaraskan `partner_id`, `unit_usaha_id`, `division`, `role`, dan `permissions` antara panduan API, OpenAPI, serta implementasi.
5. **Pencabutan akses token:** tentukan mekanisme agar token yang sudah terbit tidak tetap dapat digunakan setelah akun atau partner dinonaktifkan. Panduan menyebut modul memvalidasi JWT tanpa selalu meminta M1 pada setiap request; karena itu mekanisme revokasi, versi token, atau strategi masa berlaku perlu disepakati.
6. **Audit lintas modul:** M1 mencatat aktivitas pengelolaan akun/hak akses di M1. M2–M5 mencatat aktivitas bisnis pada modul masing-masing, kecuali nantinya disepakati kontrak event audit terpusat.
7. **Super Admin awal:** tentukan mekanisme pembuatan akun Super Admin pertama yang aman.
8. **Endpoint tambahan:** endpoint `POST /service/partners/initial-account` dan `POST /events` adalah rancangan tambahan untuk kebutuhan integrasi M3–M1; skema final wajib disetujui kedua tim dan dimasukkan ke OpenAPI bersama.

## Referensi

- Dokumen Studi Kasus Capstone Project – Kelas A.
- Panduan API Bersama v1.0 (draf).
- OpenAPI Modul M1 (`araya-business-hub-api-modul-m1.yaml` dan JSON terkait).
