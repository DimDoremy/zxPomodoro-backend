package util

import (
	"bytes"
	"encoding/json"
	"gin_spring/dao"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

//存储session结构体
type Session struct {
	Openid     string   `json:"openid"`
	Recruit    []string `json:"recruit"`
}

type Room struct {
	*Session
	StartTime time.Time	`json:"start_time"`
	EndTime time.Time `json:"end_time"`
}

//设置session进入redis
func SetSessionIntoRedis(openid []byte, json []byte, db int) {
	conn := dao.Pool.Get()
	defer conn.Close()
	//fmt.Println(openid, " ", json)
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
func GetSessionByRedis(openid []byte, db int) Session {
	conn := dao.Pool.Get()
	defer conn.Close()
	var session Session
	session = Session{}
	reply, _ := conn.Do("Select", db)
	log.Println(reply.(string))
	jsons, err := redis.String(conn.Do("GET", openid))
	if err != nil {
		log.Println("Get failed:", err.Error())
		return session
	}
	err = json.Unmarshal([]byte(jsons), &session)
	if err != nil {
		log.Println("Get failed:", err.Error())
		return session
	}
	return session
}
