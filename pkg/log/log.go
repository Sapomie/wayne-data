// @Author : liguoyu
// @Date: 2019/10/29 15:42
package log

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/go-stack/stack"
)

// DefaultCallerDepth is default caller depth
const (
	LevelNone  = "none"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelAll   = "all"

	DefaultCallerDepth = 7
)

var defaultFileLogger = NewLogger("info", 1)
var defaultLogger = NewLogger("info", 1)

func SetDefaultFileLogger(logger Logger) {
	defaultFileLogger = logger
}

func GetDefaultLogger() Logger {
	return defaultFileLogger
}

// Logger interface
type Logger interface {
	log.Logger

	Fatal(keyvals ...interface{}) error
	Error(keyvals ...interface{}) error
	Warn(keyvals ...interface{}) error
	Info(keyvals ...interface{}) error
	Debug(keyvals ...interface{}) error
}

// Fatal write fatal level log
func Fatal(keyvals ...interface{}) (err error) {
	err = defaultFileLogger.Fatal(keyvals...)
	err = defaultLogger.Fatal(keyvals...)
	return
}

// Error write error level log
func Error(keyvals ...interface{}) (err error) {
	err = defaultFileLogger.Error(keyvals...)
	err = defaultLogger.Error(keyvals...)

	return
}

// Warn write warn level log
func Warn(keyvals ...interface{}) (err error) {
	err = defaultFileLogger.Warn(keyvals...)
	err = defaultLogger.Warn(keyvals...)
	return
}

// Info write warn level log
func Info(keyvals ...interface{}) (err error) {
	err = defaultFileLogger.Info(keyvals...)
	err = defaultLogger.Info(keyvals...)
	return
}

// Debug write debug level log
func Debug(keyvals ...interface{}) (err error) {
	err = defaultFileLogger.Debug(keyvals...)
	err = defaultLogger.Debug(keyvals...)
	return
}

func getLevelOption(levelStr string) level.Option {
	switch strings.ToLower(levelStr) {
	case LevelNone:
		return level.AllowNone()
	case LevelDebug:
		return level.AllowDebug()
	case LevelInfo:
		return level.AllowInfo()
	case LevelWarn:
		return level.AllowWarn()
	case LevelError:
		return level.AllowError()
	case LevelAll:
		return level.AllowAll()
	default:
		return level.AllowAll()
	}
}

type logger struct {
	log.Logger
}

func colorFn(keyvals ...interface{}) term.FgBgColor {
	for i := 0; i < len(keyvals)-1; i += 2 {
		if keyvals[i] != "level" {
			continue
		}

		levelValue := keyvals[i+1]
		if l, ok := keyvals[i+1].(level.Value); ok {
			levelValue = l.String()
		}

		switch levelValue {
		case "debug":
			return term.FgBgColor{Fg: term.Gray}
		case "info":
			return term.FgBgColor{Fg: term.Green}
		case "warn":
			return term.FgBgColor{Fg: term.Yellow}
		case "error":
			return term.FgBgColor{Fg: term.Red}
		default:
			return term.FgBgColor{}
		}
	}
	return term.FgBgColor{}
}

func dealError(err error) {
	fmt.Println("err:", err)
}

// NewLogger return a Logger
func NewLogger(levelStr string, deltaDepth int) Logger {
	return NewFileLogger(os.Stdout, levelStr, deltaDepth)
}

func NewFileLogger(writer io.Writer, levelStr string, deltaDepth int) Logger {
	l := term.NewLogger(writer, log.NewLogfmtLogger, colorFn)

	l = log.With(l, "time", log.DefaultTimestamp)

	l = log.With(l, "caller", log.Valuer(func() interface{} {
		return fmt.Sprintf("%+v", stack.Caller(DefaultCallerDepth+deltaDepth))
	}))

	l = level.NewFilter(l, getLevelOption(levelStr))
	return &logger{l}
}

// With is a wrapper of go-kit log With
func With(lg Logger, keyvals ...interface{}) Logger {
	l, ok := lg.(*logger)
	if !ok {
		panic(errors.New("Error: param logger is not a *logger. "))
	}

	l.Logger = log.With(l.Logger, keyvals...)
	return l
}

// WithPrefix is a wrapper of go-kit log WithPrefix
func WithPrefix(lg Logger, keyvals ...interface{}) Logger {
	l, ok := lg.(*logger)
	if !ok {
		panic(errors.New("Error: param logger is not a *logger. "))
	}

	l.Logger = log.WithPrefix(l.Logger, keyvals...)
	return l
}

func (l *logger) Log(keyvals ...interface{}) error {
	return l.Logger.Log(keyvals...)
}

func (l *logger) Fatal(keyvals ...interface{}) (err error) {
	if err = level.Error(l.Logger).Log(keyvals...); err != nil {
		dealError(err)
	}
	os.Exit(1)
	return nil
}

func (l logger) Error(keyvals ...interface{}) (err error) {
	if err = level.Error(l.Logger).Log(keyvals...); err != nil {
		dealError(err)
	}
	return
}

func (l logger) Warn(keyvals ...interface{}) (err error) {
	if err = level.Warn(l.Logger).Log(keyvals...); err != nil {
		dealError(err)
	}
	return
}

func (l logger) Info(keyvals ...interface{}) (err error) {
	if err = level.Info(l.Logger).Log(keyvals...); err != nil {
		dealError(err)
	}
	return
}

func (l logger) Debug(keyvals ...interface{}) (err error) {
	if err = level.Debug(l.Logger).Log(keyvals...); err != nil {
		dealError(err)
	}
	return
}
