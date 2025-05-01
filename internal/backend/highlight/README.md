# CV Highlight Server

This server is responsible for highlighting weak areas in CVs. It works in conjunction with the AI server to identify weak areas and then highlights them in the PDF.

## Features

- Receives a PDF file and a list of weak areas
- Highlights the weak areas in the PDF with yellow annotations
- Adds popup notes with descriptions of the weak areas
- Returns the path to the highlighted PDF

## API

### POST /highlight

Highlights weak areas in a PDF file.

#### Request

```json
{
  "pdf_path": "path/to/pdf/file.pdf",
  "weak_areas": [
    {
      "text": "Weak area text",
      "page": 1,
      "x": 100,
      "y": 200,
      "width": 200,
      "height": 50,
      "description": "This area is weak because..."
    }
  ]
}
```

#### Response

```json
{
  "highlighted_pdf_path": "path/to/highlighted/pdf/file.pdf",
  "message": "PDF highlighted successfully"
}
```

## Usage

### Running the Server

```bash
go run cmd/highlight/main.go
```

The server will start on port 8083.

### Using the Client

```go
import "github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"

// Create a highlight client
highlightClient := highlight.NewClient("http://8083")

// Highlight a PDF
highlightResp, err := highlightClient.HighlightPDF("path/to/pdf/file.pdf", weakAreas)
if err != nil {
    // Handle error
}

// Get the path to the highlighted PDF
highlightedPDFPath := highlightResp.HighlightedPDFPath
```

### Processing a CV and Highlighting Weak Areas

```go
import "github.com/KietAPCS/test_recruitment_assistant/internal/backend/highlight"

// Process a CV and highlight weak areas
highlightedPDFPath, err := highlight.ProcessCVAndHighlight(
    "path/to/cv.pdf",
    "Software Engineer",
    "We are looking for a software engineer with 5+ years of experience in Go and Python.",
    "http://localhost:8081",  // AI server URL
    "http://localhost:8083"   // Highlight server URL
)
if err != nil {
    // Handle error
}

// Get the path to the highlighted PDF
fmt.Printf("CV processed and highlighted. Output saved to: %s\n", highlightedPDFPath)
```

## Dependencies

- PyMuPDF (fitz) - For PDF manipulation
- Gin - For the HTTP server
