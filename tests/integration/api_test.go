package integration

import (
	"api-stori/tests/test_utils"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

// Note: setupTestServer is now centralized in test_utils package

func TestHealthEndpoint(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/health")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Expected no error decoding JSON, got %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", result["status"])
	}
}

func TestRootEndpoint(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Expected no error decoding JSON, got %v", err)
	}

	if result["message"] == nil {
		t.Error("Expected message field in response")
	}
}

func TestMigrateEndpoint(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Create a test CSV file
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,-75.25,2024-01-15 14:45:00
3,1002,200.00,2024-01-16 09:15:00`

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add CSV file
	fileWriter, err := writer.CreateFormFile("csv_file", "test.csv")
	if err != nil {
		t.Fatalf("Expected no error creating form file, got %v", err)
	}
	fileWriter.Write([]byte(csvContent))
	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", server.URL+"/api/v1/migrate", &buf)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Expected no error decoding JSON, got %v", err)
	}

	if result["success"] != true {
		t.Error("Expected success to be true")
	}

	data := result["data"].(map[string]interface{})
	if data["total_records"] != float64(3) {
		t.Errorf("Expected 3 total records, got %v", data["total_records"])
	}
}

func TestMigrateEndpointInvalidFile(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Create multipart form data with invalid file
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add invalid file
	fileWriter, err := writer.CreateFormFile("csv_file", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error creating form file, got %v", err)
	}
	fileWriter.Write([]byte("not a csv file"))
	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", server.URL+"/api/v1/migrate", &buf)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestMigrateEndpointNoFile(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Make request without file
	req, err := http.NewRequest("POST", server.URL+"/api/v1/migrate", nil)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestBalanceEndpoint(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First, migrate some data
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,-75.25,2024-01-15 14:45:00
3,1002,200.00,2024-01-16 09:15:00`

	// Migrate data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, _ := writer.CreateFormFile("csv_file", "test.csv")
	fileWriter.Write([]byte(csvContent))
	writer.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/v1/migrate", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	client.Do(req)

	// Test balance endpoint
	resp, err := http.Get(server.URL + "/api/v1/users/1001/balance")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Expected no error decoding JSON, got %v", err)
	}

	expectedBalance := 150.50 - 75.25
	if result["balance"] != expectedBalance {
		t.Errorf("Expected balance %.2f, got %v", expectedBalance, result["balance"])
	}

	if result["total_debits"] != float64(1) {
		t.Errorf("Expected 1 debit, got %v", result["total_debits"])
	}

	if result["total_credits"] != float64(1) {
		t.Errorf("Expected 1 credit, got %v", result["total_credits"])
	}
}

func TestBalanceEndpointWithDateRange(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First, migrate some data
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,-75.25,2024-01-16 14:45:00
3,1001,200.00,2024-01-17 09:15:00`

	// Migrate data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, _ := writer.CreateFormFile("csv_file", "test.csv")
	fileWriter.Write([]byte(csvContent))
	writer.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/v1/migrate", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	client.Do(req)

	// Test balance endpoint with date range
	url := server.URL + "/api/v1/users/1001/balance?from=2024-01-16T00:00:00Z&to=2024-01-17T23:59:59Z"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Expected no error decoding JSON, got %v", err)
	}

	// Should only include transactions from Jan 16-17: -75.25 + 200.00 = 124.75
	expectedBalance := -75.25 + 200.00
	if result["balance"] != expectedBalance {
		t.Errorf("Expected balance %.2f, got %v", expectedBalance, result["balance"])
	}
}

func TestBalanceEndpointUserNotFound(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/users/9999/balance")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestBalanceEndpointInvalidUserID(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/users/invalid/balance")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestBalanceEndpointInvalidDateFormat(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/users/1001/balance?from=invalid-date")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestBalanceEndpointInvalidDateRange(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/users/1001/balance?from=2024-01-20T00:00:00Z&to=2024-01-15T23:59:59Z")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestNotFoundEndpoint(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/nonexistent")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := test_utils.SetupTestServer()
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/v1/health", "application/json", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", resp.StatusCode)
	}
}

// Helper function to create a test CSV file
func createTestCSVFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatalf("Expected no error creating temp file, got %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Expected no error writing to temp file, got %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

// Helper function to create multipart form data
func createMultipartFormData(t *testing.T, fieldName, fileName, content string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("Expected no error creating form file, got %v", err)
	}

	fileWriter.Write([]byte(content))
	writer.Close()

	return &buf, writer.FormDataContentType()
}
