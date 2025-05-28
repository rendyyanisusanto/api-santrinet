package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func TatibRoutes(r *gin.Engine) {
	TatibRoutes := r.Group("/tatib", middleware.AuthMiddleware())
	{
		TatibRoutes.GET("/", controllers.GetTatib)
	}
}
