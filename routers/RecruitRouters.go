package routers

import (
	"encoding/json"
	"fmt"
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"time"
)

const (
	Personal_RECRUIT_DB = iota
	Room_RECRUIT_DB
)

// 客户端请求使用的结构体
type requests struct {
	Openid  string `json:"openid"`
	Makerid string `json:"makerid"`
}

// 个人招募结构体
type requestPersonal struct {
	Openid  string `json:"openid"`
	Makerid string `json:"makerid"`
	Time    int64  `json:"time"`
}

// 公共招募结构体
type requestPublic struct {
	Openid string `json:"openid"`
	Makerid string `json:"makerid"`
	Score int64 `json:"score"`
	IsSpecial bool `json:"is_special"`
}

//切片转换为Map
func sliceToMap(sKey, sValue []string) map[string]string {
	mapObj := map[string]string{}
	for s_key_index := range sKey {
		mapObj[sKey[s_key_index]] = sValue[s_key_index]
	}
	return mapObj
}

// 展示所有的redis缓存
func showRedisRecruit(group *gin.RouterGroup) {
	group.GET("/show_recruit/:dbname", func(context *gin.Context) {
		iter := 0
		var scanKeys []string
		var scanValues []string
		dbname := context.Param("dbname")
		conn := dao.Pool.Get()
		defer conn.Close()
		if dbname == "public" {
			_, _ = conn.Do("Select", Room_RECRUIT_DB)
			for iter != 48 {
				scan, err := redis.Values(conn.Do("SCAN", iter))
				if err != nil {
					return
				} else {
					iters, _ := redis.Bytes(scan[0], nil)
					iter = int(iters[0])
					scanTempKeys, _ := redis.Strings(scan[1], nil)
					for _, value := range scanTempKeys {
						scanTempValue, _ := redis.String(conn.Do("Get", value))
						scanValues = append(scanValues, scanTempValue)
					}
					scanKeys = append(scanKeys, scanTempKeys...)
				}
			}
			scanMap := sliceToMap(scanKeys, scanValues)
			context.JSON(http.StatusOK, scanMap)
		} else if dbname == "personal" {
			_, _ = conn.Do("Select", Personal_RECRUIT_DB)
			for iter != 48 {
				scan, err := redis.Values(conn.Do("SCAN", iter))
				if err != nil {
					return
				} else {
					iters, _ := redis.Bytes(scan[0], nil)
					iter = int(iters[0])
					scanTempKeys, _ := redis.Strings(scan[1], nil)
					for _, value := range scanTempKeys {
						scanTempValue, _ := redis.String(conn.Do("Get", value))
						scanValues = append(scanValues, scanTempValue)
					}
					scanKeys = append(scanKeys, scanTempKeys...)
				}
				if iter == 48 {
					break
				}
			}
			scanMap := sliceToMap(scanKeys, scanValues)
			context.JSON(http.StatusOK, scanMap)
		}
	})
}

// 添加招募进入redis缓存
func addIntoRecruit(group *gin.RouterGroup) {
	group.POST("/add_into_recruit", func(context *gin.Context) {
		var rr requests
		isExist := false
		util.JsonBind(context, func() {
			byRedis := util.GetRoomByRedis([]byte(rr.Openid), Personal_RECRUIT_DB)
			// 如果没有加入任何招募，新建招募缓存
			if byRedis.Openid == "" {
				byRedis.Openid = rr.Openid
			}
			// 判断是否进入了多个重复的招募
			for _, value := range byRedis.Recruit {
				if value == rr.Makerid {
					isExist = true
					break
				}
			}
			// 如果没有重复进入
			if !isExist {
				byRedis.Recruit = append(byRedis.Recruit, rr.Makerid)
				newRedis, _ := json.Marshal(byRedis)
				util.SetRoomIntoRedis([]byte(rr.Openid), newRedis, Personal_RECRUIT_DB)
				var sessionModel util.Room
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Add recruit failed"})
			}
		}, &rr)
	})
}

