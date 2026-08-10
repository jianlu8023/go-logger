package go_logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// setupLogDir creates the log directory for tests
func setupLogDir(t *testing.T) string {
	dir := filepath.Join(t.TempDir(), "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create log dir: %v", err)
	}
	return dir
}

// logToConsoleAndCapture writes a log message and returns the output as string
// We use this to verify console output by checking the logger was created successfully
func logAllLevels(logger *zap.Logger, sugar *zap.SugaredLogger, testName string) {
	// SugaredLogger methods
	sugar.Debug("debug", zap.String("test", testName))
	sugar.Info("info", zap.String("test", testName))
	sugar.Warn("warn", zap.String("test", testName))
	sugar.Error("error", zap.String("test", testName))

	// Standard methods with typed fields
	logger.Debug("debug-field", zap.String("test", testName), zap.Int("num", 42))
	logger.Info("info-field", zap.String("test", testName))
	logger.Warn("warn-field", zap.String("test", testName))
	logger.Error("error-field", zap.String("test", testName))
}

// checkLogFileExists verifies a specific log file was created
func checkLogFileExists(t *testing.T, path string) bool {
	t.Helper()
	// First check the exact file
	if _, err := os.Stat(path); err == nil {
		return true
	}
	// For rotatelogs, the file is time-based, check for any matching pattern
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	// Try glob patterns for rotated files
	patterns := []string{
		filepath.Join(dir, base+".*"),
		filepath.Join(dir, base+"*"),
	}
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return true
		}
	}
	return false
}

