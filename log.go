package logium

import (
	"errors"
	"strings"

	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

func NewLogger(level, format string) *logrus.Logger {
	log := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Warnf("invalid log level '%s', defaulting to 'info'", level)
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	switch strings.ToLower(format) {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}

type Logger interface {
	WithError(err error) *logrus.Entry

	logrus.FieldLogger
}

type logger struct {
	*logrus.Entry
}

func (l *logger) WithError(err error) *logrus.Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return l.Entry.WithError(ae)
	}
	return l.Entry.WithError(err)
}

func NewWithBase(base *logrus.Logger) Logger {
	log := logger{
		Entry: logrus.NewEntry(base),
	}

	return &log
}
