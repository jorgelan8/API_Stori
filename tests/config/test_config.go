package config

import (
	"os"
	"sync"
	"time"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	BaseURL     string
	PathAPI     string
	HostAPI     string
	Timeout     time.Duration
	Port        string
	Host        string
	TestDataDir string
	MaxRetries  int
	RetryDelay  time.Duration
}

var (
	configInstance *TestConfig
	configOnce     sync.Once
)

// getConfig returns test configuration with sensible defaults (singleton)
func getConfig() *TestConfig {
	configOnce.Do(func() {
		baseURL := os.Getenv("TEST_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		port := os.Getenv("TEST_PORT")
		if port == "" {
			port = "8080"
		}

		host := os.Getenv("TEST_HOST")
		if host == "" {
			host = "localhost"
		}

		testDataDir := os.Getenv("TEST_DATA_DIR")
		if testDataDir == "" {
			testDataDir = "tests/testdata"
		}

		// Get PathAPI from environment or use default
		pathAPI := os.Getenv("TEST_PATH_API")
		if pathAPI == "" {
			pathAPI = "/api/v1"
		}

		// Build HostAPI from BaseURL + PathAPI
		hostAPI := baseURL + pathAPI

		configInstance = &TestConfig{
			BaseURL:     baseURL,
			PathAPI:     pathAPI,
			HostAPI:     hostAPI,
			Timeout:     30 * time.Second,
			Port:        port,
			Host:        host,
			TestDataDir: testDataDir,
			MaxRetries:  3,
			RetryDelay:  1 * time.Second,
		}
	})
	return configInstance
}

// GetTestConfig returns test configuration (for backward compatibility)
func GetTestConfig() *TestConfig {
	return getConfig()
}

// GetBaseURL returns the base URL for tests
func GetBaseURL() string {
	return getConfig().BaseURL
}

// GetTimeout returns the timeout for tests
func GetTimeout() time.Duration {
	return getConfig().Timeout
}

// GetHostAPI returns the HostAPI URL for tests (BaseURL + PathAPI)
func GetHostAPI() string {
	return getConfig().HostAPI
}

// GetPathAPI returns the PathAPI for tests (/api/v1)
func GetPathAPI() string {
	return getConfig().PathAPI
}

// GetTestDataDir returns the test data directory
func GetTestDataDir() string {
	return getConfig().TestDataDir
}

// GetMaxRetries returns the maximum number of retries
func GetMaxRetries() int {
	return getConfig().MaxRetries
}

// GetRetryDelay returns the delay between retries
func GetRetryDelay() time.Duration {
	return getConfig().RetryDelay
}

// ResetConfig resets the singleton instance (useful for testing)
func ResetConfig() {
	configInstance = nil
	configOnce = sync.Once{}
}
