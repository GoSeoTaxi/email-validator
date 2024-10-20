package usecase

import (
	"testing"

	"github.com/GoSeoTaxi/email-validator/internal/adapter"

	"github.com/GoSeoTaxi/email-validator/internal/config"
)

func TestEmailValidator_Validate(t *testing.T) {
	cfg := &config.Config{
		DNSHosts: []string{"8.8.8.8"},
	}
	mockCache := adapter.NewMockCache()
	validator := NewEmailValidator(cfg, mockCache)

	tests := []struct {
		email   string
		isValid bool
		message string
	}{
		{"test@example.com", true, "Email is valid"},
		{"test@ya.ru", true, "Email is valid"},
		{"t1@gmail.com", true, "Email is valid"},
		{"invalidemail", false, "Invalid email format"},
		{"inva@lidemail", false, "Invalid email format"},
		{"user@nonexistentdomain.tld", false, "mail domain unavailable"},
		{"user@1.y1a.ru", false, "mail domain unavailable"},
		{"test@2.ya.ru", true, "Email is valid"},
	}

	for _, tt := range tests {
		isValid, msg := validator.Validate(tt.email)
		if isValid != tt.isValid || msg != tt.message {
			t.Errorf("Expected isValid=%v, message=%q; got isValid=%v, message=%q for email %s",
				tt.isValid, tt.message, isValid, msg, tt.email)
		}
	}
}