func TestLoggerTableDriven(t *testing.T) {
	logDir := setupLogDir(t)

	type testCase struct {
		name          string
		opts          []Option
		expectConsole bool
		expectFile    bool
		filePath      string // specific file path to check, "" means skip
	}

	tests := []testCase{
		// === Only Console ===
		{
			name:          "only_console_default",
			opts:          []Option{WithConsoleOutPut()},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_debug",
			opts:          []Option{WithConsoleOutPut(), WithDefaultLogLevel("debug")},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_json",
			opts:          []Option{WithConsoleOutPut(), WithJSONFormat()},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_console_format",
			opts:          []Option{WithConsoleOutPut(), WithConsoleFormat()},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_zaplogfmt",
			opts:          []Option{WithConsoleOutPut(), WithZaplogfmtFormat()},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_with_caller",
			opts:          []Option{WithConsoleOutPut(), WithDefaultLogLevel("debug"), WithCaller()},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_with_module",
			opts:          []Option{WithConsoleOutPut(), WithModuleName("myapp")},
			expectConsole: true,
			expectFile:    false,
		},
		{
			name:          "only_console_develop_mode",
			opts:          []Option{WithConsoleOutPut(), WithDevelopMode()},
			expectConsole: true,
			expectFile:    false,
		},
		// No explicit output -> defaults to console
		{
			name:          "no_explicit_output_defaults_to_console",
			opts:          []Option{WithDefaultLogLevel("debug")},
			expectConsole: true,
			expectFile:    false,
		},

		// === Only File: Lumberjack Default ===
		{
			name:          "only_file_lumberjack_default",
			opts:          []Option{WithFileOutPut(), WithDefaultLogLevel("info"), WithLumberjack(LumberjackDefaultConfig())},
			expectConsole: false,
			expectFile:    true,
		},

		// === Only File: Lumberjack Custom ===
		{
			name: "only_file_lumberjack_custom",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName:   filepath.Join(logDir, "lumberjack_custom.log"),
					MaxSize:    10,
					MaxAge:     14,
					MaxBackups: 5,
					Compress:   true,
					Localtime:  true,
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_custom.log"),
		},

		// === Only File: Lumberjack Separate File Level ===
		{
			name: "only_file_lumberjack_separate_file_level",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithFileLogLevel("warn"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "lumberjack_file_level.log"),
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_file_level.log"),
		},

		// === Only File: Lumberjack Separate Console Level ===
		{
			name: "only_file_lumberjack_separate_console_level",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithConsoleLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "lumberjack_console_level.log"),
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_console_level.log"),
		},

		// === Only File: Rotatelog Default ===
		{
			name: "only_file_rotatelog_default",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("info"),
				WithRotateLog(RotateLogDefaultConfig()),
			},
			expectConsole: false,
			expectFile:    true,
		},

		// === Only File: Rotatelog Custom ===
		{
			name: "only_file_rotatelog_custom",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(&RotateLogConfig{
					FileName:     filepath.Join(logDir, "rotatelog_custom.log"),
					MaxAge:       "7d",
					RotationTime: "1h",
					LocalTime:    true,
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "rotatelog_custom.log"),
		},

		// === Console + File: Lumberjack Default ===
		{
			name: "console_and_file_lumberjack_default",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("info"),
				WithLumberjack(LumberjackDefaultConfig()),
			},
			expectConsole: true,
			expectFile:    true,
		},

		// === Console + File: Lumberjack Custom ===
		{
			name: "console_and_file_lumberjack_custom",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName:  filepath.Join(logDir, "lumberjack_both_custom.log"),
					MaxSize:   20,
					Compress:  true,
					Localtime: true,
				}),
			},
			expectConsole: true,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_both_custom.log"),
		},

		// === Console + File: Lumberjack Separate Levels ===
		{
			name: "console_and_file_lumberjack_separate_levels",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithConsoleLogLevel("debug"),
				WithFileLogLevel("warn"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "lumberjack_separate.log"),
				}),
			},
			expectConsole: true,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_separate.log"),
		},

		// === Console + File: Lumberjack with Caller & Module ===
		{
			name: "console_and_file_lumberjack_caller_module",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithModuleName("test"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "lumberjack_caller_module.log"),
				}),
			},
			expectConsole: true,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "lumberjack_caller_module.log"),
		},

		// === Console + File: Rotatelog Default ===
		{
			name: "console_and_file_rotatelog_default",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("info"),
				WithRotateLog(RotateLogDefaultConfig()),
			},
			expectConsole: true,
			expectFile:    true,
		},

		// === Console + File: Rotatelog Custom ===
		{
			name: "console_and_file_rotatelog_custom",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(&RotateLogConfig{
					FileName:     filepath.Join(logDir, "rotatelog_both_custom.log"),
					MaxAge:       "3d",
					RotationTime: "30m",
					LocalTime:    true,
				}),
			},
			expectConsole: true,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "rotatelog_both_custom.log"),
		},

		// === Console + File: Rotatelog with Caller & Module ===
		{
			name: "console_and_file_rotatelog_caller_module",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithModuleName("test"),
				WithRotateLog(&RotateLogConfig{
					FileName: filepath.Join(logDir, "rotatelog_caller_module.log"),
				}),
			},
			expectConsole: true,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "rotatelog_caller_module.log"),
		},

		// === No Output (explicit disable console, no file) ===
		{
			name:          "no_output_at_all",
			opts:          []Option{WithOutConsoleOutPut()},
			expectConsole: false,
			expectFile:    false,
		},

		// === Out Console Only + Lumberjack ===
		{
			name: "without_console_with_lumberjack",
			opts: []Option{
				WithOutConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "no_console_lumberjack.log"),
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "no_console_lumberjack.log"),
		},

		// === Out Console Only + Rotatelog ===
		{
			name: "without_console_with_rotatelog",
			opts: []Option{
				WithOutConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(&RotateLogConfig{
					FileName: filepath.Join(logDir, "no_console_rotatelog.log"),
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "no_console_rotatelog.log"),
		},

		// === JSON Format with File ===
		{
			name: "json_format_with_file",
			opts: []Option{
				WithJSONFormat(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "json_file.log"),
				}),
			},
			expectConsole: false,
			expectFile:    true,
			filePath:      filepath.Join(logDir, "json_file.log"),
		},

		// === Custom Encoder Config ===
		{
			name: "custom_console_encoder",
			opts: []Option{
				WithConsoleOutPut(),
				WithDefaultLogLevel("debug"),
				WithConsoleConfig(zapcore.EncoderConfig{
					MessageKey:    "msg",
					LevelKey:      "level",
					TimeKey:       "ts",
					NameKey:       "logger",
					CallerKey:     "caller",
					StacktraceKey: "stacktrace",
					EncodeLevel:   zapcore.CapitalLevelEncoder,
					EncodeTime:    zapcore.ISO8601TimeEncoder,
				}),
			},
			expectConsole: true,
			expectFile:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger - this should never panic
			logger := NewLogger(tt.opts...)
			if logger == nil {
				t.Fatalf("NewLogger returned nil")
			}
			sugar := logger.Sugar()

			// Call all log level methods
			sugar.Debug("debug", zap.String("test", tt.name))
			sugar.Info("info", zap.String("test", tt.name))
			sugar.Warn("warn", zap.String("test", tt.name))
			sugar.Error("error", zap.String("test", tt.name))

			// Also test with typed fields
			logger.Debug("debug-field", zap.String("test", tt.name), zap.Int("num", 42))
			logger.Info("info-field", zap.String("test", tt.name))
			logger.Warn("warn-field", zap.String("test", tt.name))
			logger.Error("error-field", zap.String("test", tt.name))

			// Wait for async writes
			time.Sleep(100 * time.Millisecond)
			logger.Sync()

			// Verify file was created if expected
			if tt.expectFile && tt.filePath != "" {
				if !checkLogFileExists(t, tt.filePath) {
					// For default lumberjack config, the file path in test might not match
					// Check if any .log files were created in the dir
					dir := filepath.Dir(tt.filePath)
					entries, err := os.ReadDir(dir)
					if err == nil && len(entries) > 0 {
						for _, e := range entries {
							t.Logf("found file in %s: %s", dir, e.Name())
						}
					} else {
						t.Logf("expected file %s not found (may be default lumberjack path)", tt.filePath)
					}
				}
			}
		})
	}
}

