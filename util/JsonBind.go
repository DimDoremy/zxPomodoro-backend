package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonBind(context *gin.Context, function func(), bindStruct interface{}) {
	var err error
	// 将结构体绑定json
	err = context.ShouldBind(&bindStruct)
	if err != nil {
		fmt.Println(err)
		// 如果绑定失败，返回错误信息json
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		function()
	}
}
