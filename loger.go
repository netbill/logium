package logium

import (
	"context"
	"errors"
	"time"

	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	return &Logger{Logger: logrus.New()}
}

func (l *Logger) newEntry() *Entry {
	return &Entry{Entry: logrus.NewEntry(l.Logger)}
}

func (l *Logger) WithField(key string, value any) *Entry {
	return &Entry{Entry: l.Logger.WithField(key, value)}
}

func (l *Logger) WithFields(fields Fields) *Entry {
	return &Entry{Entry: l.Logger.WithFields(logrus.Fields(fields))}
}

func (l *Logger) WithContext(ctx context.Context) *Entry {
	return &Entry{Entry: l.Logger.WithContext(ctx)}
}

func (l *Logger) WithTime(t time.Time) *Entry {
	return &Entry{Entry: l.Logger.WithTime(t)}
}

func (l *Logger) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: l.Logger.WithError(ae)}
	}
	return &Entry{Entry: l.Logger.WithError(err)}
}
