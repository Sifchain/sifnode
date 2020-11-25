package siflogger

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/tendermint/tendermint/libs/log"
)

type SifLogger struct {
	logger log.Logger
}

func New() SifLogger {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = logger.With("caller", kitlog.Caller(5))
	e := SifLogger{logger}
	return e
}

func (e SifLogger) Debug(msg string, keyvals ...interface{}) {
	e.logger.Debug(msg, keyvals)
}

func (e SifLogger) Info(msg string, keyvals ...interface{}) {
	e.logger.Info(msg, keyvals)
}

func (e SifLogger) Error(msg string, keyvals ...interface{}) {
	e.logger.Error(msg, keyvals)
}

func (e SifLogger) GetTendermintLogger(keyvals ...interface{}) log.Logger {
	return e.logger
}
