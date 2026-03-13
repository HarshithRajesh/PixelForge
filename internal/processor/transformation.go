package processor

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/HarshithRajesh/PixelForge/internal/models"

	"github.com/anthonynsimon/bild/transform"
)

type TransformRequest struct {
	Operation string
	Params    map[string]int
}

type ImageTransformation interface {
	Process(req *models.TransformRequest, img image.Image, format string) ([]byte, error)
}

type imageTransformation struct{}

func NewImageTransformation() ImageTransformation { return &imageTransformation{} }

func (i *imageTransformation) Process(req *models.TransformRequest, img image.Image, format string) ([]byte, error) {
	var res image.Image
	log.Println("Before transformation")
	switch req.Operation {
	case "resize":
		w := req.Params["width"]
		h := req.Params["height"]
		if w <= 0 || h <= 0 {
			return nil, errors.New("width and height must be greater than 0")
		}
		res = i.resize(img, w, h)
	default:
		return nil, errors.New("invalid operation as input")
	}
	log.Println("transformed")
	buf := new(bytes.Buffer)
	var err error

	switch format {
	case "png":
		err = png.Encode(buf, res)
	case "gif":
		err = gif.Encode(buf, res, nil)
	default:
		err = jpeg.Encode(buf, res, nil)
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *imageTransformation) resize(img image.Image, w, h int) image.Image {
	return transform.Resize(img, w, h, transform.Lanczos)
}
