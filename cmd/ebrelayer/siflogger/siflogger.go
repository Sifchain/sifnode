package siflogger

import (
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/tendermint/tendermint/libs/log"
)

type Level byte

const (
	Debug Level = iota
	Info
	Error
	None
)

type SifLogger struct {
	logger log.Logger
}

func colorFn(keyvals ...interface{}) term.FgBgColor {
	if keyvals[0] != level.Key() {
		panic(fmt.Sprintf("expected level key to be first, got %v", keyvals[0]))
	}

	switch keyvals[1].(level.Value).String() {
	case "debug":
		return term.FgBgColor{Fg: term.DarkGreen}
	case "error":
		return term.FgBgColor{Fg: term.Red}
	default:
		return term.FgBgColor{}
	}
}

func New() SifLogger {
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)
	logger = logger.With("caller", kitlog.Caller(5))
	e := SifLogger{logger}
	return e
}

func (e *SifLogger) SetFilterForLayer(level Level, keyvals ...interface{}) {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		panic("invalid keyvals")
	}

	options := make([]log.Option, len(keyvals)/2, len(keyvals)/2)
	for i := 0; i < len(keyvals)/2; i++ {
		switch level {
		case Debug:
			options[i] = log.AllowDebugWith(keyvals[i*2], keyvals[i*2+1])
		case Error:
			options[i] = log.AllowErrorWith(keyvals[i*2], keyvals[i*2+1])
		case Info:
			options[i] = log.AllowInfoWith(keyvals[i*2], keyvals[i*2+1])
		case None:
			options[i] = log.AllowNoneWith(keyvals[i*2], keyvals[i*2+1])
		default:
			panic("incorrect level")
		}
	}

	filter := log.NewFilter(e.logger, options...)
	e.logger = filter
}

func (e SifLogger) Debug(msg string, keyvals ...interface{}) {
	e.logger.Debug(msg, keyvals...)
}

func (e SifLogger) Info(msg string, keyvals ...interface{}) {
	e.logger.Info(msg, keyvals...)
}

func (e SifLogger) Error(msg string, keyvals ...interface{}) {
	e.logger.Error(msg, keyvals...)
}

func (e SifLogger) Tag(keyvals ...interface{}) SifLogger {
	return SifLogger{e.logger.With(keyvals...)}
}

func (e SifLogger) GetTendermintLogger(keyvals ...interface{}) log.Logger {
	return e.logger
}
