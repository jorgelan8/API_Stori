package services

import (
	"api-stori/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// MigrationService maneja la migración de datos desde archivos CSV
type MigrationService struct {
	database *MockDatabase
}

// NewMigrationService crea una nueva instancia de MigrationService
func NewMigrationService(database *MockDatabase) *MigrationService {
	return &MigrationService{
		database: database,
	}
}

// MigrationResult representa el resultado de una migración
type MigrationResult struct {
	TotalRecords   int                      `json:"total_records"`
	SuccessRecords int                      `json:"success_records"`
	ErrorRecords   int                      `json:"error_records"`
	Transactions   []models.UserTransaction `json:"transactions,omitempty"`
}

// ProcessCSV procesa un archivo CSV y migra las transacciones a la base de datos
func (ms *MigrationService) ProcessCSV(reader io.Reader) (*MigrationResult, error) {
	csvReader := csv.NewReader(reader)

	// Leer todas las líneas del CSV
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Verificar que tenga el header esperado
	header := records[0]
	expectedHeader := []string{"id", "user_id", "amount", "datetime"}
	if !ms.validateHeader(header, expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header. Expected: %v, Got: %v", expectedHeader, header)
	}

	result := &MigrationResult{
		TotalRecords: len(records) - 1, // Excluir header
		Transactions: []models.UserTransaction{},
	}

	// Procesar cada línea de datos (saltar header)
	for i, record := range records[1:] {
		transaction, err := ms.parseTransaction(record, i+2) // +2 porque empezamos desde línea 2
		if err != nil {
			result.ErrorRecords++
			fmt.Printf("Error parsing record at line %d: %v\n", i+2, err)
			continue
		}

		// Guardar en la base de datos mock
		savedTransaction, err := ms.database.SaveTransaction(transaction)
		if err != nil {
			result.ErrorRecords++
			fmt.Printf("Error saving transaction at line %d: %v\n", i+2, err)
			continue
		}

		result.SuccessRecords++
		result.Transactions = append(result.Transactions, savedTransaction)
	}

	return result, nil
}

// validateHeader verifica que el header del CSV sea correcto
func (ms *MigrationService) validateHeader(header, expected []string) bool {
	if len(header) != len(expected) {
		return false
	}

	for i, col := range header {
		if col != expected[i] {
			return false
		}
	}

	return true
}

// parseTransaction convierte una línea del CSV en una transacción
func (ms *MigrationService) parseTransaction(record []string, lineNumber int) (models.UserTransaction, error) {
	if len(record) != 4 {
		return models.UserTransaction{}, fmt.Errorf("invalid number of columns at line %d", lineNumber)
	}

	// Parsear ID
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid ID at line %d: %v", lineNumber, err)
	}

	// Parsear UserID
	userID, err := strconv.Atoi(record[1])
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid user_id at line %d: %v", lineNumber, err)
	}

	// Parsear Amount
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid amount at line %d: %v", lineNumber, err)
	}

	// Parsear DateTime
	datetime, err := time.Parse("2006-01-02 15:04:05", record[3])
	if err != nil {
		// Intentar con otro formato común
		datetime, err = time.Parse("2006-01-02T15:04:05", record[3])
		if err != nil {
			// Intentar con formato de fecha solamente
			datetime, err = time.Parse("2006-01-02", record[3])
			if err != nil {
				return models.UserTransaction{}, fmt.Errorf("invalid datetime at line %d: %v", lineNumber, err)
			}
		}
	}

	return models.UserTransaction{
		ID:       id,
		UserID:   userID,
		Amount:   amount,
		DateTime: datetime,
	}, nil
}

// GetMigrationStats retorna estadísticas de la migración
func (ms *MigrationService) GetMigrationStats() map[string]interface{} {
	allTransactions := ms.database.GetAllTransactions()

	stats := map[string]interface{}{
		"total_transactions": len(allTransactions),
	}

	// Calcular estadísticas por usuario
	userStats := make(map[int]int)
	totalAmount := 0.0

	for _, transaction := range allTransactions {
		userStats[transaction.UserID]++
		totalAmount += transaction.Amount
	}

	stats["users_count"] = len(userStats)
	stats["total_amount"] = totalAmount
	stats["user_transaction_counts"] = userStats

	return stats
}
