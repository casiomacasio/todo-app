package handler

import (
	"time"

	"github.com/casiomacasio/todo-app/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "github.com/casiomacasio/todo-app/backend/docs" 
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service     *service.Service
	redisClient *redis.Client
}

func NewHandler(service *service.Service, redisClient *redis.Client) *Handler {
	return &Handler{
		service:     service,
		redisClient: redisClient,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(GlobalRateLimitMiddleware(h.redisClient, 1000, time.Minute)) 
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	auth.Use(RateLimitIPMiddleware(h.redisClient, 5, time.Minute))

	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}

	api := router.Group("/api")
	api.Use(h.userIdentity)
	api.Use(RateLimitMiddleware(h.redisClient, 50, time.Minute)) 
	{
		lists := api.Group("/lists")
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

		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}
	router.Static("/app", "./frontend")

	return router
}