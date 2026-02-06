package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	// Logger is the global logger instance
	Logger *log.Logger
	// logFile is the log file handle
	logFile *os.File
)

// Init initializes the logger
func Init() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to temp dir
		cacheDir = os.TempDir()
	}

	logDir := filepath.Join(cacheDir, "wtx")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "wtx.log")
	
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logFile = f
	Logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
	
	return nil
}

// Close closes the log file
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// GetLogPath returns the path to the log file
func GetLogPath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "wtx", "wtx.log")
}

// InitQuiet initializes a logger that writes to nowhere (for testing)
func InitQuiet() error {
	Logger = log.New(io.Discard, "", log.LstdFlags)
	return nil
}
