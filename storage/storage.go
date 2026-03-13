// Package storage is used to save the image files
package storage

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type StorageRepository interface {
	Save(path string, userID string, file *multipart.FileHeader) error
	Read(path string, userID string) (image.Image, string, error)
	SaveTransformedImage(userID string, path string, data []byte) error
}

type storageRepository struct {
	RootDir string
}

func NewStorageRepository(rootDir string) StorageRepository {
	return &storageRepository{RootDir: rootDir}
}

func (s *storageRepository) Save(path string, userID string, file *multipart.FileHeader) error {
	fullpath := filepath.Join(s.RootDir, userID, path)
	// if err := os.MkdirAll(filepath.Dir(fullpath), os.ModePerm); err != nil {
	if err := os.MkdirAll(filepath.Dir(fullpath), 0o755); err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return errors.New("failed to open source file")
	}
	defer src.Close()
	dst, err := os.Create(fullpath)
	if err != nil {
		return errors.New("failed to save the image in the folder")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return errors.New("failed to write file content")
	}

	return nil
}

func (s *storageRepository) Read(path string, userID string) (image.Image, string, error) {
	fullpath := filepath.Join(s.RootDir, userID, path)
	file, err := os.Open(fullpath)
	if err != nil {
		return nil, "", errors.New("file not found ")
	}
	defer file.Close()
	file.Seek(0, 0)
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("Format error in decoding %w", err)
	}
	return img, format, nil
}

func (s *storageRepository) SaveTransformedImage(userID string, path string, data []byte) error {
	fullpath := filepath.Join(s.RootDir, userID, path)
	dst, err := os.Create(fullpath)
	if err != nil {
		return errors.New("failed to save the transformed Image")
	}

	defer dst.Close()
	return nil
}
