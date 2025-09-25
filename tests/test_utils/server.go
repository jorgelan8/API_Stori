package test_utils

import (
	"api-stori/internal/routes"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// SetupTestServer creates a test server for all tests
// This is a centralized function to avoid code duplication
func SetupTestServer() *httptest.Server {
	router := mux.NewRouter()
	routes.SetupRoutes(router)
	return httptest.NewServer(router)
}
