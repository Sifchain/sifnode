package siflogger

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

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

var filterNumber = 0

func LogFileLine(depth int) kitlog.Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(depth + filterNumber)
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

func New() SifLogger {
	logger := log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)
	logger = logger.With("caller", LogFileLine(6))
	e := SifLogger{logger}
	return e
}

func (e *SifLogger) SetGlobalFilter(level Level) {
	var option log.Option
	switch level {
	case Debug:
		option = log.AllowDebug()
	case Error:
		option = log.AllowError()
	case Info:
		option = log.AllowInfo()
	case None:
		option = log.AllowNone()
	default:
		panic("incorrect level")
	}
	filter := log.NewFilter(e.logger, option)
	e.logger = filter
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

	filterNumber++
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

func (e SifLogger) With(keyvals ...interface{}) log.Logger {
	return SifLogger{e.logger.With(keyvals...)}
}
