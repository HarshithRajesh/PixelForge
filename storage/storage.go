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

	//
	// fullpath := filepath.Join(s.RootDir, userID)
	// if err := os.MkdirAll(filepath.Dir(userDir), os.ModePerm); err != nil {
	// 	return errors.New("failed to create the directory")
	// }
	// dst, err := os.Create(fullpath)
	// if err != nil {
	// 	return errors.New("failed to save the image in the folder")
	// }
	// defer dst.Close()
	//
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return errors.New("failed to write file content")
	// }
	//
	// src, err := file.Open()
	// if err != nil {
	// 	return errors.New("failed to open source file")
	// }
	// defer src.Close()
	// dst, err := os.Create(fullpath)
	// if err != nil {
	// 	return errors.New("failed to save the image in the folder")
	// }
	// defer dst.Close()
	//
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return errors.New("failed to write file content")
	// }
	//
	// return nil
	// fullpath := filepath.Join(s.RootDir, userID, path)

	// 2. Get the directory portion of the full path and create it
	// This ensures the RootDir, the userID folder, and any nested folders in 'path' are created.
	// if err := os.MkdirAll(filepath.Dir(fullpath), 0o755); err != nil {
	// 	return fmt.Errorf("failed to create directories: %w", err)
	// }
	//
	// // 3. Open the uploaded file
	// src, err := file.Open()
	// if err != nil {
	// 	return fmt.Errorf("failed to open source file: %w", err)
	// }
	// defer src.Close()
	//
	// // 4. Create the destination file
	// dst, err := os.Create(fullpath)
	// if err != nil {
	// 	return fmt.Errorf("failed to create the destination file: %w", err)
	// }
	// defer dst.Close()
	//
	// // 5. Copy the contents
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return fmt.Errorf("failed to write file content: %w", err)
	// }
	//
	// return nil
}
