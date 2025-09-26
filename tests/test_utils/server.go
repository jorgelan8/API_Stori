package test_utils

import (
	"api-stori/internal/routes"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// SetupTestServer creates a test server for all tests
func SetupTestServer() *httptest.Server {
	router := mux.NewRouter()
	routes.SetupRoutesConfigDetail(router, false)
	return httptest.NewServer(router)
}
