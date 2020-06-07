package routers

import (
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func allUserData(group *gin.RouterGroup) {
	// 测试输出用户数据库全部信息
	group.GET("/alluserdata", func(context *gin.Context) {
		var userDatas []model.UserData
		dao.DB.Table("userData").Find(&userDatas)
		context.JSON(http.StatusOK, userDatas)
	})
}

func insertUserData(group *gin.RouterGroup) {
	// 注册单个用户
	group.POST("/insertuserdata", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			if !dao.DB.NewRecord(&userData) {
				// 如果记录不存在，创建新纪录
				dao.DB.Create(&userData)
				context.JSON(http.StatusOK, gin.H{
					"message": "register success",
				})
			} else {
				// 如果记录存在，返回错误信息
				context.JSON(http.StatusBadRequest, gin.H{
					"message": "user exist!",
				})
			}
		}, &userData)
	})
}

func queryByOpenid(group *gin.RouterGroup) {
	// 获取单个用户信息
	group.POST("/querybyopenid", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			dao.DB.Where("openid = ?", userData.Openid).Find(&userData)
			context.JSON(http.StatusOK, userData)
		}, &userData)
	})
}

func updateUserData(group *gin.RouterGroup) {
	// 更新单个用户信息
	group.POST("/updateuserdata", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			updateStruct := userData
			dao.DB.Model(&updateStruct).Omit("openid").Where("openid = ?", updateStruct.Openid).Updates(&updateStruct)
			context.JSON(http.StatusOK, gin.H{
				"message": updateStruct.Openid + "update success!",
			})
		}, &userData)
	})
}

func deleteUserData(group *gin.RouterGroup) {
	// 删除背包数据库指定条目
	group.POST("/deleteuserdata", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			deleteStruct := userData
			dao.DB.Where("openid = ?", deleteStruct.Openid).Delete(&deleteStruct)
			context.JSON(http.StatusOK, gin.H{
				"message": deleteStruct.Openid + "delete success!",
			})
		}, &userData)
	})
}

func UserRouter(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		allUserData(userGroup)
		insertUserData(userGroup)
		queryByOpenid(userGroup)
		updateUserData(userGroup)
		deleteUserData(userGroup)
	}
}
