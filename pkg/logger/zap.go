package logger

import (
	"github.com/mattn/go-colorable"

	"github.com/mak-alex/backlitr/pkg/conf"
	"github.com/mak-alex/backlitr/pkg/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ *zap.Logger

func getEncoder(isJSON bool) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if conf.GlobalConfig.Mode == "debug" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // for color
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case consts.Info:
		return zapcore.InfoLevel
	case consts.Warn:
		return zapcore.WarnLevel
	case consts.Debug:
		return zapcore.DebugLevel
	case consts.Error:
		return zapcore.ErrorLevel
	case consts.Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config conf.LogConf) (*zap.Logger, error) {
	var cores []zapcore.Core

	if config.EnableConsole {
		level := getZapLevel(config.ConsoleLevel)
		writer := zapcore.Lock(zapcore.AddSync(colorable.NewColorableStdout()))
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}

	if config.EnableFile {
		level := getZapLevel(config.FileLevel)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  config.MaxSize,
			Compress: config.Compress,
			MaxAge:   config.MaxAge,
		})
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	)

	return logger, nil
}
