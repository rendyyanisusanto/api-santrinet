package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func SantriRoutes(r *gin.Engine) {
	Santri := r.Group("/santri", middleware.AuthMiddleware())
	{
		Santri.GET("/id/:id", controllers.GetSantriByID)
		Santri.GET("/detail/:id", controllers.GetSantriDetailByID)
		Santri.GET("/", controllers.GetAllSantri)
	}
}
