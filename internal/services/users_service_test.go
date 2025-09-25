package services

import (
	"api-stori/internal/models"
	"testing"
	"time"
)

func TestUsersService_GetUserBalance(t *testing.T) {
	db := NewMockDatabase()
	service := NewUsersService(db)

	// Setup test data
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: baseTime},
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: baseTime.Add(24 * time.Hour)},
		{ID: 3, UserID: 1001, Amount: 200.00, DateTime: baseTime.Add(48 * time.Hour)},
		{ID: 4, UserID: 1002, Amount: 100.00, DateTime: baseTime}, // Different user
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	// Test getting balance for existing user
	balance, err := service.GetUserBalance(1001, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedBalance := 150.50 - 75.25 + 200.00
	if float64(balance.Balance) != expectedBalance {
		t.Errorf("Expected balance %.2f, got %.2f", expectedBalance, float64(balance.Balance))
	}

	if float64(balance.TotalDebits) != -75.25 {
		t.Errorf("Expected -75.25 debit, got %f", float64(balance.TotalDebits))
	}

	if float64(balance.TotalCredits) != 350.50 {
		t.Errorf("Expected 350.50 credits, got %f", float64(balance.TotalCredits))
	}

	// Test getting balance for non-existing user
	_, err = service.GetUserBalance(9999, nil, nil)
	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUsersService_GetUserBalanceWithDateRange(t *testing.T) {
	db := NewMockDatabase()
	service := NewUsersService(db)

	// Setup test data with specific dates
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

	// Test filtering by date range (Jan 16-17)
	fromDate := baseTime.Add(12 * time.Hour) // 2024-01-15 12:00
	toDate := baseTime.Add(60 * time.Hour)   // 2024-01-17 12:00

	balance, err := service.GetUserBalance(1001, &fromDate, &toDate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should only include transactions from Jan 16-17: -75.25 + 200.00 = 124.75
	expectedBalance := -75.25 + 200.00
	if float64(balance.Balance) != expectedBalance {
		t.Errorf("Expected balance %.2f, got %.2f", expectedBalance, float64(balance.Balance))
	}

	if float64(balance.TotalDebits) != -75.25 {
		t.Errorf("Expected -75.25 debit, got %f", float64(balance.TotalDebits))
	}

	if float64(balance.TotalCredits) != 200.00 {
		t.Errorf("Expected 200.00 credit, got %f", float64(balance.TotalCredits))
	}
}
