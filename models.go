package main

const (
	SlotMin     = 1
	SlotMax     = 8
	MaxNama     = 50
	MaxDurasi   = 3
	MaxKuota    = 2
	StatusAktif = "aktif"
	JenisRuang  = "ruang"
	JenisAlat   = "alat"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

type Kuota struct {
	Peminjam string `json:"peminjam"`
	Tanggal  string `json:"tanggal"`
	Terpakai int    `json:"terpakai"`
	Maksimal int    `json:"maksimal"`
	Sisa     int    `json:"sisa"`
}

type apiError struct {
	Status int
	Msg    string
}

func (e *apiError) Error() string { return e.Msg }

type RingkasanSumberDaya struct {
	SumberDayaID     int    `json:"sumber_daya_id"`
	Nama             string `json:"nama"`
	JumlahPeminjaman int    `json:"jumlah_peminjaman"`
}

type Ringkasan struct {
	Tanggal          string                `json:"tanggal"`
	JumlahPeminjaman int                   `json:"jumlah_peminjaman"`
	TotalUnitSlot    int                   `json:"total_unit_slot"`
	PerSumberDaya    []RingkasanSumberDaya `json:"per_sumber_daya"`
}
