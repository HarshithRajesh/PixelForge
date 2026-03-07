package processor

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
)

type ImageManagement interface {
	UploadImage(ctx context.Context, header *multipart.FileHeader) error
}

type imageManagement struct{}

func NewImageManagement() ImageManagement {
	return &imageManagement{}
}

func (i *imageManagement) UploadImage(ctx context.Context, header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return errors.New("failed to open the file")
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		return errors.New("file too large, reduce the size of the image and upload")
	}
	fmt.Print(header.Filename)
	fmt.Print("Image recieved")
	return nil
}
