package log

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

const (
	contextKeyErrLog contextKey = iota
	contextKeyInfoLog
	contextKeyDebugLog
	contextKeySessionDebugLog
)

type (
	contextKey int
	LoggerFunc func(...interface{}) error
	Logger     interface {
		Log(...interface{}) error
	}

	IdentifyingLogger struct {
		Logger

		id string
	}
)

func (i IdentifyingLogger) String() string {
	return i.id
}

func (f LoggerFunc) Log(keyvals ...interface{}) error {
	return f(keyvals...)
}

func With(l Logger, keyvals ...interface{}) Logger {
	return log.With(l, keyvals...)
}

func ContextWithErrLog(ctx context.Context, errLog Logger) context.Context {
	return context.WithValue(ctx, contextKeyErrLog, IdentifyingLogger{
		Logger: errLog,
		id:     "errLog",
	})
}

func ErrLog(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKeyErrLog).(Logger); !ok {
		return log.NewNopLogger()
	} else {
		return logger
	}
}

func ContextWithInfoLog(ctx context.Context, infoLog Logger) context.Context {
	return context.WithValue(ctx, contextKeyInfoLog, IdentifyingLogger{
		Logger: infoLog,
		id:     "infoLog",
	})
}

func InfoLog(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKeyInfoLog).(Logger); !ok {
		return log.NewNopLogger()
	} else {
		return logger
	}
}

func ContextWithDebugLog(ctx context.Context, debugLog Logger) context.Context {
	return context.WithValue(ctx, contextKeyDebugLog, IdentifyingLogger{
		Logger: debugLog,
		id:     "debugLog",
	})
}

func DebugLog(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKeyDebugLog).(Logger); !ok {
		return log.NewNopLogger()
	} else {
		return logger
	}
}

func ContextWithSessionDebugLog(ctx context.Context, debugLog Logger) context.Context {
	return context.WithValue(ctx, contextKeySessionDebugLog, IdentifyingLogger{
		Logger: debugLog,
		id:     "sessionDebugLog",
	})
}

func SessionDebugLog(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKeySessionDebugLog).(Logger); !ok {
		return log.NewNopLogger()
	} else {
		return logger
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode  int
	headersSent bool
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	if l.headersSent {
		return
	}
	l.statusCode = statusCode
	l.headersSent = true
	l.ResponseWriter.WriteHeader(statusCode)
}

func DebugMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			t := time.Now()
			// nolint: errcheck
			DebugLog(ctx).Log("msg", "request start",
				"requestMethod", r.Method,
				"requestURL", r.URL.String(),
				"requestHeaders", fmt.Sprintf("%+v", r.Header),
			)
			w = &loggingResponseWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(w, r)
			// nolint: errcheck
			DebugLog(ctx).Log("msg", "end",
				"statusCode", w.(*loggingResponseWriter).statusCode,
				"responseTime", time.Since(t),
				"responseHeaders", fmt.Sprintf("%+v", w.Header()),
			)
		})
	}
}

type LoggingRoundTripper struct {
	http.RoundTripper

	Logger Logger
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// nolint: errcheck
	l.Logger.Log("msg", "http request",
		"url", req.URL.String(),
	)
	if l.RoundTripper == nil {
		l.RoundTripper = http.DefaultTransport
	}
	resp, err := l.RoundTripper.RoundTrip(req)
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	// nolint: errcheck
	l.Logger.Log("msg", "http response",
		"statusCode", resp.StatusCode,
		"err", err,
		"body", string(b),
	)
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	return resp, err
}
