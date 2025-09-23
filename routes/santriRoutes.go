package routes

import (
	"santrinet-api/controllers"
	"santrinet-api/middleware"

	"github.com/gin-gonic/gin"
)

func SantriRoutes(r *gin.Engine) {
	Santri := r.Group("/santri", middleware.AuthMiddleware())
	{
		Santri.GET("/id/:id", controllers.GetSantriByID)
		Santri.GET("/detail/:id", controllers.GetSantriDetailByID)
		Santri.GET("/", controllers.GetAllSantri)
		Santri.GET("/paginated", controllers.GetSantriWithPagination)
		Santri.GET("/filter", controllers.GetSantriFiltered)      // Get semua santri dengan filter
		Santri.POST("/", controllers.CreateSantri)                // Endpoint insert data santri
		Santri.PUT("/nonaktif/:id", controllers.SoftDeleteSantri) // Soft delete santri
		Santri.PUT("/:id", controllers.UpdateSantri)              // Update data santri
	}
}
