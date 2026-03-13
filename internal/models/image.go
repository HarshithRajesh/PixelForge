package models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	UserID         uint   `gorm:"column:user_id"`
	StoredFilename string `gorm:"column:stored_filename"`
	Path           string `gorm:"column:path"`
	Size           uint64 `gorm:"column:size"`
	MimeType       string `gorm:"column:mime_type"`
	Width          int    `gorm:"column:width"`
	Height         int    `gorm:"column:height"`
}

type TransformRequest struct {
	Operation string
	Params    map[string]int
}

// controller
//    ↓
// service.TransformImage(imageID, request)
//    ↓
// repository.FindImage(imageID)
//    ↓
// storage.ReadFile()
//    ↓
// processor.Apply(operation)
//    ↓
// storage.SaveFile()
//    ↓
// repository.Save(newImage)