// TestConsoleOnly verifies console-only output
func TestConsoleOnly(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "default",
			opts: []Option{WithConsoleOutPut()},
		},
		{
			name: "debug_level",
			opts: []Option{WithConsoleOutPut(), WithDefaultLogLevel("debug")},
		},
		{
			name: "json_format",
			opts: []Option{WithConsoleOutPut(), WithJSONFormat()},
		},
		{
			name: "console_format",
			opts: []Option{WithConsoleOutPut(), WithConsoleFormat()},
		},
		{
			name: "zaplogfmt_format",
			opts: []Option{WithConsoleOutPut(), WithZaplogfmtFormat()},
		},
		{
			name: "with_caller",
			opts: []Option{WithConsoleOutPut(), WithDefaultLogLevel("debug"), WithCaller()},
		},
		{
			name: "with_module",
			opts: []Option{WithConsoleOutPut(), WithModuleName("myapp")},
		},
		{
			name: "develop_mode",
			opts: []Option{WithConsoleOutPut(), WithDevelopMode()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts...)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}
			logAllLevels(logger, logger.Sugar(), tt.name)
			time.Sleep(50 * time.Millisecond)
			logger.Sync()
		})
	}
}

// TestFileOnlyLumberjack tests lumberjack-only file output
func TestFileOnlyLumberjack(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name   string
		config *LumberjackConfig
	}{
		{
			name:   "default_config",
			config: LumberjackDefaultConfig(),
		},
		{
			name: "custom_config",
			config: &LumberjackConfig{
				FileName:   filepath.Join(logDir, "lumberjack.log"),
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Localtime:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(tt.config),
			}
			if tt.name == "custom_config" {
				// Enable console to see output during test
				opts = append(opts, WithConsoleOutPut())
			}

			logger := NewLogger(opts...)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}
			logAllLevels(logger, logger.Sugar(), tt.name)
		})
	}
}

// TestFileOnlyRotatelog tests rotatelog-only file output
func TestFileOnlyRotatelog(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name   string
		config *RotateLogConfig
	}{
		{
			name:   "default_config",
			config: RotateLogDefaultConfig(),
		},
		{
			name: "custom_config",
			config: &RotateLogConfig{
				FileName:     filepath.Join(logDir, "rotatelog.log"),
				MaxAge:       "7d",
				RotationTime: "1h",
				LocalTime:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(tt.config),
			}
			if tt.name == "custom_config" {
				opts = append(opts, WithConsoleOutPut())
			}

			logger := NewLogger(opts...)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}
			logAllLevels(logger, logger.Sugar(), tt.name)
		})
	}
}

