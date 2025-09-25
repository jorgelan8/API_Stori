package main

import (
	"api-stori/internal/config"
	"api-stori/internal/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuración
	appConfig := config.LoadConfig()

	// Crear router
	router := mux.NewRouter()

	// Configurar todas las rutas
	routes.SetupRoutes(router)

	// Iniciar servidor
	fmt.Printf("🚀 Server starting on port %s\n", appConfig.App.Port)
	fmt.Printf("📊 API Stori endpoints:\n")
	fmt.Printf("   POST /api/v1/migrate - Upload CSV file\n")
	fmt.Printf("   GET  /api/v1/users/{user_id}/balance - Get user balance\n")
	fmt.Printf("   GET  /api/v1/health - Health check\n")
	fmt.Printf("   GET  / - API information\n")

	// Mostrar configuración de email si está configurada
	if appConfig.Email.Username != "" {
		fmt.Printf("📧 Email reports enabled: %s\n", appConfig.Email.FromEmail)
	} else {
		fmt.Printf("📧 Email reports: Mock mode (no SMTP configured)\n")
	}

	if err := http.ListenAndServe(":"+appConfig.App.Port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
