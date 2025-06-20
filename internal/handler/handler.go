package handler

import (
	"github.com/casiomacasio/todo-app/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/casiomacasio/todo-app/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
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