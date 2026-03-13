// Package storage is used to save the image files
package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type StorageRepository interface {
	Save(path string, userID string, file *multipart.FileHeader) error
	Read(path string) ([]byte, error)
	SaveTransformedImage(path string, data []byte) error
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

func (s *storageRepository) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *storageRepository) SaveTransformedImage(path string, data []byte) error {
	dst, err := os.Create(path)
	if err != nil {
		return errors.New("failed to save the transformed Image")
	}
	defer dst.Close()
	return nil
}
