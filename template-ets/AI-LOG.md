# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Persiapan
- Prompt: minta persiapan belajar, lengkap dengan langkah eksekusi dan alur berpikirnya.
- Dipakai: hanya sebagai bahan belajar
- Diubah: -

## Inisisasi Setup
- Prompt: mengirim requirement PinjamKampus beserdanta kartu persona B, minta API dibuat sesuai kontrak dengan Go net/http, data di memori, disertai penjelasan alur berpikirnya
- Dipakai: kerangka awal: `go.mod`, `main.go`, `NewHandler()` di `handler.go`, `Store` dan `NewStore()` di `store.go`, tipe di `models.go`, serta helper `writeJSON`, `writeError`, `parseID`, `decodeBody`.
- Diubah: AI membuat file yang tidak disuruh seperti `main_test.go` dan saya hapus karena tidak disuruh

## B1: sumber daya
- Prompt: bagian kebutuhan B1 dan kontrak sumber daya
- Dipakai: `inputSumberDaya`, `validate()`, dan `validateNama()` di `validate.go`; `AddResource`, `ListResources`, `GetResource`, `findResource` di `store.go`; handler `createResource`, `listResources`, `getResource` di `handler.go`.
- Diubah: tidak ada

## B2-B3: buat peminjaman, daftar dengan filter, detail
- Prompt: bagian B2, B3, dan kontrak peminjaman
- Dipakai: `inputPeminjaman` dan `validate()` di `validate.go` (termasuk `validTanggal` dengan `time.Parse`); `CreateBooking`, `GetBooking`, `ListBookings`, `findBooking` di `store.go`; handler `createBooking`, `listBookings`, `getBooking` di `handler.go`.
- Diubah: tidak ada

## B5: deteksi bentrok
- Prompt: bagian B5 dan aturan bentrok
- Dipakai: `adaBentrok`, `menahanSlot`, dan `checkRules` di `store.go`. Pemeriksaan bentrok dan penyimpanan berada dalam satu kunci mutex di `CreateBooking`.
- Diubah: tidak ada

## B6: ubah dan hapus (sekaligus menjaga B7)
- Prompt: bagian B6 dan B7 
- Dipakai: `UpdateBooking`, `DeleteBooking` di `store.go`; handler `updateBooking`, `deleteBooking` di `handler.go`. Parameter `excludeID` pada `checkRules` membuat PUT tidak bentrok dengan dirinya sendiri, dan cek dilakukan sebelum data ditulis sehingga PUT yang gagal 409 tidak mengubah data lama. B7 dipenuhi oleh mutex di `Store` yang sudah ada sejak setup (tidak ada commit terpisah).
- Diubah: tidak ada 

## Fitur khas persona B: durasi maksimal 3 slot
- Prompt: kartu persona B 
- Dipakai: pemeriksaan durasi di `validate()` pada `validate.go` (menghasilkan 400 dan diperiksa sebelum bentrok) dan konstanta `MaxDurasi` di `models.go`.
- Diubah: tidak ada

## Fitur khas persona B: kuota harian dan endpoint sisa kuota
- Prompt: kartu persona B 
- Dipakai: `kuotaTerpakai` dan `Kuota` di `store.go`, pemeriksaan kuota (409) di `checkRules`, handler `getKuota` di `handler.go`, tipe `Kuota` dan konstanta `MaxKuota` di `models.go`.
- Diubah: tidak ada

## B4: error 404 dan 405 berbentuk JSON
- Prompt: saya minta commit khusus B4. Menurut kontrak B4, error harus selalu berbentuk JSON, sedangkan rute tidak dikenal dan method salah masih dijawab teks polos bawaan `net/http`.
- Dipakai: pembungkus di `NewHandler()` dan tipe `statusCatcher` di `handler.go`.
- Diubah: menyederhanakan pemanggilan `mux.Handler(r)` karena ada duplikat

## Rangkuman
- Prompt: rangkumkan semua prompt yang dikirim
- Dipakai: digunakan sebagai ringkasan dan patokan urutannya
- Diubah: tidak ada

## Change request file A
- Prompt: mengirim file A
- Dipakai: tipe Ringkasan dan RingkasanSumberDaya di models.go, fungsi Ringkasan di store.go, handler getRingkasan dan rute GET /ringkasan di handler.go
- Diubah: tidak ada 
