package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SubmitCVs(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	files := form.File["cvs"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}
	filePath := filepath.Join("uploads", "cvs")

	var uploaded []string

	for _, file := range files {
		filename := file.Filename
		ext := strings.ToLower(filepath.Ext(filename))

		if ext == ".zip" {
			// Save zip temporarily
			tempZipPath := filepath.Join(filePath, "temp_"+filename)
			if err := c.SaveUploadedFile(file, tempZipPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save ZIP file"})
				return
			}
			// Extract zip
			names, err := extractZip(tempZipPath, filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract ZIP file"})
				return
			}
			uploaded = append(uploaded, names...)
			os.Remove(tempZipPath) // Clean up
		} else {
			// Save regular file
			dst := filepath.Join(filePath, filepath.Base(filename))
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file %s", filename)})
				return
			}
			uploaded = append(uploaded, filename)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "CVs uploaded successfully",
		"uploaded": uploaded,
	})
}

func extractZip(zipPath string, destDir string) ([]string, error) {
	var extracted []string

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		fpath := filepath.Join(destDir, f.Name)

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return nil, err
		}

		// Open file in zip
		inFile, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer inFile.Close()

		outFile, err := os.Create(fpath)
		if err != nil {
			return nil, err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, inFile); err != nil {
			return nil, err
		}

		extracted = append(extracted, f.Name)
	}
	return extracted, nil
}
