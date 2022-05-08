package resizer

import (
	"bytes"
	"fmt"
	"io"
)

var _ Resizer = (*ImageResizer)(nil)

type Resizer interface {
	Resize(i io.Reader, w, h int) ([]byte, error)
}

func NewImageResizer() *ImageResizer {
	return &ImageResizer{
		&imagingProcessor{},
	}
}

type ImageResizer struct {
	processor ImageProcessor
}

func (r *ImageResizer) WithProcessor(processor ImageProcessor) *ImageResizer {
	r.processor = processor

	return r
}

func (r *ImageResizer) Resize(reader io.Reader, w, h int) ([]byte, error) {
	img, err := r.processor.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("ImageResizer decode: %w", err)
	}

	wi, hi := img.Bounds().Dx(), img.Bounds().Dy()
	targetRatio := float64(w) / float64(h)
	currentRatio := float64(wi) / float64(hi)

	if targetRatio != currentRatio {
		var cropW, cropH int
		if targetRatio > currentRatio {
			cropH = int(float64(wi) / targetRatio)
			cropW = wi
		} else {
			cropW = int(float64(hi) * targetRatio)
			cropH = hi
		}

		img = r.processor.Crop(img, cropW, cropH)
	}

	img = r.processor.Resize(img, w, h)

	buff := new(bytes.Buffer)
	if err := r.processor.Encode(img, buff); err != nil {
		return nil, fmt.Errorf("ImageResizer encode: %w", err)
	}

	return buff.Bytes(), nil
}
