package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
)

func main() {
	// Start each service in a goroutine
	go func() {
		log.Println("[API Gateway] Starting on :8081...")
		apigateway.RunServer()
	}()

	go func() {
		log.Println("[AI Services] Starting on :8083...")
		aiservices.RunServer()
	}()

	go func() {
		log.Println("[Parsing Service] Starting on :8082...")
		parsing.RunServer()
	}()

	go func() {
		log.Println("[Output Service] Starting on :8083...")
		output.RunServer()
	}()

	// Wait for termination signal (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("🔴 Shutdown signal received. Stopping all servers...")

	// Optional: Here you could call shutdown logic for each service if needed
}
