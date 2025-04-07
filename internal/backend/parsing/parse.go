// internal/backend/parsing/parse.go
package parsing

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExtractTextFromPDF(pdfPath string) (string, error) {
	pythonScriptPath := filepath.Join("internal", "backend", "parsing", "extract_pdf.py")

	outputPath := strings.TrimSuffix(pdfPath, ".pdf") + ".txt"
	cmd := exec.Command("python", pythonScriptPath, pdfPath, outputPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(err)
		return "", fmt.Errorf("error executing Python script: %v\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
