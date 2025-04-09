package main

import (
	"context"
	"os"
	"os/signal"
	"log"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"
	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/output"
)

func main() {
	
	log.Println("Starting AI servers...")
	go func() {
		aiservices.RunServer()
	}()

	log.Println("Starting Parsing server...")
	go func() {
		parsing.RunServer()
	}()
	
	log.Println("Starting Output server...")
	go func() {
		output.RunServer()
	}()

	log.Println("Starting Highlight server...")
	go func() {
		highlight.RunServer()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	log.Println("Servers shutting down...")
}
