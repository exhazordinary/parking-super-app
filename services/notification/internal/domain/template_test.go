package domain

import (
	"testing"
)

func TestNewTemplate(t *testing.T) {
	template := NewTemplate(
		"payment-success",
		ChannelPush,
		"payment.success",
		"Payment Successful",
		"Your payment of {{amount}} has been processed.",
	)

	if template.Name != "payment-success" {
		t.Errorf("expected name payment-success, got %s", template.Name)
	}
	if template.Channel != ChannelPush {
		t.Errorf("expected channel push, got %s", template.Channel)
	}
	if !template.IsActive {
		t.Error("new template should be active")
	}
	if len(template.Variables) != 1 {
		t.Errorf("expected 1 variable, got %d", len(template.Variables))
	}
	if template.Variables[0] != "amount" {
		t.Errorf("expected variable 'amount', got %s", template.Variables[0])
	}
}

func TestTemplate_Render(t *testing.T) {
	template := NewTemplate(
		"session-ended",
		ChannelSMS,
		"session.ended",
		"Parking Ended",
		"Your parking session at {{location}} has ended. Total: {{amount}}",
	)

	vars := map[string]string{
		"location": "KLCC",
		"amount":   "RM 15.00",
	}

	title, body := template.Render(vars)

	if title != "Parking Ended" {
		t.Errorf("expected title 'Parking Ended', got %s", title)
	}
	expectedBody := "Your parking session at KLCC has ended. Total: RM 15.00"
	if body != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestTemplate_Deactivate(t *testing.T) {
	template := NewTemplate("test", ChannelEmail, "test", "Test", "Test body")

	if !template.IsActive {
		t.Error("new template should be active")
	}

	template.Deactivate()

	if template.IsActive {
		t.Error("template should be inactive after deactivation")
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		text     string
		expected []string
	}{
		{"Hello {{name}}", []string{"name"}},
		{"{{a}} and {{b}}", []string{"a", "b"}},
		{"No variables here", nil},
		{"{{one}} {{two}} {{three}}", []string{"one", "two", "three"}},
		{"{{ spaced }}", []string{"spaced"}},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := extractVariables(tt.text)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d variables, got %d", len(tt.expected), len(result))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("expected variable %s, got %s", tt.expected[i], v)
				}
			}
		})
	}
}
