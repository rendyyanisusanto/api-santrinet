package routes

import (
    "santrinet-api/controllers"
    "santrinet-api/middleware"

    "github.com/gin-gonic/gin"
)

func AccountRoutes(r *gin.Engine) {
    auth := r.Group("/")
    auth.Use(middleware.AuthMiddleware())
    {
        auth.GET("/account", controllers.GetAccount)
        auth.PUT("/account", controllers.UpdateAccount)
        auth.PUT("/account/password", controllers.ChangePassword)
    }
}
