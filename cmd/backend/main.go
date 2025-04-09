package main

import (
	"context"
	"os"
	"os/signal"
	"log"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"
	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
)

func main() {
	
	log.Println("Starting AI servers...")
	go func() {
		log.Println("[API Gateway] Starting on :8080...")
		apigateway.RunServer()
	}()

	go func() {
		log.Println("[AI Services] Starting on :8081...")
		aiservices.RunServer()
	}()

	log.Println("Starting Parsing server...")
	go func() {
		log.Println("[Parsing Service] Starting on :8085...")
		parsing.RunServer()
	}()
	
	log.Println("Starting Output server...")
	go func() {
		log.Println("[Output Service] Starting on :8084...")
		output.RunServer()
	}()

	log.Println("Starting Highlight server...")
	go func() {
		log.Println("[Highlight Service] Starting on :8083...")
		highlight.RunServer()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("🔴 Shutdown signal received. Stopping all servers...")

	// Optional: Here you could call shutdown logic for each service if needed
}
