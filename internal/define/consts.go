package define

import (
	"sync"
	"time"

	rotateloggers "github.com/lestrrat-go/file-rotatelogs"
)

const (
	Lumberjack         = "lumberjack"
	LumberjackTemplate = "lumberjack:?fileName=%v&maxSize=%v&maxAge=%v&maxBackups=%v&compress=%v&localtime=%v"
)

var lumberjackMu sync.RWMutex

// LumberjackDefaults holds the default configuration for lumberjack sink.
// All fields are safe for concurrent read/write via the getter methods.
type LumberjackDefaults struct {
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Localtime  bool
}

// lumberjackDefaults is the single source of truth, guarded by lumberjackMu.
var lumberjackDefaults = LumberjackDefaults{
	FileName:   "./logs/lumberjack.log",
	MaxSize:    5,
	MaxBackups: 7,
	MaxAge:     30,
	Compress:   false,
	Localtime:  false,
}

// GetLumberjackFileName returns the current lumberjack default file name.
func GetLumberjackFileName() string {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.FileName
}

// SetLumberjackFileName sets the lumberjack default file name.
func SetLumberjackFileName(v string) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.FileName = v
}

// GetLumberjackMaxSize returns the current lumberjack default max size (MB).
func GetLumberjackMaxSize() int {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.MaxSize
}

// SetLumberjackMaxSize sets the lumberjack default max size (MB).
func SetLumberjackMaxSize(v int) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.MaxSize = v
}

// GetLumberjackMaxBackups returns the current lumberjack default max backups.
func GetLumberjackMaxBackups() int {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.MaxBackups
}

// SetLumberjackMaxBackups sets the lumberjack default max backups.
func SetLumberjackMaxBackups(v int) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.MaxBackups = v
}

// GetLumberjackMaxAge returns the current lumberjack default max age (days).
func GetLumberjackMaxAge() int {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.MaxAge
}

// SetLumberjackMaxAge sets the lumberjack default max age (days).
func SetLumberjackMaxAge(v int) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.MaxAge = v
}

// GetLumberjackCompress returns the current lumberjack default compress flag.
func GetLumberjackCompress() bool {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.Compress
}

// SetLumberjackCompress sets the lumberjack default compress flag.
func SetLumberjackCompress(v bool) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.Compress = v
}

// GetLumberjackLocaltime returns the current lumberjack default localtime flag.
func GetLumberjackLocaltime() bool {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults.Localtime
}

// SetLumberjackLocaltime sets the lumberjack default localtime flag.
func SetLumberjackLocaltime(v bool) {
	lumberjackMu.Lock()
	defer lumberjackMu.Unlock()
	lumberjackDefaults.Localtime = v
}

// SnapshotLumberjackDefaults returns a copy of the current lumberjack defaults.
// Used by the zap sink factory so each zap.Open("lumberjack:...") call gets a
// clean snapshot before merging query parameters. This avoids races when multiple
// goroutines open sinks concurrently.
func SnapshotLumberjackDefaults() LumberjackDefaults {
	lumberjackMu.RLock()
	defer lumberjackMu.RUnlock()
	return lumberjackDefaults
}

const (
	RotateLogs         = "rotatelogs"
	RotateLogsTemplate = "rotatelogs:?fileName=%v&maxAge=%v&localtime=%v&rotationTime=%v"
)

var rotatelogsMu sync.RWMutex

// RotatelogsDefaults holds the default configuration for rotatelogs sink.
type RotatelogsDefaults struct {
	BaseName     string
	RfileName    string
	RotationTime time.Duration
	RmaxAge      time.Duration
	Rlocaltime   *time.Location
	Rclock       rotateloggers.Clock
}

// rotatelogsDefaults is the single source of truth, guarded by rotatelogsMu.
var rotatelogsDefaults = RotatelogsDefaults{
	BaseName:     "./logs/rotatelogs.log",
	RfileName:    "./logs/rotatelogs.%Y-%m-%d-%H.log",
	RotationTime: 3 * time.Hour,
	RmaxAge:      24 * time.Hour,
	Rlocaltime:   time.UTC,
	Rclock:       rotateloggers.UTC,
}

// GetRotatelogsBaseName returns the current rotatelogs default base name.
func GetRotatelogsBaseName() string {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.BaseName
}

// SetRotatelogsBaseName sets the rotatelogs default base name.
func SetRotatelogsBaseName(v string) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.BaseName = v
}

// GetRotatelogsRfileName returns the current rotatelogs default filename format.
func GetRotatelogsRfileName() string {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.RfileName
}

// SetRotatelogsRfileName sets the rotatelogs default filename format.
func SetRotatelogsRfileName(v string) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.RfileName = v
}

// GetRotatelogsRotationTime returns the current rotatelogs default rotation time.
func GetRotatelogsRotationTime() time.Duration {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.RotationTime
}

// SetRotatelogsRotationTime sets the rotatelogs default rotation time.
func SetRotatelogsRotationTime(v time.Duration) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.RotationTime = v
}

// GetRotatelogsRmaxAge returns the current rotatelogs default max age.
func GetRotatelogsRmaxAge() time.Duration {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.RmaxAge
}

// SetRotatelogsRmaxAge sets the rotatelogs default max age.
func SetRotatelogsRmaxAge(v time.Duration) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.RmaxAge = v
}

// GetRotatelogsRlocaltime returns the current rotatelogs default time location.
func GetRotatelogsRlocaltime() *time.Location {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.Rlocaltime
}

// SetRotatelogsRlocaltime sets the rotatelogs default time location.
func SetRotatelogsRlocaltime(v *time.Location) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.Rlocaltime = v
}

// GetRotatelogsRclock returns the current rotatelogs default clock.
func GetRotatelogsRclock() rotateloggers.Clock {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults.Rclock
}

// SetRotatelogsRclock sets the rotatelogs default clock.
func SetRotatelogsRclock(v rotateloggers.Clock) {
	rotatelogsMu.Lock()
	defer rotatelogsMu.Unlock()
	rotatelogsDefaults.Rclock = v
}

// SnapshotRotatelogsDefaults returns a copy of the current rotatelogs defaults.
func SnapshotRotatelogsDefaults() RotatelogsDefaults {
	rotatelogsMu.RLock()
	defer rotatelogsMu.RUnlock()
	return rotatelogsDefaults
}
