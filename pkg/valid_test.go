package pkg

import (
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"ValidRussian", "Иван", true},
		{"ValidEnglish", "John", true},
		{"LowercaseStart", "иван", false},
		{"ContainsNumber", "Иван1", false},
		{"ContainsSpace", "Иван Иванов", false},
		{"Empty", "", false},
		{"OnlySymbols", "@#$%", false},
		{"SingleUpper", "A", true},
		{"HyphenName", "Жан-Поль", false}, // не пройдёт, если дефис не допускается
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidName(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
