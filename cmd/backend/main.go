package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
)

func main() {
	// Start servers concurrently
	// go func() {
	// 	evaluation.RunServer()
	// }()
	// go func() {
	// 	output.RunServer()
	// }()
	go func() {
		aiservices.RunServer()
	}()

	// Wait for interrupt signal (Ctrl+C)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	log.Println("Servers shutting down...")
}
