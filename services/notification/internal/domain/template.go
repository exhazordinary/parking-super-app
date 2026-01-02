package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Template represents a notification template
type Template struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Channel   Channel           `json:"channel"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Variables []string          `json:"variables"`
	IsActive  bool              `json:"is_active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// NewTemplate creates a new notification template
func NewTemplate(name string, channel Channel, notifType, title, body string) *Template {
	now := time.Now().UTC()
	return &Template{
		ID:        uuid.New(),
		Name:      name,
		Channel:   channel,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Variables: extractVariables(body),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Render renders the template with provided variables
func (t *Template) Render(vars map[string]string) (title, body string) {
	title = t.Title
	body = t.Body

	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		title = strings.ReplaceAll(title, placeholder, value)
		body = strings.ReplaceAll(body, placeholder, value)
	}

	return title, body
}

// Deactivate disables the template
func (t *Template) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now().UTC()
}

// extractVariables finds all {{variable}} placeholders in text
func extractVariables(text string) []string {
	var vars []string
	start := 0
	for {
		startIdx := strings.Index(text[start:], "{{")
		if startIdx == -1 {
			break
		}
		startIdx += start
		endIdx := strings.Index(text[startIdx:], "}}")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx
		varName := text[startIdx+2 : endIdx]
		vars = append(vars, strings.TrimSpace(varName))
		start = endIdx + 2
	}
	return vars
}
