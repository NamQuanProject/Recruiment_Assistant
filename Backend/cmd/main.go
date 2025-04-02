package main

import (
	"log"
	"net/http"

	"github.com/NamQuanProject/Recruiment_Assistant/Backend/internal/server"
)

func main() {
	srv := server.New()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Router))
}
