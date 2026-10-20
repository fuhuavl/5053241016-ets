package main

import (
	"net/http"
	"sync"
)

type Store struct {
	mu         sync.Mutex
	resources  []SumberDaya
	bookings   []Peminjaman
	nextResID  int
	nextBookID int
}

func NewStore() *Store {
	return &Store{nextResID: 1, nextBookID: 1}
}

func (s *Store) AddResource(r SumberDaya) SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.ID = s.nextResID
	s.nextResID++
	s.resources = append(s.resources, r)
	return r
}

func (s *Store) ListResources() []SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SumberDaya, len(s.resources))
	copy(out, s.resources)
	return out
}

func (s *Store) GetResource(id int) (SumberDaya, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.findResource(id)
}

func (s *Store) findResource(id int) (SumberDaya, bool) {
	for _, r := range s.resources {
		if r.ID == id {
			return r, true
		}
	}
	return SumberDaya{}, false
}

func (s *Store) findBooking(id int) int {
	for i, b := range s.bookings {
		if b.ID == id {
			return i
		}
	}
	return -1
}

func (s *Store) GetBooking(id int) (Peminjaman, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.findBooking(id); i >= 0 {
		return s.bookings[i], true
	}
	return Peminjaman{}, false
}

func (s *Store) ListBookings(sumberDayaID *int, tanggal, status string) []Peminjaman {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Peminjaman, 0)
	for _, b := range s.bookings {
		if sumberDayaID != nil && b.SumberDayaID != *sumberDayaID {
			continue
		}
		if tanggal != "" && b.Tanggal != tanggal {
			continue
		}
		if status != "" && b.Status != status {
			continue
		}
		out = append(out, b)
	}
	return out
}

func (s *Store) CreateBooking(p Peminjaman) (Peminjaman, *apiError) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkRules(p, 0); err != nil {
		return Peminjaman{}, err
	}
	p.ID = s.nextBookID
	s.nextBookID++
	p.Status = StatusAktif
	s.bookings = append(s.bookings, p)
	return p, nil
}

func (s *Store) UpdateBooking(id int, p Peminjaman) (Peminjaman, *apiError) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.findBooking(id)
	if i < 0 {
		return Peminjaman{}, &apiError{http.StatusNotFound, "peminjaman tidak ditemukan"}
	}
	if err := s.checkRules(p, id); err != nil {
		return Peminjaman{}, err
	}
	p.ID = id
	p.Status = s.bookings[i].Status
	s.bookings[i] = p
	return p, nil
}

func (s *Store) DeleteBooking(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.findBooking(id)
	if i < 0 {
		return false
	}
	s.bookings = append(s.bookings[:i], s.bookings[i+1:]...)
	return true
}

func menahanSlot(b Peminjaman) bool { return b.Status == StatusAktif }

func (s *Store) checkRules(p Peminjaman, excludeID int) *apiError {
	r, ok := s.findResource(p.SumberDayaID)
	if !ok {
		return badRequest("sumber_daya_id tidak ditemukan")
	}
	if p.Jumlah > r.Stok {
		return badRequest("jumlah melebihi stok sumber daya")
	}
	if s.adaBentrok(p, r.Stok, excludeID) {
		return &apiError{http.StatusConflict, "kapasitas slot terlampaui"}
	}
	return nil
}

func (s *Store) adaBentrok(p Peminjaman, stok, excludeID int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		terpakai := 0
		for _, b := range s.bookings {
			if b.ID == excludeID || !menahanSlot(b) {
				continue
			}
			if b.SumberDayaID == p.SumberDayaID && b.Tanggal == p.Tanggal &&
				b.SlotMulai <= slot && slot <= b.SlotSelesai {
				terpakai += b.Jumlah
			}
		}
		if terpakai+p.Jumlah > stok {
			return true
		}
	}
	return false
}
