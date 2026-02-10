package logium

import (
	"context"
	"errors"
	"time"

	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Entry: logrus.NewEntry(logger.Logger),
	}
}

func (e *Entry) WithField(key string, value any) *Entry {
	return &Entry{
		Entry: e.Entry.WithField(key, value),
	}
}

func (e *Entry) WithFields(fields Fields) *Entry {
	return &Entry{
		Entry: e.Entry.WithFields(logrus.Fields(fields)),
	}
}

func (e *Entry) WithContext(ctx context.Context) *Entry {
	return &Entry{
		Entry: e.Entry.WithContext(ctx),
	}
}

func (e *Entry) WithTime(t time.Time) *Entry {
	return &Entry{
		Entry: e.Entry.WithTime(t),
	}
}

func (e *Entry) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: e.Entry.WithError(ae)}
	}
	return &Entry{Entry: e.Entry.WithError(err)}
}
