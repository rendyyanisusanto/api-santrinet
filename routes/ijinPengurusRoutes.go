package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func IjinPengurusRoutes(r *gin.Engine) {
	pengasuh := r.Group("/ijin_pengurus", middleware.AuthMiddleware())
	{
		pengasuh.GET("/keluar", controllers.GetPengurusSedangKeluar)
		pengasuh.GET("/", controllers.GetAllIjinPengurus)
	}
}