// 从redis缓存退出招募
func removeFromRecruit(group *gin.RouterGroup) {
	group.POST("/remove_recruit", func(context *gin.Context) {
		var rr requests
		var number = -1
		util.JsonBind(context, func() {
			byRedis := util.GetRoomByRedis([]byte(rr.Openid), Personal_RECRUIT_DB)
			// 判断招募列表中是否存在makerid与退出的id相同的数据，并获取索引位置
			for index, value := range byRedis.Recruit {
				if value == rr.Makerid {
					number = index
					break
				}
			}
			// 如果索引位置存在，即数组序号大于等于0
			if number >= 0 {
				byRedis.Recruit = append(byRedis.Recruit[:number], byRedis.Recruit[number+1:]...)
				newRedis, _ := json.Marshal(byRedis)
				util.SetRoomIntoRedis([]byte(rr.Openid), newRedis, Personal_RECRUIT_DB)
				var sessionModel util.Room
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Add recruit failed"})
			}
		}, &rr)
	})
}

// 建立招募房间到redis缓存
func newRecruitRoom(group *gin.RouterGroup) {
	group.POST("/new_room", func(context *gin.Context) {
		nowTime := time.Now()
		var rp requestPersonal
		util.JsonBind(context, func() {
			// 如果请求得到的makerid和该用户自己的openid相同
			if rp.Makerid == rp.Openid {
				recruitSlice := []string{rp.Openid}
				// 构建Room
				makerModel := util.Room{
					Openid:    rp.Openid,                                         // 用户的openid
					Recruit:   recruitSlice,                                      // 将要加入的人员
					StartTime: nowTime,                                           // 开始时间
					EndTime:   nowTime.Add(time.Duration(rp.Time) * time.Second), // 结束时间，通过发来的time初始化
				}
				marshal, _ := json.Marshal(makerModel)
				util.SetRoomIntoRedis([]byte(rp.Openid), marshal, Personal_RECRUIT_DB)
				context.JSON(http.StatusOK, makerModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Create room failed"})
			}
		}, &rp)
	})
}

// 从redis缓存销毁招募房间
func deleteRecruitRoom(group *gin.RouterGroup) {
	group.POST("/delete_room", func(context *gin.Context) {
		nowTime := time.Now()
		var rp requestPersonal
		util.JsonBind(context, func() {
			byRedis := util.GetRoomByRedis([]byte(rp.Openid), Personal_RECRUIT_DB)
			if byRedis.EndTime.IsZero() {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Room not exist."})
			} else if byRedis.EndTime.Sub(nowTime).Milliseconds() <= 0 {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Time Out."})
			} else {
				err := util.DelRoomByRedis([]byte(rp.Openid), Personal_RECRUIT_DB)
				if err != nil {
					context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Redis failed."})
				}
				context.JSON(http.StatusOK, model.MessageBind{Message: "Delete room success."})
			}
		}, &rp)
	})
}

// 加入公共招募
func addIntoPublic(group *gin.RouterGroup) {
	group.POST("/add_into_public", func(context *gin.Context) {
		var rr requests
		isExist := false
		util.JsonBind(context, func() {
			var recruitPublic model.RecruitPublic
			dao.DB.Where("mission = ?", rr.Makerid).Find(&recruitPublic)
			fmt.Println(recruitPublic)
			if !recruitPublic.IsOpen {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Mission is not open."})
			}
			byRedis := util.GetRoomByRedis([]byte(rr.Openid), Room_RECRUIT_DB)
			// 如果没有加入任何招募，新建招募缓存
			if byRedis.Openid == "" {
				byRedis.Openid = rr.Openid
			}
			// 判断是否进入了多个重复的招募
			for _, value := range byRedis.Recruit {
				if value == rr.Makerid {
					isExist = true
					break
				}
			}
			// 如果没有重复进入
			if !isExist {
				byRedis.Recruit = append(byRedis.Recruit, rr.Makerid)
				newRedis, _ := json.Marshal(byRedis)
				util.SetRoomIntoRedis([]byte(rr.Openid), newRedis, Room_RECRUIT_DB)
				var sessionModel util.Room
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Add recruit failed"})
			}
		}, &rr)
	})
}

func RecruitRouters(router *gin.Engine) {
	recruitGroup := router.Group("/api")
	{
		showRedisRecruit(recruitGroup)
		addIntoRecruit(recruitGroup)
		removeFromRecruit(recruitGroup)
		newRecruitRoom(recruitGroup)
		deleteRecruitRoom(recruitGroup)
		addIntoPublic(recruitGroup)
	}
}
