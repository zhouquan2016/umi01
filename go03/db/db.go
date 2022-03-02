package db

import (
	"fmt"
	"go03/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	ds := conf.GetConfig().Datasource
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", *ds.Username, *ds.Password, *ds.Host, *ds.Port, *ds.Database, *ds.Query)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	namingStrategy := schema.NamingStrategy{SingularTable: true}
	gdb, err := gorm.Open(mysql.Open(dbDSN), &gorm.Config{Logger: newLogger, NamingStrategy: namingStrategy})
	if err != nil {
		panic(err)
	}
	db = gdb
}
