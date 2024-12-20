package dblogger

import (
	"fmt"
	"sync"
	"testing"
	"time"

	glog "github.com/jianlu8023/go-logger"
	gormv1 "github.com/jinzhu/gorm"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"xorm.io/xorm"
)

var (
	newLogger = glog.NewLogger(
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
				LevelKey:         "level",
				TimeKey:          "time",
				NameKey:          "name",
				CallerKey:        "caller",
				FunctionKey:      "func",
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
)

func TestNewDBLogger(t *testing.T) {

	logger := NewDBLogger(
		Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  INFO,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			ShowSql:                   true,
		},
		// WithCustomLogger(newLogger),
	)
	var wg sync.WaitGroup
	wg.Add(3)

	go func(log *Logger) {
		defer wg.Done()
		db, err := xorm.NewEngine(
			"mysql",
			"root:123456@tcp(127.0.0.1:3306)/basic",
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
			host     = "127.0.0.1"
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
		var (
			username   = "postgres"
			password   = "123456"
			host       = "127.0.0.1"
			port       = "5432"
			database   = "basic"
			searchPath = "public"
		)
		pgDsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable search_path=%v",
			host, username, password, database, port, searchPath)

		pgDB, err := gorm.Open(postgres.Open(pgDsn), &gorm.Config{
			Logger: log,
		})

		if err != nil {
			fmt.Println("pgDB connect failed", err)
			return
		}

		err = pgDB.AutoMigrate(&Basic{})
		if err != nil {
			fmt.Println("pgDB migrate failed", err)
			return
		}

		var version string
		pgDB.Raw("select version();").Scan(&version)
		fmt.Println("postgres version is", version)

		if err := pgDB.Model(&Basic{}).FirstOrCreate(&Basic{
			Name: "test",
			Age:  18,
			Sex:  0,
		}).Error; err != nil {
			fmt.Println("create failed", err)
			return
		}

		var basic Basic
		find := pgDB.Model(&Basic{}).Where(Basic{Uid: 1}).First(&basic)

		if find.Error != nil {
			fmt.Println("find failed", find.Error)
			return
		}
		fmt.Println("basic is", basic)
	}(logger)
	wg.Wait()
	newLogger.Sugar().Infof("test end ...")
}

func TestGormV1DbLogger(t *testing.T) {

	var (
		username = "root"
		password = "123456"
		host     = "127.0.0.1"
		port     = "3306"
		database = "basic"
	)
	mysqlDsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gormv1.Open("mysql", mysqlDsn)
	if err != nil {
		fmt.Println("connect db failed ", err)
		panic(err)
	} else {
		db.DB().SetMaxIdleConns(50)
		db.DB().SetMaxOpenConns(50)
		db.DB().SetConnMaxLifetime(time.Minute)
		db.Set("gorm:association_autoupdate", false).
			Set("gorm:association_autocreate", false)
		db.SingularTable(true)
		db.LogMode(true)
		db.Debug() // 打印所有日志
		db.SetLogger(NewDBLogger(
			Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  INFO,
				Colorful:                  true,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      true,
				ShowSql:                   true,
			},
			// WithCustomLogger(newLogger),
		))
	}

	if err := db.Model(&Basic{}).FirstOrCreate(&Basic{
		Name: "test",
		Age:  18,
		Sex:  0,
	}).Error; err != nil {
		fmt.Println("create failed", err)
		return
	}

	var basic Basic
	find := db.Model(&Basic{}).Where(Basic{Uid: 1}).First(&basic)

	if find.Error != nil {
		fmt.Println("find failed", find.Error)
		return
	}
	fmt.Println("basic is", basic)
}
