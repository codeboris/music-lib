package handlers

import (
	"net/http"
	"strconv"

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
		return
	}

	songs, err := h.service.GetSongList(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListsResponse{
		Data: songs,
	})
}

func (h *Handler) CreateSong(c *gin.Context) {
	var inputSong models.InputSong

	if err := c.ShouldBindJSON(&inputSong); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	songDetail, err := h.service.GetExternalData(inputSong.GroupName, inputSong.SongName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	song := h.service.PrepareSong(songDetail, inputSong)
	songId, err := h.service.CreateSong(song)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": songId,
	})
}

func (h *Handler) UpdateSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.InputUpdateSong
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateSong(songID, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})

}

func (h *Handler) DeleteSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	if err := h.service.DeleteSong(songID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
