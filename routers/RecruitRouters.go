package routers

import (
	"encoding/json"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	Personal_RECRUIT_DB = iota
	Public_RECRUIT_DB
)

//type recruit struct {
//	Block model.RecruitBlock
//	List  model.RecruitList
//}
//
//type recruits struct {
//	Blocks []model.RecruitBlock
//	Lists  []model.RecruitList
//}
//
//type recruitJSON struct {
//	Makerid string `json:"makerid"`
//	*model.RecruitBlock
//	*model.RecruitList
//}
//
//// getAllRecruit doc
//// @Summary 获取全部用户图鉴列表
//// @Description 用于测试的Api，请不要在正式开发中使用
//// @Tags recruit-routers
//// @Produce json
//// @Success 200 {object} routers.recruits "获取到的全部数据库条目"
//// @Failure 404 {string} string "404 not found"
//// @Router /allrecruit [get]
//func getAllRecruit(group *gin.RouterGroup) {
//	// 测试输出数据库全部信息
//	group.GET("/allrecruit", func(context *gin.Context) {
//		recruits := recruits{
//			Blocks: []model.RecruitBlock{},
//			Lists:  []model.RecruitList{},
//		}
//		dao.DB.Table("recruitBlock").Find(&recruits.Blocks)
//		dao.DB.Table("recruitList").Find(&recruits.Lists)
//		context.JSON(http.StatusOK, recruits)
//	})
//}
//
//// insertRecruit doc
//// @Summary 注册新招募板信息
//// @Tags recruit-routers
//// @Accept json
//// @Produce json
//// @Param handBook body routers.recruitJSON true "注册新招募板信息"
//// @Success 200 {object} model.MessageBind "注册成功"
//// @Failure 400 {object} model.MessageBind "注册失败"
//// @Router /insertrecruit [post]
//func insertRecruit(group *gin.RouterGroup) {
//	// 注册单个用户，新建图鉴列表
//	group.POST("/insertrecruit", func(context *gin.Context) {
//		var recruitJson recruitJSON
//		util.JsonBind(context, func() {
//			// 定义结构体用于存值
//			recruit := recruit{
//				Block: model.RecruitBlock{
//					MakerId:  recruitJson.Makerid,
//					BossId:   recruitJson.BossId,
//					BossLife: recruitJson.BossLife,
//				},
//				List: model.RecruitList{
//					Makerid: recruitJson.Makerid,
//					Gamer1:  recruitJson.Gamer1,
//					Gamer2:  recruitJson.Gamer2,
//					Gamer3:  recruitJson.Gamer3,
//				},
//			}
//			// 如果绑定成功，执行内容
//			if !(dao.DB.NewRecord(&recruit.Block) || dao.DB.NewRecord(&recruit.List)) {
//				// 如果记录不存在，创建新纪录
//				dao.DB.Create(&recruit.Block)
//				dao.DB.Create(&recruit.List)
//				context.JSON(http.StatusOK, model.MessageBind{Message: "register success!"})
//			} else {
//				// 如果记录存在，返回错误信息
//				context.JSON(http.StatusOK, model.MessageBind{Message: "data exist!"})
//			}
//		}, &recruitJson)
//	})
//}
//
//// recruitByMakerid doc
//// @Summary 获取单个用户图鉴列表信息
//// @Tags recruit-routers
//// @Accept json
//// @Produce json
//// @Param handBook body routers.recruitJSON true "单个用户信息"
//// @Success 200 {object} routers.recruitJSON "单个用户信息的json数据"
//// @Failure 400 {object} model.MessageBind "400 bad request"
//// @Router /recruitbymakerid [post]
//func recruitByMakerid(group *gin.RouterGroup) {
//	// 获取单个用户背包信息
//	group.POST("/recruitbymakerid", func(context *gin.Context) {
//		var recruitJson recruitJSON
//		// 定义结构体用于存值
//		recruit := recruit{
//			Block: model.RecruitBlock{},
//			List:  model.RecruitList{},
//		}
//		util.JsonBind(context, func() {
//			dao.DB.Where("makerid = ?", recruitJson.Makerid).Find(&recruit.Block)
//			dao.DB.Where("makerid = ?", recruitJson.Makerid).Find(&recruit.List)
//			context.JSON(http.StatusOK, recruitJSON{
//				Makerid:      recruit.Block.MakerId,
//				RecruitBlock: &recruit.Block,
//				RecruitList:  &recruit.List,
//			})
//		}, &recruitJson)
//	})
//}
//
//func updateRecruit(group *gin.RouterGroup) {
//	// 更新单个用户背包信息
//	group.POST("/updaterecruit", func(context *gin.Context) {
//		var recruitJson recruitJSON
//		util.JsonBind(context, func() {
//			// 如果绑定成功，执行内容
//			updateStruct := recruit{
//				Block: model.RecruitBlock{
//					MakerId:  recruitJson.Makerid,
//					BossId:   recruitJson.BossId,
//					BossLife: recruitJson.BossLife,
//				},
//				List: model.RecruitList{
//					Makerid: recruitJson.Makerid,
//					Gamer1:  recruitJson.Gamer1,
//					Gamer2:  recruitJson.Gamer2,
//					Gamer3:  recruitJson.Gamer3,
//				},
//			}
//			dao.DB.Model(&updateStruct.Block).Omit("makerid").Where("makerid = ?", recruitJson.MakerId).Updates(&updateStruct.Block)
//			dao.DB.Model(&updateStruct.List).Omit("makerid").Where("makerid = ?", recruitJson.Makerid).Updates(&updateStruct.List)
//			context.JSON(http.StatusOK, gin.H{
//				"message": updateStruct.Block.MakerId + " update success!",
//			})
//		}, &recruitJson)
//	})
//}
//
//func deleteRecruit(group *gin.RouterGroup) {
//	// 删除背包数据库指定条目
//	group.POST("/deleterecruit", func(context *gin.Context) {
//		var recruitJson recruitJSON
//		util.JsonBind(context, func() {
//			// 如果绑定成功，执行内容
//			deleteStruct := recruit{
//				Block: model.RecruitBlock{
//					MakerId:  recruitJson.Makerid,
//					BossId:   recruitJson.BossId,
//					BossLife: recruitJson.BossLife,
//				},
//				List: model.RecruitList{
//					Makerid: recruitJson.Makerid,
//					Gamer1:  recruitJson.Gamer1,
//					Gamer2:  recruitJson.Gamer2,
//					Gamer3:  recruitJson.Gamer3,
//				},
//			}
//			if !(dao.DB.NewRecord(&deleteStruct.Block) && dao.DB.NewRecord(&deleteStruct.List)) {
//				dao.DB.Where("makerid = ?", deleteStruct.Block.MakerId).Delete(&deleteStruct.Block)
//				dao.DB.Where("makerid = ?", deleteStruct.List.Makerid).Delete(&deleteStruct.List)
//				context.JSON(http.StatusOK, gin.H{
//					"message": deleteStruct.Block.MakerId + " delete success!",
//				})
//			} else {
//				context.JSON(http.StatusOK, gin.H{
//					"message": "data not exist!",
//				})
//			}
//		}, &recruitJson)
//	})
//}

