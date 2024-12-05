package handlers

import (
	"net/http"

	"github.com/codeboris/music-lib/internal/models"
	"github.com/gin-gonic/gin"
)

type getListsResponse struct {
	Data []models.Song `json:"data"`
}

func (h *Handler) GetSongList(c *gin.Context) {
	var filter models.SongFilter
	if err := c.BindQuery(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	songs, err := h.service.GetSongList(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, getListsResponse{
		Data: songs,
	})
}
