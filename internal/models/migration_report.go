package models

import (
	"time"
)

// MigrationReport representa el reporte de migración
type MigrationReport struct {
	// Información básica
	Timestamp time.Time `json:"timestamp"`
	Filename  string    `json:"filename"`
	FileSize  int64     `json:"file_size"`

	// Estadísticas de procesamiento
	TotalRecords   int           `json:"total_records"`
	SuccessRecords int           `json:"success_records"`
	ErrorRecords   int           `json:"error_records"`
	ProcessingTime time.Duration `json:"processing_time"`

	// Análisis de datos
	UsersAffected  int     `json:"users_affected"`
	TotalAmount    float64 `json:"total_amount"`
	AverageAmount  float64 `json:"average_amount"`
	LargestAmount  float64 `json:"largest_amount"`
	SmallestAmount float64 `json:"smallest_amount"`

	// Distribución temporal
	DateRange struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	} `json:"date_range"`

	// Errores específicos
	Errors []string `json:"errors,omitempty"`

	// Archivo de errores (CSV)
	ErrorFileCSV string `json:"error_file_csv,omitempty"`
}

// ReportChannel representa los canales de notificación
type ReportChannel string

const (
	EmailChannel   ReportChannel = "email"
	WebhookChannel ReportChannel = "webhook"
	LogChannel     ReportChannel = "log"
)

// ReportConfig configuración para envío de reportes
type ReportConfig struct {
	Channels []ReportChannel `json:"channels"`
	Email    EmailConfig     `json:"email"`
	Webhook  WebhookConfig   `json:"webhook"`
}

// EmailConfig configuración de email
type EmailConfig struct {
	SMTPHost     string   `json:"smtp_host"`
	SMTPPort     int      `json:"smtp_port"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	FromEmail    string   `json:"from_email"`
	ToEmails     []string `json:"to_emails"`
	Subject      string   `json:"subject"`
	TemplatePath string   `json:"template_path"`
}

// WebhookConfig configuración de webhook
type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
}
