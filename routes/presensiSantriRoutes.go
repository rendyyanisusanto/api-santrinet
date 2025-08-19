package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PresensiSantriRoutes(r *gin.Engine) {
	Presensi := r.Group("/presensi_santri", middleware.AuthMiddleware())
	{
		Presensi.GET("/santri_id/:id", controllers.GetPresensiSantriBySantriID)
		Presensi.GET("/filter", controllers.GetPresensiByKamar)
	}
}
