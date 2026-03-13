package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // Registers GIF decoder
	_ "image/jpeg" // Registers JPEG decoder
	_ "image/png"  // Registers PNG decoder
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
	"github.com/HarshithRajesh/PixelForge/storage"
	"github.com/google/uuid"
)

type ImageManagement interface {
	UploadImage(ctx context.Context, header *multipart.FileHeader, userID string) error
	ListImages(userID uint) ([]*models.Image, error)
	Transform(imageID string, userID string, req *models.TransformRequest) error
	// GetImageDimensions(data []byte) (int, int, error)
}

type imageManagement struct {
	repo        repository.UserRepository
	storageRepo storage.StorageRepository
	processor   ImageTransformation
}

func NewImageManagement(userRepo repository.UserRepository, store storage.StorageRepository) ImageManagement {
	return &imageManagement{
		repo:        userRepo,
		storageRepo: store,
	}
}

func (i *imageManagement) UploadImage(ctx context.Context, header *multipart.FileHeader, userID string) error {
	file, err := header.Open()
	if err != nil {
		return errors.New("failed to open the file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return errors.New("Failed to read the contentType")
	}

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !allowedTypes[contentType] {
		return errors.New("Wrong image type")
	}

	file.Seek(0, io.SeekStart)

	ext := filepath.Ext(header.Filename)
	newID := uuid.New().String()

	if header.Size > 5*1024*1024 {
		return errors.New("file too large, reduce the size of the image and upload")
	}

	storagePath := fmt.Sprintf("%s%s", newID, ext)

	fmt.Print(header.Filename)
	fmt.Print("Image recieved")
	err = i.storageRepo.Save(storagePath, userID, header)
	if err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return errors.New("Failed to decode the image config")
	}
	newuserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return errors.New("failed to convert userid from string to int")
	}
	imgMetadata := &models.Image{
		UserID:         uint(newuserID),
		StoredFilename: header.Filename,
		Path:           storagePath,
		Size:           uint64(header.Size), // Convert int64 to uint64
		MimeType:       header.Header.Get("Content-Type"),
		Width:          imgConfig.Width,
		Height:         imgConfig.Height,
	}
	err = i.repo.SaveImageDB(imgMetadata)
	if err != nil {
		return errors.New("Failed to save the image in Database")
	}
	return nil
}

func (i *imageManagement) ListImages(userID uint) ([]*models.Image, error) {
	return i.repo.GetAllImageData(userID)
}

func (i *imageManagement) Transform(imageID string, userID string, req *models.TransformRequest) error {
	var image models.Image
	image, err := i.repo.GetImage(userID, imageID)
	if err != nil {
		return err
	}
	data, err := i.storageRepo.Read(image.Path)
	if err != nil {
		return err
	}
	dataTransformed, err := i.processor.Process(req, data)
	if err != nil {
		return nil, err
	}
	err = i.storageRepo.SaveTransformedImage(image.Path, dataTransformed)
	if err != nil {
		return err
	}
	newuserID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return errors.New("failed to convert userid from string to int")
	}
	w, h, err := i.GetImageDimensions(dataTransformed)
	imgMetadata := &models.Image{
		UserID:         uint(newuserID),
		StoredFilename: image.StoredFilename,
		Path:           image.Path,
		Size:           uint64(len(dataTransformed)), // Convert int64 to uint64
		MimeType:       image.MimeType,
		Width:          w,
		Height:         h,
	}

	err = i.repo.SaveImageDB(imgMetadata)
	if err != nil {
		return err
	}
	return nil
}

func (i *imageManagement) GetImageDimensions(data []byte) (int, int, error) {
	// DecodeConfig reads only the image header (fast)
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}
