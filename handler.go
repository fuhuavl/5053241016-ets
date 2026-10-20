package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type API struct {
	store *Store
}

func NewHandler() http.Handler {
	a := &API{store: NewStore()}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /sumber-daya", a.createResource)
	mux.HandleFunc("GET /sumber-daya", a.listResources)
	mux.HandleFunc("GET /sumber-daya/{id}", a.getResource)
	mux.HandleFunc("POST /peminjaman", a.createBooking)
	mux.HandleFunc("GET /peminjaman", a.listBookings)
	mux.HandleFunc("GET /peminjaman/{id}", a.getBooking)
	mux.HandleFunc("PUT /peminjaman/{id}", a.updateBooking)
	mux.HandleFunc("DELETE /peminjaman/{id}", a.deleteBooking)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeAPIError(w http.ResponseWriter, err *apiError) {
	writeError(w, err.Status, err.Msg)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka")
		return 0, false
	}
	return id, true
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return false
	}
	return true
}

func (a *API) createResource(w http.ResponseWriter, r *http.Request) {
	var in inputSumberDaya
	if !decodeBody(w, r, &in) {
		return
	}
	res, err := in.validate()
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, a.store.AddResource(res))
}

func (a *API) listResources(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.ListResources())
}

func (a *API) getResource(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	res, found := a.store.GetResource(id)
	if !found {
		writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *API) createBooking(w http.ResponseWriter, r *http.Request) {
	var in inputPeminjaman
	if !decodeBody(w, r, &in) {
		return
	}
	p, err := in.validate()
	if err != nil {
		writeAPIError(w, err)
		return
	}
	created, err := a.store.CreateBooking(p)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (a *API) listBookings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var sdID *int
	if q.Has("sumber_daya_id") {
		n, err := strconv.Atoi(q.Get("sumber_daya_id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus berupa angka")
			return
		}
		sdID = &n
	}
	writeJSON(w, http.StatusOK, a.store.ListBookings(sdID, q.Get("tanggal"), q.Get("status")))
}

func (a *API) getBooking(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	b, found := a.store.GetBooking(id)
	if !found {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (a *API) updateBooking(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if _, found := a.store.GetBooking(id); !found {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	var in inputPeminjaman
	if !decodeBody(w, r, &in) {
		return
	}
	p, err := in.validate()
	if err != nil {
		writeAPIError(w, err)
		return
	}
	updated, err := a.store.UpdateBooking(id, p)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *API) deleteBooking(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if !a.store.DeleteBooking(id) {
		writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
