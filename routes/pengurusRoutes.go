package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PengurusRoutes(r *gin.Engine) {
	Pengurus := r.Group("/pengurus", middleware.AuthMiddleware())
	{
		Pengurus.GET("/", controllers.GetAllPengurus)
	}
}
