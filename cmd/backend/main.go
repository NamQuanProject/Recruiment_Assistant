package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
)

func main() {
	
	go func() {
		aiservices.RunServer()
	}()

	go func() {
		parsing.RunServer()
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
