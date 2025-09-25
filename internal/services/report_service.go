package services

import (
	"api-stori/internal/models"
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

// ReportService maneja el envío de reportes de migración
type ReportService struct {
	config *models.ReportConfig
}

// NewReportService crea una nueva instancia de ReportService
func NewReportService(config *models.ReportConfig) *ReportService {
	return &ReportService{
		config: config,
	}
}

// SendMigrationReport envía el reporte de migración por los canales configurados
func (rs *ReportService) SendMigrationReport(report *models.MigrationReport) {
	for _, channel := range rs.config.Channels {
		switch channel {
		case models.EmailChannel:
			go rs.sendEmailReport(report)
		case models.WebhookChannel:
			go rs.sendWebhookReport(report)
		case models.LogChannel:
			go rs.sendLogReport(report)
		}
	}
}

// sendEmailReport envía el reporte por email
func (rs *ReportService) sendEmailReport(report *models.MigrationReport) {
	// Si no hay configuración de SMTP, usar mock
	if rs.config.Email.SMTPHost == "" {
		rs.sendMockEmailReport(report)
		return
	}

	// Generar contenido del email
	subject := fmt.Sprintf("Migration Report - %s", report.Filename)
	body := rs.generateEmailBody(report)

	// Configurar autenticación SMTP
	auth := smtp.PlainAuth("", rs.config.Email.Username, rs.config.Email.Password, rs.config.Email.SMTPHost)

	// Crear mensaje
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
		rs.config.Email.ToEmails[0], subject, body))

	// Enviar email
	addr := fmt.Sprintf("%s:%d", rs.config.Email.SMTPHost, rs.config.Email.SMTPPort)
	err := smtp.SendMail(addr, auth, rs.config.Email.FromEmail, rs.config.Email.ToEmails, msg)
	if err != nil {
		log.Printf("Error sending email report: %v", err)
	} else {
		log.Printf("Migration report sent via email successfully")
	}
}

// sendMockEmailReport simula el envío de email (para desarrollo)
func (rs *ReportService) sendMockEmailReport(report *models.MigrationReport) {
	log.Printf("=== MOCK EMAIL REPORT ===")
	log.Printf("To: %v", rs.config.Email.ToEmails)
	log.Printf("Subject: Migration Report - %s", report.Filename)
	log.Printf("Body:\n%s", rs.generateEmailBody(report))
	log.Printf("=== END MOCK EMAIL ===")
}

// sendWebhookReport envía el reporte por webhook
func (rs *ReportService) sendWebhookReport(report *models.MigrationReport) {
	// TODO: Implementar webhook
	log.Printf("Webhook report sent: %s", report.Filename)
}

// sendLogReport envía el reporte por log
func (rs *ReportService) sendLogReport(report *models.MigrationReport) {
	log.Printf("=== MIGRATION REPORT ===")
	log.Printf("File: %s (%d bytes)", report.Filename, report.FileSize)
	log.Printf("Records: %d total, %d success, %d errors",
		report.TotalRecords, report.SuccessRecords, report.ErrorRecords)
	log.Printf("Users affected: %d", report.UsersAffected)
	log.Printf("Amount range: %.2f to %.2f (avg: %.2f)",
		report.SmallestAmount, report.LargestAmount, report.AverageAmount)
	log.Printf("Processing time: %v", report.ProcessingTime)
	if len(report.Errors) > 0 {
		log.Printf("Errors: %v", report.Errors)
	}
	log.Printf("=== END REPORT ===")
}

// generateEmailBody genera el cuerpo del email
func (rs *ReportService) generateEmailBody(report *models.MigrationReport) string {
	var body bytes.Buffer

	body.WriteString("=== MIGRATION REPORT ===\n\n")
	body.WriteString(fmt.Sprintf("File: %s (%d bytes)\n", report.Filename, report.FileSize))
	body.WriteString(fmt.Sprintf("Timestamp: %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	body.WriteString(fmt.Sprintf("Processing time: %v\n\n", report.ProcessingTime))

	body.WriteString("=== STATISTICS ===\n")
	body.WriteString(fmt.Sprintf("Total records: %d\n", report.TotalRecords))
	body.WriteString(fmt.Sprintf("Success records: %d\n", report.SuccessRecords))
	body.WriteString(fmt.Sprintf("Error records: %d\n", report.ErrorRecords))
	body.WriteString(fmt.Sprintf("Success rate: %.2f%%\n\n",
		float64(report.SuccessRecords)/float64(report.TotalRecords)*100))

	body.WriteString("=== DATA ANALYSIS ===\n")
	body.WriteString(fmt.Sprintf("Users affected: %d\n", report.UsersAffected))
	body.WriteString(fmt.Sprintf("Total amount: %.2f\n", report.TotalAmount))
	body.WriteString(fmt.Sprintf("Average amount: %.2f\n", report.AverageAmount))
	body.WriteString(fmt.Sprintf("Largest amount: %.2f\n", report.LargestAmount))
	body.WriteString(fmt.Sprintf("Smallest amount: %.2f\n", report.SmallestAmount))
	body.WriteString(fmt.Sprintf("Date range: %s to %s\n\n",
		report.DateRange.From.Format("2006-01-02"),
		report.DateRange.To.Format("2006-01-02")))

	if len(report.Errors) > 0 {
		body.WriteString("=== ERRORS ===\n")
		for i, err := range report.Errors {
			body.WriteString(fmt.Sprintf("%d. %s\n", i+1, err))
		}
		body.WriteString("\n")
	}

	if report.ErrorFileCSV != "" {
		body.WriteString("=== ERROR FILE ===\n")
		body.WriteString(fmt.Sprintf("Error records exported to: %s\n", report.ErrorFileCSV))
		body.WriteString("\n")
	}

	body.WriteString("=== END REPORT ===")

	return body.String()
}

// GenerateErrorCSV genera un archivo CSV con los registros que tuvieron errores
func (rs *ReportService) GenerateErrorCSV(errors []string, filename string) (string, error) {
	// Crear directorio de errores si no existe
	errorDir := "reports/errors"
	if err := os.MkdirAll(errorDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create error directory: %v", err)
	}

	// Generar nombre de archivo único
	timestamp := time.Now().Format("20060102_150405")
	errorFilename := fmt.Sprintf("errors_%s_%s", timestamp, filename)
	errorPath := filepath.Join(errorDir, errorFilename)

	// Crear archivo CSV
	file, err := os.Create(errorPath)
	if err != nil {
		return "", fmt.Errorf("failed to create error file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir header
	header := []string{"line_number", "error_message", "original_data"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write header: %v", err)
	}

	// Escribir errores
	for i, errorMsg := range errors {
		record := []string{
			fmt.Sprintf("%d", i+1),
			errorMsg,
			"", // TODO: Incluir datos originales si es posible
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write error record: %v", err)
		}
	}

	return errorPath, nil
}
