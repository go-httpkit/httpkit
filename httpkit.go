package httpkit

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type Handler struct {
	base http.Handler

	chain []Interceptor
}

func NewHandler(h http.Handler) *Handler {
	return &Handler{
		base: h,
	}
}

// ServeHTTP just implements the http.Handler interface by creating the middlware
// chain and applying it to the base handler defined in `serveHTTP`.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	base := h.base

	for i := len(h.chain) - 1; i >= 0; i-- {
		base = h.chain[i].Intercept(base)
	}

	base.ServeHTTP(rw, req)
}

func (h *Handler) AddInterceptor(i Interceptor) *Handler {
	h.chain = append(h.chain, i)

	return h
}

func (h *Handler) UseLogging(logger *zap.Logger) *Handler {
	return h.AddInterceptor(&LoggingMiddleware{
		Base: logger,
	})
}

type Server interface {
	Serve(context.Context) error
}

func NewServer(config *http.Server) Server {
	return &server{
		s: config,
	}
}

type server struct {
	s *http.Server
}

func (srv *server) Serve(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan error)

	go func() {
		defer cancel()

		if err := srv.s.ListenAndServe(); err != http.ErrServerClosed {
			ch <- err
			return
		}

		ch <- nil
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)

		<-c
		if err := srv.s.Shutdown(ctx); err != nil {
			// If an error happens shutting down, cancel the context because it's
			// no longer safe to assume that ListenAndServe running in the previous
			// goroutine will even return. The application is now in an invalid state
			// so it's ok to just abort here.
			ch <- err
			cancel()
		}
	}()

	err := <-ch
	<-ctx.Done()

	return err
}

func MatchPath(requestPath, testPath string) bool {
	return len(requestPath) == len(testPath) && requestPath[:len(testPath)] == testPath
}
