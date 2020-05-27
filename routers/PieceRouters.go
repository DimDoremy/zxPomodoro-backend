package routers

import (
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PieceRouter(router *gin.Engine) {
	pieceGroup := router.Group("/api")
	{
		// 测试输出背包数据库全部信息
		pieceGroup.GET("/alluserpiece", func(context *gin.Context) {
			var userPieces []model.UserPiece
			dao.DB.Table("userPiece").Find(&userPieces)
			context.JSON(http.StatusOK, userPieces)
		})

		// 注册单个用户，新建背包
		pieceGroup.POST("/insertuserpiece", func(context *gin.Context) {
			// 定义结构体用于存值
			var userPiece model.UserPiece
			util.JsonBind(context, func() {
				// 如果绑定成功，执行内容
				if !dao.DB.NewRecord(&userPiece) {
					// 如果记录不存在，创建新纪录
					dao.DB.Create(&userPiece)
					context.JSON(http.StatusOK, gin.H{
						"message": "register success",
					})
				} else {
					// 如果记录存在，返回错误信息
					context.JSON(http.StatusBadRequest, gin.H{
						"message": "piece exist!",
					})
				}
			}, &userPiece)
		})

		// 获取单个用户背包信息
		pieceGroup.POST("/piecebyopenid", func(context *gin.Context) {
			// 定义结构体用于存值
			var userPiece model.UserPiece
			util.JsonBind(context, func() {
				dao.DB.Where("openid = ?", userPiece.Openid).Find(&userPiece)
				context.JSON(http.StatusOK, userPiece)
			}, &userPiece)
		})

		// 更新单个用户背包信息
		pieceGroup.POST("/updateuserpiece", func(context *gin.Context) {
			// 定义结构体用于存值
			var userPiece model.UserPiece
			util.JsonBind(context, func() {
				// 如果绑定成功，执行内容
				updateStruct := userPiece
				dao.DB.Model(&updateStruct).Omit("openid").Where("openid = ?", updateStruct.Openid).Updates(&updateStruct)
				context.JSON(http.StatusOK, gin.H{
					"message": updateStruct.Openid + "update success!",
				})
			}, &userPiece)
		})

		// 删除背包数据库指定条目
		pieceGroup.POST("/deleteuserpiece", func(context *gin.Context) {
			// 定义结构体用于存值
			var userPiece model.UserPiece
			util.JsonBind(context, func() {
				// 如果绑定成功，执行内容
				deleteStruct := userPiece
				dao.DB.Where("openid = ?", deleteStruct.Openid).Delete(&deleteStruct)
				context.JSON(http.StatusOK, gin.H{
					"message": deleteStruct.Openid + "delete success!",
				})
			}, &userPiece)
		})
	}
}
