package log

import (
	"path/filepath"
	"strings"

	"github.com/gtkit/logger"

	"go_sleep_admin/internal/pkg/news"
	v2config "go_sleep_admin/internal/platform/config"
)

type logRuntimeConfig struct {
	Path    string
	Level   string
	Console bool
	File    bool
}

type newsSender interface {
	Send(msg string)
	SendTo(url string, msg string)
}

const (
	defaultLogPath  = "./logs/app_client"
	defaultLogLevel = "info"
)

var logInitializer = func(cfg logRuntimeConfig, sender newsSender) {
	logger.NewZap(
		logger.WithPath(cfg.Path),
		logger.WithLevel(cfg.Level),
		logger.WithConsole(cfg.Console),
		logger.WithFile(cfg.File),
		logger.WithDivision("daily"),
		logger.WithMessager(sender),
	)
}

func Init() {
	InitWithConfig(&v2config.Config{})
}

func InitWithConfig(cfg *v2config.Config) {
	runtime := buildLogRuntimeConfig(cfg)
	logInitializer(runtime, &News{notifier: news.NewWithConfig(cfg)})
}

func buildLogRuntimeConfig(cfg *v2config.Config) logRuntimeConfig {
	if cfg == nil {
		cfg = &v2config.Config{}
	}

	return logRuntimeConfig{
		Path:    resolveLogPath(cfg.Log.Path),
		Level:   resolveLogLevel(cfg.Log.Level),
		Console: cfg.Log.ConsoleStdout != 0,
		File:    cfg.Log.FileStdout != 0,
	}
}

func resolveLogPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return defaultLogPath
	}

	return filepath.Clean(path)
}

func resolveLogLevel(level string) string {
	level = strings.TrimSpace(level)
	if level == "" {
		return defaultLogLevel
	}

	return level
}

type News struct {
	notifier *news.Notifier
}

func (n *News) Send(msg string)               { n.notifier.Warn("", msg) }
func (n *News) SendTo(url string, msg string) { n.notifier.Warn(url, msg) }

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func resetLogHooksForTest() {
	logInitializer = func(cfg logRuntimeConfig, sender newsSender) {
		logger.NewZap(
			logger.WithPath(cfg.Path),
			logger.WithLevel(cfg.Level),
			logger.WithConsole(cfg.Console),
			logger.WithFile(cfg.File),
			logger.WithDivision("daily"),
			logger.WithMessager(sender),
		)
	}
}

func setLogInitializerForTest(fn func(cfg logRuntimeConfig, sender newsSender)) { logInitializer = fn }
