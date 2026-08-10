module github.com/jianlu8023/go-logger/db-logger/v2

go 1.21

require (
	github.com/jianlu8023/go-logger/v2 v2.0.2
	github.com/jinzhu/gorm v1.9.16
	go.uber.org/zap v1.28.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.30.0
	xorm.io/xorm v1.3.6
)

replace (
	github.com/jianlu8023/go-logger/v2 => ..
	github.com/jinzhu/gorm => github.com/jinzhu/gorm v1.9.16
	go.uber.org/zap => go.uber.org/zap v1.28.0
	gorm.io/driver/mysql => gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres => gorm.io/driver/postgres v1.6.0
	gorm.io/gorm => gorm.io/gorm v1.31.2

)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	xorm.io/builder v0.3.13 // indirect
)
