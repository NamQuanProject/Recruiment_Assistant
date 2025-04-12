package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/evaluation"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
)

func main() {

	go func() {
		log.Println("[AI Services] Starting on :8081...")
		aiservices.RunServer()
	}()

	go func() {
		log.Println("[Evaluation] Starting on :8082...")
		evaluation.RunServer()
	}()

	go func() {
		log.Println("[Highlight Service] Starting on :8083...")
		highlight.RunServer()
	}()

	go func() {
		log.Println("[Output Service] Starting on :8084...")
		output.RunServer()
	}()

	go func() {
		log.Println("[Parsing Service] Starting on :8085...")
		parsing.RunServer()
	}()

	go func() {
		server := highlight.NewWebServer()
		log.Println("Starting web server...")
		server.Run()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("🔴 Shutdown signal received. Stopping all servers...")

	// Optional: Here you could call shutdown logic for each service if needed
}
