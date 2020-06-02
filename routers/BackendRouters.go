package routers

import (
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func activePublic(group *gin.RouterGroup) {
	var public model.RecruitPublic
	group.POST("/active_public", func(context *gin.Context) {
		util.JsonBind(context, func() {
			context.JSON(http.StatusOK, gin.H{})
		}, &public)
	})
}

func BackendRouters(router *gin.Engine) {
	backendGroup := router.Group("/api/backend")
	{
		activePublic(backendGroup)
	}
}
