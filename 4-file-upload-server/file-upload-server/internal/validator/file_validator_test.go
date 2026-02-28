package validator

import (
	"bytes"
	"testing"
)

func TestValidateImage(t *testing.T) {
	tests := []struct {
		name      string
		content   []byte
		expectErr bool
	}{
		{
			name:      "valid jpeg",
			content:   append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 509)...),
			expectErr: false,
		},
		{
			name:      "valid png",
			content:   append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, make([]byte, 504)...),
			expectErr: false,
		},
		{
			name:      "invalid text file",
			content:   []byte("hello world"),
			expectErr: true,
		},
		{
			name:      "empty file",
			content:   []byte{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reader := bytes.NewReader(tt.content)

			err := ValidateImage(reader)

			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
