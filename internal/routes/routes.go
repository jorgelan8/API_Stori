package routes

import (
	"api-stori/internal/config"
	"api-stori/internal/handlers"
	"api-stori/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	SetupRoutesConfigDetail(router, true)
}

// SetupRoutes configura todas las rutas de la API
func SetupRoutesConfigDetail(router *mux.Router, allowSendEmail bool) {
	// Crear instancias de servicios
	mockDB := services.NewMockDatabase()
	migrationService := services.NewMigrationService(mockDB)
	usersService := services.NewUsersService(mockDB)

	// Cargar configuración desde variables de entorno
	appConfig := config.LoadConfig()

	// Configurar servicio de reportes
	reportService := services.NewReportService(appConfig.ToReportConfig())
	if !allowSendEmail {
		reportService.SetForceMockMode(true)
	}
	migrationService.SetReportService(reportService)

	// Crear handlers
	migrationHandler := handlers.NewMigrationHandler(migrationService)
	balanceHandler := handlers.NewBalanceHandler(usersService)

	// Configurar rutas de la API
	api := router.PathPrefix("/api/v1").Subrouter()

	// Migration Service routes
	api.HandleFunc("/migrate", migrationHandler.MigrateCSV).Methods("POST")

	// Balance Service routes
	api.HandleFunc("/users/{user_id}/balance", balanceHandler.GetUserBalance).Methods("GET")

	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "api-stori"}`))
	}).Methods("GET")

	// Swagger documentation routes
	api.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "api/swagger.yaml")
	}).Methods("GET")

	api.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "api/swagger.json")
	}).Methods("GET")

	// Swagger UI endpoint (redirects to swagger.yaml)
	api.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>API Stori - Swagger Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/api/v1/swagger.yaml',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.standalone
            ]
        });
    </script>
</body>
</html>
		`))
	}).Methods("GET")

	// Root endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": "API Stori - Migration & Balance Service",
			"version": "1.0.0",
			"endpoints": {
				"migrate": "POST /api/v1/migrate",
				"balance": "GET /api/v1/users/{user_id}/balance",
				"health": "GET /api/v1/health",
			},
			"documentation": {
				"swagger_ui": "/api/v1/docs",
				"openapi_yaml": "/api/v1/swagger.yaml",
				"openapi_json": "/api/v1/swagger.json"
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
}
