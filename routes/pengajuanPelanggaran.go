package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PengajuanPelanggaranRoutes(r *gin.Engine) {
	PengajuanPelanggaran := r.Group("/Pengajuanpelanggaran", middleware.AuthMiddleware())
	{
		PengajuanPelanggaran.GET("/", controllers.GetPengajuanPelanggaran)
		PengajuanPelanggaran.POST("/", controllers.AddPengajuanPelanggaran)
	}
}
