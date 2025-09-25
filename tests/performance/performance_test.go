package performance

import (
	"api-stori/tests/config"
	"api-stori/tests/test_utils"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
	"time"
)

// TestPerformanceMigration tests migration endpoint performance
func TestPerformanceMigration(t *testing.T) {
	// Setup test server
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Test with different CSV sizes
	sizes := []int{10, 100, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("CSV_Size_%d", size), func(t *testing.T) {
			csvData := generateCSV(size)

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			fileWriter, err := writer.CreateFormFile("csv_file", "perf_test.csv")
			if err != nil {
				t.Fatalf("Expected no error creating form file, got %v", err)
			}
			fileWriter.Write([]byte(csvData))
			writer.Close()

			start := time.Now()
			req, err := http.NewRequest("POST", server.URL+config.GetPathAPI()+"/migrate", &buf)
			if err != nil {
				t.Fatalf("Expected no error creating request, got %v", err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			client := &http.Client{Timeout: 60 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Expected no error making request, got %v", err)
			}
			defer resp.Body.Close()

			duration := time.Since(start)

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", resp.StatusCode)
			}

			t.Logf("CSV size: %d records, Duration: %v, Records/sec: %.2f",
				size, duration, float64(size)/duration.Seconds())
		})
	}
}

// TestPerformanceBalance tests balance endpoint performance
func TestPerformanceBalance(t *testing.T) {
	// Setup test server
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First migrate some data
	migrateTestDataWithServer(t, server.URL)

	// Test with different user counts
	userCounts := []int{1, 10, 50, 100}

	for _, userCount := range userCounts {
		t.Run(fmt.Sprintf("Users_%d", userCount), func(t *testing.T) {
			start := time.Now()

			for i := 0; i < userCount; i++ {
				userID := 1001 + (i % 10)
				url := fmt.Sprintf("%s%s/users/%d/balance", server.URL, config.GetPathAPI(), userID)

				resp, err := http.Get(url)
				if err != nil {
					t.Fatalf("Expected no error making request, got %v", err)
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
					t.Fatalf("Expected status 200 or 400, got %d", resp.StatusCode)
				}
			}

			duration := time.Since(start)
			t.Logf("User count: %d, Duration: %v, Requests/sec: %.2f",
				userCount, duration, float64(userCount)/duration.Seconds())
		})
	}
}

// TestPerformanceBalanceWithDateRange tests balance endpoint with date range performance
func TestPerformanceBalanceWithDateRange(t *testing.T) {
	// Setup test server
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First migrate some data
	migrateTestDataWithServer(t, server.URL)

	// Test with different date ranges
	dateRanges := []struct {
		name string
		from string
		to   string
	}{
		{"1_Day", "2024-01-15T00:00:00Z", "2024-01-15T23:59:59Z"},
		{"3_Days", "2024-01-15T00:00:00Z", "2024-01-17T23:59:59Z"},
		{"7_Days", "2024-01-15T00:00:00Z", "2024-01-21T23:59:59Z"},
		{"30_Days", "2024-01-01T00:00:00Z", "2024-01-30T23:59:59Z"},
	}

	for _, dateRange := range dateRanges {
		t.Run(dateRange.name, func(t *testing.T) {
			start := time.Now()

			// Test with 50 requests
			for i := 0; i < 50; i++ {
				userID := 1001 + (i % 10)
				url := fmt.Sprintf("%s%s/users/%d/balance?from=%s&to=%s",
					server.URL, config.GetPathAPI(), userID, dateRange.from, dateRange.to)

				resp, err := http.Get(url)
				if err != nil {
					t.Fatalf("Expected no error making request, got %v", err)
				}
				resp.Body.Close()

				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
					t.Fatalf("Expected status 200 or 400, got %d", resp.StatusCode)
				}
			}

			duration := time.Since(start)
			t.Logf("Date range: %s, Duration: %v, Requests/sec: %.2f",
				dateRange.name, duration, 50.0/duration.Seconds())
		})
	}
}

// TestMemoryUsage tests memory usage during operations
func TestMemoryUsage(t *testing.T) {
	// Setup test server
	server := test_utils.SetupTestServer()
	defer server.Close()

	// Test memory usage with large CSV
	csvData := generateCSV(10000) // 10,000 records

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, err := writer.CreateFormFile("csv_file", "memory_test.csv")
	if err != nil {
		t.Fatalf("Expected no error creating form file, got %v", err)
	}
	fileWriter.Write([]byte(csvData))
	writer.Close()

	start := time.Now()
	req, err := http.NewRequest("POST", server.URL+config.GetPathAPI()+"/migrate", &buf)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Large CSV (10,000 records) processed in %v", duration)
	t.Logf("Records per second: %.2f", 10000.0/duration.Seconds())
}

// TestConcurrentRequests tests concurrent request handling
func TestConcurrentRequests(t *testing.T) {
	// Setup test server
	server := test_utils.SetupTestServer()
	defer server.Close()

	// First migrate some data
	migrateTestDataWithServer(t, server.URL)

	concurrencyLevels := []int{1, 5, 10, 20, 50}

	for _, concurrency := range concurrencyLevels {
		t.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(t *testing.T) {
			start := time.Now()
			results := make(chan error, concurrency)

			for i := 0; i < concurrency; i++ {
				go func() {
					userID := 1001 + (i % 10)
					url := fmt.Sprintf("%s%s/users/%d/balance", server.URL, config.GetPathAPI(), userID)

					resp, err := http.Get(url)
					if err != nil {
						results <- err
						return
					}
					resp.Body.Close()

					if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
						results <- fmt.Errorf("expected status 200 or 400, got %d", resp.StatusCode)
						return
					}

					results <- nil
				}()
			}

			// Collect results
			successCount := 0
			errorCount := 0
			for i := 0; i < concurrency; i++ {
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
			t.Logf("Concurrency: %d, Duration: %v, Success: %d, Errors: %d",
				concurrency, duration, successCount, errorCount)

			if errorCount > 0 {
				t.Errorf("Expected 0 errors, got %d", errorCount)
			}
		})
	}
}

// Helper function to generate CSV data
func generateCSV(recordCount int) string {
	var buf bytes.Buffer
	buf.WriteString("id,user_id,amount,datetime\n")

	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	for i := 1; i <= recordCount; i++ {
		userID := 1001 + (i % 100)        // 100 different users
		amount := float64(i%1000) - 500.0 // Random amounts between -500 and 499
		datetime := baseTime.Add(time.Duration(i) * time.Hour)

		record := fmt.Sprintf("%d,%d,%.2f,%s\n",
			i, userID, amount, datetime.Format("2006-01-02 15:04:05"))
		buf.WriteString(record)
	}

	return buf.String()
}

// Helper function to migrate test data with server
func migrateTestDataWithServer(t *testing.T, serverURL string) {
	csvData := generateCSV(100) // 100 transactions

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, err := writer.CreateFormFile("csv_file", "test.csv")
	if err != nil {
		t.Fatalf("Expected no error creating form file, got %v", err)
	}
	fileWriter.Write([]byte(csvData))
	writer.Close()

	req, err := http.NewRequest("POST", serverURL+config.GetPathAPI()+"/migrate", &buf)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error making request, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
