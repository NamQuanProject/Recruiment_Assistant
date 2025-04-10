package main

import (
	"log"

	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"
)

func main() {
	// Create and run the web server
	server := highlight.NewWebServer()
	log.Println("Starting web server...")
	server.Run()
} 