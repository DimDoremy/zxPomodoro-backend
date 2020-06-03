package routers

import (
	"encoding/json"
	"gin_spring/dao"
	"gin_spring/model"
	"gin_spring/util"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"reflect"
	"time"
)

const (
	PersonalRecruitDb = iota
	RoomRecruitDb
)

// 客户端请求使用的结构体
type requests struct {
	Openid  string `json:"openid"`
	Makerid string `json:"makerid"`
	Time    int64  `json:"time"`
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
		if dbname == "room" {
			_, _ = conn.Do("Select", RoomRecruitDb)
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
			scanMap := util.SliceToMap(scanKeys, scanValues)
			context.JSON(http.StatusOK, scanMap)
		} else if dbname == "personal" {
			_, _ = conn.Do("Select", PersonalRecruitDb)
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
			scanMap := util.SliceToMap(scanKeys, scanValues)
			context.JSON(http.StatusOK, scanMap)
		}
	})
}

// 返回所有的redis缓存
func returnRedisRecruit(group *gin.RouterGroup) {
	group.POST("/return_recruit", func(context *gin.Context) {
		iter := 0
		var scanKeys []string
		var scanValues []string
		conn := dao.Pool.Get()
		defer conn.Close()
		var rp model.RequestPersonal
		util.JsonBind(context, func() {
			// 扫描Personal
			_, _ = conn.Do("Select", PersonalRecruitDb)
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
			//personalKeys := scanKeys     // 个人openid
			personalValues := scanValues // 个人参加的全部房间号
			//fmt.Println(personalKeys)
			//fmt.Println(personalValues)
			// 闭包解析key作为personal结构体
			structFunc := func(valueArrays []string) []util.Personal {
				var up util.Personal
				var upSlice []util.Personal
				for _, v := range valueArrays {
					err := json.Unmarshal([]byte(v), &up)
					if err != nil {
						return []util.Personal{}
					}
					upSlice = append(upSlice, up)
				}
				return upSlice
			}(personalValues)
			var up util.Personal
			var toMapSlice []map[string]interface{}
			for _, v := range personalValues {
				_ = json.Unmarshal([]byte(v), &up)
				toMap := util.StructToMap(up)
				toMapSlice = append(toMapSlice, toMap)
			}
			// 扫描Room
			_, _ = conn.Do("Select", RoomRecruitDb)
			// 循环个人招募信息结构体获取对应的加入房间，并根据房间号查找房间时间
			var ur util.Room
			var startTimeSlice []time.Time
			var endTimeSlice []time.Time
			//fmt.Println(structFunc)
			for _, values := range structFunc {
				// 重置scan
				scanKeys = []string{}   // 重置scanKeys
				scanValues = []string{} // 重置scanValues
				for _, v := range values.Recruit {
					tempValue, _ := redis.String(conn.Do("Get", v))
					_ = json.Unmarshal([]byte(tempValue), &ur)
					startTimeSlice = append(startTimeSlice, ur.StartTime) // 开始时间统计切片
					endTimeSlice = append(endTimeSlice, ur.EndTime)       // 结束时间统计切片
					scanValues = append(scanValues, tempValue)            // 每个房间信息
				}
				scanKeys = append(scanKeys, values.Openid) // 每个房间号
			}
			returnStruct := struct {
				Recruit   []map[string]interface{} `json:"recruit"`
				StartTime []time.Time              `json:"start_time"`
				EndTime   []time.Time              `json:"end_time"`
			}{
				toMapSlice,
				startTimeSlice,
				endTimeSlice,
			}
			// TODO: 格式 openid-recruit-time，尽量recruit和time对应
			context.JSON(http.StatusOK, returnStruct)
		}, &rp)
	})
}

