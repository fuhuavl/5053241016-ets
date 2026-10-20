package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("PinjamKampus jalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}
