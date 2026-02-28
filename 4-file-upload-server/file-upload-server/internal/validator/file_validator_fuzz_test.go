package validator

import (
	"bytes"
	"net/http"
	"testing"
)

func FuzzValidateImage(f *testing.F) {
	// ---- Seed corpus (important!) ----

	f.Add([]byte{})                                                                             // empty
	f.Add([]byte("hello world"))                                                                // random text
	f.Add(append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 509)...))                               // jpeg
	f.Add(append([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, make([]byte, 504)...)) // png

	// ---- Fuzz function ----

	f.Fuzz(func(t *testing.T, data []byte) {
		reader := bytes.NewReader(data)

		// Ensure it never panics
		err := ValidateImage(reader)

		// Optional: If it returns null, ensure MIME type is allowed
		if err == nil {
			// If validation passed, it must be one of allowed types.
			// Re-run detection to confirm logic consistency.
			if _, rErr := reader.Seek(0, 0); rErr != nil {
				t.Fatal(rErr.Error())
			}

			buf := make([]byte, 512)
			if _, rErr := reader.Read(buf); rErr != nil {
				t.Fatal(rErr.Error())
			}

			mime := http.DetectContentType(buf)

			if !allowedTypes[mime] {
				t.Fatalf("validator allowed unexpected mime: %s", mime)
			}
		}
	})

}
