package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func PoinPelanggaran(r *gin.Engine) {
	poinPelanggaran := r.Group("/PoinPelanggaran", middleware.AuthMiddleware())
	{
		poinPelanggaran.GET("/:id", controllers.GetPoinPelanggaranbyId)
		poinPelanggaran.GET("/", controllers.GetAllPoinPelanggaran)
		// Tambah lagi POST, PUT, DELETE jika ada
	}
}
