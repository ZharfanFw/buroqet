package handler

import (
	"net/http"

	"epod-service/internal/domain"
	"epod-service/internal/service"

	"github.com/gin-gonic/gin"
)

type EPODHandler struct {
	service *service.EPODService
}

func NewEPODHandler(
	service *service.EPODService,
) *EPODHandler {

	return &EPODHandler{
		service: service,
	}
}

func (h *EPODHandler) Upload(c *gin.Context) {

	var req domain.UploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.service.ProcessUpload(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}