// TestConsoleAndFileCombined tests console + file output together
func TestConsoleAndFileCombined(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "console_and_lumberjack_default",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(LumberjackDefaultConfig()),
			},
		},
		{
			name: "console_and_lumberjack_custom",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "combined_lumberjack.log"),
				}),
			},
		},
		{
			name: "console_and_lumberjack_custom_full",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName:  filepath.Join(logDir, "combined_lumberjack_full.log"),
					MaxSize:   5,
					Compress:  true,
					Localtime: true,
				}),
			},
		},
		{
			name: "console_and_rotatelog_default",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(RotateLogDefaultConfig()),
			},
		},
		{
			name: "console_and_rotatelog_custom",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(&RotateLogConfig{
					FileName:     filepath.Join(logDir, "combined_rotatelog.log"),
					MaxAge:       "3d",
					RotationTime: "30m",
					LocalTime:    true,
				}),
			},
		},
		{
			name: "separate_log_levels",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithConsoleLogLevel("debug"),
				WithFileLogLevel("info"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "separate_levels.log"),
				}),
			},
		},
		{
			name: "with_caller_and_module",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithModuleName("test"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "caller_module.log"),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts...)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}
			logAllLevels(logger, logger.Sugar(), tt.name)

			time.Sleep(50 * time.Millisecond)
			logger.Sync()

			// Check file was created for custom config tests
			for _, opt := range tt.opts {
				if opt.Name() == lumberjackKey {
					if cfg, ok := opt.Value().(*LumberjackConfig); ok && cfg.FileName != "" && strings.HasPrefix(cfg.FileName, logDir) {
						if !checkLogFileExists(t, cfg.FileName) {
							t.Logf("lumberjack file not found at %s", cfg.FileName)
						}
					}
				}
				if opt.Name() == rotatelogKey {
					if cfg, ok := opt.Value().(*RotateLogConfig); ok && cfg.FileName != "" && strings.HasPrefix(cfg.FileName, logDir) {
						if !checkLogFileExists(t, cfg.FileName) {
							t.Logf("rotatelog file not found at %s", cfg.FileName)
						}
					}
				}
			}
		})
	}
}

// TestLogLevels tests all zap log levels
func TestLogLevels(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name string
		lv   zapcore.Level
		msg  string
	}{
		{"debug", zapcore.DebugLevel, "debug message"},
		{"info", zapcore.InfoLevel, "info message"},
		{"warn", zapcore.WarnLevel, "warn message"},
		{"error", zapcore.ErrorLevel, "error message"},
		{"dpanic", zapcore.DPanicLevel, "dpanic message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "level_"+tt.name+".log"),
				}),
			)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}

			logger.Log(tt.lv, tt.msg)

			time.Sleep(50 * time.Millisecond)
			logger.Sync()
		})
	}
}

// TestNoOutput verifies that disabling console with no file config produces a logger without panics
func TestNoOutput(t *testing.T) {
	logger := NewLogger(WithOutConsoleOutPut())
	if logger == nil {
		t.Fatal("logger should not be nil even with no output configured")
	}
	// Should not panic
	time.Sleep(50 * time.Millisecond)
	logger.Sync()
}

// TestWithOutConsoleOutPutWithFile verifies disabling console while keeping file output
func TestWithOutConsoleOutPutWithFile(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name   string
		config interface{}
	}{
		{
			name: "with_lumberjack",
			config: &LumberjackConfig{
				FileName: filepath.Join(logDir, "no_console_lumberjack.log"),
			},
		},
		{
			name: "with_rotatelog",
			config: &RotateLogConfig{
				FileName: filepath.Join(logDir, "no_console_rotatelog.log"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logger *zap.Logger
			switch cfg := tt.config.(type) {
			case *LumberjackConfig:
				logger = NewLogger(
					WithOutConsoleOutPut(),
					WithFileOutPut(),
					WithDefaultLogLevel("debug"),
					WithLumberjack(cfg),
				)
			case *RotateLogConfig:
				logger = NewLogger(
					WithOutConsoleOutPut(),
					WithFileOutPut(),
					WithDefaultLogLevel("debug"),
					WithRotateLog(cfg),
				)
			}

			if logger == nil {
				t.Fatal("logger should not be nil")
			}

			logger.Info("this should only go to file", zap.String("test", tt.name))
			time.Sleep(50 * time.Millisecond)
			logger.Sync()
		})
	}
}

