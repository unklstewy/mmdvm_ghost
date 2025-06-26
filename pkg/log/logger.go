package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger = logrus.New()

// InitLogger initializes the logger with file rotation and log levels.
func InitLogger(logFile string, level string, maxSize int, maxBackups int, maxAge int) {
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set log output to a rotating file
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
	})

	// Set log format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}

// Info logs an informational message.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn logs a warning message.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Debug logs a debug message.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Fatal logs a fatal message and exits the application.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
	os.Exit(1)
}
