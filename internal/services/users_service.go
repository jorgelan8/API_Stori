package services

import (
	"api-stori/internal/models"
	"time"
)

// UsersService maneja las operaciones de negocio relacionadas con usuarios
type UsersService struct {
	database *MockDatabase
}

// NewUsersService crea una nueva instancia de UsersService
func NewUsersService(database *MockDatabase) *UsersService {
	return &UsersService{
		database: database,
	}
}

// BalanceInfo representa la información de balance de un usuario
type BalanceInfo struct {
	Balance      float64 `json:"balance"`
	TotalDebits  int     `json:"total_debits"`
	TotalCredits int     `json:"total_credits"`
}

// GetUserBalance obtiene el balance de un usuario con filtros opcionales de fecha
func (us *UsersService) GetUserBalance(userID int, fromDate, toDate *time.Time) (*BalanceInfo, error) {
	// Obtener transacciones del usuario (con filtro de fechas si se especifica)
	userTransactions := us.database.GetTransactionsByUserIDWithDateRange(userID, fromDate, toDate)

	// Si no hay transacciones, el usuario no existe
	if len(userTransactions) == 0 {
		return nil, ErrUserNotFound
	}

	// Calcular balance, débitos y créditos
	balance, totalDebits, totalCredits := us.calculateBalance(userTransactions)

	return &BalanceInfo{
		Balance:      balance,
		TotalDebits:  totalDebits,
		TotalCredits: totalCredits,
	}, nil
}

// calculateBalance calcula el balance, total de débitos y créditos
func (us *UsersService) calculateBalance(transactions []models.UserTransaction) (float64, int, int) {
	var balance float64
	var totalDebits, totalCredits int

	for _, transaction := range transactions {
		balance += transaction.Amount

		if transaction.Amount < 0 {
			totalDebits++
		} else if transaction.Amount > 0 {
			totalCredits++
		}
	}

	return balance, totalDebits, totalCredits
}
