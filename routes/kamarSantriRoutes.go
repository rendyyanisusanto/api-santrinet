package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func KamarSantriRoutes(r *gin.Engine) {
	KamarSantri := r.Group("/KamarSantri", middleware.AuthMiddleware())
	{
		KamarSantri.GET("/santri_id/:id", controllers.GetKamarSantriByID)
	}
}
