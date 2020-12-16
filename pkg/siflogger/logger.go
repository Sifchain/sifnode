package siflogger

import (
	"errors"
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
type Type byte

const (
	defaultLogLevelKey = "*"
)

const (
	Debug Level = iota
	Info
	Error
	None
)

const (
	TDFmt Type = iota
	JSON
	Custom
)

type Logger struct {
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

func logFileLine(depth int) kitlog.Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(depth + filterNumber)
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// parseLogLevel parses complex log level - comma-separated
// list of module:level pairs with an optional *:level pair (* means
// all other modules).
//
// Example:
//		parseLogLevel("module","consensus:debug,mempool:debug,*:error", "info")
func parseLogLevel(firstLayerKey string, lvl string, defaultLogLevelValue string) ([]log.Option, error) {
	if lvl == "" {
		return nil, errors.New("empty log level")
	}

	l := lvl

	// prefix simple one word levels (e.g. "info") with "*"
	if !strings.Contains(l, ":") {
		l = defaultLogLevelKey + ":" + l
	}

	options := make([]log.Option, 0)

	isDefaultLogLevelSet := false
	var option log.Option
	var err error

	list := strings.Split(l, ",")
	for _, item := range list {
		nameAndLevel := strings.Split(item, ":")

		if len(nameAndLevel) != 2 {
			return nil, fmt.Errorf("expected list in a form of \"module:level\" pairs, given pair %s, list %s", item, list)
		}

		name := nameAndLevel[0]
		level := nameAndLevel[1]

		if name == defaultLogLevelKey {
			option, err = log.AllowLevel(level)
			if err != nil {
				return nil, fmt.Errorf("failed to parse default log level (pair %s, list %s): %w", item, l, err)
			}
			options = append(options, option)
			isDefaultLogLevelSet = true
		} else {
			switch level {
			case "debug":
				option = log.AllowDebugWith(firstLayerKey, name)
			case "info":
				option = log.AllowInfoWith(firstLayerKey, name)
			case "error":
				option = log.AllowErrorWith(firstLayerKey, name)
			case "none":
				option = log.AllowNoneWith(firstLayerKey, name)
			default:
				return nil,
					fmt.Errorf("expected either \"info\", \"debug\", \"error\" or \"none\" log level, given %s (pair %s, list %s)",
						level,
						item,
						list)
			}
			options = append(options, option)
		}
	}

	// if "*" is not provided, set default global level
	if !isDefaultLogLevelSet {
		option, err = log.AllowLevel(defaultLogLevelValue)
		if err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, nil
}

func New(t Type) Logger {
	var logger log.Logger
	switch t {
	case TDFmt:
		logger = newFmtLogger(log.NewSyncWriter(os.Stdout), colorFn)
	case JSON:
		logger = newJSONLogger(log.NewSyncWriter(os.Stdout), colorFn)
	case Custom:
		logger = newCustomLogger(log.NewSyncWriter(os.Stdout), colorFn)
	default:
		panic("incorrect logger type")
	}
	if filterNumber != 0 {
		// TODO: the only problem (due not full OOP support) is to implement correct filterNumber per instance but not per package
		panic("Logger is supposed to be a singleton")
	}
	filterNumber++
	logger = logger.With("caller", logFileLine(4))
	e := Logger{logger}
	return e
}

func (e *Logger) SetGlobalFilter(level Level) {
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
	e.setNextFilter([]log.Option{option})
}

func (e *Logger) setNextFilter(options []log.Option) {
	filterNumber++
	filter := log.NewFilter(e.logger, options...)
	e.logger = filter
}

func (e *Logger) SetFilterForLayer(level Level, keyvals ...interface{}) {
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

	e.setNextFilter(options)
}

func (e *Logger) SetFilterForLayerFromConfig(firstLayer string, lvl string, defaultLogLevelValue string) error {
	options, err := parseLogLevel(firstLayer, lvl, defaultLogLevelValue)
	if err != nil {
		return err
	}
	e.setNextFilter(options)
	return nil
}

func (e Logger) Debug(msg string, keyvals ...interface{}) {
	e.logger.Debug(msg, keyvals...)
}

func (e Logger) Info(msg string, keyvals ...interface{}) {
	e.logger.Info(msg, keyvals...)
}

func (e Logger) Error(msg string, keyvals ...interface{}) {
	e.logger.Error(msg, keyvals...)
}

func (e Logger) With(keyvals ...interface{}) log.Logger {
	return Logger{e.logger.With(keyvals...)}
}
