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

	api := router.Group("/api")
	{
		lists := api.Group("/songs")
		{
			lists.GET("/", h.GetSongList)
			lists.POST("/", h.CreateSong)
			lists.PUT("/:id", h.UpdateSong)
			lists.DELETE("/:id", h.DeleteSong)

			// items := lists.Group(":id/lyrics")
			// {
			// 	items.GET("/", h.getLyrics)
			// }
		}
	}

	return router
}
