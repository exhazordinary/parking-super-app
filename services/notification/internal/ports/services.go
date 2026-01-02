package ports

import (
	"context"

	"github.com/parking-super-app/services/notification/internal/domain"
)

// Logger interface
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

func String(key, value string) Field { return Field{Key: key, Value: value} }
func Err(err error) Field            { return Field{Key: "error", Value: err} }
func Any(key string, val interface{}) Field { return Field{Key: key, Value: val} }

// PushProvider sends push notifications
type PushProvider interface {
	Send(ctx context.Context, req PushRequest) (*PushResponse, error)
}

type PushRequest struct {
	DeviceToken string
	Title       string
	Body        string
	Data        map[string]string
	Priority    string
}

type PushResponse struct {
	MessageID string
	Success   bool
	Error     string
}

// SMSProvider sends SMS messages
type SMSProvider interface {
	Send(ctx context.Context, req SMSRequest) (*SMSResponse, error)
}

type SMSRequest struct {
	PhoneNumber string
	Message     string
}

type SMSResponse struct {
	MessageID string
	Status    string
	Error     string
}

// EmailProvider sends emails
type EmailProvider interface {
	Send(ctx context.Context, req EmailRequest) (*EmailResponse, error)
}

type EmailRequest struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

type EmailResponse struct {
	MessageID string
	Status    string
	Error     string
}

// NotificationSender unified interface for sending via any channel
type NotificationSender interface {
	Send(ctx context.Context, notif *domain.Notification) error
}

// EventConsumer consumes events from message queue
type EventConsumer interface {
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	Close() error
}

type EventHandler func(ctx context.Context, event Event) error

type Event struct {
	Type    string
	Payload map[string]interface{}
}

// Common notification types
const (
	NotifTypePaymentSuccess   = "payment.success"
	NotifTypePaymentFailed    = "payment.failed"
	NotifTypeSessionStarted   = "session.started"
	NotifTypeSessionEnding    = "session.ending"
	NotifTypeSessionEnded     = "session.ended"
	NotifTypePromotion        = "promotion"
	NotifTypeAccountAlert     = "account.alert"
)
