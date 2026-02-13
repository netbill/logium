package logium

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/netbill/ape"
	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{Entry: logrus.NewEntry(logger.Logger)}
}

func (e *Entry) WithField(key string, value any) *Entry {
	return &Entry{Entry: e.Entry.WithField(key, value)}
}

func (e *Entry) WithFields(fields Fields) *Entry {
	return &Entry{Entry: e.Entry.WithFields(logrus.Fields(fields))}
}

func (e *Entry) WithContext(ctx context.Context) *Entry {
	return &Entry{Entry: e.Entry.WithContext(ctx)}
}

func (e *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Entry: e.Entry.WithTime(t)}
}

func (e *Entry) WithError(err error) *Entry {
	var ae *ape.Error
	if errors.As(err, &ae) {
		return &Entry{Entry: e.Entry.WithError(ae)}
	}
	return &Entry{Entry: e.Entry.WithError(err)}
}

func (e *Entry) WithRequest(r *http.Request) *Entry {
	return e.WithFields(Fields{
		HTTPMethodField: r.Method,
		HTTPPathField:   r.URL.Path,
	})
}

func (e *Entry) WithAccountAuthClaims(auth accountAuthClaims) *Entry {
	return e.WithFields(Fields{
		AccountIDField:        auth.GetAccountID(),
		AccountSessionIDField: auth.GetSessionID(),
	})
}

func (e *Entry) WithUploadContentClaims(tokens uploadContentClaims) *Entry {
	return e.WithFields(Fields{
		UploadAccountIdField:    tokens.GetAccountID(),
		UploadSessionIdField:    tokens.GetSessionID(),
		UploadResourceTypeField: tokens.GetResourceID(),
		UploadResourceIdField:   tokens.GetResource(),
	})
}

func (e *Entry) WithOperation(operation string) *Entry {
	return e.WithField(OperationField, operation)
}

func (e *Entry) WithService(service string) *Entry {
	return e.WithField(ServiceField, service)
}

func (e *Entry) WithComponent(component string) *Entry {
	return e.WithField(ComponentField, component)
}
