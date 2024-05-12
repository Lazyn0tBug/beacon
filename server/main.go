// main.go

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
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
		file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
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

func main() {
	// Load the TOML configuration file
	loggingConfig, lumberjackConfig, err := LoadConfigFromTOML("logging.toml")
	if err != nil {
		panic(err)
	}

	// Create the logger based on the configuration
	logger, err := newLoggerFromConfig(loggingConfig, lumberjackConfig)
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	// Now you can use the logger
	logger.Info("This is an info message")
	logger.Error("This is an error message")

	for i := 0; i < 12; i++ {
		go logger.Info(fmt.Sprintf("test log: %d", i))
	}
	time.Sleep(time.Second * 3)
}