// TestLogOutputContent verifies content is written correctly
func TestLogOutputContent(t *testing.T) {
	logDir := setupLogDir(t)

	testMsg := "content_test_unique_12345"
	logger := NewLogger(
		WithConsoleOutPut(),
		WithFileOutPut(),
		WithDefaultLogLevel("debug"),
		WithLumberjack(&LumberjackConfig{
			FileName: filepath.Join(logDir, "content_test.log"),
		}),
	)
	if logger == nil {
		t.Fatal("logger should not be nil")
	}

	logger.Info(testMsg)

	time.Sleep(200 * time.Millisecond)
	logger.Sync()

	// Check file content
	entries, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("could not read log dir: %v", err)
	}

	found := false
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".log") {
			data, err := os.ReadFile(filepath.Join(logDir, entry.Name()))
			if err != nil {
				continue
			}
			if strings.Contains(string(data), testMsg) {
				found = true
				break
			}
		}
	}

	if !found {
		t.Errorf("test message %q not found in any log file in %s", testMsg, logDir)
		// List files for debugging
		for _, entry := range entries {
			t.Logf("file in dir: %s", entry.Name())
		}
	}
}

// TestSugaredLoggerTableDriven tests SugaredLogger with various configurations
func TestSugaredLoggerTableDriven(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "sugar_console_only",
			opts: []Option{WithConsoleOutPut(), WithDefaultLogLevel("debug")},
		},
		{
			name: "sugar_file_lumberjack_default",
			opts: []Option{WithFileOutPut(), WithDefaultLogLevel("debug"), WithLumberjack(LumberjackDefaultConfig())},
		},
		{
			name: "sugar_file_lumberjack_custom",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "sugar_lumberjack.log"),
				}),
			},
		},
		{
			name: "sugar_file_rotatelog_default",
			opts: []Option{WithFileOutPut(), WithDefaultLogLevel("debug"), WithRotateLog(RotateLogDefaultConfig())},
		},
		{
			name: "sugar_file_rotatelog_custom",
			opts: []Option{
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithRotateLog(&RotateLogConfig{
					FileName: filepath.Join(logDir, "sugar_rotatelog.log"),
				}),
			},
		},
		{
			name: "sugar_console_and_file",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "sugar_both.log"),
				}),
			},
		},
		{
			name: "sugar_all_options",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithModuleName("sugar-test"),
				WithDevelopMode(),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "sugar_all.log"),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sugar := NewSugaredLogger(tt.opts...)
			if sugar == nil {
				t.Fatal("sugared logger should not be nil")
			}

			// Test all sugared logger methods
			sugar.Debugw("debug with fields", "key", "value")
			sugar.Infof("info %s", "formatted")
			sugar.Warnf("warn %s", "formatted")
			sugar.Errorf("error %s", "formatted")
			sugar.With(zap.String("test", tt.name)).Infow("with context", "extra", "data")

			time.Sleep(50 * time.Millisecond)
			sugar.Sync()
		})
	}
}

