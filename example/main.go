package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-httpkit/httpkit"
	"go.uber.org/zap"
)

// RootHandler is just an example type that implements `http.Handler`.
type RootHandler struct {
	Path         string
	ErrorHandler httpkit.ErrorHandler
}

// ServeHTTP implements the `http.Handler` interface for `RootHandler`. This is where the meat of an
// application would go, including routing if you need to do that in your app.
func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if httpkit.MatchPath(r.URL.Path, h.Path) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusNoContent)
		default:
			h.ErrorHandler(r.Context(), w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		}

		return
	}

	h.ErrorHandler(r.Context(), w, errors.New("path not found"), http.StatusNotFound)
	return
}

func main() {
	// Logger creation is left to the consumer so you can initialize it however you want to.
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("booting example app server")

	// Create an `http.Handler` using the `RootHandler` type as the base handler. Internally, the actual handler
	// used by the server is this base handler wrapped with any interceptors.
	h := httpkit.NewHandler(RootHandler{Path: "/", ErrorHandler: httpkit.DefaultErrorHandler}).UseLogging(logger)

	config := &http.Server{
		Addr: ":9999",

		Handler: h,
	}

	srv := httpkit.NewServer(config)

	if err := srv.Serve(context.Background()); err != nil {
		panic(err)
	}
}
