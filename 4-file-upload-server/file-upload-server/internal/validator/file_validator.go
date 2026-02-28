package validator

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
)

var allowedFormats = map[string]bool{
	"jpeg": true,
	"png":  true,
	"gif":  true,
}

const (
	maxWidth  = 5000
	maxHeight = 5000
)

func ValidateImage(file io.ReadSeeker) error {
	// Decode only config (does NOT fully decode image pixels)
	cfg, format, err := image.DecodeConfig(file)
	if err != nil {
		return errors.New("invalid image file")
	}
	if !allowedFormats[format] {
		return errors.New("unsupported image format")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return errors.New("invalid image dimensions")
	}
	if cfg.Width > maxWidth || cfg.Height > maxHeight {
		return errors.New("image dimensions too large")
	}

	// Reset pointer for the processing
	_, err = file.Seek(0, io.SeekCurrent)
	return err
}

func GenerateJPEG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, image.White)
		}
	}

	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func GeneratePNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.White)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err) // safe in test helper
	}

	return buf.Bytes()
}
