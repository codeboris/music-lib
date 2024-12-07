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

type createSongResponse struct {
	ID int `json:"id"`
}

func getIntParam(c *gin.Context, key string, defaultValue int) int {
	value := c.DefaultQuery(key, fmt.Sprintf("%d", defaultValue))
	intValue, err := strconv.Atoi(value)
	if err != nil || intValue < 1 {
		return defaultValue
	}
	return intValue
}

// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по всем полям.
// @Tags Песни
// @Accept json
// @Produce json
// @Param group query string false "Группа или исполнитель песни"
// @Param song query string false "Название песни"
// @Param release_date query string false "Дата выхода песни (в формате DD-MM-YYYY)"
// @Param text query string false "Фрагмент текста песни"
// @Param link query string false "Ссылка на песню"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Success 200 {object} getListsResponse "Список песен"
// @Failure 500 {object} errorResponse "Ошибка сервера"
// @Router /songs [get]
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

// @Summary Создать новую песню
// @Description Создает новую песню, заполняя необходимые поля из внешнего API по группе и названию песни.
// @Tags Песни
// @Accept json
// @Produce json
// @Param inputSong body models.InputSong true "Данные для создания песни"
// @Success 200 {object} createSongResponse "ID созданной песни"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Ошибка сервера"
// @Router /songs [post]
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

	c.JSON(http.StatusOK, createSongResponse{
		ID: songId,
	})
}

// @Summary Обновить информацию о песне
// @Description Обновляет данные о песне по ID. Требуется передать данные для обновления в формате JSON.
// @Tags Песни
// @Accept json
// @Produce json
// @Param id path int true "ID песни для обновления"
// @Param input body models.InputUpdateSong true "Данные для обновления песни"
// @Success 200 {object} statusResponse "Статус обновления"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Ошибка сервера"
// @Router /songs/{id} [put]
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

// @Summary Удалить песню по ID
// @Description Удаляет песню по заданному ID. Требуется передать ID песни в URL.
// @Tags Песни
// @Accept json
// @Produce json
// @Param id path int true "ID песни для удаления"
// @Success 200 {object} statusResponse "Статус удаления"
// @Failure 400 {object} errorResponse "Неверный запрос"
// @Failure 500 {object} errorResponse "Ошибка сервера"
// @Router /songs/{id} [delete]
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

// @Summary Получить текст песни с разбивкой на куплеты
// @Description Данный метод извлекает текст песни по ID, делит его на куплеты и возвращает их постранично.
// @Tags Песни
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы для пагинации (по умолчанию 1)"
// @Param limit query int false "Количество куплетов на странице (по умолчанию 2)"
// @Success 200 {object} getListVerses "Список куплетов"
// @Failure 400 {object} errorResponse "Неверные параметры запроса"
// @Failure 500 {object} errorResponse "Ошибка сервера"
// @Router /songs/{id}/lyrics [get]
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
