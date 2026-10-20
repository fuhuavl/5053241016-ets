# PROSES

## 1. Persona dan masalahnya
Persona B, mahasiswa yang meminjam alat kampus contohnya proyektor dan kabel. Masalahnya adalah ada orang yang meminjam terlalu banyak sehingga mahasiswa lain tidak kebagian. Yang dilakukan adalah mengimplementasi fitur yang membatasi tiap orang satu peminjaman maksimal 3 slot, maksimal 2 peminjaman per hari, dan ada endpoint untuk melihat sisa kuota

## 2. Rancangan struktur data
Semua data akan saya simpan di Store sebagai berikut:

- resources: daftar sumber daya (ruang atau alat) berupa slice. Saya pakai slice karena urutannya otomatis naik menurut id. Datanya sedikit, jadi mencari satu data cukup dengan loop.
- bookings: daftar peminjaman berupa slice
- nextResID dan nextBookID: dua penghitung id yang terpisah, karena soal minta id tiap jenis data mulai dari 1 sendiri-sendiri. Id tidak dipakai ulang setelah data dihapus.
- mu: kunci supaya data tidak rusak jika ada banyak request yang datang bersamaan

Untuk input dari pengguna, saya pakai tipe terpisah dengan pointer (*int, *string) supaya bisa membedakan field yang tidak dikirim dengan field kosong 

## Endpoint
| Method | Path | Sukses |
|---|---|---|
| POST | /sumber-daya | 201 |
| GET | /sumber-daya | 200 |
| GET | /sumber-daya/{id} | 200 |
| POST | /peminjaman | 201 |
| GET | /peminjaman?sumber_daya_id=&tanggal=&status= | 200 |
| GET | /peminjaman/{id} | 200 |
| PUT | /peminjaman/{id} | 200 |
| DELETE | /peminjaman/{id} | 204 |
| GET | /peminjam/{nama}/kuota?tanggal=YYYY-MM-DD | 200 |

## 3. Satu keputusan desain
1. cek bentrok dilakukan per slot, bukan membandingkan rentang

cara rentang hanya benar jika stoknya 1 (1 ruang hanya boleh dipakai satu 1). Untuk alat yang stoknya lebih dari 1, misalnya 3 proyektor, 3 orang boleh meminjam di slot yang sama selama totalnya tidak melebihi 3. jika cuma melihat rentang beririsan, peminjaman kedua sudah ditolak padahal stoknya masih cukup

2. penyimpanan dilakukan dalam satu kunci mutex. jika dikunci terpisah, dua request bisa sama-sama lolos pemeriksaan lalu sama-sama menyimpan, sehingga slot yang sama terpakai dua kali. Dengan satu kunci, request kedua baru memeriksa setelah request pertama selesai menyimpan

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: saat tes POST di postman respons selalu 400 Bad Request dengan pesan body JSON tidak valid, padahal servernya jalan
- Cara menemukan: ternyata saya belum menulis inputnya pada body
- Penyebab dan solusi: server menolak body kosong dengan 400 jadi ini bukan bug di kode. Solusinya mengisi body dengan memilih raw lalu JSON

### Bug 2
- Gejala: setelah menambah endpoint /ringkasan, requestnya tetap menghasilkan 404 rute tidak ditemukan, padahal kodenya sudah ada. Data hasil tes sebelumnya juga masih muncul
- Cara menemukan: saat saya run di port 8080 ternyata masih ada server lama yang masih run dari tes sebelumnya
- Penyebab dan solusi: server lama belum dimatikan dan setelah saya matikan dengan CTRL+C, server berjalan normal

## 5. Change request
tim frontend minta endpoint baru yaitu GET /ringkasan?tanggal=YYYY-MM-DD supaya browser tidak perlu mengambil semua peminjaman lalu menghitung sendiri. Oleh karena itu, perubahan yang saya buat adalah:

- models.go: menambah dua tipe untuk bentuk respons, yaitu Ringkasan dan RingkasanSumberDaya
- store.go: menambah fungsi Ringkasan yang menghitung angka-angkanya dari daftar peminjaman
- handler.go: menambah rute GET /ringkasan dan handler getRingkasan yang memeriksa tanggal lalu memanggil fungsi di store

ringkasan dihitung langsung dari daftar peminjaman setiap kali diminta, bukan disimpan sebagai angka terpisah, supaya angkanya selalu sesuai dengan data dan tidak perlu diperbarui di setiap POST, PUT, dan DELETE. Kekurangannya adalah mungkin akan lambat jika datanya ribuan.