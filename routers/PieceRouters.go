package routers

import (
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func allUserPiece(group *gin.RouterGroup) {
	// 测试输出背包数据库全部信息
	group.GET("/alluserpiece", func(context *gin.Context) {
		var userPieces []model.UserPiece
		dao.DB.Table("userPiece").Find(&userPieces)
		context.JSON(http.StatusOK, userPieces)
	})
}

func insertUserPiece(group *gin.RouterGroup) {
	// 注册单个用户，新建背包
	group.POST("/insertuserpiece", func(context *gin.Context) {
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
}

func pieceByOpenid(group *gin.RouterGroup) {
	// 获取单个用户背包信息
	group.POST("/piecebyopenid", func(context *gin.Context) {
		// 定义结构体用于存值
		var userPiece model.UserPiece
		util.JsonBind(context, func() {
			dao.DB.Where("openid = ?", userPiece.Openid).Find(&userPiece)
			context.JSON(http.StatusOK, userPiece)
		}, &userPiece)
	})
}

func updateUserPiece(group *gin.RouterGroup) {
	// 更新单个用户背包信息
	group.POST("/updateuserpiece", func(context *gin.Context) {
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
}

func deleteUserPiece(group *gin.RouterGroup) {
	// 删除背包数据库指定条目
	group.POST("/deleteuserpiece", func(context *gin.Context) {
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
func PieceRouter(router *gin.Engine) {
	pieceGroup := router.Group("/api")
	{
		allUserPiece(pieceGroup)
		insertUserPiece(pieceGroup)
		pieceByOpenid(pieceGroup)
		updateUserPiece(pieceGroup)
		deleteUserPiece(pieceGroup)
	}
}
