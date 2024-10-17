package handlers

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"trackerApp/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/signIn", h.SignIn)
		api.POST("/signUp", h.SignUp)

		protected := api.Group("/protected", h.AuthMiddleware())
		{
			tasks := api.Group("/tasks")
			{
				tasks.GET("/", h.AllTasks)
				tasks.GET("/:id", h.TaskById)
				tasks.POST("/", h.PostTask)
				tasks.PUT("/:id", h.PutTask)
				tasks.DELETE("/:id", h.DeleteTask)
			}
			protected.POST("/logout")
		}
	}
	return router
}
