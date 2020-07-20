package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

const probeUserAgent = "kube-probe"

// NewRequestLogger returns request logger.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{Logger: logrus.StandardLogger()}
}

// RequestLogger implements middleware.LogFormatter interface and is used for http/https requests and responses logging
type RequestLogger struct {
	*logrus.Logger
}

// NewLogEntry creates new logrus entry as well as logging request info
func (rl RequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	if strings.Contains(r.UserAgent(), probeUserAgent) {
		return &skipLogEntry{}
	}

	entry := requestLoggerEntry{Logger: rl}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rl.Error(err)
	}
	if err := r.Body.Close(); err != nil {
		rl.Error(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	logFields := logrus.Fields{
		"ts":          time.Now().UTC().Format(time.RFC1123),
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"uri":         r.RequestURI,
		"payload":     string(bodyBytes),
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	entry.Logger = entry.Logger.WithFields(logFields)

	return &entry
}

type requestLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (rle *requestLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	rle.Logger.WithFields(logrus.Fields{
		"status":     status,
		"length":     bytes,
		"elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	}).Debugln("request completed")
}

func (rle *requestLoggerEntry) Panic(v interface{}, stack []byte) {
	rle.Logger = rle.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

type skipLogEntry struct{}

func (sle *skipLogEntry) Write(_, _ int, _ time.Duration) {}

func (sle *skipLogEntry) Panic(_ interface{}, _ []byte) {}
