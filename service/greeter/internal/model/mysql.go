package model

import (
	"fmt"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/conf"
	"go-micro.dev/v4/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	err  error
	db   *gorm.DB
	once sync.Once
)

// InitMySQL init a GORM DB definition
func InitMySQL(c *conf.MySQL) {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User,
			c.Pass,
			c.Host,
			c.Port,
			c.DbName)

		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
		if err != nil {
			log.Fatal("init db fail, %+v", err)
		}
	})
}

// GetDb return a
func GetDb() *gorm.DB {
	if db == nil {
		log.Fatal("db is not init")
	}
	return db
}
