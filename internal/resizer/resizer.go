package resizer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/disintegration/imaging"
)

var _ Resizer = (*resizer)(nil)

type Resizer interface {
	Resize(i io.Reader, w, h int) ([]byte, error)
}

func New() Resizer {
	return &resizer{
		&imagingProcessor{},
	}
}

type resizer struct {
	processor ImageProcessor
}

func (r *resizer) WithProcessor(processor ImageProcessor) {
	r.processor = processor
}

func (r *resizer) Resize(reader io.Reader, w, h int) ([]byte, error) {
	img, err := r.processor.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("resizer decode: %w", err)
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

	img = imaging.Resize(img, w, h, imaging.Lanczos)

	buff := new(bytes.Buffer)
	if err := imaging.Encode(buff, img, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
		return nil, fmt.Errorf("resizer encode: %w", err)
	}

	return buff.Bytes(), nil
}
