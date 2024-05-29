// utils/logger.go
package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	// 用来保证Logger的线程安全
	once sync.Once
	// 用来持有全局的Logger实例
	globalLogger      *zap.Logger
	loggerInitialized = false
)

// LoggingConfig represents the logging configuration.
type LoggingConfig struct {
	Level            string   `toml:"level"`
	OutputPaths      []string `toml:"outputPaths"`
	ErrorOutputPaths []string `toml:"errorOutputPaths"`
	Encoding         string   `toml:"encoding"`
}

// LumberjackConfig represents the lumberjack-specific configuration.
type LumberjackConfig struct {
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"maxsize"`
	MaxBackups int    `toml:"maxbackups"`
	MaxAge     int    `toml:"maxage"`
	Compress   bool   `toml:"compress"`
}

// InitializeLogger 初始化全局的Logger实例
func InitializeLogger() {
	once.Do(func() {
		config_file, err := getLoggingConfigFile()
		if err != nil {
			config_file = "logging.toml"
		}
		// Load the TOML configuration file
		loggingConfig, lumberjackConfig, err := LoadConfigFromTOML(config_file)
		if err != nil {
			panic(err)
		}

		// Create the logger based on the configuration
		logger, err := newLoggerFromConfig(loggingConfig, lumberjackConfig)
		if err != nil {
			panic(err)
		}

		loggerInitialized = true
		globalLogger = logger
	})
}

// GetLogger 返回全局的Logger实例，该实例是线程安全的。
func GetLogger() *zap.Logger {
	if !loggerInitialized {
		InitializeLogger()
	}
	return globalLogger
}

// LoadConfigFromTOML loads the logging configuration from a TOML file.
func LoadConfigFromTOML(filename string) (*LoggingConfig, *LumberjackConfig, error) {
	var config struct {
		Logging    LoggingConfig
		Lumberjack LumberjackConfig
	}

	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return nil, nil, err
	}

	return &config.Logging, &config.Lumberjack, nil
}

// newLoggerFromConfig creates a new zap logger based on the provided configuration.
func newLoggerFromConfig(loggingConfig *LoggingConfig, lumberjackConfig *LumberjackConfig) (*zap.Logger, error) {
	// Parse the level
	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	err := atomicLevel.UnmarshalText([]byte(loggingConfig.Level))
	if err != nil {
		return nil, err
	}

	// Create the lumberjack syncer
	lumberJackSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   lumberjackConfig.Filename,
		MaxSize:    lumberjackConfig.MaxSize,
		MaxBackups: lumberjackConfig.MaxBackups,
		MaxAge:     lumberjackConfig.MaxAge,
		Compress:   lumberjackConfig.Compress,
	})

	// Create the encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Prepare the syncers slice
	var syncers []zapcore.WriteSyncer
	for _, outputPath := range loggingConfig.OutputPaths {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			panic("failed to create log directory: " + err.Error())
		}

		file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		// defer file.Close()
		syncers = append(syncers, zapcore.AddSync(file))
	}
	syncers = append(syncers, lumberJackSyncer)

	// Create the core with the desired options
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel,
	)

	// Create the logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
