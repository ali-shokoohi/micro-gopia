package handlers

import (
	"net/http"

	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	HomeService *services.HomeService
}

func (h *HomeHandler) Root(c *gin.Context) {
	message := h.HomeService.Root()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": message,
	})
}

func (h *HomeHandler) GetNewHandler() *HomeHandler {
	return h
}
