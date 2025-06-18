package handler

import (
	"github.com/casiomacasio/todo-app/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/log-in", h.logIn)
		auth.POST("/log-out", h.logOut)
	}
	api := router.Group("/api")
	{
		lists := api.Group("/lists", h.userIdentity)
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)
			items := lists.Group("/:id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItem)
			}
		}
		items := api.Group("/items", h.userIdentity) 
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)	
		}
	}
	return router
}