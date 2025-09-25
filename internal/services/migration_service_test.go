package services

import (
	"api-stori/internal/models"
	"strings"
	"testing"
	"time"
)

func TestMigrationService_ProcessCSV(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Test CSV content
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,-75.25,2024-01-15 14:45:00
3,1002,200.00,2024-01-16 09:15:00`

	reader := strings.NewReader(csvContent)

	result, err := service.ProcessCSV(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check result
	if result.TotalRecords != 3 {
		t.Errorf("Expected 3 total records, got %d", result.TotalRecords)
	}

	if result.SuccessRecords != 3 {
		t.Errorf("Expected 3 success records, got %d", result.SuccessRecords)
	}

	if result.ErrorRecords != 0 {
		t.Errorf("Expected 0 error records, got %d", result.ErrorRecords)
	}

	if len(result.Transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(result.Transactions))
	}

	// Verify first transaction
	firstTx := result.Transactions[0]
	if firstTx.ID != 1 {
		t.Errorf("Expected ID 1, got %d", firstTx.ID)
	}
	if firstTx.UserID != 1001 {
		t.Errorf("Expected UserID 1001, got %d", firstTx.UserID)
	}
	if firstTx.Amount != 150.50 {
		t.Errorf("Expected Amount 150.50, got %.2f", firstTx.Amount)
	}
}

func TestMigrationService_ProcessCSVWithErrors(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Test CSV with invalid data
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,invalid_user,invalid_amount,invalid_date
3,1002,200.00,2024-01-16 09:15:00`

	reader := strings.NewReader(csvContent)

	result, err := service.ProcessCSV(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 1 error record
	if result.ErrorRecords != 1 {
		t.Errorf("Expected 1 error record, got %d", result.ErrorRecords)
	}

	if result.SuccessRecords != 2 {
		t.Errorf("Expected 2 success records, got %d", result.SuccessRecords)
	}
}

func TestMigrationService_ProcessCSVEmptyFile(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Test empty CSV
	csvContent := ``
	reader := strings.NewReader(csvContent)

	_, err := service.ProcessCSV(reader)
	if err == nil {
		t.Error("Expected error for empty CSV file")
	}
}

func TestMigrationService_ProcessCSVInvalidHeader(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Test CSV with wrong header
	csvContent := `wrong,header,format
1,1001,150.50`
	reader := strings.NewReader(csvContent)

	_, err := service.ProcessCSV(reader)
	if err == nil {
		t.Error("Expected error for invalid header")
	}
}

func TestMigrationService_ProcessCSVDifferentDateFormats(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Test CSV with different date formats
	csvContent := `id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,200.00,2024-01-16T09:15:00
3,1001,300.00,2024-01-17`

	reader := strings.NewReader(csvContent)

	result, err := service.ProcessCSV(reader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.SuccessRecords != 3 {
		t.Errorf("Expected 3 success records, got %d", result.SuccessRecords)
	}
}

func TestMigrationService_GetMigrationStats(t *testing.T) {
	db := NewMockDatabase()
	service := NewMigrationService(db)

	// Add some transactions
	transactions := []models.UserTransaction{
		{ID: 1, UserID: 1001, Amount: 150.50, DateTime: time.Now()},
		{ID: 2, UserID: 1001, Amount: -75.25, DateTime: time.Now()},
		{ID: 3, UserID: 1002, Amount: 200.00, DateTime: time.Now()},
	}

	for _, tx := range transactions {
		db.SaveTransaction(tx)
	}

	stats := service.GetMigrationStats()

	// Check total transactions
	if stats["total_transactions"] != 3 {
		t.Errorf("Expected 3 total transactions, got %v", stats["total_transactions"])
	}

	// Check users count
	if stats["users_count"] != 2 {
		t.Errorf("Expected 2 users, got %v", stats["users_count"])
	}

	// Check total amount
	expectedTotal := 150.50 - 75.25 + 200.00
	if stats["total_amount"] != expectedTotal {
		t.Errorf("Expected total amount %.2f, got %v", expectedTotal, stats["total_amount"])
	}

	// Check user transaction counts
	userCounts := stats["user_transaction_counts"].(map[int]int)
	if userCounts[1001] != 2 {
		t.Errorf("Expected user 1001 to have 2 transactions, got %d", userCounts[1001])
	}
	if userCounts[1002] != 1 {
		t.Errorf("Expected user 1002 to have 1 transaction, got %d", userCounts[1002])
	}
}

// Note: parseTransaction and validateHeader are private methods
// These tests are covered indirectly through ProcessCSV tests
