/**
+ 入口函数
*/
package main

import (
	"gin_spring/dao"
	"gin_spring/routers"
	"github.com/gin-gonic/gin"
	"log"

	_ "gin_spring/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title 知夕后台Api接口
// @version 1.0-golang
// @description 用于知夕番茄钟程序的后台api文档
// @contact.name Doremy
// @contact.url https://dimdoremy.github.io/
// @host 127.0.0.1:8848
// @BasePath /api
func main() {
	// 开启控制台日志显示
	//gin.ForceConsoleColor()
	// 关闭控制台日志颜色，Windows控制台颜色显示有问题建议直接关闭
	gin.DisableConsoleColor()
	//设置release模式运行
	gin.SetMode(gin.ReleaseMode)

	// 新建gin框架
	r := gin.Default()

	// 注册swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 载入数据库
	err := dao.InitMySQL()
	if err != nil {
		log.Panic(err.Error())
	}
	// 关闭数据库
	defer dao.CloseMySQL()

	// 开启Redis连接池
	dao.InitRedisPool()
	// 关闭Redis连接池
	defer dao.Pool.Close()

	// 注册路由
	routers.UserRouter(r)
	routers.PieceRouter(r)
	routers.BookRouter(r)
	routers.RecruitRouters(r)
	routers.WechatRouter(r)
	routers.BackendRouters(r)

	// 启动框架并判断运行是否成功
	err = r.Run(":8848")
	if err != nil {
		log.Println(err.Error())
	}
}
