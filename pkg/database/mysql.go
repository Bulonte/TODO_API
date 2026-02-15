package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB全局数据库实例
var DB *gorm.DB

func InitMysql(host, port, username, password, dbName string) error {
	//创建data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败：: %v", err)
	}

	//创建连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("创建数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)           //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          //最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //连接最大存活时间

	log.Println("数据库连接成功")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseMysql() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			sqlDB.Close()
			log.Println("数据库连接关闭成功")
		}
	}
}
