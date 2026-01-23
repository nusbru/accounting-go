package validator

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.org", true},
		{"user+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"user@", false},
		{"user@.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := ValidateEmail(tt.email); got != tt.valid {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		uuid  string
		valid bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", true},
		{"invalid", false},
		{"123e4567-e89b-12d3-a456", false},
		{"123e4567-e89b-12d3-a456-426614174000-extra", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			if got := ValidateUUID(tt.uuid); got != tt.valid {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.uuid, got, tt.valid)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		{"value", false},
		{"  value  ", false},
		{"", true},
		{"   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := ValidateRequired(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequired(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		value   string
		min     int
		max     int
		wantErr bool
	}{
		{"abc", 1, 5, false},
		{"ab", 3, 5, true},  // too short
		{"abcdef", 1, 5, true}, // too long
		{"", 0, 5, false},
		{"exact", 5, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := ValidateLength(tt.value, "field", tt.min, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLength(%q, %d, %d) error = %v, wantErr %v", tt.value, tt.min, tt.max, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		value   float64
		wantErr bool
	}{
		{1.0, false},
		{0.001, false},
		{0, true},
		{-1.0, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := ValidatePositive(tt.value, "field")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositive(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}

func TestValidateOneOf(t *testing.T) {
	allowed := []string{"income", "expense", "transfer"}

	tests := []struct {
		value   string
		wantErr bool
	}{
		{"income", false},
		{"expense", false},
		{"transfer", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			err := ValidateOneOf(tt.value, "type", allowed)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOneOf(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
		})
	}
}
