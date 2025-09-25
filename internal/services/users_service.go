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

// GetUserBalance obtiene el balance de un usuario con filtros opcionales de fecha
func (us *UsersService) GetUserBalance(userID int, fromDate, toDate *time.Time) (*models.BalanceInfo, error) {
	// Obtener transacciones del usuario (con filtro de fechas si se especifica)
	userTransactions := us.database.GetTransactionsByUserIDWithDateRange(userID, fromDate, toDate)

	// Si no hay transacciones, el usuario no existe
	if len(userTransactions) == 0 {
		return nil, ErrUserNotFound
	}

	// Calcular balance, débitos y créditos
	balance, totalDebits, totalCredits := us.calculateBalance(userTransactions)

	return &models.BalanceInfo{
		Balance:      models.Float64(balance),
		TotalDebits:  models.Float64(totalDebits),
		TotalCredits: models.Float64(totalCredits),
	}, nil
}

// calculateBalance calcula el balance, total de débitos y créditos
func (us *UsersService) calculateBalance(transactions []models.UserTransaction) (float64, float64, float64) {
	var balance float64
	var totalDebits, totalCredits float64

	for _, transaction := range transactions {
		balance += transaction.Amount

		if transaction.Amount < 0 {
			totalDebits += transaction.Amount
		} else if transaction.Amount > 0 {
			totalCredits += transaction.Amount
		}
	}

	return balance, totalDebits, totalCredits
}
