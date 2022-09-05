package httpkit

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Interceptor interface {
	Intercept(http.Handler) http.Handler
}

type LoggingMiddleware struct {
	Base *zap.Logger
}

func (l *LoggingMiddleware) Intercept(next http.Handler) http.Handler {
	l.Base.Debug("running logger middleware intercept")

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		timeStart := time.Now()

		requestID := uuid.New().String()
		logger := l.Base.With(zap.String("request_id", requestID))
		ctx := SetLogger(SetRequestID(req.Context(), requestID), logger)

		lrw := &loggingResponseWriter{
			ResponseWriter: rw,
			status:         http.StatusOK,
		}

		defer func() {
			logger.With(
				zap.String("http.method", req.Method),
				zap.String("http.path", req.URL.Path),
				zap.Int("http.status", lrw.status),
				zap.Int64("http.response_time_ms", time.Since(timeStart).Milliseconds()),
			).Info("request completed")
		}()

		req = req.Clone(ctx)
		next.ServeHTTP(lrw, req)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
