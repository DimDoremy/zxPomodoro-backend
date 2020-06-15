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

// insertUserData doc
// @Summary 注册单个用户
// @Description
// @Tags user-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserData true "注册用户信息"
// @Success 200 {object} model.MessageBind "注册成功"
// @Failure 400 {object} model.MessageBind "注册失败"
// @Router /insertuserdata [post]
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
				context.JSON(http.StatusOK, model.MessageBind{
					Message: "register success",
				})
			} else {
				// 如果记录存在，返回错误信息
				context.JSON(http.StatusBadRequest, model.MessageBind{
					Message: "user exist!",
				})
			}
		}, &userData)
	})
}

// queryByOpenid doc
// @Summary 获取单个用户信息
// @Description
// @Tags user-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserData true "单个用户信息"
// @Success 200 {object} model.UserData "单个用户信息的json数据"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /querybyopenid [post]
func queryByOpenid(group *gin.RouterGroup) {
	// 获取单个用户信息
	group.POST("/querybyopenid", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			var userDataTmp model.UserData
			dao.DB.Where("openid = ?", userData.Openid).Find(&userData)
			if userDataTmp.Openid != "" {
				context.JSON(http.StatusOK, userData)
			}
		}, &userData)
	})
}

// queryByOpenid doc
// @Summary 更新单个用户信息
// @Description
// @Tags user-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserData true "更新的用户信息"
// @Success 200 {object} model.MessageBind "更新成功"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /updateuserdata [post]
func updateUserData(group *gin.RouterGroup) {
	// 更新单个用户信息
	group.POST("/updateuserdata", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			updateStruct := userData
			dao.DB.Model(&updateStruct).Omit("openid").Where("openid = ?", updateStruct.Openid).Updates(&updateStruct)
			context.JSON(http.StatusOK, model.MessageBind{
				Message: updateStruct.Openid + "update success!",
			})
		}, &userData)
	})
}

// queryByOpenid doc
// @Summary 更新单个用户信息
// @Description
// @Tags user-routers
// @Accept json
// @Produce json
// @Param handBook body model.UserData true "更新的用户信息"
// @Success 200 {object} model.MessageBind "删除成功"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /deleteuserdata [post]
func deleteUserData(group *gin.RouterGroup) {
	// 删除背包数据库指定条目
	group.POST("/deleteuserdata", func(context *gin.Context) {
		// 定义结构体用于存值
		var userData model.UserData
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			deleteStruct := userData
			dao.DB.Where("openid = ?", deleteStruct.Openid).Delete(&deleteStruct)
			context.JSON(http.StatusOK, model.MessageBind{
				Message: deleteStruct.Openid + "delete success!",
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
