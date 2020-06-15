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

// insertUserPiece doc
// @Summary 注册单个用户，新建背包
// @Description
// @Tags piece-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserPiece true "单个用户背包信息"
// @Success 200 {object} model.MessageBind "注册成功"
// @Failure 400 {object} model.MessageBind "注册失败"
// @Router /insertuserpiece [post]
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
				context.JSON(http.StatusOK, model.MessageBind{
					Message: "register success",
				})
			} else {
				// 如果记录存在，返回错误信息
				context.JSON(http.StatusBadRequest, model.MessageBind{
					Message: "piece exist!",
				})
			}
		}, &userPiece)
	})
}

// pieceByOpenid doc
// @Summary 获取单个用户背包信息
// @Description
// @Tags piece-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserPiece true "单个用户信息"
// @Success 200 {object} model.UserPiece "单个用户信息的json数据"
// @Router /piecebyopenid [post]
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

// updateUserPiece doc
// @Summary 更新单个用户背包信息
// @Description
// @Tags piece-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserPiece true "更新的用户背包信息"
// @Success 200 {object} model.MessageBind "更新成功"
// @Failure 400 {object} model.MessageBind "更新失败"
// @Router /updateuserpiece [post]
func updateUserPiece(group *gin.RouterGroup) {
	// 更新单个用户背包信息
	group.POST("/updateuserpiece", func(context *gin.Context) {
		// 定义结构体用于存值
		var userPiece model.UserPiece
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			updateStruct := userPiece
			dao.DB.Model(&updateStruct).Omit("openid").Where("openid = ?", updateStruct.Openid).Updates(&updateStruct)
			context.JSON(http.StatusOK, model.MessageBind{
				Message: updateStruct.Openid + "update success!",
			})
		}, &userPiece)
	})
}

// deleteUserPiece doc
// @Summary 删除背包数据库指定条目
// @Description
// @Tags piece-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserPiece true "单个用户信息"
// @Success 200 {object} model.MessageBind "删除成功"
// @Failure 400 {object} model.MessageBind "删除失败"
// @Router /deleteuserpiece [post]
func deleteUserPiece(group *gin.RouterGroup) {
	// 删除背包数据库指定条目
	group.POST("/deleteuserpiece", func(context *gin.Context) {
		// 定义结构体用于存值
		var userPiece model.UserPiece
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			deleteStruct := userPiece
			dao.DB.Where("openid = ?", deleteStruct.Openid).Delete(&deleteStruct)
			if deleteStruct.Openid != "" {
				context.JSON(http.StatusOK, model.MessageBind{
					Message: deleteStruct.Openid + "delete success!",
				})
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "update failed!"})
			}
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
