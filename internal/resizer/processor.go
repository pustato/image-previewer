package resizer

import (
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

var _ ImageProcessor = (*imagingProcessor)(nil)

type ImageProcessor interface {
	Decode(reader io.Reader) (image.Image, error)
	Crop(img image.Image, width, height int) image.Image
	Resize(img image.Image, width, height int) image.Image
	Encode(img image.Image, writer io.Writer) error
}

type imagingProcessor struct{}

func (i *imagingProcessor) Decode(reader io.Reader) (image.Image, error) {
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("imaging decode: %w", err)
	}

	return img, nil
}

func (i *imagingProcessor) Crop(img image.Image, width, height int) image.Image {
	return imaging.CropCenter(img, width, height)
}

func (i *imagingProcessor) Resize(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func (i *imagingProcessor) Encode(img image.Image, writer io.Writer) error {
	if err := imaging.Encode(writer, img, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
		return fmt.Errorf("imaging encode: %w", err)
	}

	return nil
}
