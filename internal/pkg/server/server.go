package server

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	shutdownTimeout   = time.Second * 5
	HeaderContentType = "Content-Type"
	HeaderSourceType  = "Source-Type"
	JsonContentType   = "application/json"
)

// Server is an http.Server wrapper.
type Server struct {
	srv *http.Server
}

// New returns a new instance of Server ready to handle requests
// using given handler.
func New(cfg *Config, handler http.Handler) Server {
	srv := http.Server{
		Addr:    cfg.addr(),
		Handler: handler,
	}

	return Server{
		srv: &srv,
	}
}

// Run runs a server. After context cancelling server will be gracefully
// shutdowned. If ListenAndServe returns http.ErrServerClosed, Run returns
// nil error.
func (s *Server) Run(ctx context.Context) error {
	logrus.Infof("http server started on %s", s.srv.Addr)

	go s.shutdownOnCancel(ctx)

	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) shutdownOnCancel(ctx context.Context) {
	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		logrus.Errorf("failed to shutdown server in %s", shutdownTimeout.String())
	}
}
