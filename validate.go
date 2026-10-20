package main

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func badRequest(msg string) *apiError { return &apiError{http.StatusBadRequest, msg} }

type inputSumberDaya struct {
	Nama  *string `json:"nama"`
	Jenis *string `json:"jenis"`
	Stok  *int    `json:"stok"`
}

type inputPeminjaman struct {
	SumberDayaID *int    `json:"sumber_daya_id"`
	Peminjam     *string `json:"peminjam"`
	Tanggal      *string `json:"tanggal"`
	SlotMulai    *int    `json:"slot_mulai"`
	SlotSelesai  *int    `json:"slot_selesai"`
	Jumlah       *int    `json:"jumlah"`
}

func validateNama(field string, p *string) (string, *apiError) {
	if p == nil {
		return "", badRequest(field + " wajib diisi")
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return "", badRequest(field + " tidak boleh kosong")
	}
	if utf8.RuneCountInString(s) > MaxNama {
		return "", badRequest(field + " maksimal 50 karakter")
	}
	return s, nil
}

func (in inputSumberDaya) validate() (SumberDaya, *apiError) {
	nama, err := validateNama("nama", in.Nama)
	if err != nil {
		return SumberDaya{}, err
	}
	if in.Jenis == nil {
		return SumberDaya{}, badRequest("jenis wajib diisi")
	}
	jenis := strings.TrimSpace(*in.Jenis)
	if jenis != JenisRuang && jenis != JenisAlat {
		return SumberDaya{}, badRequest(`jenis harus "ruang" atau "alat"`)
	}
	if in.Stok == nil {
		return SumberDaya{}, badRequest("stok wajib diisi")
	}
	if *in.Stok < 1 {
		return SumberDaya{}, badRequest("stok minimal 1")
	}
	return SumberDaya{Nama: nama, Jenis: jenis, Stok: *in.Stok}, nil
}

func validTanggal(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func (in inputPeminjaman) validate() (Peminjaman, *apiError) {
	var p Peminjaman
	if in.SumberDayaID == nil {
		return p, badRequest("sumber_daya_id wajib diisi")
	}
	peminjam, err := validateNama("peminjam", in.Peminjam)
	if err != nil {
		return p, err
	}
	if in.Tanggal == nil || !validTanggal(strings.TrimSpace(*in.Tanggal)) {
		return p, badRequest("tanggal harus berformat YYYY-MM-DD dan valid")
	}
	if in.SlotMulai == nil || in.SlotSelesai == nil {
		return p, badRequest("slot_mulai dan slot_selesai wajib diisi")
	}
	mulai, selesai := *in.SlotMulai, *in.SlotSelesai
	if mulai < SlotMin || selesai > SlotMax || mulai > selesai {
		return p, badRequest("slot harus 1 sampai 8 dan slot_mulai <= slot_selesai")
	}
	if selesai-mulai+1 > MaxDurasi {
		return p, badRequest("durasi maksimal 3 slot")
	}
	if in.Jumlah == nil || *in.Jumlah < 1 {
		return p, badRequest("jumlah minimal 1")
	}
	return Peminjaman{
		SumberDayaID: *in.SumberDayaID,
		Peminjam:     peminjam,
		Tanggal:      strings.TrimSpace(*in.Tanggal),
		SlotMulai:    mulai,
		SlotSelesai:  selesai,
		Jumlah:       *in.Jumlah,
	}, nil
}
