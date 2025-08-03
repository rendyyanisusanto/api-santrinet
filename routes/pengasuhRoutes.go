package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PengasuhRoutes(r *gin.Engine) {
	pengasuh := r.Group("/pengasuh", middleware.AuthMiddleware())
	{
		pengasuh.GET("/id/:id", controllers.GetPengasuhByID)
		pengasuh.GET("/", controllers.GetAllPengasuh)
	}
}
