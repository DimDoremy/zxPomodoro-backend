package routers

import (
	"fmt"
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAllHandBook(group *gin.RouterGroup) {
	// 测试输出背包数据库全部信息
	group.GET("/allhandbook", func(context *gin.Context) {
		var handBooks []model.HandBook
		dao.DB.Table("handBook").Find(&handBooks)
		context.JSON(http.StatusOK, handBooks)
	})
}

// insertHandBook doc
// @Summary 注册新用户新建图鉴信息
// @Tags book-routers
// @Accept json
// @Produce json
// @Param handBook body model.HandBook true "注册用户信息"
// @Success 200 {object} model.MessageBind "注册成功"
// @Failure 400 {object} model.MessageBind "注册失败"
// @Router /inserthandbook [post]
func insertHandBook(group *gin.RouterGroup) {
	// 注册单个用户，新建背包
	group.POST("/inserthandbook", func(context *gin.Context) {
		// 定义结构体用于存值
		var handBook model.HandBook
		util.JsonBind(context, func() {
			//fmt.Println(handBook)
			//fmt.Println(dao.DB.NewRecord(&handBook))
			// 如果绑定成功，执行内容
			if !dao.DB.NewRecord(&handBook) {
				// 如果记录不存在，创建新纪录
				dao.DB.Create(&handBook)
				context.JSON(http.StatusOK, model.MessageBind{Message: "register success"})
			} else {
				// 如果记录存在，返回错误信息
				context.JSON(http.StatusBadRequest, model.MessageBind{
					Message: "register failed",
				})
			}
		}, &handBook)
	})
}

// handBookByOpenid doc
// @Summary 获取单个用户图鉴列表信息
// @Tags book-routers
// @Accept json
// @Produce json
// @Param handBook body model.HandBook true "单个用户信息"
// @Success 200 {object} model.HandBook "单个用户图鉴信息的json数据"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /handbookbyopenid [post]
func handBookByOpenid(group *gin.RouterGroup) {
	// 获取单个用户背包信息
	group.POST("/handbookbyopenid", func(context *gin.Context) {
		// 定义结构体用于存值
		var handBook model.HandBook
		util.JsonBind(context, func() {
			var handBookTmp model.HandBook
			dao.DB.Where("openid = ?", handBook.Openid).Find(&handBookTmp)
			if handBookTmp.Openid != "" {
				context.JSON(http.StatusOK, handBookTmp)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{
					Message: "openid not exist!",
				})
			}
		}, &handBook)
	})
}

// updateHandBook doc
// @Summary 更新单个用户图鉴信息
// @Tags book-routers
// @Accept json
// @Produce json
// @Param handBook body model.HandBook true "更新的用户信息"
// @Success 200 {object} model.MessageBind "更新成功"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /updatehandbook [post]
func updateHandBook(group *gin.RouterGroup) {
	// 更新单个用户背包信息
	group.POST("/updatehandbook", func(context *gin.Context) {
		// 定义结构体用于存值
		var handBook model.HandBook
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			updateStruct := handBook
			dao.DB.Model(&updateStruct).Where("openid = ?", handBook.Openid).Updates(&updateStruct)
			fmt.Println(updateStruct)
			if updateStruct.Openid != "" {
				context.JSON(http.StatusOK, model.MessageBind{Message: updateStruct.Openid + "update success!"})
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "update failed!"})
			}
		}, &handBook)
	})
}

// deleteHandBook doc
// @Summary 删除图鉴数据库的指定条目
// @Tags book-routers
// @Accept json
// @Produce json
// @Param handBook body model.HandBook true "删除用户信息"
// @Success 200 {object} model.MessageBind "删除成功"
// @Failure 400 {object} model.MessageBind "400 bad request"
// @Router /deletehandbook [post]
func deleteHandBook(group *gin.RouterGroup) {
	// 删除背包数据库指定条目
	group.POST("/deletehandbook", func(context *gin.Context) {
		// 定义结构体用于存值
		var handBook model.HandBook
		util.JsonBind(context, func() {
			// 如果绑定成功，执行内容
			var deleteStruct model.HandBook
			dao.DB.Where("openid = ?", handBook.Openid).Delete(&deleteStruct)
			fmt.Println(deleteStruct)
			if deleteStruct.Openid != "" {
				context.JSON(http.StatusOK, model.MessageBind{Message: deleteStruct.Openid + "delete success!"})
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "delete failed!"})
			}
		}, &handBook)
	})
}

func BookRouter(router *gin.Engine) {
	bookGroup := router.Group("/api")
	{
		getAllHandBook(bookGroup)
		insertHandBook(bookGroup)
		handBookByOpenid(bookGroup)
		updateHandBook(bookGroup)
		deleteHandBook(bookGroup)
	}
}
