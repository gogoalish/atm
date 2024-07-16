package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gogoalish/atm/config"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server
	notify chan error
	// cfg    *config.Config
}

func New(cfg *config.Config, handler http.Handler) *Server {
	httpServer := &http.Server{
		Handler: handler,
		// ReadTimeout:  _defaultReadTimeout,
		// WriteTimeout: _defaultWriteTimeout,
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}
	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
		// shutdownTimeout: _defaultShutdownTimeout,
	}

	// for _, opt := range opts {
	// 	opt(s)
	// }

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	// Timeout
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
