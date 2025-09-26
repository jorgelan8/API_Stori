package load

import (
	"api-stori/tests/config"
	"api-stori/tests/test_utils"
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestLoadMigration tests the migration endpoint under load
// Uses pre-built multipart data for better memory efficiency
func TestLoadMigration(t *testing.T) {
	// Setup test server using centralized function
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Generate large CSV data using centralized function
	csvData := test_utils.GenerateTestCSV(100) // Reduced for testing

	// Test with multiple concurrent requests
	concurrency := 10          // Número de goroutines concurrentes
	requestsPerGoroutine := 25 // Requests por goroutine
	totalRequests := concurrency * requestsPerGoroutine

	start := time.Now()
	results := make(chan error, totalRequests)

	// Crear goroutines que ejecuten múltiples requests cada una
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			for j := 0; j < requestsPerGoroutine; j++ {
				// Create multipart data for each individual request (thread-safe)
				multipartData, contentType := test_utils.CreateMultipartFormDataPerRequest(csvData, "large_test.csv")

				req, err := http.NewRequest("POST", server.URL+config.GetPathAPI()+"/migrate", bytes.NewReader(multipartData))
				if err != nil {
					results <- err
					continue
				}
				req.Header.Set("Content-Type", contentType)

				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					results <- err
					continue
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					results <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
					continue
				}

				results <- nil
			}
		}(i)
	}

	// Collect results
	successCount := 0
	errorCount := 0
	t.Logf("Waiting for %d results...", totalRequests)

	for i := 0; i < totalRequests; i++ {
		select {
		case err := <-results:
			if err != nil {
				errorCount++
				t.Logf("Request failed: %v", err)
			} else {
				successCount++
				t.Logf("Request %d successful", i+1)
			}
		case <-time.After(30 * time.Second):
			t.Errorf("Test timed out after 30 seconds. Received %d/%d results", successCount+errorCount, totalRequests)
			return
		}
	}

	duration := time.Since(start)
	t.Logf("Load test completed in %v", duration)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Failed requests: %d", errorCount)
	t.Logf("Requests per second: %.2f", float64(successCount)/duration.Seconds())

	if errorCount > 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
}

// TestLoadBalance tests the balance endpoint under load
func TestLoadBalance(t *testing.T) {
	// Setup test server using centralized function
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First, migrate some data using centralized function
	test_utils.MigrateTestData(t, server.URL, 50)

	// Test with multiple concurrent requests
	concurrency := 10           // Número de goroutines concurrentes
	requestsPerGoroutine := 250 // Requests por goroutine
	totalRequests := concurrency * requestsPerGoroutine

	start := time.Now()
	results := make(chan error, totalRequests)

	// Crear goroutines que ejecuten múltiples requests cada una
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			for j := 0; j < requestsPerGoroutine; j++ {
				// Test different user IDs
				userID := 1001 + ((workerID*requestsPerGoroutine + j) % 10)
				url := fmt.Sprintf("%s%s/users/%d/balance", server.URL, config.GetPathAPI(), userID)

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Get(url)
				if err != nil {
					results <- err
					continue
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
					results <- fmt.Errorf("expected status 200 or 400, got %d", resp.StatusCode)
					continue
				}

				results <- nil
			}
		}(i)
	}

	// Collect results
	successCount := 0
	errorCount := 0
	for i := 0; i < totalRequests; i++ {
		select {
		case err := <-results:
			if err != nil {
				errorCount++
				t.Logf("Request failed: %v", err)
			} else {
				successCount++
			}
		case <-time.After(30 * time.Second):
			t.Error("Test timed out")
			return
		}
	}

	duration := time.Since(start)
	t.Logf("Balance load test completed in %v", duration)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Failed requests: %d", errorCount)
	t.Logf("Requests per second: %.2f", float64(successCount)/duration.Seconds())

	if errorCount > 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
}

// TestLoadBalanceWithDateRange tests the balance endpoint with date range under load
func TestLoadBalanceWithDateRange(t *testing.T) {
	// Setup test server using centralized function
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First, migrate some data using centralized function
	test_utils.MigrateTestData(t, server.URL, 50)

	// Test with multiple concurrent requests
	concurrency := 15           // Número de goroutines concurrentes
	requestsPerGoroutine := 250 // Requests por goroutine
	totalRequests := concurrency * requestsPerGoroutine

	start := time.Now()
	results := make(chan error, totalRequests)

	// Crear goroutines que ejecuten múltiples requests cada una
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			for j := 0; j < requestsPerGoroutine; j++ {
				// Test different user IDs and date ranges
				userID := 1001 + ((workerID*requestsPerGoroutine + j) % 10)
				fromDate := "2024-01-15T00:00:00Z"
				toDate := "2024-01-20T23:59:59Z"
				url := fmt.Sprintf("%s%s/users/%d/balance?from=%s&to=%s", server.URL, config.GetPathAPI(), userID, fromDate, toDate)

				resp, err := http.Get(url)
				if err != nil {
					results <- err
					continue
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
					results <- fmt.Errorf("expected status 200 or 400, got %d", resp.StatusCode)
					continue
				}

				results <- nil
			}
		}(i)
	}

	// Collect results
	successCount := 0
	errorCount := 0
	for i := 0; i < totalRequests; i++ {
		select {
		case err := <-results:
			if err != nil {
				errorCount++
				t.Logf("Request failed: %v", err)
			} else {
				successCount++
			}
		case <-time.After(30 * time.Second):
			t.Error("Test timed out")
			return
		}
	}

	duration := time.Since(start)
	t.Logf("Balance with date range load test completed in %v", duration)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Failed requests: %d", errorCount)
	t.Logf("Requests per second: %.2f", float64(successCount)/duration.Seconds())

	if errorCount > 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
}
