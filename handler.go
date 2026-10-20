package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
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
