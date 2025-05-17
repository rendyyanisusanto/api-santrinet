package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func RekamMedisRoutes(r *gin.Engine) {
	RekamMedis := r.Group("/rekam_medis", middleware.AuthMiddleware())
	{
		RekamMedis.GET("/santri_id/:id", controllers.GetRekamMedisBySantriId)
		RekamMedis.GET("/id/:id", controllers.GetRekamMedisById)
	}
}
