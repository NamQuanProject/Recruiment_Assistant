package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SubmitCVsHandler(c *gin.Context) {
	currentPathBytes, err := os.ReadFile("storage/current.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read storage/current.txt"})
		return
	}
	basePath := strings.TrimSpace(string(currentPathBytes))
	cvsFolder := filepath.Join(basePath, "original", "cvs")
	parseCvs := filepath.Join(basePath, "parse", "cvs")

	if err := os.MkdirAll(cvsFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create folder to store CVs"})
		return
	}
	if err := os.MkdirAll(parseCvs, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create parse folder"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid file"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if ext == ".pdf" {
		dst := filepath.Join(cvsFolder, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDF file"})
			return
		}

		parsedst := filepath.Join(parseCvs, strings.TrimSuffix(file.Filename, ".pdf")+".txt")
		if err := processCV(dst, parsedst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse PDF"})
			return
		}

		if err := evaluate(basePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate CV: " + err.Error()})
			return
		}
		// c.JSON(http.StatusOK, gin.H{"message": "PDF uploaded and parsed successfully", "path": dst})
		fmt.Print("message", "PDF uploaded and parsed successfully")
		fmt.Print("path", dst)

	} else if ext == ".zip" {
		tempZipPath := filepath.Join(os.TempDir(), file.Filename)
		if err := c.SaveUploadedFile(file, tempZipPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ZIP file :" + err.Error()})
			return
		}

		zipReader, err := zip.OpenReader(tempZipPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read ZIP file"})
			return
		}
		defer zipReader.Close()

		var savedFiles []string
		for _, f := range zipReader.File {
			if strings.HasSuffix(strings.ToLower(f.Name), ".pdf") {
				fileName := filepath.Base(f.Name)
				outPath := filepath.Join(cvsFolder, fileName)

				rc, err := f.Open()
				if err != nil {
					continue
				}

				outFile, err := os.Create(outPath)
				if err != nil {
					rc.Close()
					continue
				}

				if _, err := io.Copy(outFile, rc); err != nil {
					outFile.Close()
					rc.Close()
					continue
				}
				outFile.Close()
				rc.Close()

				parsedst := filepath.Join(parseCvs, strings.TrimSuffix(fileName, ".pdf")+".txt")
				if err := processCV(outPath, parsedst); err != nil {
					continue
				}

				savedFiles = append(savedFiles, outPath)
			}
		}

		if err := evaluate(basePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate CV"})
			return
		}

		// c.JSON(http.StatusOK, gin.H{
		// 	"message":    "ZIP processed",
		// 	"pdf_count":  len(savedFiles),
		// 	"saved_pdfs": savedFiles,
		// })
		fmt.Print("message", "ZIP processed")
		fmt.Print("pdf_count", len(savedFiles))
		fmt.Print("saved_pdfs", savedFiles)

	}

	req := struct {
		EvaluationFolder string `json:"evaluation_folder"`
	}{
		EvaluationFolder: filepath.Join(basePath, "evaluation"),
	}
	// call output
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare request"})
		return
	}
	URL := os.Getenv("OUTPUT_URL")
	if URL == "" {
		URL = "http://localhost:8080" // Default URL for local testing
	}
	// URL := "https://output-service.onrender.com"
	// resp, err := http.Post("http://localhost:8084/output", "application/json", bytes.NewBuffer(reqBody))
	resp, err := http.Post(fmt.Sprintf("%s/output", URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call evaluation server"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Evaluation server returned error: " + resp.Status})
		return
	}
	log.Printf("CVs evaluation successfully")
	c.JSON(http.StatusOK, gin.H{"final_out_path": "internal/backend/output/output.json"})
	// return
}

func processCV(pdfPath, outPath string) error {

	parseRequest := struct {
		InputPath  string `json:"input_path" binding:"required"`
		OutputPath string `json:"output_path" binding:"required"`
	}{
		InputPath:  pdfPath,
		OutputPath: outPath,
	}

	reqBody, err := json.Marshal(parseRequest)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}

	// Call parsing server
	URL := os.Getenv("PARSE_URL")
	if URL == "" {
		URL = "http://localhost:8080" // Default URL for local testing
	}
	// URL := "https://parsing-service.onrender.com"
	// resp, err := http.Post("http://localhost:8085/parse/cv", "application/json", bytes.NewBuffer(reqBody))
	resp, err := http.Post(fmt.Sprintf("%s/parse/cv", URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call parsing server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("parsing server returned error: %v", resp.Status)
	}

	log.Printf("JD parsed successfully")
	return nil
}

func evaluate(path string) error {

	evalRequest := struct {
		InputPath string `json:"input_path" binding:"required"`
	}{
		InputPath: path,
	}

	reqBody, err := json.Marshal(evalRequest)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %v", err)
	}

	URL := os.Getenv("EVAL_URL")
	if URL == "" {
		URL = "http://localhost:8080" // Default URL for local testing
	}
	// URL := "https://evaluation-service-dytd.onrender.com"

	// resp, err := http.Post("http://localhost:8082/evaluate", "application/json", bytes.NewBuffer(reqBody))
	resp, err := http.Post(fmt.Sprintf("%s/evaluate", URL), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to call evaluation server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("evaluation server returned error: %v", resp.Status)
	}

	log.Printf("CVs evalutaion successfully")
	return nil
}
