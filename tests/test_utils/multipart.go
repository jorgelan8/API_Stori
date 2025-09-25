package test_utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

// CreateMultipartFormData creates multipart form data in a thread-safe way
// This is a centralized function to avoid code duplication across tests
func CreateMultipartFormData(csvData, filename string) ([]byte, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile("csv_file", filename)
	if err != nil {
		panic(fmt.Sprintf("Error creating form file: %v", err))
	}

	fileWriter.Write([]byte(csvData))
	writer.Close()

	return buf.Bytes(), writer.FormDataContentType()
}

// CreateMultipartFormDataPerRequest creates multipart form data for each individual request
// This is the most robust approach for concurrent testing
func CreateMultipartFormDataPerRequest(csvData, filename string) ([]byte, string) {
	// For now, we use the same implementation as CreateMultipartFormData
	// In the future, we could add request-specific optimizations
	return CreateMultipartFormData(csvData, filename)
}
