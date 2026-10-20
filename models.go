package main

const (
	SlotMin     = 1
	SlotMax     = 8
	MaxNama     = 50
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

type apiError struct {
	Status int
	Msg    string
}

func (e *apiError) Error() string { return e.Msg }
