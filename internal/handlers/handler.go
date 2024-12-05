package handlers

import (
	"github.com/codeboris/music-lib/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("api", h.GetSongList)

	api := router.Group("/api")
	{
		lists := api.Group("/songs")
		{
			lists.GET("/", h.GetSongList)
			// lists.POST("/", h.createSong)
			// lists.PUT("/:id", h.updateSong)
			// lists.DELETE("/:id", h.deleteSong)

			// items := lists.Group(":id/lyrics")
			// {
			// 	items.GET("/", h.getLyrics)
			// }
		}
	}

	return router
}
