package resizer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"testing"

	mockresizer "github.com/pustato/image-previewer/internal/resizer/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var anyWriter = mock.MatchedBy(func(_ io.Writer) bool { return true })

func newImageStub(w, h int) *imageStub {
	return &imageStub{
		rect: image.Rect(0, 0, w, h),
	}
}

type imageStub struct {
	rect image.Rectangle
}

func (t *imageStub) ColorModel() color.Model {
	return nil
}

func (t *imageStub) Bounds() image.Rectangle {
	return t.rect
}

func (t *imageStub) At(x, y int) color.Color {
	return nil
}

func TestImageResizer_Resize_Crop(t *testing.T) {
	t.Parallel()

	testData := []struct {
		iw, ih, w, h, cropW, cropH int
		expectCrop                 bool
	}{
		{
			4000, 3000, 50, 50, 3000, 3000, true,
		},

		{
			1024, 768, 2000, 50, 1024, 25, true,
		},

		{
			1024, 768, 50, 1000, 38, 768, true,
		},

		{
			1024, 768, 1023, 768, 1023, 768, true,
		},

		{
			1024, 768, 1025, 768, 1024, 767, true,
		},

		{
			1024, 768, 1024, 767, 1024, 767, true,
		},

		{
			1024, 768, 1024, 769, 1022, 768, true,
		},

		{
			1000, 2000, 100, 200, 1000, 2000, false,
		},
	}

	for i, td := range testData {
		td := td

		t.Run(fmt.Sprintf("crop case %d", i), func(t *testing.T) {
			t.Parallel()

			img := newImageStub(td.iw, td.ih)
			croppedImg := newImageStub(td.cropW, td.cropH)
			resizedImg := newImageStub(td.w, td.h)
			in := new(bytes.Buffer)

			processor := &mockresizer.ImageProcessor{}
			processor.
				On("Decode", in).
				Once().
				Return(img, nil)

			if td.expectCrop {
				processor.
					On("Crop", img, td.cropW, td.cropH).
					Once().
					Return(croppedImg)
			} else {
				croppedImg = img
			}

			processor.
				On("Resize", croppedImg, td.w, td.h).
				Once().
				Return(resizedImg)

			processor.
				On("Encode", resizedImg, anyWriter).
				Once().
				Return(nil)

			resizer := NewImageResizer().WithProcessor(processor)

			_, err := resizer.Resize(in, td.w, td.h)
			require.NoError(t, err)
		})
	}
}

func TestImageResizer_Resize_Errors(t *testing.T) {
	t.Run("decode error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		in := new(bytes.Buffer)

		processor := &mockresizer.ImageProcessor{}
		processor.
			On("Decode", in).
			Once().
			Return(nil, expectedErr)

		resizer := NewImageResizer().WithProcessor(processor)

		_, err := resizer.Resize(in, 0, 0)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("encode error", func(t *testing.T) {
		img := newImageStub(100, 100)
		resizedImg := newImageStub(1000, 1000)
		in := new(bytes.Buffer)
		expectedErr := errors.New("test error")

		processor := &mockresizer.ImageProcessor{}
		processor.
			On("Decode", in).
			Once().
			Return(img, nil)

		processor.
			On("Resize", img, 1000, 1000).
			Once().
			Return(resizedImg)

		processor.
			On("Encode", resizedImg, anyWriter).
			Once().
			Return(expectedErr)

		resizer := NewImageResizer().WithProcessor(processor)
		_, err := resizer.Resize(in, 1000, 1000)
		require.Error(t, err)
		require.ErrorIs(t, err, expectedErr)
	})
}
