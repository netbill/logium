package logium

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	return &Logger{Logger: logrus.New()}
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

func (l *Logger) WithRequest(r *http.Request) *Entry {
	return l.WithFields(Fields{
		HTTPMethodField: r.Method,
		HTTPPathField:   r.URL.Path,
	})
}

type AccountAuth interface {
	GetAccountID() uuid.UUID
	GetSessionID() uuid.UUID
}

func (l *Logger) WithAccount(auth AccountAuth) *Entry {
	return l.WithFields(Fields{
		AccountIDField:        auth.GetAccountID(),
		AccountSessionIDField: auth.GetSessionID(),
	})
}

type uploadSessionContent interface {
	GetUploadOwnerAccountID() uuid.UUID
	GetUploadSessionID() uuid.UUID
	GetUploadResourceID() string
	GetUploadResource() string
}

func (l *Logger) WithUploadSession(tokens uploadSessionContent) *Entry {
	return l.WithFields(Fields{
		UploadOwnerAccountIdField: tokens.GetUploadOwnerAccountID(),
		UploadSessionIdField:      tokens.GetUploadSessionID(),
		UploadResourceTypeField:   tokens.GetUploadResourceID(),
		UploadResourceIdField:     tokens.GetUploadResource(),
	})
}