// 添加招募进入redis缓存
func addIntoRecruit(group *gin.RouterGroup) {
	group.POST("/add_into_recruit", func(context *gin.Context) {
		var rr requests
		isExist := false
		util.JsonBind(context, func() {
			byRedis := util.GetByRedis([]byte(rr.Openid), PersonalRecruitDb)
			_, types := byRedis.(util.Personal)
			if !types {
				messageBind := byRedis.(model.MessageBind)
				if messageBind.Message == "Get failed:redigo: nil returned" {
					byRedis = util.Personal{
						Openid:  "",
						Recruit: []string{},
						Score:   []int64{},
					}
				} else {
					context.JSON(http.StatusBadRequest, messageBind)
					return
				}
			}
			tempPersonal := util.Personal{
				Openid:  byRedis.(util.Personal).Openid,
				Recruit: byRedis.(util.Personal).Recruit,
				Score:   byRedis.(util.Personal).Score,
			}
			// 如果没有加入任何招募，新建招募缓存
			if tempPersonal.Openid == "" {
				tempPersonal.Openid = rr.Openid
			}
			// 判断加入的招募是否存在
			funcKeys := func() []string {
				iter := 0
				var scanKeys []string
				conn := dao.Pool.Get()
				defer conn.Close()
				_, _ = conn.Do("Select", RoomRecruitDb)
				for iter != 48 {
					scan, err := redis.Values(conn.Do("Scan", iter))
					if err != nil {
						return scanKeys
					} else {
						iters, _ := redis.Bytes(scan[0], nil)
						iter = int(iters[0])
						scanTempKeys, _ := redis.Strings(scan[1], nil)
						scanKeys = append(scanKeys, scanTempKeys...)
					}
				}
				return scanKeys
			}
			for index, value := range funcKeys() {
				// 如果扫描存在这个makerid，跳出扫描
				if rr.Makerid == value {
					break
				}
				// 如果全部扫描结束都没有跳出
				if index == len(funcKeys())-1 {
					context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Room is not exist."})
					return
				}
			}
			// 判断是否进入了多个重复的招募
			for _, value := range tempPersonal.Recruit {
				if value == rr.Makerid {
					isExist = true
					break
				}
			}
			// 如果没有重复进入
			if !isExist {
				tempPersonal.Recruit = append(tempPersonal.Recruit, rr.Makerid)
				tempPersonal.Score = append(tempPersonal.Score, 0)
				newRedis, _ := json.Marshal(tempPersonal)
				util.SetIntoRedis([]byte(rr.Openid), newRedis, PersonalRecruitDb)
				var sessionModel util.Personal
				_ = json.Unmarshal(newRedis, &sessionModel)
				// 添加用户至房间
				getByRedis := util.GetByRedis([]byte(rr.Makerid), RoomRecruitDb)
				recruitByRedis := getByRedis.(util.Room).Recruit
				recruitByRedis = append(recruitByRedis, rr.Openid)
				room := util.Room{
					Openid:    getByRedis.(util.Room).Openid,
					Recruit:   recruitByRedis,
					StartTime: getByRedis.(util.Room).StartTime,
					EndTime:   getByRedis.(util.Room).EndTime,
				}
				marshal, _ := json.Marshal(room)
				util.SetIntoRedis([]byte(rr.Makerid), marshal, RoomRecruitDb)
				// 发送json包
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
			byRedis := util.GetByRedis([]byte(rr.Openid), PersonalRecruitDb)
			// 判断interface的变量类型
			_, types := byRedis.(util.Personal)
			if !types {
				messageBind := byRedis.(model.MessageBind)
				if messageBind.Message == "Get failed:redigo: nil returned" {
					byRedis = util.Personal{
						Openid:  "",
						Recruit: []string{},
						Score:   []int64{},
					}
				} else {
					context.JSON(http.StatusBadRequest, messageBind)
					return
				}
			}
			tempPersonal := util.Personal{
				Openid:  byRedis.(util.Personal).Openid,
				Recruit: byRedis.(util.Personal).Recruit,
				Score:   byRedis.(util.Personal).Score,
			}
			// 判断招募列表中是否存在makerid与退出的id相同的数据，并获取索引位置
			for index, value := range tempPersonal.Recruit {
				if value == rr.Makerid {
					number = index
					break
				}
			}
			// 如果索引位置存在，即数组序号大于等于0
			if number >= 0 {
				tempPersonal.Recruit = append(tempPersonal.Recruit[:number], tempPersonal.Recruit[number+1:]...)
				tempPersonal.Score = append(tempPersonal.Score[:number], tempPersonal.Score[number+1:]...)
				// 如果删除后为空，则将其销毁
				if len(tempPersonal.Recruit) <= 0 {
					err := util.DeleteFromRedis([]byte(rr.Openid), PersonalRecruitDb)
					if err != nil {
						context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete recruit failed"})
					}
					context.JSON(http.StatusOK, model.MessageBind{Message: "Delete success"})
					return
				}
				newRedis, _ := json.Marshal(tempPersonal)
				util.SetIntoRedis([]byte(rr.Openid), newRedis, PersonalRecruitDb)
				var sessionModel util.Personal
				_ = json.Unmarshal(newRedis, &sessionModel)
				context.JSON(http.StatusOK, sessionModel)
			} else {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Exit recruit failed"})
			}
		}, &rr)
	})
}

// 建立招募房间到redis缓存
func newRecruitRoom(group *gin.RouterGroup) {
	group.POST("/new_room", func(context *gin.Context) {
		nowTime := time.Now()
		var rp model.RequestPersonal
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
				util.SetIntoRedis([]byte(rp.Openid), marshal, RoomRecruitDb)
				// 添加用户至房间
				getByRedis := util.GetByRedis([]byte(rp.Makerid), PersonalRecruitDb)
				recruitByRedis := getByRedis.(util.Personal).Recruit
				scoreByRedis := getByRedis.(util.Personal).Score
				recruitByRedis = append(recruitByRedis, rp.Openid)
				scoreByRedis = append(scoreByRedis, 0)
				personal := util.Personal{
					Openid:  rp.Makerid,
					Recruit: recruitByRedis,
					Score:   scoreByRedis,
				}
				marshal, _ = json.Marshal(personal)
				util.SetIntoRedis([]byte(rp.Makerid), marshal, PersonalRecruitDb)
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
		var rp model.RequestPersonal
		util.JsonBind(context, func() {
			byRedis := util.GetByRedis([]byte(rp.Openid), RoomRecruitDb)
			if byRedis.(util.Room).EndTime.IsZero() {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Room not exist."})
			} else if byRedis.(util.Room).EndTime.Sub(nowTime).Milliseconds() <= 0 {
				context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Time Out."})
			} else {
				number := -1
				for _, v := range byRedis.(util.Room).Recruit {
					personalByRedis := util.GetByRedis([]byte(v), PersonalRecruitDb)
					// 判断interface的变量类型
					_, types := byRedis.(util.Personal)
					if !types {
						messageBind := byRedis.(model.MessageBind)
						if messageBind.Message == "Get failed:redigo: nil returned" {
							byRedis = util.Personal{
								Openid:  "",
								Recruit: []string{},
								Score:   []int64{},
							}
						} else {
							context.JSON(http.StatusBadRequest, messageBind)
							return
						}
					}
					tempPersonal := util.Personal{
						Openid:  personalByRedis.(util.Personal).Openid,
						Recruit: personalByRedis.(util.Personal).Recruit,
						Score:   personalByRedis.(util.Personal).Score,
					}
					// 判断招募列表中是否存在makerid与退出的id相同的数据，并获取索引位置
					for index, value := range tempPersonal.Recruit {
						if value == rp.Makerid {
							number = index
							break
						}
					}
					// 如果索引位置存在，即数组序号大于等于0
					if number >= 0 {
						tempPersonal.Recruit = append(tempPersonal.Recruit[:number], tempPersonal.Recruit[number+1:]...)
						tempPersonal.Score = append(tempPersonal.Score[:number], tempPersonal.Score[number+1:]...)
						// 如果删除后为空，则将其销毁
						if len(tempPersonal.Recruit) <= 0 {
							err := util.DeleteFromRedis([]byte(rp.Openid), PersonalRecruitDb)
							if err != nil {
								context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete recruit failed"})
							}
							context.JSON(http.StatusOK, model.MessageBind{Message: "Delete success"})
							return
						}
						newRedis, _ := json.Marshal(tempPersonal)
						util.SetIntoRedis([]byte(rp.Openid), newRedis, PersonalRecruitDb)
						var sessionModel util.Personal
						_ = json.Unmarshal(newRedis, &sessionModel)
					}
				}
				err := util.DeleteFromRedis([]byte(rp.Openid), RoomRecruitDb)
				if err != nil {
					context.JSON(http.StatusBadRequest, model.MessageBind{Message: "Delete room failed:Redis failed."})
				}
				context.JSON(http.StatusOK, model.MessageBind{Message: "Delete room success."})
			}
		}, &rp)
	})
}

// 接收客户端完成的积分
func acceptRecruit(group *gin.RouterGroup) {
	group.POST("/accept_recruit", func(context *gin.Context) {
		var rr requests
		util.JsonBind(context, func() {
			// 获取个人参加招募缓存
			byRedis := util.GetByRedis([]byte(rr.Openid), PersonalRecruitDb)
			// 判断是否获取到的是错误信息
			_, types := byRedis.(util.Personal)
			if !types {
				messageBind := byRedis.(model.MessageBind)
				context.JSON(http.StatusBadRequest, messageBind)
				return
			}
			personalValue := reflect.ValueOf(byRedis)
			recruit := personalValue.Field(1).Interface().([]string)
			score := personalValue.Field(2).Interface().([]int64)
			for i, v := range recruit {
				if v == rr.Makerid {
					score[i] += rr.Time
					break
				}
			}
			tempPersonal := util.Personal{
				Openid:  byRedis.(util.Personal).Openid,
				Recruit: recruit,
				Score:   score,
			}
			newRedis, _ := json.Marshal(tempPersonal)
			util.SetIntoRedis([]byte(rr.Openid), newRedis, PersonalRecruitDb)
			context.JSON(http.StatusOK, tempPersonal)
		}, &rr)
	})
}

// 结算请求
func completeRecruit(group *gin.RouterGroup) {
	group.POST("/complete_recruit", func(context *gin.Context) {
		var rr requests
		util.JsonBind(context, func() {
			playerScore := 0
			byRedis := util.GetByRedis([]byte(rr.Openid), PersonalRecruitDb)
			// 判断是否获取到的是错误信息
			_, types := byRedis.(util.Personal)
			if !types {
				messageBind := byRedis.(model.MessageBind)
				context.JSON(http.StatusBadRequest, messageBind)
				return
			}
			personalValue := reflect.ValueOf(byRedis)
			recruit := personalValue.Field(1).Interface().([]string)
			score := personalValue.Field(2).Interface().([]int64)
			for i, v := range recruit {
				if v == rr.Makerid {
					playerScore = int(score[i])
					break
				}
			}
			byRedis = util.GetByRedis([]byte(rr.Makerid), RoomRecruitDb)
			_, types = byRedis.(util.Room)
			if !types {
				messageBind := byRedis.(model.MessageBind)
				context.JSON(http.StatusBadRequest, messageBind)
				return
			}
			roomValue := reflect.ValueOf(byRedis)
			roomStart := roomValue.Field(2).Interface().(time.Time)
			roomEnd := roomValue.Field(3).Interface().(time.Time)
			returnStruct := struct {
				Openid      string        `json:"openid"`
				Makerid     string        `json:"makerid"`
				PlayerScore time.Duration `json:"player_score"`
				AllScore    time.Duration `json:"all_score"`
			}{
				rr.Openid,
				rr.Makerid,
				time.Duration(playerScore),
				roomEnd.Sub(roomStart),
			}
			context.JSON(http.StatusOK, returnStruct)
		}, &rr)
	})
}

func RecruitRouters(router *gin.Engine) {
	recruitGroup := router.Group("/api")
	{
		showRedisRecruit(recruitGroup)
		returnRedisRecruit(recruitGroup)
		addIntoRecruit(recruitGroup)
		removeFromRecruit(recruitGroup)
		newRecruitRoom(recruitGroup)
		deleteRecruitRoom(recruitGroup)
		acceptRecruit(recruitGroup)
		completeRecruit(recruitGroup)
	}
}
