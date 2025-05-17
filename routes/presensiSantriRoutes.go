package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PresensiSantriRoutes(r *gin.Engine) {
	RekamMedis := r.Group("/presensi_santri", middleware.AuthMiddleware())
	{
		RekamMedis.GET("/santri_id/:id", controllers.GetPresensiSantriBySantriID)
	}
}
