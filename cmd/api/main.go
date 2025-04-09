package main

import (
	"sync"

	"github.com/KietAPCS/test_recruitment_assistant/internal/aiservices"
	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway"
	"github.com/KietAPCS/test_recruitment_assistant/internal/backend/parsing"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		apigateway.RunServer()
	}()

	go func() {
		defer wg.Done()
		parsing.RunServer()
	}()

	go func() {
		defer wg.Done()
		aiservices.RunServer()
	}()

	wg.Wait() // Blocks main until both servers exit
}
