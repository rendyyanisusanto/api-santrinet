package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func KamarRoutes(r *gin.Engine) {
	kamar := r.Group("/kamar", middleware.AuthMiddleware())
	{
		kamar.GET("/", controllers.GetKamar)
	}
}
