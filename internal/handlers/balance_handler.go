package handlers

import (
	"api-stori/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BalanceHandler maneja las requests del endpoint de balance
type BalanceHandler struct {
	usersService *services.UsersService
}

// NewBalanceHandler crea una nueva instancia de BalanceHandler
func NewBalanceHandler(usersService *services.UsersService) *BalanceHandler {
	return &BalanceHandler{
		usersService: usersService,
	}
}

// BalanceResponse representa la respuesta del endpoint de balance
type BalanceResponse struct {
	Balance      float64 `json:"balance"`
	TotalDebits  int     `json:"total_debits"`
	TotalCredits int     `json:"total_credits"`
}

// GetUserBalance maneja el endpoint GET /users/{user_id}/balance
func (h *BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método sea GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer user_id de la URL
	userIDStr := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	// Obtener parámetros de fecha
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var fromDate, toDate *time.Time

	// Parsear fecha "from" si se proporciona
	if fromStr != "" {
		parsedFrom, err := time.Parse("2006-01-02T15:04:05Z", fromStr)
		if err != nil {
			http.Error(w, "Invalid 'from' date format. Expected: YYYY-MM-DDTHH:MM:SSZ", http.StatusBadRequest)
			return
		}
		fromDate = &parsedFrom
	}

	// Parsear fecha "to" si se proporciona
	if toStr != "" {
		parsedTo, err := time.Parse("2006-01-02T15:04:05Z", toStr)
		if err != nil {
			http.Error(w, "Invalid 'to' date format. Expected: YYYY-MM-DDTHH:MM:SSZ", http.StatusBadRequest)
			return
		}
		toDate = &parsedTo
	}

	// Validar que from sea anterior a to si ambos se proporcionan
	if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
		http.Error(w, "Invalid date range: 'from' date must be before 'to' date", http.StatusBadRequest)
		return
	}

	// Obtener balance del usuario usando el servicio
	balanceInfo, err := h.usersService.GetUserBalance(userID, fromDate, toDate)
	if err != nil {
		if err == services.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Preparar respuesta
	response := BalanceResponse{
		Balance:      balanceInfo.Balance,
		TotalDebits:  balanceInfo.TotalDebits,
		TotalCredits: balanceInfo.TotalCredits,
	}

	// Escribir respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
