package routes

import (
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/analyze", analyze)

	return router
}
