package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSongList(c *gin.Context) {
	list, err := h.service.FetchSongs()
	if err != nil {
		log.Fatalf("Ошибка получения списка песен: %v", err)
	}

	c.JSON(http.StatusOK, list)
}
