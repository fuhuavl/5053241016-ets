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
