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
		Presensi.POST("/update", controllers.UpdatePresensiSantri)
		
		// Routes untuk laporan
		Presensi.GET("/laporan/rekap-harian", controllers.GetLaporanRekapHarian)
		Presensi.GET("/laporan/bulanan-per-santri", controllers.GetLaporanBulananPerSantri)
		Presensi.GET("/laporan/kehadiran-per-waktu", controllers.GetLaporanKehadiranPerWaktu)
	}
}
