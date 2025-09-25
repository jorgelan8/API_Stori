package test_utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"
)

// GenerateTestCSV generates CSV data for testing
// This is a centralized function to avoid code duplication across tests
func GenerateTestCSV(recordCount int) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write([]string{"id", "user_id", "amount", "datetime"})

	// Write records
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	for i := 1; i <= recordCount; i++ {
		userID := 1001 + (i % 100)        // 100 different users
		amount := float64(i%1000) - 500.0 // Random amounts between -500 and 499
		datetime := baseTime.Add(time.Duration(i) * time.Hour)

		record := []string{
			strconv.Itoa(i),
			strconv.Itoa(userID),
			fmt.Sprintf("%.2f", amount),
			datetime.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.String()
}

// GenerateTestCSVWithUsers generates CSV data with specific user distribution
func GenerateTestCSVWithUsers(recordCount, userCount int) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write([]string{"id", "user_id", "amount", "datetime"})

	// Write records
	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	for i := 1; i <= recordCount; i++ {
		userID := 1001 + (i % userCount)  // Distribute across specified number of users
		amount := float64(i%1000) - 500.0 // Random amounts between -500 and 499
		datetime := baseTime.Add(time.Duration(i) * time.Hour)

		record := []string{
			strconv.Itoa(i),
			strconv.Itoa(userID),
			fmt.Sprintf("%.2f", amount),
			datetime.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.String()
}
