package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PelanggaranRoutes(r *gin.Engine) {
	Pelanggaran := r.Group("/pelanggaran", middleware.AuthMiddleware())
	{
		Pelanggaran.GET("/santri_id/:id", controllers.GetPelanggaranById)
		Pelanggaran.GET("/detail/:id", controllers.GetDetailPelanggaranById)
		Pelanggaran.GET("/total/:id", controllers.GetTotalPelanggaranBySantriID)
		Pelanggaran.GET("/total-hari-ini/:id", controllers.GetTotalPelanggaranHariIniBySantriID)
	}
}
