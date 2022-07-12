package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"runtime"
)

type ILogger interface {
	Debug(args ...interface{})
	DebugF(format string, args ...interface{})

	Info(args ...interface{})
	InfoF(format string, args ...interface{})

	Warning(args ...interface{})
	WarningF(format string, args ...interface{})

	Error(args ...interface{})
	ErrorF(format string, args ...interface{})

	FatalF(format string, args ...interface{})

	MiddlewareLogging(next http.Handler) http.Handler
}

type logger struct {
	l *logrus.Logger
}

func New(loggingLevel string, tsFormat string) (*logger, error) {
	l := logrus.New()

	l.Out = os.Stdout
	l.SetReportCaller(true)

	if lvl, err := logrus.ParseLevel(loggingLevel); err != nil {
		return nil, err
	} else {
		l.SetLevel(lvl)
	}

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: tsFormat,
		PrettyPrint:     false,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})
	return &logger{
		l,
	}, nil

}

func (l *logger) Debug(args ...interface{}) {
	l.l.Debug(args...)
}

func (l *logger) DebugF(format string, args ...interface{}) {
	l.l.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.l.Info(args...)
}

func (l *logger) InfoF(format string, args ...interface{}) {
	l.l.Infof(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.l.Warning(args...)
}

func (l *logger) WarningF(format string, args ...interface{}) {
	l.l.Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.l.Error(args...)
}

func (l *logger) ErrorF(format string, args ...interface{}) {
	l.l.Errorf(format, args...)
}

func (l *logger) FatalF(format string, args ...interface{}) {
	l.l.Fatalf(format, args...)
}
