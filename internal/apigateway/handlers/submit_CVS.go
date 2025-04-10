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
		c.JSON(http.StatusOK, gin.H{"message": "PDF uploaded and parsed successfully", "path": dst})

		return

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

		c.JSON(http.StatusOK, gin.H{
			"message":    "ZIP processed",
			"pdf_count":  len(savedFiles),
			"saved_pdfs": savedFiles,
		})

		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF or ZIP files are supported"})
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
	resp, err := http.Post("http://localhost:8085/parse/cv", "application/json", bytes.NewBuffer(reqBody))
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

	resp, err := http.Post("http://localhost:8082/evaluate", "application/json", bytes.NewBuffer(reqBody))
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
