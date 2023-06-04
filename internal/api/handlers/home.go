package handlers

import (
	"net/http"

	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	HomeService *services.HomeService
}

// Root godoc
// @Summary      Show the server's status
// @Description  home page
// @Tags         root
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.HttpSuccess
// @Router       / [get]
func (h *HomeHandler) Root(c *gin.Context) {
	message := h.HomeService.Root()
	c.JSON(http.StatusOK, &dto.HttpSuccess{
		Status:  "success",
		Message: message,
	})
}

func (h *HomeHandler) GetNewHandler() *HomeHandler {
	return h
}
