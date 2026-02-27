package validator

import (
	"errors"
	"io"
	"net/http"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func ValidateImage(file io.Reader) error {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return err
	}

	fileType := http.DetectContentType(buffer)

	if !allowedTypes[fileType] {
		return errors.New("invalid file type")
	}

	return nil
}
