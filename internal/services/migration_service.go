package services

import (
	"api-stori/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// MigrationService maneja la migración de datos desde archivos CSV
type MigrationService struct {
	database      *MockDatabase
	reportService *ReportService
}

// NewMigrationService crea una nueva instancia de MigrationService
func NewMigrationService(database *MockDatabase) *MigrationService {
	// Crear ReportService por defecto (solo logs, modo mock)
	defaultConfig := &models.ReportConfig{
		Channels: []models.ReportChannel{models.LogChannel},
		Email: models.EmailConfig{
			SMTPHost: "", // Sin SMTP = modo mock por defecto
		},
	}
	defaultReportService := NewReportServiceWithMockMode(defaultConfig)

	return &MigrationService{
		database:      database,
		reportService: defaultReportService,
	}
}

// SetReportService establece el servicio de reportes
func (ms *MigrationService) SetReportService(reportService *ReportService) {
	ms.reportService = reportService
}

// GetReportService devuelve el servicio de reportes
func (ms *MigrationService) GetReportService() *ReportService {
	return ms.reportService
}

// MigrationStats representa las estadísticas de migración (usado tanto para procesamiento como respuesta)
type MigrationStats struct {
	TotalRecords   int      `json:"total_records"`
	SuccessRecords int      `json:"success_records"`
	ErrorRecords   int      `json:"error_records"`
	Errors         []string `json:"errors,omitempty"`

	// Campos internos para cálculos (no se serializan en JSON)
	UsersAffected  map[int]bool
	TotalAmount    float64
	LargestAmount  float64
	SmallestAmount float64
	FirstDate      time.Time
	LastDate       time.Time
}

// NewMigrationStats crea una nueva instancia de estadísticas
func NewMigrationStats() *MigrationStats {
	return &MigrationStats{
		UsersAffected: make(map[int]bool),
		Errors:        []string{},
	}
}

// UpdateSuccess actualiza las estadísticas para una transacción exitosa
func (ms *MigrationStats) UpdateSuccess(transaction models.UserTransaction) {
	ms.SuccessRecords++
	ms.UsersAffected[transaction.UserID] = true
	ms.TotalAmount += transaction.Amount

	// Actualizar montos
	if transaction.Amount > ms.LargestAmount {
		ms.LargestAmount = transaction.Amount
	}
	if transaction.Amount < ms.SmallestAmount || ms.SmallestAmount == 0 {
		ms.SmallestAmount = transaction.Amount
	}

	// Actualizar fechas
	if ms.FirstDate.IsZero() || transaction.DateTime.Before(ms.FirstDate) {
		ms.FirstDate = transaction.DateTime
	}
	if ms.LastDate.IsZero() || transaction.DateTime.After(ms.LastDate) {
		ms.LastDate = transaction.DateTime
	}
}

// UpdateError actualiza las estadísticas para una transacción con error
func (ms *MigrationStats) UpdateError(lineNumber int, err error) {
	ms.ErrorRecords++
	errorMsg := fmt.Sprintf("Line %d: %v", lineNumber, err)
	ms.Errors = append(ms.Errors, errorMsg)
}

// ProcessCSV procesa un archivo CSV y migra las transacciones a la base de datos
func (ms *MigrationService) ProcessCSV(reader io.Reader) (*MigrationStats, error) {
	// Capturar tiempo de inicio
	startTime := time.Now()

	csvReader := csv.NewReader(reader)

	// Leer todas las líneas del CSV
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	// Verificar que tenga el header esperado
	header := records[0]
	expectedHeader := []string{"id", "user_id", "amount", "datetime"}
	if !ms.validateHeader(header, expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header. Expected: %v, Got: %v", expectedHeader, header)
	}

	// Inicializar estadísticas en línea
	stats := NewMigrationStats()
	stats.TotalRecords = len(records) - 1 // Excluir header

	// Procesar cada línea de datos (saltar header) - ESTADÍSTICAS EN LÍNEA
	for i, record := range records[1:] {
		lineNumber := i + 2 // +2 porque empezamos desde línea 2

		// Parsear transacción
		transaction, err := ms.parseTransaction(record, lineNumber)
		if err != nil {
			stats.UpdateError(lineNumber, err)
			fmt.Printf("Error parsing record at line %d: %v\n", lineNumber, err)
			continue
		}

		// Guardar en la base de datos mock
		savedTransaction, err := ms.database.SaveTransaction(transaction)
		if err != nil {
			stats.UpdateError(lineNumber, err)
			fmt.Printf("Error saving transaction at line %d: %v\n", lineNumber, err)
			continue
		}

		// Actualizar estadísticas en línea (NO almacenar en memoria)
		stats.UpdateSuccess(savedTransaction)
	}

	// Calcular tiempo de procesamiento real
	processingTime := time.Since(startTime)

	// Enviar reporte de migración (asíncrono)
	if ms.reportService != nil {
		go func() {
			report := ms.generateMigrationReportFromStats(stats, "uploaded_file.csv", 0, processingTime)
			ms.reportService.SendMigrationReport(report)
		}()
	}

	return stats, nil
}

// validateHeader verifica que el header del CSV sea correcto
func (ms *MigrationService) validateHeader(header, expected []string) bool {
	if len(header) != len(expected) {
		return false
	}

	for i, col := range header {
		if col != expected[i] {
			return false
		}
	}

	return true
}

// parseTransaction convierte una línea del CSV en una transacción
func (ms *MigrationService) parseTransaction(record []string, lineNumber int) (models.UserTransaction, error) {
	if len(record) != 4 {
		return models.UserTransaction{}, fmt.Errorf("invalid number of columns at line %d", lineNumber)
	}

	// Parsear ID
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid ID at line %d: %v", lineNumber, err)
	}

	// Parsear UserID
	userID, err := strconv.Atoi(record[1])
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid user_id at line %d: %v", lineNumber, err)
	}

	// Parsear Amount
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return models.UserTransaction{}, fmt.Errorf("invalid amount at line %d: %v", lineNumber, err)
	}

	// Parsear DateTime
	datetime, err := time.Parse("2006-01-02 15:04:05", record[3])
	if err != nil {
		// Intentar con otro formato común
		datetime, err = time.Parse("2006-01-02T15:04:05", record[3])
		if err != nil {
			// Intentar con formato de fecha solamente
			datetime, err = time.Parse("2006-01-02", record[3])
			if err != nil {
				return models.UserTransaction{}, fmt.Errorf("invalid datetime at line %d: %v", lineNumber, err)
			}
		}
	}

	return models.UserTransaction{
		ID:       id,
		UserID:   userID,
		Amount:   amount,
		DateTime: datetime,
	}, nil
}

