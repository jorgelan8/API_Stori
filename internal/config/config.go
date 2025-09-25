package config

import (
	"api-stori/internal/models"
	"os"
	"strconv"
	"strings"
)

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() *Config {
	return &Config{
		App:    loadAppConfig(),
		Email:  loadEmailConfig(),
		Report: loadReportConfig(),
	}
}

// Config contiene toda la configuración de la aplicación
type Config struct {
	App    AppConfig
	Email  EmailConfig
	Report ReportConfig
}

// AppConfig configuración de la aplicación
type AppConfig struct {
	Port string
	Host string
}

// EmailConfig configuración de email
type EmailConfig struct {
	SMTPHost  string
	SMTPPort  int
	Username  string
	Password  string
	FromEmail string
	ToEmails  []string
}

// ReportConfig configuración de reportes
type ReportConfig struct {
	Channels []models.ReportChannel
	Subject  string
}

// loadAppConfig carga la configuración de la aplicación
func loadAppConfig() AppConfig {
	return AppConfig{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
	}
}

// loadEmailConfig carga la configuración de email
func loadEmailConfig() EmailConfig {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return EmailConfig{
		SMTPHost:  getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:  port,
		Username:  getEnv("SMTP_USERNAME", ""),
		Password:  getEnv("SMTP_PASSWORD", ""),
		FromEmail: getEnv("FROM_EMAIL", ""),
		ToEmails:  parseEmailList(getEnv("TO_EMAILS", "admin@api-stori.com")),
	}
}

// loadReportConfig carga la configuración de reportes
func loadReportConfig() ReportConfig {
	channelsStr := getEnv("REPORT_CHANNELS", "email,log")
	channels := parseReportChannels(channelsStr)

	return ReportConfig{
		Channels: channels,
		Subject:  getEnv("REPORT_SUBJECT", "Migration Report - API Stori"),
	}
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseEmailList parsea una lista de emails separados por comas
func parseEmailList(emailsStr string) []string {
	emails := strings.Split(emailsStr, ",")
	var result []string
	for _, email := range emails {
		email = strings.TrimSpace(email)
		if email != "" {
			result = append(result, email)
		}
	}
	return result
}

// parseReportChannels parsea una lista de canales de reporte
func parseReportChannels(channelsStr string) []models.ReportChannel {
	channels := strings.Split(channelsStr, ",")
	var result []models.ReportChannel
	for _, channel := range channels {
		channel = strings.TrimSpace(channel)
		switch channel {
		case "email":
			result = append(result, models.EmailChannel)
		case "webhook":
			result = append(result, models.WebhookChannel)
		case "log":
			result = append(result, models.LogChannel)
		}
	}
	return result
}

// ToReportConfig convierte la configuración a models.ReportConfig
func (c *Config) ToReportConfig() *models.ReportConfig {
	return &models.ReportConfig{
		Channels: c.Report.Channels,
		Email: models.EmailConfig{
			SMTPHost:  c.Email.SMTPHost,
			SMTPPort:  c.Email.SMTPPort,
			Username:  c.Email.Username,
			Password:  c.Email.Password,
			FromEmail: c.Email.FromEmail,
			ToEmails:  c.Email.ToEmails,
			Subject:   c.Report.Subject,
		},
	}
}
