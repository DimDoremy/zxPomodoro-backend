package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 全局变量定义
var (
	DB *gorm.DB
)

func InitMySQL() error {
	var err error
	DB, err = gorm.Open("mysql", "User:passwdd@tcp(127.0.0.1:3306)/dbname?charset=utf8")
	if err != nil {
		// 暂时定义发生错误之间中断
		return err
	}
	DB.SingularTable(true)
	// 预设延迟关闭数据库连接
	//defer db.Close()
	return DB.DB().Ping()
}

func CloseMySQL() {
	DB.Close()
}
