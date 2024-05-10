package ws

import (
	"context"
	"github.com/topfreegames/pitaya/v2"
	"nunu-layout-chat/pkg/log"
	"time"
)

type Server struct {
	logger *log.Logger
	app    pitaya.Pitaya
}
type Option func(s *Server)

func NewServer(logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithPitayaApp(app pitaya.Pitaya) Option {
	return func(s *Server) {
		s.app = app
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app.Start()
	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.app.Shutdown()

	s.logger.Sugar().Info("Server exiting")
	return nil
}
