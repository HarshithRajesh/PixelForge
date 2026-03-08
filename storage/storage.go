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
	Save(path string, file *multipart.FileHeader) error
}

type storageRepository struct {
	RootDir string
}

func NewStorageRepository(rootDir string) StorageRepository {
	return &storageRepository{RootDir: rootDir}
}

func (s *storageRepository) Save(path string, file *multipart.FileHeader) error {
	fullpath := filepath.Join(s.RootDir, path)
	if err := os.MkdirAll(filepath.Dir(fullpath), os.ModePerm); err != nil {
		return errors.New("failed to create the directory")
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
