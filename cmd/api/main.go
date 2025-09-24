package main

import (
	"api-stori/internal/routes"
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

	// Configurar todas las rutas
	routes.SetupRoutes(router)

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
