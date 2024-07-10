package db

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"

	glog "github.com/jianlu8023/go-logger"
	"github.com/jianlu8023/go-logger/pkg/db/model"
)

type TableInfo struct {
	TableName   string
	TableEngine string
	TableRows   int
}

func TestNewDevelopDBLogger(t *testing.T) {
	fmt.Println("--- TestNewDevelopDBLogger")
	newLogger := glog.NewLogger(&glog.Config{
		DevelopMode: false,
		LogLevel:    "info",
		Mode:        []string{"stdout", "file", "date"},
	})
	newLogger.Info("db info log")
	logger := NewProductionDBLogger(Config{
		Logger:                    newLogger,
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormlogger.Info,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	})

	var (
		username = "root"
		password = "123456"
		host     = "192.168.58.110"
		port     = "3306"
		database = "basic"
	)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})

	if err != nil {
		fmt.Println("db connect failed", err)
		return
	}

	err = db.AutoMigrate(&model.Basic{})
	if err != nil {

		fmt.Println("db migrate failed", err)
		return
	}

	var version string
	db.Raw("select version()").Scan(&version)
	fmt.Println("mysql version is", version)
	var tableInfos []TableInfo
	// 查询数据库中的所有表信息
	db.Raw("SELECT TABLE_NAME, ENGINE, TABLE_ROWS FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = 'upchain_dev'").Scan(&tableInfos)
	for _, info := range tableInfos {
		fmt.Printf("表名: %s, 引擎: %s, 行数: %d\n", info.TableName, info.TableEngine, info.TableRows)
	}

	fmt.Println("----------------------------------")

	var basic model.Basic
	find := db.Model(&model.Basic{}).Where(model.Basic{
		Uid: 1,
	}).First(&basic)

	if find.Error != nil {

		fmt.Println("find failed", find.Error)
		return
	}

	fmt.Println("basic is", basic)

}
