package routes

import (
	"api-stori/internal/handlers"
	"api-stori/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupMigrationRoutes configura las rutas del servicio de migración
func SetupMigrationRoutes(router *mux.Router, migrationService *services.MigrationService) {
	// Crear handler
	migrationHandler := handlers.NewMigrationHandler(migrationService)

	// Configurar rutas
	api := router.PathPrefix("/api/v1").Subrouter()

	// Endpoint principal de migración
	api.HandleFunc("/migrate", migrationHandler.MigrateCSV).Methods("POST")

	// Endpoint de health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "migration-service"}`))
	}).Methods("GET")
}
