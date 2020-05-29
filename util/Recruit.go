package util

import (
	"bytes"
	"encoding/json"
	"gin_spring/dao"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

// Room结构体
type Room struct {
	Openid    string    `json:"openid"`
	Recruit   []string  `json:"recruit"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// 写入redis
func SetRoomIntoRedis(openid []byte, json []byte, db int) {
	conn := dao.Pool.Get()
	defer conn.Close()
	reply, err := conn.Do("Select", db)
	if err != nil {
		log.Println("Change DB failed:", err.Error())
		return
	}
	reply, err = conn.Do("Set", openid, json)
	log.Println(reply.(string))
	if err != nil {
		log.Println("Set failed:", err.Error())
		return
	}
	var str []byte
	str, err = redis.Bytes(conn.Do("GET", openid))
	if err != nil {
		log.Println("Get failed:", err.Error())
		return
	}
	if !bytes.Equal(str, json) {
		log.Println("Data failed")
		return
	}
}

//获得redis存储的session
func GetRoomByRedis(openid []byte, db int) Room {
	conn := dao.Pool.Get()
	defer conn.Close()
	var room Room
	room = Room{}
	_, _ = conn.Do("Select", db)
	//reply, _ := conn.Do("Select", db)
	//log.Println(reply.(string))
	jsons, err := redis.String(conn.Do("GET", openid))
	if err != nil {
		log.Println("Get failed:", err.Error())
		return room
	}
	err = json.Unmarshal([]byte(jsons), &room)
	if err != nil {
		log.Println("Get failed:", err.Error())
		return room
	}
	return room
}

//删除redis的指定键值对
func DelRoomByRedis(openid []byte, db int) error {
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
