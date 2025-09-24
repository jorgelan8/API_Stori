package routes

import (
	"api-stori/internal/handlers"
	"api-stori/internal/services"

	"github.com/gorilla/mux"
)

// SetupBalanceRoutes configura las rutas del servicio de balance
func SetupBalanceRoutes(router *mux.Router, usersService *services.UsersService) {
	// Crear handler
	balanceHandler := handlers.NewBalanceHandler(usersService)

	// Configurar rutas
	api := router.PathPrefix("/api/v1").Subrouter()

	// Endpoint de balance de usuario
	api.HandleFunc("/users/{user_id}/balance", balanceHandler.GetUserBalance).Methods("GET")
}
