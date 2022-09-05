package httpkit

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type ErrorHandler func(ctx context.Context, rw http.ResponseWriter, err error, status int)

func DefaultErrorHandler(ctx context.Context, rw http.ResponseWriter, err error, status int) {
	logger := GetLogger(ctx)

	resp := struct {
		Error     string `json:"error"`
		RequestID string `json:"request_id"`
	}{
		Error:     err.Error(),
		RequestID: GetRequestID(ctx),
	}

	buf, err := json.Marshal(resp)
	if err != nil {
		logger.Error(
			"unable to marshal error response json",
			zap.Error(err),
		)
		status = http.StatusInternalServerError
	}

	rw.WriteHeader(status)
	rw.Write(buf)
}
