package external

import (
	"log"
	"os"

	"github.com/parking-super-app/services/notification/internal/ports"
)

type StdLogger struct {
	logger *log.Logger
}

func NewStdLogger() *StdLogger {
	return &StdLogger{logger: log.New(os.Stdout, "", log.LstdFlags)}
}

func (l *StdLogger) Debug(msg string, fields ...ports.Field) {
	l.logger.Printf("[DEBUG] %s %s", msg, formatFields(fields))
}

func (l *StdLogger) Info(msg string, fields ...ports.Field) {
	l.logger.Printf("[INFO] %s %s", msg, formatFields(fields))
}

func (l *StdLogger) Warn(msg string, fields ...ports.Field) {
	l.logger.Printf("[WARN] %s %s", msg, formatFields(fields))
}

func (l *StdLogger) Error(msg string, fields ...ports.Field) {
	l.logger.Printf("[ERROR] %s %s", msg, formatFields(fields))
}

func formatFields(fields []ports.Field) string {
	if len(fields) == 0 {
		return ""
	}
	result := ""
	for _, f := range fields {
		switch v := f.Value.(type) {
		case string:
			result += f.Key + "=" + v + " "
		case error:
			result += f.Key + "=" + v.Error() + " "
		}
	}
	return result
}