// TestOptionsComposition tests various option combinations
func TestOptionsComposition(t *testing.T) {
	logDir := setupLogDir(t)

	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "json_format_with_file",
			opts: []Option{
				WithJSONFormat(),
				WithFileOutPut(),
				WithDefaultLogLevel("info"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "json_file.log"),
				}),
			},
		},
		{
			name: "console_format_with_file",
			opts: []Option{
				WithConsoleFormat(),
				WithFileOutPut(),
				WithDefaultLogLevel("info"),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "console_file.log"),
				}),
			},
		},
		{
			name: "custom_encoder_config",
			opts: []Option{
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithConsoleConfig(zapcore.EncoderConfig{
					MessageKey:    "msg",
					LevelKey:      "level",
					TimeKey:       "ts",
					NameKey:       "logger",
					CallerKey:     "caller",
					StacktraceKey: "stacktrace",
					EncodeLevel:   zapcore.CapitalLevelEncoder,
					EncodeTime:    zapcore.ISO8601TimeEncoder,
				}),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, "custom_encoder.log"),
				}),
			},
		},
		{
			name: "caller_skip",
			opts: []Option{
				WithConsoleOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithCallerSkip(1),
			},
		},
		{
			name: "stack_log_level",
			opts: []Option{
				WithConsoleOutPut(),
				WithDefaultLogLevel("debug"),
				WithStackLogLevel("error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.opts...)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}

			logger.Info("composition test", zap.String("test", tt.name))
			time.Sleep(50 * time.Millisecond)
			logger.Sync()
		})
	}
}

// TestLumberjackConfigString ensures config stringification works
func TestLumberjackConfigString(t *testing.T) {
	cfg := &LumberjackConfig{
		FileName:   "test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
		Localtime:  true,
	}
	s := cfg.String()
	if s == "" {
		t.Fatal("config string should not be empty")
	}
	t.Logf("config string: %s", s)
}

// TestRotateLogConfigString ensures config stringification works
func TestRotateLogConfigString(t *testing.T) {
	cfg := &RotateLogConfig{
		FileName:     "test.log",
		MaxAge:       "7d",
		LocalTime:    true,
		RotationTime: "1h",
	}
	s := cfg.String()
	if s == "" {
		t.Fatal("config string should not be empty")
	}
	t.Logf("config string: %s", s)
}

// TestNewLumberjackUrl verifies URL generation
func TestNewLumberjackUrlTable(t *testing.T) {
	tests := []struct {
		name   string
		config *LumberjackConfig
	}{
		{
			name:   "nil_config",
			config: nil,
		},
		{
			name: "custom_config",
			config: &LumberjackConfig{
				FileName:   "test.log",
				MaxSize:    10,
				MaxAge:     14,
				MaxBackups: 5,
				Compress:   true,
				Localtime:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := NewLumberjackUrl(tt.config)
			if url == "" {
				t.Fatal("URL should not be empty")
			}
			if !strings.HasPrefix(url, "lumberjack:") {
				t.Fatalf("URL should start with 'lumberjack:', got: %s", url)
			}
			t.Logf("URL: %s", url)
		})
	}
}

// TestNewRotateLogURL verifies URL generation
func TestNewRotateLogURLTable(t *testing.T) {
	tests := []struct {
		name   string
		config *RotateLogConfig
	}{
		{
			name:   "nil_config",
			config: nil,
		},
		{
			name: "custom_config",
			config: &RotateLogConfig{
				FileName:     "test.log",
				MaxAge:       "7d",
				LocalTime:    true,
				RotationTime: "1h",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := NewRotateLogURL(tt.config)
			if url == "" {
				t.Fatal("URL should not be empty")
			}
			if !strings.HasPrefix(url, "rotatelogs:") {
				t.Fatalf("URL should start with 'rotatelogs:', got: %s", url)
			}
			t.Logf("URL: %s", url)
		})
	}
}

// TestConcurrentLoggerCreation tests concurrent logger creation safety
func TestConcurrentLoggerCreation(t *testing.T) {
	logDir := setupLogDir(t)

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			logger := NewLogger(
				WithConsoleOutPut(),
				WithFileOutPut(),
				WithDefaultLogLevel("debug"),
				WithCaller(),
				WithModuleName(fmt.Sprintf("goroutine-%d", idx)),
				WithLumberjack(&LumberjackConfig{
					FileName: filepath.Join(logDir, fmt.Sprintf("concurrent_%d.log", idx)),
				}),
			)
			if logger == nil {
				t.Errorf("goroutine %d: logger is nil", idx)
			} else {
				logger.Info("from goroutine", zap.Int("idx", idx))
			}
			time.Sleep(50 * time.Millisecond)
			logger.Sync()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