// generateMigrationReportFromStats genera un reporte basado en estadísticas en línea
func (ms *MigrationService) generateMigrationReportFromStats(stats *MigrationStats, filename string, fileSize int64, processingTime time.Duration) *models.MigrationReport {
	// Calcular promedio basado en transacciones exitosas
	averageAmount := float64(0)
	if stats.SuccessRecords > 0 {
		averageAmount = stats.TotalAmount / float64(stats.SuccessRecords)
	}

	// Generar archivo CSV de errores si hay errores
	var errorFileCSV string
	if len(stats.Errors) > 0 && ms.reportService != nil {
		if errorPath, err := ms.reportService.GenerateErrorCSV(stats.Errors, filename); err == nil {
			errorFileCSV = errorPath
		}
	}

	report := &models.MigrationReport{
		Timestamp:      time.Now(),
		Filename:       filename,
		FileSize:       fileSize,
		TotalRecords:   stats.TotalRecords,
		SuccessRecords: stats.SuccessRecords,
		ErrorRecords:   stats.ErrorRecords,
		ProcessingTime: processingTime,
		UsersAffected:  len(stats.UsersAffected),
		TotalAmount:    stats.TotalAmount,
		AverageAmount:  averageAmount,
		LargestAmount:  stats.LargestAmount,
		SmallestAmount: stats.SmallestAmount,
		Errors:         stats.Errors,
		ErrorFileCSV:   errorFileCSV,
	}

	// Configurar rango de fechas
	report.DateRange.From = stats.FirstDate
	report.DateRange.To = stats.LastDate

	return report
}