// 申请加入使用的结构体
type requests struct {
	Openid  string `json:"openid"`
	Makerid string `json:"makerid"`
}

type requestPersonal struct {
	*requests
	Time time.Time `json:"time"`
}

// 添加招募进入redis缓存
func addIntoRecruit(group *gin.RouterGroup) {
	group.POST("/addintorecruit", func(context *gin.Context) {
		var rr requests
		isExist := false
		util.JsonBind(context, func() {
			//fmt.Println(rr.Openid, " ", rr.Makerid)
			redis := util.GetSessionByRedis([]byte(rr.Openid), Personal_RECRUIT_DB)
			// 如果没有加入任何招募，新建招募缓存
			if redis.Openid == "" {
				redis.Openid = rr.Openid
			}
			for _, value := range redis.Recruit {
				if value == rr.Makerid {
					isExist = true
					break
				}
			}
			if !isExist {
				redis.Recruit = append(redis.Recruit, rr.Makerid)
				newRedis, _ := json.Marshal(redis)
				util.SetSessionIntoRedis([]byte(rr.Openid), newRedis, Personal_RECRUIT_DB)
				var sessionModel util.Session
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Add recruit failed"})
			}
		}, &rr)
	})
}

// 从redis缓存删除招募
func removeFromRecruit(group *gin.RouterGroup) {
	group.POST("/removerecruit", func(context *gin.Context) {
		var rr requests
		var number = -1
		util.JsonBind(context, func() {
			redis := util.GetSessionByRedis([]byte(rr.Openid), Personal_RECRUIT_DB)
			for index, value := range redis.Recruit {
				if value == rr.Makerid {
					number = index
					break
				}
			}
			if number >= 0 {
				redis.Recruit = append(redis.Recruit[:number], redis.Recruit[number+1:]...)
				newRedis, _ := json.Marshal(redis)
				util.SetSessionIntoRedis([]byte(rr.Openid), newRedis, Personal_RECRUIT_DB)
				var sessionModel util.Session
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Add recruit failed"})
			}
		}, &rr)
	})
}

////redis缓存,处理招募列表的信息
//func newRecruitRoom(group *gin.RouterGroup) {
//	group.POST("/newroom", func(context *gin.Context) {
//		var rp requestPersonal
//		util.JsonBind(context, func() {
//			if rr.Makerid == rr.Openid {
//				recruitSlice := []string{rp.Openid}
//				makerModel := util.Room{
//					Session: &util.Session{
//						Openid:  rp.Openid,
//						Recruit: recruitSlice,
//					},
//					Time: rp.Time,
//				}
//				marshal, _ := json.Marshal(makerModel)
//				util.SetSessionIntoRedis([]byte(rp.Openid), marshal, Public_RECRUIT_DB)
//				context.JSON(http.StatusOK, makerModel)
//			} else {
//				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Create room failed"})
//			}
//		}, &rp)
//	})
//}
//
////redis销毁
//func deleteRecruitRoom(group *gin.RouterGroup) {
//	group.POST("/deleteroom", func(context *gin.Context) {
//		var rr requests
//		util.JsonBind(context, func() {
//
//		}, &rr)
//	})
//}

func RecruitRouters(router *gin.Engine) {
	recruitGroup := router.Group("/api")
	{
		//getAllRecruit(recruitGroup)
		//insertRecruit(recruitGroup)
		//recruitByMakerid(recruitGroup)
		//updateRecruit(recruitGroup)
		//deleteRecruit(recruitGroup)
		addIntoRecruit(recruitGroup)
		removeFromRecruit(recruitGroup)
	}
}
