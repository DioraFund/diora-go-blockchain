package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger represents the application logger
type Logger struct {
	*logrus.Logger
	name string
}

// NewLogger creates a new logger instance
func NewLogger(name string) *Logger {
	log := logrus.New()

	// Set default configuration
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	return &Logger{
		Logger: log,
		name:   name,
	}
}

// NewLoggerWithConfig creates a new logger with custom configuration
func NewLoggerWithConfig(name string, level string, format string, output string, maxSize int) *Logger {
	log := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// Set formatter
	switch format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output
	switch output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		// Create logs directory
		logDir := filepath.Join(os.Getenv("HOME"), ".diora", "logs")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.SetOutput(os.Stdout)
		} else {
			logFile := filepath.Join(logDir, name+".log")
			if maxSize > 0 {
				log.SetOutput(&lumberjack.Logger{
					Filename:   logFile,
					MaxSize:    maxSize,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
				})
			} else {
				file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					log.SetOutput(os.Stdout)
				} else {
					log.SetOutput(file)
				}
			}
		}
	default:
		log.SetOutput(os.Stdout)
	}

	return &Logger{
		Logger: log,
		name:   name,
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.WithField(key, value),
		name:   l.name,
	}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		Logger: l.Logger.WithFields(fields),
		name:   l.name,
	}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.WithError(err),
		name:   l.name,
	}
}

// GetWriter returns the logger's writer
func (l *Logger) GetWriter() io.Writer {
	return l.Logger.Out
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	l.Logger.SetLevel(logLevel)
	return nil
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

// SetOutput sets the output destination
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

// SetFormatter sets the log formatter
func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.SetFormatter(formatter)
}
