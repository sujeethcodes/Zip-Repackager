package main

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"zip-repackager/entity"
)

func main() {
	if len(os.Args) != 3 {
		exitWithError("Usage: repackager <input.zip> <output.zip>")
	}

	inputZip := os.Args[1]
	outputZip := os.Args[2]

	fileMap, err := processInputZip(inputZip)
	if err != nil {
		exitWithError(err.Error())
	}

	err = createOutputZip(outputZip, fileMap)
	if err != nil {
		exitWithError(err.Error())
	}
}

func processInputZip(inputZip string) (map[string]entity.FileData, error) {
	reader, err := zip.OpenReader(inputZip)
	if err != nil {
		return nil, fmt.Errorf("Failed to open input zip: %v", err)
	}
	defer reader.Close()

	fileMap := make(map[string]entity.FileData)

	for _, f := range reader.File {
		if f.FileInfo().IsDir() || f.Mode()&os.ModeSymlink != 0 {
			continue // skip directories and symlinks
		}

		fileOnlyName := filepath.Base(f.Name)

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("Failed to open file %s: %v", f.Name, err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("Failed to read file %s: %v", f.Name, err)
		}

		hash := sha256.Sum256(content)
		fileSize := int64(len(content))

		existing, exists := fileMap[fileOnlyName]
		if !exists || fileSize > existing.FileSize {
			fileMap[fileOnlyName] = entity.FileData{
				Name:     fileOnlyName,
				Content:  content,
				FileSize: fileSize,
				SHA256:   hash,
			}
		}
	}
	return fileMap, nil
}

func createOutputZip(outputZip string, fileMap map[string]entity.FileData) error {
	outFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("Failed to create output zip: %v", err)
	}
	defer outFile.Close()

	writer := zip.NewWriter(outFile)
	defer writer.Close()

	for _, file := range fileMap {
		header := &zip.FileHeader{
			Name:   file.Name,
			Method: zip.Store, // no compression
		}
		header.SetMode(0644)

		w, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("Failed to write file %s: %v", file.Name, err)
		}
		_, err = w.Write(file.Content)
		if err != nil {
			return fmt.Errorf("Failed to write content for %s: %v", file.Name, err)
		}

		// Recompute hash and verify
		newHash := sha256.Sum256(file.Content)
		if newHash != file.SHA256 {
			return fmt.Errorf("Hash mismatch for file %s", file.Name)
		}
	}
	return nil
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
