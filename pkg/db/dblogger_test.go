package db

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	glog "github.com/jianlu8023/go-logger"
	"github.com/jianlu8023/go-logger/pkg/db/model"
)

func TestNewDBLogger(t *testing.T) {
	newLogger := glog.NewLogger(
		&glog.Config{
			DevelopMode: false,
			LogLevel:    "info",
		},
		glog.WithConsoleFormat(),
		glog.WithLumberjack(&glog.LumberjackConfig{
			FileName:  "./logs/lumberjack-db.log",
			Localtime: true,
		}),
		glog.WithRotateLog(&glog.RotateLogConfig{
			FileName:  "./logs/rotatelog-db.log",
			LocalTime: true,
		}),
	)
	logger := NewDBLogger(Config{
		Logger:                    newLogger,
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  Info,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      true,
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
	db.Raw("select version();").Scan(&version)
	fmt.Println("mysql version is", version)

	if err := db.Model(&model.Basic{}).FirstOrCreate(&model.Basic{
		Name: "test",
		Age:  18,
		Sex:  0,
	}).Error; err != nil {
		fmt.Println("create failed", err)
		return
	}

	var basic model.Basic
	find := db.Model(&model.Basic{}).Where(model.Basic{Uid: 1}).First(&basic)

	if find.Error != nil {
		fmt.Println("find failed", find.Error)
		return
	}
	fmt.Println("basic is", basic)
}
