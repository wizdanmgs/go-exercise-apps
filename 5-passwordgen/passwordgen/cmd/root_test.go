package cmd

import (
	"strings"
	"testing"
)

func TestBuildCharset(t *testing.T) {
	tests := []struct {
		name           string
		includeSymbols bool
		expectSymbols  bool
	}{
		{
			name:           "with symbols",
			includeSymbols: true,
			expectSymbols:  true,
		},
		{
			name:           "without symbols",
			includeSymbols: false,
			expectSymbols:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charset := buildCharset(tt.includeSymbols)

			hasSymbol := strings.ContainsAny(charset, symbols)

			if tt.expectSymbols && !hasSymbol {
				t.Errorf("expected symbols in charset, got none")
			}

			if !tt.expectSymbols && hasSymbol {
				t.Errorf("expected no symbols in charset, but found some")
			}
		})
	}
}

func TestGeneratePassword_Length(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 8", 8},
		{"length 16", 16},
		{"length 32", 32},
	}

	charset := buildCharset(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw, err := generatePassword(tt.length, charset)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(pw) != tt.length {
				t.Errorf("expected length %d, got %d", tt.length, len(pw))
			}
		})
	}
}

func TestGeneratePassword_CharsetCompliance(t *testing.T) {
	tests := []struct {
		name           string
		includeSymbols bool
	}{
		{"with symbols", true},
		{"without symbols", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charset := buildCharset(tt.includeSymbols)

			pw, err := generatePassword(50, charset)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, r := range pw {
				if !strings.ContainsRune(charset, r) {
					t.Errorf("found character not in charset: %q", r)
				}
			}
		})
	}
}
