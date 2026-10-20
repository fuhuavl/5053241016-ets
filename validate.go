package main

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

func badRequest(msg string) *apiError { return &apiError{http.StatusBadRequest, msg} }

type inputSumberDaya struct {
	Nama  *string `json:"nama"`
	Jenis *string `json:"jenis"`
	Stok  *int    `json:"stok"`
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
