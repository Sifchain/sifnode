package siflogger

import (
	"io"

	kitlog "github.com/go-kit/kit/log"
	kitlevel "github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	msgKey = "_msg" // "_" prefixed to avoid collisions
)

type sifImplLogger struct {
	srcLogger kitlog.Logger
}

var _ log.Logger = (*sifImplLogger)(nil)

func newFmtLogger(w io.Writer, colorFn func(keyvals ...interface{}) term.FgBgColor) log.Logger {
	return &sifImplLogger{term.NewLogger(w, log.NewTMFmtLogger, colorFn)}
}

func (l *sifImplLogger) Info(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Info(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *sifImplLogger) Debug(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Debug(l.srcLogger)
	if err := kitlog.With(lWithLevel, msgKey, msg).Log(keyvals...); err != nil {
		errLogger := kitlevel.Error(l.srcLogger)
		kitlog.With(errLogger, msgKey, msg).Log("err", err)
	}
}

func (l *sifImplLogger) Error(msg string, keyvals ...interface{}) {
	lWithLevel := kitlevel.Error(l.srcLogger)
	lWithMsg := kitlog.With(lWithLevel, msgKey, msg)
	if err := lWithMsg.Log(keyvals...); err != nil {
		lWithMsg.Log("err", err)
	}
}

func (l *sifImplLogger) With(keyvals ...interface{}) log.Logger {
	return &sifImplLogger{kitlog.With(l.srcLogger, keyvals...)}
}
