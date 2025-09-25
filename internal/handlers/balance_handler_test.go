package handlers

import (
	"api-stori/internal/models"
	"api-stori/internal/services"
	"api-stori/tests/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestBalanceHandler_GetUserBalance(t *testing.T) {
	// Setup
	db := services.NewMockDatabase()
	usersService := services.NewUsersService(db)
	handler := NewBalanceHandler(usersService)

	// Add test data
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: baseTime},
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: baseTime.Add(24 * time.Hour)},
		{ID: 3, UserID: 1002, Amount: 200.00, DateTime: baseTime},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	// Test cases
	tests := []struct {
		name            string
		userID          string
		expectedStatus  int
		expectedBalance float64
	}{
		{
			name:            "Valid user ID",
			userID:          "1001",
			expectedStatus:  http.StatusOK,
			expectedBalance: 150.50 - 75.25,
		},
		{
			name:           "Non-existent user ID",
			userID:         "9999",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID format",
			userID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty user ID",
			userID:         "   ", // Use spaces instead of empty string
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with URL parameters
			req, err := http.NewRequest("GET", config.GetPathAPI()+"/users/"+tt.userID+"/balance", nil)
			if err != nil {
				t.Fatalf("Expected no error creating request, got %v", err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create router and add route with parameter
			router := mux.NewRouter()
			router.HandleFunc(config.GetPathAPI()+"/users/{user_id}/balance", handler.GetUserBalance).Methods("GET")

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// If successful, check response body
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Expected no error decoding response, got %v", err)
				}

				balance, ok := response["balance"].(float64)
				if !ok {
					t.Error("Expected balance field in response")
				} else if balance != tt.expectedBalance {
					t.Errorf("Expected balance %.2f, got %.2f", tt.expectedBalance, balance)
				}
			}
		})
	}
}

func TestBalanceHandler_GetUserBalanceWithDateRange(t *testing.T) {
	// Setup
	db := services.NewMockDatabase()
	usersService := services.NewUsersService(db)
	handler := NewBalanceHandler(usersService)

	// Add test data with specific dates
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: baseTime},                     // 2024-01-15
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: baseTime.Add(24 * time.Hour)}, // 2024-01-16
		{ID: 3, UserID: 1001, Amount: 200.00, DateTime: baseTime.Add(48 * time.Hour)}, // 2024-01-17
		{ID: 4, UserID: 1001, Amount: 50.75, DateTime: baseTime.Add(72 * time.Hour)},  // 2024-01-18
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	// Test with date range
	req, err := http.NewRequest("GET", config.GetPathAPI()+"/users/1001/balance?from=2024-01-16T00:00:00Z&to=2024-01-17T23:59:59Z", nil)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
		router.HandleFunc(config.GetPathAPI()+"/users/{user_id}/balance", handler.GetUserBalance).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Expected no error decoding response, got %v", err)
	}

	// Should only include transactions from Jan 16-17: -75.25 + 200.00 = 124.75
	expectedBalance := -75.25 + 200.00
	balance, ok := response["balance"].(float64)
	if !ok {
		t.Error("Expected balance field in response")
	} else if balance != expectedBalance {
		t.Errorf("Expected balance %.2f, got %.2f", expectedBalance, balance)
	}
}

func TestBalanceHandler_GetUserBalanceInvalidDateFormat(t *testing.T) {
	// Setup
	db := services.NewMockDatabase()
	usersService := services.NewUsersService(db)
	handler := NewBalanceHandler(usersService)

	// Test with invalid date format
	req, err := http.NewRequest("GET", config.GetPathAPI()+"/users/1001/balance?from=invalid-date", nil)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
		router.HandleFunc(config.GetPathAPI()+"/users/{user_id}/balance", handler.GetUserBalance).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestBalanceHandler_GetUserBalanceInvalidDateRange(t *testing.T) {
	// Setup
	db := services.NewMockDatabase()
	usersService := services.NewUsersService(db)
	handler := NewBalanceHandler(usersService)

	// Test with invalid date range (from > to)
	req, err := http.NewRequest("GET", config.GetPathAPI()+"/users/1001/balance?from=2024-01-20T00:00:00Z&to=2024-01-15T23:59:59Z", nil)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
		router.HandleFunc(config.GetPathAPI()+"/users/{user_id}/balance", handler.GetUserBalance).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestBalanceHandler_GetUserBalanceWithSpecialCharactersInURL(t *testing.T) {
	// Setup
	db := services.NewMockDatabase()
	usersService := services.NewUsersService(db)
	handler := NewBalanceHandler(usersService)

	// Test with special characters in URL (this should fail gracefully)
	req, err := http.NewRequest("GET", config.GetPathAPI()+"/users/1001/balance?from=2024-01-15T00:00:00Z&to=2024-01-20T23:59:59Z&extra=/path/with/slashes", nil)
	if err != nil {
		t.Fatalf("Expected no error creating request, got %v", err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
		router.HandleFunc(config.GetPathAPI()+"/users/{user_id}/balance", handler.GetUserBalance).Methods("GET")
	router.ServeHTTP(rr, req)

	// Should still work because Gorilla Mux properly extracts the user_id parameter
	// regardless of query parameters
	if rr.Code != http.StatusOK && rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 200 or 400, got %d", rr.Code)
	}
}
