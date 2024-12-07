package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codeboris/music-lib/internal/models"
	"github.com/gin-gonic/gin"
)

type getListsResponse struct {
	Data []models.Song `json:"data"`
}

type getListVerses struct {
	Data []string `json:"data"`
}

func getIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.DefaultQuery(key, fmt.Sprintf("%d", defaultValue))
	intValue, err := strconv.Atoi(value)
	if err != nil || intValue < 1 {
		return defaultValue
	}
	return intValue
}

// @Summary Get All List
// @Tags list
// @Description get all list
// @ID get-all-list
// @Accept  json
// @Produce  json
// @Success 200 {object} getListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs [get]
func (h *Handler) GetSongList(c *gin.Context) {
	filter := models.SongFilter{
		Group:       c.DefaultQuery("group", ""),
		Song:        c.DefaultQuery("song", ""),
		ReleaseDate: c.DefaultQuery("release_date", ""),
		Text:        c.DefaultQuery("text", ""),
		Link:        c.DefaultQuery("link", ""),
		Page:        getIntParam(c, "page", 1),
		Limit:       getIntParam(c, "limit", 10),
	}

	songList, err := h.service.GetSongList(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListsResponse{
		Data: songList,
	})
}

// @Summary Create song
// @Tags list
// @Description create song list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body models.InputSong true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/songs [post]
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

func (h *Handler) GetLyrics(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid page param")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "2"))
	if err != nil || limit < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid limit param")
		return
	}

	lyrics, err := h.service.GetLyricsList(songID, page, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getListVerses{
		Data: lyrics,
	})
}
