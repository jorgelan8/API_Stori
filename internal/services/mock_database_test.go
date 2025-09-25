package services

import (
	"api-stori/internal/models"
	"testing"
	"time"
)

func TestMockDatabase_SaveTransaction(t *testing.T) {
	db := NewMockDatabase()

	// Test saving a transaction
	transaction := models.UserTransaction{
		ID:       1,
		UserID:   1001,
		Amount:   150.50,
		DateTime: time.Now(),
	}

	saved, err := db.SaveTransaction(transaction)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if saved.ID != transaction.ID {
		t.Errorf("Expected ID %d, got %d", transaction.ID, saved.ID)
	}

	// Test auto-increment ID
	transaction2 := models.UserTransaction{
		UserID:   1002,
		Amount:   -75.25,
		DateTime: time.Now(),
	}

	saved2, err := db.SaveTransaction(transaction2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if saved2.ID == 0 {
		t.Error("Expected auto-incremented ID, got 0")
	}
}

func TestMockDatabase_GetTransaction(t *testing.T) {
	db := NewMockDatabase()

	// Save a transaction
	transaction := models.UserTransaction{
		ID:       1,
		UserID:   1001,
		Amount:   150.50,
		DateTime: time.Now(),
	}
	db.SaveTransaction(transaction)

	// Test getting existing transaction
	retrieved, exists := db.GetTransaction(1)
	if !exists {
		t.Error("Expected transaction to exist")
	}
	if retrieved.ID != transaction.ID {
		t.Errorf("Expected ID %d, got %d", transaction.ID, retrieved.ID)
	}

	// Test getting non-existing transaction
	_, exists = db.GetTransaction(999)
	if exists {
		t.Error("Expected transaction to not exist")
	}
}

func TestMockDatabase_GetTransactionsByUserID(t *testing.T) {
	db := NewMockDatabase()

	// Save multiple transactions for different users
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: time.Now()},
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: time.Now()},
		{ID: 3, UserID: 1002, Amount: 200.00, DateTime: time.Now()},
		{ID: 4, UserID: 1001, Amount: 50.75, DateTime: time.Now()},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	// Test getting transactions for user 1001
	userTransactions := db.GetTransactionsByUserID(1001)
	if len(userTransactions) != 3 {
		t.Errorf("Expected 3 transactions for user 1001, got %d", len(userTransactions))
	}

	// Test getting transactions for user 1002
	userTransactions = db.GetTransactionsByUserID(1002)
	if len(userTransactions) != 1 {
		t.Errorf("Expected 1 transaction for user 1002, got %d", len(userTransactions))
	}

	// Test getting transactions for non-existing user
	userTransactions = db.GetTransactionsByUserID(9999)
	if len(userTransactions) != 0 {
		t.Errorf("Expected 0 transactions for non-existing user, got %d", len(userTransactions))
	}
}

func TestMockDatabase_GetTransactionsByUserIDWithDateRange(t *testing.T) {
	db := NewMockDatabase()

	// Create transactions with specific dates
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: baseTime},
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: baseTime.Add(24 * time.Hour)},
		{ID: 3, UserID: 1001, Amount: 200.00, DateTime: baseTime.Add(48 * time.Hour)},
		{ID: 4, UserID: 1001, Amount: 50.75, DateTime: baseTime.Add(72 * time.Hour)},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	// Test filtering by date range
	fromDate := baseTime.Add(12 * time.Hour) // 2024-01-15 12:00
	toDate := baseTime.Add(60 * time.Hour)   // 2024-01-17 12:00

	filtered := db.GetTransactionsByUserIDWithDateRange(1001, &fromDate, &toDate)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 transactions in date range, got %d", len(filtered))
	}

	// Test filtering with only from date
	filtered = db.GetTransactionsByUserIDWithDateRange(1001, &fromDate, nil)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 transactions from date, got %d", len(filtered))
	}

	// Test filtering with only to date
	filtered = db.GetTransactionsByUserIDWithDateRange(1001, nil, &toDate)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 transactions to date, got %d", len(filtered))
	}

	// Test filtering with no dates
	filtered = db.GetTransactionsByUserIDWithDateRange(1001, nil, nil)
	if len(filtered) != 4 {
		t.Errorf("Expected 4 transactions with no date filter, got %d", len(filtered))
	}
}

func TestMockDatabase_GetAllTransactions(t *testing.T) {
	db := NewMockDatabase()

	// Save multiple transactions
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: time.Now()},
		{ID: 2, UserID: 1002, Amount: -75.25, DateTime: time.Now()},
		{ID: 3, UserID: 1001, Amount: 200.00, DateTime: time.Now()},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	allTransactions := db.GetAllTransactions()
	if len(allTransactions) != 3 {
		t.Errorf("Expected 3 total transactions, got %d", len(allTransactions))
	}
}

func TestMockDatabase_GetTransactionCount(t *testing.T) {
	db := NewMockDatabase()

	// Initially should be 0
	count := db.GetTransactionCount()
	if count != 0 {
		t.Errorf("Expected 0 transactions initially, got %d", count)
	}

	// Save some transactions
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: time.Now()},
		{ID: 2, UserID: 1002, Amount: -75.25, DateTime: time.Now()},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	count = db.GetTransactionCount()
	if count != 2 {
		t.Errorf("Expected 2 transactions, got %d", count)
	}
}

func TestMockDatabase_ClearTransactions(t *testing.T) {
	db := NewMockDatabase()

	// Save some transactions
	transaction := models.UserTransaction{
		ID:       1,
		UserID:   1001,
		Amount:   150.50,
		DateTime: time.Now(),
	}
	db.SaveTransaction(transaction)

	// Verify it exists
	if db.GetTransactionCount() != 1 {
		t.Error("Expected 1 transaction before clear")
	}

	// Clear transactions
	db.ClearTransactions()

	// Verify it's cleared
	if db.GetTransactionCount() != 0 {
		t.Error("Expected 0 transactions after clear")
	}

	// Verify next ID is reset
	transaction2 := models.UserTransaction{
		UserID:   1002,
		Amount:   200.00,
		DateTime: time.Now(),
	}
	saved, _ := db.SaveTransaction(transaction2)
	if saved.ID != 1 {
		t.Errorf("Expected ID to reset to 1, got %d", saved.ID)
	}
}
