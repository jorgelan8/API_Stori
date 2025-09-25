package test_utils

import (
	"api-stori/tests/config"
	"bytes"
	"net/http"
	"testing"
	"time"
)

// MigrateTestData migrates test data to the specified server
// This is a centralized function to avoid code duplication across tests
func MigrateTestData(t *testing.T, serverURL string, recordCount int) {
	csvData := GenerateTestCSV(recordCount)
	multipartData, contentType := CreateMultipartFormData(csvData, "test.csv")

	req, err := http.NewRequest("POST", serverURL+config.GetPathAPI()+"/migrate", bytes.NewReader(multipartData))
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
