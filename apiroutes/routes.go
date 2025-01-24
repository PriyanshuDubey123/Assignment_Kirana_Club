package apiroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/PriyanshuDubey123/Assignment_Kirana_Club/internals/service"
)

func StoreVisitServiceRoutes(router *gin.Engine) {
	router.POST("/api/submit", service.SubmitJobHandler())
	router.GET("/api/status", service.GetJobInfoHandler())
}
