package main

import (
	"sync"

	"github.com/KietAPCS/test_recruitment_assistant/internal/apigateway"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		apigateway.RunServer()
	}()

	wg.Wait() // Blocks main until both servers exit
}
