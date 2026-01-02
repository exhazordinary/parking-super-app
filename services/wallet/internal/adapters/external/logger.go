package external

import (
	"log"
	"os"

	"github.com/parking-super-app/services/wallet/internal/ports"
)

// StdLogger is a simple logger that writes to stdout.
// In production, this should be replaced with a structured logger like Zap or Zerolog.
type StdLogger struct {
	logger *log.Logger
	fields []ports.Field
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		fields: nil,
	}
}

func (l *StdLogger) formatFields() string {
	if len(l.fields) == 0 {
		return ""
	}
	result := " "
	for _, f := range l.fields {
		result += f.Key + "=" + formatValue(f.Value) + " "
	}
	return result
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case error:
		return val.Error()
	default:
		return ""
	}
}

func (l *StdLogger) Debug(msg string, fields ...ports.Field) {
	l.logger.Printf("[DEBUG] %s%s", msg, l.formatWithFields(fields))
}

func (l *StdLogger) Info(msg string, fields ...ports.Field) {
	l.logger.Printf("[INFO] %s%s", msg, l.formatWithFields(fields))
}

func (l *StdLogger) Warn(msg string, fields ...ports.Field) {
	l.logger.Printf("[WARN] %s%s", msg, l.formatWithFields(fields))
}

func (l *StdLogger) Error(msg string, fields ...ports.Field) {
	l.logger.Printf("[ERROR] %s%s", msg, l.formatWithFields(fields))
}

func (l *StdLogger) WithFields(fields ...ports.Field) ports.Logger {
	newLogger := &StdLogger{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
	return newLogger
}

func (l *StdLogger) formatWithFields(fields []ports.Field) string {
	allFields := append(l.fields, fields...)
	if len(allFields) == 0 {
		return ""
	}
	result := " "
	for _, f := range allFields {
		result += f.Key + "=" + formatValue(f.Value) + " "
	}
	return result
}
