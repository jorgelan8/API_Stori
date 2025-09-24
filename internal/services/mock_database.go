package services

import (
	"api-stori/internal/models"
	"sync"
)

// MockDatabase simula una base de datos en memoria
type MockDatabase struct {
	transactions map[int]models.UserTransaction
	nextID       int
	mutex        sync.RWMutex
}

// NewMockDatabase crea una nueva instancia de MockDatabase
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		transactions: make(map[int]models.UserTransaction),
		nextID:       1,
	}
}

// SaveTransaction guarda una transacción en el mock de base de datos
func (db *MockDatabase) SaveTransaction(transaction models.UserTransaction) (models.UserTransaction, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Si no tiene ID, asignar uno nuevo
	if transaction.ID == 0 {
		transaction.ID = db.nextID
		db.nextID++
	}

	// Guardar la transacción
	db.transactions[transaction.ID] = transaction

	return transaction, nil
}

// SaveTransactions guarda múltiples transacciones en el mock de base de datos
func (db *MockDatabase) SaveTransactions(transactions []models.UserTransaction) ([]models.UserTransaction, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var savedTransactions []models.UserTransaction

	for _, transaction := range transactions {
		// Si no tiene ID, asignar uno nuevo
		if transaction.ID == 0 {
			transaction.ID = db.nextID
			db.nextID++
		}

		// Guardar la transacción
		db.transactions[transaction.ID] = transaction
		savedTransactions = append(savedTransactions, transaction)
	}

	return savedTransactions, nil
}

// GetTransaction obtiene una transacción por ID
func (db *MockDatabase) GetTransaction(id int) (models.UserTransaction, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	transaction, exists := db.transactions[id]
	return transaction, exists
}

// GetTransactionsByUserID obtiene todas las transacciones de un usuario
func (db *MockDatabase) GetTransactionsByUserID(userID int) []models.UserTransaction {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	var userTransactions []models.UserTransaction
	for _, transaction := range db.transactions {
		if transaction.UserID == userID {
			userTransactions = append(userTransactions, transaction)
		}
	}

	return userTransactions
}

// GetAllTransactions obtiene todas las transacciones
func (db *MockDatabase) GetAllTransactions() []models.UserTransaction {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	var allTransactions []models.UserTransaction
	for _, transaction := range db.transactions {
		allTransactions = append(allTransactions, transaction)
	}

	return allTransactions
}

// GetTransactionCount retorna el número total de transacciones
func (db *MockDatabase) GetTransactionCount() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.transactions)
}

// ClearTransactions limpia todas las transacciones (útil para testing)
func (db *MockDatabase) ClearTransactions() {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.transactions = make(map[int]models.UserTransaction)
	db.nextID = 1
}
