package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
    // Define the path to the React build folder
    frontendDir := "D:/Projects/Milestones/Second Year/Second Semester/Hackathon 2025 [GDGoC]/source/Recruiment_Assistant/web/build"

    // Check if the build folder exists
    if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
        log.Fatalf("Build folder not found. Please run 'npm run build' in the 'web' folder first.")
    }

    // Serve static files from the React build folder
    fs := http.FileServer(http.Dir(frontendDir))
    http.Handle("/", fs)

    // Start the server on port 3000
    log.Println("Starting frontend server on http://localhost:3000")
    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}