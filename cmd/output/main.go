package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
)

func main() {

	go func() {
		log.Println("[Output Service] Starting on :8084...")
		output.RunServer()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("🔴 Shutdown signal received. Stopping all servers...")

	// Optional: Here you could call shutdown logic for each service if needed
}
