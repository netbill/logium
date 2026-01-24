package logium

import (
	"errors"

	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type Entry struct {
	*logrus.Entry
}

func New() *Logger {
	return &Logger{Logger: logrus.New()}
}

func NewWithBase(base *logrus.Logger) *Logger {
	return &Logger{Logger: base}
}

func (l *Logger) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: l.Logger.WithError(ae)}
	}
	return &Entry{Entry: l.Logger.WithError(err)}
}

func (e *Entry) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: e.Entry.WithError(ae)}
	}
	return &Entry{Entry: e.Entry.WithError(err)}
}
