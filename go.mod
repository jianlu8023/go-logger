module github.com/jianlu8023/go-logger/v2

go 1.21

require (
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty v0.0.20
	github.com/sykesm/zap-logfmt v0.0.4
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

replace (
	github.com/lestrrat-go/file-rotatelogs => github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.1.15
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.24
	github.com/sykesm/zap-logfmt => github.com/sykesm/zap-logfmt v0.0.4
	go.uber.org/zap => go.uber.org/zap v1.28.0
	// github.com/ugorji/go => github.com/ugorji/go v1.2.6
	gopkg.in/natefinch/lumberjack.v2 => gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
