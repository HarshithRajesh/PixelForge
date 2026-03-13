package processor

import (
	"errors"

	"github.com/HarshithRajesh/PixelForge/internal/models"
)

type TransformRequest struct {
	Operation string
	Params    map[string]int
}

type ImageTransformation interface {
	Process(req *models.TransformRequest, data []byte) ([]byte, error)
}

type imageTransformation struct{}

func NewImageTransformation() ImageTransformation { return &imageTransformation{} }

func (i *imageTransformation) Process(req *models.TransformRequest, data []byte) ([]byte, error) {
	switch req.Operation {
	case "resize":
		return i.resize()
	case "crop":
		return i.crop()
	default:
		return nil, errors.New("invalid operation as input")
	}
}

func (i *imageTransformation) resize() ([]byte, error) { return nil, nil }
func (i *imageTransformation) crop() ([]byte, error)   { return nil, nil }
