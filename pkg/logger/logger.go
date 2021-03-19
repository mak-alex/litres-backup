package logger

import (
	"errors"
	"fmt"
	"github.com/mak-alex/litres-backup/pkg/conf"
	"go.uber.org/zap"
)

var Work *zap.Logger

const (
	InstanceZapLogger int = iota
)

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

//NewLogger returns an instance of logger
func NewLogger(config conf.LogConf, loggerInstance int) *zap.Logger {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(config)
		if err != nil {
			fmt.Printf("%s", err)
		}
		return logger
	default:
		fmt.Printf("%s", errInvalidLoggerInstance)
	}
	return nil
}
