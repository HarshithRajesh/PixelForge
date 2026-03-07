package handler

import (
	"net/http"

	"github.com/HarshithRajesh/PixelForge/internal/processor"
	"github.com/gin-gonic/gin"
)

type ImageManagementHandler struct {
	imgService processor.ImageManagement
}

func NewImageManagementHandler(imgService processor.ImageManagement) *ImageManagementHandler {
	return &ImageManagementHandler{imgService: imgService}
}

func (h *ImageManagementHandler) ImageUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "getting file error"})
		return
	}

	err = h.imgService.UploadImage(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "file was not sent to processor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "getting file "})
	return
}
