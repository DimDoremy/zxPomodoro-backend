package dao

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//使用redis连接池
var Pool *redis.Pool

//初始化连接池
func InitRedisPool() {
	Pool = &redis.Pool{
		//创建连接池方法
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			_, err = conn.Do("AUTH", "密码")
			if err != nil {
				_ = conn.Close()
				return nil, err
			}
			return conn, err
		},
		//测试连接池方法
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         8,                 //最大空闲连接处，打开8个连接并保持
		MaxActive:       100000,            //最大激活连接数，0为自动按需分配
		IdleTimeout:     300 * time.Second, //最大空闲等待时间，超过后空闲连接将会关闭
		Wait:            true,              //当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
		MaxConnLifetime: 0,
	}
}
