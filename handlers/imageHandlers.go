package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type imageHandlers struct{}

func NewImageHandlers() *imageHandlers {
	return &imageHandlers{}
}

// HandleGetImageById godoc
// @Summary Get image by ID
// @Description Retrieves an image file by its ID
// @Tags images
// @Produce octet-stream
// @Param imageId path string true "Image ID"
// @Success 200 {file} file "Binary image file"
// @Failure 400 {object} string "Invalid image id"
// @Failure 500 {object} string "Internal server error"
// @Router /images/{imageId} [get]
func (h *imageHandlers) HandleGetImageById(c *gin.Context) {
	imageId := c.Param("imageId")
	if imageId == "" {
		c.JSON(http.StatusBadRequest, "Invalid image id")
		return
	}

	fileName := filepath.Base(imageId)
	byteFile, err := os.ReadFile(fmt.Sprintf("images/%s", imageId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", byteFile)
}