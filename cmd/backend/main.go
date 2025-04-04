package main

import (
	"log"
	"os"
	"os/signal"
	"context"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/evaluation"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
)

func main() {
	// Start servers concurrently
	go func() {
		evaluation.RunServer()
	}()
	go func() {
		output.RunServer()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	log.Println("Servers shutting down...")
}