package main

import (
	"api-stori/internal/routes"
	"api-stori/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Obtener puerto de variable de entorno o usar default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Crear router
	router := mux.NewRouter()

	// Crear instancias de servicios
	mockDB := services.NewMockDatabase()
	migrationService := services.NewMigrationService(mockDB)
	usersService := services.NewUsersService(mockDB)

	// Configurar rutas de migración
	routes.SetupMigrationRoutes(router, migrationService)

	// Configurar rutas de balance
	routes.SetupBalanceRoutes(router, usersService)

	// Endpoint raíz
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "API Stori - Migration & Balance Service",
			"version": "1.0.0",
			"endpoints": {
				"migrate": "POST /api/v1/migrate",
				"balance": "GET /api/v1/users/{user_id}/balance",
				"health": "GET /api/v1/health"
			}
		}`))
	}).Methods("GET")

	// 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Endpoint not found", "status": 404}`))
	})

	// 405 handler
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "status": 405}`))
	})

	// Iniciar servidor
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("📊 API Stori endpoints:\n")
	fmt.Printf("   POST /api/v1/migrate - Upload CSV file\n")
	fmt.Printf("   GET  /api/v1/users/{user_id}/balance - Get user balance\n")
	fmt.Printf("   GET  /api/v1/health - Health check\n")
	fmt.Printf("   GET  / - API information\n")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
