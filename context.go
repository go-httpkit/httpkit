package httpkit

import (
	"context"

	"go.uber.org/zap"
)

type ContextKey string

var (
	keyRequestID ContextKey = "github.com/go-httpkit/httpkit:request_id"
	keyLogger    ContextKey = "github.com/go-httpkit/httpkit:logger"
)

func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(keyRequestID).(string)

	return requestID
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(keyLogger).(*zap.Logger)

	return logger
}

func SetLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, logger)
}
