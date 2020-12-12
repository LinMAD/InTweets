package core

import "github.com/sirupsen/logrus"

// Logger implementation
type Logger struct {
	*logrus.Logger
}

// Critical ...
func (l *Logger) Critical(args ...interface{}) {
	l.Error(args...)
}

// Criticalf ...
func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

// Notice ...
func (l *Logger) Notice(args ...interface{}) {
	l.Info(args...)
}

// Noticef ...
func (l *Logger) Noticef(format string, args ...interface{}) {
	l.Infof(format, args...)
}
