package validator

import (
	"bytes"
	"image"
	"testing"
)

func FuzzValidateImage(f *testing.F) {
	// ---- Seed corpus (important!) ----

	f.Add([]byte{})              // empty
	f.Add([]byte("hello world")) // random text

	// Add real valid JPEG Seed
	validJPEG := GenerateJPEG()
	f.Add(validJPEG)

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

			_, _, decodeErr := image.DecodeConfig(reader)
			if decodeErr != nil {
				t.Fatalf("validator accepted undecodable image")
			}
		}
	})

}
