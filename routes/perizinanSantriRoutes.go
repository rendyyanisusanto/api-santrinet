package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PerizinanSantriRoutes(r *gin.Engine) {
	PengajuanPelanggaran := r.Group("/PerizinanSantri", middleware.AuthMiddleware())
	{
		PengajuanPelanggaran.GET("/", controllers.GetPerizinanSantri)
		PengajuanPelanggaran.POST("/", controllers.AddPerizinanSantri)
	}
}
