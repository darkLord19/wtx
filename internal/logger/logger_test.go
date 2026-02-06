package logger

import (
	"os"
	"strings"
	"testing"
)

func setupLoggerEnv(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	os.Setenv("XDG_CACHE_HOME", tmpDir)
	// For macOS fallback
	os.Setenv("HOME", tmpDir)

	return tmpDir
}

func TestInit(t *testing.T) {
	// Save env
	origCache := os.Getenv("XDG_CACHE_HOME")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CACHE_HOME", origCache)
		os.Setenv("HOME", origHome)
		Close() // ensure closed
	}()

	setupLoggerEnv(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify log file exists
	logPath := GetLogPath()
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log file not created at %s", logPath)
	}

	// Test logging
	Info("Test info message")
	Error("Test error message")

	// Close to flush
	if err := Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Read file content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "[INFO] Test info message") {
		t.Error("Log file missing info message")
	}
	if !strings.Contains(s, "[ERROR] Test error message") {
		t.Error("Log file missing error message")
	}
}

func TestInitQuiet(t *testing.T) {
	if err := InitQuiet(); err != nil {
		t.Fatalf("InitQuiet() error = %v", err)
	}

	// Should not panic
	Info("Quiet info")
}
