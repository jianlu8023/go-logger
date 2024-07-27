package dblogger

import (
	"fmt"
	"sync"
	"testing"
	"time"

	glog "github.com/jianlu8023/go-logger"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xorm.io/xorm"
)

func TestNewDBLogger(t *testing.T) {
	newLogger := glog.NewLogger(
		&glog.Config{
			DevelopMode: false,
			LogLevel:    "info",
			Caller:      true,
			StackLevel:  "error",
			ModuleName:  "[db]",
		},
		glog.WithConsoleFormat(),
		glog.WithConsoleConfig(
			zapcore.EncoderConfig{
				MessageKey:       "msg",
				LevelKey:         "",
				TimeKey:          "",
				NameKey:          "",
				CallerKey:        "",
				FunctionKey:      "",
				StacktraceKey:    "stacktrace",
				SkipLineEnding:   false,
				LineEnding:       zapcore.DefaultLineEnding,
				EncodeLevel:      glog.CustomColorCapitalLevelEncoder,
				EncodeTime:       glog.CustomTimeEncoder,
				EncodeDuration:   zapcore.SecondsDurationEncoder,
				EncodeCaller:     zapcore.ShortCallerEncoder,
				EncodeName:       zapcore.FullNameEncoder,
				ConsoleSeparator: "\t",
			},
		),
	)
	logger := NewDBLogger(Config{
		Logger:                    newLogger,
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  INFO,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      true,
		ShowSql:                   true,
	})
	var wg sync.WaitGroup
	wg.Add(3)

	go func(log *Logger) {
		defer wg.Done()
		db, err := xorm.NewEngine(
			"mysql",
			"root:123456@tcp(192.168.209.128:3306)/basic",
		)
		db.SetLogger(log)
		if err != nil {
			fmt.Println("xorm connect failed", err)
			return
		}

		defer func() {
			if err := db.Close(); err != nil {
				fmt.Println("xorm close failed", err)
			}
		}()
		version, err := db.DBVersion()
		if err != nil {
			fmt.Println("db version failed", err)
			return
		}
		fmt.Println("xorm version is", version)
		err = db.Ping()

		if err != nil {
			fmt.Println("db ping failed", err)
			return
		}
		var b Basic
		get, err := db.Where("uid=?", 1).Get(&b)

		if err != nil || get == false {
			fmt.Println("get failed", err)
			return
		}
		fmt.Println("xorm get basic ", b)
	}(logger)

	go func(log *Logger) {
		defer wg.Done()
		var (
			username = "root"
			password = "123456"
			host     = "192.168.209.128"
			port     = "3306"
			database = "basic"
		)
		mysqlDsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)

		mysqlDB, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{
			Logger: log,
		})

		if err != nil {
			fmt.Println("mysqlDB connect failed", err)
			return
		}

		err = mysqlDB.AutoMigrate(&Basic{})
		if err != nil {

			fmt.Println("mysqlDB migrate failed", err)
			return
		}

		var version string
		mysqlDB.Raw("select version();").Scan(&version)
		fmt.Println("mysql version is", version)

		if err := mysqlDB.Model(&Basic{}).FirstOrCreate(&Basic{
			Name: "test",
			Age:  18,
			Sex:  0,
		}).Error; err != nil {
			fmt.Println("create failed", err)
			return
		}

		var basic Basic
		find := mysqlDB.Model(&Basic{}).Where(Basic{Uid: 1}).First(&basic)

		if find.Error != nil {
			fmt.Println("find failed", find.Error)
			return
		}
		fmt.Println("basic is", basic)
	}(logger)

	go func(log *Logger) {
		defer wg.Done()
		//var (
		//	username   = "postgres"
		//	password   = "123456"
		//	host       = "192.168.58.110"
		//	port       = "5432"
		//	database   = "basic"
		//	searchPath = "public"
		//)
		//pgDsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable search_path=%v",
		//	host, username, password, database, port, searchPath)
		//
		//pgDB, err := gorm.Open(postgres.Open(pgDsn), &gorm.Config{
		//	Logger: log,
		//})
		//
		//if err != nil {
		//	fmt.Println("pgDB connect failed", err)
		//	return
		//}
		//
		//err = pgDB.AutoMigrate(&Basic{})
		//if err != nil {
		//	fmt.Println("pgDB migrate failed", err)
		//	return
		//}
		//
		//var version string
		//pgDB.Raw("select version();").Scan(&version)
		//fmt.Println("postgres version is", version)
		//
		//if err := pgDB.Model(&Basic{}).FirstOrCreate(&Basic{
		//	Name: "test",
		//	Age:  18,
		//	Sex:  0,
		//}).Error; err != nil {
		//	fmt.Println("create failed", err)
		//	return
		//}
		//
		//var basic Basic
		//find := pgDB.Model(&Basic{}).Where(Basic{Uid: 1}).First(&basic)
		//
		//if find.Error != nil {
		//	fmt.Println("find failed", find.Error)
		//	return
		//}
		//fmt.Println("basic is", basic)
	}(logger)
	wg.Wait()
	newLogger.Sugar().Infof("test end ...")
}
