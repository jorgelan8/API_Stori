package handlers

import (
	"api-stori/internal/services"
	"net/http"
)

// MigrationHandler maneja las requests del endpoint de migración
type MigrationHandler struct {
	migrationService *services.MigrationService
}

// NewMigrationHandler crea una nueva instancia de MigrationHandler
func NewMigrationHandler(migrationService *services.MigrationService) *MigrationHandler {
	return &MigrationHandler{
		migrationService: migrationService,
	}
}

// MigrateCSV maneja el endpoint POST /migrate
func (h *MigrationHandler) MigrateCSV(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el Content-Type sea multipart/form-data
	contentType := r.Header.Get("Content-Type")
	if contentType == "" || (contentType != "multipart/form-data" && contentType[:19] != "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	// Parsear el formulario multipart
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		http.Error(w, "Error parsing multipart form => "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el archivo CSV
	file, header, err := r.FormFile("csv_file")
	if err != nil {
		http.Error(w, "Error retrieving CSV file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verificar que sea un archivo CSV
	if header.Header.Get("Content-Type") != "text/csv" &&
		header.Filename[len(header.Filename)-4:] != ".csv" {
		http.Error(w, "File must be a CSV file", http.StatusBadRequest)
		return
	}

	// Procesar el archivo CSV
	_, err = h.migrationService.ProcessCSV(file)
	if err != nil {
		http.Error(w, "Error processing CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver solo código HTTP 200 OK sin body
	w.WriteHeader(http.StatusOK)
}
