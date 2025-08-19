package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func JadwalPresensiRoutes(r *gin.Engine) {
	jadwalPresensi := r.Group("/jadwalPresensi", middleware.AuthMiddleware())
	{
		jadwalPresensi.GET("/", controllers.GetJadwalPresensi)
	}
}
