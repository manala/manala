package web

import (
	"fmt"
	"github.com/caarlos0/log"
	"github.com/go-chi/chi/v5/middleware"
	internalLog "manala/internal/log"
	"net/http"
	"time"
)

func NewLogger(log *internalLog.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(
		&LogFormatter{
			log: log,
		},
	)
}

type LogFormatter struct {
	log *internalLog.Logger
}

func (formatter *LogFormatter) NewLogEntry(request *http.Request) middleware.LogEntry {
	entry := formatter.log.
		WithoutPadding().
		WithField("request", fmt.Sprintf("%s %s%s %s", request.Method, request.Host, request.RequestURI, request.Proto))

	return &LogEntry{
		entry: entry,
	}
}

type LogEntry struct {
	entry *log.Entry
}

func (entry *LogEntry) Write(status, bytes int, _ http.Header, _ time.Duration, _ interface{}) {
	entry.entry.
		WithField("status", status).
		WithField("bytes", bytes).
		Debug("web server request")
}

func (entry *LogEntry) Panic(v interface{}, stack []byte) {
	entry.entry = entry.entry.
		WithField("stack", stack).
		WithField("panic", fmt.Sprintf("%+v", v))
}
