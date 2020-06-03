package util

import (
	"bytes"
	"encoding/json"
	"gin_spring/dao"
	"gin_spring/model"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"time"
)

// Personal结构体
type Personal struct {
	Openid  string   `json:"openid"`
	Recruit []string `json:"recruit"`
	Score   []int64  `json:"score"`
}

// Room结构体
type Room struct {
	Openid    string    `json:"openid"`
	Recruit   []string  `json:"recruit"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// 切片转换为Map
func SliceToMap(sKey, sValue []string) map[string]string {
	mapObj := map[string]string{}
	for sKeyIndex := range sKey {
		mapObj[sKey[sKeyIndex]] = sValue[sKeyIndex]
	}
	return mapObj
}

//反射结构体转map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// 写入redis
func SetIntoRedis(openid []byte, json []byte, db int) interface{} {
	conn := dao.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("Select", db)
	if err != nil {
		return model.MessageBind{Message: "Set failed:" + err.Error()}
	}
	_, err = conn.Do("Set", openid, json)
	if err != nil {
		return model.MessageBind{Message: "Set failed:" + err.Error()}
	}
	var str []byte
	str, err = redis.Bytes(conn.Do("GET", openid))
	if err != nil {
		return model.MessageBind{Message: "Set failed:" + err.Error()}
	}
	if !bytes.Equal(str, json) {
		return model.MessageBind{Message: "Date failed:"}
	}
	return nil
}

//获得redis存储的session
func GetByRedis(openid []byte, db int) interface{} {
	conn := dao.Pool.Get()
	defer conn.Close()
	_, e := conn.Do("Select", db)
	if e != nil {
		return model.MessageBind{Message: "Get failed:" + e.Error()}
	}
	jsons, err := redis.String(conn.Do("GET", openid))
	if err != nil {
		return model.MessageBind{Message: "Get failed:" + err.Error()}
	}
	if db == 0 { // 如果是Personal
		var personal Personal
		personal = Personal{}
		err = json.Unmarshal([]byte(jsons), &personal)
		if err != nil {
			return model.MessageBind{Message: "Get failed:" + err.Error()}
		}
		return personal
	} else if db == 1 { // 如果是Room
		var room Room
		room = Room{}
		err = json.Unmarshal([]byte(jsons), &room)
		if err != nil {
			return model.MessageBind{Message: "Get failed:" + err.Error()}
		}
		return room
	}
	return nil
}

//删除redis的指定键值对
func DeleteFromRedis(openid []byte, db int) error {
	conn := dao.Pool.Get()
	defer conn.Close()
	_, _ = conn.Do("Select", db)
	_, err := redis.String(conn.Do("GET", openid))
	if err != nil {
		return err
	}
	_, err = conn.Do("Del", openid)
	if err != nil {
		return err
	}
	return nil
}
